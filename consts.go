// Package mysubaru provides constants and configuration for the MySubaru API client.
// This package contains API endpoints, error codes, feature mappings, and other
// constants used throughout the MySubaru Go client library.
package mysubaru

// API Configuration
var MOBILE_API_VERSION = "/g2v31"

var MOBILE_API_SERVER = map[string]string{
	"USA":  "https://mobileapi.prod.subarucs.com",
	"CAN":  "https://mobileapi.ca.prod.subarucs.com",
	"TEST": "http://127.0.0.1:56765",
}

var MOBILE_APP = map[string]string{
	"USA": "com.subaru.telematics.app.remote",
	"CAN": "ca.subaru.telematics.remote",
}

var WEB_API_SERVER = map[string]string{
	"USA": "https://www.mysubaru.com",
	"CAN": "https://www.mysubaru.ca",
}

// API Endpoints
var apiURLs = map[string]string{
	// Web API endpoints
	"WEB_API_LOGIN":              "/login",
	"WEB_API_LIST_DEVICES":       "/listMyDevices.json",
	"WEB_API_AUTHORIZE_DEVICE":   "/profile/updateDeviceEntry.json",
	"WEB_API_NAME_DEVICE":        "/profile/addDeviceName.json",
	"WEB_API_EDIT_NAME_DEVICE":   "/profile/editDeviceName.json",
	"WEB_API_VERIFY_NAME_DEVICE": "/profile/verifyDeviceName.json",

	// Authentication endpoints
	"API_2FA_CONTACT":           "/twoStepAuthContacts.json",
	"API_2FA_SEND_VERIFICATION": "/twoStepAuthSendVerification.json",
	"API_2FA_AUTH_VERIFY":       "/twoStepAuthVerify.json",
	"API_LOGIN":                 "/login.json",
	"API_REFRESH_VEHICLES":      "/refreshVehicles.json",
	"API_SELECT_VEHICLE":        "/selectVehicle.json",
	"API_VALIDATE_SESSION":      "/validateSession.json",

	// Device management
	"API_AUTHORIZE_DEVICE": "/authenticateDevice.json",
	"API_NAME_DEVICE":      "/nameThisDevice.json",

	// Vehicle data endpoints
	"API_VEHICLE_STATUS": "/vehicleStatus.json",
	"API_VEHICLE_HEALTH": "/vehicleHealth.json",
	"API_CONDITION":      "/service/api_gen/condition/execute.json",
	"API_LOCATE":         "/service/api_gen/locate/execute.json",

	// Remote service endpoints
	"API_LOCK":               "/service/api_gen/lock/execute.json",
	"API_LOCK_CANCEL":        "/service/api_gen/lock/cancel.json",
	"API_UNLOCK":             "/service/api_gen/unlock/execute.json",
	"API_UNLOCK_CANCEL":      "/service/api_gen/unlock/cancel.json",
	"API_HORN_LIGHTS":        "/service/api_gen/hornLights/execute.json",
	"API_HORN_LIGHTS_CANCEL": "/service/api_gen/hornLights/cancel.json",
	"API_HORN_LIGHTS_STOP":   "/service/api_gen/hornLights/stop.json",
	"API_LIGHTS":             "/service/api_gen/lightsOnly/execute.json",
	"API_LIGHTS_CANCEL":      "/service/api_gen/lightsOnly/cancel.json",
	"API_LIGHTS_STOP":        "/service/api_gen/lightsOnly/stop.json",

	// Generation-specific endpoints
	"API_G1_LOCATE_UPDATE":      "/service/g1/vehicleLocate/execute.json",
	"API_G1_LOCATE_STATUS":      "/service/g1/vehicleLocate/status.json",
	"API_G1_HORN_LIGHTS_STATUS": "/service/g1/hornLights/status.json",
	"API_G2_LOCATE_UPDATE":      "/service/g2/vehicleStatus/execute.json",
	"API_G2_LOCATE_STATUS":      "/service/g2/vehicleStatus/locationStatus.json",
	"API_REMOTE_SVC_STATUS":     "/service/g2/remoteService/status.json",
	"API_G2_SEND_POI":           "/service/g2/sendPoi/execute.json",

	// Advanced features
	"API_G2_SPEEDFENCE": "/service/g2/speedFence/execute.json",
	"API_G2_GEOFENCE":   "/service/g2/geoFence/execute.json",
	"API_G2_CURFEW":     "/service/g2/curfew/execute.json",

	// Remote engine start
	"API_G2_REMOTE_ENGINE_START":        "/service/g2/engineStart/execute.json",
	"API_G2_REMOTE_ENGINE_START_CANCEL": "/service/g2/engineStart/cancel.json",
	"API_G2_REMOTE_ENGINE_STOP":         "/service/g2/engineStop/execute.json",

	// Climate control
	"API_G2_FETCH_RES_QUICK_START_SETTINGS": "/service/g2/remoteEngineQuickStartSettings/fetch.json",
	"API_G2_FETCH_RES_USER_PRESETS":         "/service/g2/remoteEngineStartSettings/fetch.json",
	"API_G2_FETCH_RES_SUBARU_PRESETS":       "/service/g2/climatePresetSettings/fetch.json",
	"API_G2_SAVE_RES_SETTINGS":              "/service/g2/remoteEngineStartSettings/save.json",
	"API_G2_SAVE_RES_QUICK_START_SETTINGS":  "/service/g2/remoteEngineQuickStartSettings/save.json",

	// EV-specific endpoints
	"API_EV_CHARGE_NOW":             "/service/g2/phevChargeNow/execute.json",
	"API_EV_FETCH_CHARGE_SETTINGS":  "/service/g2/phevGetTimerSettings/execute.json",
	"API_EV_SAVE_CHARGE_SETTINGS":   "/service/g2/phevSendTimerSetting/execute.json",
	"API_EV_DELETE_CHARGE_SCHEDULE": "/service/g2/phevDeleteTimerSetting/execute.json",
}

// Error codes and messages
var API_ERRORS = map[string]string{
	// G2 API errors
	"API_ERROR_SOA_403":                 "403-soa-unableToParseResponseBody",
	"API_ERROR_NO_ACCOUNT":              "accountNotFound",
	"API_ERROR_INVALID_ACCOUNT":         "invalidAccount",
	"API_ERROR_INVALID_CREDENTIALS":     "InvalidCredentials",
	"API_ERROR_INVALID_TOKEN":           "InvalidToken",
	"API_ERROR_PASSWORD_WARNING":        "passwordWarning",
	"API_ERROR_TOO_MANY_ATTEMPTS":       "tooManyAttempts",
	"API_ERROR_ACCOUNT_LOCKED":          "accountLocked",
	"API_ERROR_NO_VEHICLES":             "noVehiclesOnAccount",
	"API_ERROR_VEHICLE_SETUP":           "VEHICLESETUPERROR",
	"API_ERROR_VEHICLE_NOT_IN_ACCOUNT":  "vehicleNotInAccount",
	"API_ERROR_SERVICE_ALREADY_STARTED": "ServiceAlreadyStarted",

	// G1 API errors
	"API_ERROR_G1_NO_SUBSCRIPTION":         "SXM40004",
	"API_ERROR_G1_STOLEN_VEHICLE":          "SXM40005",
	"API_ERROR_G1_INVALID_PIN":             "SXM40006",
	"API_ERROR_G1_SERVICE_ALREADY_STARTED": "SXM40009",
	"API_ERROR_G1_PIN_LOCKED":              "SXM40017",
}

var APP_ERRORS = map[string]string{
	"SUBSCRIPTION_REQUIRED": "active STARLINK Security Plus subscription required",
}

// Vehicle features mapping
var features = map[string]string{
	// Telematics generations
	"g1": "Generation #1",
	"g2": "Generation #2",
	"g3": "Generation #3",

	// Safety and driver assistance
	"BSD":      "Blind-Spot Detection",
	"RHSF":     "Rear High-Speed Function / Reverse Automatic Braking / Rear Cross-Traffic Alert",
	"EYESIGHT": "EyeSight Exclusive Advanced Driver-Assist System",
	"ACCS":     "Adaptive Cruise Control",
	"REARBRK":  "Reverse Auto Braking",

	// Infotainment and connectivity
	"11.6MMAN":   "11.6-inch Infotainment System",
	"NAV_TOMTOM": "TomTom Navigation",
	"SXM360L":    "SiriusXM with 360L",

	// Comfort and convenience
	"PWAAADWWAP":     "Power Windows",
	"PANPM-TUIRWAOC": "Power Moonroof",
	"PANPM-DG2G":     "Panoramic Moonroof",
	"WDWSTAT":        "Window Status",
	"MOONSTAT":       "Moonroof Status",
	"RTGU":           "Remote Trunk / Rear Gate Unlock",
	"RVFS":           "Remote Vehicle Find System",
	"DOOR_LU_STAT":   "Door Lock/Unlock Status",

	// Climate and heating
	"RES":          "Remote Engine Start",
	"RESCC":        "Remote Engine Start with Climate Control",
	"RCC":          "Remote Climate Control",
	"RES_HVAC_HFS": "Heated Front Seats",
	"RES_HVAC_VFS": "Vented Front Seats",

	// EV features
	"PHEV": "Electric Vehicle",

	// Other features
	"VALET":  "Valet Parking",
	"TIF_35": "Tire Pressure Front 35",
	"TIR_33": "Tire Pressure Rear 35",
	"TLD":    "Tire Pressure Low Detection",
	"RPOI":   "Remote Geo Point of Interest",
}

// Vehicle trouble codes mapping
var troubles = map[string]string{
	"ABS_MIL":    "Anti-Lock Braking System",
	"AHBL_MIL":   "Automatic Headlight Beam Leveler",
	"ATF_MIL":    "Automatic Transmission Oil Temperature",
	"AWD_MIL":    "All-Wheel Drive / Symmetrical Full-Time",
	"BSDRCT_MIL": "Blind-Spot Detection",
	"CEL_MIL":    "Check Engine Light",
	"EBD_MIL":    "Brake System / Electronic Brake Force Distribution",
	"EOL_MIL":    "Engine Oil Level",
	"EPAS_MIL":   "Power Steering / Electric Power Assisted Steering",
	"EPB_MIL":    "Parking Brake",
	"ESS_MIL":    "EyeSight Exclusive Advanced Driver-Assist System",
	"HEV_MIL":    "Hybrid System",
	"HEVCM_MIL":  "Hybrid Charge System",
	"ISS_MIL":    "Auto Start Stop (Idling Stop System)",
	"OPL_MIL":    "Oil Pressure",
	"RAB_MIL":    "Reverse Auto Braking",
	"SRH_MIL":    "Steering Responsive Headlights (SRH)",
	"SRS_MIL":    "Airbag System",
	"TEL_MIL":    "MySubaru Emergency Services",
	"TPMS_MIL":   "Tire Pressure",
	"VDC_MIL":    "Vehicle Dynamics Control",
	"WASH_MIL":   "Windshield Washer Fluid Level",
}

// Invalid/erroneous values that should be ignored
var badValues = []any{
	"NOT_EQUIPPED",
	"UNKNOWN",
	"unknown",
	"None",
	"16383",
	"65535",
	"-64",
	"",
	0,
	float64(0),
	nil,
}

// HTTP methods and API constants
const (
	GET  = "GET"
	POST = "POST"

	// Service request constants
	SERVICE_REQ_ID = "serviceRequestId"

	// Temperature constants (Fahrenheit)
	TEMP_F_MIN = 60
	TEMP_F_MAX = 85
	TEMP_F     = "climateZoneFrontTemp"

	// Temperature constants (Celsius)
	TEMP_C_MIN = 15
	TEMP_C_MAX = 30
	TEMP_C     = "climateZoneFrontTempCelsius"

	// Climate control constants
	CLIMATE         = "climateSettings"
	CLIMATE_DEFAULT = "climateSettings"
	RUNTIME         = "runTimeMinutes"
	RUNTIME_DEFAULT = "10"

	// Air mode options
	MODE              = "climateZoneFrontAirMode"
	MODE_DEFROST      = "WINDOW"
	MODE_FEET_DEFROST = "FEET_WINDOW"
	MODE_FACE         = "FACE"
	MODE_FEET         = "FEET"
	MODE_SPLIT        = "FEET_FACE_BALANCED"
	MODE_AUTO         = "AUTO"

	// Seat heating options
	HEAT_SEAT_LEFT  = "heatedSeatFrontLeft"
	HEAT_SEAT_RIGHT = "heatedSeatFrontRight"
	HEAT_SEAT_HI    = "HIGH_HEAT"
	HEAT_SEAT_MED   = "MEDIUM_HEAT"
	HEAT_SEAT_LOW   = "LOW_HEAT"
	HEAT_SEAT_OFF   = "OFF"

	// Rear defrost options
	REAR_DEFROST     = "heatedRearWindowActive"
	REAR_DEFROST_ON  = "true"
	REAR_DEFROST_OFF = "false"

	// Fan speed options
	FAN_SPEED      = "climateZoneFrontAirVolume"
	FAN_SPEED_LOW  = "2"
	FAN_SPEED_MED  = "4"
	FAN_SPEED_HI   = "7"
	FAN_SPEED_AUTO = "AUTO"

	// Air circulation options
	RECIRCULATE     = "outerAirCirculation"
	RECIRCULATE_OFF = "outsideAir"
	RECIRCULATE_ON  = "recirculation"

	// Rear AC options
	REAR_AC     = "airConditionOn"
	REAR_AC_ON  = "true"
	REAR_AC_OFF = "false"

	// Start configuration options
	START_CONFIG             = "startConfiguration"
	START_CONFIG_DEFAULT_EV  = "START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION"
	START_CONFIG_DEFAULT_RES = "START_ENGINE_ALLOW_KEY_IN_IGNITION"

	// Door unlock options
	WHICH_DOOR   = "unlockDoorType"
	ALL_DOORS    = "ALL_DOORS_CMD"
	DRIVERS_DOOR = "FRONT_LEFT_DOOR_CMD"

	// Location data constants
	HEADING       = "heading"
	LATITUDE      = "latitude"
	LONGITUDE     = "longitude"
	LOCATION_TIME = "locationTimestamp"
	SPEED         = "speed"
	BAD_LATITUDE  = 90.0
	BAD_LONGITUDE = 180.0

	// Vehicle status constants
	AVG_FUEL_CONSUMPTION = "AVG_FUEL_CONSUMPTION"
	BATTERY_VOLTAGE      = "BATTERY_VOLTAGE"
	DIST_TO_EMPTY        = "DISTANCE_TO_EMPTY_FUEL"

	// Door position constants
	DOOR_BOOT_POSITION        = "DOOR_BOOT_POSITION"
	DOOR_ENGINE_HOOD_POSITION = "DOOR_ENGINE_HOOD_POSITION"
	DOOR_FRONT_LEFT_POSITION  = "DOOR_FRONT_LEFT_POSITION"
	DOOR_FRONT_RIGHT_POSITION = "DOOR_FRONT_RIGHT_POSITION"
	DOOR_REAR_LEFT_POSITION   = "DOOR_REAR_LEFT_POSITION"
	DOOR_REAR_RIGHT_POSITION  = "DOOR_REAR_RIGHT_POSITION"

	// Door lock status constants
	DOOR_BOOT_LOCK_STATUS        = "DOOR_BOOT_LOCK_STATUS"
	DOOR_FRONT_LEFT_LOCK_STATUS  = "DOOR_FRONT_LEFT_LOCK_STATUS"
	DOOR_FRONT_RIGHT_LOCK_STATUS = "DOOR_FRONT_RIGHT_LOCK_STATUS"
	DOOR_REAR_LEFT_LOCK_STATUS   = "DOOR_REAR_LEFT_LOCK_STATUS"
	DOOR_REAR_RIGHT_LOCK_STATUS  = "DOOR_REAR_RIGHT_LOCK_STATUS"

	// EV-specific constants
	EV_CHARGER_STATE_TYPE         = "EV_CHARGER_STATE_TYPE"
	EV_CHARGE_SETTING_AMPERE_TYPE = "EV_CHARGE_SETTING_AMPERE_TYPE"
	EV_CHARGE_VOLT_TYPE           = "EV_CHARGE_VOLT_TYPE"
	EV_DISTANCE_TO_EMPTY          = "EV_DISTANCE_TO_EMPTY"
	EV_IS_PLUGGED_IN              = "EV_IS_PLUGGED_IN"
	EV_STATE_OF_CHARGE_MODE       = "EV_STATE_OF_CHARGE_MODE"
	EV_STATE_OF_CHARGE_PERCENT    = "EV_STATE_OF_CHARGE_PERCENT"
	EV_TIME_TO_FULLY_CHARGED      = "EV_TIME_TO_FULLY_CHARGED"
	EV_TIME_TO_FULLY_CHARGED_UTC  = "EV_TIME_TO_FULLY_CHARGED_UTC"

	// General vehicle constants
	EXTERNAL_TEMP      = "EXT_EXTERNAL_TEMP"
	ODOMETER           = "ODOMETER"
	POSITION_TIMESTAMP = "POSITION_TIMESTAMP"
	TIMESTAMP          = "TIMESTAMP"

	// Tire pressure constants
	TIRE_PRESSURE_FL = "TYRE_PRESSURE_FRONT_LEFT"
	TIRE_PRESSURE_FR = "TYRE_PRESSURE_FRONT_RIGHT"
	TIRE_PRESSURE_RL = "TYRE_PRESSURE_REAR_LEFT"
	TIRE_PRESSURE_RR = "TYRE_PRESSURE_REAR_RIGHT"

	// Vehicle state constants
	VEHICLE_STATE = "VEHICLE_STATE_TYPE"

	// Window status constants
	WINDOW_FRONT_LEFT_STATUS  = "WINDOW_FRONT_LEFT_STATUS"
	WINDOW_FRONT_RIGHT_STATUS = "WINDOW_FRONT_RIGHT_STATUS"
	WINDOW_REAR_LEFT_STATUS   = "WINDOW_REAR_LEFT_STATUS"
	WINDOW_REAR_RIGHT_STATUS  = "WINDOW_REAR_RIGHT_STATUS"
	WINDOW_SUNROOF_STATUS     = "WINDOW_SUNROOF_STATUS"

	// Vehicle state values
	CHARGING           = "CHARGING"
	LOCKED_CONNECTED   = "LOCKED_CONNECTED"
	UNLOCKED_CONNECTED = "UNLOCKED_CONNECTED"
	DOOR_OPEN          = "OPEN"
	DOOR_CLOSED        = "CLOSED"
	WINDOW_OPEN        = "OPEN"
	WINDOW_CLOSED      = "CLOSE"
	IGNITION_ON        = "IGNITION_ON"
	NOT_EQUIPPED       = "NOT_EQUIPPED"

	// Vehicle status JSON field mappings
	VS_AVG_FUEL_CONSUMPTION = "avgFuelConsumptionLitersPer100Kilometers"
	VS_DIST_TO_EMPTY        = "distanceToEmptyFuelKilometers"
	VS_TIMESTAMP            = "eventDate"
	VS_LATITUDE             = "latitude"
	VS_LONGITUDE            = "longitude"
	VS_HEADING              = "positionHeadingDegree"
	VS_ODOMETER             = "odometerValueKilometers"
	VS_VEHICLE_STATE        = "vehicleStateType"
	VS_TIRE_PRESSURE_FL     = "tirePressureFrontLeft"
	VS_TIRE_PRESSURE_FR     = "tirePressureFrontRight"
	VS_TIRE_PRESSURE_RL     = "tirePressureRearLeft"
	VS_TIRE_PRESSURE_RR     = "tirePressureRearRight"

	// Erroneous sensor values
	BAD_AVG_FUEL_CONSUMPTION     = "16383"
	BAD_DISTANCE_TO_EMPTY_FUEL   = "16383"
	BAD_EV_TIME_TO_FULLY_CHARGED = "65535"
	BAD_TIRE_PRESSURE            = "32767"
	BAD_ODOMETER                 = "None"
	BAD_EXTERNAL_TEMP            = "-64.0"
	UNKNOWN                      = "UNKNOWN"
	VENTED                       = "VENTED"
	LOCATION_VALID               = "location_valid"

	// Timestamp formats
	TIMESTAMP_FMT          = "%Y-%m-%dT%H:%M:%S%z"
	POSITION_TIMESTAMP_FMT = "%Y-%m-%dT%H:%M:%SZ"

	// Error codes
	ERROR_SOA_403                    = "403-soa-unableToParseResponseBody"
	ERROR_SOA_404                    = "404-soa-unableToParseResponseBody"
	ERROR_INVALID_CREDENTIALS        = "InvalidCredentials"
	ERROR_SERVICE_ALREADY_STARTED    = "ServiceAlreadyStarted"
	ERROR_INVALID_ACCOUNT            = "invalidAccount"
	ERROR_PASSWORD_WARNING           = "passwordWarning"
	ERROR_ACCOUNT_LOCKED             = "accountLocked"
	ERROR_NO_VEHICLES                = "noVehiclesOnAccount"
	ERROR_NO_ACCOUNT                 = "accountNotFound"
	ERROR_TOO_MANY_ATTEMPTS          = "tooManyAttempts"
	ERROR_VEHICLE_NOT_IN_ACCOUNT     = "vehicleNotInAccount"
	ERROR_G1_NO_SUBSCRIPTION         = "SXM40004"
	ERROR_G1_STOLEN_VEHICLE          = "SXM40005"
	ERROR_G1_INVALID_PIN             = "SXM40006"
	ERROR_G1_SERVICE_ALREADY_STARTED = "SXM40009"
	ERROR_G1_PIN_LOCKED              = "SXM40017"

	// Vehicle data dictionary keys
	VEHICLE_ATTRIBUTES            = "attributes"
	VEHICLE_STATUS                = "status"
	VEHICLE_ID                    = "id"
	VEHICLE_NAME                  = "nickname"
	VEHICLE_API_GEN               = "api_gen"
	VEHICLE_LOCK                  = "lock"
	VEHICLE_LAST_UPDATE           = "last_update_time"
	VEHICLE_LAST_FETCH            = "last_fetch_time"
	VEHICLE_FEATURES              = "features"
	VEHICLE_SUBSCRIPTION_FEATURES = "subscriptionFeatures"
	VEHICLE_SUBSCRIPTION_STATUS   = "subscriptionStatus"

	// Vehicle feature constants
	FEATURE_PHEV          = "PHEV"
	FEATURE_REMOTE_START  = "RES"
	FEATURE_G1_TELEMATICS = "g1"
	FEATURE_G2_TELEMATICS = "g2"
	FEATURE_G3_TELEMATICS = "g3"
	FEATURE_REMOTE        = "REMOTE"
	FEATURE_SAFETY        = "SAFETY"
	FEATURE_ACTIVE        = "ACTIVE"

	// Update intervals (in seconds)
	DEFAULT_UPDATE_INTERVAL = 7200
	DEFAULT_FETCH_INTERVAL  = 300
)
