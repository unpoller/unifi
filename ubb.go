package unifi

import (
	"encoding/json"
	"time"
)

// UBBStat holds the "stat" data for a switch.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type UBBStat struct {
	*Bb
}

type Bb struct {
	SiteID                        string    `json:"site_id"`
	O                             string    `json:"o"`
	Oid                           string    `json:"oid"`
	Bb                            string    `json:"bb"`
	Time                          FlexInt   `json:"time"`
	Datetime                      time.Time `fake:"{recent_time}"                     json:"datetime"`
	UserWifi0RxPackets            FlexInt   `json:"user-wifi0-rx_packets"`
	UserTerra2RxPackets           FlexInt   `json:"user-terra2-rx_packets"`
	UserRxPackets                 FlexInt   `json:"user-rx_packets"`
	Wifi0RxPackets                FlexInt   `json:"wifi0-rx_packets"`
	Terra2RxPackets               FlexInt   `json:"terra2-rx_packets"`
	RxPackets                     FlexInt   `json:"rx_packets"`
	UserWifi0RxBytes              FlexInt   `json:"user-wifi0-rx_bytes"`
	UserTerra2RxBytes             FlexInt   `json:"user-terra2-rx_bytes"`
	UserRxBytes                   FlexInt   `json:"user-rx_bytes"`
	Wifi0RxBytes                  FlexInt   `json:"wifi0-rx_bytes"`
	Terra2RxBytes                 FlexInt   `json:"terra2-rx_bytes"`
	RxBytes                       FlexInt   `json:"rx_bytes"`
	UserWifi0RxErrors             FlexInt   `json:"user-wifi0-rx_errors"`
	UserTerra2RxErrors            FlexInt   `json:"user-terra2-rx_errors"`
	UserRxErrors                  FlexInt   `json:"user-rx_errors"`
	Wifi0RxErrors                 FlexInt   `json:"wifi0-rx_errors"`
	Terra2RxErrors                FlexInt   `json:"terra2-rx_errors"`
	RxErrors                      FlexInt   `json:"rx_errors"`
	UserWifi0RxDropped            FlexInt   `json:"user-wifi0-rx_dropped"`
	UserTerra2RxDropped           FlexInt   `json:"user-terra2-rx_dropped"`
	UserRxDropped                 FlexInt   `json:"user-rx_dropped"`
	Wifi0RxDropped                FlexInt   `json:"wifi0-rx_dropped"`
	Terra2RxDropped               FlexInt   `json:"terra2-rx_dropped"`
	RxDropped                     FlexInt   `json:"rx_dropped"`
	UserWifi0RxCrypts             FlexInt   `json:"user-wifi0-rx_crypts"`
	UserTerra2RxCrypts            FlexInt   `json:"user-terra2-rx_crypts"`
	UserRxCrypts                  FlexInt   `json:"user-rx_crypts"`
	Wifi0RxCrypts                 FlexInt   `json:"wifi0-rx_crypts"`
	Terra2RxCrypts                FlexInt   `json:"terra2-rx_crypts"`
	RxCrypts                      FlexInt   `json:"rx_crypts"`
	UserWifi0RxFrags              FlexInt   `json:"user-wifi0-rx_frags"`
	UserTerra2RxFrags             FlexInt   `json:"user-terra2-rx_frags"`
	UserRxFrags                   FlexInt   `json:"user-rx_frags"`
	Wifi0RxFrags                  FlexInt   `json:"wifi0-rx_frags"`
	Terra2RxFrags                 FlexInt   `json:"terra2-rx_frags"`
	RxFrags                       FlexInt   `json:"rx_frags"`
	UserWifi0TxPackets            FlexInt   `json:"user-wifi0-tx_packets"`
	UserTerra2TxPackets           FlexInt   `json:"user-terra2-tx_packets"`
	UserTxPackets                 FlexInt   `json:"user-tx_packets"`
	Wifi0TxPackets                FlexInt   `json:"wifi0-tx_packets"`
	Terra2TxPackets               FlexInt   `json:"terra2-tx_packets"`
	TxPackets                     FlexInt   `json:"tx_packets"`
	UserWifi0TxBytes              FlexInt   `json:"user-wifi0-tx_bytes"`
	UserTerra2TxBytes             FlexInt   `json:"user-terra2-tx_bytes"`
	UserTxBytes                   FlexInt   `json:"user-tx_bytes"`
	Wifi0TxBytes                  FlexInt   `json:"wifi0-tx_bytes"`
	Terra2TxBytes                 FlexInt   `json:"terra2-tx_bytes"`
	TxBytes                       FlexInt   `json:"tx_bytes"`
	UserWifi0TxErrors             FlexInt   `json:"user-wifi0-tx_errors"`
	UserTerra2TxErrors            FlexInt   `json:"user-terra2-tx_errors"`
	UserTxErrors                  FlexInt   `json:"user-tx_errors"`
	Wifi0TxErrors                 FlexInt   `json:"wifi0-tx_errors"`
	Terra2TxErrors                FlexInt   `json:"terra2-tx_errors"`
	TxErrors                      FlexInt   `json:"tx_errors"`
	UserWifi0TxDropped            FlexInt   `json:"user-wifi0-tx_dropped"`
	UserTerra2TxDropped           FlexInt   `json:"user-terra2-tx_dropped"`
	UserTxDropped                 FlexInt   `json:"user-tx_dropped"`
	Wifi0TxDropped                FlexInt   `json:"wifi0-tx_dropped"`
	Terra2TxDropped               FlexInt   `json:"terra2-tx_dropped"`
	TxDropped                     FlexInt   `json:"tx_dropped"`
	UserWifi0TxRetries            FlexInt   `json:"user-wifi0-tx_retries"`
	UserTerra2TxRetries           FlexInt   `json:"user-terra2-tx_retries"`
	UserTxRetries                 FlexInt   `json:"user-tx_retries"`
	Wifi0TxRetries                FlexInt   `json:"wifi0-tx_retries"`
	Terra2TxRetries               FlexInt   `json:"terra2-tx_retries"`
	TxRetries                     FlexInt   `json:"tx_retries"`
	UserWifi0MacFilterRejections  FlexInt   `json:"user-wifi0-mac_filter_rejections"`
	UserTerra2MacFilterRejections FlexInt   `json:"user-terra2-mac_filter_rejections"`
	UserMacFilterRejections       FlexInt   `json:"user-mac_filter_rejections"`
	Wifi0MacFilterRejections      FlexInt   `json:"wifi0-mac_filter_rejections"`
	Terra2MacFilterRejections     FlexInt   `json:"terra2-mac_filter_rejections"`
	MacFilterRejections           FlexInt   `json:"mac_filter_rejections"`
	UserWifi0WifiTxAttempts       FlexInt   `json:"user-wifi0-wifi_tx_attempts"`
	UserTerra2WifiTxAttempts      FlexInt   `json:"user-terra2-wifi_tx_attempts"`
	UserWifiTxAttempts            FlexInt   `json:"user-wifi_tx_attempts"`
	Wifi0WifiTxAttempts           FlexInt   `json:"wifi0-wifi_tx_attempts"`
	Terra2WifiTxAttempts          FlexInt   `json:"terra2-wifi_tx_attempts"`
	WifiTxAttempts                FlexInt   `json:"wifi_tx_attempts"`
	UserWifi0WifiTxDropped        FlexInt   `json:"user-wifi0-wifi_tx_dropped"`
	UserTerra2WifiTxDropped       FlexInt   `json:"user-terra2-wifi_tx_dropped"`
	UserWifiTxDropped             FlexInt   `json:"user-wifi_tx_dropped"`
	Wifi0WifiTxDropped            FlexInt   `json:"wifi0-wifi_tx_dropped"`
	Terra2WifiTxDropped           FlexInt   `json:"terra2-wifi_tx_dropped"`
	WifiTxDropped                 FlexInt   `json:"wifi_tx_dropped"`
	Bytes                         FlexInt   `json:"bytes"`
	Duration                      FlexInt   `json:"duration"`
	UserWifi0Ath0RxPackets        FlexInt   `json:"user-wifi0-ath0-rx_packets"`
	UserWifi0Ath0RxBytes          FlexInt   `json:"user-wifi0-ath0-rx_bytes"`
	UserWifi0Ath0TxPackets        FlexInt   `json:"user-wifi0-ath0-tx_packets"`
	UserWifi0Ath0TxBytes          FlexInt   `json:"user-wifi0-ath0-tx_bytes"`
	UserTerra2Wlan0RxPackets      FlexInt   `json:"user-terra2-wlan0-rx_packets"`
	UserTerra2Wlan0RxBytes        FlexInt   `json:"user-terra2-wlan0-rx_bytes"`
	UserTerra2Wlan0TxPackets      FlexInt   `json:"user-terra2-wlan0-tx_packets"`
	UserTerra2Wlan0TxBytes        FlexInt   `json:"user-terra2-wlan0-tx_bytes"`
	UserTerra2Wlan0TxDropped      FlexInt   `json:"user-terra2-wlan0-tx_dropped"`
	UserTerra2Wlan0RxErrors       FlexInt   `json:"user-terra2-wlan0-rx_errors"`
	UserTerra2Wlan0TxErrors       FlexInt   `json:"user-terra2-wlan0-tx_errors"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Switch Stat data.
func (v *UBBStat) UnmarshalJSON(data []byte) error {
	var n struct {
		Bb `json:"bb"`
	}

	v.Bb = &n.Bb

	err := json.Unmarshal(data, v.Bb) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}

	return nil
}

type AntennaTable struct {
	Default   FlexBool   `json:"default"`
	ID        FlexInt    `json:"id"`
	Name      FlexString `json:"name"`
	Wifi0Gain FlexInt    `json:"wifi0_gain"`
	Wifi1Gain FlexInt    `json:"wifi1_gain"`
}

type UBB struct {
	AdoptableWhenUpgraded     FlexBool        `json:"adoptable_when_upgraded"`
	Adopted                   FlexBool        `json:"adopted"`
	AdoptedAt                 FlexInt         `json:"adopted_at"`
	AdoptedByClient           string          `json:"adopted_by_client"`
	AdoptionCompleted         FlexBool        `json:"adoption_completed"`
	AnonID                    string          `json:"anon_id"`
	AntennaTable              []AntennaTable  `fakesize:"1"                        json:"antenna_table"`
	Architecture              string          `json:"architecture"`
	AtfEnabled                FlexBool        `json:"atf_enabled"`
	BandsteeringMode          string          `json:"bandsteering_mode"`
	Bytes                     FlexInt         `json:"bytes"`
	BytesD                    FlexInt         `json:"bytes-d"`
	BytesR                    FlexInt         `json:"bytes-r"`
	Cfgversion                string          `json:"cfgversion"`
	ConfigNetwork             *ConfigNetwork  `json:"config_network"`
	ConnectRequestIP          string          `json:"connect_request_ip"`
	ConnectRequestPort        string          `json:"connect_request_port"`
	ConnectedAt               FlexInt         `json:"connected_at"`
	ConnectionNetworkID       string          `json:"connection_network_id"`
	ConnectionNetworkName     string          `json:"connection_network_name"`
	CountryCode               FlexInt         `json:"country_code"`
	DeviceID                  string          `json:"device_id"`
	DisconnectedAt            FlexInt         `json:"disconnected_at"`
	DisconnectionReason       string          `json:"disconnection_reason"`
	DisplayableVersion        string          `json:"displayable_version"`
	DownlinkLLDPMacs          []LLDPTable     `fakesize:"1"                        json:"downlink_lldp_macs"`
	DownlinkTable             []DownlinkTable `fakesize:"1"                        json:"downlink_table"`
	EthernetTable             []EthernetTable `fakesize:"1"                        json:"ethernet_table"`
	FixedApAvailable          FlexBool        `json:"fixed_ap_available"`
	FwCaps                    FlexInt         `json:"fw_caps"`
	GatewayMac                string          `json:"gateway_mac"`
	GeneralTemperature        FlexInt         `json:"general_temperature"`
	HasEth1                   FlexBool        `json:"has_eth1"`
	HasFan                    FlexBool        `json:"has_fan"`
	HasSpeaker                FlexBool        `json:"has_speaker"`
	HasTemperature            FlexBool        `json:"has_temperature"`
	HashID                    string          `json:"hash_id"`
	HwCaps                    FlexInt         `json:"hw_caps"`
	ID                        string          `json:"_id"`
	IP                        string          `json:"ip"`
	InformIP                  string          `json:"inform_ip"`
	InformURL                 string          `json:"inform_url"`
	IsAccessPoint             FlexBool        `json:"is_access_point"`
	Isolated                  FlexBool        `json:"isolated"`
	KernelVersion             string          `json:"kernel_version"`
	KnownCfgversion           string          `json:"known_cfgversion"`
	LastConnectionNetworkID   string          `json:"last_connection_network_id"`
	LastConnectionNetworkName string          `json:"last_connection_network_name"`
	LastSeen                  FlexInt         `json:"last_seen"`
	LastUplink                struct {
		UplinkMac        string  `json:"uplink_mac"`
		UplinkRemotePort FlexInt `json:"uplink_remote_port"`
		PortIdx          FlexInt `json:"port_idx"`
		Type             string  `json:"type"`
	} `json:"last_uplink"`
	LcmTrackerEnabled          FlexBool          `json:"lcm_tracker_enabled"`
	LedOverride                string            `json:"led_override"`
	LedOverrideColor           string            `json:"led_override_color"`
	LedOverrideColorBrightness FlexInt           `json:"led_override_color_brightness"`
	LicenseState               string            `json:"license_state"`
	LinkCapacity               FlexInt           `json:"link_capacity"`
	LinkQuality                FlexInt           `json:"link_quality"`
	LinkQualityCurrent         FlexInt           `json:"link_quality_current"`
	LldpTable                  []LLDPTable       `json:"lldp_table"`
	Locating                   FlexBool          `json:"locating"`
	Mac                        string            `json:"mac"`
	MeshStaVapEnabled          FlexBool          `json:"mesh_sta_vap_enabled"`
	MgmtNetworkID              string            `json:"mgmt_network_id"`
	Mode                       string            `json:"mode"`
	Model                      string            `json:"model"`
	ModelInEol                 FlexBool          `json:"model_in_eol"`
	ModelInLts                 FlexBool          `json:"model_in_lts"`
	ModelIncompatible          FlexBool          `json:"model_incompatible"`
	Name                       string            `json:"name"`
	NextInterval               FlexInt           `json:"next_interval"`
	NumSta                     FlexInt           `json:"num_sta"`
	P2PStats                   *P2PStats         `json:"p2p_stats"`
	PeerUbb                    *UBB              `fake:"-"                             json:"peer_ubb"`
	ProvisionedAt              FlexInt           `json:"provisioned_at"`
	RadioTable                 []RadioTable      `fakesize:"1"                         json:"radio_table"`
	RadioTableStats            []RadioTableStats `fakesize:"1"                         json:"radio_table_stats"`
	RebootDuration             FlexInt           `json:"reboot_duration"`
	RequiredVersion            string            `json:"required_version"`
	RxBytes                    FlexInt           `json:"rx_bytes"`
	RxBytesD                   FlexInt           `json:"rx_bytes-d"`
	SafeForAutoupgrade         FlexBool          `json:"safe_for_autoupgrade"`
	Serial                     string            `json:"serial"`
	SetupID                    string            `json:"setup_id"`
	SiteID                     string            `fake:"{uuid}"                        json:"site_id"`
	SiteName                   string            `fake:"{company}"                     json:"site_name"`
	SourceName                 string            `fake:"{animal}"                      json:"source_name"`
	StartConnectedMillis       FlexInt           `json:"start_connected_millis"`
	StartDisconnectedMillis    FlexInt           `json:"start_disconnected_millis"`
	StartupTimestamp           FlexInt           `json:"startup_timestamp"`
	Stat                       *UBBStat          `json:"stat"`
	State                      FlexInt           `json:"state"`
	SysErrorCaps               FlexInt           `json:"sys_error_caps"`
	SysStats                   *SysStats         `json:"sys_stats"`
	SyslogKey                  string            `json:"syslog_key"`
	SystemStats                *SystemStats      `json:"system-stats"`
	TwoPhaseAdopt              FlexBool          `json:"two_phase_adopt"`
	TxBytes                    FlexInt           `json:"tx_bytes"`
	TxBytesD                   FlexInt           `json:"tx_bytes-d"`
	Type                       string            `fake:"{lexify:ubb}"                  json:"type"`
	UbbBssid                   string            `json:"ubb_bssid"`
	UbbIsAp                    FlexBool          `json:"ubb_is_ap"`
	UbbPairID                  string            `json:"ubb_pair_id"`
	UbbPairName                string            `json:"ubb_pair_name"`
	UbbPeerState               FlexInt           `json:"ubb_peer_state"`
	UbbPsk                     string            `json:"ubb_psk"`
	UbbSsid                    string            `json:"ubb_ssid"`
	Unsupported                FlexBool          `json:"unsupported"`
	UnsupportedReason          FlexInt           `json:"unsupported_reason"`
	Upgradable                 FlexBool          `json:"upgradable"`
	UpgradeDuration            FlexInt           `json:"upgrade_duration"`
	UpgradeState               FlexInt           `json:"upgrade_state"`
	UpgradeToFirmware          string            `json:"upgrade_to_firmware"`
	Uplink                     *Uplink           `json:"uplink"`
	Uptime                     FlexInt           `json:"uptime"`
	Uptime0                    FlexInt           `json:"_uptime"`
	UserNumSta                 FlexInt           `json:"user-num_sta"`
	UserWlanNumSta             FlexInt           `json:"user-wlan-num_sta"`
	VapTable                   []VapTable        `fakesize:"1"                         json:"vap_table"`
	Version                    string            `json:"version"`
	WifiCaps                   FlexInt           `json:"wifi_caps"`
	WlangroupIDNa              string            `json:"wlangroup_id_na"`
	XAesGcm                    FlexBool          `json:"x_aes_gcm"`
	XAuthkey                   string            `json:"x_authkey"`
	XHasSSHHostkey             FlexBool          `json:"x_has_ssh_hostkey"`
	site                       *Site             `fake:"-"`

	// unknown types, but appear in user capture
	PortStats       []any `fake:"-" json:"port_stats"`
	SSHSessionTable []any `fake:"-" json:"ssh_session_table"`
	VwireTable      []any `fake:"-" json:"vwire_table"`
	VwireVapTable   []any `fake:"-" json:"vwire_vap_table"`
}

type P2PStats struct {
	RXRate     FlexInt `json:"rx_rate"`
	Throughput FlexInt `json:"throughput"`
	TXRate     FlexInt `json:"tx_rate"`
}

type LLDPTable struct {
	ChassisID     string   `json:"chassis_id"`
	IsWired       FlexBool `json:"is_wired"`
	LocalPortIdx  FlexInt  `json:"local_port_idx,omitempty"`
	LocalPortName string   `json:"local_port_name"`
	PortID        string   `json:"port_id"`
}
