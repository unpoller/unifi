package unifi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	gofakeit.AddFuncLookup("port", gofakeit.Info{
		Category:    "custom",
		Description: "Random Unifi Port integer value",
		Example:     "8443",
		Output:      "int",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			return r.Int31n(65535), nil
		},
	})

	gofakeit.AddFuncLookup("timestamp", gofakeit.Info{
		Category:    "custom",
		Description: "Random timestamp value",
		Example:     "123456",
		Output:      "int64",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			return gofakeit.DateRange(time.Now().Add(time.Hour-2), time.Now()).Unix(), nil
		},
	})

	gofakeit.AddFuncLookup("timestamps", gofakeit.Info{
		Category:    "custom",
		Description: "Random timestamp value",
		Example:     "123456",
		Output:      "[]int64",
		Params: []gofakeit.Param{
			{
				Field:       "length",
				Display:     "number of items to generate",
				Type:        "uint",
				Optional:    false,
				Default:     "2",
				Description: "The number of ints to generate",
			},
		},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			l, err := info.GetUint(m, "length")
			if err != nil {
				return nil, err
			}
			result := make([]int64, 0)
			for i := 0; i < int(l); i++ {
				result = append(result, gofakeit.DateRange(time.Now().Add(time.Hour-2), time.Now()).Unix())
			}
			return result, nil
		},
	})
}

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

type UnifiClient interface {
	// GetAlarms returns Alarms for a list of Sites.
	GetAlarms(sites []*Site) ([]*Alarm, error)
	// GetAlarmsSite retreives the Alarms for a single Site.
	GetAlarmsSite(site *Site) ([]*Alarm, error)
	// GetAnomalies returns Anomalies for a list of Sites.
	GetAnomalies(sites []*Site, timeRange ...time.Time) ([]*Anomaly, error)
	// GetAnomaliesSite retreives the Anomalies for a single Site.
	GetAnomaliesSite(site *Site, timeRange ...time.Time) ([]*Anomaly, error)
	// GetClients returns a response full of clients' data from the UniFi Controller.
	GetClients(sites []*Site) ([]*Client, error)
	// GetClientsDPI garners dpi data for clients.
	GetClientsDPI(sites []*Site) ([]*DPITable, error)
	// GetDevices returns a response full of devices' data from the UniFi Controller.
	GetDevices(sites []*Site) (*Devices, error)
	// GetUSWs returns all switches, an error, or nil if there are no switches.
	GetUSWs(site *Site) ([]*USW, error)
	// GetUAPs returns all access points, an error, or nil if there are no APs.
	GetUAPs(site *Site) ([]*UAP, error)
	// GetUDMs returns all dream machines, an error, or nil if there are no UDMs.
	GetUDMs(site *Site) ([]*UDM, error)
	// GetUXGs returns all 10Gb gateways, an error, or nil if there are no UXGs.
	GetUXGs(site *Site) ([]*UXG, error)
	// GetUSGs returns all 1Gb gateways, an error, or nil if there are no USGs.
	GetUSGs(site *Site) ([]*USG, error)
	// GetEvents returns a response full of UniFi Events for the last 1 hour from multiple sites.
	GetEvents(sites []*Site, hours time.Duration) ([]*Event, error)
	// GetSiteEvents retrieves the last 1 hour's worth of events from a single site.
	GetSiteEvents(site *Site, hours time.Duration) ([]*Event, error)
	// GetIDS returns Intrusion Detection Systems events for a list of Sites.
	// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
	// Events between start and end are returned. End defaults to time.Now().
	GetIDS(sites []*Site, timeRange ...time.Time) ([]*IDS, error)
	// GetIDSSite retrieves the Intrusion Detection System Data for a single Site.
	// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
	// Events between start and end are returned. End defaults to time.Now().
	GetIDSSite(site *Site, timeRange ...time.Time) ([]*IDS, error)
	// GetNetworks returns a response full of network data from the UniFi Controller.
	GetNetworks(sites []*Site) ([]Network, error)
	// GetSites returns a list of configured sites on the UniFi controller.
	GetSites() ([]*Site, error)
	// GetSiteDPI garners dpi data for sites.
	GetSiteDPI(sites []*Site) ([]*DPITable, error)
	// GetRogueAPs returns RogueAPs for a list of Sites.
	// Use GetRogueAPsSite if you want more control.
	GetRogueAPs(sites []*Site) ([]*RogueAP, error)
	// GetRogueAPsSite returns RogueAPs for a single Site.
	GetRogueAPsSite(site *Site) ([]*RogueAP, error)
	// Login is a helper method. It can be called to grab a new authentication cookie.
	Login() error
	// Logout closes the current session.
	Logout() error
	// GetServerData sets the controller's version and UUID. Only call this if you
	// previously called Login and suspect the controller version has changed.
	GetServerData() error
	// GetUsers returns a response full of clients that connected to the UDM within the provided amount of time
	// using the insight historical connection data set.
	GetUsers(sites []*Site, hours int) ([]*User, error)
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	*Config
	server       *ServerStatus
	csrf         string
	fingerprints fingerprints
	new          bool
}

// ensure Unifi implements UnifiClient fully, will fail to compile otherwise
var _ UnifiClient = &Unifi{}

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

// ServerStatus is the /status endpoint from the Unifi controller.
type ServerStatus struct {
	Up            FlexBool `json:"up"`
	ServerVersion string   `json:"server_version"`
	UUID          string   `json:"uuid" fake:"{uuid}"`
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

// Fake implements gofakeit Fake interface
func (f *FlexInt) Fake(faker *gofakeit.Faker) interface{} {
	randValue := faker.Rand.Float64()
	opts := []interface{}{
		randValue,
		int64(randValue),
		strconv.FormatInt(int64(randValue), 10),
		strconv.FormatFloat(randValue, 'f', 8, 64),
	}

	return opts[faker.Rand.Intn(2)]
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

// Fake implements gofakeit Fake interface
func (f *FlexBool) Fake(faker *gofakeit.Faker) interface{} {
	opts := []interface{}{
		"true",
		true,
		"false",
		false,
	}
	return opts[faker.Rand.Intn(4)]
}

// DownlinkTable is part of a UXG and UDM output.
type DownlinkTable struct {
	PortIdx    FlexInt  `json:"port_idx" fake:"{port}"`
	Speed      FlexInt  `json:"speed"`
	FullDuplex FlexBool `json:"full_duplex"`
	Mac        string   `json:"mac" fake:"{macaddress}"`
}

// ConfigNetwork comes from gateways.
type ConfigNetwork struct {
	Type string `json:"type" fake:"{randomstring:[wan,lan,vlan]}"`
	IP   string `json:"ip" fake:"{ipv4address}"`
}

type EthernetTable struct {
	Mac     string  `json:"mac" fake:"{macaddress}"`
	NumPort FlexInt `json:"num_port" fake:"{port}"`
	Name    string  `json:"name" fake:"{animal}"`
}

// Port is a physical connection on a USW or Gateway.
// Not every port has the same capabilities.
type Port struct {
	AggregatedBy       FlexBool   `json:"aggregated_by"`
	Autoneg            FlexBool   `json:"autoneg,omitempty"`
	BytesR             FlexInt    `json:"bytes-r"`
	DNS                []string   `json:"dns,omitempty" fakesize:"5"`
	Dot1XMode          string     `json:"dot1x_mode"`
	Dot1XStatus        string     `json:"dot1x_status"`
	Enable             FlexBool   `json:"enable"`
	FlowctrlRx         FlexBool   `json:"flowctrl_rx"`
	FlowctrlTx         FlexBool   `json:"flowctrl_tx"`
	FullDuplex         FlexBool   `json:"full_duplex"`
	IP                 string     `json:"ip,omitempty" fake:"{ipv4address}"`
	Ifname             string     `json:"ifname,omitempty" fake:"{randomstring:[wlan0,wlan1,lan0,lan1,vlan1,vlan0,vlan2]}"`
	IsUplink           FlexBool   `json:"is_uplink"`
	Mac                string     `json:"mac,omitempty" fake:"{macaddress}"`
	MacTable           []MacTable `json:"mac_table,omitempty" fakesize:"5"`
	Jumbo              FlexBool   `json:"jumbo,omitempty"`
	Masked             FlexBool   `json:"masked"`
	Media              string     `json:"media"`
	Name               string     `json:"name" fake:"{animal}"`
	NetworkName        string     `json:"network_name,omitempty" fake:"{animal}"`
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
