# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- Initial project scaffold: Cobra CLI, `internal/` diagnostic packages, docs, and CI.
- `kubediag diagnose`, `kubediag nodes`, `kubediag pods` — v1.0 diagnostics (cluster health, node
  conditions, pod health) with terminal, JSON, and HTML report output.
- `kubediag report`, `kubediag version` — report export and version commands.
- `kubediag fix`, `kubediag ingress`, `kubediag storage`, `kubediag network`, `kubediag security` —
  stubbed commands for features planned in later roadmap versions.
