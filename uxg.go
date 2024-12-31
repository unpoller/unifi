package unifi

// UXG represents all the data from the Ubiquiti Controller for a UniFi 10Gb Gateway.
// The UDM shares several structs/type-data with USW and USG.
type UXG struct {
	AdoptableWhenUpgraded      FlexBool             `json:"adoptable_when_upgraded"`
	Adopted                    FlexBool             `fake:"{constFlexBool:true}"          json:"adopted"`
	AdoptedByClient            string               `json:"adopted_by_client"`
	AdoptionCompleted          FlexBool             `json:"adoption_completed"`
	AnonID                     string               `json:"anon_id"`
	Architecture               string               `json:"architecture"`
	BoardRev                   FlexInt              `json:"board_rev"`
	Bytes                      FlexInt              `json:"bytes"`
	Cfgversion                 string               `json:"cfgversion"`
	ConfigNetwork              *ConfigNetwork       `json:"config_network"`
	ConfigNetworkLan           *ConfigNetworkLan    `json:"config_network_lan"`
	ConnectedAt                FlexInt              `json:"connected_at"`
	ConnectionNetworkName      string               `json:"connection_network_name"`
	ConnectRequestIP           string               `json:"connect_request_ip"`
	ConnectRequestPort         string               `json:"connect_request_port"`
	ConsideredLostAt           FlexInt              `json:"considered_lost_at"`
	DeviceID                   string               `fake:"{uuid}"                        json:"device_id"`
	DisplayableVersion         string               `json:"displayable_version"`
	DownlinkTable              []*DownlinkTable     `fakesize:"5"                         json:"downlink_table"`
	EthernetOverrides          []*EthernetOverrides `fakesize:"5"                         json:"ethernet_overrides"`
	EthernetTable              []*EthernetTable     `fakesize:"5"                         json:"ethernet_table"`
	FwCaps                     FlexInt              `json:"fw_caps"`
	GeoInfo                    map[string]*GeoInfo  `fakesize:"5"                         json:"geo_info"`
	GuestKicks                 FlexInt              `json:"guest_kicks"`
	GuestLanNumSta             FlexInt              `json:"guest-lan-num_sta"`
	GuestNumSta                FlexInt              `json:"guest-num_sta"`
	GuestToken                 string               `json:"guest_token"`
	GuestWlanNumSta            FlexInt              `json:"guest-wlan-num_sta"`
	HasEth1                    FlexBool             `json:"has_eth1"`
	HasFan                     FlexBool             `json:"has_fan"`
	HashID                     string               `json:"hash_id"`
	HasSpeaker                 FlexBool             `json:"has_speaker"`
	HasTemperature             FlexBool             `json:"has_temperature"`
	HwCaps                     FlexInt              `json:"hw_caps"`
	ID                         string               `fake:"{uuid}"                        json:"_id"`
	InformIP                   string               `fake:"{ipv4address}"                 json:"inform_ip"`
	InformURL                  string               `fake:"{url}"                         json:"inform_url"`
	Internet                   FlexBool             `json:"internet"`
	IP                         string               `fake:"{ipv4address}"                 json:"ip"`
	IsAccessPoint              FlexBool             `json:"is_access_point"`
	KernelVersion              string               `json:"kernel_version"`
	KnownCfgversion            string               `json:"known_cfgversion"`
	LanNumSta                  FlexInt              `json:"lan-num_sta"`
	LastSeen                   FlexInt              `json:"last_seen"`
	LastWanIP                  string               `json:"last_wan_ip"`
	LcmBrightness              FlexInt              `json:"lcm_brightness"`
	LcmBrightnessOverride      FlexBool             `json:"lcm_brightness_override"`
	LcmIdleTimeoutOverride     FlexBool             `json:"lcm_idle_timeout_override"`
	LcmNightModeBegins         string               `json:"lcm_night_mode_begins"`
	LcmNightModeEnabled        FlexBool             `json:"lcm_night_mode_enabled"`
	LcmNightModeEnds           string               `json:"lcm_night_mode_ends"`
	LedOverride                string               `json:"led_override"`
	LedOverrideColor           string               `json:"led_override_color"`
	LedOverrideColorBrightness FlexInt              `json:"led_override_color_brightness"`
	LedState                   *LedState            `json:"led_state"`
	LicenseState               string               `json:"license_state"`
	Locating                   FlexBool             `fake:"{constFlexBool:false}"         json:"locating"`
	Mac                        string               `fake:"{macaddress}"                  json:"mac"`
	ManufacturerID             FlexInt              `json:"manufacturer_id"`
	MinInformIntervalSeconds   FlexInt              `json:"min_inform_interval_seconds"`
	Model                      string               `json:"model"`
	ModelIncompatible          FlexBool             `json:"model_incompatible"`
	ModelInEol                 FlexBool             `json:"model_in_eol"`
	ModelInLts                 FlexBool             `json:"model_in_lts"`
	Name                       string               `fake:"{animal}"                      json:"name"`
	NetworkTable               NetworkTable         `json:"network_table"`
	NextHeartbeatAt            FlexInt              `json:"next_heartbeat_at"`
	NextInterval               FlexInt              `json:"next_interval"`
	NumDesktop                 FlexInt              `json:"num_desktop"`
	NumHandheld                FlexInt              `json:"num_handheld"`
	NumMobile                  FlexInt              `json:"num_mobile"`
	NumSta                     FlexInt              `json:"num_sta"`
	OutdoorModeOverride        string               `json:"outdoor_mode_override"`
	OutdoorPowerCycleEnabled   FlexBool             `json:"outlet_power_cycle_enabled"`
	Overheating                FlexBool             `json:"overheating"`
	PortTable                  []Port               `json:"port_table"`
	ProvisionedAt              FlexInt              `json:"provisioned_at"`
	RequiredVersion            string               `json:"required_version"`
	RollUpgrade                FlexBool             `json:"rollupgrade"`
	RulesetInterfaces          interface{}          `json:"ruleset_interfaces"`
	RxBytes                    FlexInt              `json:"rx_bytes"`
	Serial                     string               `json:"serial"`
	SetupID                    string               `json:"setup_id"`
	site                       *Site
	SiteID                     string                  `fake:"{uuid}"                    json:"site_id"`
	SiteName                   string                  `json:"-"`
	SourceName                 string                  `json:"-"`
	SpeedtestStatus            SpeedtestStatus         `json:"speedtest-status"`
	SpeedtestStatusSaved       FlexBool                `json:"speedtest-status-saved"`
	StartConnectedMillis       FlexInt                 `json:"start_connected_millis"`
	StartDisconnectedMillis    FlexInt                 `json:"start_disconnected_millis"`
	StartupTimestamp           FlexInt                 `json:"startup_timestamp"`
	Stat                       *UXGStat                `json:"stat"`
	State                      FlexInt                 `json:"state"`
	Storage                    []*Storage              `json:"storage"`
	SwitchCaps                 *SwitchCaps             `json:"switch_caps"`
	SyslogKey                  string                  `json:"syslog_key"`
	SysStats                   SysStats                `json:"sys_stats"`
	SystemStats                SystemStats             `json:"system-stats"`
	TeleportVersion            string                  `json:"teleport_version"`
	Temperatures               []Temperature           `json:"temperatures"`
	TwoPhaseAdopt              FlexBool                `json:"two_phase_adopt"`
	TxBytes                    FlexInt                 `json:"tx_bytes"`
	Type                       string                  `fake:"{lexify:uxg}"              json:"type"`
	UdapiCaps                  FlexInt                 `json:"udapi_caps"`
	UnderscoreUptime           FlexInt                 `json:"_uptime"`
	Unsupported                FlexBool                `json:"unsupported"`
	UnsupportedReason          FlexInt                 `json:"unsupported_reason"`
	UpgradeState               FlexInt                 `json:"upgrade_state"`
	Uplink                     Uplink                  `json:"uplink"`
	Uptime                     FlexInt                 `json:"uptime"`
	UptimeStats                map[string]*UptimeStats `json:"uptime_stats"`
	UserLanNumSta              FlexInt                 `json:"user-lan-num_sta"`
	UserNumSta                 FlexInt                 `json:"user-num_sta"`
	UserWlanNumSta             FlexInt                 `json:"user-wlan-num_sta"`
	UsgCaps                    FlexInt                 `json:"usg_caps"`
	Version                    string                  `json:"version"`
	Wan1                       Wan                     `json:"wan1"`
	Wan2                       Wan                     `json:"wan2"`
	WifiCaps                   FlexInt                 `json:"wifi_caps"`
	WlanNumSta                 FlexInt                 `json:"wlan-num_sta"`
}

// ConfigNetworkLan is part of a UXG, maybe others.
type ConfigNetworkLan struct {
	DhcpEnabled FlexBool `json:"dhcp_enabled"`
	Vlan        int      `json:"vlan"`
}

// LedState is incuded with newer devices.
type LedState struct {
	Pattern string  `json:"pattern"`
	Tempo   FlexInt `json:"tempo"`
}

// GeoInfo is incuded with certain devices.
type GeoInfo struct {
	Accuracy        FlexInt `json:"accuracy"`
	Address         string  `fake:"{address}"        json:"address"`
	Asn             FlexInt `json:"asn"`
	City            string  `fake:"{city}"           json:"city"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
	CountryName     string  `fake:"{country}"        json:"country_name"`
	IspName         string  `json:"isp_name"`
	IspOrganization string  `json:"isp_organization"`
	Latitude        FlexInt `json:"latitude"`
	Longitude       FlexInt `json:"longitude"`
	Timezone        string  `json:"timezone"`
}

// UptimeStats is incuded with certain devices.
type UptimeStats struct {
	Availability   FlexInt `json:"availability"`
	LatencyAverage FlexInt `json:"latency_average"`
	TimePeriod     FlexInt `json:"time_period"`
}

// UXGStat holds the "stat" data for a 10Gb gateway.
type UXGStat struct {
	*Gw `json:"gw"`
	*Sw `json:"sw"`
}
