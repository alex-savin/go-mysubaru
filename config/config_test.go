package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	yamlBytes = []byte(`
mysubaru:
  credentials:
    username: username@mysubaru.golang
    password: "PASSWORD"
    pin: "1234"
    deviceid: AaBbCcDdEeFf0123456789
    devicename: MySubaru Golang Client
  region: USA
  auto_reconnect: true
timezone: "America/New_York"
logging:
  level: INFO
  output: TEXT
  source: false
`)

	jsonBytes = []byte(`{
    "mysubaru": {
        "credentials": {
            "username": "username@mysubaru.golang",
            "password": "PASSWORD",
            "pin": "1234",
            "deviceid": "AaBbCcDdEeFf0123456789",
            "devicename": "MySubaru Golang Client"
        },
        "region": "USA",
        "auto_reconnect": true
    },
    "timezone": "America/New_York",
    "logging": {
        "level": "INFO",
        "output": "TEXT",
        "source": false
    }
}
`)
	parsedOptions = Config{
		MySubaru: MySubaru{
			Credentials: Credentials{
				Username:   "username@mysubaru.golang",
				Password:   "PASSWORD",
				PIN:        "1234",
				DeviceID:   "AaBbCcDdEeFf0123456789",
				DeviceName: "MySubaru Golang Client",
			},
			Region:        "USA",
			AutoReconnect: true,
		},
		TimeZone: "America/New_York",
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
)

func TestFromBytesEmptyL(t *testing.T) {
	_, err := FromBytes([]byte{})
	require.NoError(t, err)
}

func TestFromBytesYAML(t *testing.T) {
	o, err := FromBytes(yamlBytes)
	require.NoError(t, err)
	require.Equal(t, parsedOptions, *o)
}

func TestFromBytesYAMLError(t *testing.T) {
	_, err := FromBytes(append(yamlBytes, 'a'))
	require.Error(t, err)
}

func TestFromBytesJSON(t *testing.T) {
	o, err := FromBytes(jsonBytes)
	require.NoError(t, err)
	require.Equal(t, parsedOptions, *o)
}

func TestFromBytesJSONError(t *testing.T) {
	_, err := FromBytes(append(jsonBytes, 'a'))
	require.Error(t, err)
}
