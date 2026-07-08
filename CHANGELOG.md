# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.10.0] - 2026-07-08

### Added
- Native GTK Webkit window desktop client integration using `webview_go` (V2.14).
- Dedicated `/docs` folder with detailed guides for setup, compiling/building, and architecture.
- Page favicon integration using the project's anchor logo in `index.html` (V2.10).
- Cross-platform DOCKER_HOST suggestions for Linux, Windows Named Pipes, TCP and SSH in `.env.example` (V2.13).

### Changed
- Translated the entire UI (labels, buttons, tables, logs) to English (V2.12).
- Cleaned up redundant comments in codebase and translated remaining technical comments to English (V2.11).
- Restructured `README.md` to be a concise bilingual landing page pointing to `/docs` guides.

## [2.9.0] - 2026-07-08

### Added
- Environment variables support through local `.env` files.
- Configurable listen port using the `PORT` variable with a fallback default to `8080`.
- Configuration template file `.env.example`.

## [2.8.0] - 2026-07-08

### Changed
- Refactored domain layer to segregate the monolithic repository interface into role-based interfaces (ISP).
- Split HTTP router logic into resource-focused handlers (SRP) to improve code maintainability and testability.

## [2.7.0] - 2026-07-08

### Added
- Local Docker images list at the top of the dashboard showing IDs, tags, sizes, and creation dates.
- Inline details panel opening directly below the selected container row in the table.
- Interactive toggle icon to expand/collapse Compose cards and the Images card.

### Fixed
- Terminal websocket connection deadlock causing black screen.
- Layout squishing on tables when detail panel is active by adopting full-width stacked views.

### Changed
- CPU and RAM metrics layout moved to the top of the details panel, keeping telemetry visible even on the terminal tab.