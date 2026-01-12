# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/alex-savin/go-mysubaru/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/alex-savin/go-mysubaru/releases/tag/v1.0.0
