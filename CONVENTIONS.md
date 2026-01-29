# UniFi Go Library - Coding Conventions

This file provides coding conventions and project context for AI assistants working with this codebase.

## Project Overview

This is **github.com/unpoller/unifi/v5**: a Go library that connects to a Ubiquiti UniFi controller and **pulls** data (clients, devices, sites, alarms, events, etc.). It does **not** update or change controller settings by design.

**Entry point**: `unifi.NewUnifi(config *Config)` → authenticated `*Unifi` client.

**Typical usage flow**: `Config` → `NewUnifi(config)` → `GetSites()` → `GetClients(sites)` / `GetDevices(sites)`. Optionally `GetAlarms`, `GetEvents`, `GetAnomalies`, etc.

## Go Coding Standards

- **Linter**: `.golangci.yaml` – `nlreturn`, `revive`, `tagalign`, `testpackage`, `wsl_v5`. Max issues 0; `fix: true`.
- **wsl_v5**: Use blank lines between logical blocks; keep `if`/`for` bodies short.
- **Tests**: Use `package unifi` with `// nolint: testpackage` where needed. Prefer `t.Parallel()`, `assert`/`require` from `github.com/stretchr/testify`. Mock via `UnifiClient` and `mocks.MockUnifi`.
- **Errors**: Wrap with `fmt.Errorf("context: %w", err)`. Use package sentinels (`ErrAuthenticationFailed`, etc.).
- **Imports**: Group stdlib, then third-party.

Run `golangci-lint run` and `go test ./...` before committing.

## UniFi API Patterns

### Local UniFi Controller API (Primary)

**Authentication**: Cookie-based via `/api/login` (standard) or `/api/auth/login` (UDM Pro). Alternative: API key via `X-API-Key` header. Library handles both via `Config.User`/`Config.Pass` or `Config.APIKey`.

**Endpoint pattern**: `/api/s/{site}/...` where `{site}` is typically `"default"`. UDM Pro requires `/proxy/network` prefix.

**Response format**: `{"data": [...], "meta": {"rc": "ok"}}`. Errors: `{"data": [], "meta": {"msg": "api.err.LoginRequired", "rc": "error"}}`. Handle `meta.count` for truncated results.

**Key endpoints** (constants in `types.go`):
- `/api/s/{site}/stat/device` - Devices
- `/api/s/{site}/stat/sta` - Active clients
- `/api/s/{site}/stat/event` - Events
- `/api/s/{site}/list/alarm` - Alarms
- `/api/s/{site}/stat/stadpi` - Client DPI stats
- `/api/s/{site}/stat/sitedpi` - Site DPI stats
- `/api/s/{site}/stat/anomalies` - Anomalies

**Code patterns**:
- Use `u.GetData(apiPath, &response)` for GET requests
- API paths use `fmt.Sprintf(APIDevicePath, site.Name)` etc.
- Most fetches are per-site: `GetSites()` first, then loop `GetDevices(sites)`, `GetClients(sites)`, etc.
- Single-site helpers like `GetUAPs(site)` also exist
- `GetDevices` returns `*Devices` with device slices (UAPs, USWs, USGs, UDMs, UXGs, PDUs, UBBs, UCIs)
- Parsing uses `parseDevices` with `json.RawMessage` and type detection

### Remote Site Manager API (Cloud)

**Base URL**: `https://api.ui.com/v1/`

**Authentication**: API key in `X-API-Key` header. Currently **read-only**.

**Rate limits**: Early Access (EA): 100 requests/minute; v1 stable: 10,000 requests/minute.

**Implementation**: See `remote.go` (`RemoteAPIClient`). Used for discovering consoles and sites managed via UniFi Site Manager (cloud).

## Architecture

- **Package layout**: Single package `unifi` for core logic. `main/` for CLI. `mocks/` implements `UnifiClient` for testing.
- **UnifiClient**: Interface in `types.go`; both `*Unifi` and `mocks.MockUnifi` implement it. Use it for dependency injection.
- **Config**: `Config` has `URL`, `User`, `Pass`, optional `APIKey`, `ErrorLog`/`DebugLog`, `SSLCert`, `VerifySSL`.
- **Logger**: `ErrorLog` / `DebugLog` accept `(format, args...)`; use `log.Printf` or custom.

## Adding Features

When adding features:
- Follow existing patterns (e.g. `GetDevices` / `GetClients`)
- Use `GetData` for HTTP requests
- Add tests following existing test patterns
- For new device/client types: extend `Devices` and `parseDevices`, add getters following `GetUAPs`/`GetUSWs` style, and update mocks
- Implement all `UnifiClient` interface methods when adding mocks or new client types
- Check with `var _ unifi.UnifiClient = &YourClient{}`
