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
	"strings"
	"sync"
	"time"

	"github.com/alex-savin/go-mysubaru/config"
	pkgerrors "github.com/pkg/errors"
	"resty.dev/v3"
)

// Client represents a MySubaru API client that interacts with the MySubaru API.
type Client struct {
	credentials     config.Credentials
	httpClient      *resty.Client
	country         string  // USA | CA
	contactMethods  dataMap // List of contact methods for 2FA
	currentVin      string
	listOfVins      []string
	isAuthenticated bool
	isRegistered    bool
	isAlive         bool
	updateInterval  int // 7200
	fetchInterval   int // 360
	logger          *slog.Logger
	metrics         config.MetricsRecorder
	sync.RWMutex
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
		updateInterval: 7200,
		fetchInterval:  360,
		logger:         config.Logger,
		metrics:        metrics,
	}

	httpClient := resty.New()
	httpClient.
		SetBaseURL(MOBILE_API_SERVER[client.country]).
		SetHeaders(map[string]string{
			"User-Agent":       "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Origin":           "file://",
			"X-Requested-With": MOBILE_APP[client.country],
			"Accept-Language":  "en-US,en;q=0.9",
			"Accept-Encoding":  "gzip, deflate",
			"Accept":           "*/*"},
		)

	client.httpClient = httpClient

	// Don't authenticate during initialization - let Authenticate() method handle it
	// if ok, err := client.auth(); !ok {
	// 	client.logger.Error("error while executing auth request", "request", "auth", "error", err.Error())
	// 	return nil, errors.New("error while executing auth request: " + err.Error())
	// }

	return client, nil
}

// auth authenticates the client with the MySubaru API using the provided credentials.
func (c *Client) auth() (bool, error) {
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
	resp, err := c.execute(POST, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing auth request", "request", "auth", "error", err.Error())
		return false, errors.New("error while executing auth request: " + err.Error())
	}
	c.logger.Debug("http request output", "request", "auth", "body", resp)

	// Check if response contains HTML instead of JSON (API error page)
	if isHTMLResponse(resp.Data) {
		c.logger.Error("received HTML error page instead of JSON", "request", "auth", "response_start", getResponsePreview(resp.Data, 200))
		return false, errHTMLResponse("auth")
	}

	var sd SessionData
	err = json.Unmarshal(resp.Data, &sd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "auth", "error", err.Error())
	}
	// client.logger.Debug("unmarshaled json data", "request", "auth", "type", "sessionData", "body", sd)

	if !sd.DeviceRegistered {
		err := c.getContactMethods()
		if err != nil {
			c.logger.Error("error while getting contact methods", "request", "auth", "error", err.Error())
			return false, errors.New("error while getting contact methods: " + err.Error())
		}

		c.logger.Error("device is not registered", "request", "auth", "deviceId", c.credentials.DeviceID)
		return false, errors.New("device is not registered: " + c.credentials.DeviceID)
	}

	if sd.DeviceRegistered && sd.RegisteredDevicePermanent {
		c.isAuthenticated = true
		c.isRegistered = true
		c.isAlive = true
	}
	c.logger.Debug("MySubaru API client authenticated")

	if len(sd.Vehicles) > 0 {
		for _, vehicle := range sd.Vehicles {
			c.listOfVins = append(c.listOfVins, vehicle.Vin)
		}
		c.currentVin = c.listOfVins[0]
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
	httpClient := resty.New()
	httpClient.
		SetBaseURL(MOBILE_API_SERVER[c.country]).
		SetHeaders(map[string]string{
			"User-Agent":       "Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.191030.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36",
			"Origin":           "file://",
			"X-Requested-With": MOBILE_APP[c.country],
			"Accept-Language":  "en-US,en;q=0.9",
			"Accept-Encoding":  "gzip, deflate",
			"Accept":           "*/*"},
		)
	c.httpClient = httpClient
	c.isAlive = false
}

// SelectVehicle selects a vehicle by its VIN. If no VIN is provided, it uses the current VIN.
func (c *Client) SelectVehicle(vin string) (*VehicleData, error) {
	if vin == "" {
		vin = c.currentVin
	}
	if err := ValidateVIN(vin); err != nil {
		return nil, err
	}

	params := map[string]string{
		"vin": vin,
		"_":   timestamp()}
	reqURL := MOBILE_API_VERSION + apiURLs["API_SELECT_VEHICLE"]
	resp, err := c.execute(GET, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing SelectVehicle request", "request", "SelectVehicle", "error", err.Error())
		return nil, errors.New("error while executing SelectVehicle request: " + err.Error())
	}
	// c.logger.Debug("http request output", "request", "SelectVehicle", "body", resp)

	// Check if response contains HTML instead of JSON (API error page)
	if isHTMLResponse(resp.Data) {
		c.logger.Error("received HTML error page instead of JSON", "request", "SelectVehicle", "response_start", getResponsePreview(resp.Data, 200))
		return nil, errHTMLResponse("SelectVehicle")
	}

	var vd VehicleData
	err = json.Unmarshal(resp.Data, &vd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "SelectVehicle", "error", err.Error())
		return nil, errors.New("error while parsing json while vehicle selection")
	}
	// c.logger.Debug("http request output", "request", "SelectVehicle", "body", resp)
	return &vd, nil
}

// GetVehicles retrieves a list of vehicles associated with the client's account.
func (c *Client) GetVehicles() ([]*Vehicle, error) {
	var vehicles []*Vehicle
	for _, vin := range c.listOfVins {
		vehicle, err := c.GetVehicleByVin(vin)
		if err != nil {
			c.logger.Error("cannot get vehicle data", "request", "GetVehicles", "error", err.Error())
			return nil, errors.New("cannot get vehicle data: " + err.Error())
		}
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, nil
}

// GetVehicleByVin retrieves a vehicle by its VIN from the client's list of vehicles.
func (c *Client) GetVehicleByVin(vin string) (*Vehicle, error) {
	var vehicle *Vehicle
	if slices.Contains(c.listOfVins, vin) {
		params := map[string]string{
			"vin": vin,
			"_":   timestamp()}
		reqURL := MOBILE_API_VERSION + apiURLs["API_SELECT_VEHICLE"]
		resp, err := c.execute(GET, reqURL, params, false)
		if err != nil {
			c.logger.Error("error while executing GetVehicleByVin request", "request", "GetVehicleByVin", "error", err.Error())
			return nil, errors.New("error while executing GetVehicleByVin request: " + err.Error())
		}
		// c.logger.Debug("http request output", "request", "GetVehicleByVin", "body", resp)

		// Check if response contains HTML instead of JSON (API error page)
		if isHTMLResponse(resp.Data) {
			c.logger.Error("received HTML error page instead of JSON", "request", "GetVehicleByVin", "response_start", getResponsePreview(resp.Data, 200))
			return nil, errHTMLResponse("GetVehicleByVin")
		}

		var vd VehicleData
		err = json.Unmarshal(resp.Data, &vd)
		if err != nil {
			c.logger.Error("error while parsing json", "request", "GetVehicleByVin", "error", err.Error())
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

		vehicle.GetVehicleStatus()
		vehicle.GetVehicleCondition()
		vehicle.GetVehicleHealth()

		// Get climate presets - log errors but don't fail the entire vehicle creation
		if err := vehicle.GetClimatePresets(); err != nil {
			c.logger.Warn("failed to get climate presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		if err := vehicle.GetClimateUserPresets(); err != nil {
			c.logger.Warn("failed to get climate user presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		if err := vehicle.GetClimateQuickPresets(); err != nil {
			c.logger.Warn("failed to get climate quick presets during vehicle initialization", "vin", vin, "error", err.Error())
		}

		return vehicle, nil
	}
	c.logger.Error("vin code is not in the list of the available vin codes", "request", "GetVehicleByVIN")
	return nil, errors.New("vin code is not in the list of the available vin codes")
}

func (c *Client) RefreshVehicles() error {
	params := map[string]string{}
	reqURL := MOBILE_API_VERSION + apiURLs["API_REFRESH_VEHICLES"]
	resp, err := c.execute(GET, reqURL, params, false)
	if err != nil {
		c.logger.Error("error while executing RefreshVehicles request", "request", "RefreshVehicles", "error", err.Error())
		return pkgerrors.Wrap(err, "RefreshVehicles request failed")
	}
	c.logger.Debug("http request output", "request", "RefreshVehicles", "body", resp)

	// Check if response contains HTML instead of JSON (API error page)
	if isHTMLResponse(resp.Data) {
		c.logger.Error("received HTML error page instead of JSON", "request", "RefreshVehicles", "response_start", getResponsePreview(resp.Data, 200))
		return pkgerrors.Wrap(errHTMLResponse("RefreshVehicles"), "RefreshVehicles request failed")
	}

	var sd SessionData
	err = json.Unmarshal(resp.Data, &sd)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "RefreshVehicles", "error", err.Error())
		return pkgerrors.Wrap(err, "failed to parse RefreshVehicles response")
	}

	// Update client state
	if sd.DeviceRegistered && sd.RegisteredDevicePermanent {
		c.isAuthenticated = true
		c.isRegistered = true
		c.isAlive = true
	}
	c.listOfVins = []string{}
	for _, vehicle := range sd.Vehicles {
		c.listOfVins = append(c.listOfVins, vehicle.Vin)
	}
	if len(c.listOfVins) > 0 {
		c.currentVin = c.listOfVins[0]
	}

	c.logger.Info("vehicles refreshed successfully", "count", len(c.listOfVins))
	return nil
}

// RequestAuthCode requests an authentication code for two-factor authentication (2FA).
// (?!^).(?=.*@)
// (?!^): This is a negative lookbehind assertion. It ensures that the matched character is not at the beginning of the string.
// .: This matches any single character (except newline, by default).
// (?=.*@): This is a positive lookahead assertion. It ensures that the matched character is followed by any characters (.*) and then an "@" symbol. This targets the username part of the email address.
func (c *Client) RequestAuthCode(email string) error {
	email, err := emailMasking(email)
	if err != nil {
		c.logger.Error("error while hiding email", "request", "RequestAuthCode", "error", err.Error())
		return errors.New("error while hiding email: " + err.Error())
	}

	if !containsValueInStruct(c.contactMethods, email) {
		c.logger.Error("email is not in the list of contact methods", "request", "RequestAuthCode", "email", email)
		return errors.New("email is not in the list of contact methods: " + email)
	}

	params := map[string]string{
		"contactMethod":      email,
		"languagePreference": "EN"}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_2FA_SEND_VERIFICATION"]
	resp, err := c.execute(POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing RequestAuthCode request", "request", "RequestAuthCode", "error", err.Error())
		return errors.New("error while executing RequestAuthCode request: " + err.Error())
	}
	c.logger.Debug("http request output", "request", "RequestAuthCode", "body", resp)

	return nil
}

// SubmitAuthCode submits the authentication code received from the RequestAuthCode method.
func (c *Client) SubmitAuthCode(code string, permanent bool) error {
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
	resp, err := c.execute(POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing SubmitAuthCode request", "request", "SubmitAuthCode", "error", err.Error())
		return errors.New("error while executing SubmitAuthCode request: " + err.Error())
	}
	c.logger.Debug("http request output", "request", "SubmitAuthCode", "body", resp)

	// Device registration does not always immediately take effect
	time.Sleep(time.Second * 3)

	// Reauthenticate after submitting the code
	if ok, err := c.auth(); !ok {
		c.logger.Error("error while executing auth request", "request", "auth", "error", err.Error())
		return errors.New("error while executing auth request: " + err.Error())
	}

	return nil
}

// getContactMethods retrieves the available contact methods for two-factor authentication (2FA).
// {"success":true,"dataName":"dataMap","data":{"userName":"a**x@savin.nyc","email":"t***a@savin.nyc"}}
func (c *Client) getContactMethods() error {
	// // Validate session before executing the request
	// if !c.validateSession() {
	// 	c.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
	// 	return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	// }

	params := map[string]string{}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_2FA_CONTACT"]
	resp, err := c.execute(POST, reqUrl, params, false)
	if err != nil {
		c.logger.Error("error while executing getContactMethods request", "request", "getContactMethods", "error", err.Error())
		return errors.New("error while executing getContactMethods request: " + err.Error())
	}
	c.logger.Debug("http request output", "request", "getContactMethods", "body", resp)

	// Check if response contains HTML instead of JSON (API error page)
	if isHTMLResponse(resp.Data) {
		c.logger.Error("received HTML error page instead of JSON", "request", "getContactMethods", "response_start", getResponsePreview(resp.Data, 200))
		return errHTMLResponse("getContactMethods")
	}

	var dm dataMap
	err = json.Unmarshal(resp.Data, &dm)
	if err != nil {
		c.logger.Error("error while parsing json", "request", "getContactMethods", "error", err.Error())
		return errors.New("error while parsing json while getting contact methods: " + err.Error())
	}
	c.contactMethods = dm
	c.logger.Debug("contact methods successfully retrieved", "request", "getContactMethods", "methods", dm)

	return nil
}

// RemoteUnlock unlocks the vehicle remotely.
func (c *Client) RemoteUnlock(vin string) error {
	if !slices.Contains(c.listOfVins, vin) {
		return pkgerrors.New("VIN not in list")
	}

	if !c.validateSession() {
		return pkgerrors.New("session is not valid")
	}

	vData, err := c.SelectVehicle(vin)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to select vehicle")
	}

	if !slices.Contains(vData.SubscriptionFeatures, FEATURE_REMOTE) {
		return pkgerrors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}
	reqURL := MOBILE_API_VERSION + apiURLs["API_UNLOCK"]
	params := map[string]string{"vin": vin}
	resp, err := c.execute(POST, reqURL, params, false)
	if err != nil {
		return pkgerrors.Wrap(err, "RemoteUnlock request failed")
	}
	if !resp.Success {
		return pkgerrors.New("RemoteUnlock failed")
	}
	c.logger.Info("Vehicle unlocked successfully", "vin", vin)
	return nil
}

// execute executes an HTTP request based on the method, URL, and parameters provided.
func (c *Client) execute(method string, url string, params map[string]string, j bool) (*Response, error) {
	return c.executeWithRetry(method, url, params, j, 3)
}

// executeWithRetry executes an HTTP request with retry logic
func (c *Client) executeWithRetry(method string, url string, params map[string]string, j bool, maxRetries int) (*Response, error) {
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
			time.Sleep(waitTime)
		}

		resp, err := c.executeOnce(method, url, params, j)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on non-retryable errors
		if !IsRetryableError(err) {
			break
		}

		// Handle InvalidToken by re-authenticating before retry
		var apiErr APIError
		if errors.As(err, &apiErr) && apiErr.Code == API_ERRORS["API_ERROR_INVALID_TOKEN"] {
			c.logger.Info("InvalidToken error detected, attempting re-authentication before retry")
			if !c.reauthenticateAndSelect() {
				c.logger.Error("re-authentication failed during retry")
				break
			}
		}

		c.logger.Warn("request failed, will retry", "attempt", attempt+1, "maxRetries", maxRetries, "error", err.Error())
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
		return nil, pkgerrors.New("unsupported HTTP method: " + method)
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
	if _, known := API_ERRORS[er.ErrorLabel]; known {
		c.logger.Error("known API error", "label", er.ErrorLabel, "description", er.ErrorDescription)
	} else {
		c.logger.Error("unknown API error", "label", er.ErrorLabel, "description", er.ErrorDescription)
	}
}

// handleVehicleSetupError handles the VEHICLESETUPERROR case and returns response/error accordingly.
func (c *Client) handleVehicleSetupError(r *Response, method, url string, duration time.Duration) (*Response, error) {
	// With vehicle data: treat as success (user needs to complete setup but data is functional)
	if r.DataName == "vehicle" {
		c.logger.Warn("VEHICLESETUPERROR received but vehicle data is present; treating as success",
			"errorCode", r.ErrorCode, "dataName", r.DataName)
		c.metrics.RecordRequest(method, url, duration, true)
		c.isAlive = true
		return r, nil
	}
	// Without vehicle data: reset session and allow retry
	c.logger.Warn("VEHICLESETUPERROR received without vehicle data; resetting session",
		"errorCode", r.ErrorCode, "dataName", r.DataName)
	c.resetSession()
	c.metrics.RecordRequest(method, url, duration, false)
	return nil, APIError{Code: r.ErrorCode, Message: "VEHICLESETUPERROR: session reset, please retry", Retryable: true}
}

func (c *Client) executeOnce(method string, url string, params map[string]string, j bool) (*Response, error) {
	start := time.Now()

	c.Lock()
	defer c.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := c.httpClient.R().SetContext(ctx)
	resp, err := c.sendRequest(req, method, url, params, j)
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
		return nil, pkgerrors.Wrap(err, "failed to read response body")
	}
	c.logger.Debug("received HTTP response", "method", method, "url", url, "status", resp.Status(), "body", string(resBytes))

	c.httpClient.SetCookies(resp.Cookies())

	r, ok := c.parseResponse(resBytes)
	if !ok {
		c.metrics.RecordRequest(method, url, duration, false)
		c.metrics.RecordError("parse_error")
		c.isAlive = false
		return nil, pkgerrors.New("failed to parse response")
	}

	if resp.IsSuccess() && r.Success {
		c.metrics.RecordRequest(method, url, duration, true)
		c.isAlive = true
		return &r, nil
	}

	c.handleAPIError(&r)

	if resp.IsSuccess() && !r.Success {
		c.logger.Warn("API returned success=false despite HTTP 200",
			"method", method, "url", url, "errorCode", r.ErrorCode, "dataName", r.DataName)

		if r.ErrorCode == API_ERRORS["API_ERROR_VEHICLE_SETUP"] {
			return c.handleVehicleSetupError(&r, method, url, duration)
		}

		c.metrics.RecordRequest(method, url, duration, false)
		c.isAlive = false
		if r.ErrorCode != "" {
			return nil, APIError{Code: r.ErrorCode, Message: fmt.Sprintf("API request failed: %s", r.ErrorCode), Retryable: true}
		}
		return nil, APIError{Code: "API_SUCCESS_FALSE", Message: "API request failed with success=false", Retryable: true}
	}

	c.metrics.RecordRequest(method, url, duration, false)
	c.isAlive = false
	return nil, pkgerrors.Errorf("request failed with status %s", resp.Status())
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

// ValidateSession checks if the current session is valid by making a request to the vehicle status API.
func (c *Client) validateSession() bool {
	if c == nil {
		return false
	}
	reqURL := MOBILE_API_VERSION + apiURLs["API_VALIDATE_SESSION"]
	resp, err := c.execute(GET, reqURL, map[string]string{}, false)
	if err != nil {
		c.logger.Error("error while executing validateSession request", "request", "validateSession", "error", err.Error())
		return c.reauthenticateAndSelect()
	}
	if resp == nil {
		c.logger.Warn("validateSession returned nil response; forcing re-auth")
		return c.reauthenticateAndSelect()
	}
	c.logger.Debug("http request output", "request", "validateSession", "body", resp)

	if resp.Success {
		if _, err := c.SelectVehicle(c.currentVin); err == nil {
			return true
		} else {
			errMsg := err.Error()
			c.logger.Warn("select vehicle failed during session validation; attempting re-auth", "request", "validateSession", "error", errMsg)
			return c.reauthenticateAndSelect()
		}
	}

	return c.reauthenticateAndSelect()
}

// reauthenticateAndSelect forces re-authentication and re-selects the current VIN.
// Returns true on success, false if either step fails.
func (c *Client) reauthenticateAndSelect() bool {
	if ok, err := c.auth(); !ok || err != nil {
		if err != nil {
			c.logger.Error("error while re-authenticating", "request", "validateSession", "error", err.Error())
		} else {
			c.logger.Error("reauthentication failed", "request", "validateSession")
		}
		return false
	}

	// Ensure we have a VIN to select. Fall back to the first known VIN after auth.
	if c.currentVin == "" && len(c.listOfVins) > 0 {
		c.currentVin = c.listOfVins[0]
	}

	if _, err := c.SelectVehicle(c.currentVin); err != nil {
		c.logger.Error("error while selecting vehicle", "request", "validateSession", "error", err.Error())
		return false
	}

	return true
}

// Authenticate attempts to authenticate the client and returns detailed authentication status
func (c *Client) Authenticate() (bool, error, bool) {
	// Try to authenticate
	ok, err := c.auth()

	// If authentication failed and the error indicates device is not registered,
	// return the 2FA required flag
	if !ok && err != nil && (strings.Contains(err.Error(), "device is not registered") || strings.Contains(err.Error(), "error while executing auth request")) {
		return false, err, true // Device registration required
	}

	return ok, err, false // 2FA not required
}
