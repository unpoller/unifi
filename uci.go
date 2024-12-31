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
	RequiredVersion          string           `json:"required_version"`
	PortTable                []*Port          `json:"port_table"`
	LicenseState             string           `json:"license_state"`
	TwoPhaseAdopt            FlexBool         `json:"two_phase_adopt"`
	ConnectedAt              FlexInt          `json:"connected_at"`
	Type                     string           `json:"type"`
	InformIP                 string           `json:"inform_ip"`
	CfgVersion               string           `json:"cfgversion"`
	BoardRev                 FlexInt          `json:"board_rev"`
	Mac                      string           `json:"mac"`
	SetupID                  string           `json:"setup_id"`
	ProvisionedAt            FlexInt          `json:"provisioned_at"`
	HwCaps                   FlexInt          `json:"hw_caps"`
	InformURL                string           `json:"inform_url"`
	RebootDuration           FlexInt          `json:"reboot_duration"`
	UpgradeDuration          FlexInt          `json:"upgrade_duration"`
	EthernetTable            []*EthernetTable `json:"ethernet_table"`
	ConfigNetwork            *ConfigNetwork   `json:"config_network"`
	Unsupported              FlexBool         `json:"unsupported"`
	BleCaps                  FlexInt          `json:"ble_caps"`
	SysErrorCaps             FlexInt          `json:"sys_error_caps"`
	SyslogKey                string           `json:"syslog_key"`
	Model                    string           `json:"model"`
	DisconnectedAt           FlexInt          `json:"disconnected_at"`
	Architecture             string           `json:"architecture"`
	ManufacturerID           FlexInt          `json:"manufacturer_id"`
	ModelIncompatible        FlexBool         `json:"model_incompatible"`
	IP                       string           `json:"ip"`
	ModelInEol               FlexBool         `json:"model_in_eol"`
	Version                  string           `json:"version"`
	AdoptionCompleted        FlexBool         `json:"adoption_completed"`
	UnsupportedReason        FlexInt          `json:"unsupported_reason"`
	CiStateTable             *CIStateTable    `json:"ci_state_table"`
	AnonID                   string           `json:"anon_id"`
	AdoptedByClient          string           `json:"adopted_by_client"`
	ModelInLts               FlexBool         `json:"model_in_lts"`
	KernelVersion            string           `json:"kernel_version"`
	Serial                   string           `json:"serial"`
	SiteID                   string           `json:"site_id"`
	AdoptedAt                FlexInt          `json:"adopted_at"`
	FwCaps                   FlexInt          `json:"fw_caps"`
	ID                       string           `json:"_id"`
	Adopted                  FlexBool         `json:"adopted"`
	Internet                 FlexBool         `json:"internet"`
	HashID                   string           `json:"hash_id"`
	DeviceID                 string           `json:"device_id"`
	State                    FlexInt          `json:"state"`
	StartDisconnectedMillis  FlexInt          `json:"start_disconnected_millis"`
	UciConnectTo             string           `json:"uci_connect_to"`
	WanPort                  string           `json:"wan_port"`
	WanNetworkGroup          string           `json:"wan_networkgroup"`
	IspName                  string           `json:"isp_name"`
	LastSeen                 FlexInt          `json:"last_seen"`
	MinInformIntervalSeconds FlexInt          `json:"min_inform_interval_seconds"`
	Upgradable               FlexBool         `json:"upgradable"`
	AdoptableWhenUpgraded    FlexBool         `json:"adoptable_when_upgraded"`
	RollUpgrade              FlexBool         `json:"rollupgrade"`
	KnownCfgVersion          string           `json:"known_cfgversion"`
	Uptime                   FlexInt          `json:"uptime"`
	Uptime0                  FlexInt          `json:"_uptime"`
	Locating                 FlexBool         `json:"locating"`
	StartConnectedMillis     FlexInt          `json:"start_connected_millis"`
	NextInterval             FlexInt          `json:"next_interval"`
	SysStats                 *SysStats        `json:"sys_stats"`
	SystemStats              *SystemStats     `json:"system_stats"`
	LLDPTable                []*LLDPTable     `json:"lldp_table"`
	DisplayableVersion       string           `json:"displayable_version"`
	StartupTimestamp         FlexInt          `json:"startup_timestamp"`
	IsAccessPoint            FlexBool         `json:"is_access_point"`
	SafeForAutoUpgrade       FlexBool         `json:"safe_for_autoupgrade"`
	DownlinkTable            []*DownlinkTable `json:"downlink_table"`
	PrevNonBusyState         FlexInt          `json:"prev_non_busy_state"`
	Name                     string           `json:"name"`
	LcmBrightness            FlexInt          `json:"lcm_brightness"`
	Stat                     *UCIStat         `json:"stat"`
	TxBytes                  FlexInt          `json:"tx_bytes"`
	RxBytes                  FlexInt          `json:"rx_bytes"`
	Bytes                    FlexInt          `json:"bytes"`
	NumSta                   FlexInt          `json:"num_sta"`
}
