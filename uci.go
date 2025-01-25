package unifi

import "encoding/json"

// UCIStat holds the "stat" data for a switch.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type UCIStat struct {
	*Sw
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Switch Stat data.
func (v *UCIStat) UnmarshalJSON(data []byte) error {
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

// CIStateTable holds the CI state table
type CIStateTable struct {
	CIState      string `json:"ci_state"`
	CISwDlStatus string `json:"ci_sw_dl_status"`
	CIMac        string `json:"ci_mac"`
	CIVersion    string `json:"ci_version"`
	CIMode       string `json:"ci_mode"`
}

type UCI struct {
	AdoptableWhenUpgraded    FlexBool        `json:"adoptable_when_upgraded"`
	Adopted                  FlexBool        `json:"adopted"`
	AdoptedAt                FlexInt         `json:"adopted_at"`
	AdoptedByClient          string          `json:"adopted_by_client"`
	AdoptionCompleted        FlexBool        `json:"adoption_completed"`
	AnonID                   string          `json:"anon_id"`
	Architecture             string          `json:"architecture"`
	BleCaps                  FlexInt         `json:"ble_caps"`
	BoardRev                 FlexInt         `json:"board_rev"`
	Bytes                    FlexInt         `json:"bytes"`
	CfgVersion               string          `json:"cfgversion"`
	CiStateTable             *CIStateTable   `json:"ci_state_table"`
	ConfigNetwork            *ConfigNetwork  `json:"config_network"`
	ConnectedAt              FlexInt         `json:"connected_at"`
	DeviceID                 string          `json:"device_id"`
	DisconnectedAt           FlexInt         `json:"disconnected_at"`
	DisplayableVersion       string          `json:"displayable_version"`
	DownlinkTable            []DownlinkTable `fakesize:"1"                       json:"downlink_table"`
	EthernetTable            []EthernetTable `fakesize:"1"                       json:"ethernet_table"`
	FwCaps                   FlexInt         `json:"fw_caps"`
	HashID                   string          `json:"hash_id"`
	HwCaps                   FlexInt         `json:"hw_caps"`
	ID                       string          `json:"_id"`
	IP                       string          `json:"ip"`
	InformIP                 string          `json:"inform_ip"`
	InformURL                string          `json:"inform_url"`
	Internet                 FlexBool        `json:"internet"`
	IsAccessPoint            FlexBool        `json:"is_access_point"`
	IspName                  string          `json:"isp_name"`
	KernelVersion            string          `json:"kernel_version"`
	KnownCfgVersion          string          `json:"known_cfgversion"`
	LLDPTable                []LLDPTable     `fakesize:"1"                       json:"lldp_table"`
	LastSeen                 FlexInt         `json:"last_seen"`
	LcmBrightness            FlexInt         `json:"lcm_brightness"`
	LicenseState             string          `json:"license_state"`
	Locating                 FlexBool        `json:"locating"`
	Mac                      string          `json:"mac"`
	ManufacturerID           FlexInt         `json:"manufacturer_id"`
	MinInformIntervalSeconds FlexInt         `json:"min_inform_interval_seconds"`
	Model                    string          `json:"model"`
	ModelInEol               FlexBool        `json:"model_in_eol"`
	ModelInLts               FlexBool        `json:"model_in_lts"`
	ModelIncompatible        FlexBool        `json:"model_incompatible"`
	Name                     string          `json:"name"`
	NextInterval             FlexInt         `json:"next_interval"`
	NumSta                   FlexInt         `json:"num_sta"`
	PortTable                []Port          `fakesize:"1"                       json:"port_table"`
	PrevNonBusyState         FlexInt         `json:"prev_non_busy_state"`
	ProvisionedAt            FlexInt         `json:"provisioned_at"`
	RebootDuration           FlexInt         `json:"reboot_duration"`
	RequiredVersion          string          `json:"required_version"`
	RollUpgrade              FlexBool        `json:"rollupgrade"`
	RxBytes                  FlexInt         `json:"rx_bytes"`
	SafeForAutoUpgrade       FlexBool        `json:"safe_for_autoupgrade"`
	Serial                   string          `json:"serial"`
	SetupID                  string          `json:"setup_id"`
	SiteID                   string          `fake:"{uuid}"                      json:"site_id"`
	SiteName                 string          `fake:"{company}"                   json:"site_name"`
	SourceName               string          `fake:"{animal}"                    json:"source_name"`
	StartConnectedMillis     FlexInt         `json:"start_connected_millis"`
	StartDisconnectedMillis  FlexInt         `json:"start_disconnected_millis"`
	StartupTimestamp         FlexInt         `json:"startup_timestamp"`
	Stat                     *UCIStat        `json:"stat"`
	State                    FlexInt         `json:"state"`
	SysErrorCaps             FlexInt         `json:"sys_error_caps"`
	SysStats                 *SysStats       `json:"sys_stats"`
	SyslogKey                string          `json:"syslog_key"`
	SystemStats              *SystemStats    `json:"system_stats"`
	TwoPhaseAdopt            FlexBool        `json:"two_phase_adopt"`
	TxBytes                  FlexInt         `json:"tx_bytes"`
	Type                     string          `fake:"{lexify:uci}"                json:"type"`
	UciConnectTo             string          `json:"uci_connect_to"`
	Unsupported              FlexBool        `json:"unsupported"`
	UnsupportedReason        FlexInt         `json:"unsupported_reason"`
	Upgradable               FlexBool        `json:"upgradable"`
	UpgradeDuration          FlexInt         `json:"upgrade_duration"`
	Uptime                   FlexInt         `json:"uptime"`
	Uptime0                  FlexInt         `json:"_uptime"`
	Version                  string          `json:"version"`
	WanNetworkGroup          string          `json:"wan_networkgroup"`
	WanPort                  string          `json:"wan_port"`
	site                     *Site           `fake:"-"`
}
