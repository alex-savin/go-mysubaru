package mysubaru

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Custom error types for better error handling
type APIError struct {
	Code      string
	Message   string
	Retryable bool
}

func (e APIError) Error() string {
	return fmt.Sprintf("MySubaru API error [%s]: %s", e.Code, e.Message)
}

func (e APIError) IsRetryable() bool {
	return e.Retryable
}

// Common API errors
var (
	ErrInvalidCredentials   = APIError{Code: "INVALID_CREDENTIALS", Message: "Invalid username or password", Retryable: false}
	ErrAccountLocked        = APIError{Code: "ACCOUNT_LOCKED", Message: "Account is locked", Retryable: false}
	ErrSessionExpired       = APIError{Code: "SESSION_EXPIRED", Message: "Session has expired", Retryable: true}
	ErrRateLimited          = APIError{Code: "RATE_LIMITED", Message: "Too many requests", Retryable: true}
	ErrNetworkError         = APIError{Code: "NETWORK_ERROR", Message: "Network communication failed", Retryable: true}
	ErrSubscriptionRequired = APIError{Code: "SUBSCRIPTION_REQUIRED", Message: "Required subscription not active", Retryable: false}
)

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	var apiErr APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}
	return false
}

// Response represents the structure of a response from the MySubaru API.
type Response struct {
	Success   bool            `json:"success"`             // true | false
	ErrorCode string          `json:"errorCode,omitempty"` // string | Error message if Success is false
	DataName  string          `json:"dataName,omitempty"`  // string | Describes the structure which is included in Data field
	Data      json.RawMessage `json:"data"`                // Data struct
}

// apiErrorMessages maps API error codes to human-readable error messages.
var apiErrorMessages = map[string]string{
	// G2 API errors
	"API_ERROR_NO_ACCOUNT":              "Account not found",
	"API_ERROR_INVALID_ACCOUNT":         "Invalid Account",
	"API_ERROR_INVALID_CREDENTIALS":     "Invalid Credentials",
	"API_ERROR_INVALID_TOKEN":           "Invalid Token",
	"API_ERROR_PASSWORD_WARNING":        "Mutiple failed login attempts, password warning",
	"API_ERROR_TOO_MANY_ATTEMPTS":       "Too many attempts, please try again later",
	"API_ERROR_ACCOUNT_LOCKED":          "Account Locked",
	"API_ERROR_NO_VEHICLES":             "No vehicles found for the account",
	"API_ERROR_VEHICLE_SETUP":           "Vehicle setup is not complete",
	"API_ERROR_VEHICLE_NOT_IN_ACCOUNT":  "Vehicle not in account",
	"API_ERROR_SERVICE_ALREADY_STARTED": "Service already started",
	"API_ERROR_SOA_403":                 "Unable to parse response body, SOA 403 error",
	// G1 API errors
	"API_ERROR_G1_NO_SUBSCRIPTION":         "No subscription found for the vehicle",
	"API_ERROR_G1_STOLEN_VEHICLE":          "Car is reported as stolen",
	"API_ERROR_G1_INVALID_PIN":             "Invalid PIN",
	"API_ERROR_G1_PIN_LOCKED":              "PIN is locked",
	"API_ERROR_G1_SERVICE_ALREADY_STARTED": "Service already started",
}

// parse parses the JSON response from the MySubaru API into a Response struct.
func (r *Response) parse(b []byte, logger *slog.Logger) (*Response, error) {
	err := json.Unmarshal(b, &r)
	if err != nil {
		logger.Error("error while parsing json", "error", err.Error())
		return nil, errors.New("error while parsing json: " + err.Error())
	}

	if !r.Success && r.ErrorCode != "" {
		logger.Error("error in response", "errorCode", r.ErrorCode, "dataName", r.DataName)
		if msg := getAPIErrorMessage(r.ErrorCode); msg != "" {
			return r, errors.New("error in response: " + msg)
		}
		return r, errors.New("error in response: " + r.ErrorCode)
	}

	return r, nil
}

// getAPIErrorMessage looks up the error message for a given API error code.
func getAPIErrorMessage(errorCode string) string {
	for key, msg := range apiErrorMessages {
		if API_ERRORS[key] == errorCode {
			return msg
		}
	}
	return ""
}

// Request represents the structure of a request to the MySubaru API.
type Request struct {
	Vin                       string  `json:"vin"`                                 //
	Pin                       string  `json:"pin"`                                 //
	Delay                     int     `json:"delay,string,omitempty"`              //
	ForceKeyInCar             *bool   `json:"forceKeyInCar,string,omitempty"`      //
	UnlockDoorType            *string `json:"unlockDoorType,omitempty"`            // [ ALL_DOORS_CMD | FRONT_LEFT_DOOR_CMD | ALL_DOORS_CMD ]
	Horn                      *string `json:"horn,omitempty"`                      //
	ClimateSettings           *string `json:"climateSettings,omitempty"`           //
	ClimateZoneFrontTemp      *string `json:"climateZoneFrontTemp,omitempty"`      //
	ClimateZoneFrontAirMode   *string `json:"climateZoneFrontAirMode,omitempty"`   //
	ClimateZoneFrontAirVolume *string `json:"climateZoneFrontAirVolume,omitempty"` //
	HeatedSeatFrontLeft       *string `json:"heatedSeatFrontLeft,omitempty"`       //
	HeatedSeatFrontRight      *string `json:"heatedSeatFrontRight,omitempty"`      //
	HeatedRearWindowActive    *string `json:"heatedRearWindowActive,omitempty"`    //
	OuterAirCirculation       *string `json:"outerAirCirculation,omitempty"`       //
	AirConditionOn            *string `json:"airConditionOn,omitempty"`            //
	RunTimeMinutes            *string `json:"runTimeMinutes,omitempty"`            //
	StartConfiguration        *string `json:"startConfiguration,omitempty"`        //
}

// account .
type account struct {
	MarketID      int      `json:"marketId"`
	AccountKey    int      `json:"accountKey"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	ZipCode       string   `json:"zipCode"`
	ZipCode5      string   `json:"zipCode5"`
	LastLoginDate UnixTime `json:"lastLoginDate"`
	CreatedDate   UnixTime `json:"createdDate"`
}

// Customer .
type Customer struct {
	SessionCustomer SessionCustomer `json:"sessionCustomer,omitempty"` // struct | Only by performing a RefreshVehicles request
	Email           string          `json:"email"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	Zip             string          `json:"zip"`
	OemCustID       string          `json:"oemCustId"`
	Phone           string          `json:"phone"`
}

// SessionCustomer .
type SessionCustomer struct {
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Title            string `json:"title,omitempty"`
	Suffix           string `json:"suffix,omitempty"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	Address2         string `json:"address2,omitempty"`
	City             string `json:"city"`
	State            string `json:"state"`
	Zip              string `json:"zip"`
	CellularPhone    string `json:"cellularPhone,omitempty"`
	WorkPhone        string `json:"workPhone,omitempty"`
	HomePhone        string `json:"homePhone,omitempty"`
	CountryCode      string `json:"countryCode"`
	RelationshipType any    `json:"relationshipType,omitempty"`
	Gender           string `json:"gender,omitempty"`
	DealerCode       any    `json:"dealerCode,omitempty"`
	OemCustID        string `json:"oemCustId"`
	CreateMysAccount any    `json:"createMysAccount,omitempty"`
	SourceSystemCode string `json:"sourceSystemCode"`
	Vehicles         []struct {
		Vin                       string `json:"vin"`
		SiebelVehicleRelationship string `json:"siebelVehicleRelationship"` // TM Subscriber | Previous TM Subscriber | Previous Owner
		Primary                   bool   `json:"primary"`                   // true | false
		OemCustID                 string `json:"oemCustId"`                 // CRM-41PLM-5TYE | 1-8K7OBOJ | 1-8JY3UVS | CRM-44UFUA14-V
		Status                    string `json:"status,omitempty"`          // "Active" | "Draft" | "Inactive"
	} `json:"vehicles"`
	Phone                  string `json:"phone,omitempty"`
	Zip5Digits             string `json:"zip5Digits"`
	PrimaryPersonalCountry string `json:"primaryPersonalCountry"`
}

// DataMap .
// "dataName": "dataMap"
type dataMap struct {
	Username string `json:"userName"`
	Email    string `json:"email"`
}

// SessionData .
// "dataName": "sessionData"
type SessionData struct {
	Account                            account       `json:"account"`
	PasswordToken                      string        `json:"passwordToken"`
	ResetPassword                      bool          `json:"resetPassword"`
	SessionID                          string        `json:"sessionId"`
	SessionChanged                     bool          `json:"sessionChanged"`
	DeviceID                           string        `json:"deviceId"`
	DeviceRegistered                   bool          `json:"deviceRegistered"`
	RegisteredDevicePermanent          bool          `json:"registeredDevicePermanent"`
	Vehicles                           []VehicleData `json:"vehicles"`
	VehicleInactivated                 bool          `json:"vehicleInactivated"`
	RightToRepairEnabled               bool          `json:"rightToRepairEnabled"`
	RightToRepairStates                string        `json:"rightToRepairStates"`
	CurrentVehicleIndex                int           `json:"currentVehicleIndex"`
	HandoffToken                       string        `json:"handoffToken"`
	EnableXtime                        bool          `json:"enableXtime"`
	TermsAndConditionsAccepted         bool          `json:"termsAndConditionsAccepted"`
	RightToRepairStartYear             int           `json:"rightToRepairStartYear"`
	DigitalGlobeConnectID              string        `json:"digitalGlobeConnectId"`
	DigitalGlobeImageTileService       string        `json:"digitalGlobeImageTileService"`
	DigitalGlobeTransparentTileService string        `json:"digitalGlobeTransparentTileService"`
	TomtomKey                          string        `json:"tomtomKey"`
	SatelliteViewEnabled               bool          `json:"satelliteViewEnabled"`
}

// Vehicle .
// "dataName": "vehicle"
type VehicleData struct {
	Customer                   Customer    `json:"customer"`                    // Customer struct
	OemCustID                  string      `json:"oemCustId"`                   // CRM-631-HQN48K
	UserOemCustID              string      `json:"userOemCustId"`               // CRM-631-HQN48K
	Active                     bool        `json:"active"`                      // true | false
	Email                      string      `json:"email"`                       // null | email@address.com
	FirstName                  string      `json:"firstName,omitempty"`         // null | First Name
	LastName                   string      `json:"lastName,omitempty"`          // null | Last Name
	Zip                        string      `json:"zip"`                         // 12345
	Phone                      string      `json:"phone,omitempty"`             // null | 123-456-7890
	StolenVehicle              bool        `json:"stolenVehicle"`               // true | false
	VehicleName                string      `json:"vehicleName"`                 // Subaru Outback LXT
	Features                   []string    `json:"features"`                    // "11.6MMAN", "ABS_MIL", "ACCS", "AHBL_MIL", "ATF_MIL", "AWD_MIL", "BSD", "BSDRCT_MIL", "CEL_MIL", "EBD_MIL", "EOL_MIL", "EPAS_MIL", "EPB_MIL", "ESS_MIL", "EYESIGHT", "ISS_MIL", "NAV_TOMTOM", "OPL_MIL", "RAB_MIL", "RCC", "REARBRK", "RES", "RESCC", "RHSF", "RPOI", "RPOIA", "SRH_MIL", "SRS_MIL", "TEL_MIL", "TPMS_MIL", "VDC_MIL", "WASH_MIL", "g2"
	Vin                        string      `json:"vin"`                         // 4Y1SL65848Z411439
	VehicleKey                 int64       `json:"vehicleKey"`                  // 3832950
	Nickname                   string      `json:"nickname"`                    // Subaru Outback LXT
	ModelName                  string      `json:"modelName"`                   // Outback
	ModelYear                  string      `json:"modelYear"`                   // 2020
	ModelCode                  string      `json:"modelCode"`                   // LDJ
	ExtDescrip                 string      `json:"extDescrip"`                  // Abyss Blue Pearl (ext color)
	IntDescrip                 string      `json:"intDescrip"`                  // Gray (int color)
	TransCode                  string      `json:"transCode"`                   // CVT
	EngineSize                 float64     `json:"engineSize"`                  // 2.4
	Phev                       bool        `json:"phev"`                        // null
	CachedStateCode            string      `json:"cachedStateCode"`             // NJ
	LicensePlate               string      `json:"licensePlate"`                // NJ
	LicensePlateState          string      `json:"licensePlateState"`           // ABCDEF
	SubscriptionStatus         string      `json:"subscriptionStatus"`          // ACTIVE
	SubscriptionFeatures       []string    `json:"subscriptionFeatures"`        // "[ REMOTE ], [ SAFETY ], [ Retail | Finance3 | RetailPHEV ]""
	SubscriptionPlans          []string    `json:"subscriptionPlans,omitempty"` // []
	VehicleGeoPosition         GeoPosition `json:"vehicleGeoPosition"`          // GeoPosition struct
	AccessLevel                int         `json:"accessLevel"`                 // -1
	VehicleMileage             int         `json:"vehicleMileage,omitempty"`    // null
	CrmRightToRepair           bool        `json:"crmRightToRepair"`            // true | false
	AuthorizedVehicle          bool        `json:"authorizedVehicle"`           // false | true
	NeedMileagePrompt          bool        `json:"needMileagePrompt"`           // false | true
	RemoteServicePinExist      bool        `json:"remoteServicePinExist"`       // true | false
	NeedEmergencyContactPrompt bool        `json:"needEmergencyContactPrompt"`  // false | true
	Show3GSunsetBanner         bool        `json:"show3gSunsetBanner"`          // false | true
	Provisioned                bool        `json:"provisioned"`                 // true | false
	TimeZone                   string      `json:"timeZone"`                    // America/New_York
	SunsetUpgraded             bool        `json:"sunsetUpgraded"`              // true | false
	PreferredDealer            string      `json:"preferredDealer,omitempty"`   // null |
	VehicleBranded             bool        `json:"vehicleBranded"`
}

// GeoPosition .
type GeoPosition struct {
	Latitude  float64     `json:"latitude"`          // 40.700184
	Longitude float64     `json:"longitude"`         // -74.401375
	Speed     int         `json:"speed,omitempty"`   // 62
	Heading   int         `json:"heading,omitempty"` // 155
	Timestamp CustomTime1 `json:"timestamp"`         // "2021-12-22T13:14:47"
}

// VehicleStatus .
type VehicleStatus struct {
	VehicleId                                int64    `json:"vhsId"`                                        // + 9969776690 5198812434
	OdometerValue                            int      `json:"odometerValue"`                                // + 23787
	OdometerValueKm                          int      `json:"odometerValueKilometers"`                      // + 38273
	EventDate                                UnixTime `json:"eventDate"`                                    // + 1701896993000
	EventDateStr                             string   `json:"eventDateStr"`                                 // + 2023-12-06T21:09+0000
	EventDateCarUser                         UnixTime `json:"eventDateCarUser"`                             // + 1701896993000
	EventDateStrCarUser                      string   `json:"eventDateStrCarUser"`                          // + 2023-12-06T21:09+0000
	Latitude                                 float64  `json:"latitude"`                                     // + 40.700183
	Longitude                                float64  `json:"longitude"`                                    // + -74.401372
	Heading                                  int      `json:"positionHeadingDegree,string"`                 // + "154"
	DistanceToEmptyFuelMiles                 float64  `json:"distanceToEmptyFuelMiles"`                     // + 209.4
	DistanceToEmptyFuelKilometers            int      `json:"distanceToEmptyFuelKilometers"`                // + 337
	DistanceToEmptyFuelMiles10s              int      `json:"distanceToEmptyFuelMiles10s"`                  // + 210
	DistanceToEmptyFuelKilometers10s         int      `json:"distanceToEmptyFuelKilometers10s"`             // + 340
	AvgFuelConsumptionMpg                    float64  `json:"avgFuelConsumptionMpg"`                        // + 18.4
	AvgFuelConsumptionLitersPer100Kilometers float64  `json:"avgFuelConsumptionLitersPer100Kilometers"`     // + 12.8
	RemainingFuelPercent                     int      `json:"remainingFuelPercent"`                         // + 82
	TirePressureFrontLeft                    int      `json:"tirePressureFrontLeft,string,omitempty"`       // + "2275"
	TirePressureFrontRight                   int      `json:"tirePressureFrontRight,string,omitempty"`      // + "2344"
	TirePressureRearLeft                     int      `json:"tirePressureRearLeft,string,omitempty"`        // + "2413"
	TirePressureRearRight                    int      `json:"tirePressureRearRight,string,omitempty"`       // + "2344"
	TirePressureFrontLeftPsi                 float64  `json:"tirePressureFrontLeftPsi,string,omitempty"`    // + "33"
	TirePressureFrontRightPsi                float64  `json:"tirePressureFrontRightPsi,string,omitempty"`   // + "34"
	TirePressureRearLeftPsi                  float64  `json:"tirePressureRearLeftPsi,string,omitempty"`     // + "35"
	TirePressureRearRightPsi                 float64  `json:"tirePressureRearRightPsi,string,omitempty"`    // + "34"
	TyreStatusFrontLeft                      string   `json:"tyreStatusFrontLeft"`                          // + "UNKNOWN"
	TyreStatusFrontRight                     string   `json:"tyreStatusFrontRight"`                         // + "UNKNOWN"
	TyreStatusRearLeft                       string   `json:"tyreStatusRearLeft"`                           // + "UNKNOWN"
	TyreStatusRearRight                      string   `json:"tyreStatusRearRight"`                          // + "UNKNOWN"
	EvStateOfChargePercent                   float64  `json:"evStateOfChargePercent,omitempty"`             // + null
	EvDistanceToEmptyMiles                   int      `json:"evDistanceToEmptyMiles,omitempty"`             // + null
	EvDistanceToEmptyKilometers              int      `json:"evDistanceToEmptyKilometers,omitempty"`        // + null
	EvDistanceToEmptyByStateMiles            int      `json:"evDistanceToEmptyByStateMiles,omitempty"`      // + null
	EvDistanceToEmptyByStateKilometers       int      `json:"evDistanceToEmptyByStateKilometers,omitempty"` // + null
	VehicleStateType                         string   `json:"vehicleStateType"`                             // + "IGNITION_OFF | IGNITION_ON"
	WindowFrontLeftStatus                    string   `json:"windowFrontLeftStatus"`                        // CLOSE | VENTED | OPEN
	WindowFrontRightStatus                   string   `json:"windowFrontRightStatus"`                       // CLOSE | VENTED | OPEN
	WindowRearLeftStatus                     string   `json:"windowRearLeftStatus"`                         // CLOSE | VENTED | OPEN
	WindowRearRightStatus                    string   `json:"windowRearRightStatus"`                        // CLOSE | VENTED | OPEN
	WindowSunroofStatus                      string   `json:"windowSunroofStatus"`                          // CLOSE | SLIDE_PARTLY_OPEN | OPEN | TILT
	DoorBootPosition                         string   `json:"doorBootPosition"`                             // CLOSED | OPEN
	DoorEngineHoodPosition                   string   `json:"doorEngineHoodPosition"`                       // CLOSED | OPEN
	DoorFrontLeftPosition                    string   `json:"doorFrontLeftPosition"`                        // CLOSED | OPEN
	DoorFrontRightPosition                   string   `json:"doorFrontRightPosition"`                       // CLOSED | OPEN
	DoorRearLeftPosition                     string   `json:"doorRearLeftPosition"`                         // CLOSED | OPEN
	DoorRearRightPosition                    string   `json:"doorRearRightPosition"`                        // CLOSED | OPEN
	DoorBootLockStatus                       string   `json:"doorBootLockStatus"`                           // LOCKED | UNLOCKED
	DoorFrontLeftLockStatus                  string   `json:"doorFrontLeftLockStatus"`                      // LOCKED | UNLOCKED
	DoorFrontRightLockStatus                 string   `json:"doorFrontRightLockStatus"`                     // LOCKED | UNLOCKED
	DoorRearLeftLockStatus                   string   `json:"doorRearLeftLockStatus"`                       // LOCKED | UNLOCKED
	DoorRearRightLockStatus                  string   `json:"doorRearRightLockStatus"`                      // LOCKED | UNLOCKED
}

// VehicleCondition .
// "dataName":"remoteServiceStatus"
// "remoteServiceType":"condition"
type VehicleCondition struct {
	VehicleStateType           string  `json:"vehicleStateType"`                 // "IGNITION_OFF | IGNITION_ON"
	AvgFuelConsumption         float64 `json:"avgFuelConsumption,omitempty"`     // null | 18.4
	AvgFuelConsumptionUnit     string  `json:"avgFuelConsumptionUnit"`           // "MPG"
	DistanceToEmptyFuel        int     `json:"distanceToEmptyFuel,omitempty"`    // null | 160
	DistanceToEmptyFuelUnit    string  `json:"distanceToEmptyFuelUnit"`          // "MILES"
	RemainingFuelPercent       int     `json:"remainingFuelPercent,string"`      // "66"
	Odometer                   int     `json:"odometer"`                         // 92
	OdometerUnit               string  `json:"odometerUnit"`                     // "MILES"
	TirePressureFrontLeft      float64 `json:"tirePressureFrontLeft,omitempty"`  // null | 36
	TirePressureFrontLeftUnit  string  `json:"tirePressureFrontLeftUnit"`        // "PSI"
	TirePressureFrontRight     float64 `json:"tirePressureFrontRight,omitempty"` // null | 36
	TirePressureFrontRightUnit string  `json:"tirePressureFrontRightUnit"`       // "PSI",
	TirePressureRearLeft       float64 `json:"tirePressureRearLeft,omitempty"`   // null | 36
	TirePressureRearLeftUnit   string  `json:"tirePressureRearLeftUnit"`         // "PSI"
	TirePressureRearRight      float64 `json:"tirePressureRearRight,omitempty"`  // null | 36
	TirePressureRearRightUnit  string  `json:"tirePressureRearRightUnit"`        // "PSI"
	DoorBootPosition           string  `json:"doorBootPosition"`                 // "CLOSED | OPEN"
	DoorEngineHoodPosition     string  `json:"doorEngineHoodPosition"`           // "CLOSED | OPEN"
	DoorFrontLeftPosition      string  `json:"doorFrontLeftPosition"`            // "CLOSED | OPEN"
	DoorFrontRightPosition     string  `json:"doorFrontRightPosition"`           // "CLOSED | OPEN"
	DoorRearLeftPosition       string  `json:"doorRearLeftPosition"`             // "CLOSED | OPEN"
	DoorRearRightPosition      string  `json:"doorRearRightPosition"`            // "CLOSED | OPEN"
	WindowFrontLeftStatus      string  `json:"windowFrontLeftStatus"`            // "CLOSE | VENTED | OPEN"
	WindowFrontRightStatus     string  `json:"windowFrontRightStatus"`           // "CLOSE | VENTED | OPEN"
	WindowRearLeftStatus       string  `json:"windowRearLeftStatus"`             // "CLOSE | VENTED | OPEN"
	WindowRearRightStatus      string  `json:"windowRearRightStatus"`            // "CLOSE | VENTED | OPEN"
	WindowSunroofStatus        string  `json:"windowSunroofStatus"`              // "CLOSE | VENTED | OPEN"
	EvDistanceToEmpty          int     `json:"evDistanceToEmpty,omitempty"`      // null,
	EvDistanceToEmptyUnit      string  `json:"evDistanceToEmptyUnit,omitempty"`  // null,
	EvChargerStateType         string  `json:"evChargerStateType,omitempty"`     // null,
	EvIsPluggedIn              bool    `json:"evIsPluggedIn,omitempty"`          // null,
	EvStateOfChargeMode        string  `json:"evStateOfChargeMode,omitempty"`    // null,
	EvTimeToFullyCharged       string  `json:"evTimeToFullyCharged,omitempty"`   // null,
	EvStateOfChargePercent     int     `json:"evStateOfChargePercent,omitempty"` // null,
	LastUpdatedTime            string  `json:"lastUpdatedTime"`                  // "2023-04-10T17:50:54+0000",
}

// ClimateProfile represents a climate control profile for a Subaru vehicle.
type ClimateProfile struct {
	Name                      string `json:"name"`
	VehicleType               string `json:"vehicleType,omitempty"`       // vehicleType                  [ gas | phev ]
	PresetType                string `json:"presetType"`                  // presetType                   [ subaruPreset | userPreset ]
	StartConfiguration        string `json:"startConfiguration"`          // startConfiguration           [ START_ENGINE_ALLOW_KEY_IN_IGNITION (gas) | START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION (phev) ]
	RunTimeMinutes            int    `json:"runTimeMinutes,string"`       // runTimeMinutes               [ 0 | 1 | 5 | 10 ]
	HeatedRearWindowActive    string `json:"heatedRearWindowActive"`      // heatedRearWindowActive:      [ false | true ]
	HeatedSeatFrontRight      string `json:"heatedSeatFrontRight"`        // heatedSeatFrontRight:        [ OFF | LOW_HEAT | MEDIUM_HEAT | HIGH_HEAT ]
	HeatedSeatFrontLeft       string `json:"heatedSeatFrontLeft"`         // heatedSeatFrontLeft:         [ OFF | LOW_HEAT | MEDIUM_HEAT | HIGH_HEAT ]
	ClimateZoneFrontTemp      int    `json:"climateZoneFrontTemp,string"` // climateZoneFrontTemp:        [ for _ in range(60, 85 + 1)] // climateZoneFrontTempCelsius: [for _ in range(15, 30 + 1) ]
	ClimateZoneFrontAirMode   string `json:"climateZoneFrontAirMode"`     // climateZoneFrontAirMode:     [ WINDOW | FEET_WINDOW | FACE | FEET | FEET_FACE_BALANCED | AUTO ]
	ClimateZoneFrontAirVolume string `json:"climateZoneFrontAirVolume"`   // climateZoneFrontAirVolume:   [ AUTO | 2 | 4 | 7 ]
	OuterAirCirculation       string `json:"outerAirCirculation"`         // airConditionOn:              [ auto | outsideAir | true ]
	AirConditionOn            string `json:"airConditionOn"`              // airConditionOn:              [ false | true ]
	CanEdit                   string `json:"canEdit"`                     // canEdit                      [ false | true ]
	Disabled                  string `json:"disabled"`                    // disabled                     [ false | true ]
}

type ClimateProfiles map[string]ClimateProfile

// GeoLocation represents the geographical location of a Subaru vehicle.
type GeoLocation struct {
	Latitude  float64     `json:"latitude"`          // 40.700184
	Longitude float64     `json:"longitude"`         // -74.401375
	Heading   int         `json:"heading,omitempty"` // 189
	Speed     int         `json:"speed,omitempty"`   // 0.00
	Updated   CustomTime1 `json:"timestamp"`         // "2025-07-08T19:05:07"
}

// ServiceRequest .
// "dataName": "remoteServiceStatus"
type ServiceRequest struct {
	ServiceRequestID   string          `json:"serviceRequestId,omitempty"` // 4S4BTGND8L3137058_1640294426029_19_@NGTP
	Vin                string          `json:"vin"`                        // 4S4BTGND8L3137058
	Success            bool            `json:"success"`                    // false | true // Could be in the false state while the executed request in the progress
	Cancelled          bool            `json:"cancelled"`                  // false | true
	RemoteServiceType  string          `json:"remoteServiceType"`          // vehicleStatus | condition | locate | unlock | lock | lightsOnly | engineStart | engineStop | phevChargeNow
	RemoteServiceState string          `json:"remoteServiceState"`         // started | finished | stopping
	SubState           string          `json:"subState,omitempty"`         // null
	ErrorCode          string          `json:"errorCode,omitempty"`        // null:null
	Result             json.RawMessage `json:"result,omitempty"`           // struct
	UpdateTime         UnixTime        `json:"updateTime,omitempty"`       // timestamp // is empty if the request is started
}

// parse parses the JSON response from the MySubaru API into a ServiceRequest struct.
func (sr *ServiceRequest) parse(b []byte, logger *slog.Logger) error {
	err := json.Unmarshal(b, &sr)
	if err != nil {
		logger.Error("error while parsing json", "request", "GetVehicleCondition", "error", err.Error())
	}
	if !sr.Success && sr.ErrorCode != "" {
		logger.Error("error in response", "request", "GetVehicleCondition", "errorCode", sr.ErrorCode, "remoteServiceType", sr.RemoteServiceType)
		switch sr.ErrorCode {
		case API_ERRORS["API_ERROR_SERVICE_ALREADY_STARTED"]:
			return errors.New("error in response: Service already started")
		case API_ERRORS["API_ERROR_VEHICLE_NOT_IN_ACCOUNT"]:
			return errors.New("error in response: Vehicle not in account")
		case API_ERRORS["API_ERROR_SOA_403"]:
			return errors.New("error in response: Unable to parse response body, SOA 403 error")
		default:
			return errors.New("error in response: " + sr.ErrorCode)
		}
	}
	return nil
}

type VehicleHealth struct {
	VehicleHealthItems []VehicleHealthItem `json:"vehicleHealthItems"`
	LastUpdatedDate    int64               `json:"lastUpdatedDate"`
}

type VehicleHealthItem struct {
	WarningCode int        `json:"warningCode"`       // internal code used by MySubaru, not documented
	B2cCode     string     `json:"b2cCode"`           // oilTemp | airbag | oilLevel | etc.
	FeatureCode string     `json:"featureCode"`       // SRS_MIL | CEL_MIL | ATF_MIL | etc.
	IsTrouble   bool       `json:"isTrouble"`         // false | true
	OnDaiID     int        `json:"onDaiId"`           // Has a number, probably internal record id
	OnDates     []UnixTime `json:"onDates,omitempty"` // List of the timestamps
}

type ErrorResponse struct {
	ErrorLabel       string `json:"errorLabel"`                 // "404-soa-unableToParseResponseBody"
	ErrorDescription string `json:"errorDescription,omitempty"` // null
}

// UnixTime is a wrapper around time.Time that allows us to marshal and unmarshal Unix timestamps
type UnixTime struct {
	time.Time
}

// UnmarshalJSON is the method that satisfies the Unmarshaller interface
// Note that it uses a pointer receiver. It needs this because it will be modifying the embedded time.Time instance
func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	u.Time = time.Unix(timestamp, 0)
	return nil
}

// MarshalJSON turns our time.Time back into an int
func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", u.Unix())), nil
}

// CustomTime1 "2021-12-22T13:14:47" is a custom type for unmarshalling time strings without timezone
type CustomTime1 struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface for CustomTime1
func (ct *CustomTime1) UnmarshalJSON(b []byte) error {
	return unmarshalTime(b, "2006-01-02T15:04:05", &ct.Time)
}

// MarshalJSON implements the json.Marshaler interface for CustomTime1
func (ct CustomTime1) MarshalJSON() ([]byte, error) {
	return marshalTime(ct.Time)
}

// CustomTime2 "2023-04-10T17:50:54+0000" is a custom type for unmarshalling time strings with timezone offset
type CustomTime2 struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface for CustomTime2
func (ct *CustomTime2) UnmarshalJSON(b []byte) error {
	return unmarshalTime(b, "2006-01-02T15:04:05-0700", &ct.Time)
}

// MarshalJSON implements the json.Marshaler interface for CustomTime2
func (ct CustomTime2) MarshalJSON() ([]byte, error) {
	return marshalTime(ct.Time)
}

// unmarshalTime is a helper function for parsing time strings with custom layouts
func unmarshalTime(b []byte, layout string, t *time.Time) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		*t = time.Time{}
		return nil
	}

	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("failed to parse time %q with layout %q: %w", s, layout, err)
	}

	*t = parsedTime
	return nil
}

// marshalTime is a helper function for marshaling time to JSON
func marshalTime(t time.Time) ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02T15:04:05Z07:00"))), nil
}
