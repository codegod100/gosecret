# Changelog

All notable changes to GoSecret will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Simple key/value interface for secret management
- `set` command to store secrets with optional value parameter
- `get` command to retrieve secrets by key
- `delete` command to remove secrets by key
- `list` command to show all stored secrets with optional filtering
- Compatibility aliases (`store`, `lookup`, `clear`) for secret-tool users
- Integration with GNOME Keyring via D-Bus Secret Service API
- Secure password prompting when no value provided
- Support for stdin password input
- Cross-platform binary builds for Linux, macOS, and Windows
- Comprehensive documentation and examples

### Technical Details
- Written in Go with minimal dependencies
- Uses `github.com/godbus/dbus/v5` for D-Bus communication
- Uses `golang.org/x/term` for secure terminal input
- Stores secrets with `application=gosecret` attribute for organization
- Maintains compatibility with freedesktop.org Secret Service specification

## [1.0.0] - Initial Release

### Added
- Complete rewrite of secret-tool functionality in Go
- Simplified command-line interface
- Cross-platform support
- GitHub Actions for automated builds and releases