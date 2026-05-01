# Discovered UniFi Controller Endpoints

This document lists API endpoints observed while browsing a UniFi controller (UDM-style, with `/proxy/network` prefix). Data was captured using the XHR/fetch capture script (e.g. in repo `scripts/` or Desktop `scripts/`). Paths are relative to the controller base URL (e.g. `https://192.168.1.1`).

**Conventions**

- **Legacy API**: `/proxy/network/api/s/default/...` — classic `stat/*` and `rest/*` endpoints; responses often use `{"meta":{"rc":"ok"},"data":[...]}`.
- **v2 API**: `/proxy/network/v2/api/...` — newer site-scoped endpoints; responses vary (arrays or objects).
- **Site**: `default` appears as the site name; replace with your site name if different.

---

## Auth & session

| Method | Path | Data exposed |
|--------|------|----------------|
| POST | `/api/auth/login` | Login (request body: `username`, `password`, `token`, `rememberMe`). Response: user + session (e.g. `unique_id`, `first_name`, cookies). |
| GET | `/api/users/self` | Current user (controller-level). Returns 401 when not authenticated. |
| GET | `/proxy/users/api/v2/user/self` | Current user (v2). |
| PUT | `/proxy/network/api/self` | Update self; response `{"meta":{"rc":"ok"},"data":[]}`. |

---

## System & controller info

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/api/system` | System info (controller). |
| GET | `/api/system/syslog/settings` | Syslog settings (may 403). |
| GET | `/proxy/network/v2/api/info` | Network app API info. |
| GET | `/proxy/network/v2/api/timezones` | Timezone list. |
| GET | `/proxy/network/api/s/default/stat/sysinfo` | System info (site). |
| GET | `/proxy/network/api/s/default/stat/health` | Health summary. |
| GET | `/proxy/network/api/s/default/stat/sdn` | SDN/controller status: `enabled`, `connected`, `is_udm`, `ubic_env`, etc. |
| GET | `/proxy/network/api/s/default/get/setting` | Site settings. |
| GET | `/proxy/network/api/ui-data` | UI config (e.g. mapbox, feature flags). |
| GET | `/proxy/network/v2/api/site/default/models` | Device/site models. |
| GET | `/proxy/network/v2/api/site/default/described-features` | Feature descriptions. |
| GET | `/proxy/network/v2/api/system/event/{type}/first` | First occurrence of system event (e.g. `SETUP_COMPLETED`). |
| GET | `/proxy/network/v2/api/fingerprint_devices/0`, `/1` | Device fingerprint data. |

---

## Devices

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/stat/device` | Full device list (UAP, USW, UDM, etc.). |
| GET | `/proxy/network/api/s/default/stat/device-basic` | Basic device list: `mac`, `state`, `adopted`, `type`, `model`, `name`, etc. |
| GET | `/proxy/network/v2/api/site/default/device` | Devices (v2). |
| GET | `/proxy/network/v2/api/site/default/device/wireless-links` | Wireless links between devices (array). |
| GET | `/proxy/network/v2/api/site/default/device-tags` | Device tags. |
| GET | `/proxy/network/v2/api/site/default/apgroups` | AP groups: `_id`, `name`, `device_macs`, etc. |
| GET | `/proxy/network/api/s/default/stat/current-channel` | Current channel usage. |
| GET | `/proxy/network/api/s/default/stat/ccode` | Country code. |

---

## Clients

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/clients/active` | Active clients. |
| GET | `/proxy/network/v2/api/site/default/clients/history` | Client history (array). |
| POST | `/proxy/network/v2/api/site/default/clients/metadata` | Client metadata: fingerprint, hostname, mac, name, wlanconf_id, etc. |

---

## Networks & LAN

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/networkconf` | Network configs (VLANs, etc.). |
| GET | `/proxy/network/v2/api/site/default/lan/enriched-configuration` | LAN enriched config. |
| GET | `/proxy/network/v2/api/site/default/lan/defaults` | LAN defaults. |
| GET | `/proxy/network/v2/api/site/default/lan/mdns` | mDNS config. |
| GET | `/proxy/network/v2/api/site/default/global/config/network` | Global network config: `default_security_posture`, `mdns_enabled_for`, etc. |
| GET | `/proxy/network/v2/api/site/default/object-oriented-network-configs` | Object-oriented network configs (array). |
| GET | `/proxy/network/v2/api/site/default/network-members-groups` | Network member groups. |
| GET | `/proxy/network/v2/api/site/default/excluded-ips/` | Excluded IPs: `excluded_ip_client_info`, `unidentified_excluded_ip_info`. |

---

## DHCP & DNS

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/dhcpoption` | DHCP options (`meta`/`data`). |
| GET | `/proxy/network/v2/api/site/default/active-leases` | DHCP active leases. |
| GET | `/proxy/network/v2/api/site/default/static-dns` | Static DNS entries. |
| GET | `/proxy/network/v2/api/site/default/static-dns/devices` | Static DNS devices. |

---

## WLAN & wireless

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/wlanconf` | WLAN configs. |
| GET | `/proxy/network/v2/api/site/default/wlan/enriched-configuration` | WLAN enriched config. |
| GET | `/proxy/network/v2/api/site/default/wlan/defaults` | WLAN defaults. |
| GET | `/proxy/network/v2/api/site/default/wlan-capabilities` | e.g. `6ghz_band_supported`, `wpa3_supported`. |
| GET | `/proxy/network/v2/api/site/default/wifi-stats/details` | WiFi stats details. |
| GET | `/proxy/network/v2/api/site/default/wifi-stats/radios` | Radio details: `bytes`, `client_signal_avg`, `device_mac`, `interference_avg`, `utilization_avg`, etc. |
| GET | `/proxy/network/v2/api/site/default/wifiman` | WiFiman data. |
| POST | `/proxy/network/v2/api/site/default/wifi-connectivity/events` | WiFi connectivity events (paginated). |
| POST | `/proxy/network/v2/api/site/default/wifi-connectivity/events/filter-data` | Filter options for connectivity events. |
| POST | `/proxy/network/v2/api/site/default/wifi-connectivity/roaming/topology` | Roaming topology: `clients`, `edges`, `vertices`. |

---

## WAN & load balancing

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/wan/networkgroups` | WAN network groups: id, name, `load_balancing_mode`, `port_info`, `priority`, `uptime`, etc. |
| GET | `/proxy/network/v2/api/site/default/wan/enriched-configuration` | WAN enriched config. |
| GET | `/proxy/network/v2/api/site/default/wan/defaults` | WAN defaults. |
| GET | `/proxy/network/v2/api/site/default/wan/load-balancing/configuration` | Load balancing config: `mode`, `wan_interfaces` (priority, weight). |
| GET | `/proxy/network/v2/api/site/default/wan/load-balancing/status` | Per-WAN state: `ACTIVE`, `BACKUP`, etc. |
| GET | `/proxy/network/v2/api/site/default/wan/magic/configuration` | Magic WAN config (data usage, enabled). |
| GET | `/proxy/network/v2/api/site/default/wan/magic/subscription` | Magic subscription: `subscribed`, `traffic_usage`, etc. |
| GET | `/proxy/network/v2/api/site/default/wan-slas` | WAN SLAs (array). |

---

## Firewall & security

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/firewallgroup` | Firewall groups. |
| GET | `/proxy/network/api/s/default/rest/firewallrule` | Firewall rules. |
| GET | `/proxy/network/v2/api/site/default/firewall/zone` | Firewall zones. |
| GET | `/proxy/network/v2/api/site/default/firewall/zone-matrix` | Zone matrix. |
| GET | `/proxy/network/v2/api/site/default/firewall-policies` | Firewall policies. |
| GET | `/proxy/network/v2/api/site/default/firewall-app-blocks` | App blocks (array). |
| GET | `/proxy/network/v2/api/site/default/acl-rules` | ACL rules (array). |
| GET | `/proxy/network/v2/api/site/default/content-filtering` | Content filtering (array). |
| GET | `/proxy/network/v2/api/site/default/content-filtering/categories` | Content filtering categories. |
| GET | `/proxy/network/v2/api/site/default/ssl-inspection/setting` | SSL inspection setting (state, certs). |
| GET | `/proxy/network/v2/api/site/default/ssl-inspection/setting/defaults` | SSL inspection defaults. |
| GET | `/proxy/network/v2/api/site/default/ssl-inspection/certificates` | SSL inspection certificates. |
| GET | `/proxy/network/v2/api/site/default/ssl-inspection/certificates/active` | Active certificate. |

---

## Routing, NAT, VPN

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/routing` | Routing table. |
| GET | `/proxy/network/api/s/default/rest/portforward` | Port forward rules. |
| GET | `/proxy/network/api/s/default/stat/portforward` | Port forward stat. |
| GET | `/proxy/network/v2/api/site/default/trafficrules` | Traffic rules (array). |
| GET | `/proxy/network/v2/api/site/default/trafficroutes` | Traffic routes (array). |
| GET | `/proxy/network/v2/api/site/default/nat` | NAT config (array). |
| GET | `/proxy/network/v2/api/site/default/vpn/connections` | VPN connections: `connections[]`. |
| GET | `/proxy/network/v2/api/site/default/magicsitetositevpn/configs` | Magic Site-to-Site VPN (may 403). |
| GET | `/proxy/network/v2/api/site/default/bgp/config/all` | BGP config (array). |
| GET | `/proxy/network/v2/api/site/default/ospf/router` | OSPF router (array). |
| GET | `/proxy/network/v2/api/site/default/wireguard/users` | WireGuard users (array). |

---

## DPI, QoS, RADIUS

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/dpiapp` | DPI apps. |
| GET | `/proxy/network/api/s/default/rest/dpigroup` | DPI groups (name, site_id, _id). |
| GET | `/proxy/network/v2/api/site/default/qos-rules` | QoS rules: destination apps/ports, schedule, objective, etc. |
| GET | `/proxy/network/api/s/default/rest/radiusprofile` | RADIUS profiles (auth/acct servers, vlan, etc.). |
| GET | `/proxy/network/v2/api/site/default/radius/users` | RADIUS users: `_id`, `name`, `vlan`, `tunnel_type`, etc. |

---

## Traffic & flows

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/traffic` | Traffic data. |
| POST | `/proxy/network/v2/api/site/default/traffic-flows` | Traffic flows (paginated/filtered). |
| GET | `/proxy/network/v2/api/site/default/traffic-flows/filter-data` | Filter options for traffic flows. |
| GET | `/proxy/network/v2/api/site/default/country-traffic` | Country-level traffic. |
| POST | `/proxy/network/v2/api/site/default/app-traffic-rate` | App traffic rate. |

---

## System log & events

| Method | Path | Data exposed |
|--------|------|----------------|
| POST | `/proxy/network/v2/api/site/default/system-log/critical` | Critical system log entries (array). |
| POST | `/proxy/network/v2/api/site/default/system-log/all` | Full system log (paginated). |
| POST | `/proxy/network/v2/api/site/default/system-log/count` | Log count. |
| POST | `/proxy/network/v2/api/site/default/system-log/filter-data` | Log filter options. |
| GET | `/proxy/network/v2/api/site/default/system-log/setting` | System log settings. |

---

## Dashboard & topology

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/aggregated-dashboard` | Dashboard aggregates. |
| GET | `/proxy/network/v2/api/site/default/topology` | Topology (devices/links). |
| GET | `/proxy/network/api/s/default/stat/widget/warnings` | Widget warnings: upgradable devices, firmware status, EOL count, etc. |
| GET | `/proxy/network/v2/api/site/default/stacking` | Stacking info (array). |
| GET | `/proxy/network/v2/api/site/default/mclag-groups` | MC-LAG groups (array). |

---

## Ports & switching

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/api/s/default/rest/portconf` | Port configs. |
| GET | `/proxy/network/v2/api/site/default/ports/port-anomalies` | Port anomalies (array). |
| POST | `/proxy/network/v2/api/site/default/ports/mac-tables` | MAC tables per port: `mac`, `ports[].port_idx`, `mac_table[]`. |

---

## Settings (defaults & config)

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/settings/connectivity/defaults` | Connectivity defaults (uplink_type, etc.). |
| GET | `/proxy/network/v2/api/site/default/settings/global_switch/defaults` | Global switch: STP, DHCP snoop, dot1x, jumboframe, etc. |
| GET | `/proxy/network/v2/api/site/default/settings/usg/defaults` | USG defaults (timeouts, DNS verification). |
| GET | `/proxy/network/v2/api/site/default/settings/ntp/defaults` | NTP servers. |
| GET | `/proxy/network/v2/api/site/default/settings/doh/defaults` | DoH: `server_names`, `state`. |
| GET | `/proxy/network/v2/api/site/default/settings/doh/available-server-names` | DoH server list. |
| GET | `/proxy/network/v2/api/site/default/settings/netflow/defaults` | Netflow: enabled, port, sampling, version. |
| GET | `/proxy/network/v2/api/site/default/settings/wifiai/defaults` | WiFi AI defaults. |
| GET | `/proxy/network/v2/api/site/default/settings/element_adopt/defaults` | Element adopt (e.g. `enabled`). |
| GET | `/proxy/network/v2/api/site/default/settings/ips/available-categories` | IPS categories. |
| GET | `/proxy/network/v2/api/site/default/settings/ips/advanced-filtering-defaults` | IPS advanced filtering defaults. |

---

## Other site features

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/network/v2/api/site/default/hotspot/info` | Hotspot info (e.g. `show_hotspot_manager`). |
| GET | `/proxy/network/api/stat/s/default/hotspotpackages` | Hotspot packages. |
| GET | `/proxy/network/v2/api/site/default/ulp/users_groups` | ULP users and groups. |
| GET | `/proxy/network/v2/api/site/default/vendor-ids` | Vendor ID list (OUI prefixes). |
| GET | `/proxy/network/v2/api/site/default/site-feature-migration` | Feature migration timestamps. |
| GET | `/proxy/network/v2/api/site/default/shadowmode/info` | Shadow mode info. |
| GET | `/proxy/network/v2/api/site/default/features/{FEATURE}/exists` | Feature flag (e.g. `UPS_ADOPTED`, `AFC_CAPABLE_AP_ADOPTED`) → `feature_exists`. |
| GET | `/proxy/network/api/s/default/rest/dynamicdns` | Dynamic DNS configs. |
| GET | `/proxy/network/api/s/default/rest/scheduletask` | Schedule tasks. |
| GET | `/proxy/network/api/s/default/rest/usergroup` | User groups (QoS, name, _id). |
| GET | `/proxy/network/api/s/default/stat/rogueap` | Rogue AP list. |

---

## Users & access (UniFi OS / proxy/users)

| Method | Path | Data exposed |
|--------|------|----------------|
| GET | `/proxy/users/api/v2/info` | Users API info. |
| GET | `/proxy/users/api/v2/org` | Organization: org_id, name, domain, subdomain. |
| GET | `/proxy/users/api/v2/ucore/controllers` | Controllers: type, name, version, port, installState. |
| GET | `/proxy/users/api/v2/user_groups` | User groups (unique_id, name, up_ids). |
| GET | `/proxy/users/api/v2/users/admin/uos` | Admin UOS. |
| GET | `/proxy/users/api/v2/users/search` | User search. |
| GET | `/proxy/users/api/v2/identity/workspace/info` | Workspace: name, invitation_expire_days, etc. |
| GET | `/proxy/users/api/v2/identity/network/info` | Identity network info. |
| GET | `/proxy/users/api/v2/identity/ownership` | Identity ownership. |
| GET | `/proxy/users/api/v2/custom_roles` | Custom roles. |
| GET | `/proxy/users/api/v2/permission_manifests` | Permission manifests. |
| GET | `/proxy/users/api/v2/csv/capabilities` | CSV export capabilities. |
| GET | `/proxy/users/access/api/v2/access/feature` | Access feature (may 502). |
| GET | `/proxy/users/access/api/v2/access/info` | Access info (may 502). |
| GET | `/proxy/users/access/api/v2/settings` | Access settings (may 502). |

---

## App assets (non-API)

| Path pattern | Data exposed |
|--------------|----------------|
| `/app-assets/network/react/data/locales/en/*.json` | UI locale strings (activity, clients, devices, logs, settings, etc.). Filenames include content hashes. |

---

## External / third-party

| Path / host | Data exposed |
|-------------|----------------|
| `https://static.ui.com/fingerprint/0/devicelist.json` | Device list for fingerprinting. |
| `https://static.ui.com/fingerprint/ui/public.json` | Public fingerprint config. |
| `https://api.maptiler.com/...` (or `/tiles/countries/tiles.json`) | Map tiles. |
| Firebase (`firebase.googleapis.com`, `firebaseinstallations.googleapis.com`, `firebaseremoteconfig.googleapis.com`) | Identity/remote config. |

---

---

## Network 10.x API changes (tested on 10.3.58, UDM-Pro)

### Removed endpoints (404 on Network 10.x)

| Endpoint | Previously used for | Status |
|----------|---------------------|--------|
| `/proxy/network/api/s/default/stat/event` | UniFi Events (GetEvents) | **Removed** — no replacement in integration/v1 |
| `/proxy/network/api/s/default/stat/ips/event` | IDS/IPS Events (GetIDS) | **Removed** — no replacement in integration/v1 |

**Workaround**: Use `save_syslog` (v2 system-log POST API) instead of `save_events`. IDS events have no current replacement.

### Status endpoint regression

`/proxy/network/status` no longer returns `server_version` in Network 10.x (field is absent). Use `/proxy/network/api/s/default/stat/sysinfo` → `data[0].version` for controller version detection instead.

### Valid system-log categories (Network 10.x)

The `/proxy/network/v2/api/site/default/system-log/all` endpoint (POST) only accepts these `categories` values on Network 10.x (see [unpoller/unifi#198](https://github.com/unpoller/unifi/issues/198)):

`SECURITY`, `UNIFI_DEVICES`, `SOFTWARE_UPDATES`, `VPN`, `POWER`, `UNIFI_ETHERNET_PORTS`, `CLIENT_DEVICES`, `UNKNOWN`, `AUDIT`, `INTERNET_AND_WAN`

Values `MONITORING` and `SYSTEM` cause a `400 Bad Request`.

---

## Official integration/v1 API (Network 9.3.43+)

A formally supported REST API with OpenAPI spec. Requires `X-API-Key` header auth (not cookie auth). Uses UUID-format site IDs (not site names like `"default"`). Spec available at `/proxy/network/api-docs/integration.json`. Full versioned specs at <https://github.com/beezly/unifi-apis>.

Base path: `/proxy/network/integration/v1/`

### Endpoints (as of 10.3.58, from `/proxy/network/api-docs/integration.json`)

This API supports both reads and writes. `siteId` is the UUID from `GET /v1/sites` (not the short name `"default"`).

| Methods | Path | Summary |
|---------|------|---------|
| GET | `/v1/info` | Get Application Info |
| GET | `/v1/countries` | List Countries |
| GET | `/v1/dpi/categories` | List DPI Application Categories |
| GET | `/v1/dpi/applications` | List DPI Applications |
| GET | `/v1/pending-devices` | List Devices Pending Adoption |
| GET | `/v1/sites` | List Local Sites |
| GET, POST | `/v1/sites/{siteId}/acl-rules` | List / Create ACL Rules |
| GET, PUT | `/v1/sites/{siteId}/acl-rules/ordering` | Get / Reorder ACL Rule Ordering |
| GET, PUT, DELETE | `/v1/sites/{siteId}/acl-rules/{aclRuleId}` | Get / Update / Delete ACL Rule |
| GET | `/v1/sites/{siteId}/clients` | List Connected Clients |
| GET | `/v1/sites/{siteId}/clients/{clientId}` | Get Connected Client Details |
| POST | `/v1/sites/{siteId}/clients/{clientId}/actions` | Execute Client Action |
| GET | `/v1/sites/{siteId}/device-tags` | List Device Tags |
| GET, POST | `/v1/sites/{siteId}/devices` | List / Adopt Devices |
| GET, DELETE | `/v1/sites/{siteId}/devices/{deviceId}` | Get / Unadopt Device |
| POST | `/v1/sites/{siteId}/devices/{deviceId}/actions` | Execute Device Action |
| POST | `/v1/sites/{siteId}/devices/{deviceId}/interfaces/ports/{portIdx}/actions` | Execute Port Action |
| GET | `/v1/sites/{siteId}/devices/{deviceId}/statistics/latest` | Get Latest Device Statistics |
| GET, POST | `/v1/sites/{siteId}/dns/policies` | List / Create DNS Policies |
| GET, PUT, DELETE | `/v1/sites/{siteId}/dns/policies/{dnsPolicyId}` | Get / Update / Delete DNS Policy |
| GET, POST | `/v1/sites/{siteId}/firewall/policies` | List / Create Firewall Policies |
| GET, PUT | `/v1/sites/{siteId}/firewall/policies/ordering` | Get / Reorder Firewall Policies |
| GET, PUT, DELETE, PATCH | `/v1/sites/{siteId}/firewall/policies/{firewallPolicyId}` | Get / Update / Delete / Patch Firewall Policy |
| GET, POST | `/v1/sites/{siteId}/firewall/zones` | List / Create Firewall Zones |
| GET, PUT, DELETE | `/v1/sites/{siteId}/firewall/zones/{firewallZoneId}` | Get / Update / Delete Firewall Zone |
| GET, POST, DELETE | `/v1/sites/{siteId}/hotspot/vouchers` | List / Generate / Delete Vouchers |
| GET, DELETE | `/v1/sites/{siteId}/hotspot/vouchers/{voucherId}` | Get / Delete Voucher |
| GET, POST | `/v1/sites/{siteId}/networks` | List / Create Networks |
| GET, PUT, DELETE | `/v1/sites/{siteId}/networks/{networkId}` | Get / Update / Delete Network |
| GET | `/v1/sites/{siteId}/networks/{networkId}/references` | Get Network References |
| GET | `/v1/sites/{siteId}/radius/profiles` | List RADIUS Profiles |
| GET | `/v1/sites/{siteId}/switching/lags` | List LAGs |
| GET | `/v1/sites/{siteId}/switching/lags/{lagId}` | Get LAG Details |
| GET | `/v1/sites/{siteId}/switching/mc-lag-domains` | List MC-LAG Domains |
| GET | `/v1/sites/{siteId}/switching/mc-lag-domains/{mcLagDomainId}` | Get MC-LAG Domain |
| GET | `/v1/sites/{siteId}/switching/switch-stacks` | List Switch Stacks |
| GET | `/v1/sites/{siteId}/switching/switch-stacks/{switchStackId}` | Get Switch Stack |
| GET, POST | `/v1/sites/{siteId}/traffic-matching-lists` | List / Create Traffic Matching Lists |
| GET, PUT, DELETE | `/v1/sites/{siteId}/traffic-matching-lists/{id}` | Get / Update / Delete Traffic Matching List |
| GET | `/v1/sites/{siteId}/vpn/servers` | List VPN Servers |
| GET | `/v1/sites/{siteId}/vpn/site-to-site-tunnels` | List Site-To-Site VPN Tunnels |
| GET | `/v1/sites/{siteId}/wans` | List WAN Interfaces |
| GET, POST | `/v1/sites/{siteId}/wifi/broadcasts` | List / Create WiFi Broadcasts |
| GET, PUT, DELETE | `/v1/sites/{siteId}/wifi/broadcasts/{wifiBroadcastId}` | Get / Update / Delete WiFi Broadcast |

**Not available in integration/v1**: events, IDS/IPS, alarms, anomalies, DPI stats, rogue APs, system log.

---

*Generated from a capture session; endpoints may vary by controller version and role. Use the capture script to record your own session and extend this list.*
