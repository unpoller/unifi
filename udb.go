package unifi

import (
	"encoding/json"
	"time"
)

// UDB represents a UniFi Device Bridge (UDB) device.
// The UDB product range includes: UDB-Switch, UDB-Pro, UDB-Pro-Sector.
// UDB-Switch is a hybrid device combining switch ports (8 PoE ports:
// 7x 2.5GbE + 1x 10GbE) with WiFi 7 wireless bridge capability
// (5GHz + 6GHz radios).
type UDB struct {
	site                  *Site
	AdoptableWhenUpgraded FlexBool         `json:"adoptable_when_upgraded,omitempty"`
	Adopted               FlexBool         `fake:"{constFlexBool:true}"              json:"adopted"`
	AdoptionCompleted     FlexBool         `json:"adoption_completed"`
	Anomalies             FlexInt          `json:"anomalies"`
	Architecture          string           `json:"architecture"`
	BoardRev              FlexInt          `json:"board_rev"`
	Bytes                 FlexInt          `json:"bytes"`
	BytesD                FlexInt          `json:"bytes-d"`
	BytesR                FlexInt          `json:"bytes-r"`
	Cfgversion            string           `json:"cfgversion"`
	ConfigNetwork         *ConfigNetwork   `json:"config_network"`
	ConnectRequestIP      string           `json:"connect_request_ip"`
	ConnectRequestPort    string           `json:"connect_request_port"`
	ConnectedAt           FlexInt          `json:"connected_at"`
	ConnectionNetworkName string           `json:"connection_network_name"`
	CountryCode           FlexInt          `json:"country_code"`
	DeviceID              string           `json:"device_id"`
	DisplayableVersion    string           `json:"displayable_version"`
	Dot1XPortctrlEnabled  FlexBool         `json:"dot1x_portctrl_enabled"`
	DownlinkTable         []*DownlinkTable `json:"downlink_table"`
	EthernetTable         []*EthernetTable `json:"ethernet_table"`
	FanLevel              FlexInt          `json:"fan_level"`
	FlowctrlEnabled       FlexBool         `json:"flowctrl_enabled"`
	FwCaps                FlexInt          `json:"fw_caps"`
	GatewayMac            string           `json:"gateway_mac"`
	GeneralTemperature    FlexInt          `json:"general_temperature"`
	GuestNumSta           FlexInt          `json:"guest-num_sta"`
	GuestWlanNumSta       FlexInt          `json:"guest-wlan-num_sta"`
	HasFan                FlexBool         `json:"has_fan"`
	HasTemperature        FlexBool         `json:"has_temperature"`
	HwCaps                FlexInt          `json:"hw_caps"`
	ID                    string           `json:"_id"`
	IP                    string           `json:"ip"`
	InformIP              string           `json:"inform_ip"`
	InformURL             string           `json:"inform_url"`
	Internet              FlexBool         `json:"internet"`
	IsAccessPoint         FlexBool         `json:"is_access_point"`
	JumboframeEnabled     FlexBool         `json:"jumboframe_enabled"`
	KernelVersion         string           `json:"kernel_version"`
	KnownCfgversion       string           `json:"known_cfgversion"`
	LCMTrackerEnabled     FlexBool         `json:"lcm_tracker_enabled"`
	LastSeen              FlexInt          `json:"last_seen"`
	LastUplink            struct {
		UplinkMac string `json:"uplink_mac"`
	} `json:"last_uplink"`
	LedOverride         string   `json:"led_override"`
	LicenseState        string   `json:"license_state"`
	Locating            FlexBool `fake:"{constFlexBool:false}" json:"locating"`
	Mac                 string   `json:"mac"`
	Tags                []string `json:"tags"` // Device tags assigned to this device
	ManufacturerID      FlexInt  `json:"manufacturer_id"`
	MeshStaVapEnabled   FlexBool `json:"mesh_sta_vap_enabled"`
	Model               string   `json:"model"`
	ModelInEOL          FlexBool `json:"model_in_eol"`
	ModelInLTS          FlexBool `json:"model_in_lts"`
	ModelIncompatible   FlexBool `json:"model_incompatible"`
	Name                string   `fake:"{animal}"              json:"name"`
	NextInterval        FlexInt  `json:"next_interval"`
	NumSta              FlexInt  `json:"num_sta"`
	OutdoorModeOverride string   `json:"outdoor_mode_override"`
	Overheating         FlexBool `json:"overheating"`
	PortOverrides       []struct {
		Name       string  `json:"name,omitempty"`
		PoeMode    string  `json:"poe_mode,omitempty"`
		PortIdx    FlexInt `json:"port_idx"`
		PortconfID string  `json:"portconf_id"`
	} `json:"port_overrides"`
	PortTable               []Port          `json:"port_table"`
	PowerSource             string          `json:"power_source"`
	PowerSourceCtrlEnabled  FlexBool        `json:"power_source_ctrl_enabled"`
	PowerSourceVoltage      string          `json:"power_source_voltage"`
	PreviousNonBusyState    FlexInt         `json:"prev_non_busy_state"`
	ProvisionedAt           FlexInt         `json:"provisioned_at"`
	RadioTable              RadioTable      `json:"radio_table"`
	RadioTableStats         RadioTableStats `json:"radio_table_stats"`
	RequiredVersion         string          `json:"required_version"`
	Rollupgrade             FlexBool        `json:"rollupgrade,omitempty"`
	RxBytes                 FlexInt         `json:"rx_bytes"`
	RxBytesD                FlexInt         `json:"rx_bytes-d"`
	Satisfaction            FlexInt         `json:"satisfaction"`
	Serial                  string          `json:"serial"`
	SetupID                 string          `json:"setup_id"`
	SiteID                  string          `json:"site_id"`
	SiteName                string          `json:"-"`
	SourceName              string          `json:"-"`
	StartConnectedMillis    FlexInt         `json:"start_connected_millis"`
	StartDisconnectedMillis FlexInt         `json:"start_disconnected_millis"`
	StartupTimestamp        FlexInt         `json:"startup_timestamp"`
	Stat                    UDBStat         `json:"stat"`
	State                   FlexInt         `json:"state"`
	StpPriority             FlexInt         `json:"stp_priority"`
	StpVersion              string          `json:"stp_version"`
	SwitchCaps              *SwitchCaps     `json:"switch_caps"`
	SysErrorCaps            FlexInt         `json:"sys_error_caps"`
	SysStats                SysStats        `json:"sys_stats"`
	SystemStats             SystemStats     `json:"system-stats"`
	TotalMaxPower           FlexInt         `json:"total_max_power"`
	TotalRxBytes            FlexInt         `json:"total_rx_bytes"`
	TotalTxBytes            FlexInt         `json:"total_tx_bytes"`
	TwoPhaseAdopt           FlexBool        `json:"two_phase_adopt"`
	TxBytes                 FlexInt         `json:"tx_bytes"`
	TxBytesD                FlexInt         `json:"tx_bytes-d"`
	Type                    string          `fake:"{randomstring:[udb]}"      json:"type"`
	Unsupported             FlexBool        `json:"unsupported"`
	UnsupportedReason       FlexInt         `json:"unsupported_reason"`
	Upgradable              FlexBool        `json:"upgradable,omitempty"`
	Upgradeable             FlexBool        `json:"upgradeable"`
	Uplink                  Uplink          `json:"uplink"`
	UplinkDepth             FlexInt         `json:"uplink_depth"`
	Uptime                  FlexInt         `json:"uptime"`
	UserNumSta              FlexInt         `json:"user-num_sta"`
	UserWlanNumSta          FlexInt         `json:"user-wlan-num_sta"`
	VapTable                VapTable        `json:"vap_table"`
	Version                 string          `json:"version"`
	VwireEnabled            FlexBool        `json:"vwireEnabled"`
	VwireTable              []interface{}   `json:"vwire_table"`
	VwireVapTable           []interface{}   `json:"vwire_vap_table"`
	WifiCaps                FlexInt         `json:"wifi_caps"`
	WlangroupIDNa           string          `json:"wlangroup_id_na"`
	WlangroupIDNg           string          `json:"wlangroup_id_ng"`
}

// UDBStat holds the "stat" data for a UDB (UniFi Device Bridge) device.
// This is split out because of a JSON data format change from 5.10 to 5.11.
// Contains device bridge stats (Db) and switch stats (Sw).
type UDBStat struct {
	*Db
	*Sw
}

// Db holds UDB-specific statistics (device bridge metrics).
// Named "db" to match the JSON key in the UniFi API response.
type Db struct {
	SiteID      string    `fake:"{uuid}"        json:"site_id"`
	O           string    `json:"o"`
	Oid         string    `json:"oid"`
	Db          string    `json:"db"`
	Time        FlexInt   `json:"time"`
	Datetime    time.Time `fake:"{recent_time}" json:"datetime"`
	Bytes       FlexInt   `json:"bytes"`
	Duration    FlexInt   `json:"duration"`
	RxBytes     FlexInt   `json:"rx_bytes"`
	RxPackets   FlexInt   `json:"rx_packets"`
	RxErrors    FlexInt   `json:"rx_errors"`
	RxDropped   FlexInt   `json:"rx_dropped"`
	RxCrypts    FlexInt   `json:"rx_crypts"`
	RxFrags     FlexInt   `json:"rx_frags"`
	RxMulticast FlexInt   `json:"rx_multicast"`
	RxBroadcast FlexInt   `json:"rx_broadcast"`
	TxBytes     FlexInt   `json:"tx_bytes"`
	TxPackets   FlexInt   `json:"tx_packets"`
	TxErrors    FlexInt   `json:"tx_errors"`
	TxDropped   FlexInt   `json:"tx_dropped"`
	TxRetries   FlexInt   `json:"tx_retries"`
	TxMulticast FlexInt   `json:"tx_multicast"`
	TxBroadcast FlexInt   `json:"tx_broadcast"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted UDB Stat data.
// Supports both nested (5.11+) and flat (5.10) JSON formats.
func (v *UDBStat) UnmarshalJSON(data []byte) error {
	var nested struct {
		Db `json:"db"`
		Sw `json:"sw"`
	}

	v.Db = &nested.Db
	v.Sw = &nested.Sw

	// Try flat format first (controller version 5.10)
	err := json.Unmarshal(data, v.Sw)
	if err != nil {
		// Try nested format (controller version 5.11+)
		return json.Unmarshal(data, &nested)
	}

	// Also try to unmarshal db stats from flat format
	_ = json.Unmarshal(data, v.Db)

	return nil
}
