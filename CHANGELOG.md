# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-21

### Added
- Backend datasource plugin with Go backend and React frontend
- Config editor with API URL, username, secret, and integration code fields
- Query editor with entity type selector, time field mapping, and filter input
- Support for querying Tickets, Companies, Contacts, and Resources
- Health check via Autotask Zone Information API
- Time range filtering — map Grafana time range to Autotask date fields
- Autotask API authentication (Basic auth + API headers)
- Rate limiting via autotask-go client library
- Docker-based development environment with hot reload
- Cross-compilation for Linux and macOS (amd64 + arm64)
