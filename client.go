package mysubaru

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alex-savin/go-mysubaru/v2/config"
	"resty.dev/v3"
)

// ErrDeviceNotRegistered indicates the device must complete 2FA/registration
// before authentication can succeed. Callers can detect it with errors.Is.
var ErrDeviceNotRegistered = errors.New("device is not registered")

// apiVersionRetryLimit caps how many times the client auto-increments the API
// version in response to HTTP 404s before giving up (mirrors subarulink).
const apiVersionRetryLimit = 5

// requestTimeout bounds each individual HTTP attempt (the caller's context can
// impose a stricter overall deadline across retries).
const requestTimeout = 30 * time.Second

// sessionValidityWindow is how long a session is trusted after the last
// successful API response before validateSession makes a real round-trip.
// The backend expires idle sessions after 5 minutes; 4 leaves a safety margin.
const sessionValidityWindow = 4 * time.Minute

var (
	// apiVersionPrefixRe matches the leading /g2vNN version segment of a path.
	apiVersionPrefixRe = regexp.MustCompile(`^/g2v\d+`)
	// apiVersionNumRe matches the trailing version number of a version prefix.
	apiVersionNumRe = regexp.MustCompile(`\d+$`)
)

// Client represents a MySubaru API client that interacts with the MySubaru API.
type Client struct {
	credentials config.Credentials
	httpClient  *resty.Client
	country     string // USA | CA
	// stateMu guards session state mutated by auth/re-auth on background
	// goroutines while pollers/commands read it: contactMethods, currentVin,
	// listOfVins. It is distinct from reqMu (the request-serialization mutex),
	// and is never acquired while that mutex is held (and vice-versa).
	stateMu        sync.RWMutex
	contactMethods dataMap // List of contact methods for 2FA
	currentVin     string
	listOfVins     []string
	// Liveness/auth flags are written from both inside and outside the request
	// lock, so they are atomic to stay race-free without lock-ordering concerns.
	isAuthenticated atomic.Bool
	isRegistered    atomic.Bool
	isAlive         atomic.Bool
	// apiVer holds the current API version prefix (e.g. "/g2v33"), auto-bumped on
	// 404. apiBumps counts bumps against apiVersionRetryLimit. Both atomic so the
	// request path reads them without taking a lock.
	apiVer         atomic.Pointer[string]
	apiBumps       atomic.Int32
	updateInterval int // seconds, DEFAULT_UPDATE_INTERVAL
	fetchInterval  int // seconds, DEFAULT_FETCH_INTERVAL
	logger         *slog.Logger
	metrics        config.MetricsRecorder
	// baseURL is the resolved API host: the config override when set, otherwise
	// the regional default from mobileAPIServer.
	baseURL string
	// lastValidated holds the unix time of the last proof the session is alive
	// (any successful API response). validateSession uses it to skip redundant
	// validate+select round-trips within sessionValidityWindow.
	lastValidated atomic.Int64
	// reqMu serializes all HTTP requests. The MySubaru backend is a stateful,
	// cookie-scoped session (the selected vehicle is server-side session state),
	// so requests are deliberately one-at-a-time. It also guards httpClient,
	// which resetSession swaps for a fresh cookie jar.
	reqMu sync.RWMutex
}

// session-state accessors (guarded by stateMu) ------------------------------

func (c *Client) getCurrentVin() string {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.currentVin
}

func (c *Client) setCurrentVin(vin string) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.currentVin = vin
}

func (c *Client) getVins() []string {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return slices.Clone(c.listOfVins)
}

func (c *Client) hasVin(vin string) bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return slices.Contains(c.listOfVins, vin)
}

// setVins replaces the VIN list and returns the first VIN (or "" if empty).
func (c *Client) setVins(vins []string) string {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.listOfVins = vins
	if len(vins) > 0 {
		c.currentVin = vins[0]
		return vins[0]
	}
	return ""
}

func (c *Client) getContactMethodsData() dataMap {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.contactMethods
}

func (c *Client) setContactMethodsData(dm dataMap) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.contactMethods = dm
}

// httpC returns the current resty client under the request lock. resetSession
// can swap the pointer, so callers that use it outside executeOnce must capture
// it through this accessor to avoid racing that swap.
func (c *Client) httpC() *resty.Client {
	c.reqMu.RLock()
	defer c.reqMu.RUnlock()
	return c.httpClient
}

// getAPIVersion returns the current API version prefix, falling back to the
// package default when unset (e.g. a Client built without New, as in tests).
func (c *Client) getAPIVersion() string {
	if p := c.apiVer.Load(); p != nil && *p != "" {
		return *p
	}
	return MOBILE_API_VERSION
}

// applyAPIVersion rewrites a path's leading /g2vNN segment to the client's
// current version, so call sites can keep using the package default while the
// client transparently follows version bumps. Non-versioned paths pass through.
func (c *Client) applyAPIVersion(url string) string {
	if apiVersionPrefixRe.MatchString(url) {
		return apiVersionPrefixRe.ReplaceAllString(url, c.getAPIVersion())
	}
	return url
}

// bumpAPIVersion increments the trailing number of the current API version
// (e.g. /g2v33 -> /g2v34) after a 404, persisting it for subsequent requests.
// Returns false once apiVersionRetryLimit is reached so the caller stops.
func (c *Client) bumpAPIVersion() bool {
	if c.apiBumps.Load() >= apiVersionRetryLimit {
		return false
	}
	cur := c.getAPIVersion()
	m := apiVersionNumRe.FindString(cur)
	if m == "" {
		return false
	}
	n, err := strconv.Atoi(m)
	if err != nil {
		return false
	}
	next := cur[:len(cur)-len(m)] + strconv.Itoa(n+1)
	c.apiVer.Store(&next)
	c.apiBumps.Add(1)
	return true
}

// newHTTPClient builds a resty client with the client's base URL and the
// mobile-app headers for its country. Used both at construction and when
// resetSession needs a fresh cookie jar.
func (c *Client) newHTTPClient() *resty.Client {
	httpClient := resty.New()
	httpClient.
		SetBaseURL(c.baseURL).
		SetHeaders(map[string]string{
			"User-Agent":       "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Origin":           "file://",
			"X-Requested-With": mobileApp[c.country],
			"Accept-Language":  "en-US,en;q=0.9",
			"Accept-Encoding":  "gzip, deflate",
			"Accept":           "*/*"},
		)
	return httpClient
}

// New function creates a New MySubaru API client
func New(config *config.Config) (*Client, error) {
	metrics := config.Metrics
	if metrics == nil {
		metrics = &NoOpMetricsRecorder{}
	}

	client := &Client{
		credentials:    config.MySubaru.Credentials,
		country:        config.MySubaru.Region,
		updateInterval: DEFAULT_UPDATE_INTERVAL,
		fetchInterval:  DEFAULT_FETCH_INTERVAL,
		logger:         config.Logger,
		metrics:        metrics,
	}
	client.baseURL = config.MySubaru.BaseURL
	if client.baseURL == "" {
		client.baseURL = mobileAPIServer[client.country]
	}
	initialVersion := MOBILE_API_VERSION
	client.apiVer.Store(&initialVersion)

	client.httpClient = client.newHTTPClient()

	// Don't authenticate during initialization - let Authenticate() method handle it
	return client, nil
}

// auth authenticates the client with the MySubaru API using the provided credentials.
func (c *Client) auth(ctx context.Context) (bool, error) {
	params := map[string]string{
		"env":           "cloudprod",
		"deviceType":    "android",
		"loginUsername": c.credentials.Username,
		"password":      c.credentials.Password,
		"deviceId":      c.credentials.DeviceID,
		"passwordToken": "",
		"selectedVin":   "",
		"pushToken":     ""}
	reqURL := MOBILE_API_VERSION + apiURLs["API_LOGIN"]
	resp, err := c.execute(ctx, POST, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing auth request", "request", "auth", "error", err.Error())
		return false, fmt.Errorf("error while executing auth request: %w", err)
	}
	c.logger.Debug("http request output", "request", "auth", "body", resp)

	var sd SessionData
	err = json.Unmarshal(resp.Data, &sd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "auth", "error", err.Error())
		return false, fmt.Errorf("failed to parse auth response: %w", err)
	}

	if !sd.DeviceRegistered {
		err := c.getContactMethods(ctx)
		if err != nil {
			c.logger.Error("error while getting contact methods", "request", "auth", "error", err.Error())
			return false, errors.New("error while getting contact methods: " + err.Error())
		}

		c.logger.Error("device is not registered", "request", "auth", "deviceId", c.credentials.DeviceID)
		return false, fmt.Errorf("%w: %s", ErrDeviceNotRegistered, c.credentials.DeviceID)
	}

	if sd.DeviceRegistered && sd.RegisteredDevicePermanent {
		c.isAuthenticated.Store(true)
		c.isRegistered.Store(true)
		c.isAlive.Store(true)
	}
	c.logger.Debug("MySubaru API client authenticated")

	if len(sd.Vehicles) > 0 {
		vins := make([]string, 0, len(sd.Vehicles))
		for _, vehicle := range sd.Vehicles {
			vins = append(vins, vehicle.Vin)
		}
		c.setVins(vins) // replaces (not appends) so re-auth doesn't duplicate VINs
	} else {
		errNoVehicles := errors.New("there are no vehicles associated with the account")
		c.logger.Error("there are no vehicles associated with the account", "request", "auth", "error", errNoVehicles.Error())
		return false, errNoVehicles
	}
	return true, nil
}

// resetSession clears the current session by removing cookies and resetting session state.
// This is useful when the API returns errors like VEHICLESETUPERROR that indicate
// a stale or corrupted session state on the Subaru backend.
func (c *Client) resetSession() {
	c.logger.Warn("resetting session - clearing cookies and session state")
	// Create a new HTTP client to clear all cookies
	c.httpClient = c.newHTTPClient()
	c.isAlive.Store(false)
	c.lastValidated.Store(0)
}

// SelectVehicle selects a vehicle by its VIN. If no VIN is provided, it uses the current VIN.
func (c *Client) SelectVehicle(ctx context.Context, vin string) (*VehicleData, error) {
	if vin == "" {
		vin = c.getCurrentVin()
	}
	if err := ValidateVIN(vin); err != nil {
		return nil, err
	}

	params := map[string]string{
		"vin": vin,
		"_":   timestamp()}
	reqURL := MOBILE_API_VERSION + apiURLs["API_SELECT_VEHICLE"]
	resp, err := c.execute(ctx, GET, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing SelectVehicle request", "request", "SelectVehicle", "error", err.Error())
		return nil, fmt.Errorf("error while executing SelectVehicle request: %w", err)
	}

	var vd VehicleData
	err = json.Unmarshal(resp.Data, &vd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "SelectVehicle", "error", err.Error())
		return nil, fmt.Errorf("error while parsing vehicle selection response: %w", err)
	}
	return &vd, nil
}

// GetVehicles retrieves a list of vehicles associated with the client's account.
func (c *Client) GetVehicles(ctx context.Context) ([]*Vehicle, error) {
	var vehicles []*Vehicle
	for _, vin := range c.getVins() {
		vehicle, err := c.GetVehicleByVin(ctx, vin)
		if err != nil {
			c.logger.Error("cannot get vehicle data", "request", "GetVehicles", "error", err.Error())
			return nil, fmt.Errorf("cannot get vehicle data: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, nil
}

// GetVehicleByVin retrieves a vehicle by its VIN from the client's list of vehicles.
func (c *Client) GetVehicleByVin(ctx context.Context, vin string) (*Vehicle, error) {
	var vehicle *Vehicle
	if c.hasVin(vin) {
		params := map[string]string{
			"vin": vin,
			"_":   timestamp()}
		reqURL := MOBILE_API_VERSION + apiURLs["API_SELECT_VEHICLE"]
		resp, err := c.execute(ctx, GET, reqURL, params, false)
		if err != nil {
			c.logger.Error("error while executing GetVehicleByVin request", "request", "GetVehicleByVin", "error", err.Error())
			return nil, fmt.Errorf("error while executing GetVehicleByVin request: %w", err)
		}

		var vd VehicleData
		err = json.Unmarshal(resp.Data, &vd)
		if err != nil {
			c.logger.Error("error while parsing json", "request", "GetVehicleByVin", "error", err.Error())
			return nil, fmt.Errorf("failed to parse GetVehicleByVin response: %w", err)
		}
		// c.logger.Debug("http request output", "request", "GetVehicleByVin", "body", resp)

		vehicle = &Vehicle{
			Vin:                  vin,
			CarName:              vd.VehicleName,
			CarNickname:          vd.Nickname,
			ModelName:            vd.ModelName,
			ModelYear:            vd.ModelYear,
			ModelCode:            vd.ModelCode,
			ExtDescrip:           vd.ExtDescrip,
			IntDescrip:           vd.IntDescrip,
			TransCode:            vd.TransCode,
			EngineSize:           vd.EngineSize,
			VehicleKey:           vd.VehicleKey,
			LicensePlate:         vd.LicensePlate,
			LicensePlateState:    vd.LicensePlateState,
			Features:             vd.Features,
			SubscriptionFeatures: vd.SubscriptionFeatures,
			client:               c,
		}
		vehicle.Doors = make(map[string]Door)
		vehicle.Windows = make(map[string]Window)
		vehicle.Tires = make(map[string]Tire)
		vehicle.ClimateProfiles = make(map[string]ClimateProfile)
		vehicle.Troubles = make(map[string]Trouble)

		// Populate vehicle state - log errors but don't fail the entire vehicle creation
		if err := vehicle.GetVehicleStatus(ctx); err != nil {
			c.logger.Warn("failed to get vehicle status during vehicle initialization", "vin", vin, "error", err.Error())
		}
		if err := vehicle.GetVehicleCondition(ctx); err != nil {
			c.logger.Warn("failed to get vehicle condition during vehicle initialization", "vin", vin, "error", err.Error())
		}
		if err := vehicle.GetVehicleHealth(ctx); err != nil {
			c.logger.Warn("failed to get vehicle health during vehicle initialization", "vin", vin, "error", err.Error())
		}

		// Get climate presets - log errors but don't fail the entire vehicle creation
		if err := vehicle.GetClimatePresets(ctx); err != nil {
			c.logger.Warn("failed to get climate presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		if err := vehicle.GetClimateUserPresets(ctx); err != nil {
			c.logger.Warn("failed to get climate user presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		if err := vehicle.GetClimateQuickPresets(ctx); err != nil {
			c.logger.Warn("failed to get climate quick presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		return vehicle, nil
	}
	c.logger.Error("vin code is not in the list of the available vin codes", "request", "GetVehicleByVIN")
	return nil, errors.New("vin code is not in the list of the available vin codes")
}

func (c *Client) RefreshVehicles(ctx context.Context) error {
	params := map[string]string{}
	reqURL := MOBILE_API_VERSION + apiURLs["API_REFRESH_VEHICLES"]
	resp, err := c.execute(ctx, GET, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing RefreshVehicles request", "request", "RefreshVehicles", "error", err.Error())
		return fmt.Errorf("RefreshVehicles request failed: %w", err)
	}
	c.logger.Debug("http request output", "request", "RefreshVehicles", "body", resp)

	var sd SessionData
	err = json.Unmarshal(resp.Data, &sd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "RefreshVehicles", "error", err.Error())
		return fmt.Errorf("failed to parse RefreshVehicles response: %w", err)
	}

	// Update client state
	if sd.DeviceRegistered && sd.RegisteredDevicePermanent {
		c.isAuthenticated.Store(true)
		c.isRegistered.Store(true)
		c.isAlive.Store(true)
	}
	vins := make([]string, 0, len(sd.Vehicles))
	for _, vehicle := range sd.Vehicles {
		vins = append(vins, vehicle.Vin)
	}
	c.setVins(vins)

	c.logger.Info("vehicles refreshed successfully", "count", len(vins))
	return nil
}

// RequestAuthCode requests an authentication code for two-factor authentication (2FA).
// (?!^).(?=.*@)
// (?!^): This is a negative lookbehind assertion. It ensures that the matched character is not at the beginning of the string.
// .: This matches any single character (except newline, by default).
// (?=.*@): This is a positive lookahead assertion. It ensures that the matched character is followed by any characters (.*) and then an "@" symbol. This targets the username part of the email address.
func (c *Client) RequestAuthCode(ctx context.Context, email string) error {
	email, err := emailMasking(email)
	if err != nil {
		c.logger.Error("error while hiding email", "request", "RequestAuthCode", "error", err.Error())
		return fmt.Errorf("error while hiding email: %w", err)
	}

	if !containsValueInStruct(c.getContactMethodsData(), email) {
		c.logger.Error("email is not in the list of contact methods", "request", "RequestAuthCode", "email", email)
		return errors.New("email is not in the list of contact methods: " + email)
	}

	params := map[string]string{
		"contactMethod":      email,
		"languagePreference": "EN"}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_2FA_SEND_VERIFICATION"]
	resp, err := c.execute(ctx, POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing RequestAuthCode request", "request", "RequestAuthCode", "error", err.Error())
		return fmt.Errorf("error while executing RequestAuthCode request: %w", err)
	}
	c.logger.Debug("http request output", "request", "RequestAuthCode", "body", resp)

	return nil
}

// SubmitAuthCode submits the authentication code received from the RequestAuthCode method.
func (c *Client) SubmitAuthCode(ctx context.Context, code string, permanent bool) error {
	regex := regexp.MustCompile(`^\d{6}$`)
	if !regex.MatchString(code) {
		c.logger.Error("invalid verification code format", "request", "SubmitAuthCode", "code", code)
		return errors.New("invalid verification code format, must be 6 digits")
	}

	params := map[string]string{
		"deviceId":         c.credentials.DeviceID,
		"deviceName":       c.credentials.DeviceName,
		"verificationCode": code}
	if permanent {
		params["rememberDevice"] = "on"
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_2FA_AUTH_VERIFY"]
	resp, err := c.execute(ctx, POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing SubmitAuthCode request", "request", "SubmitAuthCode", "error", err.Error())
		return fmt.Errorf("error while executing SubmitAuthCode request: %w", err)
	}
	c.logger.Debug("http request output", "request", "SubmitAuthCode", "body", resp)

	// Device registration does not always immediately take effect
	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	// Reauthenticate after submitting the code
	if ok, err := c.auth(ctx); !ok {
		c.logger.Error("error while executing auth request", "request", "auth", "error", err.Error())
		return fmt.Errorf("error while executing auth request: %w", err)
	}

	return nil
}

// getContactMethods retrieves the available contact methods for two-factor authentication (2FA).
// {"success":true,"dataName":"dataMap","data":{"userName":"a**x@savin.nyc","email":"t***a@savin.nyc"}}
func (c *Client) getContactMethods(ctx context.Context) error {
	params := map[string]string{}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_2FA_CONTACT"]
	resp, err := c.execute(ctx, POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing getContactMethods request", "request", "getContactMethods", "error", err.Error())
		return fmt.Errorf("error while executing getContactMethods request: %w", err)
	}
	c.logger.Debug("http request output", "request", "getContactMethods", "body", resp)

	var dm dataMap
	err = json.Unmarshal(resp.Data, &dm)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "getContactMethods", "error", err.Error())
		return fmt.Errorf("error while parsing contact methods response: %w", err)
	}
	c.setContactMethodsData(dm)
	c.logger.Debug("contact methods successfully retrieved", "request", "getContactMethods", "methods", dm)

	return nil
}

// RemoteUnlock unlocks the vehicle remotely.
func (c *Client) RemoteUnlock(ctx context.Context, vin string) error {
	if !c.hasVin(vin) {
		return errors.New("VIN not in list")
	}

	if !c.validateSession(ctx) {
		return errors.New("session is not valid")
	}

	vData, err := c.SelectVehicle(ctx, vin)
	if err != nil {
		return fmt.Errorf("failed to select vehicle: %w", err)
	}

	if !slices.Contains(vData.SubscriptionFeatures, FEATURE_REMOTE) {
		return errors.New(appErrors["SUBSCRIPTION_REQUIRED"])
	}
	reqURL := MOBILE_API_VERSION + apiURLs["API_UNLOCK"]
	params := map[string]string{"vin": vin}
	resp, err := c.execute(ctx, POST, reqURL, params, false)
	if err != nil {
		return fmt.Errorf("RemoteUnlock request failed: %w", err)
	}
	if !resp.Success {
		return errors.New("RemoteUnlock failed")
	}
	c.logger.Info("Vehicle unlocked successfully", "vin", vin)
	return nil
}

// execute executes an HTTP request based on the method, URL, and parameters provided.
func (c *Client) execute(ctx context.Context, method string, url string, params map[string]string, j bool) (*Response, error) {
	return c.executeWithRetry(ctx, method, url, params, j, 3)
}

// executeWithRetry executes an HTTP request with retry logic
func (c *Client) executeWithRetry(ctx context.Context, method string, url string, params map[string]string, j bool, maxRetries int) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Record retry attempt
			c.metrics.RecordRetry(url, attempt)

			// Exponential backoff: wait 1s, 2s, 4s, etc.
			waitTime := time.Duration(1<<uint(attempt-1)) * time.Second
			if waitTime > 30*time.Second {
				waitTime = 30 * time.Second
			}
			c.logger.Debug("retrying request after backoff", "attempt", attempt, "wait", waitTime)
			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := c.executeOnce(ctx, method, url, params, j)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry when the caller's context is gone or on non-retryable errors
		if ctx.Err() != nil || !IsRetryableError(err) {
			break
		}

		// Handle InvalidToken by re-authenticating before retry
		var apiErr APIError
		if errors.As(err, &apiErr) && apiErr.Code == apiErrors["API_ERROR_INVALID_TOKEN"] {
			c.logger.Info("InvalidToken error detected, attempting re-authentication before retry")
			if !c.reauthenticateAndSelect(ctx) {
				c.logger.Error("re-authentication failed during retry")
				break
			}
		}

		c.logger.Debug("request failed, will retry", "attempt", attempt+1, "maxRetries", maxRetries, "error", err.Error())
	}

	return nil, lastErr
}

// executeOnce performs a single HTTP request attempt
// sendRequest executes the HTTP request based on method type.
func (c *Client) sendRequest(req *resty.Request, method, url string, params map[string]string, j bool) (*resty.Response, error) {
	switch method {
	case GET:
		req.SetQueryParams(params)
		return req.Get(url)
	case POST:
		if j {
			req.SetBody(params)
		} else {
			req.SetFormData(params)
		}
		return req.Post(url)
	default:
		return nil, errors.New("unsupported HTTP method: " + method)
	}
}

// handleAPIError logs error response details if available.
func (c *Client) handleAPIError(r *Response) {
	if r.DataName != "errorResponse" {
		return
	}
	var er ErrorResponse
	if jsonErr := json.Unmarshal(r.Data, &er); jsonErr != nil {
		c.logger.Error("error parsing error response", "error", jsonErr.Error())
		return
	}
	// apiErrors maps symbolic names to wire labels, so the incoming label must be
	// matched against the map values (via the wireToSymbol reverse index). Known
	// errors are routinely handled/retried by the caller (e.g. InvalidToken
	// triggers re-auth), so keep them at debug; the retry layer logs a real error
	// if recovery ultimately fails.
	if isKnownAPIErrorLabel(er.ErrorLabel) {
		c.logger.Debug("known API error", "label", er.ErrorLabel, "description", er.ErrorDescription)
	} else {
		c.logger.Warn("unknown API error", "label", er.ErrorLabel, "description", er.ErrorDescription)
	}
}

// handleVehicleSetupError handles the VEHICLESETUPERROR case and returns response/error accordingly.
func (c *Client) handleVehicleSetupError(r *Response, method, url string, duration time.Duration) (*Response, error) {
	// With vehicle data: treat as success (user needs to complete setup but data is functional)
	if r.DataName == "vehicle" {
		c.logger.Debug("VEHICLESETUPERROR received but vehicle data is present; treating as success",
			"errorCode", r.ErrorCode, "dataName", r.DataName)
		c.metrics.RecordRequest(method, url, duration, true)
		c.isAlive.Store(true)
		return r, nil
	}
	// Without vehicle data: reset session and allow retry
	c.logger.Warn("VEHICLESETUPERROR received without vehicle data; resetting session",
		"errorCode", r.ErrorCode, "dataName", r.DataName)
	c.resetSession()
	c.metrics.RecordRequest(method, url, duration, false)
	return nil, APIError{Code: r.ErrorCode, Message: "VEHICLESETUPERROR: session reset, please retry", Retryable: true}
}

func (c *Client) executeOnce(ctx context.Context, method string, url string, params map[string]string, j bool) (*Response, error) {
	start := time.Now()

	c.reqMu.Lock()
	defer c.reqMu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	// Subaru retires old API versions with a 404. Send against the client's
	// current version, and on a 404 bump the version and retry in-place so the
	// client follows version transitions without code changes.
	var resp *resty.Response
	var err error
	for {
		versionedURL := c.applyAPIVersion(url)
		req := c.httpClient.R().SetContext(ctx)
		resp, err = c.sendRequest(req, method, versionedURL, params, j)
		if err == nil && resp.StatusCode() == 404 {
			prev := c.getAPIVersion()
			if c.bumpAPIVersion() {
				_ = resp.Body.Close()
				c.logger.Warn("API version returned 404; bumping and retrying",
					"url", versionedURL, "from", prev, "to", c.getAPIVersion())
				continue
			}
		}
		break
	}
	if err != nil {
		c.metrics.RecordRequest(method, url, time.Since(start), false)
		c.metrics.RecordError("network_error")
		c.logger.Error("error while executing HTTP request", "method", method, "url", url, "error", err.Error())
		return nil, ErrNetworkError
	}

	resBytes, err := io.ReadAll(resp.Body)
	duration := time.Since(start)
	if err != nil {
		c.metrics.RecordRequest(method, url, duration, false)
		c.metrics.RecordError("response_read_error")
		c.logger.Error("error while reading response body", "error", err.Error())
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	c.logger.Debug("received HTTP response", "method", method, "url", url, "status", resp.Status(), "body", string(resBytes))

	c.httpClient.SetCookies(resp.Cookies())

	// The backend serves HTML error pages (login redirects, maintenance pages)
	// instead of JSON when the session is invalid or the API has changed; catch
	// that here once so individual call sites don't have to.
	if isHTMLResponse(resBytes) {
		c.metrics.RecordRequest(method, url, duration, false)
		c.metrics.RecordError("html_response")
		c.isAlive.Store(false)
		c.lastValidated.Store(0)
		c.logger.Error("received HTML error page instead of JSON", "method", method, "url", url, "response_start", getResponsePreview(resBytes, 200))
		return nil, errHTMLResponse(url)
	}

	r, ok := c.parseResponse(resBytes)
	if !ok {
		c.metrics.RecordRequest(method, url, duration, false)
		c.metrics.RecordError("parse_error")
		c.isAlive.Store(false)
		return nil, errors.New("failed to parse response")
	}

	if resp.IsStatusSuccess() && r.Success {
		c.metrics.RecordRequest(method, url, duration, true)
		c.isAlive.Store(true)
		// Any successful API response resets the backend's idle-session timer, so
		// it doubles as proof of session validity (see validateSession).
		c.lastValidated.Store(time.Now().Unix())
		return &r, nil
	}

	c.handleAPIError(&r)

	if resp.IsStatusSuccess() && !r.Success {
		// handleAPIError (above) already classifies and logs the specific error;
		// this is redundant diagnostic detail, so keep it at debug. A genuine,
		// unrecovered failure still surfaces via the retry-exhaustion error path.
		c.logger.Debug("API returned success=false despite HTTP 200",
			"method", method, "url", url, "errorCode", r.ErrorCode, "dataName", r.DataName)

		if r.ErrorCode == apiErrors["API_ERROR_VEHICLE_SETUP"] {
			return c.handleVehicleSetupError(&r, method, url, duration)
		}

		c.metrics.RecordRequest(method, url, duration, false)
		c.isAlive.Store(false)
		if r.ErrorCode != "" {
			// Map the wire code to a typed error (NegativeAckError, PINLockedError,
			// retryable APIError, etc.) so callers can use errors.As/Is. The retry
			// layer keys off APIError.Code/Retryable, both preserved by ParseAPIError.
			parsedErr := ParseAPIError(r.ErrorCode)
			if IsSessionError(parsedErr) {
				// The cached session-validity window no longer holds.
				c.lastValidated.Store(0)
			}
			return nil, parsedErr
		}
		return nil, APIError{Code: "API_SUCCESS_FALSE", Message: "API request failed with success=false", Retryable: true}
	}

	c.metrics.RecordRequest(method, url, duration, false)
	c.isAlive.Store(false)
	return nil, fmt.Errorf("request failed with status %s", resp.Status())
}

// parseResponse parses the JSON response from the MySubaru API into a Response struct.
func (c *Client) parseResponse(b []byte) (Response, bool) {
	var r Response
	err := json.Unmarshal(b, &r)
	if err != nil {
		c.logger.Error("error while parsing json", "error", err.Error())
		return r, false
	}
	return r, true
}

// validateSession checks that the current session is still valid, trusting the
// recent-activity window first: every successful API response proves the
// session alive for sessionValidityWindow, so within it no round-trip is made.
// Past the window it calls validateSession.json, re-selects the current VIN,
// and falls back to full re-authentication when either step fails.
func (c *Client) validateSession(ctx context.Context) bool {
	if c == nil {
		return false
	}
	if last := c.lastValidated.Load(); last > 0 && time.Since(time.Unix(last, 0)) < sessionValidityWindow {
		return true
	}
	reqURL := MOBILE_API_VERSION + apiURLs["API_VALIDATE_SESSION"]
	resp, err := c.execute(ctx, GET, reqURL, map[string]string{}, false)
	if err != nil {
		c.logger.Error("error while executing validateSession request", "request", "validateSession", "error", err.Error())
		return c.reauthenticateAndSelect(ctx)
	}
	if resp == nil {
		c.logger.Warn("validateSession returned nil response; forcing re-auth")
		return c.reauthenticateAndSelect(ctx)
	}
	c.logger.Debug("http request output", "request", "validateSession", "body", resp)

	if resp.Success {
		if _, err := c.SelectVehicle(ctx, c.getCurrentVin()); err == nil {
			return true
		} else {
			errMsg := err.Error()
			c.logger.Warn("select vehicle failed during session validation; attempting re-auth", "request", "validateSession", "error", errMsg)
			return c.reauthenticateAndSelect(ctx)
		}
	}

	return c.reauthenticateAndSelect(ctx)
}

// reauthenticateAndSelect forces re-authentication and re-selects the current VIN.
// Returns true on success, false if either step fails.
func (c *Client) reauthenticateAndSelect(ctx context.Context) bool {
	if ok, err := c.auth(ctx); !ok || err != nil {
		if err != nil {
			c.logger.Error("error while re-authenticating", "request", "validateSession", "error", err.Error())
		} else {
			c.logger.Error("reauthentication failed", "request", "validateSession")
		}
		return false
	}

	// Ensure we have a VIN to select. Fall back to the first known VIN after auth.
	if c.getCurrentVin() == "" {
		if vins := c.getVins(); len(vins) > 0 {
			c.setCurrentVin(vins[0])
		}
	}

	if _, err := c.SelectVehicle(ctx, c.getCurrentVin()); err != nil {
		c.logger.Error("error while selecting vehicle", "request", "validateSession", "error", err.Error())
		return false
	}

	return true
}

// Authenticate attempts to authenticate the client. needs2FA reports that the
// device must complete 2FA/device registration (via RequestAuthCode and
// SubmitAuthCode) before authentication can succeed; transport/parse errors are
// NOT treated as 2FA-required — they are genuine failures.
func (c *Client) Authenticate(ctx context.Context) (ok bool, needs2FA bool, err error) {
	ok, err = c.auth(ctx)

	if !ok && errors.Is(err, ErrDeviceNotRegistered) {
		return false, true, err // Device registration required
	}

	return ok, false, err
}

// GetAppStatus queries the MySubaru availability/maintenance gate
// (appStatus.json). The mobile app checks it before login; a success=false
// reply means the backend is in a maintenance window, which is reported as
// (false, nil) — not an error — so callers can distinguish "API down for
// maintenance" from a transport failure. No authentication is required.
func (c *Client) GetAppStatus(ctx context.Context) (bool, error) {
	reqURL := MOBILE_API_VERSION + apiURLs["API_APP_STATUS"]
	// Single attempt: a "down for maintenance" reply must not be ground through
	// the retry/backoff loop before being reported.
	resp, err := c.executeWithRetry(ctx, GET, reqURL, map[string]string{"_": timestamp()}, false, 0)
	if err != nil {
		var apiErr APIError
		if errors.As(err, &apiErr) && apiErr.Code != ErrNetworkError.Code {
			// The backend answered with success=false: unavailable, not broken.
			c.logger.Warn("MySubaru API reports unavailable", "request", "GetAppStatus", "code", apiErr.Code)
			return false, nil
		}
		c.logger.Error("error while executing GetAppStatus request", "request", "GetAppStatus", "error", err.Error())
		return false, fmt.Errorf("GetAppStatus request failed: %w", err)
	}
	return resp.Success, nil
}

// Logout terminates the session on the MySubaru backend (invalidateSession.json)
// and clears the client's local auth state. Device registration is permanent and
// survives logout. A session that is already expired or invalid on the backend
// is treated as a successful logout.
func (c *Client) Logout(ctx context.Context) error {
	reqURL := MOBILE_API_VERSION + apiURLs["API_INVALIDATE_SESSION"]
	// Single attempt: retrying (or re-authenticating) while tearing the session
	// down would defeat the purpose.
	_, err := c.executeWithRetry(ctx, GET, reqURL, map[string]string{"_": timestamp()}, false, 0)

	c.isAuthenticated.Store(false)
	c.isAlive.Store(false)
	c.lastValidated.Store(0)

	if err != nil && !IsSessionError(err) {
		c.logger.Warn("error while invalidating session", "request", "Logout", "error", err.Error())
		return fmt.Errorf("Logout request failed: %w", err)
	}
	c.logger.Debug("session invalidated", "request", "Logout")
	return nil
}
