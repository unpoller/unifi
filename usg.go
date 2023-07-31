package unifi

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// USG represents all the data from the Ubiquiti Controller for a Unifi Security Gateway.
type USG struct {
	AdoptableWhenUpgraded FlexBool             `json:"adoptable_when_upgraded"`
	Adopted               FlexBool             `fake:"{constFlexBool:true}"    json:"adopted"`
	BoardRev              FlexInt              `json:"board_rev"`
	Bytes                 FlexInt              `json:"bytes"`
	Cfgversion            string               `fake:"{appversion}"            json:"cfgversion"`
	ConfigNetwork         *ConfigNetwork       `json:"config_network"`
	ConnectRequestIP      string               `fake:"{ipv4address}"           json:"connect_request_ip"`
	ConnectRequestPort    string               `json:"connect_request_port"`
	DeviceID              string               `fake:"{uuid}"                  json:"device_id"`
	EthernetOverrides     []*EthernetOverrides `fakesize:"5"                   json:"ethernet_overrides"`
	EthernetTable         []*EthernetTable     `fakesize:"5"                   json:"ethernet_table"`
	FwCaps                FlexInt              `json:"fw_caps"`
	GuestNumSta           FlexInt              `json:"guest-num_sta"`
	GuestToken            string               `json:"guest_token"`
	HwCaps                FlexInt              `json:"hw_caps"`
	ID                    string               `fake:"{uuid}"                  json:"_id"`
	InformIP              string               `fake:"{ipv4address}"           json:"inform_ip"`
	InformURL             string               `fake:"{url}"                   json:"inform_url"`
	IP                    string               `fake:"{ipv4address}"           json:"ip"`
	KnownCfgversion       string               `fake:"{appversion}"            json:"known_cfgversion"`
	LastSeen              FlexInt              `json:"last_seen"`
	LedOverride           string               `json:"led_override"`
	LicenseState          string               `json:"license_state"`
	Locating              FlexBool             `fake:"{constFlexBool:false}"   json:"locating"`
	Mac                   string               `fake:"{macaddress}"            json:"mac"`
	Model                 string               `json:"model"`
	Name                  string               `fake:"{animal}"                json:"name"`
	NetworkTable          NetworkTable         `json:"network_table"`
	NumDesktop            FlexInt              `json:"num_desktop"`
	NumHandheld           FlexInt              `json:"num_handheld"`
	NumMobile             FlexInt              `json:"num_mobile"`
	NumSta                FlexInt              `json:"num_sta"`
	OutdoorModeOverride   string               `json:"outdoor_mode_override"`
	PortTable             []*Port              `fakesize:"5"                   json:"port_table"`
	RequiredVersion       string               `fake:"{appversion}"            json:"required_version"`
	Rollupgrade           FlexBool             `json:"rollupgrade"`
	RxBytes               FlexInt              `json:"rx_bytes"`
	Serial                string               `json:"serial"`
	site                  *Site
	SiteID                string               `fake:"{uuid}"                 json:"site_id"`
	SiteName              string               `json:"-"`
	SourceName            string               `json:"-"`
	SpeedtestStatus       SpeedtestStatus      `json:"speedtest-status"`
	SpeedtestStatusSaved  FlexBool             `json:"speedtest-status-saved"`
	Stat                  USGStat              `json:"stat"`
	State                 FlexInt              `json:"state"`
	SysStats              SysStats             `json:"sys_stats"`
	SystemStats           SystemStats          `json:"system-stats"`
	Temperatures          []Temperature        `fakesize:"5"                  json:"temperatures,omitempty"`
	TxBytes               FlexInt              `json:"tx_bytes"`
	Type                  string               `fake:"{lexify:usg}"           json:"type"`
	Unsupported           FlexBool             `json:"unsupported"`
	UnsupportedReason     FlexInt              `json:"unsupported_reason"`
	Upgradable            FlexBool             `json:"upgradable"`
	Uplink                Uplink               `json:"uplink"`
	Uptime                FlexInt              `json:"uptime"`
	UserNumSta            FlexInt              `json:"user-num_sta"`
	UsgCaps               FlexInt              `json:"usg_caps"`
	Version               string               `fake:"{appversion}"           json:"version"`
	Wan1                  Wan                  `json:"wan1"`
	Wan2                  Wan                  `json:"wan2"`
}

// Uplink is the Internet connection (or uplink) on a UniFi device.
type Uplink struct {
	BytesR           FlexInt  `json:"bytes-r"`
	Drops            FlexInt  `json:"drops"`
	Enable           FlexBool `json:"enable,omitempty"`
	FullDuplex       FlexBool `json:"full_duplex"`
	Gateways         []string `fakesize:"5"                       json:"gateways,omitempty"`
	IP               string   `fake:"{ipv4address}"               json:"ip"`
	Latency          FlexInt  `json:"latency"`
	Mac              string   `json:"mac,omitempty"`
	MaxSpeed         FlexInt  `json:"max_speed"`
	Media            string   `json:"media"`
	Name             string   `fake:"{animal}"                    json:"name"`
	Nameservers      []string `json:"nameservers"`
	Netmask          string   `json:"netmask"`
	NumPort          FlexInt  `json:"num_port"`
	PortIdx          FlexInt  `json:"port_idx"`
	RxBytes          FlexInt  `json:"rx_bytes"`
	RxBytesR         FlexInt  `json:"rx_bytes-r"`
	RxDropped        FlexInt  `json:"rx_dropped"`
	RxErrors         FlexInt  `json:"rx_errors"`
	RxMulticast      FlexInt  `json:"rx_multicast"`
	RxPackets        FlexInt  `json:"rx_packets"`
	RxRate           FlexInt  `json:"rx_rate"`
	Speed            FlexInt  `json:"speed"`
	SpeedtestLastrun FlexInt  `json:"speedtest_lastrun,omitempty"`
	SpeedtestPing    FlexInt  `json:"speedtest_ping,omitempty"`
	SpeedtestStatus  string   `json:"speedtest_status,omitempty"`
	TxBytes          FlexInt  `json:"tx_bytes"`
	TxBytesR         FlexInt  `json:"tx_bytes-r"`
	TxDropped        FlexInt  `json:"tx_dropped"`
	TxErrors         FlexInt  `json:"tx_errors"`
	TxPackets        FlexInt  `json:"tx_packets"`
	TxRate           FlexInt  `json:"tx_rate"`
	Type             string   `json:"type"`
	Up               FlexBool `json:"up"`
	Uptime           FlexInt  `json:"uptime"`
	XputDown         FlexInt  `json:"xput_down,omitempty"`
	XputUp           FlexInt  `json:"xput_up,omitempty"`
}

// Wan is a Wan interface on a USG or UDM.
type Wan struct {
	Autoneg     FlexBool `json:"autoneg"`
	BytesR      FlexInt  `json:"bytes-r"`
	DNS         []string `fakesize:"5"                json:"dns"` // may be deprecated
	Enable      FlexBool `json:"enable"`
	FlowctrlRx  FlexBool `json:"flowctrl_rx"`
	FlowctrlTx  FlexBool `json:"flowctrl_tx"`
	FullDuplex  FlexBool `json:"full_duplex"`
	Gateway     string   `json:"gateway"` // may be deprecated
	IP          string   `fake:"{ipv4address}"        json:"ip"`
	Ifname      string   `json:"ifname"`
	IsUplink    FlexBool `json:"is_uplink"`
	Mac         string   `fake:"{macaddress}"         json:"mac"`
	MaxSpeed    FlexInt  `json:"max_speed"`
	Media       string   `json:"media"`
	Name        string   `fake:"{animal}"             json:"name"`
	Netmask     string   `json:"netmask"` // may be deprecated
	NumPort     int      `json:"num_port"`
	PortIdx     int      `json:"port_idx"`
	PortPoe     FlexBool `json:"port_poe"`
	RxBroadcast FlexInt  `json:"rx_broadcast"`
	RxBytes     FlexInt  `json:"rx_bytes"`
	RxBytesR    FlexInt  `json:"rx_bytes-r"`
	RxDropped   FlexInt  `json:"rx_dropped"`
	RxErrors    FlexInt  `json:"rx_errors"`
	RxMulticast FlexInt  `json:"rx_multicast"`
	RxPackets   FlexInt  `json:"rx_packets"`
	RxRate      FlexInt  `json:"rx_rate"`
	Speed       FlexInt  `json:"speed"`
	SpeedCaps   FlexInt  `json:"speed_caps"`
	TxBroadcast FlexInt  `json:"tx_broadcast"`
	TxBytes     FlexInt  `json:"tx_bytes"`
	TxBytesR    FlexInt  `json:"tx_bytes-r"`
	TxDropped   FlexInt  `json:"tx_dropped"`
	TxErrors    FlexInt  `json:"tx_errors"`
	TxMulticast FlexInt  `json:"tx_multicast"`
	TxPackets   FlexInt  `json:"tx_packets"`
	TxRate      FlexInt  `json:"tx_rate"`
	Type        string   `json:"type"`
	Up          FlexBool `fake:"{constFlexBool:true}" json:"up"`
}

// SpeedtestStatus is the speed test info on a USG or UDM.
type SpeedtestStatus struct {
	Latency         FlexInt          `json:"latency"`
	Rundate         FlexInt          `json:"rundate"`
	Runtime         FlexInt          `json:"runtime"`
	Server          *SpeedtestServer `json:"server"`
	ServerDesc      string           `json:"server_desc,omitempty"`
	SourceInterface string           `json:"source_interface"`
	StatusDownload  FlexInt          `json:"status_download"`
	StatusPing      FlexInt          `json:"status_ping"`
	StatusSummary   FlexInt          `json:"status_summary"`
	StatusUpload    FlexInt          `json:"status_upload"`
	XputDownload    FlexInt          `json:"xput_download"`
	XputUpload      FlexInt          `json:"xput_upload"`
}

type SpeedtestServer struct {
	Cc          string  `json:"cc"`
	City        string  `fake:"{city}"    json:"city"`
	Country     string  `fake:"{country}" json:"country"`
	Lat         FlexInt `json:"lat"`
	Lon         FlexInt `json:"lon"`
	Provider    string  `json:"provider"`
	ProviderURL string  `fake:"{url}"     json:"provider_url"`
}

type TempStatusByName map[string]*FlexTemp

func (t TempStatusByName) Fake(faker *gofakeit.Faker) interface{} {
	return TempStatusByName{
		"cpu":     NewFlexTemp(float64(faker.Rand.Int63n(100))),
		"sys":     NewFlexTemp(float64(faker.Rand.Int63n(100))),
		"probe":   NewFlexTemp(float64(faker.Rand.Int63n(100))),
		"memory":  NewFlexTemp(float64(faker.Rand.Int63n(100))),
		"network": NewFlexTemp(float64(faker.Rand.Int63n(100))),
	}
}

// SystemStats is system info for a UDM, USG, USW.
type SystemStats struct {
	CPU    FlexInt `json:"cpu"`
	Mem    FlexInt `json:"mem"`
	Uptime FlexInt `json:"uptime"`
	// This exists on at least USG4, may others, maybe not.
	// {"Board (CPU)":"51 C","Board (PHY)":"51 C","CPU":"72 C","PHY":"77 C"}
	Temps TempStatusByName `json:"temps,omitempty"`
}

// SysStats is load info for a UDM, USG, USW.
type SysStats struct {
	Loadavg1  FlexInt `json:"loadavg_1"`
	Loadavg15 FlexInt `json:"loadavg_15"`
	Loadavg5  FlexInt `json:"loadavg_5"`
	MemBuffer FlexInt `json:"mem_buffer"`
	MemTotal  FlexInt `json:"mem_total"`
	MemUsed   FlexInt `json:"mem_used"`
}

// USGStat holds the "stat" data for a gateway.
// This is split out because of a JSON data format change from 5.10 to 5.11.
type USGStat struct {
	*Gw
}

// Gw is a subtype of USGStat to make unmarshalling of different controller versions possible.
type Gw struct {
	Datetime     time.Time `fake:"{recent_time}"           json:"datetime"`
	Duration     FlexInt   `json:"duration"`
	Gw           string    `json:"gw"`
	LanRxBytes   FlexInt   `json:"lan-rx_bytes"`
	LanRxDropped FlexInt   `json:"lan-rx_dropped"`
	LanRxErrors  FlexInt   `json:"lan-rx_errors,omitempty"`
	LanRxPackets FlexInt   `json:"lan-rx_packets"`
	LanTxBytes   FlexInt   `json:"lan-tx_bytes"`
	LanTxPackets FlexInt   `json:"lan-tx_packets"`
	O            string    `json:"o"`
	Oid          string    `json:"oid"`
	SiteID       string    `json:"site_id"`
	Time         FlexInt   `json:"time"`
	WanRxBytes   FlexInt   `json:"wan-rx_bytes"`
	WanRxDropped FlexInt   `json:"wan-rx_dropped"`
	WanRxErrors  FlexInt   `json:"wan-rx_errors,omitempty"`
	WanRxPackets FlexInt   `json:"wan-rx_packets"`
	WanTxBytes   FlexInt   `json:"wan-tx_bytes"`
	WanTxPackets FlexInt   `json:"wan-tx_packets"`
}

// UnmarshalJSON unmarshalls 5.10 or 5.11 formatted Gateway Stat data.
func (v *USGStat) UnmarshalJSON(data []byte) error {
	var n struct {
		Gw `json:"gw"`
	}

	v.Gw = &n.Gw

	err := json.Unmarshal(data, v.Gw) // controller version 5.10.
	if err != nil {
		return json.Unmarshal(data, &n) // controller version 5.11.
	}

	return nil
}
