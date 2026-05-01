# Potential Features to Add to unpoller/unifi

Based on analysis of API response examples, here are features that could be added.

Status key: **Done** = shipped, **Partial** = some coverage, **Open** = not started.

---

## Part A — Legacy / v2 API gaps

### A1. System Information (`/api/s/{site}/stat/sysinfo`)

**Status**: Done (`sysinfo.go`)

---

### A2. WAN Status (`/api/s/{site}/stat/status`)

**Status**: Open

**Data**: WAN interface names, states (ACTIVE/BACKUP), network groups (WAN, WAN2).

**Use case**: WAN failover monitoring.

**Implementation**: `GetWANStatus(site *Site) (*WANStatus, error)`

---

### A3. Firewall Policies (`/proxy/network/v2/api/site/{site}/firewall-policies`)

**Status**: Done (`firewall_policies.go`)

---

### A4. UPS Device List (`/api/s/{site}/stat/ups-devices`)

**Status**: Partial (PDU/UPS come from `/stat/device`; this endpoint is a lightweight selector format)

**Data**: Device MAC, image URL, label, site ID — for UI selection lists, not full stats.

**Implementation**: `GetUPSDeviceList(site *Site) ([]*UPSDeviceSelector, error)`

---

### A5. Enhanced Network Configuration (`/api/s/{site}/rest/networkconf`)

**Status**: Partial (basic `Network` struct exists; many fields missing)

**Missing fields**: WAN failover (`wan_load_balance_type`, `wan_failover_priority`), IPv6 (`ipv6_enabled`, `dhcpdv6_*`), VPN (`vpn_type`, `sdwan_remote_site_id`), DHCP advanced (`dhcpd_ntp_enabled`, `dhcpd_wins_enabled`), misc (`firewall_zone_id`, `mdns_enabled`) — `domain_name` is already present as `DomainName`.

**Implementation**: Extend `Network` struct in-place (additive fields, no breakage).

---

### A6. Port Forwarding (`/api/s/{site}/rest/portforward`)

**Status**: Open

**Data**: Port forward rules — source/destination ports, protocols, enabled flag.

**Implementation**: `GetPortForwards(site *Site) ([]*PortForward, error)`

---

### A7. SSL Certificate Info (`/api/s/{site}/stat/active`)

**Status**: Open

**Data**: Certificate chain, status, valid-from/to dates, issuer, fingerprint.

**Use case**: Certificate expiration monitoring.

**Implementation**: `GetSSLCertificate(site *Site) (*SSLCertificate, error)`

---

## Part B — Integration/v1 API coverage

The controller exposes a formally supported REST API at `/proxy/network/integration/v1/` (available since Network 9.3.43, spec at `/proxy/network/api-docs/integration.json`). It requires `X-API-Key` auth and uses **UUID site IDs** (not the legacy short name `"default"`). The library currently reads from **none** of these endpoints.

This API has 44 paths (source: `/proxy/network/api-docs/integration.json` on 10.3.58). Of those, 36 are GET operations; the remainder are write-only (POST/PUT/DELETE/PATCH). The 36 GETs collapse into 17 features (B0–B16) because list and single-item endpoints for the same resource are grouped together. Write operations are excluded — the library is read-only by design.

### B0. Infrastructure: site ID resolution (prerequisite for all B-phase work)

**Status**: Open

**Problem**: Integration/v1 per-site endpoints use a UUID `siteId` from `GET /v1/sites`, not the legacy `Site.Name` short name. These are different identifiers.

**Solution**: New type and getter:

```go
// IntegrationSite holds identity fields from GET /v1/sites.
// InternalReference matches legacy Site.Name for cross-lookup.
type IntegrationSite struct {
    ID                string `json:"id"`                // UUID — pass to all integration/v1 calls
    InternalReference string `json:"internalReference"` // equals legacy Site.Name ("default")
    Name              string `json:"name"`
}

func (u *Unifi) GetIntegrationSites() ([]*IntegrationSite, error)
```

Callers who already have `[]*Site` from `GetSites()` join on `site.Name == integrationSite.InternalReference` to get the UUID. All per-site integration/v1 getters take `*IntegrationSite`.

**Files**: new `integration_sites.go`; path constant `APIIntegrationSitesPath` in `types.go`; add to `UnifiClient` interface; add to `mocks/`.

---

### B1. Device statistics (`GET /v1/sites/{siteId}/devices/{deviceId}/statistics/latest`)

**Status**: Open

**Data**: `uptimeSec`, `cpuUtilizationPct`, `memoryUtilizationPct`, `loadAverage{1,5,15}Min`, `lastHeartbeatAt`, `nextHeartbeatAt`, per-radio `txRetriesPct`/`frequencyGHz`, per-uplink `txRateBps`/`rxRateBps`.

**Use case**: Per-device CPU/memory/uptime metrics — currently only available embedded in the large legacy device payload. This is a clean dedicated stats endpoint.

**Implementation**:
```go
GetIntegrationDeviceStats(site *IntegrationSite, deviceID string) (*IntegrationDeviceStats, error)
// Convenience bulk getter: enumerates device IDs via GET /v1/sites/{siteId}/devices first.
GetAllIntegrationDeviceStats(site *IntegrationSite) ([]*IntegrationDeviceStats, error)
```

**Files**: `integration_devices.go`

---

### B2. WiFi broadcasts (`GET /v1/sites/{siteId}/wifi/broadcasts`)

**Status**: Open

**Data**: `id`, `name`, `enabled`, `network` (VLAN reference), `securityConfiguration` (type: WPA2/WPA3/open, PSK, RADIUS profile ref), `broadcastingDeviceFilter`. Single-broadcast getter at `/{wifiBroadcastId}`.

**Use case**: SSID inventory — the library has no WiFi broadcast / SSID data today.

**Implementation**: `GetWifiBroadcasts(site *IntegrationSite) ([]*WifiBroadcast, error)`

**Files**: `wifi_broadcasts.go`

---

### B3. Firewall zones (`GET /v1/sites/{siteId}/firewall/zones`)

**Status**: Open

**Data**: `id`, `name`, `networkIds` (list of network UUIDs in the zone), `metadata` (origin: system/user).

**Use case**: The existing `GetFirewallPolicies()` returns policies that reference zone IDs — without zone data, those IDs are opaque. This completes the firewall picture.

**Implementation**: `GetFirewallZones(site *IntegrationSite) ([]*FirewallZone, error)`

**Files**: `firewall_zones.go`

---

### B4. ACL rules (`GET /v1/sites/{siteId}/acl-rules`)

**Status**: Open

**Data**: `id`, `enabled`, `name`, `description`, `action` (ALLOW/BLOCK), `enforcingDeviceFilter`, `index`, `sourceFilter`. Single-rule getter at `/{aclRuleId}`. Ordering at `/ordering`.

**Use case**: Network access-control visibility; currently no ACL data in the library.

**Implementation**: `GetACLRules(site *IntegrationSite) ([]*ACLRule, error)`

**Files**: `acl_rules.go`

---

### B5. Networks (`GET /v1/sites/{siteId}/networks`)

**Status**: Open

**Data**: `id`, `name`, `enabled`, `vlanId`, `management` (gateway/switch-managed/unmanaged), `dhcpGuarding`. More structured than the legacy `/rest/networkconf` response. Single-network getter at `/{networkId}`; references at `/{networkId}/references`.

**Use case**: Clean network inventory; complements the legacy `GetNetworks()` which remains for backward compat.

**Implementation**: `GetIntegrationNetworks(site *IntegrationSite) ([]*IntegrationNetwork, error)`

**Files**: `integration_networks.go`

---

### B6. WAN interfaces (`GET /v1/sites/{siteId}/wans`)

**Status**: Open

**Data**: `id`, `name`. Lightweight list of WAN interface identifiers.

**Use case**: Enumerate WAN interfaces to drive the existing v2 WAN enriched/SLA calls by interface ID.

**Implementation**: `GetIntegrationWANs(site *IntegrationSite) ([]*IntegrationWAN, error)`

**Files**: extend `wan.go` or new `integration_wans.go`

---

### B7. VPN servers + site-to-site tunnels

**Status**: Open

**Data**:
- `GET /v1/sites/{siteId}/vpn/servers` → `id`, `name`, `enabled`, `type` (L2TP/OpenVPN/WireGuard/UID)
- `GET /v1/sites/{siteId}/vpn/site-to-site-tunnels` → `id`, `name`, `type` (IPSec/OpenVPN/WireGuard), `metadata`

**Use case**: VPN infrastructure visibility; complements the existing `GetMagicSiteToSiteVPN()` (v2 API).

**Implementation**:
```go
GetVPNServers(site *IntegrationSite) ([]*VPNServer, error)
GetSiteToSiteTunnels(site *IntegrationSite) ([]*SiteToSiteTunnel, error)
```

**Files**: `vpn_servers.go` (or extend `vpn.go`)

---

### B8. Switching topology (LAGs, MC-LAG domains, switch stacks)

**Status**: Open

**Data**:
- `GET /v1/sites/{siteId}/switching/lags` → `id`, `type` (local/global), `members` (device ID + port indexes), `metadata`
- `GET /v1/sites/{siteId}/switching/mc-lag-domains` → `id`, `name`, `peers` (role + device ID + link ports), `lags`, `metadata`
- `GET /v1/sites/{siteId}/switching/switch-stacks` → `id`, `name`, `members` (device IDs), `lags`, `metadata`

**Use case**: Switching topology not available via any current endpoint. Required for understanding LAG/MLAG/stack relationships between switches.

**Implementation**:
```go
GetLAGs(site *IntegrationSite) ([]*LAG, error)
GetMCLAGDomains(site *IntegrationSite) ([]*MCLAGDomain, error)
GetSwitchStacks(site *IntegrationSite) ([]*SwitchStack, error)
```

**Files**: `switching.go`

---

### B9. DNS policies (`GET /v1/sites/{siteId}/dns/policies`)

**Status**: Open

**Data**: `id`, `enabled`, `type` (forward-domain/block/allow/etc.), `domain`. Single-policy getter at `/{dnsPolicyId}`.

**Use case**: DNS filtering / split-horizon visibility.

**Implementation**: `GetDNSPolicies(site *IntegrationSite) ([]*DNSPolicy, error)`

**Files**: `dns_policies.go`

---

### B10. RADIUS profiles (`GET /v1/sites/{siteId}/radius/profiles`)

**Status**: Open

**Data**: `id`, `name`, `metadata`. Reference data used by WiFi enterprise security configs.

**Implementation**: `GetRADIUSProfiles(site *IntegrationSite) ([]*RADIUSProfile, error)`

**Files**: `radius_profiles.go`

---

### B11. Traffic matching lists (`GET /v1/sites/{siteId}/traffic-matching-lists`)

**Status**: Open

**Data**: `id`, `name`, `type` (IPv4/IPv6/port). Used as references in firewall policy traffic filters.

**Implementation**: `GetTrafficMatchingLists(site *IntegrationSite) ([]*TrafficMatchingList, error)`

**Files**: `traffic_matching_lists.go`

---

### B12. Hotspot vouchers (`GET /v1/sites/{siteId}/hotspot/vouchers`)

**Status**: Open

**Data**: `id`, `code`, `name`, `authorizedGuestLimit`, `authorizedGuestCount`, `activatedAt`, `expiresAt`, `timeLimitMinutes`, `dataUsageLimitMBytes`, rate limits.

**Use case**: Guest portal monitoring — voucher usage, expiry, capacity.

**Implementation**: `GetHotspotVouchers(site *IntegrationSite) ([]*HotspotVoucher, error)`

**Files**: `hotspot_vouchers.go`

---

### B13. DPI application catalogue (`GET /v1/dpi/applications`, `GET /v1/dpi/categories`)

**Status**: Open

**Data**: `id`, `name` for each application and category. Reference tables that decode the integer IDs in existing DPI stat responses.

**Implementation**:
```go
GetDPIApplications() ([]*DPIApplication, error)
GetDPICategories() ([]*DPICategory, error)
```

**Files**: extend `dpi.go`

---

### B14. Pending devices (`GET /v1/pending-devices`)

**Status**: Open

**Data**: `macAddress`, `ipAddress`, `model`, `state`, `supported`, `firmwareVersion`, `firmwareUpdatable`, `features`.

**Use case**: Adoption queue monitoring — see what devices are waiting to be adopted.

**Implementation**: `GetPendingDevices() ([]*PendingDevice, error)`

**Files**: extend `devices.go` or new `pending_devices.go`

---

### B15. Application info (`GET /v1/info`)

**Status**: Open

**Data**: `applicationVersion` string.

**Use case**: Lightweight version check; complements `GetServerData()` which populates `ServerStatus`, after which `u.ServerStatus.MajorVersion()` returns the numeric major version (promoted to `u.MajorVersion()` via embedding, but requires `GetServerData()` to have been called first).

**Implementation**: `GetIntegrationInfo() (*IntegrationInfo, error)`

**Files**: `integration_sites.go` (same file as B0, trivially small)

---

### B16. Countries (`GET /v1/countries`)

**Status**: Open

**Data**: `code`, `name` per country.

**Use case**: Reference data for geo-based firewall policy filters (region filter). Low priority.

**Implementation**: `GetCountries() ([]*Country, error)`

**Files**: `countries.go`

---

## Part B — Implementation order

### Phase 0 (prerequisite)
- **B0** — `GetIntegrationSites()` + `IntegrationSite` type

### Phase 1 (high monitoring value, low complexity)
- **B1** — Device statistics (CPU/mem/uptime)
- **B2** — WiFi broadcasts (SSID inventory)
- **B3** — Firewall zones (completes firewall policy context)
- **B4** — ACL rules

### Phase 2 (network topology)
- **B5** — Networks (integration/v1)
- **B6** — WAN interfaces
- **B7** — VPN servers + site-to-site tunnels
- **B8** — Switching topology (LAGs, MC-LAGs, stacks)
- **B9** — DNS policies

### Phase 3 (reference and config data)
- **B10** — RADIUS profiles
- **B11** — Traffic matching lists
- **B12** — Hotspot vouchers
- **B13** — DPI app/category catalogue
- **B14** — Pending devices
- **B15** — Application info
- **B16** — Countries

---

## Part B — Implementation notes

- **Auth**: Integration/v1 requires `X-API-Key` — cookie auth returns 401. Every B-phase getter must guard with an early return when `u.APIKey == ""` (a new sentinel `ErrAPIKeyRequired` should be added). The `setHeaders` helper already sends the key when present; no transport changes needed.
- **Pagination**: Integration/v1 list responses use `{offset, limit, count, totalCount, data}` — a different envelope from the legacy `{meta, data}` that `GetData` is shaped for. Paginated getters must use `GetJSON` + a shared integration-pagination helper that loops until `offset+count >= totalCount`, rather than `GetData`.
- **Global vs. per-site**: Endpoints without a `{siteId}` segment (B13–B16) take no site parameter. Don't add `*IntegrationSite` to these getters.
- **Path constants**: All integration/v1 path constants go in `types.go` under a new `// Integration/v1 API paths` comment block, following the same naming convention (`APIIntegration*Path`).
- **Interface**: Every new getter must be added to `UnifiClient` in `types.go` and implemented in `mocks/`.
- **Write operations excluded**: POST/PUT/DELETE/PATCH are out of scope — the library is read-only by design.
- **Availability gate**: Integration/v1 requires Network 9.3.43+. Callers should check `u.ServerStatus.MajorVersion() >= 9` (after `GetServerData()`) before calling integration/v1 getters, or handle `ErrEndpointNotFound` gracefully if the endpoint is absent on older controllers.
