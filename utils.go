package mysubaru

import (
	"errors"
	"fmt"
	"math"
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// isHTMLResponse checks if the response body contains HTML instead of JSON.
// This can happen when the API returns an error page instead of a proper JSON response.
// Returns true if the response appears to be HTML.
func isHTMLResponse(body []byte) bool {
	return strings.HasPrefix(strings.TrimSpace(string(body)), "<")
}

// errHTMLResponse creates a standard error for when API returns HTML instead of JSON.
func errHTMLResponse(request string) error {
	return errors.New("API returned HTML error page instead of JSON for " + request + " - session may be invalid or API may have changed")
}

// getResponsePreview returns first n characters of response for logging, safely handling short responses.
func getResponsePreview(body []byte, maxLen int) string {
	s := string(body)
	if len(s) < maxLen {
		return s
	}
	return s[:maxLen]
}

// timestamp is a function
func timestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
}

// urlToGen .
func urlToGen(url string, gen string) string {
	var re = regexp.MustCompile(`api_gen`)
	// dirty trick for current G3
	if gen == "g3" {
		gen = "g2"
	}
	url = re.ReplaceAllString(url, gen)

	return url
}

// VinCheck - Vehicle Identification Number check digit validation
// Parameter: string - 17 digit VIN
// Return:
//
//	1- boolean - Validity flag. Set to true if VIN check digit is correct, false otherwise.
//	2- string - Valid VIN. Same VIN passed as parameter but with the correct check digit on it.
func vinCheck(vin string) (bool, string) {
	var valid = false
	vin = strings.ToUpper(vin)
	var retVin = vin

	if len(vin) == 17 {
		traSum := transcodeDigits(vin)
		checkNum := math.Mod(float64(traSum), 11)
		var checkDigit byte
		if checkNum == 10 {
			checkDigit = byte('X')
		} else {
			checkDigitTemp := strconv.Itoa(int(checkNum))
			checkDigit = checkDigitTemp[len(checkDigitTemp)-1]
		}
		if retVin[8] == checkDigit {
			valid = true
		}
		retVin = retVin[:8] + string(checkDigit) + retVin[9:]
	} else {
		valid = false
		retVin = ""
	}

	return valid, retVin
}

// transcodeDigits transcodes VIN digits to a numeric value
func transcodeDigits(vin string) int {
	var digitSum = 0
	var code int
	for i, chr := range vin {
		code = 0

		switch chr {
		case 'A', 'J', '1':
			code = 1
		case 'B', 'K', 'S', '2':
			code = 2
		case 'C', 'L', 'T', '3':
			code = 3
		case 'D', 'M', 'U', '4':
			code = 4
		case 'E', 'N', 'V', '5':
			code = 5
		case 'F', 'W', '6':
			code = 6
		case 'G', 'P', 'X', '7':
			code = 7
		case 'H', 'Y', '8':
			code = 8
		case 'R', 'Z', '9':
			code = 9
		case 'I', 'O', 'Q':
			code = 0
		}
		switch i + 1 {
		case 1, 11:
			digitSum += code * 8
		case 2, 12:
			digitSum += code * 7
		case 3, 13:
			digitSum += code * 6
		case 4, 14:
			digitSum += code * 5
		case 5, 15:
			digitSum += code * 4
		case 6, 16:
			digitSum += code * 3
		case 7, 17:
			digitSum += code * 2
		case 8:
			digitSum += code * 10
		case 9:
			digitSum += code * 0
		case 10:
			digitSum += code * 9
		}
	}

	return digitSum
}

// emailMasking takes an email address as input and returns a version of the email
// with the username part partially hidden for privacy. Only the first and last
// characters of the username are visible, with the middle characters replaced by asterisks.
// The function validates the email format before processing.
// Returns the obfuscated email or an error if the input is not a valid email address.
func emailMasking(email string) (string, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("invalid email address: %s", email)
	}

	re1 := regexp.MustCompile(`^(.*?)@(.*)$`)
	matches := re1.FindStringSubmatch(email)

	var username, domain string
	if len(matches) == 3 { // Expecting the full match, username, and domain
		username = matches[1]
		domain = matches[2]
	} else {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	re2 := regexp.MustCompile(`(.)(.*)(.)`)

	replacedString := re2.ReplaceAllStringFunc(username, func(s string) string {
		firstChar := string(s[0])
		lastChar := string(s[len(s)-1])
		middleCharsCount := len(s) - 2

		if middleCharsCount < 0 { // Should not happen with the length check above, but for robustness
			return s
		}
		return firstChar + strings.Repeat("*", middleCharsCount) + lastChar
	})

	return replacedString + "@" + domain, nil
}

// containsValueInStruct checks if any string field in the given struct 's' contains the specified 'search' substring (case-insensitive).
// It returns true if at least one string field contains the substring, and false otherwise.
// If 's' is not a struct, it returns false.
func containsValueInStruct(s any, search string) bool {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return false // Not a struct
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			if strings.Contains(strings.ToLower(field.String()), strings.ToLower(search)) {
				return true
			}
		}
	}

	return false
}

// Input validation functions

// ValidateVIN validates a Vehicle Identification Number
func ValidateVIN(vin string) error {
	if len(vin) != 17 {
		return fmt.Errorf("VIN must be exactly 17 characters, got %d", len(vin))
	}

	// Check for valid characters (alphanumeric, no I, O, Q)
	validChars := regexp.MustCompile(`^[A-HJ-NPR-Z0-9]+$`)
	if !validChars.MatchString(strings.ToUpper(vin)) {
		return fmt.Errorf("VIN contains invalid characters")
	}

	// Validate check digit
	valid, _ := vinCheck(vin)
	if !valid {
		return fmt.Errorf("VIN check digit is invalid")
	}

	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address format: %w", err)
	}
	return nil
}

// ValidatePIN validates a 4-digit PIN
func ValidatePIN(pin string) error {
	if len(pin) != 4 {
		return fmt.Errorf("PIN must be exactly 4 digits, got %d", len(pin))
	}

	if matched, _ := regexp.MatchString(`^\d{4}$`, pin); !matched {
		return fmt.Errorf("PIN must contain only digits")
	}

	return nil
}

// ValidateCoordinates validates latitude and longitude
func ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", lat)
	}

	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", lon)
	}

	return nil
}

// ValidateTimeRange validates start and end times in HH:MM format
func ValidateTimeRange(startTime, endTime string) error {
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(startTime) {
		return fmt.Errorf("start time must be in HH:MM format, got %s", startTime)
	}
	if !timeRegex.MatchString(endTime) {
		return fmt.Errorf("end time must be in HH:MM format, got %s", endTime)
	}
	return nil
}

// ValidateSpeedLimit validates speed limit in mph
func ValidateSpeedLimit(speed int) error {
	if speed < 5 || speed > 140 {
		return fmt.Errorf("speed limit must be between 5 and 140 mph, got %d", speed)
	}
	return nil
}

// ValidateGeoFenceRadius validates geofence radius in meters
func ValidateGeoFenceRadius(radius int) error {
	if radius < 100 || radius > 10000 {
		return fmt.Errorf("geofence radius must be between 100 and 10000 meters, got %d", radius)
	}
	return nil
}

// ValidateDaysOfWeek validates array of days (0=Sunday, 6=Saturday)
func ValidateDaysOfWeek(days []int) error {
	for _, day := range days {
		if day < 0 || day > 6 {
			return fmt.Errorf("day of week must be between 0 (Sunday) and 6 (Saturday), got %d", day)
		}
	}
	return nil
}

// ValidateDeviceID validates device ID format
func ValidateDeviceID(deviceID string) error {
	if len(deviceID) == 0 {
		return fmt.Errorf("device ID cannot be empty")
	}
	if len(deviceID) > 100 {
		return fmt.Errorf("device ID too long, maximum 100 characters")
	}
	return nil
}

// ValidateDeviceName validates device name
func ValidateDeviceName(deviceName string) error {
	if len(deviceName) == 0 {
		return fmt.Errorf("device name cannot be empty")
	}
	if len(deviceName) > 50 {
		return fmt.Errorf("device name too long, maximum 50 characters")
	}
	return nil
}

// timeTrack .
// func timeTrack(name string) {
// 	start := time.Now()
// 	fmt.Printf("%s took %v\n", name, time.Since(start))
// }
