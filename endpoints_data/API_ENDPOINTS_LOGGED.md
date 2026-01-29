# API Endpoints Logged from Browser Session

This file logs all UniFi API endpoints discovered while navigating the UniFi interface.

I logged an entire session from login, clicked around each page, and settings menu, i may have missed some items feel free to add to this
---

## API Endpoints Discovered

### SSO/Authentication APIs
- `GET https://sso.ui.com/api/sso/v1/user/self?settings=true`
- `GET https://sso.ui.com/api/sso/v1/user/self?include=settings,roles`
- `GET https://sso.ui.com/api/sso/v1/account-chooser`
- `GET https://sso.ui.com/api/sso/v1/search/email/{email}` (e.g., `briangates1998@gmail.com`)
- `POST https://sso.ui.com/api/sso/v1/passwordless/options`
- `POST https://sso.ui.com/api/sso/v1/passwordless/login`
- `GET https://sso.ui.com/api/sso/v1/user/self/settings`
- `GET https://sso.ui.com/api/sso/v1/user/self/trusted_devices`
- `GET https://sso.ui.com/api/sso/v1/user/self/mfa`
- `GET https://sso.ui.com/api/sso/v1/user/self/roles`
- `GET https://sso.ui.com/api/sso/v1/user/self/sessions`

### Enterprise/Organization APIs
- `GET https://enterprise.svc.ui.com/api/v1/organizations?status=active`
- `GET https://enterprise.svc.ui.com/api/v2/organizations`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/domains/{domain}`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/{orgId}/hosts?type=uos&withUserData=true&pageSize=350`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/{orgId}/hosts?type=netserver&withUserData=true&pageSize=350`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/{orgId}/sites?includeStatistics=true`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/{orgId}/hosts/{hostId}?type=uos&withUserData=true`
- `GET https://enterprise.svc.ui.com/api/v2/organizations/{orgId}/identityhub`
- `GET https://enterprise.svc.ui.com/api/v1/organizations/{orgId}/credential-server`

### Cloud Access APIs
- `POST https://cloudaccess.svc.ui.com/create-credentials`
- `GET https://cloudaccess.svc.ui.com/colocated-hosts`
- `GET https://cloudaccess.svc.ui.com/vantage-points`
- `GET https://cloudaccess.svc.ui.com/api/v1/hosts?type=netserver&withUserData=true&pageSize=350`

### Billing APIs
- `GET https://billing.svc.ui.com/api/v2/customers/self/payment-methods?type=card&limit=100&withIsInUse=true`
- `GET https://billing.svc.ui.com/api/v2/customers/self`
- `GET https://billing.svc.ui.com/api/v2/customers/self/tax-ids`
- `GET https://billing.svc.ui.com/api/v2/customers/self/billing-details`
- `GET https://billing.svc.ui.com/api/v2/products?limit=100`
- `GET https://billing.svc.ui.com/api/v2/customers/self/subscriptions?limit=100&expand[]=data.latest_invoice`
- `GET https://billing.svc.ui.com/api/v2/customers/self/invoices?status=open&limit=100&expand[]=data.charge`
- `GET https://billing.svc.ui.com/api/v2/customers/self/invoices?status=draft&limit=100`
- `GET https://billing.svc.ui.com/api/v2/accounts/business-profile`
- `GET https://billing.svc.ui.com/api/v2/accounts/tax-ids`
- `GET https://billing.svc.ui.com/api/v2/customers/self/invoices/available-years`

### IPS (Intrusion Prevention System) APIs
- `GET https://ips.svc.ui.com/api/v2/subscriptions/devices?subscription_item=wan_magic`
- `GET https://ips.svc.ui.com/api/v2/subscriptions/devices?subscription_item=rule_set_pro`

### Other Service APIs
- `GET https://status.ui.com/api/v2/scheduled-maintenances/upcoming.json`
- `GET https://status.ui.com/api/v2/scheduled-maintenances/active.json`
- `GET https://status.ui.com/api/v2/incidents/unresolved.json`
- `GET https://cap.svc.ui.com/api/v1/devices/integrations/ssoa`
- `GET https://umr.svc.ui.com/api/v1/users/integrations/ssoa/subscriptions`
- `GET https://account.svc.ui.com/zendesk/user`
- `GET https://account.svc.ui.com/chat/queue`
- `GET https://megaphone.svc.ui.com/v1/announcements`
- `GET https://ues.svc.ui.com/api/v1/info/subscriptions`
- `GET https://setup.svc.ui.com/api/v1/site-managers/devices?order_by=added_at&order=desc`
- `GET https://config.svc.ui.com/cloudAccessConfig.json`
- `GET https://geo.svc.ui.com/geo`
- `GET https://static.ui.com/fingerprint/ui/public.json`

### Local Controller Endpoints (briangates.ui.com)
- `GET https://briangates.ui.com/consoles/{hostId}/network/{site}` - Main UI routes (not direct API)
- `GET https://briangates.ui.com/ping` - Health check
- `HEAD https://briangates.ui.com/` - Health check

**Note**: The cloud UI (briangates.ui.com) uses WebSocket connections and internal APIs that aren't visible in HTTP requests. The actual local controller API endpoints (`/api/s/{site}/stat/device`, etc.) are likely accessed via WebSocket or internal API calls that don't show up in the network tab.

---

## Local Controller API Endpoints (Direct Connection)

### Session 2 - Direct Local Controller Access
- **URL**: `https://192.168.1.1/login`
- **Date**: January 28, 2026
- **Method**: Direct browser access with XHR capture
- **Total unique API endpoints discovered**: 170

### Authentication & System APIs
- `GET /api/users/self` - Current user info
- `POST /api/auth/login` - Login (UDM Pro style)
- `GET /api/system` - System information
- `GET /api/system/syslog/settings` - Syslog settings
- `GET /api/firmware/update` - Firmware update info
- `POST /api/controllers/checkUpdates` - Check for updates
- `GET /api/cloud/backup/settings/list` - Cloud backup settings
- `GET /api/integrations` - Integrations list
- `GET /api/availableIntegrations` - Available integrations
- `GET /api/notifications/settings` - Notification settings
- `GET /api/userCertificates` - User certificates
- `GET /api/device/groups` - Device groups

### Alarms API (v2)
- `GET /api/v2/alarms/profiles` - Alarm profiles
- `GET /api/v2/alarms/network` - Network alarms
- `GET /api/v2/alarms/network/manifest` - Alarm manifest
- `GET /api/v2/alarms/me/notifications/mute` - Muted notifications
- `GET /proxy/network/v2/api/alarm-manager/scope/sites` - Alarm scope: sites
- `GET /proxy/network/v2/api/alarm-manager/scope/devices` - Alarm scope: devices
- `GET /proxy/network/v2/api/alarm-manager/scope/clients` - Alarm scope: clients
- `GET /proxy/network/v2/api/alarm-manager/scope/poe-switch-devices` - Alarm scope: POE switches
- `GET /proxy/network/v2/api/alarm-manager/scope/ups-devices` - Alarm scope: UPS devices
- `GET /proxy/network/v2/api/alarm-manager/scope/psu-devices` - Alarm scope: PSU devices

### Network API (v1 - Legacy)
- `GET /proxy/network/api/s/{site}/stat/device` - Devices (already implemented)
- `GET /proxy/network/api/s/{site}/stat/device-basic` - Basic device info
- `GET /proxy/network/api/s/{site}/stat/health` - Health status
- `GET /proxy/network/api/s/{site}/stat/rogueap` - Rogue APs
- `GET /proxy/network/api/s/{site}/stat/portforward` - Port forwarding stats
- `GET /proxy/network/api/s/{site}/stat/current-channel` - Current channel info
- `GET /proxy/network/api/s/{site}/stat/spectrum-scan/{mac}` - Spectrum scan results
- `GET /proxy/network/api/s/{site}/stat/sysinfo` - System info
- `GET /proxy/network/api/s/{site}/stat/widget/warnings` - Widget warnings
- `GET /proxy/network/api/s/{site}/stat/widget/subnet-suggest` - Subnet suggestions
- `GET /proxy/network/api/s/{site}/stat/ccode` - Country code
- `GET /proxy/network/api/s/{site}/stat/sdn` - SDN info
- `GET /proxy/network/api/s/{site}/get/setting` - Site settings
- `GET /proxy/network/api/self/sites` - User's sites
- `GET /proxy/network/api/ui-data` - UI data
- `PUT /proxy/network/api/self` - Update self

### Network Configuration (REST)
- `GET /proxy/network/api/s/{site}/rest/networkconf` - Network configs (already implemented)
- `GET /proxy/network/api/s/{site}/rest/wlanconf` - WLAN configs
- `GET /proxy/network/api/s/{site}/rest/firewallrule` - Firewall rules
- `GET /proxy/network/api/s/{site}/rest/firewallgroup` - Firewall groups
- `GET /proxy/network/api/s/{site}/rest/portforward` - Port forwarding configs
- `GET /proxy/network/api/s/{site}/rest/portconf` - Port configs
- `GET /proxy/network/api/s/{site}/rest/routing` - Routing configs
- `GET /proxy/network/api/s/{site}/rest/usergroup` - User groups
- `GET /proxy/network/api/s/{site}/rest/radiusprofile` - RADIUS profiles
- `GET /proxy/network/api/s/{site}/rest/dhcpoption` - DHCP options
- `GET /proxy/network/api/s/{site}/rest/dynamicdns` - Dynamic DNS configs
- `GET /proxy/network/api/s/{site}/rest/dpiapp` - DPI apps
- `GET /proxy/network/api/s/{site}/rest/dpigroup` - DPI groups
- `GET /proxy/network/api/s/{site}/rest/scheduletask` - Scheduled tasks

### Network API (v2 - Modern)
- `GET /proxy/network/v2/api/info` - API info
- `GET /proxy/network/v2/api/info?preferred_site_name={site}` - Site-specific info
- `GET /proxy/network/v2/api/site/{site}/models` - Device models
- `GET /proxy/network/v2/api/site/{site}/aggregated-dashboard` - Dashboard data
- `GET /proxy/network/v2/api/site/{site}/described-features` - Feature descriptions
- `GET /proxy/network/v2/api/system/event/SETUP_COMPLETED/first` - Setup completion event
- `GET /proxy/network/v2/api/timezones` - Available timezones
- `GET /proxy/network/v2/api/log-levels/defaults` - Log level defaults
- `GET /proxy/network/v2/api/settings/super-mgmt/defaults` - Super management defaults
- `GET /proxy/network/v2/api/unifi-core/general-info` - UniFi Core general info

### Clients API (v2)
- `GET /proxy/network/v2/api/site/{site}/clients/active` - Active clients
- `GET /proxy/network/v2/api/site/{site}/clients/history` - Client history
- `POST /proxy/network/v2/api/site/{site}/clients/metadata` - Client metadata

### Devices API (v2)
- `GET /proxy/network/v2/api/site/{site}/device` - Devices (enhanced)
- `GET /proxy/network/v2/api/site/{site}/device-tags` - Device tags
- `GET /proxy/network/v2/api/site/{site}/device/wireless-links` - Wireless links
- `GET /proxy/network/v2/api/site/{site}/device/{mac}/replacement-candidates` - Device replacement candidates
- `GET /proxy/network/v2/api/fingerprint_devices/{id}` - Device fingerprinting

### WAN/LAN Configuration
- `GET /proxy/network/v2/api/site/{site}/wan/enriched-configuration` - WAN config
- `GET /proxy/network/v2/api/site/{site}/wan/defaults` - WAN defaults
- `GET /proxy/network/v2/api/site/{site}/wan/WAN/isp-status` - ISP status
- `GET /proxy/network/v2/api/site/{site}/wan/networkgroups` - WAN network groups
- `GET /proxy/network/v2/api/site/{site}/wan/load-balancing/configuration` - Load balancing config
- `GET /proxy/network/v2/api/site/{site}/wan/load-balancing/status` - Load balancing status
- `GET /proxy/network/v2/api/site/{site}/wan/magic/configuration` - Magic WAN config
- `GET /proxy/network/v2/api/site/{site}/wan/magic/subscription` - Magic WAN subscription
- `GET /proxy/network/v2/api/site/{site}/wan-slas` - WAN SLAs
- `GET /proxy/network/v2/api/site/{site}/lan/enriched-configuration` - LAN config
- `GET /proxy/network/v2/api/site/{site}/lan/defaults` - LAN defaults
- `GET /proxy/network/v2/api/site/{site}/lan/mdns` - mDNS settings

### WiFi/WLAN API
- `GET /proxy/network/v2/api/site/{site}/wlan/enriched-configuration` - WLAN config
- `GET /proxy/network/v2/api/site/{site}/wlan/defaults` - WLAN defaults
- `GET /proxy/network/v2/api/site/{site}/wlan-capabilities` - WLAN capabilities
- `GET /proxy/network/v2/api/site/{site}/apgroups` - AP groups
- `GET /proxy/network/v2/api/site/{site}/wifi-stats/details` - WiFi stats details
- `GET /proxy/network/v2/api/site/{site}/wifi-stats/radios` - Radio stats
- `GET /proxy/network/v2/api/site/{site}/wifi-connectivity` - WiFi connectivity
- `POST /proxy/network/v2/api/site/{site}/wifi-connectivity/events` - WiFi connectivity events
- `POST /proxy/network/v2/api/site/{site}/wifi-connectivity/events/filter-data` - Filter data
- `POST /proxy/network/v2/api/site/{site}/wifi-connectivity/roaming/topology` - Roaming topology
- `GET /proxy/network/v2/api/site/{site}/wifiman` - WiFiMan info
- `GET /proxy/network/v2/api/site/{site}/radio-ai/isolation-matrix` - Radio AI isolation matrix

### Traffic & DPI API
- `GET /proxy/network/v2/api/site/{site}/traffic` - Traffic stats (already implemented)
- `POST /proxy/network/v2/api/site/{site}/app-traffic-rate` - App traffic rates
- `GET /proxy/network/v2/api/site/{site}/country-traffic` - Country traffic (already implemented)
- `POST /proxy/network/v2/api/site/{site}/traffic-flows` - Traffic flows
- `GET /proxy/network/v2/api/site/{site}/traffic-flows/filter-data` - Traffic flow filters
- `GET /proxy/network/v2/api/site/{site}/trafficrules` - Traffic rules
- `GET /proxy/network/v2/api/site/{site}/trafficroutes` - Traffic routes

### Firewall & Security API
- `GET /proxy/network/v2/api/site/{site}/firewall/zone` - Firewall zones
- `GET /proxy/network/v2/api/site/{site}/firewall/zone-matrix` - Zone matrix
- `GET /proxy/network/v2/api/site/{site}/firewall-policies` - Firewall policies
- `GET /proxy/network/v2/api/site/{site}/firewall-app-blocks` - App blocks
- `GET /proxy/network/v2/api/site/{site}/acl-rules` - ACL rules
- `GET /proxy/network/v2/api/site/{site}/nat` - NAT rules
- `GET /proxy/network/v2/api/site/{site}/qos-rules` - QoS rules
- `GET /proxy/network/v2/api/site/{site}/content-filtering` - Content filtering
- `GET /proxy/network/v2/api/site/{site}/content-filtering/categories` - Content filtering categories
- `GET /proxy/network/v2/api/site/{site}/ssl-inspection/certificates` - SSL inspection certs
- `GET /proxy/network/v2/api/site/{site}/ssl-inspection/certificates/active` - Active SSL certs
- `GET /proxy/network/v2/api/site/{site}/ssl-inspection/setting` - SSL inspection settings
- `GET /proxy/network/v2/api/site/{site}/ssl-inspection/setting/defaults` - SSL inspection defaults
- `GET /proxy/network/v2/api/site/{site}/settings/ips/advanced-filtering-defaults` - IPS advanced filtering
- `GET /proxy/network/v2/api/site/{site}/settings/ips/available-categories` - IPS categories
- `GET /proxy/network/v2/api/site/{site}/settings/doh/defaults` - DoH defaults
- `GET /proxy/network/v2/api/site/{site}/settings/doh/available-server-names` - DoH servers

### VPN API
- `GET /proxy/network/v2/api/site/{site}/vpn/connections` - VPN connections
- `GET /proxy/network/v2/api/site/{site}/magicsitetositevpn/configs` - Magic Site-to-Site VPN
- `GET /proxy/network/v2/api/site/{site}/wireguard/users` - WireGuard users
- `POST /proxy/network/v2/api/site/{site}/teleport/invitation-history` - Teleport invitations

### Network Management
- `GET /proxy/network/v2/api/site/{site}/global/config/network` - Global network config
- `GET /proxy/network/v2/api/site/{site}/network-members-groups` - Network member groups
- `GET /proxy/network/v2/api/site/{site}/object-oriented-network-configs` - OO network configs
- `GET /proxy/network/v2/api/site/{site}/static-dns` - Static DNS
- `GET /proxy/network/v2/api/site/{site}/static-dns/devices` - Static DNS devices
- `GET /proxy/network/v2/api/site/{site}/active-leases` - Active DHCP leases
- `GET /proxy/network/v2/api/site/{site}/excluded-ips/` - Excluded IPs
- `GET /proxy/network/v2/api/site/{site}/vendor-ids` - Vendor IDs

### System Logs API (v2)
- `GET /proxy/network/v2/api/site/{site}/system-log/setting` - System log settings
- `POST /proxy/network/v2/api/site/{site}/system-log/all` - All system logs (already implemented)
- `POST /proxy/network/v2/api/site/{site}/system-log/critical` - Critical system logs
- `POST /proxy/network/v2/api/site/{site}/system-log/count` - Log count
- `POST /proxy/network/v2/api/site/{site}/system-log/filter-data` - Log filter data
- `GET /proxy/network/v2/api/site/{site}/system-log/device/{mac}` - Device-specific system logs

### Topology & Network Topology
- `GET /proxy/network/v2/api/site/{site}/topology` - Network topology
- `GET /proxy/network/v2/api/site/{site}/mclag-groups` - MCLAG groups
- `GET /proxy/network/v2/api/site/{site}/stacking` - Stacking info

### Advanced Features
- `GET /proxy/network/v2/api/site/{site}/features/{feature}/exists` - Feature existence checks
- `GET /proxy/network/v2/api/site/{site}/site-feature-migration` - Feature migration info
- `GET /proxy/network/v2/api/site/{site}/shadowmode/info` - Shadow mode info
- `GET /proxy/network/v2/api/site/{site}/settings/usg/defaults` - USG defaults
- `GET /proxy/network/v2/api/site/{site}/settings/ntp/defaults` - NTP defaults
- `GET /proxy/network/v2/api/site/{site}/settings/netflow/defaults` - NetFlow defaults
- `GET /proxy/network/v2/api/site/{site}/settings/global_switch/defaults` - Global switch defaults

### UPS (Uninterruptible Power Supply) API
- `GET /proxy/network/v2/api/ups/consoles` - UPS consoles information

### Routing Protocols
- `GET /proxy/network/v2/api/site/{site}/ospf/router` - OSPF router config
- `GET /proxy/network/v2/api/site/{site}/bgp/config/all` - BGP config

### Port Management
- `GET /proxy/network/v2/api/site/{site}/ports/port-anomalies` - Port anomalies
- `POST /proxy/network/v2/api/site/{site}/ports/mac-tables` - MAC tables

### RADIUS & Authentication
- `GET /proxy/network/v2/api/site/{site}/radius/users` - RADIUS users
- `GET /proxy/network/v2/api/site/{site}/hotspot/info` - Hotspot info
- `GET /proxy/network/api/stat/s/{site}/hotspotpackages` - Hotspot packages
- `GET /proxy/network/v2/api/site/{site}/ulp/users_groups` - ULP users/groups

### Users API (v2)
- `GET /proxy/users/api/v2/user/self` - Current user
- `GET /proxy/users/api/v2/users` - All users
- `GET /proxy/users/api/v2/users/{user}/uos` - User UOS info
- `GET /proxy/users/api/v2/user/{id}/keys` - User API keys
- `GET /proxy/users/api/v2/info` - Users API info

### Other Services
- `GET /proxy/innerspace/api/project` - InnerSpace project (3D visualization)
- `GET /proxy/innerspace/{id}/{id}` - InnerSpace resources
- `GET /proxy/innerspace/{id}/scan.gltf` - 3D scan data (GLTF)
- `GET /proxy/innerspace/{id}/scan.bin` - 3D scan data (binary)

---

## Summary

**Total Unique Endpoints**: 217 (47 cloud + 170 local)

**Endpoint Categories**:
- **Cloud APIs**: 47 endpoints
  - SSO/Authentication: 11 endpoints
  - Enterprise/Organization: 9 endpoints  
  - Cloud Access: 4 endpoints
  - Billing: 12 endpoints
  - IPS: 2 endpoints
  - Other Services: 9 endpoints

- **Local Controller APIs**: 174 endpoints (170 from first session + 4 new)
  - Authentication & System: 12 endpoints
  - Alarms: 10 endpoints
  - Network API (v1): 16 endpoints
  - Network Configuration: 14 endpoints
  - Network API (v2): 10 endpoints
  - Clients: 3 endpoints
  - Devices: 4 endpoints
  - WAN/LAN: 13 endpoints
  - WiFi/WLAN: 11 endpoints
  - Traffic & DPI: 6 endpoints
  - Firewall & Security: 15 endpoints
  - VPN: 4 endpoints
  - Network Management: 8 endpoints
  - System Logs: 5 endpoints
  - Topology: 3 endpoints
  - Advanced Features: 7 endpoints
  - Routing: 2 endpoints
  - Port Management: 2 endpoints
  - RADIUS & Authentication: 4 endpoints
  - Users: 5 endpoints
  - Other Services: 3 endpoints

## Notes
- **Discovery Method**: These endpoints were discovered by monitoring browser network requests during a session on a local UniFi controller
- **Testing**: Output has been tested on most endpoints to verify functionality
- **Data Availability**: Response data has been collected but is not published as it contains sensitive information and represents a large volume of data
- Many endpoints use organization IDs and host IDs in the path
- Some endpoints support query parameters like `pageSize`, `limit`, `expand[]`, `include`, etc.
- The local controller appears to be accessed via `briangates.ui.com` (cloud access URL) OR directly via IP
- **Important**: The cloud UI uses WebSocket connections for real-time data. The traditional REST API endpoints (`/api/s/{site}/stat/device`, etc.) are accessible via direct HTTP requests when accessing the controller directly
- Most endpoints discovered are cloud/enterprise APIs, but we now have comprehensive local controller REST API endpoints
- Many v2 API endpoints provide enhanced functionality over v1 endpoints
- Some endpoints require POST with filter data for pagination/filtering (e.g., system-log, traffic-flows)
