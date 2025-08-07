# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Fixed
  
### Changed

## [0.6.0] - 2025-08-07

### Added
- metric for canceled tasks (thanks to @uhelm)

## [0.5.1] - 2023-11-21

### Added
- bump go version
- use builtin log/slog package

### Changed
- better parsing of flags

## [0.5.0] - 2023-07-27

### Added
- new metric about last task execution duration

### Changed
- change the flags name to respect other exporters naming

## [0.4.0] - 2023-05-24

### Added
- documentation about configuration params

### Changed
- made name of parameters similar to other exporters
- simplified project layout

## [0.3.1] - 2023-05-18

### Fixed
- wrong handling of -log.format
- wrong source in log messages
  
## [0.3.0] - 2023-05-17

### Changed

- upgrade dependency (slog)
- set version in dev mode
  
## [0.2.2] - 2023-05-17

### Changed

- change directory name to support better `go install`
  
## [0.2.1] - 2023-05-17

### Added

- instructions on how to configure Huey2
- package godoc
  
### Fixed

### Changed

- print only version with `--version`
- use dot notation for flags family
- add more UTs
  
### Removed

## [0.2.0] - 2023-05-14

### Added

- option to show version
  
## [0.1.4] - 2023-05-13

### Changed

- improve release page

## [0.1.3] - 2023-05-13

### Added

- add release changelog categories
  
## [0.1.2] - 2023-05-13

### Fixed

- correct name in README file
  
## [0.1.1] - 2023-05-13

### Added

- GitHub Actions
  
## [0.1.0] - 2023-05-05

### Added

- Initial implementation
