package unifi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrCannotUnmarshalFlexInt = fmt.Errorf("cannot unmarshal to FlexInt")

// This is a list of unifi API paths.
// The %s in each string must be replaced with a Site.Name.
const (
	// APIRogueAP shows your neighbors' wifis.
	APIRogueAP string = "/api/s/%s/stat/rogueap"
	// APIStatusPath shows Controller version.
	APIStatusPath string = "/status"
	// APIEventPath contains UniFi Event data.
	APIEventPath string = "/api/s/%s/stat/event"
	// APISiteList is the path to the api site list.
	APISiteList string = "/api/stat/sites"
	// APISiteDPI is site DPI data.
	APISiteDPI string = "/api/s/%s/stat/sitedpi"
	// APISiteDPI is site DPI data.
	APIClientDPI string = "/api/s/%s/stat/stadpi"
	// APIClientPath is Unifi Clients API Path.
	APIClientPath string = "/api/s/%s/stat/sta"
	// APIAllUserPath is Unifi Insight all previous Clients API Path.
	APIAllUserPath string = "/api/s/%s/stat/alluser"
	// APINetworkPath is where we get data about Unifi networks.
	APINetworkPath string = "/api/s/%s/rest/networkconf"
	// APIDevicePath is where we get data about Unifi devices.
	APIDevicePath string = "/api/s/%s/stat/device"
	// APILoginPath is Unifi Controller Login API Path.
	APILoginPath string = "/api/login"
	// APILoginPathNew is how we log into UDM 5.12.55+.
	APILoginPathNew string = "/api/auth/login"
	// APILogoutPath is how we logout from UDM.
	APILogoutPath string = "/api/logout"
	// APIEventPathIDS returns Intrusion Detection/Prevention Systems Events.
	APIEventPathIDS string = "/api/s/%s/stat/ips/event"
	// APIEventPathAlarms contains the site alarms.
	APIEventPathAlarms string = "/api/s/%s/list/alarm"
	// APIPrefixNew is the prefix added to the new API paths; except login. duh.
	APIPrefixNew string = "/proxy/network"
	// APIAnomaliesPath returns site anomalies.
	APIAnomaliesPath string = "/api/s/%s/stat/anomalies"
	APICommandPath   string = "/api/s/%s/cmd"
	APIDevMgrPath    string = APICommandPath + "/devmgr"
)

// path returns the correct api path based on the new variable.
// new is based on the unifi-controller output. is it new or old output?
func (u *Unifi) path(path string) string {
	if u.new {
		if path == APILoginPath {
			return APILoginPathNew
		}

		if !strings.HasPrefix(path, APIPrefixNew) && path != APILoginPathNew {
			return APIPrefixNew + path
		}
	}

	return path
}

// Logger is a base type to deal with changing log outputs. Create a logger
// that matches this interface to capture debug and error logs.
type Logger func(msg string, fmt ...interface{})

// discardLogs is the default debug logger.
func discardLogs(msg string, v ...interface{}) {
	// do nothing.
}

// Devices contains a list of all the unifi devices from a controller.
// Contains Access points, security gateways and switches.
type Devices struct {
	UAPs []*UAP
	USGs []*USG
	USWs []*USW
	UDMs []*UDM
	UXGs []*UXG
	PDUs []*PDU
}

// Config is the data passed into our library. This configures things and allows
// us to connect to a controller and write log messages. Optional SSLCert is used
// for ssl cert pinning; provide the content of a PEM to validate the server's cert.
type Config struct {
	User      string
	Pass      string
	URL       string
	SSLCert   [][]byte
	ErrorLog  Logger
	DebugLog  Logger
	Timeout   time.Duration // how long to wait for replies, default: forever.
	VerifySSL bool
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	*Config
	*server
	csrf         string
	fingerprints fingerprints
	new          bool
}

type fingerprints []string

// Contains returns true if the fingerprint is in the list.
func (f fingerprints) Contains(s string) bool {
	for i := range f {
		if s == f[i] {
			return true
		}
	}

	return false
}

// server is the /status endpoint from the Unifi controller.
type server struct {
	Up            FlexBool `json:"up"`
	ServerVersion string   `json:"server_version"`
	UUID          string   `json:"uuid"`
}

// FlexInt provides a container and unmarshalling for fields that may be
// numbers or strings in the Unifi API.
type FlexInt struct {
	Val float64
	Txt string
}

func NewFlexInt(v float64) *FlexInt {
	return &FlexInt{
		Val: v,
		Txt: strconv.FormatFloat(v, 'f', -1, 64),
	}
}

// UnmarshalJSON converts a string or number to an integer.
// Generally, do call this directly, it's used in the json interface.
func (f *FlexInt) UnmarshalJSON(b []byte) error {
	var unk interface{}

	if err := json.Unmarshal(b, &unk); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	switch i := unk.(type) {
	case float64:
		f.Val = i
		f.Txt = strconv.FormatFloat(i, 'f', -1, 64)
	case string:
		f.Txt = i
		f.Val, _ = strconv.ParseFloat(i, 64)
	case nil:
		f.Txt = "0"
		f.Val = 0
	default:
		return fmt.Errorf("%v: %w", b, ErrCannotUnmarshalFlexInt)
	}

	return nil
}

func (f *FlexInt) Int() int {
	return int(f.Val)
}

func (f *FlexInt) Int64() int64 {
	return int64(f.Val)
}

func (f *FlexInt) String() string {
	return f.Txt
}

func (f *FlexInt) Add(o *FlexInt) {
	f.Val += o.Val
	f.Txt = strconv.FormatFloat(f.Val, 'f', -1, 64)
}

func (f *FlexInt) AddFloat64(v float64) {
	f.Val += v
	f.Txt = strconv.FormatFloat(f.Val, 'f', -1, 64)
}

// FlexBool provides a container and unmarshalling for fields that may be
// boolean or strings in the Unifi API.
type FlexBool struct {
	Val bool
	Txt string
}

// UnmarshalJSON method converts armed/disarmed, yes/no, active/inactive or 0/1 to true/false.
// Really it converts ready, ok, up, t, armed, yes, active, enabled, 1, true to true. Anything else is false.
func (f *FlexBool) UnmarshalJSON(b []byte) error {
	f.Txt = strings.Trim(string(b), `"`)
	f.Val = f.Txt == "1" || strings.EqualFold(f.Txt, "true") || strings.EqualFold(f.Txt, "yes") ||
		strings.EqualFold(f.Txt, "t") || strings.EqualFold(f.Txt, "armed") || strings.EqualFold(f.Txt, "active") ||
		strings.EqualFold(f.Txt, "enabled") || strings.EqualFold(f.Txt, "ready") || strings.EqualFold(f.Txt, "up") ||
		strings.EqualFold(f.Txt, "ok")

	return nil
}

func (f *FlexBool) String() string {
	return f.Txt
}

// DownlinkTable is part of a UXG and UDM output.
type DownlinkTable struct {
	PortIdx    FlexInt  `json:"port_idx"`
	Speed      FlexInt  `json:"speed"`
	FullDuplex FlexBool `json:"full_duplex"`
	Mac        string   `json:"mac"`
}

// ConfigNetwork comes from gateways.
type ConfigNetwork struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type EthernetTable struct {
	Mac     string  `json:"mac"`
	NumPort FlexInt `json:"num_port"`
	Name    string  `json:"name"`
}

// Port is a physical connection on a USW or Gateway.
// Not every port has the same capabilities.
type Port struct {
	AggregatedBy       FlexBool   `json:"aggregated_by"`
	Autoneg            FlexBool   `json:"autoneg,omitempty"`
	BytesR             FlexInt    `json:"bytes-r"`
	DNS                []string   `json:"dns,omitempty"`
	Dot1XMode          string     `json:"dot1x_mode"`
	Dot1XStatus        string     `json:"dot1x_status"`
	Enable             FlexBool   `json:"enable"`
	FlowctrlRx         FlexBool   `json:"flowctrl_rx"`
	FlowctrlTx         FlexBool   `json:"flowctrl_tx"`
	FullDuplex         FlexBool   `json:"full_duplex"`
	IP                 string     `json:"ip,omitempty"`
	Ifname             string     `json:"ifname,omitempty"`
	IsUplink           FlexBool   `json:"is_uplink"`
	Mac                string     `json:"mac,omitempty"`
	MacTable           []MacTable `json:"mac_table,omitempty"`
	Jumbo              FlexBool   `json:"jumbo,omitempty"`
	Masked             FlexBool   `json:"masked"`
	Media              string     `json:"media"`
	Name               string     `json:"name"`
	NetworkName        string     `json:"network_name,omitempty"`
	Netmask            string     `json:"netmask,omitempty"`
	NumPort            FlexInt    `json:"num_port,omitempty"`
	OpMode             string     `json:"op_mode"`
	PoeCaps            FlexInt    `json:"poe_caps"`
	PoeClass           string     `json:"poe_class,omitempty"`
	PoeCurrent         FlexInt    `json:"poe_current,omitempty"`
	PoeEnable          FlexBool   `json:"poe_enable,omitempty"`
	PoeGood            FlexBool   `json:"poe_good,omitempty"`
	PoeMode            string     `json:"poe_mode,omitempty"`
	PoePower           FlexInt    `json:"poe_power,omitempty"`
	PoeVoltage         FlexInt    `json:"poe_voltage,omitempty"`
	PortDelta          PortDelta  `json:"port_delta,omitempty"`
	PortIdx            FlexInt    `json:"port_idx"`
	PortPoe            FlexBool   `json:"port_poe"`
	PortconfID         string     `json:"portconf_id"`
	RxBroadcast        FlexInt    `json:"rx_broadcast"`
	RxBytes            FlexInt    `json:"rx_bytes"`
	RxBytesR           FlexInt    `json:"rx_bytes-r"`
	RxDropped          FlexInt    `json:"rx_dropped"`
	RxErrors           FlexInt    `json:"rx_errors"`
	RxMulticast        FlexInt    `json:"rx_multicast"`
	RxPackets          FlexInt    `json:"rx_packets"`
	RxRate             FlexInt    `json:"rx_rate,omitempty"`
	Satisfaction       FlexInt    `json:"satisfaction,omitempty"`
	SatisfactionReason FlexInt    `json:"satisfaction_reason"`
	SFPCompliance      string     `json:"sfp_compliance"`
	SFPCurrent         FlexInt    `json:"sfp_current"`
	SFPFound           FlexBool   `json:"sfp_found"`
	SFPPart            string     `json:"sfp_part"`
	SFPRev             string     `json:"sfp_rev"`
	SFPRxfault         FlexBool   `json:"sfp_rxfault"`
	SFPRxpower         FlexInt    `json:"sfp_rxpower"`
	SFPSerial          string     `json:"sfp_serial"`
	SFPTemperature     FlexInt    `json:"sfp_temperature"`
	SFPTxfault         FlexBool   `json:"sfp_txfault"`
	SFPTxpower         FlexInt    `json:"sfp_txpower"`
	SFPVendor          string     `json:"sfp_vendor"`
	SFPVoltage         FlexInt    `json:"sfp_voltage"`
	Speed              FlexInt    `json:"speed"`
	SpeedCaps          FlexInt    `json:"speed_caps"`
	StpPathcost        FlexInt    `json:"stp_pathcost"`
	StpState           string     `json:"stp_state"`
	TxBroadcast        FlexInt    `json:"tx_broadcast"`
	TxBytes            FlexInt    `json:"tx_bytes"`
	TxBytesR           FlexInt    `json:"tx_bytes-r"`
	TxDropped          FlexInt    `json:"tx_dropped"`
	TxErrors           FlexInt    `json:"tx_errors"`
	TxMulticast        FlexInt    `json:"tx_multicast"`
	TxPackets          FlexInt    `json:"tx_packets"`
	TxRate             FlexInt    `json:"tx_rate,omitempty"`
	Type               string     `json:"type,omitempty"`
	Up                 FlexBool   `json:"up"`
}
