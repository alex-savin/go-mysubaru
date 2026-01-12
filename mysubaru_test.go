package mysubaru

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnixTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTime  time.Time
		wantError bool
	}{
		{
			name:     "valid unix timestamp",
			input:    "1700000000",
			wantTime: time.Unix(1700000000, 0),
		},
		{
			name:      "invalid string",
			input:     "\"notanumber\"",
			wantError: true,
		},
		{
			name:      "empty input",
			input:     "",
			wantError: true,
		},
		{
			name:      "float value",
			input:     "1700000000.123",
			wantError: true,
		},
		{
			name:     "zero timestamp",
			input:    "0",
			wantTime: time.Unix(0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ut UnixTime
			err := ut.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantError {
				t.Errorf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && !ut.Time.Equal(tt.wantTime) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", ut.Time, tt.wantTime)
			}
		})
	}
}

func TestUnixTime_UnmarshalJSON_withJSONUnmarshal(t *testing.T) {
	type testStruct struct {
		Time UnixTime `json:"time"`
	}
	input := `{"time":1700000000}`
	var ts testStruct
	err := json.Unmarshal([]byte(input), &ts)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	want := time.Unix(1700000000, 0)
	if !ts.Time.Time.Equal(want) {
		t.Errorf("UnmarshalJSON() got = %v, want %v", ts.Time.Time, want)
	}
}
func TestUnixTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  string
	}{
		{
			name:  "epoch",
			input: time.Unix(0, 0),
			want:  "0",
		},
		{
			name:  "positive unix time",
			input: time.Unix(1700000000, 0),
			want:  "1700000000",
		},
		{
			name:  "negative unix time",
			input: time.Unix(-100, 0),
			want:  "-100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := UnixTime{Time: tt.input}
			got, err := ut.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestUnixTime_MarshalJSON_withJSONMarshal(t *testing.T) {
	type testStruct struct {
		Time UnixTime `json:"time"`
	}
	ts := testStruct{Time: UnixTime{Time: time.Unix(1700000000, 0)}}
	b, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	want := `{"time":1700000000}`
	if string(b) != want {
		t.Errorf("json.Marshal() = %s, want %s", string(b), want)
	}
}

// func TestResponse_parse(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		wantErr  error
// 		wantCode string
// 		wantLog  string
// 	}{
// 		{
// 			name:    "success response",
// 			input:   `{"success":true,"dataName":"foo","data":{}}`,
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "invalid json",
// 			input:   `{"success":tru`,
// 			wantErr: errors.New("error while parsing json:"),
// 			wantLog: "error while parsing json",
// 		},
// 		{
// 			name:     "API_ERROR_NO_ACCOUNT",
// 			input:    `{"success":false,"errorCode":"noAccount","dataName":"errorResponse","data":{}}`,
// 			wantErr:  errors.New("error in response: Account not found"),
// 			wantCode: "noAccount",
// 			wantLog:  "error in response",
// 		},
// 		{
// 			name:     "API_ERROR_INVALID_CREDENTIALS",
// 			input:    `{"success":false,"errorCode":"invalidCredentials","dataName":"errorResponse","data":{}}`,
// 			wantErr:  errors.New("error in response: Invalid Credentials"),
// 			wantCode: "invalidCredentials",
// 			wantLog:  "error in response",
// 		},
// 		{
// 			name:     "API_ERROR_SOA_403",
// 			input:    `{"success":false,"errorCode":"404-soa-unableToParseResponseBody","dataName":"errorResponse","data":{}}`,
// 			wantErr:  errors.New("error in response: Unable to parse response body, SOA 403 error"),
// 			wantCode: "404-soa-unableToParseResponseBody",
// 			wantLog:  "error in response",
// 		},
// 		{
// 			name:     "unknown error code",
// 			input:    `{"success":false,"errorCode":"somethingElse","dataName":"errorResponse","data":{}}`,
// 			wantErr:  errors.New("error in response: somethingElse"),
// 			wantCode: "somethingElse",
// 			wantLog:  "error in response",
// 		},
// 		{
// 			name:    "no errorCode but not success",
// 			input:   `{"success":false,"dataName":"errorResponse","data":{}}`,
// 			wantErr: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var resp Response
// 			logger := slog.New(slog.NewTextHandler(nil, nil))
// 			got, err := resp.parse([]byte(tt.input), logger)
// 			if tt.wantErr != nil {
// 				if err == nil {
// 					t.Fatalf("expected error, got nil")
// 				}
// 				if !contains(err.Error(), tt.wantErr.Error()) {
// 					t.Errorf("parse() error = %v, want %v", err, tt.wantErr)
// 				}
// 			} else if err != nil {
// 				t.Errorf("parse() unexpected error: %v", err)
// 			}
// 			if tt.wantCode != "" && got != nil && got.ErrorCode != tt.wantCode {
// 				t.Errorf("parse() got.ErrorCode = %v, want %v", got.ErrorCode, tt.wantCode)
// 			}
// 		})
// 	}
// }

// // contains is a helper for substring matching.
// func contains(s, substr string) bool {
// 	return bytes.Contains([]byte(s), []byte(substr))
// }

func TestCustomTime1_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTime  time.Time
		wantError bool
	}{
		{
			name:     "valid time without timezone",
			input:    `"2021-12-22T13:14:47"`,
			wantTime: time.Date(2021, 12, 22, 13, 14, 47, 0, time.UTC),
		},
		{
			name:     "null value",
			input:    "null",
			wantTime: time.Time{},
		},
		{
			name:      "invalid format",
			input:     `"2021-12-22"`,
			wantError: true,
		},
		{
			name:      "empty string",
			input:     `""`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime1
			err := ct.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantError {
				t.Errorf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && !ct.Time.Equal(tt.wantTime) {
				t.Errorf("UnmarshalJSON() got time = %v, want %v", ct.Time, tt.wantTime)
			}
		})
	}
}

func TestCustomTime1_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		time     CustomTime1
		expected string
	}{
		{
			name:     "zero time",
			time:     CustomTime1{Time: time.Time{}},
			expected: "null",
		},
		{
			name:     "valid time",
			time:     CustomTime1{Time: time.Date(2021, 12, 22, 13, 14, 47, 0, time.UTC)},
			expected: `"2021-12-22T13:14:47Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.time.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() got = %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestCustomTime2_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTime  time.Time
		wantError bool
	}{
		{
			name:     "valid time with timezone",
			input:    `"2023-04-10T17:50:54+0000"`,
			wantTime: time.Date(2023, 4, 10, 17, 50, 54, 0, time.UTC),
		},
		{
			name:     "valid time with negative timezone",
			input:    `"2023-04-10T17:50:54-0700"`,
			wantTime: time.Date(2023, 4, 10, 17, 50, 54, 0, time.FixedZone("", -7*3600)),
		},
		{
			name:     "null value",
			input:    "null",
			wantTime: time.Time{},
		},
		{
			name:      "invalid format",
			input:     `"2023-04-10T17:50:54"`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime2
			err := ct.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantError {
				t.Errorf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && !ct.Time.Equal(tt.wantTime) {
				t.Errorf("UnmarshalJSON() got time = %v, want %v", ct.Time, tt.wantTime)
			}
		})
	}
}

func TestCustomTime2_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		time     CustomTime2
		expected string
	}{
		{
			name:     "zero time",
			time:     CustomTime2{Time: time.Time{}},
			expected: "null",
		},
		{
			name:     "valid time",
			time:     CustomTime2{Time: time.Date(2023, 4, 10, 17, 50, 54, 0, time.UTC)},
			expected: `"2023-04-10T17:50:54Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.time.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() got = %s, want %s", string(data), tt.expected)
			}
		})
	}
}
