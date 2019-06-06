package unifi

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u USG) Points() ([]*influx.Point, error) {
	tags := map[string]string{
		"id":                     u.ID,
		"mac":                    u.Mac,
		"device_type":            u.Stat.O,
		"device_oid":             u.Stat.Oid,
		"site_id":                u.SiteID,
		"site_name":              u.SiteName,
		"adopted":                u.Adopted.Txt,
		"name":                   u.Name,
		"adopt_ip":               u.AdoptIP,
		"adopt_url":              u.AdoptURL,
		"cfgversion":             u.Cfgversion,
		"config_network_ip":      u.ConfigNetwork.IP,
		"config_network_type":    u.ConfigNetwork.Type,
		"connect_request_ip":     u.ConnectRequestIP,
		"connect_request_port":   u.ConnectRequestPort,
		"default":                u.Default.Txt,
		"device_id":              u.DeviceID,
		"discovered_via":         u.DiscoveredVia,
		"guest_token":            u.GuestToken,
		"inform_ip":              u.InformIP,
		"known_cfgversion":       u.KnownCfgversion,
		"led_override":           u.LedOverride,
		"locating":               u.Locating.Txt,
		"model":                  u.Model,
		"outdoor_mode_override":  u.OutdoorModeOverride,
		"serial":                 u.Serial,
		"type":                   u.Type,
		"version_incompatible":   u.VersionIncompatible.Txt,
		"usg_caps":               strconv.FormatFloat(u.UsgCaps, 'f', 6, 64),
		"speedtest-status-saved": u.SpeedtestStatusSaved.Txt,
	}
	fields := map[string]interface{}{
		"ip":                             u.IP,
		"bytes":                          u.Bytes,
		"last_seen":                      u.LastSeen,
		"license_state":                  u.LicenseState,
		"fw_caps":                        u.FwCaps,
		"guest-num_sta":                  u.GuestNumSta,
		"rx_bytes":                       u.RxBytes.Val,
		"tx_bytes":                       u.TxBytes.Val,
		"uptime":                         u.Uptime,
		"considered_lost_at":             u.ConsideredLostAt,
		"next_heartbeat_at":              u.NextHeartbeatAt,
		"roll_upgrade":                   u.Rollupgrade.Val,
		"state":                          u.State,
		"upgradable":                     u.Upgradable.Val,
		"user-num_sta":                   u.UserNumSta,
		"version":                        u.Version,
		"num_desktop":                    u.NumDesktop,
		"num_handheld":                   u.NumHandheld,
		"num_mobile":                     u.NumMobile,
		"speedtest-status_latency":       u.SpeedtestStatus.Latency,
		"speedtest-status_rundate":       u.SpeedtestStatus.Rundate,
		"speedtest-status_runtime":       u.SpeedtestStatus.Runtime,
		"speedtest-status_download":      u.SpeedtestStatus.StatusDownload,
		"speedtest-status_ping":          u.SpeedtestStatus.StatusPing,
		"speedtest-status_summary":       u.SpeedtestStatus.StatusSummary,
		"speedtest-status_upload":        u.SpeedtestStatus.StatusUpload,
		"speedtest-status_xput_download": u.SpeedtestStatus.XputDownload,
		"speedtest-status_xput_upload":   u.SpeedtestStatus.XputUpload,
		// have two WANs? mmmm, go ahead and add it. ;)
		"config_network_wan_type": u.ConfigNetworkWan.Type,
		"wan1_bytes-r":            u.Wan1.BytesR,
		"wan1_enable":             u.Wan1.Enable.Val,
		"wan1_full_duplex":        u.Wan1.FullDuplex.Val,
		"wan1_purpose":            "uplink", // because it should have a purpose.
		"wan1_gateway":            u.Wan1.Gateway,
		"wan1_ifname":             u.Wan1.Ifname,
		"wan1_ip":                 u.Wan1.IP,
		"wan1_mac":                u.Wan1.Mac,
		"wan1_max_speed":          u.Wan1.MaxSpeed,
		"wan1_name":               u.Wan1.Name,
		"wan1_netmask":            u.Wan1.Netmask,
		"wan1_rx_bytes":           u.Wan1.RxBytes.Val,
		"wan1_rx_bytes-r":         u.Wan1.RxBytesR,
		"wan1_rx_dropped":         u.Wan1.RxDropped,
		"wan1_rx_errors":          u.Wan1.RxErrors,
		"wan1_rx_multicast":       u.Wan1.RxMulticast,
		"wan1_rx_packets":         u.Wan1.RxPackets,
		"wan1_type":               u.Wan1.Type,
		"wan1_speed":              u.Wan1.Speed,
		"wan1_up":                 u.Wan1.Up.Val,
		"wan1_tx_bytes":           u.Wan1.TxBytes.Val,
		"wan1_tx_bytes-r":         u.Wan1.TxBytesR,
		"wan1_tx_dropped":         u.Wan1.TxDropped,
		"wan1_tx_errors":          u.Wan1.TxErrors,
		"wan1_tx_packets":         u.Wan1.TxPackets,
		"loadavg_1":               u.SysStats.Loadavg1,
		"loadavg_5":               u.SysStats.Loadavg5,
		"loadavg_15":              u.SysStats.Loadavg15,
		"mem_used":                u.SysStats.MemUsed,
		"mem_buffer":              u.SysStats.MemBuffer,
		"mem_total":               u.SysStats.MemTotal,
		"cpu":                     u.SystemStats.CPU,
		"mem":                     u.SystemStats.Mem,
		"system_uptime":           u.SystemStats.Uptime,
		"stat_duration":           u.Stat.Duration,
		"stat_datetime":           u.Stat.Datetime,
		"gw":                      u.Stat.Gw,
		"false":                   "false", // to fill holes in graphs.
		"lan-rx_bytes":            u.Stat.LanRxBytes,
		"lan-rx_packets":          u.Stat.LanRxPackets,
		"lan-tx_bytes":            u.Stat.LanTxBytes,
		"lan-tx_packets":          u.Stat.LanTxPackets,
		"wan-rx_bytes":            u.Stat.WanRxBytes,
		"wan-rx_dropped":          u.Stat.WanRxDropped,
		"wan-rx_packets":          u.Stat.WanRxPackets,
		"wan-tx_bytes":            u.Stat.WanTxBytes,
		"wan-tx_packets":          u.Stat.WanTxPackets,
		"uplink_name":             u.Uplink.Name,
		"uplink_latency":          u.Uplink.Latency,
		"uplink_speed":            u.Uplink.Speed,
		"uplink_num_ports":        u.Uplink.NumPort,
		"uplink_max_speed":        u.Uplink.MaxSpeed,
	}
	pt, err := influx.NewPoint("usg", tags, fields, time.Now())
	if err != nil {
		return nil, err
	}
	points := []*influx.Point{pt}
	for _, p := range u.NetworkTable {
		tags := map[string]string{
			"device_name":               u.Name,
			"device_id":                 u.ID,
			"device_mac":                u.Mac,
			"name":                      p.Name,
			"dhcpd_dns_enabled":         p.DhcpdDNSEnabled.Txt,
			"dhcpd_enabled":             p.DhcpdEnabled.Txt,
			"dhcpd_ntp_enabled":         p.DhcpdNtpEnabled.Txt,
			"dhcpd_time_offset_enabled": p.DhcpdTimeOffsetEnabled.Txt,
			"dhcp_relay_enabledy":       p.DhcpRelayEnabled.Txt,
			"dhcpd_gateway_enabled":     p.DhcpdGatewayEnabled.Txt,
			"dhcpd_wins_enabled":        p.DhcpdWinsEnabled.Txt,
			"dhcpguard_enabled":         p.DhcpguardEnabled.Txt,
			"enabled":                   p.Enabled.Txt,
			"vlan_enabled":              p.VlanEnabled.Txt,
			"attr_no_delete":            p.AttrNoDelete.Txt,
			"upnp_lan_enabled":          p.UpnpLanEnabled.Txt,
			"igmp_snooping":             p.IgmpSnooping.Txt,
			"is_guest":                  p.IsGuest.Txt,
			"is_nat":                    p.IsNat.Txt,
			"networkgroup":              p.Networkgroup,
			"site_id":                   p.SiteID,
		}
		fields := map[string]interface{}{
			"dhcpd_ip_1":             p.DhcpdIP1,
			"domain_name":            p.DomainName,
			"dhcpd_start":            p.DhcpdStart,
			"dhcpd_stop":             p.DhcpdStop,
			"ip":                     p.IP,
			"ip_subnet":              p.IPSubnet,
			"mac":                    p.Mac,
			"name":                   p.Name,
			"num_sta":                p.NumSta,
			"purpose":                p.Purpose,
			"rx_bytes":               p.RxBytes.Val,
			"rx_packets":             p.RxPackets,
			"tx_bytes":               p.TxBytes.Val,
			"tx_packets":             p.TxPackets,
			"up":                     p.Up.Txt,
			"vlan":                   p.Vlan,
			"dhcpd_ntp_1":            p.DhcpdNtp1,
			"dhcpd_unifi_controller": p.DhcpdUnifiController,
			"ipv6_interface_type":    p.Ipv6InterfaceType,
			"attr_hidden_id":         p.AttrHiddenID,
		}
		pt, err = influx.NewPoint("usg_networks", tags, fields, time.Now())
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
