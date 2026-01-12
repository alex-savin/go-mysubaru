# MySubaru Go Client Example

This directory contains a comprehensive example demonstrating all features of the MySubaru Go client library.

## Setup

1. **Copy the sample configuration**:
   ```bash
   cp ../config.sample.yaml config.yaml
   ```

2. **Edit `config.yaml`** with your MySubaru credentials:
   ```yaml
   mysubaru:
     credentials:
       username: your-email@example.com
       password: "your-password"
       pin: "1234"
       deviceid: your-device-id
       devicename: My Go App
     region: USA
   ```

## Running the Example

### Build the example:
```bash
go build -o example-app .
```

### Run with default config:
```bash
./example-app
```

### Run with custom config:
```bash
./example-app path/to/your/config.yaml
```

### Show help:
```bash
./example-app --help
```

## Configuration File Format

The `config.yaml` file supports the following structure:

```yaml
mysubaru:
  credentials:
    username: user@email.com          # MySubaru account email
    password: "Secr#TPassW0rd"       # MySubaru account password
    pin: "PIN"                       # 4-digit vehicle PIN
    deviceid: GENERATE-DEVICE-ID     # Unique device identifier
    devicename: Golang Integration   # Human-readable device name
  region: USA                        # Region: USA or CAN
  auto_reconnect: true               # Auto-reconnect on connection loss
timezone: "America/New_York"         # Timezone for logging
logging:
  level: INFO                        # Log level: DEBUG, INFO, WARN, ERROR
  output: JSON                       # Output format: JSON or TEXT
  source: false                      # Include source file info in logs
```

## Features Demonstrated

The example showcases all major MySubaru API features:

### 🚗 Basic Vehicle Operations
- Authentication and session management
- Vehicle discovery and selection
- Vehicle status (odometer, fuel level, location)
- Vehicle condition (doors, windows, tires)
- Vehicle health diagnostics

### 🔧 Remote Control
- Lock/unlock doors
- Engine start/stop
- Horn and lights control

### 🌡️ Climate Control
- Climate presets (factory, user, quick-start)
- Temperature and fan settings

### 🛡️ Safety Features (G2 Telematics)
- Geofencing setup and management
- Speed fence configuration
- Curfew scheduling

### ⚡ Electric Vehicle Features
- Charge control and scheduling
- EV status monitoring

## Error Handling

The example includes comprehensive error handling for:
- Configuration file errors
- Authentication failures
- Network connectivity issues
- API rate limiting
- Subscription requirement errors

## Security Notes

- Never commit your `config.yaml` file to version control
- Use strong, unique passwords
- Consider using environment variables for sensitive data in production
- The example uses test credentials by default - update with real credentials to test actual API calls

## Troubleshooting

### Configuration Errors
```
Failed to load configuration: open config.yaml: no such file or directory
```
**Solution**: Copy `config.sample.yaml` to `config.yaml` and update with your credentials.

### Authentication Errors
```
Failed to create client: authentication failed
```
**Solution**: Verify your MySubaru credentials in `config.yaml`.

### Network Errors
```
Failed to get vehicles: connection timeout
```
**Solution**: Check your internet connection and MySubaru service status.

## Advanced Usage

### Custom Configuration
```go
// Load custom config file
cfg, err := loadConfig("production.yaml")
if err != nil {
    log.Fatal(err)
}
```

### Environment Variables
For production deployments, consider using environment variables:
```bash
export MYSUBARU_USERNAME="user@example.com"
export MYSUBARU_PASSWORD="password"
export MYSUBARU_DEVICE_ID="device-id"
```

### Logging Configuration
Customize logging output:
```yaml
logging:
  level: DEBUG
  output: TEXT
  source: true
```

## Related

- [Main README](../README.md) - Go client library documentation