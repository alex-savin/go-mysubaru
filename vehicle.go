package mysubaru

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultClimateTemp = "65"
	DefaultRunTime     = "10"
	DefaultFanSpeed    = "7"
	DefaultAirMode     = "FEET_WINDOW"
	DefaultSeatHeat    = "HIGH_COOL"
	DefaultCirculation = "outsideAir"
	DefaultStartConfig = "START_ENGINE_ALLOW_KEY_IN_IGNITION"

	// Service request constants
	MaxServiceRequestAttempts = 15
	ServiceRequestPollDelay   = 5 * time.Second
)

// Vehicle represents a Subaru vehicle with various attributes and methods to interact with it.
type Vehicle struct {
	CarId                int64
	Vin                  string   // SELECT CAR REQUEST > "vin": "4S4BTGND8L3137058"
	CarName              string   // SELECT CAR REQUEST > "vehicleName": "Subaru Outback LXT"
	CarNickname          string   // SELECT CAR REQUEST > "nickname": "Subaru Outback LXT"
	ExtDescrip           string   // SELECT CAR REQUEST > "extDescrip": "Abyss Blue Pearl"
	IntDescrip           string   // SELECT CAR REQUEST > "intDescrip": "Gray"
	ModelName            string   // SELECT CAR REQUEST > "modelName": "Outback",
	ModelYear            string   // SELECT CAR REQUEST > "modelYear": "2020"
	ModelCode            string   // SELECT CAR REQUEST > "modelCode": "LDJ"
	TransCode            string   // SELECT CAR REQUEST > "transCode": "CVT"
	EngineSize           float64  // SELECT CAR REQUEST > "engineSize": 2.4
	VehicleKey           int64    // SELECT CAR REQUEST > "vehicleKey": 3832950
	LicensePlate         string   // SELECT CAR REQUEST > "licensePlate": "8KV8"
	LicensePlateState    string   // SELECT CAR REQUEST > "licensePlateState": "NJ"
	Features             []string // SELECT CAR REQUEST > "features": ["ATF_MIL","11.6MMAN","ABS_MIL","CEL_MIL","ACCS","RCC","REARBRK","TEL_MIL","VDC_MIL","TPMS_MIL","WASH_MIL","BSDRCT_MIL","OPL_MIL","EYESIGHT","RAB_MIL","SRS_MIL","ESS_MIL","RESCC","EOL_MIL","BSD","EBD_MIL","EPB_MIL","RES","RHSF","AWD_MIL","NAV_TOMTOM","ISS_MIL","RPOIA","EPAS_MIL","RPOI","AHBL_MIL","SRH_MIL","g2"],
	SubscriptionFeatures []string // SELECT CAR REQUEST > "subscriptionFeatures": ["REMOTE","SAFETY","Retail"]
	SubscriptionStatus   string   // SELECT CAR REQUEST > "subscriptionStatus": "ACTIVE"
	EngineState          string   // STATUS REQUEST     > "vehicleStateType": "IGNITION_OFF"
	Odometer             struct {
		Miles      int // STATUS REQUEST > "odometerValue": 24999
		Kilometers int // STATUS REQUEST > "odometerValueKilometers": 40223
	}
	DistanceToEmpty struct {
		Miles         int // STATUS REQUEST > "distanceToEmptyFuelMiles": 149.75
		Kilometers    int // STATUS REQUEST > "distanceToEmptyFuelKilometers": 241
		Miles10s      int // STATUS REQUEST > "distanceToEmptyFuelMiles10s": 150
		Kilometers10s int // STATUS REQUEST > "distanceToEmptyFuelKilometers10s": 240
		Percentage    int // > "remainingFuelPercent": 66
	}
	FuelConsumptionAvg struct {
		MPG     float64 // STATUS REQUEST > "avgFuelConsumptionMpg": 18.5
		LP100Km float64 // STATUS REQUEST > "avgFuelConsumptionLitersPer100Kilometers": 12.7
	}
	ClimateProfiles map[string]ClimateProfile
	Doors           map[string]Door    // CONDITION REQUEST >
	Windows         map[string]Window  // CONDITION REQUEST >
	Tires           map[string]Tire    // CONDITION AND STATUS REQUEST >
	Troubles        map[string]Trouble //
	GeoLocation     GeoLocation
	// EV-specific fields
	EVStatus struct {
		StateOfChargePercent        int    // Battery charge percentage (0-100)
		DistanceToEmptyMiles        int    // Electric range in miles
		DistanceToEmptyKm           int    // Electric range in kilometers
		DistanceToEmptyByStateMiles int    // Electric range by state in miles
		DistanceToEmptyByStateKm    int    // Electric range by state in kilometers
		IsPluggedIn                 bool   // Whether vehicle is plugged in
		ChargerStateType            string // Charger state (e.g., "CHARGING", "NOT_CHARGING")
		StateOfChargeMode           string // Charge mode
		TimeToFullyCharged          string // Time remaining to full charge
	}
	Updated time.Time
	client  *Client
}

// MarshalJSON provides custom JSON serialization for Vehicle.
// It includes a computed "EV" field for backwards compatibility with clients
// that check this field to determine if the vehicle is electric.
func (v Vehicle) MarshalJSON() ([]byte, error) {
	type VehicleAlias Vehicle // Alias to avoid infinite recursion

	return json.Marshal(&struct {
		VehicleAlias
		EV bool `json:"EV"` // Computed field for backwards compatibility
	}{
		VehicleAlias: VehicleAlias(v),
		EV:           slices.Contains(v.Features, FEATURE_PHEV),
	})
}

// Door represents a door of a Subaru vehicle with its position, sub-position, status, and lock state.
type Door struct {
	Position    string // front | rear | boot | enginehood
	SubPosition string // right | left
	Status      string // CLOSED | OPEN
	Lock        string // LOCKED | UNLOCKED
	Updated     time.Time
}

// Window represents a window of a Subaru vehicle with its position, sub-position, status, and last updated time.
type Window struct {
	Position    string
	SubPosition string
	Status      string // CLOSE | VENTED | OPEN
	Updated     time.Time
}

// Tire represents a tire of a Subaru vehicle with its position, sub-position, pressure, pressure in PSI, and last updated time.
type Tire struct {
	Position    string
	SubPosition string
	Pressure    int
	PressurePsi int
	Updated     time.Time
	// Status string
}

// Trouble represents a trouble or issue with a Subaru vehicle, containing a description of the trouble.
type Trouble struct {
	Description string
}

func normalizeClimateProfile(rp map[string]any) ClimateProfile {
	toString := func(v any) string {
		switch t := v.(type) {
		case string:
			return t
		case float64:
			return strconv.FormatFloat(t, 'f', -1, 64)
		case json.Number:
			return t.String()
		case bool:
			return strconv.FormatBool(t)
		default:
			return fmt.Sprint(t)
		}
	}

	toInt := func(v any) int {
		switch t := v.(type) {
		case json.Number:
			if i, err := t.Int64(); err == nil {
				return int(i)
			}
		case float64:
			return int(t)
		case string:
			if i, err := strconv.Atoi(strings.TrimSpace(t)); err == nil {
				return i
			}
		}
		return 0
	}

	return ClimateProfile{
		Name:                      toString(rp["name"]),
		VehicleType:               toString(rp["vehicleType"]),
		PresetType:                toString(rp["presetType"]),
		StartConfiguration:        toString(rp["startConfiguration"]),
		RunTimeMinutes:            toInt(rp["runTimeMinutes"]),
		HeatedRearWindowActive:    toString(rp["heatedRearWindowActive"]),
		HeatedSeatFrontRight:      toString(rp["heatedSeatFrontRight"]),
		HeatedSeatFrontLeft:       toString(rp["heatedSeatFrontLeft"]),
		ClimateZoneFrontTemp:      toInt(rp["climateZoneFrontTemp"]),
		ClimateZoneFrontAirMode:   toString(rp["climateZoneFrontAirMode"]),
		ClimateZoneFrontAirVolume: toString(rp["climateZoneFrontAirVolume"]),
		OuterAirCirculation:       toString(rp["outerAirCirculation"]),
		AirConditionOn:            toString(rp["airConditionOn"]),
		CanEdit:                   toString(rp["canEdit"]),
		Disabled:                  toString(rp["disabled"]),
	}
}

func (v *Vehicle) String() string {
	var vString string
	vString += "=== INFORMATION =====================\n"
	vString += "Nickname: " + v.CarNickname + "\n"
	vString += "Car Name: " + v.CarName + "\n"
	vString += "Model: " + v.ModelName + "\n"

	vString += "=== ODOMETER =====================\n"
	vString += "Miles: " + strconv.Itoa(v.Odometer.Miles) + "\n"
	vString += "Kilometers: " + strconv.Itoa(v.Odometer.Kilometers) + "\n"

	vString += "=== DISTANCE TO EMPTY =====================\n"
	vString += "Miles: " + strconv.Itoa(v.DistanceToEmpty.Miles) + "\n"
	vString += "Kilometers: " + strconv.Itoa(v.DistanceToEmpty.Kilometers) + "\n"

	vString += "=== FUEL =============================\n"
	vString += "Tank (%): " + fmt.Sprintf("%v", v.DistanceToEmpty.Percentage) + "\n"
	vString += "MPG: " + fmt.Sprintf("%v", v.FuelConsumptionAvg.MPG) + "\n"
	vString += "Litres per 100 km: " + fmt.Sprintf("%v", v.FuelConsumptionAvg.LP100Km) + "\n"

	if v.IsEV() {
		vString += "=== ELECTRIC VEHICLE =================\n"
		vString += "Battery (%): " + fmt.Sprintf("%d", v.EVStatus.StateOfChargePercent) + "\n"
		vString += "Range (Miles): " + fmt.Sprintf("%d", v.EVStatus.DistanceToEmptyMiles) + "\n"
		vString += "Range (Km): " + fmt.Sprintf("%d", v.EVStatus.DistanceToEmptyKm) + "\n"
		vString += "Range by State (Miles): " + fmt.Sprintf("%d", v.EVStatus.DistanceToEmptyByStateMiles) + "\n"
		vString += "Range by State (Km): " + fmt.Sprintf("%d", v.EVStatus.DistanceToEmptyByStateKm) + "\n"
		vString += "Plugged In: " + fmt.Sprintf("%t", v.EVStatus.IsPluggedIn) + "\n"
		vString += "Charger State: " + v.EVStatus.ChargerStateType + "\n"
		vString += "Charge Mode: " + v.EVStatus.StateOfChargeMode + "\n"
		vString += "Time to Full Charge: " + v.EVStatus.TimeToFullyCharged + "\n"
	}

	vString += "=== GPS LOCATION ==============\n"
	vString += "Latitude: " + fmt.Sprintf("%v", v.GeoLocation.Latitude) + "\n"
	vString += "Longitude: " + fmt.Sprintf("%v", v.GeoLocation.Longitude) + "\n"
	vString += "Heading: " + fmt.Sprintf("%v", v.GeoLocation.Heading) + "\n"

	vString += "=== WINDOWS ===================\n"
	for k, v := range v.Windows {
		vString += fmt.Sprintf("%s >> %+v\n", k, v)
	}

	vString += "=== DOORS =====================\n"
	for k, v := range v.Doors {
		vString += fmt.Sprintf("%s >> %+v\n", k, v)
	}

	vString += "=== TIRES =====================\n"
	for k, v := range v.Tires {
		vString += fmt.Sprintf("%s >> %+v\n", k, v)
	}

	vString += "=== CLIMATE PROFILES ==========\n"
	for k, v := range v.ClimateProfiles {
		vString += fmt.Sprintf("%s >> %+v\n", k, v)
	}

	vString += "=== TROUBLES =====================\n"
	for k, v := range v.Troubles {
		vString += fmt.Sprintf("%s >> %+v\n", k, v)
	}

	vString += "=== FEATURES =====================\n"
	for i, f := range v.Features {
		if !strings.HasSuffix(f, "_MIL") {
			if _, ok := features[f]; ok {
				vString += fmt.Sprintf("%d >> %+v || %s\n", i+1, f, features[f])
			} else {
				vString += fmt.Sprintf("%d >> %+v\n", i+1, f)
			}
		}
	}
	return vString
}

// Lock
// Sends a command to lock doors.
func (v *Vehicle) Lock() (chan string, error) {
	params := map[string]string{
		"delay":         "0",
		"vin":           v.Vin,
		"pin":           v.client.credentials.PIN,
		"forceKeyInCar": "false"}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_LOCK"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// Unlock
// Send command to unlock doors.
func (v *Vehicle) Unlock() (chan string, error) {
	params := map[string]string{
		"delay":          "0",
		"vin":            v.Vin,
		"pin":            v.client.credentials.PIN,
		"unlockDoorType": "ALL_DOORS_CMD"} // FRONT_LEFT_DOOR_CMD | ALL_DOORS_CMD
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_UNLOCK"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// EngineStart
// Sends a command to start engine and set climate control.
func (v *Vehicle) EngineStart(run, delay int, horn bool) (chan string, error) {
	return v.EngineStartWithProfile(run, delay, horn, "")
}

// EngineStartWithProfile starts the engine using either a selected climate profile or defaults.
// If profileName matches an entry in ClimateProfiles, its values override defaults.
func (v *Vehicle) EngineStartWithProfile(run, delay int, horn bool, profileName string) (chan string, error) {
	// Validate run time parameter
	validRunTimes := []int{0, 1, 5, 10}
	if !slices.Contains(validRunTimes, run) {
		return nil, fmt.Errorf("run time must be one of %v minutes, got %d", validRunTimes, run)
	}

	// Validate delay parameter (reasonable bounds)
	if delay < 0 || delay > 60 {
		return nil, fmt.Errorf("delay must be between 0 and 60 minutes, got %d", delay)
	}

	// Defaults
	startConfig := START_CONFIG_DEFAULT_RES
	if v.IsEV() {
		startConfig = START_CONFIG_DEFAULT_EV
	}

	params := map[string]string{
		"delay":                     strconv.Itoa(delay),
		"vin":                       v.Vin,
		"pin":                       v.client.credentials.PIN,
		"horn":                      strconv.FormatBool(horn),
		"climateSettings":           "climateSettings",
		"climateZoneFrontTemp":      DefaultClimateTemp,
		"climateZoneFrontAirMode":   DefaultAirMode,
		"climateZoneFrontAirVolume": DefaultFanSpeed,
		"heatedSeatFrontLeft":       DefaultSeatHeat,
		"heatedSeatFrontRight":      DefaultSeatHeat,
		"heatedRearWindowActive":    "false",
		"outerAirCirculation":       DefaultCirculation,
		"airConditionOn":            "true",
		"runTimeMinutes":            strconv.Itoa(run),
		"startConfiguration":        startConfig,
	}

	// Apply profile overrides if present
	if profileName != "" {
		if cp, ok := v.ClimateProfiles[profileName]; ok {
			applyClimateProfile(params, cp)
		}
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_REMOTE_ENGINE_START"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// applyClimateProfile applies climate profile settings to the params map.
func applyClimateProfile(params map[string]string, cp ClimateProfile) {
	if cp.ClimateZoneFrontTemp != 0 {
		params["climateZoneFrontTemp"] = strconv.Itoa(cp.ClimateZoneFrontTemp)
	}
	if cp.ClimateZoneFrontAirMode != "" {
		params["climateZoneFrontAirMode"] = cp.ClimateZoneFrontAirMode
	}
	if cp.ClimateZoneFrontAirVolume != "" {
		params["climateZoneFrontAirVolume"] = cp.ClimateZoneFrontAirVolume
	}
	if cp.HeatedSeatFrontLeft != "" {
		params["heatedSeatFrontLeft"] = cp.HeatedSeatFrontLeft
	}
	if cp.HeatedSeatFrontRight != "" {
		params["heatedSeatFrontRight"] = cp.HeatedSeatFrontRight
	}
	if cp.HeatedRearWindowActive != "" {
		params["heatedRearWindowActive"] = cp.HeatedRearWindowActive
	}
	if cp.OuterAirCirculation != "" {
		params["outerAirCirculation"] = cp.OuterAirCirculation
	}
	if cp.AirConditionOn != "" {
		params["airConditionOn"] = cp.AirConditionOn
	}
	if cp.RunTimeMinutes != 0 {
		params["runTimeMinutes"] = strconv.Itoa(cp.RunTimeMinutes)
	}
	if cp.StartConfiguration != "" {
		params["startConfiguration"] = cp.StartConfiguration
	}
}

// EngineStop
// Sends a command to stop engine.
func (v *Vehicle) EngineStop() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_REMOTE_ENGINE_STOP"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// LightsStart
// Sends a command to flash lights.
func (v *Vehicle) LightsStart() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_LIGHTS"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// LightsStop
// Sends a command to stop flash lights.
func (v *Vehicle) LightsStop() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_LIGHTS_STOP"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// HornStart
// Send command to sound horn.
func (v *Vehicle) HornStart() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_HORN_LIGHTS"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// HornStop
// Send command to sound horn.
func (v *Vehicle) HornStop() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_HORN_LIGHTS_STOP"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// LockCancel
// Cancel an ongoing lock operation.
func (v *Vehicle) LockCancel() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_LOCK_CANCEL"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// UnlockCancel
// Cancel an ongoing unlock operation.
func (v *Vehicle) UnlockCancel() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_UNLOCK_CANCEL"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// EngineStartCancel
// Cancel an ongoing engine start operation.
func (v *Vehicle) EngineStartCancel() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_REMOTE_ENGINE_START_CANCEL"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// LightsCancel
// Cancel an ongoing lights operation.
func (v *Vehicle) LightsCancel() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_LIGHTS_CANCEL"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// HornLightsCancel
// Cancel an ongoing horn and lights operation.
func (v *Vehicle) HornLightsCancel() (chan string, error) {
	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_HORN_LIGHTS_CANCEL"], v.getAPIGen())
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]
	if v.getAPIGen() == FEATURE_G1_TELEMATICS {
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_HORN_LIGHTS_STATUS"]
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// ChargeOn
// Sends a command to start charging the EV.
func (v *Vehicle) ChargeOn() (chan string, error) {
	if !v.IsEV() {
		v.client.logger.Error("vehicle is not an EV")
		return nil, errors.New("vehicle is not an EV")
	}

	params := map[string]string{
		"delay": "0",
		"vin":   v.Vin,
		"pin":   v.client.credentials.PIN}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_EV_CHARGE_NOW"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// GetEVChargeSettings
// Retrieves the EV charging settings and schedules.
func (v *Vehicle) GetEVChargeSettings() error {
	if !v.IsEV() {
		v.client.logger.Error("vehicle is not an EV")
		return errors.New("vehicle is not an EV")
	}

	if err := v.validateSubscriptionAndSession(); err != nil {
		return err
	}

	v.ensureVehicleSelected()

	reqUrl := MOBILE_API_VERSION + apiURLs["API_EV_FETCH_CHARGE_SETTINGS"]
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error executing GetEVChargeSettings request", "error", err.Error())
		return err
	}

	// TODO: Parse and store EV charge settings
	v.client.logger.Debug("EV charge settings response", "data", string(resp.Data))

	return nil
}

// SaveEVChargeSettings
// Saves or updates the EV charging settings.
func (v *Vehicle) SaveEVChargeSettings(settings map[string]string) error {
	if !v.IsEV() {
		v.client.logger.Error("vehicle is not an EV")
		return errors.New("vehicle is not an EV")
	}

	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != v.client.currentVin {
		v.selectVehicle()
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_EV_SAVE_CHARGE_SETTINGS"]
	resp, err := v.client.execute(POST, reqUrl, settings, true)
	if err != nil {
		v.client.logger.Error("error executing SaveEVChargeSettings request", "error", err.Error())
		return err
	}

	v.client.logger.Debug("Save EV charge settings response", "data", string(resp.Data))

	return nil
}

// DeleteEVChargeSchedule
// Deletes an EV charging schedule.
func (v *Vehicle) DeleteEVChargeSchedule(scheduleID string) error {
	if !v.IsEV() {
		v.client.logger.Error("vehicle is not an EV")
		return errors.New("vehicle is not an EV")
	}

	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != v.client.currentVin {
		v.selectVehicle()
	}

	params := map[string]string{
		"scheduleId": scheduleID,
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_EV_DELETE_CHARGE_SCHEDULE"]
	resp, err := v.client.execute(POST, reqUrl, params, true)
	if err != nil {
		v.client.logger.Error("error executing DeleteEVChargeSchedule request", "error", err.Error())
		return err
	}

	v.client.logger.Debug("Delete EV charge schedule response", "data", string(resp.Data))

	return nil
}

// GetLocation retrieves the current location of the vehicle.
// If force is true, it sends a locate command to get real-time position.
// If force is false, it reports the last known location from Subaru's records.
// Returns a channel that will receive status updates about the location request.
func (v *Vehicle) GetLocation(force bool) (chan string, error) {
	var reqUrl, pollingUrl string
	var params map[string]string
	if force { // Sends a locate command to the vehicle to get real time position
		reqUrl = MOBILE_API_VERSION + apiURLs["API_G2_LOCATE_UPDATE"]
		pollingUrl = MOBILE_API_VERSION + apiURLs["API_G2_LOCATE_STATUS"]
		params = map[string]string{
			"vin": v.Vin,
			"pin": v.client.credentials.PIN}
		if v.getAPIGen() == FEATURE_G1_TELEMATICS {
			reqUrl = MOBILE_API_VERSION + apiURLs["API_G1_LOCATE_UPDATE"]
			pollingUrl = MOBILE_API_VERSION + apiURLs["API_G1_LOCATE_STATUS"]
		}
	} else { // Reports the last location the vehicle has reported to Subaru
		params = map[string]string{
			"vin": v.Vin,
			"pin": v.client.credentials.PIN}
		reqUrl = MOBILE_API_VERSION + urlToGen(apiURLs["API_LOCATE"], v.getAPIGen())
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()
	return ch, nil
}

// GetClimatePresets connects to the MySubaru API to download available climate presets.
// It first attempts to establish a connection with the MySubaru API.
// If successful and climate presets are found for the user's vehicle,
// it downloads them. If no presets are available, or if the connection fails,
// appropriate handling should be implemented within the function.
func (v *Vehicle) GetClimatePresets() error {
	if err := v.validateSubscriptionAndSession(); err != nil {
		return err
	}

	v.ensureVehicleSelected()
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_FETCH_RES_SUBARU_PRESETS"]
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error executing GetClimatePresets request", "error", err.Error())
		return err
	}

	result := v.parseClimateData(string(resp.Data))

	var cProfiles []ClimateProfile
	err = json.Unmarshal([]byte(result), &cProfiles)
	if err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetClimatePresets", "error", err.Error())
		return err
	}

	v.processClimateProfiles(cProfiles)
	v.Updated = time.Now()
	return nil
}

// GetClimateQuickPresets
// Used while user uses "quick start engine" button in the app
func (v *Vehicle) GetClimateQuickPresets() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != (v.client).currentVin {
		v.selectVehicle()
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_FETCH_RES_QUICK_START_SETTINGS"]
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error executing GetClimateQuickPresets request", "error", err.Error())
		return err
	}

	if resp == nil {
		v.client.logger.Error("received nil response from GetClimateQuickPresets request")
		return errors.New("received nil response from API")
	}

	// v.client.logger.Debug("http request output", "request", "GetClimateQuickPresets", "body", resp)

	re1 := regexp.MustCompile(`\"`)
	result := re1.ReplaceAllString(string(resp.Data), "")
	re2 := regexp.MustCompile(`\\`)
	result = re2.ReplaceAllString(result, `"`) // \u0022

	var cp ClimateProfile
	err = json.Unmarshal([]byte(result), &cp)
	if err != nil {
		v.client.logger.Error("error while parsing climate quick presets json", "request", "GetClimateQuickPresets", "error", err.Error())
		return err
	}

	re := regexp.MustCompile(`([A-Z])`)
	cpn := strings.ToLower("quick_" + re.ReplaceAllString(cp.PresetType, "_$1") + "_" + strings.ReplaceAll(cp.Name, " ", "_"))

	v.ClimateProfiles[cpn] = cp
	v.Updated = time.Now()
	return nil
}

// UpdateClimateQuickPresets
// Updates the quick climate presets by fetching them from the MySubaru API.
// {"success":true,"data":null}
func (v *Vehicle) UpdateClimateQuickPresets() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != (v.client).currentVin {
		v.selectVehicle()
	}

	params := map[string]string{
		"name":                      "Cooling",
		"runTimeMinutes":            "10",
		"climateSettings":           "climateSettings",                    // climateSettings
		"climateZoneFrontTemp":      "65",                                 // 60-86
		"climateZoneFrontAirMode":   "FEET_WINDOW",                        // FEET_FACE_BALANCED | FEET_WINDOW | WINDOW | FEET
		"climateZoneFrontAirVolume": "7",                                  // 1-7
		"heatedSeatFrontLeft":       "HIGH_COOL",                          // OFF | LOW_HEAT | MEDIUM_HEAT | HIGH_HEAT | LOW_COOL | MEDIUM_COOL |  HIGH_COOL
		"heatedSeatFrontRight":      "HIGH_COOL",                          // ---//---
		"heatedRearWindowActive":    "false",                              // boolean
		"outerAirCirculation":       "outsideAir",                         // outsideAir | recirculation
		"airConditionOn":            "false",                              // boolean
		"startConfiguration":        "START_ENGINE_ALLOW_KEY_IN_IGNITION", // START_ENGINE_ALLOW_KEY_IN_IGNITION | ONLY FOR PHEV > START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_SAVE_RES_QUICK_START_SETTINGS"]
	resp, _ := v.client.execute(POST, reqUrl, params, true)

	v.client.logger.Debug("http request output", "request", "UpdateClimateUserPresets", "body", resp)

	return nil
}

// GetClimateUserPresets retrieves user-defined climate presets from the MySubaru API.
// These are custom presets created by the user through the Subaru app.
func (v *Vehicle) GetClimateUserPresets() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != (v.client).currentVin {
		v.selectVehicle()
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_FETCH_RES_USER_PRESETS"]
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error executing GetClimateUserPresets request", "error", err.Error())
		return err
	}

	if resp == nil {
		v.client.logger.Error("received nil response from GetClimateUserPresets request")
		return errors.New("received nil response from API")
	}

	re1 := regexp.MustCompile(`\"`)
	result := re1.ReplaceAllString(string(resp.Data), "")
	re2 := regexp.MustCompile(`\\`)
	result = re2.ReplaceAllString(result, `"`) // \u0022

	var rawProfiles []map[string]any
	if err := json.Unmarshal([]byte(result), &rawProfiles); err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetClimateUserPresets", "error", err.Error())
		return err
	}

	profiles := make([]ClimateProfile, 0, len(rawProfiles))
	for _, rp := range rawProfiles {
		profiles = append(profiles, normalizeClimateProfile(rp))
	}

	v.processClimateProfiles(profiles)
	v.Updated = time.Now()
	return nil
}

// UpdateClimateUserPresets
// Updates the user's climate presets by fetching them from the MySubaru API.
func (v *Vehicle) UpdateClimateUserPresets() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != (v.client).currentVin {
		v.selectVehicle()
	}
	params := map[string]string{
		"presetType":                "userPreset",
		"name":                      "Cooling",
		"runTimeMinutes":            "10",
		"climateZoneFrontTemp":      "65",
		"climateZoneFrontAirMode":   "FEET_FACE_BALANCED",
		"climateZoneFrontAirVolume": "7",
		"outerAirCirculation":       "outsideAir",
		"heatedRearWindowActive":    "false",
		"heatedSeatFrontLeft":       "HIGH_COOL",
		"airConditionOn":            "false",
		"startConfiguration":        "START_ENGINE_ALLOW_KEY_IN_IGNITION",
		// "canEdit":                   "true",
		// "disabled":                  "false",
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_SAVE_RES_SETTINGS"]
	resp, _ := v.client.execute(POST, reqUrl, params, false)

	v.client.logger.Debug("http request output", "request", "UpdateClimateUserPresets", "body", resp)

	return nil
}

// SaveClimateUserPresets saves a list of user-defined climate presets to MySubaru.
// This overwrites all existing user presets with the provided list.
// Maximum of 4 user presets are allowed.
func (v *Vehicle) SaveClimateUserPresets(presets []ClimateProfile) error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	if len(presets) > 4 {
		return errors.New("maximum of 4 user presets allowed")
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != v.client.currentVin {
		v.selectVehicle()
	}

	// Ensure all presets have required fields
	for i := range presets {
		presets[i].PresetType = "userPreset"
		presets[i].CanEdit = "true"
		presets[i].Disabled = "false"
		if v.IsEV() {
			presets[i].StartConfiguration = "START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION"
		} else {
			presets[i].StartConfiguration = "START_ENGINE_ALLOW_KEY_IN_IGNITION"
		}
	}

	// Convert presets to JSON for the request body
	presetsJSON, err := json.Marshal(presets)
	if err != nil {
		v.client.logger.Error("error marshaling presets to JSON", "error", err.Error())
		return err
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_SAVE_RES_SETTINGS"]
	resp, err := v.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(presetsJSON).
		Post(v.client.httpClient.BaseURL() + reqUrl)
	if err != nil {
		v.client.logger.Error("error saving climate presets", "error", err.Error())
		return err
	}

	v.client.logger.Debug("http request output", "request", "SaveClimateUserPresets", "status", resp.StatusCode())

	// Refresh the presets after saving
	return v.GetClimateUserPresets()
}

// DeleteClimateUserPreset removes a user-defined climate preset by name.
// This works by fetching all user presets, removing the target, and saving the updated list.
func (v *Vehicle) DeleteClimateUserPreset(presetName string) error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Ensure we have the latest presets
	if err := v.GetClimateUserPresets(); err != nil {
		return fmt.Errorf("failed to fetch user presets: %w", err)
	}

	// Find and collect user presets, excluding the one to delete
	var updatedPresets []ClimateProfile
	found := false

	for _, profile := range v.ClimateProfiles {
		if strings.EqualFold(profile.PresetType, "userPreset") {
			if strings.EqualFold(profile.Name, presetName) {
				found = true
				continue // Skip this preset (delete it)
			}
			updatedPresets = append(updatedPresets, profile)
		}
	}

	if !found {
		return fmt.Errorf("user preset '%s' not found", presetName)
	}

	v.client.logger.Info("deleting climate preset", "name", presetName, "remaining", len(updatedPresets))

	// Save the updated list (without the deleted preset)
	return v.SaveClimateUserPresets(updatedPresets)
}

// GetVehicleStatus .
// updateVehicleFromStatus updates basic vehicle fields from VehicleStatus data.
func (v *Vehicle) updateVehicleFromStatus(vs *VehicleStatus) {
	v.EngineState = vs.VehicleStateType
	v.Odometer.Miles = vs.OdometerValue
	v.Odometer.Kilometers = vs.OdometerValueKm
	v.DistanceToEmpty.Miles = int(vs.DistanceToEmptyFuelMiles)
	v.DistanceToEmpty.Kilometers = vs.DistanceToEmptyFuelKilometers
	v.DistanceToEmpty.Miles10s = vs.DistanceToEmptyFuelMiles10s
	v.DistanceToEmpty.Kilometers10s = vs.DistanceToEmptyFuelKilometers10s
	if vs.RemainingFuelPercent >= 0 && vs.RemainingFuelPercent <= 100 {
		v.DistanceToEmpty.Percentage = vs.RemainingFuelPercent
	}
	v.FuelConsumptionAvg.MPG = float64(vs.AvgFuelConsumptionMpg)
	v.FuelConsumptionAvg.LP100Km = float64(vs.AvgFuelConsumptionLitersPer100Kilometers)
	v.GeoLocation.Latitude = float64(vs.Latitude)
	v.GeoLocation.Longitude = float64(vs.Longitude)
	v.GeoLocation.Heading = vs.Heading
}

// updateEVStatusFromStatus updates EV-specific fields if this is an EV.
func (v *Vehicle) updateEVStatusFromStatus(vs *VehicleStatus) {
	if !v.IsEV() {
		return
	}
	v.EVStatus.StateOfChargePercent = int(vs.EvStateOfChargePercent)
	v.EVStatus.DistanceToEmptyMiles = vs.EvDistanceToEmptyMiles
	v.EVStatus.DistanceToEmptyKm = vs.EvDistanceToEmptyKilometers
	v.EVStatus.DistanceToEmptyByStateMiles = vs.EvDistanceToEmptyByStateMiles
	v.EVStatus.DistanceToEmptyByStateKm = vs.EvDistanceToEmptyByStateKilometers
}

// isPartField checks if a field name represents a parseable vehicle part.
func isPartField(name string) bool {
	return (strings.HasPrefix(name, "Door") && strings.HasSuffix(name, "Position")) ||
		(strings.HasPrefix(name, "Door") && strings.HasSuffix(name, "LockStatus")) ||
		(strings.HasPrefix(name, "Window") && strings.HasSuffix(name, "Status")) ||
		strings.HasPrefix(name, "TirePressure")
}

func (v *Vehicle) GetVehicleStatus() error {
	if err := v.validateSubscriptionAndSession(); err != nil {
		return err
	}

	v.ensureVehicleSelected()
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_VEHICLE_STATUS"], v.getAPIGen())
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error while executing GetVehicleStatus request", "request", "GetVehicleStatus", "error", err.Error())
		return err
	}

	var vs VehicleStatus
	if err = json.Unmarshal(resp.Data, &vs); err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetVehicleStatus", "error", err.Error())
	}

	v.updateVehicleFromStatus(&vs)
	v.updateEVStatusFromStatus(&vs)

	val := reflect.ValueOf(vs)
	typeOfS := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if isBadValue(val.Field(i).Interface()) {
			continue
		}
		name := typeOfS.Field(i).Name
		if isPartField(name) {
			v.parseParts(name, val.Field(i).Interface())
		}
	}
	v.Updated = time.Now()
	return nil
}

func isBadValue(val any) bool {
	switch t := val.(type) {
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return true
		}
		upper := strings.ToUpper(s)
		if upper == "UNKNOWN" || upper == "NOT_EQUIPPED" {
			return true
		}
		if s == "None" || s == "16383" || s == "65535" || s == "-64" {
			return true
		}
	case int:
		if t == 0 {
			return true
		}
	case float64:
		if t == 0 {
			return true
		}
	case nil:
		return true
	}
	return false
}

// GetVehicleCondition retrieves the current condition/status of various vehicle components
// such as doors, windows, and tires from the MySubaru API.
func (v *Vehicle) GetVehicleCondition() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != (v.client).currentVin {
		v.selectVehicle()
	}
	reqUrl := MOBILE_API_VERSION + urlToGen(apiURLs["API_CONDITION"], v.getAPIGen())
	resp, err := v.client.execute(GET, reqUrl, map[string]string{}, false)
	if err != nil {
		v.client.logger.Error("error executing GetVehicleCondition request", "error", err.Error())
		return err
	}

	if resp == nil {
		v.client.logger.Error("received nil response from GetVehicleCondition request")
		return errors.New("received nil response from API")
	}

	// v.client.logger.Info("http request output", "request", "GetVehicleCondition", "body", resp)

	var sr ServiceRequest
	err = json.Unmarshal(resp.Data, &sr)
	if err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetVehicleCondition", "error", err.Error())
		return err
	}
	// v.client.logger.Debug("http request output", "request", "GetVehicleCondition", "body", resp)

	var vc VehicleCondition
	err = json.Unmarshal(sr.Result, &vc)
	if err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetVehicleCondition", "error", err.Error())
	}
	// v.client.logger.Debug("http request output", "request", "GetVehicleCondition", "body", resp)

	// Parse EV-specific fields if this is an EV
	if v.IsEV() {
		v.EVStatus.StateOfChargePercent = vc.EvStateOfChargePercent
		v.EVStatus.IsPluggedIn = vc.EvIsPluggedIn
		v.EVStatus.ChargerStateType = vc.EvChargerStateType
		v.EVStatus.StateOfChargeMode = vc.EvStateOfChargeMode
		v.EVStatus.TimeToFullyCharged = vc.EvTimeToFullyCharged
	}

	val := reflect.ValueOf(vc)
	typeOfS := val.Type()

	for i := 0; i < val.NumField(); i++ {
		// v.client.logger.Debug("vehicle condition >> parsing a car part", "field", typeOfS.Field(i).Name, "value", val.Field(i).Interface(), "type", val.Field(i).Type())
		if isBadValue(val.Field(i).Interface()) {
			continue
		}

		name := typeOfS.Field(i).Name
		// Lock status must come from vehicleStatus; condition endpoint is position-only.
		if strings.HasPrefix(name, "Door") && strings.HasSuffix(name, "Position") ||
			strings.HasPrefix(name, "Window") && strings.HasSuffix(name, "Status") {
			v.parseParts(name, val.Field(i).Interface())
		}
		// if strings.HasPrefix(name, "TirePressure") {
		// 	v.parseParts(name, val.Field(i).Interface())
		// }
	}
	v.Updated = time.Now()
	return nil
}

// GetVehicleHealth
// Retrieves the vehicle health status from MySubaru API.
func (v *Vehicle) GetVehicleHealth() error {
	if err := v.validateSubscriptionAndSession(); err != nil {
		return err
	}

	v.ensureVehicleSelected()
	params := map[string]string{
		"vin": v.Vin,
		"_":   timestamp()}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_VEHICLE_HEALTH"]
	resp, err := v.client.execute(GET, reqUrl, params, false)
	if err != nil {
		v.client.logger.Error("error executing GetVehicleHealth request", "error", err.Error())
		return err
	}

	if resp == nil {
		v.client.logger.Error("received nil response from GetVehicleHealth request")
		return errors.New("received nil response from API")
	}

	// v.client.logger.Debug("http request output", "request", "GetVehicleHealth", "body", resp)

	var vh VehicleHealth
	err = json.Unmarshal(resp.Data, &vh)
	if err != nil {
		v.client.logger.Error("error while parsing json", "request", "GetVehicleHealth", "error", err.Error())
		return err
	}
	// v.client.logger.Debug("http request output", "request", "GetVehicleHealth", "vehicle health", vh)

	for i, vhi := range vh.VehicleHealthItems {
		// v.client.logger.Debug("vehicle health item", "id", i, "item", vhi)
		if vhi.IsTrouble {
			if _, ok := troubles[vhi.FeatureCode]; ok {
				t := Trouble{
					Description: troubles[vhi.FeatureCode],
				}
				v.Troubles[vhi.FeatureCode] = t
				v.client.logger.Debug("found troubled vehicle health item", "id", i, "item", vhi.FeatureCode, "description", troubles[vhi.FeatureCode])
			}
		}
	}
	return nil
}

// GetFeaturesList logs all vehicle features and their descriptions to the debug logger.
// This is primarily used for debugging and understanding what features are available.
func (v *Vehicle) GetFeaturesList() {
	for _, f := range v.Features {
		if _, ok := features[f]; ok {
			v.client.logger.Debug("vehicle features", "id", f, "feature", features[f])
		} else {
			v.client.logger.Debug("vehicle features", "id", f)
		}
	}
}

// executeServiceRequest
// Executes a service request to the Subaru API and handles the response.
func (v *Vehicle) executeServiceRequest(params map[string]string, reqUrl, pollingUrl string, ch chan string, attempt int) error {
	if attempt >= MaxServiceRequestAttempts {
		v.client.logger.Error("maximum attempts reached for service request", "request", reqUrl, "attempts", attempt)
		ch <- "error"
		return errors.New("maximum attempts reached for service request")
	}

	// Check if the vehicle has a valid subscription for remote services
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	// Validate session before executing the request
	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	if v.Vin != v.client.currentVin {
		v.selectVehicle()
	}

	var resp *Response
	var err error
	if attempt == 1 {
		resp, err = v.client.execute(POST, reqUrl, params, true)
		if err != nil {
			v.client.logger.Error("error while executing service request", "request", reqUrl, "error", err.Error())
			ch <- "error"
			return err
		}
	} else {
		resp, err = v.client.execute(GET, pollingUrl, params, false)
		if err != nil {
			v.client.logger.Error("error while executing service request status polling", "request", reqUrl, "error", err.Error())
			ch <- "error"
			return err
		}
	}

	// dataName field has the list of the states [ remoteServiceStatus | errorResponse ]
	if resp.DataName == "remoteServiceStatus" {
		if sr, ok := v.parseServiceRequest([]byte(resp.Data)); ok {
			ch <- sr.RemoteServiceState
			switch sr.RemoteServiceState {

			case "finished":
				// Finished RemoteServiceState Service Request does not include Service Request ID
				v.client.logger.Debug("Remote service request completed successfully")

			case "started":
				time.Sleep(ServiceRequestPollDelay)
				v.client.logger.Debug("MySubaru API reports remote service request (started) is in progress", "id", sr.ServiceRequestID)
				v.executeServiceRequest(map[string]string{"serviceRequestId": sr.ServiceRequestID}, reqUrl, pollingUrl, ch, attempt+1)

			case "stopping":
				time.Sleep(ServiceRequestPollDelay)
				v.client.logger.Debug("MySubaru API reports remote service request (stopping) is in progress", "id", sr.ServiceRequestID)
				v.executeServiceRequest(map[string]string{"serviceRequestId": sr.ServiceRequestID}, reqUrl, pollingUrl, ch, attempt+1)

			default:
				v.client.logger.Debug("MySubaru API reports remote service request (default)")
				v.executeServiceRequest(map[string]string{"serviceRequestId": sr.ServiceRequestID}, reqUrl, pollingUrl, ch, attempt+1)
			}
			return nil
		}
		v.client.logger.Error("error while parsing service request json", "request", reqUrl, "response", resp.Data)
		return errors.New("error while parsing service request json")
	}
	return errors.New("response is not a service request")
}

// parseServiceRequest parses the JSON response from a service request into a ServiceRequest struct.
// Returns the parsed ServiceRequest and a boolean indicating success.
func (v *Vehicle) parseServiceRequest(b []byte) (ServiceRequest, bool) {
	var sr ServiceRequest
	err := json.Unmarshal(b, &sr)
	if err != nil {
		v.client.logger.Error("error while parsing service request json", "error", err.Error())
		return sr, false
	}
	return sr, true
}

// selectVehicle selects this vehicle in the client's session if it's not already selected.
// This ensures that subsequent API calls operate on the correct vehicle.
func (v *Vehicle) selectVehicle() {
	if v.client.currentVin != v.Vin {
		vData, err := (v.client).SelectVehicle(v.Vin)
		if err != nil {
			v.client.logger.Debug("cannot get vehicle data")
		}
		v.SubscriptionStatus = vData.SubscriptionStatus
		v.GeoLocation.Latitude = vData.VehicleGeoPosition.Latitude
		v.GeoLocation.Longitude = vData.VehicleGeoPosition.Longitude
		v.GeoLocation.Heading = vData.VehicleGeoPosition.Heading
		v.GeoLocation.Speed = vData.VehicleGeoPosition.Speed
		v.GeoLocation.Updated = vData.VehicleGeoPosition.Timestamp
		v.Updated = time.Now()
	}
}

// getAPIGen returns the Subaru telematics API generation (g1, g2, g3) for this vehicle
// based on the features present in the vehicle configuration.
func (v *Vehicle) getAPIGen() string {
	if slices.Contains(v.Features, FEATURE_G1_TELEMATICS) {
		return "g1"
	}
	if slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return "g2"
	}
	if slices.Contains(v.Features, FEATURE_G3_TELEMATICS) {
		return "g3"
	}
	return "unknown"
}

// IsEV returns true if this vehicle is an electric vehicle (PHEV).
func (v *Vehicle) IsEV() bool {
	return slices.Contains(v.Features, FEATURE_PHEV)
}

// getRemoteOptionsStatus returns true if this vehicle has remote service options available
// (requires appropriate subscription features).
func (v *Vehicle) getRemoteOptionsStatus() bool {
	return slices.Contains(v.SubscriptionFeatures, FEATURE_REMOTE)
}

// parseParts parses vehicle component data from API responses and updates the corresponding
// vehicle structures (doors, windows, tires) based on the field name and value.
func (v *Vehicle) parseParts(name string, value any) {
	re := regexp.MustCompile(`([Dd]oor|[Ww]indow|[Tt]ire)(?:[Pp]ressure)?([Ff]ront|[Rr]ear|[Bb]oot|[Ee]ngine[Hh]ood|[Ss]unroof)([Ll]eft|[Rr]ight)?([Pp]osition|[Ss]tatus|[Ll]ock[Ss]tatus|[Pp]si)?`)
	grps := re.FindStringSubmatch(name)

	if len(grps) < 2 {
		return
	}

	pn := strings.ToLower(grps[1] + "_" + grps[2])
	if len(grps[3]) > 0 {
		pn = pn + "_" + strings.ToLower(grps[3])
	}

	partType := strings.ToLower(grps[1])
	switch partType {
	case "door":
		v.parseDoor(pn, grps, value)
	case "window":
		v.parseWindow(pn, grps, value)
	case "tire":
		v.parseTire(pn, grps, value)
	}
}

// parseDoor handles door-specific parsing logic.
func (v *Vehicle) parseDoor(pn string, grps []string, value any) {
	d, exists := v.Doors[pn]
	if !exists {
		d = Door{
			Position:    grps[2],
			SubPosition: grps[3],
		}
	}

	v.applyDoorValue(&d, grps[4], value)
	d.Updated = time.Now()
	v.Doors[pn] = d
}

// applyDoorValue applies the value to the appropriate door field based on the field type.
func (v *Vehicle) applyDoorValue(d *Door, fieldType string, value any) {
	s, ok := value.(string)
	if !ok {
		return
	}
	normalized := strings.ToUpper(strings.TrimSpace(s))

	switch fieldType {
	case "Position":
		d.Status = normalized
	case "LockStatus":
		d.Lock = normalized
	}
}

// parseWindow handles window-specific parsing logic.
func (v *Vehicle) parseWindow(pn string, grps []string, value any) {
	w, exists := v.Windows[pn]
	if !exists {
		w = Window{
			Position:    grps[2],
			SubPosition: grps[3],
		}
	}

	if s, ok := value.(string); ok {
		w.Status = s
	}
	w.Updated = time.Now()
	v.Windows[pn] = w
}

// parseTire handles tire-specific parsing logic.
func (v *Vehicle) parseTire(pn string, grps []string, value any) {
	t, exists := v.Tires[pn]
	if !exists {
		t = Tire{
			Position:    grps[2],
			SubPosition: grps[3],
		}
	}

	pressure := toInt(value)
	if grps[4] == "Psi" {
		t.PressurePsi = pressure
	} else {
		t.Pressure = pressure
	}
	t.Updated = time.Now()
	v.Tires[pn] = t
}

// toInt converts a value to int, supporting int and float64 types.
func toInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	}
	return 0
}

// SetGeoFence sets up a geofence for the vehicle with specified parameters.
// This is a G2-only feature that requires Safety Plus subscription.
// Parameters:
//   - latitude: Center latitude of the geofence
//   - longitude: Center longitude of the geofence
//   - radius: Radius in meters (typically 100-1000)
//   - name: Name for the geofence
//   - enabled: Whether the geofence is active
//   - entryAlert: Alert when vehicle enters the geofence
//   - exitAlert: Alert when vehicle exits the geofence
func (v *Vehicle) SetGeoFence(latitude, longitude float64, radius int, name string, enabled, entryAlert, exitAlert bool) (chan string, error) {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return nil, errors.New("geofence feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return nil, errors.New("geofence feature requires Safety Plus subscription")
	}

	// Validate input parameters
	if err := ValidateCoordinates(latitude, longitude); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	if err := ValidateGeoFenceRadius(radius); err != nil {
		return nil, fmt.Errorf("invalid radius: %w", err)
	}

	if name == "" {
		return nil, errors.New("geofence name cannot be empty")
	}

	params := map[string]string{
		"delay":      "0",
		"vin":        v.Vin,
		"pin":        v.client.credentials.PIN,
		"latitude":   fmt.Sprintf("%.6f", latitude),
		"longitude":  fmt.Sprintf("%.6f", longitude),
		"radius":     strconv.Itoa(radius),
		"name":       name,
		"enabled":    strconv.FormatBool(enabled),
		"entryAlert": strconv.FormatBool(entryAlert),
		"exitAlert":  strconv.FormatBool(exitAlert),
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_GEOFENCE"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// UpdateGeoFence updates an existing geofence with new parameters.
// Parameters:
//   - fenceId: ID of the geofence to update
//   - latitude: New center latitude (optional, use 0 to keep current)
//   - longitude: New center longitude (optional, use 0 to keep current)
//   - radius: New radius in meters (optional, use 0 to keep current)
//   - name: New name (optional, use empty string to keep current)
//   - enabled: New enabled status
//   - entryAlert: New entry alert setting
//   - exitAlert: New exit alert setting
func (v *Vehicle) UpdateGeoFence(fenceId string, latitude, longitude float64, radius int, name string, enabled, entryAlert, exitAlert bool) (chan string, error) {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return nil, errors.New("geofence feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return nil, errors.New("geofence feature requires Safety Plus subscription")
	}

	params := map[string]string{
		"delay":      "0",
		"vin":        v.Vin,
		"pin":        v.client.credentials.PIN,
		"fenceId":    fenceId,
		"enabled":    strconv.FormatBool(enabled),
		"entryAlert": strconv.FormatBool(entryAlert),
		"exitAlert":  strconv.FormatBool(exitAlert),
	}

	// Only include optional parameters if they are provided
	if latitude != 0 {
		params["latitude"] = fmt.Sprintf("%.6f", latitude)
	}
	if longitude != 0 {
		params["longitude"] = fmt.Sprintf("%.6f", longitude)
	}
	if radius != 0 {
		params["radius"] = strconv.Itoa(radius)
	}
	if name != "" {
		params["name"] = name
	}

	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_GEOFENCE"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// DeleteGeoFence removes a geofence from the vehicle.
// Parameters:
//   - fenceId: ID of the geofence to delete
func (v *Vehicle) DeleteGeoFence(fenceId string) (chan string, error) {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return nil, errors.New("geofence feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return nil, errors.New("geofence feature requires Safety Plus subscription")
	}

	params := map[string]string{
		"delay":   "0",
		"vin":     v.Vin,
		"pin":     v.client.credentials.PIN,
		"fenceId": fenceId,
		"delete":  "true",
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_GEOFENCE"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// GetGeoFenceStatus retrieves the current status of all geofences for the vehicle.
// This method fetches information about active geofences and their current status.
func (v *Vehicle) GetGeoFenceStatus() error {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return errors.New("geofence feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return errors.New("geofence feature requires Safety Plus subscription")
	}

	// This would typically make a GET request to retrieve geofence status
	// For now, we'll implement it as a placeholder that could be expanded
	// based on the actual API response structure
	return errors.New("GetGeoFenceStatus not yet implemented - requires API research")
}

// SetSpeedFence sets up a speed fence for the vehicle.
// This is a G2-only feature that alerts when the vehicle exceeds a speed limit.
// Parameters:
//   - speedLimit: Speed limit in mph
//   - enabled: Whether the speed fence is active
//   - persistent: Whether to keep the setting across restarts
func (v *Vehicle) SetSpeedFence(speedLimit int, enabled, persistent bool) (chan string, error) {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return nil, errors.New("speed fence feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return nil, errors.New("speed fence feature requires Safety Plus subscription")
	}

	// Validate input parameters
	if err := ValidateSpeedLimit(speedLimit); err != nil {
		return nil, fmt.Errorf("invalid speed limit: %w", err)
	}

	params := map[string]string{
		"delay":      "0",
		"vin":        v.Vin,
		"pin":        v.client.credentials.PIN,
		"speedLimit": strconv.Itoa(speedLimit),
		"enabled":    strconv.FormatBool(enabled),
		"persistent": strconv.FormatBool(persistent),
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_SPEEDFENCE"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// SetCurfew sets up a curfew for the vehicle.
// This is a G2-only feature that restricts vehicle operation during specified times.
// Parameters:
//   - startTime: Start time in HH:MM format (24-hour)
//   - endTime: End time in HH:MM format (24-hour)
//   - daysOfWeek: Array of days (0=Sunday, 6=Saturday)
//   - enabled: Whether the curfew is active
func (v *Vehicle) SetCurfew(startTime, endTime string, daysOfWeek []int, enabled bool) (chan string, error) {
	if !slices.Contains(v.Features, FEATURE_G2_TELEMATICS) {
		return nil, errors.New("curfew feature requires G2 telematics")
	}

	if !slices.Contains(v.SubscriptionFeatures, FEATURE_SAFETY) {
		return nil, errors.New("curfew feature requires Safety Plus subscription")
	}

	// Validate input parameters
	if err := ValidateTimeRange(startTime, endTime); err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	if err := ValidateDaysOfWeek(daysOfWeek); err != nil {
		return nil, fmt.Errorf("invalid days of week: %w", err)
	}

	// Convert daysOfWeek array to comma-separated string
	daysStr := ""
	for i, day := range daysOfWeek {
		if i > 0 {
			daysStr += ","
		}
		daysStr += strconv.Itoa(day)
	}

	params := map[string]string{
		"delay":      "0",
		"vin":        v.Vin,
		"pin":        v.client.credentials.PIN,
		"startTime":  startTime,
		"endTime":    endTime,
		"daysOfWeek": daysStr,
		"enabled":    strconv.FormatBool(enabled),
	}
	reqUrl := MOBILE_API_VERSION + apiURLs["API_G2_CURFEW"]
	pollingUrl := MOBILE_API_VERSION + apiURLs["API_REMOTE_SVC_STATUS"]

	ch := make(chan string)
	go func() {
		defer close(ch)
		v.executeServiceRequest(params, reqUrl, pollingUrl, ch, 1)
	}()

	return ch, nil
}

// validateSubscriptionAndSession checks if the vehicle has remote options and validates the session
func (v *Vehicle) validateSubscriptionAndSession() error {
	if !v.getRemoteOptionsStatus() {
		v.client.logger.Error(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
		return errors.New(APP_ERRORS["SUBSCRIPTION_REQUIRED"])
	}

	if !v.client.validateSession() {
		v.client.logger.Error(APP_ERRORS["SESSION_EXPIRED"])
		return errors.New(APP_ERRORS["SESSION_EXPIRED"])
	}

	return nil
}

// ensureVehicleSelected ensures the current vehicle is selected in the client
func (v *Vehicle) ensureVehicleSelected() {
	if v.Vin != v.client.currentVin {
		v.selectVehicle()
	}
}

// processClimateProfiles processes and stores climate profiles based on vehicle type
func (v *Vehicle) processClimateProfiles(cProfiles []ClimateProfile) {
	if len(cProfiles) == 0 {
		v.client.logger.Debug("couldn't find any climate presets")
		return
	}

	if v.ClimateProfiles == nil {
		v.ClimateProfiles = make(map[string]ClimateProfile)
	}

	for _, cp := range cProfiles {
		re := regexp.MustCompile(`([A-Z])`)
		cpn := strings.ToLower(re.ReplaceAllString(cp.PresetType, "_$1") + "_" + strings.ReplaceAll(cp.Name, " ", "_"))

		// Always keep user presets; be permissive on vehicleType to avoid dropping valid entries.
		allowed := cp.VehicleType == ""
		if v.IsEV() {
			allowed = allowed || cp.VehicleType == "phev"
		} else {
			allowed = allowed || cp.VehicleType == "gas"
		}
		if strings.EqualFold(cp.PresetType, "userPreset") {
			allowed = true
		}
		if allowed {
			v.ClimateProfiles[cpn] = cp
		}
	}
}

// parseClimateData parses the climate data from the API response.
func (v *Vehicle) parseClimateData(data string) string {
	re1 := regexp.MustCompile(`\"`)
	result := re1.ReplaceAllString(data, "")
	re2 := regexp.MustCompile(`\\`)
	result = re2.ReplaceAllString(result, `"`)
	return result
}

// DetectModelFromCode detects the Subaru model and trim from the model code.
func DetectModelFromCode(modelCode string) (model string, trim string) {
	if len(modelCode) < 3 {
		return "Unknown", "Unknown"
	}

	// Convert to uppercase for consistency
	code := strings.ToUpper(modelCode)

	// Extract model indicator (first character)
	modelIndicator := string(code[0])

	// Extract trim indicator (second and third characters combined)
	trimIndicator := string(code[1:3])

	// Model mapping based on first character
	modelMap := map[string]string{
		"P": "Outback",
		"S": "Forester",
		"L": "Legacy",
		"C": "Crosstrek",
		"K": "Crosstrek", // Older model code for Crosstrek (pre-2023)
		"A": "Ascent",
		"W": "WRX",
		"I": "Impreza",
		"B": "BRZ",
		"T": "Solterra",
	}

	model = modelMap[modelIndicator]
	if model == "" {
		model = "Unknown"
	}

	// Trim mapping based on model and second/third characters
	trimMap := map[string]map[string]string{
		"Outback": {
			"DL": "Limited XT",
			"FL": "Limited",
			"CL": "Convenience",
			"BL": "Base",
			"DH": "Limited XT Hybrid",
			"FH": "Limited Hybrid",
			"CH": "Convenience Hybrid",
			"BH": "Base Hybrid",
		},
		"Forester": {
			"DL": "Limited",
			"FL": "Premier",
			"CL": "Convenience",
			"BL": "Base",
			"DH": "Limited Hybrid",
			"FH": "Premier Hybrid",
			"CH": "Convenience Hybrid",
			"BH": "Base Hybrid",
		},
		"Legacy": {
			"DL": "Limited",
			"FL": "Limited XT",
			"CL": "Convenience",
			"BL": "Base",
		},
		"Crosstrek": {
			"DL": "Limited",
			"FL": "Premier",
			"CL": "Convenience",
			"BL": "Base",
			"DH": "Limited Hybrid",
			"FH": "Premier Hybrid",
			"CH": "Convenience Hybrid",
			"BH": "Base Hybrid",
			"RH": "Convenience", // Older trim code for 2017+ Crosstrek
		},
		"Ascent": {
			"DL": "Limited",
			"FL": "Premier",
			"CL": "Convenience",
			"BL": "Base",
		},
		"WRX": {
			"DL": "Limited",
			"FL": "STI",
			"CL": "Convenience",
			"BL": "Base",
		},
		"Impreza": {
			"DL": "Limited",
			"FL": "Premier",
			"CL": "Convenience",
			"BL": "Base",
		},
		"BRZ": {
			"DL": "Limited",
			"FL": "STI",
			"CL": "Convenience",
			"BL": "Base",
		},
		"Solterra": {
			"DL": "Limited",
			"FL": "Premier",
			"CL": "Convenience",
			"BL": "Base",
		},
	}

	if modelTrims, exists := trimMap[model]; exists {
		if trimName, exists := modelTrims[trimIndicator]; exists {
			trim = trimName
		} else {
			trim = "Unknown"
		}
	} else {
		trim = "Unknown"
	}

	return model, trim
}

// GetModelInfo returns the detected model name and trim level based on the model code.
func (v *Vehicle) GetModelInfo() (model string, trim string) {
	return DetectModelFromCode(v.ModelCode)
}
