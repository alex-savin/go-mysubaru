# MySubaru Go Client

[![CI](https://github.com/alex-savin/go-mysubaru/actions/workflows/ci.yml/badge.svg)](https://github.com/alex-savin/go-mysubaru/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex-savin/go-mysubaru)](https://goreportcard.com/report/github.com/alex-savin/go-mysubaru)
[![Go Reference](https://pkg.go.dev/badge/github.com/alex-savin/go-mysubaru.svg)](https://pkg.go.dev/github.com/alex-savin/go-mysubaru)

A Go client library for interacting with the MySubaru API. Supports authentication, vehicle status, remote commands, and more.

## Installation

```bash
go get github.com/alex-savin/go-mysubaru
```

## Quick Start

```go
package main

import (
    "log"
    "log/slog"

    "github.com/alex-savin/go-mysubaru"
    "github.com/alex-savin/go-mysubaru/config"
)

func main() {
    cfg := &config.Config{
        MySubaru: config.MySubaru{
            Credentials: config.Credentials{
                Username:   "your-email@example.com",
                Password:   "your-password",
                PIN:        "1234",
                DeviceID:   "your-device-id",
                DeviceName: "My Go App",
            },
            Region: "USA", // or "CAN"
        },
        Logger: slog.Default(),
    }

    client, err := mysubaru.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Authenticate
    ok, authErr, needs2FA := client.Authenticate()
    if needs2FA {
        log.Fatal("Device not registered; 2FA verification required")
    }
    if !ok || authErr != nil {
        log.Fatal("Authentication failed:", authErr)
    }

    // Get vehicles
    vehicles, err := client.GetVehicles()
    if err != nil {
        log.Fatal(err)
    }

    for _, v := range vehicles {
        log.Printf("Vehicle: %s (%s)", v.CarNickname, v.Vin)
    }
}
```

## Configuration

### From File (YAML or JSON)

```go
data, _ := os.ReadFile("config.yaml")
cfg, _ := config.FromBytes(data)
client, _ := mysubaru.New(cfg)
```

### Example config.yaml

```yaml
mysubaru:
  credentials:
    username: your-email@example.com
    password: "your-password"
    pin: "1234"
    deviceid: your-device-id
    devicename: My Go App
  region: USA

logging:
  level: info
  output: TEXT  # or JSON
```

## Metrics

The client supports pluggable metrics collection via the `MetricsRecorder` interface:

```go
type MetricsRecorder interface {
    RecordRequest(method, endpoint string, duration time.Duration, success bool)
    RecordError(errorType string)
    RecordRetry(endpoint string, attempt int)
}
```

### Prometheus Example

```go
import (
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/alex-savin/go-mysubaru/config"
)

type PrometheusMetrics struct {
    requestDuration *prometheus.HistogramVec
    requestTotal    *prometheus.CounterVec
    errorTotal      *prometheus.CounterVec
    retryTotal      *prometheus.CounterVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
    m := &PrometheusMetrics{
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "mysubaru_request_duration_seconds",
                Help:    "Duration of MySubaru API requests",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "endpoint"},
        ),
        requestTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "mysubaru_requests_total",
                Help: "Total MySubaru API requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        errorTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "mysubaru_errors_total",
                Help: "Total MySubaru API errors by type",
            },
            []string{"error_type"},
        ),
        retryTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "mysubaru_retries_total",
                Help: "Total MySubaru API retry attempts",
            },
            []string{"endpoint", "attempt"},
        ),
    }

    prometheus.MustRegister(m.requestDuration, m.requestTotal, m.errorTotal, m.retryTotal)
    return m
}

func (p *PrometheusMetrics) RecordRequest(method, endpoint string, duration time.Duration, success bool) {
    status := "success"
    if !success {
        status = "failure"
    }
    p.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
    p.requestTotal.WithLabelValues(method, endpoint, status).Inc()
}

func (p *PrometheusMetrics) RecordError(errorType string) {
    p.errorTotal.WithLabelValues(errorType).Inc()
}

func (p *PrometheusMetrics) RecordRetry(endpoint string, attempt int) {
    p.retryTotal.WithLabelValues(endpoint, strconv.Itoa(attempt)).Inc()
}
```

### Using Metrics

```go
cfg := &config.Config{
    MySubaru: config.MySubaru{
        Credentials: config.Credentials{...},
        Region:      "USA",
    },
    Logger:  slog.Default(),
    Metrics: NewPrometheusMetrics(),
}

client, _ := mysubaru.New(cfg)
```

If no `Metrics` is provided, a no-op recorder is used.

## API Reference

### Client Methods

#### Authentication

| Method | Description |
|--------|-------------|
| `Authenticate() (bool, error, bool)` | Authenticates with MySubaru. Returns (success, error, needsVerification) |
| `RequestAuthCode(email string) error` | Requests a 2FA verification code to be sent |
| `SubmitAuthCode(code string, permanent bool) error` | Submits the 2FA verification code |

#### Vehicle Management

| Method | Description |
|--------|-------------|
| `GetVehicles() ([]*Vehicle, error)` | Returns all vehicles associated with the account |
| `GetVehicleByVin(vin string) (*Vehicle, error)` | Returns a specific vehicle by VIN |
| `SelectVehicle(vin string) (*VehicleData, error)` | Selects a vehicle for subsequent operations |
| `RefreshVehicles() error` | Refreshes vehicle data from the API |

### Vehicle Methods

#### Remote Commands

All remote commands return `(chan string, error)`. The channel receives status updates as the command progresses.

```go
// Lock/Unlock
vehicle.Lock()
vehicle.Unlock()

// Remote Start with climate settings
vehicle.EngineStart(runMinutes, delayMinutes, honkHorn)
vehicle.EngineStartWithProfile(runMinutes, delayMinutes, honkHorn, profileName)
vehicle.EngineStop()

// Horn & Lights
vehicle.HornStart()
vehicle.HornStop()
vehicle.LightsStart()
vehicle.LightsStop()

// EV Charging (PHEV/BEV only)
vehicle.ChargeOn()
```

#### Vehicle Information

```go
// Basic info
vehicle.Vin              // "4S4BTGND8L3137058"
vehicle.CarName          // "Subaru Outback LXT"
vehicle.CarNickname      // "My Outback"
vehicle.ModelYear        // "2020"
vehicle.ModelName        // "Outback"

// Status
vehicle.EngineState      // "IGNITION_OFF", "RUNNING", etc.
vehicle.Odometer.Miles   // 24999
vehicle.DistanceToEmpty.Miles      // 149
vehicle.DistanceToEmpty.Percentage // 66

// Fuel economy
vehicle.FuelConsumptionAvg.MPG     // 18.5
vehicle.FuelConsumptionAvg.LP100Km // 12.7

// Component states
vehicle.Doors["FrontLeft"].Status   // "CLOSED"
vehicle.Doors["FrontLeft"].Locked   // true
vehicle.Windows["FrontLeft"].Status // "CLOSED"
vehicle.Tires["FrontLeft"].PressurePsi // 32.5

// Location
vehicle.GeoLocation.Latitude  // 40.7128
vehicle.GeoLocation.Longitude // -74.0060
vehicle.GeoLocation.Heading   // 180

// EV Status (for PHEV/BEV)
vehicle.IsEV()                           // true/false
vehicle.EVStatus.StateOfChargePercent    // 80
vehicle.EVStatus.DistanceToEmptyMiles    // 25
vehicle.EVStatus.IsPluggedIn             // true
vehicle.EVStatus.ChargerStateType        // "CHARGING"

// Trouble codes
vehicle.Troubles["P0301"]  // Trouble{Code: "P0301", Description: "Cylinder 1 Misfire"}
```

#### Climate Profiles

```go
// Get available climate presets (Subaru defaults + user presets)
vehicle.GetClimatePresets()

// Get user-defined presets only
vehicle.GetClimateUserPresets()

// Access presets
for name, profile := range vehicle.ClimateProfiles {
    fmt.Printf("%s: %d°F, Fan %s\n", name, profile.ClimateZoneFrontTemp, profile.ClimateZoneFrontAirVolume)
}

// Delete a user preset by name
err := vehicle.DeleteClimateUserPreset("My Winter Preset")

// Save a list of user presets (overwrites existing)
presets := []mysubaru.ClimateProfile{
    {
        Name:                      "Morning Commute",
        RunTimeMinutes:            10,
        ClimateZoneFrontTemp:      72,
        ClimateZoneFrontAirMode:   "AUTO",
        ClimateZoneFrontAirVolume: "AUTO",
        HeatedSeatFrontLeft:       "MEDIUM_HEAT",
        HeatedSeatFrontRight:      "OFF",
        HeatedRearWindowActive:    "true",
        OuterAirCirculation:       "outsideAir",
        AirConditionOn:            "false",
    },
}
err := vehicle.SaveClimateUserPresets(presets)
```

### Example: Remote Start with Status Tracking

```go
package main

import (
    "fmt"
    "log"

    "github.com/alex-savin/go-mysubaru"
    "github.com/alex-savin/go-mysubaru/config"
)

func main() {
    // Setup client...
    cfg := &config.Config{...}
    client, _ := mysubaru.New(cfg)
    client.Authenticate()

    vehicles, _ := client.GetVehicles()
    vehicle := vehicles[0]

    // Start engine with climate profile
    statusChan, err := vehicle.EngineStartWithProfile(10, 0, false, "Winter")
    if err != nil {
        log.Fatal(err)
    }

    // Track command progress
    for status := range statusChan {
        fmt.Printf("Status: %s\n", status)
        // Outputs: "started", "pending", "success" or "failed"
    }
}
```

### Example: Monitor Vehicle Location

```go
// Get current location (cached)
fmt.Printf("Location: %f, %f\n", vehicle.GeoLocation.Latitude, vehicle.GeoLocation.Longitude)

// Force refresh from vehicle
statusChan, err := vehicle.GetLocation(true)
if err != nil {
    log.Fatal(err)
}

for status := range statusChan {
    fmt.Printf("Location update: %s\n", status)
}

// Now access updated location
fmt.Printf("Updated: %f, %f (heading %d°)\n", 
    vehicle.GeoLocation.Latitude, 
    vehicle.GeoLocation.Longitude,
    vehicle.GeoLocation.Heading)
```

### Example: Check Door and Window States

```go
vehicles, _ := client.GetVehicles()
vehicle := vehicles[0]

// Check all doors
for name, door := range vehicle.Doors {
    status := "closed"
    if door.Status != "CLOSED" {
        status = "OPEN"
    }
    lockStatus := "unlocked"
    if door.Locked {
        lockStatus = "locked"
    }
    fmt.Printf("%s: %s, %s\n", name, status, lockStatus)
}

// Check all windows
for name, window := range vehicle.Windows {
    fmt.Printf("%s: %s\n", name, window.Status)
}

// Check tire pressures
for name, tire := range vehicle.Tires {
    fmt.Printf("%s: %.1f PSI\n", name, tire.PressurePsi)
}
```

### Example: Manage Climate Presets

```go
vehicles, _ := client.GetVehicles()
vehicle := vehicles[0]

// Fetch all climate presets (Subaru defaults + user-defined)
vehicle.GetClimatePresets()

// List available presets
fmt.Println("Available climate presets:")
for name, profile := range vehicle.ClimateProfiles {
    fmt.Printf("  %s: %d°F, Fan: %s, Mode: %s\n", 
        profile.Name,
        profile.ClimateZoneFrontTemp,
        profile.ClimateZoneFrontAirVolume,
        profile.ClimateZoneFrontAirMode)
}

// Create a new user preset
newPreset := mysubaru.ClimateProfile{
    Name:                      "My Morning Commute",
    RunTimeMinutes:            10,
    ClimateZoneFrontTemp:      72,
    ClimateZoneFrontAirMode:   "AUTO",
    ClimateZoneFrontAirVolume: "AUTO",
    HeatedSeatFrontLeft:       "MEDIUM_HEAT",
    HeatedSeatFrontRight:      "OFF",
    HeatedRearWindowActive:    "true",
    OuterAirCirculation:       "outsideAir",
    AirConditionOn:            "false",
}

// Get existing user presets and add the new one
vehicle.GetClimateUserPresets()
var userPresets []mysubaru.ClimateProfile
for _, p := range vehicle.ClimateProfiles {
    if p.PresetType == "userPreset" {
        userPresets = append(userPresets, p)
    }
}
userPresets = append(userPresets, newPreset)

// Save (max 4 user presets allowed)
if err := vehicle.SaveClimateUserPresets(userPresets); err != nil {
    log.Printf("Failed to save presets: %v", err)
}

// Delete a preset by name
if err := vehicle.DeleteClimateUserPreset("My Morning Commute"); err != nil {
    log.Printf("Failed to delete preset: %v", err)
}
```

### Error Handling

The library provides typed errors for common API failures:

```go
import "errors"

_, err, _ := client.Authenticate()
if err != nil {
    var apiErr mysubaru.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case "INVALID_CREDENTIALS":
            log.Fatal("Wrong username or password")
        case "ACCOUNT_LOCKED":
            log.Fatal("Account is locked, contact Subaru support")
        case "RATE_LIMITED":
            log.Println("Too many requests, retrying...")
            // Retry logic
        }
    }
}

// Check if error is retryable
if mysubaru.IsRetryableError(err) {
    // Safe to retry
}
```

## Features

- **Authentication**: Login, 2FA verification, session management
- **Vehicle Status**: Odometer, fuel level, tire pressure, door/window states
- **Remote Commands**: Lock/unlock, remote start/stop, horn/lights
- **EV Support**: Battery status, charging control
- **Health Reports**: Trouble codes, maintenance alerts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT
