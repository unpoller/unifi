package unifi

import "encoding/json"

// PDU is the Smart Power PDU line of products
type PDU struct {
	site                     *Site
	ID                       string           `fake:"{uuid}"                      json:"_id"`
	AdoptIP                  string           `fake:"{ipv4address}"               json:"adopt_ip"`
	AdoptURL                 string           `fake:"{url}"                       json:"adopt_url"`
	AdoptableWhenUpgraded    FlexBool         `json:"adoptable_when_upgraded"`
	Adopted                  FlexBool         `json:"adopted"`
	Anomalies                FlexInt          `json:"anomalies"`
	AnonID                   string           `fake:"{uuid}"                      json:"anon_id"`
	Architecture             string           `json:"architecture"`
	BoardRev                 FlexInt          `json:"board_rev"`
	Bytes                    FlexInt          `json:"bytes"`
	CfgVersion               string           `fake:"{appversion}"                json:"cfgversion"`
	ConfigNetwork            *ConfigNetwork   `json:"config_network"`
	ConnectRequestIP         string           `fake:"{ipv4address}"               json:"connect_request_ip"`
	ConnectRequestPort       FlexInt          `json:"connect_request_port"`
	ConnectedAt              FlexInt          `json:"connected_at"`
	ConnectionNetworkName    string           `json:"connection_network_name"`
	Default                  FlexBool         `json:"default"`
	DeviceID                 string           `fake:"{uuid}"                      json:"device_id"`
	DiscoveredVia            string           `json:"discovered_via"`
	DisplayableVersion       string           `fake:"{appversion}"                json:"displayable_version"`
	Dot1xPortCtrlEnabled     FlexBool         `json:"dot1x_portctrl_enabled"`
	DownlinkTable            []*DownlinkTable `fakesize:"5"                       json:"downlink_table"`
	EthernetTable            []*EthernetTable `fakesize:"5"                       json:"ethernet_table"`
	FlowctrlEnabled          FlexBool         `json:"flowctrl_enabled"`
	FwCaps                   FlexInt          `json:"fw_caps"`
	GatewayMac               string           `fake:"{macaddress}"                json:"gateway_mac"`
	GuestNumSta              FlexInt          `json:"guest-num_sta"`
	HasFan                   FlexBool         `json:"has_fan"`
	HasTemperature           FlexBool         `json:"has_temperature"`
	HashID                   string           `json:"hash_id"`
	HwCaps                   FlexInt          `json:"hw_caps"`
	InformIP                 string           `fake:"{ipv4address}"               json:"inform_ip"`
	InformURL                string           `fake:"{url}"                       json:"inform_url"`
	Internet                 FlexBool         `json:"internet"`
	IP                       string           `fake:"{ipv4address}"               json:"ip"`
	JumboframeEnabled        FlexBool         `json:"jumboframe_enabled"`
	KernelVersion            string           `fake:"{appversion}"                json:"kernel_version"`
	KnownCfgVersion          string           `fake:"{appversion}"                json:"known_cfgversion"`
	LastSeen                 FlexInt          `json:"last_seen"`
	LastUplink               Uplink           `json:"last_uplink"`
	LcmBrightness            FlexInt          `json:"lcm_brightness"`
	LcmBrightnessOverride    FlexBool         `json:"lcm_brightness_override"`
	LcmNightModeBegins       string           `json:"lcm_night_mode_begins"`
	LcmNightModeEnabled      FlexBool         `json:"lcm_night_mode_enabled"`
	LcmNightModeEnds         string           `json:"lcm_night_mode_ends"`
	LicenseState             string           `json:"license_state"`
	Locating                 FlexBool         `json:"locating"`
	Mac                      string           `fake:"{macaddress}"                json:"mac"`
	ManufacturerID           FlexInt          `json:"manufacturer_id"`
	MinIfnromIntervalSeconds FlexInt          `json:"min_inform_interval_seconds"`
	Model                    string           `json:"model"`
	ModelInEOL               FlexBool         `json:"model_in_eol"`
	ModelInLTS               FlexBool         `json:"model_in_lts"`
	ModelIncompatible        FlexBool         `json:"model_incompatible"`
	Name                     string           `json:"name"`
	NextInterval             FlexInt          `json:"next_interval"`
	NumSta                   FlexInt          `json:"num_sta"`
	OutletACPowerBudget      FlexInt          `json:"outlet_ac_power_budget"`
	OutletACPowerConsumption FlexInt          `json:"outlet_ac_power_consumption"`
	OutletEnabled            FlexBool         `json:"outlet_enabled"`
	OutletOverrides          []OutletOverride `fakesize:"5"                       json:"outlet_overrides"`
	OutletTable              []OutletTable    `fakesize:"5"                       json:"outlet_table"`
	Overheating              FlexBool         `json:"overheating"`
	PortTable                []Port           `fakesize:"5"                       json:"port_table"`
	PowerSource              FlexInt          `json:"power_source"`
	PowerSourceCtrlEnabled   FlexBool         `json:"power_source_ctrl_enabled"`
	PrevNonBusyState         FlexInt          `json:"prev_non_busy_state"`
	ProvisionedAt            FlexInt          `json:"provisioned_at"`
	RequiredVersion          string           `fake:"{appversion}"                json:"required_version"`
	RollUpgrade              FlexBool         `json:"rollupgrade"`
	RxBytes                  FlexInt          `json:"rx_bytes"`
	Satisfaction             FlexInt          `json:"satisfaction"`
	Serial                   string           `json:"serial"`
	SetupID                  string           `fake:"{uuid}"                      json:"setup_id"`
	SiteID                   string           `fake:"{uuid}"                      json:"site_id"`
	SiteName                 string           `json:"site_name"`
	SourceName               string           `json:"source_name"`
	StartConnectedMillis     FlexInt          `json:"start_connected_millis"`
	StartDisconnectedMillis  FlexInt          `json:"start_disconnected_millis"`
	StartupTimestamp         FlexInt          `json:"startup_timestamp"`
	Stat                     PDUStat          `json:"stat"`
	State                    FlexInt          `json:"state"`
	StpPriority              FlexInt          `json:"stp_priority"`
	StpVersion               string           `fake:"{appversion}"                json:"stp_version"`
	SwitchCaps               *SwitchCaps      `json:"switch_caps"`
	SysErrorCaps             FlexInt          `json:"sys_error_caps"`
	SysStats                 SysStats         `json:"sys_stats"`
	SyslogKey                string           `json:"syslog_key"`
	SystemStats              SystemStats      `json:"system-stats"`
	TotalMaxPower            FlexInt          `json:"total_max_power"`
	TwoPhaseAdopt            FlexBool         `json:"two_phase_adopt"`
	TxBytes                  FlexInt          `json:"tx_bytes"`
	Type                     string           `fake:"{lexify:pdu}"                json:"type"`
	Unsupported              FlexBool         `json:"unsupported"`
	UnsupportedReason        FlexInt          `json:"unsupported_reason"`
	Upgradeable              FlexBool         `json:"upgradable"`
	Uplink                   Uplink           `json:"uplink"`
	UplinkDepth              FlexBool         `json:"uplink_depth"`
	Uptime                   FlexInt          `json:"uptime"`
	UserNumSta               FlexInt          `json:"user-num_sta"`
	Version                  string           `fake:"{appversion}"                json:"version"`
}

// OutletOverride hold the PDU outlet override data.
type OutletOverride struct {
	CycleEnabled FlexBool `json:"cycle_enabled"`
	Index        FlexInt  `json:"index"`
	Name         string   `json:"name"`
	RelayState   FlexBool `json:"relay_state"`
}

// OutletTable hold the PDU outlet data.
type OutletTable struct {
	CycleEnabled      FlexBool `json:"cycle_enabled"`
	Index             FlexInt  `json:"index"`
	Name              string   `json:"name"`
	RelayState        FlexBool `json:"relay_state"`
	OutletCaps        FlexInt  `json:"outlet_caps"`
	OutletCurrent     FlexInt  `json:"outlet_current"`
	OutletPower       FlexInt  `json:"outlet_power"`
	OutletPowerFactor FlexInt  `json:"outlet_power_factor"`
	OutletVoltage     FlexInt  `json:"outlet_voltage"`
}

// PDUStat holds the "stat" data for a pdu.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type PDUStat struct {
	*Sw
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Switch Stat data.
func (v *PDUStat) UnmarshalJSON(data []byte) error {
	var n struct {
		Sw `json:"sw"`
	}

	v.Sw = &n.Sw

	err := json.Unmarshal(data, v.Sw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}

	return nil
}
