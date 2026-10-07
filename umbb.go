package unifi

// UMBB represents a UniFi Mobile Broadband device, such as the U5G Max (UMBBE630).
// These are cellular (4G/5G) WAN modems adopted as their own device, with type "umbb".
// The modem, radio and SIM state live in the nested "mbb" object.
type UMBB struct {
	site                    *Site
	AdoptableWhenUpgraded   FlexBool         `json:"adoptable_when_upgraded,omitempty"`
	Adopted                 FlexBool         `fake:"{constFlexBool:true}"              json:"adopted"`
	AdoptionCompleted       FlexBool         `json:"adoption_completed"`
	Architecture            string           `json:"architecture"`
	BoardRev                FlexInt          `json:"board_rev"`
	CanActivateEsim         FlexBool         `json:"can_activate_esim"`
	Cfgversion              string           `json:"cfgversion"`
	ConfigNetwork           *ConfigNetwork   `json:"config_network"`
	ConnectRequestIP        string           `json:"connect_request_ip"`
	ConnectRequestPort      string           `json:"connect_request_port"`
	ConnectedAt             FlexInt          `json:"connected_at"`
	ConnectionNetworkName   string           `json:"connection_network_name"`
	DeviceID                string           `json:"device_id"`
	DisconnectedAt          FlexInt          `json:"disconnected_at"`
	DisplayableVersion      string           `json:"displayable_version"`
	EthernetTable           []*EthernetTable `json:"ethernet_table"`
	FwCaps                  FlexInt          `json:"fw_caps"`
	GatewayMac              string           `json:"gateway_mac"`
	HasFan                  FlexBool         `json:"has_fan"`
	HasTemperature          FlexBool         `json:"has_temperature"`
	HwCaps                  FlexInt          `json:"hw_caps"`
	ID                      string           `json:"_id"`
	IP                      string           `fake:"{ipv4address}"                     json:"ip"`
	InformIP                string           `json:"inform_ip"`
	InformURL               string           `json:"inform_url"`
	Internet                FlexBool         `json:"internet"`
	IsAccessPoint           FlexBool         `json:"is_access_point"`
	KernelVersion           string           `json:"kernel_version"`
	KnownCfgversion         string           `json:"known_cfgversion"`
	LastConfigAppliedOK     FlexBool         `json:"last_config_applied_successfully"`
	LastSeen                FlexInt          `json:"last_seen"`
	Locating                FlexBool         `fake:"{constFlexBool:false}"             json:"locating"`
	Mac                     string           `json:"mac"`
	ManufacturerID          FlexInt          `json:"manufacturer_id"`
	Mbb                     MBB              `json:"mbb"`
	MbbOverrides            MBBOverrides     `json:"mbb_overrides"`
	Model                   string           `json:"model"`
	ModelInEOL              FlexBool         `json:"model_in_eol"`
	ModelInLTS              FlexBool         `json:"model_in_lts"`
	ModelIncompatible       FlexBool         `json:"model_incompatible"`
	Name                    string           `fake:"{animal}"                          json:"name"`
	NextInterval            FlexInt          `json:"next_interval"`
	PreviousFirmwareVersion string           `json:"previous_firmware_version"`
	ProvisionedAt           FlexInt          `json:"provisioned_at"`
	RequiredVersion         string           `json:"required_version"`
	Rollupgrade             FlexBool         `json:"rollupgrade,omitempty"`
	Serial                  string           `json:"serial"`
	SetupID                 string           `json:"setup_id"`
	Shortname               string           `json:"shortname"`
	SiteID                  string           `json:"site_id"`
	SiteName                string           `json:"-"`
	SourceName              string           `json:"-"`
	StartConnectedMillis    FlexInt          `json:"start_connected_millis"`
	StartDisconnectedMillis FlexInt          `json:"start_disconnected_millis"`
	StartupTimestamp        FlexInt          `json:"startup_timestamp"`
	State                   FlexInt          `json:"state"`
	SysErrorCaps            FlexInt          `json:"sys_error_caps"`
	SysStats                SysStats         `json:"sys_stats"`
	SystemStats             SystemStats      `json:"system-stats"`
	Tags                    []string         `json:"tags"` // Device tags assigned to this device
	TwoPhaseAdopt           FlexBool         `json:"two_phase_adopt"`
	Type                    string           `fake:"{randomstring:[umbb]}"             json:"type"`
	Unsupported             FlexBool         `json:"unsupported"`
	UnsupportedReason       FlexInt          `json:"unsupported_reason"`
	Upgradable              FlexBool         `json:"upgradable,omitempty"`
	Uplink                  Uplink           `json:"uplink"`
	Uptime                  FlexInt          `json:"uptime"`
	Version                 string           `json:"version"`
}

// MBBOverrides holds the operator-configured mobile broadband settings.
type MBBOverrides struct {
	PrimarySlot FlexInt          `json:"primary_slot"`
	Sim         []MBBSimOverride `json:"sim"`
}

// MBBSimOverride is the data-plan configuration for one SIM slot.
// Live usage against this plan is reported on MBBSim.
type MBBSimOverride struct {
	DataLimitEnabled         FlexBool `json:"data_limit_enabled"`
	DataSoftLimitBytes       FlexInt  `json:"data_soft_limit_bytes"`
	DataSoftLimitDisplayUnit string   `json:"data_soft_limit_display_unit"`
	DataWarningThreshold     FlexInt  `json:"data_warning_threshold"`
	ResetDate                FlexInt  `json:"reset_date"`
	ResetPolicy              string   `json:"reset_policy"`
	Slot                     FlexInt  `json:"slot"`
}

// MBB is the mobile broadband modem state reported by a UMBB device.
type MBB struct {
	Esim         MBBEsim         `json:"esim"`
	GeoInfo      MBBGeoInfo      `json:"geo_info"`
	IMEI         string          `json:"imei"`
	IPSettings   MBBIPSettings   `json:"ip_settings"`
	IPv6Settings MBBIPv6Settings `json:"ipv6_settings"`
	Mode         string          `json:"mode"` // e.g. "failover"
	NetworkScan  MBBNetworkScan  `json:"network_scan"`
	Radio        MBBRadio        `json:"radio"`
	Sim          []MBBSim        `fakesize:"2"         json:"sim"`
	State        string          `json:"state"` // e.g. "ready"
}

// MBBRadio is the live cellular radio state: technology, band, signal and carrier aggregation.
// When attached over 5G, the plain and _nr signal values carry the same reading.
// Snr and SnrNr are fractional dB (for example 17.6), so they use FlexFloat.
type MBBRadio struct {
	SA5GMode          FlexBool        `json:"5g_sa_mode"`
	Band              string          `json:"band"` // e.g. "n78" or "b3"
	CaLte             []MBBCarrierLte `json:"ca_lte"`
	CaNr              []MBBCarrierNr  `json:"ca_nr"`
	CellID            FlexInt         `json:"cell_id"`
	Channel           FlexInt         `json:"channel"`
	CurrentSlot       FlexInt         `json:"current_slot"`
	HasCoverage       FlexBool        `json:"has_coverage"`
	HplmnDenied       FlexBool        `json:"hplmn_denied"`
	HplmnServing      FlexBool        `json:"hplmn_serving"`
	IsVerizon         FlexBool        `json:"is_verizon"`
	LteBands          []FlexInt       `json:"lte_bands"`
	MaxBitrateDl      FlexInt         `json:"max_bitrate_dl"`
	MaxBitrateDlNr    FlexInt         `json:"max_bitrate_dl_nr"`
	MaxBitrateUl      FlexInt         `json:"max_bitrate_ul"`
	MaxBitrateUlNr    FlexInt         `json:"max_bitrate_ul_nr"`
	Mcc               FlexInt         `json:"mcc"`
	MccCountry        string          `json:"mcc_cc2"`
	Mnc               FlexInt         `json:"mnc"`
	NetworkOperator   string          `json:"networkoperator"`
	Nr5GBands         []FlexInt       `json:"nr5g_bands"`
	Pci               FlexInt         `json:"pci"`
	Rat               string          `json:"rat"` // e.g. "5G" or "LTE"
	Rat5GUW           FlexBool        `json:"rat_5g_uw"`
	RatCaps           []string        `json:"rat_caps"`
	RatModeActive     string          `json:"rat_mode_active"`
	RegistrationState FlexInt         `json:"registration_state"`
	Roaming           FlexBool        `json:"roaming"`
	Rsrp              FlexInt         `json:"rsrp"`
	RsrpNr            FlexInt         `json:"rsrp_nr"`
	Rsrq              FlexInt         `json:"rsrq"`
	RsrqNr            FlexInt         `json:"rsrq_nr"`
	Signal            FlexInt         `json:"signal"` // 0-5 bars
	SignalPercent     FlexInt         `json:"signal_percent"`
	Snr               FlexFloat       `json:"snr"`
	SnrNr             FlexFloat       `json:"snr_nr"`
}

// MBBCarrierNr is one 5G NR component carrier. UlBwMhz is null on downlink-only carriers.
type MBBCarrierNr struct {
	Band    FlexInt  `json:"band"`
	DlArfcn FlexInt  `json:"dl_arfcn"`
	DlBwMhz FlexInt  `json:"dl_bw_mhz"`
	Primary FlexBool `json:"primary"`
	UlArfcn FlexInt  `json:"ul_arfcn"`
	UlBwMhz FlexInt  `json:"ul_bw_mhz"`
}

// MBBCarrierLte is one LTE component carrier.
// LTE reports the channel as dl_earfcn / ul_earfcn. NR uses dl_arfcn / ul_arfcn
// on MBBCarrierNr, so this is not the same struct. A missing uplink bandwidth is null.
type MBBCarrierLte struct {
	Band     FlexInt  `json:"band"`
	DlEarfcn FlexInt  `json:"dl_earfcn"`
	DlBwMhz  FlexInt  `json:"dl_bw_mhz"`
	Primary  FlexBool `json:"primary"`
	UlEarfcn FlexInt  `json:"ul_earfcn"`
	UlBwMhz  FlexInt  `json:"ul_bw_mhz"`
}

// MBBSim is the state of one SIM slot (physical or eSIM). RxBytes and TxBytes are
// cumulative counters, sent by the controller as strings.
// Network reject arrives as NetworkReject (a single string such as "ip") on some
// firmware, and as NetworkRejectType, NetworkRejectText, and NetworkRejectAge on
// others. The age unit is not documented; one device held the same value for many
// hours, so it is not a running age in seconds.
type MBBSim struct {
	Active              FlexBool          `json:"active"`
	Asn                 FlexInt           `json:"asn"`
	CardPresent         FlexBool          `json:"card_present"`
	ConnectionInfo      MBBConnectionInfo `json:"connection_info"`
	CurrentApn          MBBApn            `json:"current_apn"`
	DataLimited         FlexBool          `json:"data_limited"`
	DataWarning         FlexBool          `json:"data_warning"`
	DisplayState        string            `json:"display_state"` // e.g. "operational", "no-sim"
	DisplayStateElapsed FlexInt           `json:"display_state_elapsed"`
	Esim                FlexBool          `json:"esim"`
	HasCarrier          FlexBool          `json:"has_carrier"`
	HasDataPlan         FlexBool          `json:"has_data_plan"`
	ICCID               string            `json:"iccid"`
	Incompatible        FlexBool          `json:"incompatible"`
	Mcc                 FlexInt           `json:"mcc"`
	Metered             FlexBool          `json:"metered"`
	Mnc                 FlexInt           `json:"mnc"`
	NetworkReject       string            `json:"network_reject"`
	NetworkRejectAge    FlexInt           `json:"network_reject_age"`
	NetworkRejectText   string            `json:"network_reject_text"`
	NetworkRejectType   string            `json:"network_reject_type"`
	OperationInProgress FlexBool          `json:"operation_in_progress"`
	PinBlocked          FlexBool          `json:"pin_blocked"`
	PinLock             FlexBool          `json:"pin_lock"`
	RxBytes             FlexInt           `json:"rxbytes"`
	Slot                FlexInt           `json:"slot"`
	Spn                 string            `json:"spn"` // Service provider name, e.g. "EE"
	TxBytes             FlexInt           `json:"txbytes"`
}

// MBBConnectionInfo describes the data session on a SIM.
type MBBConnectionInfo struct {
	Asn       FlexInt `json:"asn"`
	InetState FlexInt `json:"inet_state"`
	IPType    string  `json:"ip_type"`
	Timestamp FlexInt `json:"timestamp"`
}

// MBBApn is the APN in use on a SIM. Credentials are deliberately not decoded.
type MBBApn struct {
	Apn         string   `json:"apn"`
	AuthType    string   `json:"auth_type"`
	ClatEnabled FlexBool `json:"clat_enabled"`
	PdpType     string   `json:"pdp_type"`
	Roaming     FlexBool `json:"roaming"`
}

// MBBEsim describes the embedded SIM.
type MBBEsim struct {
	EID   string `json:"eid"`
	ICCID string `json:"iccid"`
}

// MBBGeoInfo is the public address and ISP the modem is seen from.
type MBBGeoInfo struct {
	Address string `fake:"{ipv4address}" json:"address"`
	ISP     string `json:"isp"`
}

// MBBIPSettings is the IPv4 configuration of the cellular data session.
type MBBIPSettings struct {
	IPv4Address string  `fake:"{ipv4address}" json:"ipv4_address"`
	IPv4DNS     string  `json:"ipv4_dns"`
	IPv4Gateway string  `json:"ipv4_gateway"`
	IPv4Netmask string  `json:"ipv4_netmask"`
	MTU         FlexInt `json:"mtu"`
}

// MBBIPv6Settings is the IPv6 configuration of the cellular data session.
type MBBIPv6Settings struct {
	IPv6Address string  `json:"ipv6_address"`
	IPv6DNS     string  `json:"ipv6_dns"`
	IPv6DNS2    string  `json:"ipv6_dns2"`
	IPv6Gateway string  `json:"ipv6_gateway"`
	MTU         FlexInt `json:"mtu"`
}

// MBBNetworkScan reports the last manual operator scan.
type MBBNetworkScan struct {
	ScanStatus struct {
		LastScanTimestamp FlexInt  `json:"last_scan_timestamp"`
		OperatorCount     FlexInt  `json:"operator_count"`
		ScanCompleted     FlexBool `json:"scan_completed"`
		ScanInProgress    FlexBool `json:"scan_in_progress"`
		ScanSuccess       FlexBool `json:"scan_success"`
	} `json:"scan_status"`
}
