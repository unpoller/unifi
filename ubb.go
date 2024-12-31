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
	Datetime                      time.Time `json:"datetime"`
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

type UBB struct {
	SiteName                  string `json:"-"`
	SourceName                string `json:"-"`
	site                      *Site
	RequiredVersion           string           `json:"required_version"`
	HasSpeaker                FlexBool         `json:"has_speaker"`
	LicenseState              string           `json:"license_state"`
	MeshStaVapEnabled         FlexBool         `json:"mesh_sta_vap_enabled"`
	Type                      string           `json:"type"`
	SetupID                   string           `json:"setup_id"`
	HwCaps                    FlexInt          `json:"hw_caps"`
	RebootDuration            FlexInt          `json:"reboot_duration"`
	ConfigNetwork             *ConfigNetwork   `json:"config_network"`
	SyslogKey                 string           `json:"syslog_key"`
	Model                     string           `json:"model"`
	LcmTrackerEnabled         FlexBool         `json:"lcm_tracker_enabled"`
	BandsteeringMode          string           `json:"bandsteering_mode"`
	IP                        string           `json:"ip"`
	LastConnectionNetworkName string           `json:"last_connection_network_name"`
	UbbIsAp                   FlexBool         `json:"ubb_is_ap"`
	LedOverrideColor          string           `json:"led_override_color"`
	Version                   string           `json:"version"`
	UnsupportedReason         FlexInt          `json:"unsupported_reason"`
	AdoptionCompleted         FlexBool         `json:"adoption_completed"`
	AnonID                    string           `json:"anon_id"`
	LastConnectionNetworkID   string           `json:"last_connection_network_id"`
	CountryCode               FlexInt          `json:"country_code"`
	WlangroupIDNa             string           `json:"wlangroup_id_na"`
	UbbSsid                   string           `json:"ubb_ssid"`
	AntennaTable              []any            `json:"antenna_table"`
	WifiCaps                  FlexInt          `json:"wifi_caps"`
	SiteID                    string           `json:"site_id"`
	AdoptedAt                 FlexInt          `json:"adopted_at"`
	Name                      string           `json:"name"`
	FwCaps                    FlexInt          `json:"fw_caps"`
	ID                        string           `json:"_id"`
	UbbBssid                  string           `json:"ubb_bssid"`
	MgmtNetworkID             string           `json:"mgmt_network_id"`
	GatewayMac                string           `json:"gateway_mac"`
	AtfEnabled                FlexBool         `json:"atf_enabled"`
	RadioTable                []*RadioTable    `json:"radio_table"`
	ConnectedAt               FlexInt          `json:"connected_at"`
	TwoPhaseAdopt             FlexBool         `json:"two_phase_adopt"`
	InformIP                  string           `json:"inform_ip"`
	Cfgversion                string           `json:"cfgversion"`
	Mac                       string           `json:"mac"`
	ProvisionedAt             FlexInt          `json:"provisioned_at"`
	InformURL                 string           `json:"inform_url"`
	EthernetTable             []*EthernetTable `json:"ethernet_table"`
	UpgradeDuration           FlexInt          `json:"upgrade_duration"`
	Unsupported               FlexBool         `json:"unsupported"`
	SysErrorCaps              FlexInt          `json:"sys_error_caps"`
	LastUplink                struct {
		UplinkMac        string  `json:"uplink_mac"`
		UplinkRemotePort FlexInt `json:"uplink_remote_port"`
		PortIdx          FlexInt `json:"port_idx"`
		Type             string  `json:"type"`
	} `json:"last_uplink"`
	LedOverride                string             `json:"led_override"`
	DisconnectedAt             FlexInt            `json:"disconnected_at"`
	UbbPairName                string             `json:"ubb_pair_name"`
	Architecture               string             `json:"architecture"`
	XAesGcm                    FlexBool           `json:"x_aes_gcm"`
	HasFan                     FlexBool           `json:"has_fan"`
	HasEth1                    FlexBool           `json:"has_eth1"`
	ModelIncompatible          FlexBool           `json:"model_incompatible"`
	XAuthkey                   string             `json:"x_authkey"`
	ModelInEol                 FlexBool           `json:"model_in_eol"`
	HasTemperature             FlexBool           `json:"has_temperature"`
	AdoptedByClient            string             `json:"adopted_by_client"`
	ModelInLts                 FlexBool           `json:"model_in_lts"`
	UbbPairID                  string             `json:"ubb_pair_id"`
	KernelVersion              string             `json:"kernel_version"`
	Serial                     string             `json:"serial"`
	UbbPsk                     string             `json:"ubb_psk"`
	FixedApAvailable           FlexBool           `json:"fixed_ap_available"`
	LedOverrideColorBrightness FlexInt            `json:"led_override_color_brightness"`
	Adopted                    FlexBool           `json:"adopted"`
	HashID                     string             `json:"hash_id"`
	DeviceID                   string             `json:"device_id"`
	Uplink                     *Uplink            `json:"uplink"`
	State                      FlexInt            `json:"state"`
	StartDisconnectedMillis    FlexInt            `json:"start_disconnected_millis"`
	UpgradeState               FlexInt            `json:"upgrade_state"`
	StartConnectedMillis       FlexInt            `json:"start_connected_millis"`
	LastSeen                   FlexInt            `json:"last_seen"`
	NextInterval               FlexInt            `json:"next_interval"`
	DisconnectionReason        string             `json:"disconnection_reason"`
	Upgradable                 FlexBool           `json:"upgradable"`
	AdoptableWhenUpgraded      FlexBool           `json:"adoptable_when_upgraded"`
	UpgradeToFirmware          string             `json:"upgrade_to_firmware"`
	KnownCfgversion            string             `json:"known_cfgversion"`
	Uptime                     FlexInt            `json:"uptime"`
	Uptime0                    FlexInt            `json:"_uptime"`
	Locating                   FlexBool           `json:"locating"`
	ConnectRequestIP           string             `json:"connect_request_ip"`
	ConnectRequestPort         string             `json:"connect_request_port"`
	SysStats                   *SysStats          `json:"sys_stats"`
	SystemStats                *SystemStats       `json:"system-stats"`
	SSHSessionTable            []any              `json:"ssh_session_table"`
	LldpTable                  []*LLDPTable       `json:"lldp_table"`
	DisplayableVersion         string             `json:"displayable_version"`
	ConnectionNetworkID        string             `json:"connection_network_id"`
	ConnectionNetworkName      string             `json:"connection_network_name"`
	StartupTimestamp           FlexInt            `json:"startup_timestamp"`
	IsAccessPoint              FlexBool           `json:"is_access_point"`
	SafeForAutoupgrade         FlexBool           `json:"safe_for_autoupgrade"`
	UbbPeerState               FlexInt            `json:"ubb_peer_state"`
	LinkQuality                FlexInt            `json:"link_quality"`
	LinkQualityCurrent         FlexInt            `json:"link_quality_current"`
	GeneralTemperature         FlexInt            `json:"general_temperature"`
	LinkCapacity               FlexInt            `json:"link_capacity"`
	Mode                       string             `json:"mode"`
	P2PStats                   *P2PStats          `json:"p2p_stats"`
	Isolated                   FlexBool           `json:"isolated"`
	RadioTableStats            []*RadioTableStats `json:"radio_table_stats"`
	PortStats                  []any              `json:"port_stats"`
	VwireTable                 []any              `json:"vwire_table"`
	VwireVapTable              []any              `json:"vwire_vap_table"`
	DownlinkLLDPMacs           []*LLDPTable       `json:"downlink_lldp_macs"`
	VapTable                   []*VapTable        `json:"vap_table"`
	DownlinkTable              []any              `json:"downlink_table"`
	BytesD                     FlexInt            `json:"bytes-d"`
	TxBytesD                   FlexInt            `json:"tx_bytes-d"`
	RxBytesD                   FlexInt            `json:"rx_bytes-d"`
	BytesR                     FlexInt            `json:"bytes-r"`
	PeerUbb                    *UBB               `json:"peer_ubb"`
	Stat                       *UBBStat           `json:"stat"`
	TxBytes                    FlexInt            `json:"tx_bytes"`
	RxBytes                    FlexInt            `json:"rx_bytes"`
	Bytes                      FlexInt            `json:"bytes"`
	NumSta                     FlexInt            `json:"num_sta"`
	UserNumSta                 FlexInt            `json:"user-num_sta"`
	UserWlanNumSta             FlexInt            `json:"user-wlan-num_sta"`
	XHasSSHHostkey             FlexBool           `json:"x_has_ssh_hostkey"`
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
