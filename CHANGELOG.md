# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- **Valet status on vehicles without valet mode**: `GetValetModeStatus` and
  `GetValetModeSettings` no longer fail with `json: cannot unmarshal string into
  Go value of type mysubaru.ValetModeSettings` when the backend returns the
  `data` field as a plain string (or null) for a vehicle that doesn't have valet
  mode provisioned. Those responses are now reported as a disabled (zero-value)
  configuration instead of an error.

## [2.0.0] - 2026-06-12

### ⚠️ Breaking Changes

- **`context.Context` everywhere**: every `Client` and `Vehicle` API method now
  takes a `ctx context.Context` as its first argument. The context bounds the
  whole operation — per-attempt HTTP timeout, retry backoff waits, and remote
  command polling — and cancelling it stops them.
- **`Authenticate` signature**: `Authenticate() (bool, error, bool)` is now the
  idiomatic `Authenticate(ctx) (ok bool, needs2FA bool, err error)`.
- **Exported mutable package state removed**: `MOBILE_API_SERVER`, `MOBILE_APP`,
  `WEB_API_SERVER`, `API_ERRORS`, and `APP_ERRORS` are now unexported. Use the
  typed errors (`APIError`, `IsSessionError`, …) for error classification and
  `config.MySubaru.BaseURL` to point at a different host. The internal
  `"TEST"` region entry was removed from production code.
- **`Client` no longer embeds `sync.RWMutex`**: the accidental public
  `Lock`/`Unlock`/`RLock`/`RUnlock` methods are gone (now an unexported field).
- Dropped the deprecated `github.com/pkg/errors` dependency; all errors are
  wrapped with `fmt.Errorf("...: %w", err)` and remain `errors.Is/As`-friendly.

### Added

- **MySubaru app 3.2.7 API updates** (from the decompiled APK reference)
  - `GetAppStatus(ctx)` - query the app availability / maintenance gate
    (`appStatus.json`, new in 3.2.7); distinguishes a maintenance window
    from a transport failure
  - `Logout(ctx)` - invalidate the session on the backend
    (`invalidateSession.json`) and clear local auth state
  - New endpoint constants: `API_APP_STATUS`, `API_INVALIDATE_SESSION`,
    `API_PROFILE_2FA_VERIFY` (PIN reset via 2FA),
    `API_MICRO_VEHICLE_ACCOUNT_ATTRIBUTES`
  - `MICROSERVICE_API_SERVER` base URLs for the JWT-only `/micro/` endpoints
- `config.MySubaru.BaseURL` (`base_url` in YAML/JSON) to override the regional
  mobile-API host (QA environments, local mocks)

### Changed

- **Session-validity caching**: any successful API response now counts as proof
  of session liveness for 4 minutes (the backend expires idle sessions at 5).
  Within that window, command preambles skip the `validateSession.json` +
  `selectVehicle.json` round-trips, removing up to 3 HTTP calls per command
  during bursts.
- **Centralized HTML-error-page detection**: the transport now detects HTML
  responses (login redirects, maintenance pages) once in `executeOnce` for all
  ~50 endpoints, instead of ad-hoc checks at a handful of call sites.
- Vehicle fetch/save endpoints share a single `fetchInto` helper (subscription
  check, session validation, vehicle selection, parse), removing several
  hundred lines of duplicated boilerplate.
- `GetVehicleByVin()` now returns an error when the vehicle payload fails to
  parse (previously continued with empty vehicle data) and logs warnings when
  initial status/condition/health fetches fail
- Deduplicated HTTP client construction between `New()` and session reset
- `New()` defaults `updateInterval`/`fetchInterval` from
  `DEFAULT_UPDATE_INTERVAL`/`DEFAULT_FETCH_INTERVAL` instead of diverging
  hard-coded literals

### Fixed

- `TIR_33` feature description said "Rear 35" instead of "Rear 33"
- Typo in the password-warning message ("Mutiple" → "Multiple")
- Removed dead Python-style `strftime` format constants (`TIMESTAMP_FMT`,
  `POSITION_TIMESTAMP_FMT`) left over from the subarulink port

### Notes

- API version `g2v33` (3.2.7 moved all `mobileapi` base URLs from `g2v32`;
  the client already tracks this and auto-bumps on 404)
- The module path is now `github.com/alex-savin/go-mysubaru/v2`, as Go modules
  require for v2.0.0+ releases.

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

[Unreleased]: https://github.com/alex-savin/go-mysubaru/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/alex-savin/go-mysubaru/compare/v1.1.0...v2.0.0
[1.1.0]: https://github.com/alex-savin/go-mysubaru/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/alex-savin/go-mysubaru/releases/tag/v1.0.0
