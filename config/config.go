package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// MetricsRecorder defines the interface for recording metrics
type MetricsRecorder interface {
	// RecordRequest records an API request
	RecordRequest(method, endpoint string, duration time.Duration, success bool)

	// RecordError records an error occurrence
	RecordError(errorType string)

	// RecordRetry records a retry attempt
	RecordRetry(endpoint string, attempt int)
}

const (
	LoggingOutputJson = "JSON"
	LoggingOutputText = "TEXT"
)

// Config .
type Config struct {
	MySubaru MySubaru
	TimeZone string
	Logger   *slog.Logger
	Metrics  MetricsRecorder
}

// config defines the structure of configuration data to be parsed from a config source.
type config struct {
	MySubaru MySubaru `json:"mysubaru" yaml:"mysubaru"`
	TimeZone string   `json:"timezone" yaml:"timezone"`
	Logging  *Logging `json:"logging" yaml:"logging"`
}

// MySubaru .
type MySubaru struct {
	Credentials   Credentials `json:"credentials" yaml:"credentials"`
	Region        string      `json:"region" yaml:"region"`
	AutoReconnect bool        `json:"auto_reconnect" yaml:"auto_reconnect"`
}

// Credentials .
type Credentials struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	PIN        string `json:"pin" yaml:"pin"`
	DeviceID   string `json:"deviceid" yaml:"deviceid"`
	DeviceName string `json:"devicename" yaml:"devicename"`
}

// Logging .
type Logging struct {
	Level  string `json:"level" yaml:"level"`
	Output string `json:"output" yaml:"output"`
	Source bool   `json:"source,omitempty" yaml:"source,omitempty"`
}

func (l Logging) ToLogger() *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(l.Level)); err != nil {
		level = slog.LevelInfo
	}

	var handler slog.Handler
	switch l.Output {
	case LoggingOutputJson:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: l.Source, Level: level})
	case LoggingOutputText:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: l.Source, Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: l.Source, Level: level})
	}

	return slog.New(handler)
}

// FromBytes unmarshals a byte slice of JSON or YAML config data into a valid server options value.
func FromBytes(b []byte) (*Config, error) {
	c := new(config)
	o := Config{}

	if len(b) == 0 {
		return nil, nil
	}

	if b[0] == '{' {
		err := json.Unmarshal(b, c)
		if err != nil {
			return nil, err
		}
	} else {
		err := yaml.Unmarshal(b, c)
		if err != nil {
			return nil, err
		}
	}

	o.MySubaru = c.MySubaru
	o.TimeZone = c.TimeZone
	if c.Logging != nil {
		o.Logger = c.Logging.ToLogger()
	} else {
		o.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return &o, nil
}
