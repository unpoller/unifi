package unifi

// UDM represents all the data from the Ubiquiti Controller for a Unifi Dream Machine.
// The UDM shares several structs/type-data with USW and USG.
type UDM struct {
	site                               *Site
	AFTEnabled                         FlexBool             `json:"atf_enabled"`
	AdoptIP                            string               `fake:"{ipv4address}"                          json:"adopt_ip"`
	AdoptManual                        FlexBool             `json:"adopt_manual"`
	AdoptState                         FlexInt              `json:"adopt_state"`
	AdoptStatus                        FlexInt              `json:"adopt_status"`
	AdoptTries                         FlexInt              `json:"adopt_tries"`
	AdoptURL                           string               `fake:"{url}"                                  json:"adopt_url"`
	AdoptableWhenUpgraded              FlexBool             `fake:"{constFlexBool:true}"                   json:"adoptable_when_upgraded"`
	Adopted                            FlexBool             `fake:"{constFlexBool:true}"                   json:"adopted"`
	AdoptionCompleted                  FlexBool             `fake:"{constFlexBool:true}"                   json:"adoption_completed"`
	Architecture                       string               `json:"architecture"`
	BandsteeringMode                   string               `json:"bandsteering_mode"`
	BoardRev                           FlexInt              `json:"board_rev"`
	Bytes                              FlexInt              `json:"bytes"`
	BytesD                             FlexInt              `json:"bytes-d"`
	BytesR                             FlexInt              `json:"bytes-r"`
	Cfgversion                         string               `fake:"{appversion}"                           json:"cfgversion"`
	ConfigNetwork                      *ConfigNetwork       `json:"config_network"`
	ConnectRequestIP                   string               `fake:"{ipv4address}"                          json:"connect_request_ip"`
	ConnectRequestPort                 string               `json:"connect_request_port"`
	ConnectedAt                        FlexInt              `json:"connected_at"`
	ConnectionNetworkName              string               `json:"connection_network_name"`
	Default                            FlexBool             `json:"default"`
	DeviceDomain                       string               `json:"device_domain"`
	DeviceID                           string               `fake:"{uuid}"                                 json:"device_id"`
	DiscoveredVia                      string               `json:"discovered_via"`
	DisplayableVersion                 string               `fake:"{appversion}"                           json:"displayable_version"`
	Dot1XPortctrlEnabled               FlexBool             `json:"dot1x_portctrl_enabled"`
	DownlinkTable                      []*DownlinkTable     `fakesize:"5"                                  json:"downlink_table"`
	EthernetOverrides                  []*EthernetOverrides `fakesize:"5"                                  json:"ethernet_overrides"`
	EthernetTable                      []*EthernetTable     `fakesize:"5"                                  json:"ethernet_table"`
	FlowctrlEnabled                    FlexBool             `json:"flowctrl_enabled"`
	FwCaps                             FlexInt              `json:"fw_caps"`
	GeoInfo                            map[string]GeoInfo   `fakesize:"5"                                  json:"geo_info"`
	GuestKicks                         FlexInt              `json:"guest_kicks"`
	GuestLanNumSta                     FlexInt              `json:"guest-lan-num_sta"` // USW
	GuestNumSta                        FlexInt              `json:"guest-num_sta"`     // USG
	GuestToken                         string               `json:"guest_token"`
	GuestWlanNumSta                    FlexInt              `json:"guest-wlan-num_sta"` // UAP
	HasEth1                            FlexBool             `json:"has_eth1"`
	HasFan                             FlexBool             `json:"has_fan"`
	HasSpeaker                         FlexBool             `json:"has_speaker"`
	HasTemperature                     FlexBool             `json:"has_temperature"`
	HwCaps                             FlexInt              `json:"hw_caps"`
	ID                                 string               `fake:"{uuid}"                                 json:"_id"`
	IP                                 string               `fake:"{ipv4address}"                          json:"ip"`
	InformIP                           string               `fake:"{ipv4address}"                          json:"inform_ip"`
	InformURL                          string               `fake:"{url}"                                  json:"inform_url"`
	Internet                           FlexBool             `json:"internet"`
	IsAccessPoint                      FlexBool             `json:"is_access_point"`
	JumboframeEnabled                  FlexBool             `json:"jumboframe_enabled"`
	KernelVersion                      string               `fake:"{appversion}"                           json:"kernel_version"`
	KnownCfgversion                    string               `fake:"{appversion}"                           json:"known_cfgversion"`
	LanIP                              string               `fake:"{ipv4address}"                          json:"lan_ip"`
	LanNumSta                          FlexInt              `json:"lan-num_sta"` // USW
	LastLteFailoverTransitionTimestamp FlexInt              `json:"last_lte_failover_transition_timestamp"`
	LastSeen                           FlexInt              `json:"last_seen"`
	LastWlanIP                         string               `fake:"{ipv4address}"                          json:"last_wan_ip"`
	LcmBrightness                      FlexInt              `json:"lcm_brightness"`
	LcmNightModeBegins                 string               `json:"lcm_night_mode_begins"`
	LcmNightModeEnabled                FlexBool             `json:"lcm_night_mode_enabled"`
	LcmNightModeEnds                   string               `json:"lcm_night_mode_ends"`
	LcmTrackerEnabled                  FlexBool             `json:"lcm_tracker_enabled"`
	LcmTrackerSeed                     string               `json:"lcm_tracker_seed"`
	LicenseState                       string               `json:"license_state"`
	Locating                           FlexBool             `fake:"{constFlexBool:false}"                  json:"locating"`
	Mac                                string               `fake:"{macaddress}"                           json:"mac"`
	Tags                               []string             `json:"tags"` // Device tags assigned to this device
	ManufacturerID                     FlexInt              `json:"manufacturer_id"`
	MinInformIntervalSeconds           FlexInt              `json:"min_inform_interval_seconds"`
	Model                              string               `json:"model"`
	ModelInEOL                         FlexBool             `json:"model_in_eol"`
	ModelInLTS                         FlexBool             `json:"model_in_lts"`
	ModelIncompatible                  FlexBool             `json:"model_incompatible"`
	Name                               string               `fake:"{animal}"                               json:"name"`
	NetworkTable                       NetworkTable         `json:"network_table"`
	NextInterval                       FlexInt              `json:"next_interval"`
	NumDesktop                         FlexInt              `json:"num_desktop"`  // USG
	NumHandheld                        FlexInt              `json:"num_handheld"` // USG
	NumMobile                          FlexInt              `json:"num_mobile"`   // USG
	NumSta                             FlexInt              `json:"num_sta"`      // USG
	Overheating                        FlexBool             `json:"overheating"`
	PortOverrides                      []struct {
		PortIdx    FlexInt `json:"port_idx"`
		PortconfID string  `json:"portconf_id"`
	} `fakesize:"5"                       json:"port_overrides"`
	PortTable              []Port           `fakesize:"5"                       json:"port_table"`
	PowerSourceCtrlEnabled FlexBool         `json:"power_source_ctrl_enabled"`
	ProvisionedAt          FlexInt          `json:"provisioned_at"`
	RadioTable             *RadioTable      `json:"radio_table,omitempty"`
	RadioTableStats        *RadioTableStats `json:"radio_table_stats,omitempty"`
	RequiredVersion        string           `fake:"{appversion}"                json:"required_version"`
	RollUpgrade            FlexBool         `json:"rollupgrade"`
	RulesetInterfaces      interface{}      `json:"ruleset_interfaces"`
	/* struct {
		Br0  string `json:"br0"`
		Eth0 string `json:"eth0"`
		Eth1 string `json:"eth1"`
		Eth2 string `json:"eth2"`
		Eth3 string `json:"eth3"`
		Eth4 string `json:"eth4"`
		Eth5 string `json:"eth5"`
		Eth6 string `json:"eth6"`
		Eth7 string `json:"eth7"`
		Eth8 string `json:"eth8"`
	} */
	RxBytes                   FlexInt         `json:"rx_bytes"`
	RxBytesD                  FlexInt         `json:"rx_bytes-d"`
	Serial                    string          `json:"serial"`
	SetupProvisionCompleted   FlexBool        `json:"setup_provision_completed"`
	SetupProvisionTracking    FlexBool        `json:"setup_provision_tracking"`
	SiteID                    string          `fake:"{uuid}"                    json:"site_id"`
	SiteName                  string          `json:"-"`
	SourceName                string          `json:"-"`
	SpeedtestStatus           SpeedtestStatus `json:"speedtest-status"`
	SpeedtestStatusSaved      FlexBool        `json:"speedtest-status-saved"`
	StartupConnectedMillis    FlexInt         `json:"start_connected_millis"`
	StartupDisconnectedMillis FlexInt         `json:"start_disconnected_millis"`
	StartupTimestamp          FlexInt         `json:"startup_timestamp"`
	Stat                      UDMStat         `json:"stat"`
	State                     FlexInt         `json:"state"`
	Storage                   []*Storage      `fakesize:"5"                     json:"storage"`
	StpPriority               FlexInt         `json:"stp_priority"`
	StpVersion                string          `fake:"{appversion}"              json:"stp_version"`
	SwitchCaps                struct {
		MaxMirrorSessions    FlexInt `json:"max_mirror_sessions"`
		MaxAggregateSessions FlexInt `json:"max_aggregate_sessions"`
	} `json:"switch_caps"`
	SysStats        SysStats      `json:"sys_stats"`
	SyslogKey       string        `json:"syslog_key"`
	SystemStats     SystemStats   `json:"system-stats"`
	TeleportVersion FlexInt       `json:"teleport_version"`
	Temperatures    []Temperature `fakesize:"5"            json:"temperatures,omitempty"`
	TwoPhaseAdopt   FlexBool      `json:"two_phase_adopt"`
	TxBytes         FlexInt       `json:"tx_bytes"`
	TxBytesD        FlexInt       `json:"tx_bytes-d"`
	Type            string        `fake:"{lexify:udm}"     json:"type"`
	UdapiCaps       FlexInt       `json:"udapi_caps"`
	UnifiCare       struct {
		ActivationDismissed FlexBool `json:"activation_dismissed"`
		ActivationEnd       FlexInt  `json:"activation_end"`
		ActivationUrl       string   `fake:"{url}"                json:"activation_url"` //nolint:revive
		CoverageEnd         FlexInt  `json:"coverage_end"`
		CoverageStart       FlexInt  `json:"coverage_start"`
		Registration        FlexInt  `json:"registration"`
		RmaUrl              string   `fake:"{url}"                json:"rma_url"` //nolint:revive
		State               string   `json:"state"`
		TrackingUrl         string   `fake:"{url}"                json:"tracking_url"` //nolint:revive
	} `json:"unifi_care"`
	Unsupported       FlexBool      `json:"unsupported"`
	UnsupportedReason FlexInt       `json:"unsupported_reason"`
	UpgradeState      FlexInt       `json:"upgrade_state"`
	Upgradeable       FlexBool      `json:"upgradable"`
	Uplink            Uplink        `json:"uplink"`
	Uptime            FlexInt       `json:"uptime"`
	UserLanNumSta     FlexInt       `json:"user-lan-num_sta"`  // USW
	UserNumSta        FlexInt       `json:"user-num_sta"`      // USG
	UserWlanNumSta    FlexInt       `json:"user-wlan-num_sta"` // UAP
	UsgCaps           FlexInt       `json:"usg_caps"`
	VapTable          *VapTable     `json:"vap_table"`
	Version           string        `fake:"{appversion}"       json:"version"`
	VwireTable        []interface{} `json:"vwire_table"`
	Wan1              Wan           `json:"wan1"`
	Wan2              Wan           `json:"wan2"`
	WifiCaps          FlexInt       `json:"wifi_caps"`
	WlanNumSta        FlexInt       `json:"wlan-num_sta"` // UAP
	WlangroupIDNa     string        `json:"wlangroup_id_na"`
	WlangroupIDNg     string        `json:"wlangroup_id_ng"`
	XInformAuthkey    string        `json:"x_inform_authkey"`
}

type EthernetOverrides struct {
	Ifname       string `json:"ifname"`
	Networkgroup string `json:"networkgroup"`
}

// NetworkTable is the list of networks on a gateway.
// Not all gateways have all features.
type NetworkTable []struct {
	ActiveDhcpLeaseCount   FlexInt   `json:"active_dhcp_lease_count"`
	AttrHiddenID           string    `json:"attr_hidden_id"`
	AttrNoDelete           FlexBool  `json:"attr_no_delete"`
	AutoScaleEnabled       FlexBool  `json:"auto_scale_enabled"`
	DPIStatsTable          *DPITable `json:"dpistats_table"`
	DhcpRelayEnabled       FlexBool  `json:"dhcp_relay_enabled"`
	DhcpdDNS1              string    `json:"dhcpd_dns_1"`
	DhcpdDNS2              string    `json:"dhcpd_dns_2"`
	DhcpdDNS3              string    `json:"dhcpd_dns_3"`
	DhcpdDNS4              string    `json:"dhcpd_dns_4"`
	DhcpdDNSEnabled        FlexBool  `json:"dhcpd_dns_enabled"`
	DhcpdEnabled           FlexBool  `json:"dhcpd_enabled"`
	DhcpdGatewayEnabled    FlexBool  `json:"dhcpd_gateway_enabled"`
	DhcpdLeasetime         FlexInt   `json:"dhcpd_leasetime"`
	DhcpdStart             string    `json:"dhcpd_start"`
	DhcpdStop              string    `json:"dhcpd_stop"`
	DhcpdTimeOffsetEnabled FlexBool  `json:"dhcpd_time_offset_enabled"`
	Dhcpdv6Enabled         FlexBool  `json:"dhcpdv6_enabled"`
	DomainName             string    `json:"domain_name"`
	Enabled                FlexBool  `json:"enabled"`
	GatewayInterfaceName   string    `json:"gateway_interface_name"`
	ID                     string    `fake:"{uuid}"                    json:"_id"`
	IP                     string    `fake:"{ipv4address}"             json:"ip"`
	IPSubnet               string    `json:"ip_subnet"`
	Ipv6InterfaceType      string    `json:"ipv6_interface_type"`
	Ipv6PdStart            string    `json:"ipv6_pd_start"`
	Ipv6PdStop             string    `json:"ipv6_pd_stop"`
	Ipv6RaEnabled          FlexBool  `json:"ipv6_ra_enabled"`
	IsGuest                FlexBool  `json:"is_guest"`
	IsNat                  FlexBool  `json:"is_nat"`
	LteLanEnabled          FlexBool  `json:"lte_lan_enabled"`
	Mac                    string    `fake:"{macaddress}"              json:"mac"`
	Name                   string    `fake:"{animal}"                  json:"name"`
	Networkgroup           string    `json:"networkgroup"`
	NumSta                 FlexInt   `json:"num_sta"`
	Purpose                string    `json:"purpose"`
	RxBytes                FlexInt   `json:"rx_bytes"`
	RxPackets              FlexInt   `json:"rx_packets"`
	SiteID                 string    `fake:"{uuid}"                    json:"site_id"`
	TxBytes                FlexInt   `json:"tx_bytes"`
	TxPackets              FlexInt   `json:"tx_packets"`
	Up                     FlexBool  `json:"up"`
	VlanEnabled            FlexBool  `json:"vlan_enabled"`
}

// Storage is hard drive into for a device with storage.
type Storage struct {
	MountPoint string  `json:"mount_point"`
	Name       string  `fake:"{randomstring:[foo,bar]}" json:"name"`
	Size       FlexInt `json:"size"`
	Type       string  `fake:"{randomstring:[foo,bar]}" json:"type"`
	Used       FlexInt `json:"used"`
}

type Temperature struct {
	Name  string  `fake:"{randomstring:[cpu,sys]}" json:"name"`
	Type  string  `fake:"{randomstring:[foo,bar]}" json:"type"`
	Value float64 `json:"value"`
}

// UDMStat holds the "stat" data for a dream machine.
// A dream machine is a USG + USW + Controller.
type UDMStat struct {
	*Gw `json:"gw"`
	*Sw `json:"sw"`
	*Ap `json:"ap,omitempty"`
}
