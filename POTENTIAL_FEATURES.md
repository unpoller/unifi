# Potential Features to Add to unpoller/unifi

Based on analysis of API response examples, here are features that could be added:

## 1. System Information (`/api/s/{site}/stat/sysinfo`)

**Status**: Not implemented

**Data Available**:
- Controller version, build, previous version
- Timezone, hostname, IP addresses
- Ports (inform, https, portal_http)
- Uptime
- Data retention settings
- Update availability
- Device type (UDRULT, etc.)
- Debug settings
- SSO app ID
- Anonymous controller ID
- WebRTC support
- Unsupported device count/list

**Use Case**: Monitoring controller health, version tracking, uptime monitoring

**Implementation**: Add `GetSysInfo(site *Site) (*SysInfo, error)` method

---

## 2. WAN Status (`/api/s/{site}/stat/status`)

**Status**: Not implemented

**Data Available**:
- WAN interface names
- WAN states (ACTIVE, BACKUP)
- WAN network groups (WAN, WAN2, etc.)

**Use Case**: Monitor WAN failover status, active/backup WAN interfaces

**Implementation**: Add `GetWANStatus(site *Site) (*WANStatus, error)` method

---

## 3. Firewall Policies (`/api/s/{site}/rest/firewall-policies`)

**Status**: Not implemented

**Data Available**:
- Firewall rules with zone-based configuration
- Actions (ALLOW, BLOCK)
- Source/destination zones
- Protocols (all, udp, icmpv6, etc.)
- Port matching
- IP version (BOTH, IPV4, IPV6)
- Connection states
- ICMP types
- Logging settings
- Schedule (ALWAYS mode shown)
- Predefined vs custom rules
- Index/priority

**Use Case**: Firewall rule auditing, security monitoring, rule analysis

**Implementation**: Add `GetFirewallPolicies(site *Site) ([]*FirewallPolicy, error)` method

**Note**: This is a significant feature with complex zone-based firewall data

---

## 4. UPS Devices List (`/api/s/{site}/stat/ups-devices`)

**Status**: Partially implemented (PDU/UPS devices come from `/stat/device`, but this endpoint provides a different format)

**Data Available**:
- Device selector format with image, label, metadata
- Device MAC addresses
- Site IDs
- Device identification for UI selection

**Use Case**: Quick UPS device lookup, device selection lists

**Implementation**: Add `GetUPSDeviceList(site *Site) ([]*UPSDeviceSelector, error)` method

**Note**: Different from full device data - this is a lightweight selector format

---

## 5. Enhanced Network Configuration

**Status**: Partially implemented (basic Network struct exists, but missing many fields)

**Missing Fields from `/api/s/{site}/rest/networkconf`**:
- WAN-specific: `wan_type`, `wan_dhcp_options`, `wan_load_balance_type`, `wan_failover_priority`, `wan_provider_capabilities`, `wan_ip_aliases`, `wan_smartq_enabled`, `wan_dns_preference`, `wan_vlan_enabled`, `wan_ipv6_dns_preference`, `wan_ipv6_dns1/2`, `ipv6_wan_delegation_type`, `wan_dhcpv6_pd_size_auto`
- IPv6: `ipv6_enabled`, `ipv6_setting_preference`, `ipv6_interface_type`, `ipv6_pd_start/stop`, `ipv6_pd_auto_prefixid_enabled`, `ipv6_ra_enabled`, `ipv6_ra_preferred_lifetime`, `ipv6_ra_priority`, `ipv6_client_address_assignment`, `dhcpdv6_*` fields
- VPN: `vpn_type`, `sdwan_remote_site_id`, `remote_vpn_subnets`, `ifname` (for site-vpn)
- Advanced: `firewall_zone_id`, `routing_table_id`, `external_id`, `igmp_proxy_*`, `mac_override_enabled`, `setting_preference`, `report_wan_event`, `attr_no_delete`, `attr_hidden_id`, `attr_no_edit`
- DHCP: `dhcpdv6_*` fields, `dhcpd_wpad_url`, `dhcpd_tftp_server`, `dhcpd_boot_enabled`, `dhcpd_ntp_enabled`, `dhcpd_wins_enabled`
- Other: `lte_lan_enabled`, `upnp_lan_enabled`, `mdns_enabled`, `auto_scale_enabled`, `dhcpguard_enabled`, `dhcpd_conflict_checking`, `dhcpd_time_offset_enabled`, `dhcpd_gateway_enabled`, `dhcp_relay_enabled`, `nat_outbound_ip_addresses`, `single_network_lan`, `gateway_type`, `domain_name`

**Use Case**: Full network configuration monitoring, WAN failover configuration, IPv6 support, VPN monitoring

**Implementation**: Extend `Network` struct with additional fields, or create `NetworkExtended` struct

---

## 6. Port Forwarding (`/api/s/{site}/rest/portforward`)

**Status**: Not implemented (endpoint exists but returns empty array in example)

**Data Available**: (when configured)
- Port forward rules
- Source/destination ports
- Protocols
- Enabled status

**Use Case**: Port forward rule auditing

**Implementation**: Add `GetPortForwards(site *Site) ([]*PortForward, error)` method

---

## 7. SSL Certificate Info (`/api/s/{site}/stat/active`)

**Status**: Not implemented

**Data Available**:
- Certificate details (root, intermediate)
- Certificate status (active, valid)
- Certificate type (generated)
- Valid from/to dates
- Issuer/subject information
- Fingerprint, serial number

**Use Case**: SSL certificate monitoring, expiration tracking

**Implementation**: Add `GetSSLCertificate(site *Site) (*SSLCertificate, error)` method

---

## Priority Recommendations

### High Priority
1. **SysInfo** - Useful for monitoring and health checks
2. **WAN Status** - Important for failover monitoring
3. **Enhanced Network Config** - Many missing fields that users might need

### Medium Priority
4. **Firewall Policies** - Complex but valuable for security auditing
5. **SSL Certificate** - Useful for certificate management

### Low Priority
6. **UPS Device List** - Different format, may not be needed if device endpoint covers it
7. **Port Forwarding** - Less commonly used, endpoint was empty in example

---

## Implementation Notes

- All endpoints follow the standard `/api/s/{site}/...` pattern
- Most return `{"meta": {"rc": "ok"}, "data": [...]}` format
- Some endpoints (like `status`) return different formats
- Need to handle UDM Pro `/proxy/network` prefix
- Consider backward compatibility when extending existing structs
