# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-07

Initial public release.

### Added

- Canonical, versioned email event taxonomy covering origin, gateway, mailbox,
  engagement, and status event types.
- Reason taxonomy for classifying failures with sender/recipient attribution.
- Pluggable architecture for **enrichers** and **classifiers**.
- Default toolkit (`toolkit.NewDefaultEventClassifier`) to quickly classify events
  from vendor facts.
- Bounce classification via Sisimai (`libsisimai.org/sisimai/v5`).

### Documentation

- See `README.md` for installation, usage examples, and architecture.
