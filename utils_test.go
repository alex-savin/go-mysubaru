package mysubaru

import (
	"regexp"
	"strconv"
	"testing"
	"time"
)

// timestamp returns the current time in milliseconds since epoch as a string.
func TestTimestamp(t *testing.T) {
	ts1 := timestamp()
	time.Sleep(1 * time.Millisecond)
	ts2 := timestamp()

	// Should be numeric
	if _, err := strconv.ParseInt(ts1, 10, 64); err != nil {
		t.Errorf("timestamp() returned non-numeric string: %s", ts1)
	}
	// Should be increasing
	if ts1 >= ts2 {
		t.Errorf("timestamp() not increasing: %s >= %s", ts1, ts2)
	}
}

// timestamp returns the current time in milliseconds since epoch as a string.
func TestTimestamp_Format(t *testing.T) {
	ts := timestamp()
	matched, err := regexp.MatchString(`^\d+$`, ts)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Errorf("timestamp() = %q, want only digits", ts)
	}
}

// urlToGen replaces "api_gen" in the URL with the specified generation.
func TestUrlToGen(t *testing.T) {
	tests := []struct {
		url, gen, want string
	}{
		{"https://host/api_gen/endpoint", "g1", "https://host/g1/endpoint"},
		{"https://host/api_gen/endpoint", "g2", "https://host/g2/endpoint"},
		{"https://host/api_gen/endpoint", "g3", "https://host/g2/endpoint"}, // g3 special case
		{"https://host/api_gen/api_gen", "g1", "https://host/g1/g1"},
		{"https://host/other/endpoint", "g1", "https://host/other/endpoint"},
	}
	for _, tt := range tests {
		got := urlToGen(tt.url, tt.gen)
		if got != tt.want {
			t.Errorf("urlToGen(%q, %q) = %q, want %q", tt.url, tt.gen, got, tt.want)
		}
	}
}

// vinCheck validates the VIN check digit and returns the corrected VIN.
func TestVinCheck_Valid(t *testing.T) {
	// Example valid VIN: 1HGCM82633A004352 (check digit is '3')
	vin := "1HGCM82633A004352"
	valid, corrected := vinCheck(vin)
	if !valid {
		t.Errorf("vinCheck(%q) = false, want true", vin)
	}
	if corrected != vin {
		t.Errorf("vinCheck(%q) corrected VIN = %q, want %q", vin, corrected, vin)
	}
}

// TestVinCheck_InvalidCheckDigit tests a VIN with an incorrect check digit.
func TestVinCheck_InvalidCheckDigit(t *testing.T) {
	vin := "1HGCM82633A004352"
	// Change check digit (9th char) to '9'
	badVin := vin[:8] + "9" + vin[9:]
	valid, corrected := vinCheck(badVin)
	if valid {
		t.Errorf("vinCheck(%q) = true, want false", badVin)
	}
	// Should correct to original VIN
	if corrected != vin {
		t.Errorf("vinCheck(%q) corrected VIN = %q, want %q", badVin, corrected, vin)
	}
}

// TestVinCheck_WrongLength tests a VIN that is not 17 characters long.
func TestVinCheck_WrongLength(t *testing.T) {
	vin := "1234567890123456" // 16 chars
	valid, corrected := vinCheck(vin)
	if valid {
		t.Errorf("vinCheck(%q) = true, want false", vin)
	}
	if corrected != "" {
		t.Errorf("vinCheck(%q) corrected VIN = %q, want empty string", vin, corrected)
	}
}

func TestValidateVIN_Valid(t *testing.T) {
	vin := "1HGCM82633A004352"
	if err := ValidateVIN(vin); err != nil {
		t.Fatalf("ValidateVIN(%q) returned error: %v", vin, err)
	}
}

func TestValidateVIN_InvalidLength(t *testing.T) {
	vin := "1234567890123456"
	if err := ValidateVIN(vin); err == nil {
		t.Fatalf("ValidateVIN(%q) expected length error, got nil", vin)
	}
}

func TestValidateVIN_InvalidCharacters(t *testing.T) {
	vin := "1HGCM82633A00O352" // contains 'O'
	if err := ValidateVIN(vin); err == nil {
		t.Fatalf("ValidateVIN(%q) expected invalid character error, got nil", vin)
	}
}

func TestValidateVIN_InvalidCheckDigit(t *testing.T) {
	vin := "1HGCM82633A004352"
	badVin := vin[:8] + "9" + vin[9:]
	if err := ValidateVIN(badVin); err == nil {
		t.Fatalf("ValidateVIN(%q) expected check digit error, got nil", badVin)
	}
}

// transcodeDigits computes the sum of the VIN digits according to the VIN rules.
func TestTranscodeDigits(t *testing.T) {
	// Use a known VIN and manually compute the sum
	vin := "1HGCM82633A004352"
	sum := transcodeDigits(vin)
	// Precomputed sum for this VIN is 311 (from online VIN calculator)
	want := 311
	if sum != want {
		t.Errorf("transcodeDigits(%q) = %d, want %d", vin, sum, want)
	}
}

// TestVinCheck_XCheckDigit tests a VIN with 'X' as the check digit.
func TestVinCheck_XCheckDigit(t *testing.T) {
	// VIN with check digit 'X'
	vin := "1M8GDM9AXKP042788"
	valid, corrected := vinCheck(vin)
	if !valid {
		t.Errorf("vinCheck(%q) = false, want true", vin)
	}
	if corrected != vin {
		t.Errorf("vinCheck(%q) corrected VIN = %q, want %q", vin, corrected, vin)
	}
}

// TestUrlToGen_NoApiGen tests the case where the URL does not contain "api_gen".
func TestUrlToGen_NoApiGen(t *testing.T) {
	url := "https://host/endpoint"
	gen := "g1"
	got := urlToGen(url, gen)
	if got != url {
		t.Errorf("urlToGen(%q, %q) = %q, want %q", url, gen, got, url)
	}
}

func TestEmailHidder(t *testing.T) {
	tests := []struct {
		email    string
		expected string
		wantErr  bool
	}{
		{"alex@example.com", "a**x@example.com", false},
		{"a@example.com", "a@example.com", false},
		{"ab@example.com", "ab@example.com", false},
		{"", "", true},
		{"notanemail", "", true},
	}

	for _, tt := range tests {
		got, err := emailMasking(tt.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("emailHidder(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("emailHidder(%q) = %q, want %q", tt.email, got, tt.expected)
		}
	}
}
func TestContainsValueInStruct(t *testing.T) {
	type TestStruct struct {
		Name    string
		Address string
		Age     int
		Note    string
	}
	tests := []struct {
		s      any
		search string
		want   bool
	}{
		{
			s:      TestStruct{Name: "Alice", Address: "123 Main St", Age: 30, Note: "VIP customer"},
			search: "alice",
			want:   true,
		},
		{
			s:      TestStruct{Name: "Bob", Address: "456 Elm St", Age: 25, Note: "Regular"},
			search: "elm",
			want:   true,
		},
		{
			s:      TestStruct{Name: "Charlie", Address: "789 Oak St", Age: 40, Note: "VIP"},
			search: "vip",
			want:   true,
		},
		{
			s:      TestStruct{Name: "Diana", Address: "101 Pine St", Age: 22, Note: ""},
			search: "xyz",
			want:   false,
		},
		{
			s:      TestStruct{Name: "", Address: "", Age: 0, Note: ""},
			search: "",
			want:   true, // empty string is contained in all strings
		},
		{
			s:      struct{ Foo int }{Foo: 42},
			search: "42",
			want:   false,
		},
		{
			s:      "not a struct",
			search: "struct",
			want:   false,
		},
		{
			s:      struct{ S string }{S: "CaseInsensitive"},
			search: "caseinsensitive",
			want:   true,
		},
	}

	for i, tt := range tests {
		got := containsValueInStruct(tt.s, tt.search)
		if got != tt.want {
			t.Errorf("Test %d: containsStringInStruct(%#v, %q) = %v, want %v", i, tt.s, tt.search, got, tt.want)
		}
	}
}
