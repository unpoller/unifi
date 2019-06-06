package unifi

import "encoding/json"

// USG represents all the data from the Ubiquiti Controller for a Unifi Security Gateway.
type USG struct {
	ID            string   `json:"_id"`
	UUptime       float64  `json:"_uptime"`
	AdoptIP       string   `json:"adopt_ip"`
	AdoptURL      string   `json:"adopt_url"`
	Adopted       FlexBool `json:"adopted"`
	Bytes         float64  `json:"bytes"`
	Cfgversion    string   `json:"cfgversion"`
	ConfigNetwork struct {
		IP   string `json:"ip"`
		Type string `json:"type"`
	} `json:"config_network"`
	ConfigNetworkWan struct {
		Type string `json:"type"`
	} `json:"config_network_wan"`
	ConnectRequestIP   string   `json:"connect_request_ip"`
	ConnectRequestPort string   `json:"connect_request_port"`
	ConsideredLostAt   float64  `json:"considered_lost_at"`
	Default            FlexBool `json:"default"`
	DeviceID           string   `json:"device_id"`
	DiscoveredVia      string   `json:"discovered_via"`
	EthernetTable      []struct {
		Mac     string  `json:"mac"`
		Name    string  `json:"name"`
		NumPort float64 `json:"num_port"`
	} `json:"ethernet_table"`
	FwCaps          float64  `json:"fw_caps"`
	GuestNumSta     float64  `json:"guest-num_sta"`
	GuestToken      string   `json:"guest_token"`
	InformIP        string   `json:"inform_ip"`
	InformURL       string   `json:"inform_url"`
	IP              string   `json:"ip"`
	KnownCfgversion string   `json:"known_cfgversion"`
	LastSeen        float64  `json:"last_seen"`
	LedOverride     string   `json:"led_override"`
	LicenseState    string   `json:"license_state"`
	Locating        FlexBool `json:"locating"`
	Mac             string   `json:"mac"`
	Model           string   `json:"model"`
	Name            string   `json:"name"`
	NetworkTable    []struct {
		ID                     string      `json:"_id"`
		DhcpdDNSEnabled        FlexBool    `json:"dhcpd_dns_enabled"`
		DhcpdEnabled           FlexBool    `json:"dhcpd_enabled"`
		DhcpdIP1               string      `json:"dhcpd_ip_1,omitempty"`
		DhcpdLeasetime         json.Number `json:"dhcpd_leasetime,Number"`
		DhcpdStart             string      `json:"dhcpd_start"`
		DhcpdStop              string      `json:"dhcpd_stop"`
		DhcpdWinsEnabled       FlexBool    `json:"dhcpd_wins_enabled,omitempty"`
		DhcpguardEnabled       FlexBool    `json:"dhcpguard_enabled,omitempty"`
		DomainName             string      `json:"domain_name"`
		Enabled                FlexBool    `json:"enabled"`
		IgmpSnooping           FlexBool    `json:"igmp_snooping,omitempty"`
		IP                     string      `json:"ip"`
		IPSubnet               string      `json:"ip_subnet"`
		IsGuest                FlexBool    `json:"is_guest"`
		IsNat                  FlexBool    `json:"is_nat"`
		Mac                    string      `json:"mac"`
		Name                   string      `json:"name"`
		Networkgroup           string      `json:"networkgroup"`
		NumSta                 float64     `json:"num_sta"`
		Purpose                string      `json:"purpose"`
		RxBytes                FlexInt     `json:"rx_bytes"`
		RxPackets              float64     `json:"rx_packets"`
		SiteID                 string      `json:"site_id"`
		TxBytes                FlexInt     `json:"tx_bytes"`
		TxPackets              float64     `json:"tx_packets"`
		Up                     FlexBool    `json:"up"`
		Vlan                   string      `json:"vlan,omitempty"`
		VlanEnabled            FlexBool    `json:"vlan_enabled"`
		DhcpRelayEnabled       FlexBool    `json:"dhcp_relay_enabled,omitempty"`
		DhcpdGatewayEnabled    FlexBool    `json:"dhcpd_gateway_enabled,omitempty"`
		DhcpdNtp1              string      `json:"dhcpd_ntp_1,omitempty"`
		DhcpdNtpEnabled        FlexBool    `json:"dhcpd_ntp_enabled,omitempty"`
		DhcpdTimeOffsetEnabled FlexBool    `json:"dhcpd_time_offset_enabled,omitempty"`
		DhcpdUnifiController   string      `json:"dhcpd_unifi_controller,omitempty"`
		Ipv6InterfaceType      string      `json:"ipv6_interface_type,omitempty"`
		AttrHiddenID           string      `json:"attr_hidden_id,omitempty"`
		AttrNoDelete           FlexBool    `json:"attr_no_delete,omitempty"`
		UpnpLanEnabled         FlexBool    `json:"upnp_lan_enabled,omitempty"`
	} `json:"network_table"`
	NextHeartbeatAt     float64 `json:"next_heartbeat_at"`
	NumDesktop          float64 `json:"num_desktop"`
	NumHandheld         float64 `json:"num_handheld"`
	NumMobile           float64 `json:"num_mobile"`
	NumSta              float64 `json:"num_sta"`
	OutdoorModeOverride string  `json:"outdoor_mode_override"`
	PortTable           []struct {
		DNS         []string `json:"dns,omitempty"`
		Enable      FlexBool `json:"enable"`
		FullDuplex  FlexBool `json:"full_duplex"`
		Gateway     string   `json:"gateway,omitempty"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxDropped   float64  `json:"rx_dropped"`
		RxErrors    float64  `json:"rx_errors"`
		RxMulticast float64  `json:"rx_multicast"`
		RxPackets   float64  `json:"rx_packets"`
		Speed       float64  `json:"speed"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxDropped   float64  `json:"tx_dropped"`
		TxErrors    float64  `json:"tx_errors"`
		TxPackets   float64  `json:"tx_packets"`
		Up          FlexBool `json:"up"`
	} `json:"port_table"`
	Rollupgrade     FlexBool `json:"rollupgrade"`
	RxBytes         FlexInt  `json:"rx_bytes"`
	Serial          string   `json:"serial"`
	SiteID          string   `json:"site_id"`
	SiteName        string   `json:"-"`
	SpeedtestStatus struct {
		Latency        float64 `json:"latency"`
		Rundate        float64 `json:"rundate"`
		Runtime        float64 `json:"runtime"`
		StatusDownload float64 `json:"status_download"`
		StatusPing     float64 `json:"status_ping"`
		StatusSummary  float64 `json:"status_summary"`
		StatusUpload   float64 `json:"status_upload"`
		XputDownload   float64 `json:"xput_download"`
		XputUpload     float64 `json:"xput_upload"`
	} `json:"speedtest-status"`
	SpeedtestStatusSaved FlexBool `json:"speedtest-status-saved"`
	Stat                 struct {
		Datetime     string  `json:"datetime"`
		Duration     float64 `json:"duration"`
		Gw           string  `json:"gw"`
		LanRxBytes   float64 `json:"lan-rx_bytes"`
		LanRxPackets float64 `json:"lan-rx_packets"`
		LanTxBytes   float64 `json:"lan-tx_bytes"`
		LanTxPackets float64 `json:"lan-tx_packets"`
		O            string  `json:"o"`
		Oid          string  `json:"oid"`
		SiteID       string  `json:"site_id"`
		Time         float64 `json:"time"`
		WanRxBytes   float64 `json:"wan-rx_bytes"`
		WanRxDropped float64 `json:"wan-rx_dropped"`
		WanRxPackets float64 `json:"wan-rx_packets"`
		WanTxBytes   float64 `json:"wan-tx_bytes"`
		WanTxPackets float64 `json:"wan-tx_packets"`
	} `json:"stat"`
	State    float64 `json:"state"`
	SysStats struct {
		Loadavg1  float64 `json:"loadavg_1,string"`
		Loadavg15 float64 `json:"loadavg_15,string"`
		Loadavg5  float64 `json:"loadavg_5,string"`
		MemBuffer float64 `json:"mem_buffer"`
		MemTotal  float64 `json:"mem_total"`
		MemUsed   float64 `json:"mem_used"`
	} `json:"sys_stats"`
	SystemStats struct {
		CPU    float64 `json:"cpu,string"`
		Mem    float64 `json:"mem,string"`
		Uptime float64 `json:"uptime,string"`
	} `json:"system-stats"`
	TxBytes    FlexInt  `json:"tx_bytes"`
	Type       string   `json:"type"`
	Upgradable FlexBool `json:"upgradable"`
	Uplink     struct {
		BytesR           float64  `json:"bytes-r"`
		Drops            float64  `json:"drops"`
		Enable           FlexBool `json:"enable"`
		FullDuplex       FlexBool `json:"full_duplex"`
		Gateways         []string `json:"gateways"`
		IP               string   `json:"ip"`
		Latency          float64  `json:"latency"`
		Mac              string   `json:"mac"`
		MaxSpeed         float64  `json:"max_speed"`
		Name             string   `json:"name"`
		Nameservers      []string `json:"nameservers"`
		Netmask          string   `json:"netmask"`
		NumPort          float64  `json:"num_port"`
		RxBytes          FlexInt  `json:"rx_bytes"`
		RxBytesR         float64  `json:"rx_bytes-r"`
		RxDropped        float64  `json:"rx_dropped"`
		RxErrors         float64  `json:"rx_errors"`
		RxMulticast      float64  `json:"rx_multicast"`
		RxPackets        float64  `json:"rx_packets"`
		Speed            float64  `json:"speed"`
		SpeedtestLastrun float64  `json:"speedtest_lastrun"`
		SpeedtestPing    float64  `json:"speedtest_ping"`
		SpeedtestStatus  string   `json:"speedtest_status"`
		TxBytes          FlexInt  `json:"tx_bytes"`
		TxBytesR         float64  `json:"tx_bytes-r"`
		TxDropped        float64  `json:"tx_dropped"`
		TxErrors         float64  `json:"tx_errors"`
		TxPackets        float64  `json:"tx_packets"`
		Type             string   `json:"type"`
		Up               FlexBool `json:"up"`
		Uptime           float64  `json:"uptime"`
		XputDown         float64  `json:"xput_down"`
		XputUp           float64  `json:"xput_up"`
	} `json:"uplink"`
	Uptime              float64  `json:"uptime"`
	UserNumSta          float64  `json:"user-num_sta"`
	UsgCaps             float64  `json:"usg_caps"`
	Version             string   `json:"version"`
	VersionIncompatible FlexBool `json:"version_incompatible"`
	Wan1                struct {
		BytesR      float64  `json:"bytes-r"`
		DNS         []string `json:"dns"`
		Enable      FlexBool `json:"enable"`
		FullDuplex  FlexBool `json:"full_duplex"`
		Gateway     string   `json:"gateway"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		MaxSpeed    float64  `json:"max_speed"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxBytesR    float64  `json:"rx_bytes-r"`
		RxDropped   float64  `json:"rx_dropped"`
		RxErrors    float64  `json:"rx_errors"`
		RxMulticast float64  `json:"rx_multicast"`
		RxPackets   float64  `json:"rx_packets"`
		Speed       float64  `json:"speed"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxBytesR    float64  `json:"tx_bytes-r"`
		TxDropped   float64  `json:"tx_dropped"`
		TxErrors    float64  `json:"tx_errors"`
		TxPackets   float64  `json:"tx_packets"`
		Type        string   `json:"type"`
		Up          FlexBool `json:"up"`
	} `json:"wan1"`
	Wan2 struct {
		BytesR      float64  `json:"bytes-r"`
		DNS         []string `json:"dns"`
		Enable      FlexBool `json:"enable"`
		FullDuplex  FlexBool `json:"full_duplex"`
		Gateway     string   `json:"gateway"`
		Ifname      string   `json:"ifname"`
		IP          string   `json:"ip"`
		Mac         string   `json:"mac"`
		MaxSpeed    float64  `json:"max_speed"`
		Name        string   `json:"name"`
		Netmask     string   `json:"netmask"`
		RxBytes     FlexInt  `json:"rx_bytes"`
		RxBytesR    float64  `json:"rx_bytes-r"`
		RxDropped   float64  `json:"rx_dropped"`
		RxErrors    float64  `json:"rx_errors"`
		RxMulticast float64  `json:"rx_multicast"`
		RxPackets   float64  `json:"rx_packets"`
		Speed       float64  `json:"speed"`
		TxBytes     FlexInt  `json:"tx_bytes"`
		TxBytesR    float64  `json:"tx_bytes-r"`
		TxDropped   float64  `json:"tx_dropped"`
		TxErrors    float64  `json:"tx_errors"`
		TxPackets   float64  `json:"tx_packets"`
		Type        string   `json:"type"`
		Up          FlexBool `json:"up"`
	} `json:"wan2"`
}
