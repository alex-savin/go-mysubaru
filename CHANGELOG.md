# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-02-23

### Added

- **Valet Mode Support**
  - `GetValetModeStatus()` - Get current valet mode status
  - `GetValetModeSettings()` - Retrieve valet mode configuration
  - `ValetModeStart()` - Activate valet mode
  - `ValetModeStop()` - Deactivate valet mode
  - `SaveValetModeSettings()` - Save valet mode configuration

- **Geo-fence (Boundary Alerts) Support**
  - `GetGeoFenceSettings()` - Get geo-fence configuration
  - `SaveGeoFenceSettings()` - Save geo-fence configuration
  - `ActivateGeoFence()` - Enable boundary alerts
  - `DeactivateGeoFence()` - Disable boundary alerts

- **Speed-fence (Speed Alerts) Support**
  - `GetSpeedFenceSettings()` - Get speed alert configuration
  - `SaveSpeedFenceSettings()` - Save speed alert configuration
  - `ActivateSpeedFence()` - Enable speed alerts
  - `DeactivateSpeedFence()` - Disable speed alerts

- **Curfew Alerts Support**
  - `GetCurfewSettings()` - Get curfew configuration
  - `SaveCurfewSettings()` - Save curfew configuration
  - `ActivateCurfew()` - Enable curfew alerts
  - `DeactivateCurfew()` - Disable curfew alerts

- **Trip Tracker Support**
  - `GetTrips()` - Retrieve trip history
  - `TripLogStart()` - Start trip logging
  - `TripLogStop()` - Stop trip logging
  - `DeleteTrip()` - Remove a trip from history

- **POI/Destination Support**
  - `SendPOI()` - Send point of interest to vehicle navigation
  - `GetFavoritePOIs()` - Retrieve saved favorite destinations
  - `SaveFavoritePOI()` - Save a favorite destination

- **Roadside Assistance Support**
  - `GetRoadsideAssistance()` - Get roadside assistance info and status
  - `RequestRoadsideAssistance()` - Request roadside assistance

- **Vehicle Safety Information**
  - `GetRecalls()` - Retrieve open recalls for the vehicle
  - `GetWarningLights()` - Get active warning lights/indicators

- **Enhanced Error Handling**
  - `NegativeAckError` type for vehicle-side command rejections
  - `PINLockedError` type with timeout information
  - `IsNegativeAckError()` helper function
  - `IsPINLockedError()` helper function
  - `ParseAPIError()` function to convert error codes to typed errors
  - Added comprehensive error variables:
    - `ErrInvalidSession`, `ErrNoVehicles`, `ErrVehicleNotInAccount`
    - `ErrStolenVehicle`, `ErrInvalidPIN`, `ErrServiceInProgress`
    - `ErrTokenGenFailed`, `ErrAccIsOn`, `ErrDoorNotClosed`
    - `ErrEngineRunning`, `ErrHoodNotClosed`, `ErrIgnitionOn`
    - `ErrKeyInIgnition`, `ErrRemoteStartActive`, `ErrTrunkNotClosed`
    - `ErrLowBattery`, `ErrAlarmActive`, `ErrValetModeOn`

- **New API Endpoints** (~100+ new endpoints)
  - Valet mode endpoints
  - Geo-fence/boundary alert endpoints
  - Speed-fence/speed alert endpoints
  - Curfew alert endpoints
  - Trip tracker endpoints
  - POI/destination endpoints
  - Roadside assistance endpoints
  - Notification endpoints
  - Maintenance schedule endpoints
  - Dealer and appointment endpoints
  - Recall and event endpoints
  - Profile and authorized users endpoints
  - JWT token endpoints

- **New Test Fixtures**
  - `valetModeSettings.json`, `geoFenceSettings.json`, `speedFenceSettings.json`
  - `curfewSettings.json`, `trips.json`, `favoritePOIs.json`
  - `roadsideAssistance.json`, `recalls.json`, `warningLights.json`
  - `remoteServiceStatus.json`

### Changed

- Updated API version from `/g2v31` to `/g2v32`
- Added ~45 new API error codes including all `NegativeAcknowledge_*` codes
- Expanded `apiErrorMessages` map with comprehensive error mappings
- Updated test fixtures and added comprehensive tests for new features

### Fixed

- Improved error handling with specific error types for better debugging

## [1.0.0] - 2026-01-11

### Added

- Initial release of the MySubaru Go client library
- Authentication with username/password and PIN
- Two-factor authentication (2FA) support for device registration
- Vehicle management and selection
- Vehicle status retrieval (fuel level, odometer, tire pressure, etc.)
- Vehicle health reports with trouble code detection
- Remote commands support:
  - Lock/Unlock doors
  - Start/Stop engine (remote start)
  - Horn and lights
  - Climate control presets
- Location tracking and polling
- Session management with automatic refresh
- Configuration via YAML or JSON files
- Structured logging with slog
- Comprehensive test coverage
- CI/CD pipeline with GitHub Actions
- Example application with usage demonstrations

### Security

- Credentials stored securely in configuration
- Session token management with automatic renewal
- PIN verification for sensitive operations

[Unreleased]: https://github.com/alex-savin/go-mysubaru/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/alex-savin/go-mysubaru/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/alex-savin/go-mysubaru/releases/tag/v1.0.0
