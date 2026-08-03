# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [1.3.0] - 2025-10-26

### Added

- Add IsPreferredPaymentMethodInvalidError (#29)
- Add frontend status component and improve typing definitions (#25)

### Changed

- Change Github actions to run on pull requests (#26)

### Fixed

- Fix missing categories by matching short categories code in plan code (#28)
- Fix README badge

## [1.2.0] - 2025-01-21

### Added

- Add nextjs frontend (#7)
- Add error IsPlanNotFoundError and matcher kimsufi.IsAvailabilityNotFoundError and kimsufi.IsPreferredPaymentMethodNotSetError
- Add possibility to try many options combinations at once
- Add more datacenters name and --list-datacenters flag to check command

### Changed

- Update order command, add capability to try all datacenters
- Upgrade Go dependencies
- Update order command, add capability to try all datacenters

### Fixed

- Show unknown categories
- Fix kimsufi.WithAuth

## [1.1.0] - 2024-12-01

### Added

- Add Golang CLI
- Add check command --list-options flag
- Add --human flag in list command, and unify --help text
- Add Github test action badge
- Build all binaries in CI

### Changed

- Update order command, better help and return error from printPrices
- Update USAGE and add command examples
- Update README examples
- Let goreleaser handle builds
- Let goreleaser publish binaries
- Keep Makefile simple
- Update Go dependencies: go mod tidy
- Update goreleaser to reflect Makefile
- Restore .github/workflows/release.yaml
- Update Makefile: output into build directory
- Nancy ignore dependencies not-included in final binary

### Removed

- Remove Bash scripts
- Remove CircleCI config
- Remove RUN_IN_CI.md
- Remove unneeded assets

### Fixed

- Fix CI: install dependencies
- Fix releaser Github action
- Fix options.Merge: check existing Family
- Fix linting errors

[Unreleased]: https://github.com/TheoBrigitte/kimsufi-notifier/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/TheoBrigitte/kimsufi-notifier/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/TheoBrigitte/kimsufi-notifier/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/TheoBrigitte/kimsufi-notifier/releases/tag/v1.1.0
