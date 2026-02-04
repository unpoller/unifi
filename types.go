package unifi

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	gofakeit.AddFuncLookup("port", gofakeit.Info{
		Category:    "custom",
		Description: "Random Unifi Port integer value",
		Example:     "8443",
		Output:      "int",
		Generate: func(r *rand.Rand, _ *gofakeit.MapParams, _ *gofakeit.Info) (interface{}, error) {
			return r.Int31n(65535), nil
		},
	})

	gofakeit.AddFuncLookup("timestamp", gofakeit.Info{
		Category:    "custom",
		Description: "Recent timestamp value",
		Example:     "123456",
		Output:      "int64",
		Generate: func(_ *rand.Rand, _ *gofakeit.MapParams, _ *gofakeit.Info) (interface{}, error) {
			return gofakeit.DateRange(time.Now().Add(-time.Second*59), time.Now().Add(-time.Second)).Unix(), nil
		},
	})

	gofakeit.AddFuncLookup("recent_time", gofakeit.Info{
		Category:    "custom",
		Description: "Recent time.Time value",
		Example:     "time.Now().Add(-time.Second)",
		Output:      "time.Time",
		Generate: func(_ *rand.Rand, _ *gofakeit.MapParams, _ *gofakeit.Info) (interface{}, error) {
			return gofakeit.DateRange(time.Now().Add(-time.Second*59), time.Now().Add(-time.Second)), nil
		},
	})

	gofakeit.AddFuncLookup("timestamps", gofakeit.Info{
		Category:    "custom",
		Description: "Recent timestamp values",
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
		Generate: func(_ *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
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

	gofakeit.AddFuncLookup("constFlexBool", gofakeit.Info{
		Category:    "custom",
		Description: "Configured FlexBool",
		Example:     "FlexBool{Val: false, Txt: \"false\"}",
		Output:      "FlexBool",
		Params: []gofakeit.Param{
			{
				Field:       "value",
				Display:     "value",
				Type:        "bool",
				Optional:    true,
				Default:     "false",
				Description: "The default value",
			},
		},
		Generate: func(_ *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			l, err := info.GetBool(m, "value")
			if err != nil {
				return nil, err
			}

			return *NewFlexBool(l), nil
		},
	})

	gofakeit.AddFuncLookup("tempStatusByName", gofakeit.Info{
		Category:    "custom",
		Description: "Configured TempStatusByName",
		Example:     "TempStatusByName{...}",
		Output:      "TempStatusByName",
		Generate: func(r *rand.Rand, _ *gofakeit.MapParams, _ *gofakeit.Info) (interface{}, error) {
			return TempStatusByName{
				"cpu":     NewFlexTemp(float64(r.Int31n(100))),
				"sys":     NewFlexTemp(float64(r.Int31n(100))),
				"probe":   NewFlexTemp(float64(r.Int31n(100))),
				"memory":  NewFlexTemp(float64(r.Int31n(100))),
				"network": NewFlexTemp(float64(r.Int31n(100))),
			}, nil
		},
	})
}

var ErrCannotUnmarshalFlexInt = fmt.Errorf("cannot unmarshal to FlexInt")
var ErrCannotUnmarshalFlexString = fmt.Errorf("cannot unmarshal to FlexString")

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
	// APIClientHistoryPath is Unifi Clients History API Path.
	APIClientHistoryPath string = "/v2/api/site/%s/clients/history?%s"
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
	APIEventPathIDS string = "/api/s/%s/stat/ips/event" //nolint:revive
	// APIEventPathAlarms contains the site alarms.
	APIEventPathAlarms string = "/api/s/%s/list/alarm"
	// APIPrefixNew is the prefix added to the new API paths; except login. duh.
	APIPrefixNew string = "/proxy/network"
	// APIAnomaliesPath returns site anomalies.
	APIAnomaliesPath string = "/api/s/%s/stat/anomalies"
	// APISystemLogPath returns system log events (v2 API).
	APISystemLogPath          string = "/v2/api/site/%s/system-log/all"
	APICommandPath            string = "/api/s/%s/cmd"
	APIDevMgrPath             string = APICommandPath + "/devmgr"
	APIClientTrafficPath      string = "/v2/api/site/%s/traffic?start=%d&end=%d&includeUnidentified=%t"
	APIClientTrafficByMacPath string = "/v2/api/site/%s/traffic/%s?start=%d&end=%d&includeUnidentified=%t&mac=%s"
	APICountryTrafficPath     string = "/v2/api/site/%s/country-traffic?start=%d&end=%d"
	APIAggregatedDashboard    string = "/v2/api/site/%s/aggregated-dashboard?historySeconds=%d"
	// APIProtectLogPath returns Protect system log events.
	APIProtectLogPath string = "/proxy/protect/api/events/system-logs"
	// APIProtectEventsPath is the base path for Protect events (for thumbnails, etc.).
	APIProtectEventsPath string = "/proxy/protect/api/events"
	// APIDeviceTagsPath returns device tags for a site.
	APIDeviceTagsPath string = "/proxy/network/v2/api/site/%s/device-tags"
	// APIActiveDHCPLeasesPath returns active DHCP leases for a site.
	APIActiveDHCPLeasesPath string = "/proxy/network/v2/api/site/%s/active-leases"
	// APIWANEnrichedConfigPath returns enriched WAN configuration with statistics.
	APIWANEnrichedConfigPath string = "/proxy/network/v2/api/site/%s/wan/enriched-configuration"
	// APIWANISPStatusPath returns WAN interface status (ACTIVE/BACKUP).
	APIWANISPStatusPath string = "/proxy/network/v2/api/site/%s/wan/%s/isp-status"
	// APIWANLoadBalancingStatusPath returns load balancing status for WAN interfaces.
	APIWANLoadBalancingStatusPath string = "/proxy/network/v2/api/site/%s/wan/load-balancing/status"
	// APIWANLoadBalancingConfigPath returns load balancing configuration for WAN interfaces.
	APIWANLoadBalancingConfigPath string = "/proxy/network/v2/api/site/%s/wan/load-balancing/configuration"
	// APIWANSLAsPath returns WAN SLA monitoring data (latency, packet loss, jitter).
	APIWANSLAsPath string = "/proxy/network/v2/api/site/%s/wan-slas"
	// APISysinfoPath returns controller system info and health (UniFi OS).
	APISysinfoPath string = "/api/s/%s/stat/sysinfo"
)

// path returns the correct api path based on the new variable.
// new is based on the unifi-controller output. is it new or old output?
func (u *Unifi) path(path string) string {
	if u.new {
		if path == APILoginPath {
			return APILoginPathNew
		}

		// Don't add prefix if path already has /proxy/network or starts with /proxy/ (e.g., /proxy/protect/)
		if !strings.HasPrefix(path, APIPrefixNew) && !strings.HasPrefix(path, "/proxy/") && path != APILoginPathNew {
			return APIPrefixNew + path
		}
	}

	return path
}

// Logger is a base type to deal with changing log outputs. Create a logger
// that matches this interface to capture debug and error logs.
type Logger func(msg string, fmt ...interface{})

// discardLogs is the default debug logger.
func discardLogs(_ string, _ ...interface{}) {
	// do nothing.
}

// Devices contains a list of all the unifi devices from a controller.
// Contains Access points, security gateways and switches.
type Devices struct {
	UAPs []*UAP `fakesize:"5"`
	USGs []*USG `fakesize:"5"`
	USWs []*USW `fakesize:"5"`
	UDMs []*UDM `fakesize:"5"`
	UXGs []*UXG `fakesize:"5"`
	PDUs []*PDU `fakesize:"5"`
	UBBs []*UBB `fakesize:"5"`
	UCIs []*UCI `fakesize:"5"`
}

// Config is the data passed into our library. This configures things and allows
// us to connect to a controller and write log messages. Optional SSLCert is used
// for ssl cert pinning; provide the content of a PEM to validate the server's cert.
type Config struct {
	User      string
	Pass      string
	APIKey    string
	URL       string
	SSLCert   [][]byte
	ErrorLog  Logger
	DebugLog  Logger
	Timeout   time.Duration // how long to wait for replies, default: forever.
	VerifySSL bool
}

type UnifiClient interface { //nolint: revive
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
	// GetClients returns a response full of client history data from the UniFi Controller.
	GetClientHistory(sites []*Site, opts *ClientHistoryOpts) ([]*ClientHistory, error)
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
	// GetUBBs returns all UBB devices, an error, or nil if there are no UBBs.
	GetUBBs(site *Site) ([]*UBB, error)
	// GetUCIs returns all UCI devices, an error, or nil if there are no UCIs.
	GetUCIs(site *Site) ([]*UCI, error)
	// GetPDUs returns all PDU devices, an error, or nil if there are no PDUs.
	GetPDUs(site *Site) ([]*PDU, error)
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
	// GetActiveDHCPLeases returns active DHCP leases for the given sites.
	GetActiveDHCPLeases(sites []*Site) ([]*DHCPLease, error)
	// GetActiveDHCPLeasesWithAssociations returns active DHCP leases enriched with client and device associations.
	GetActiveDHCPLeasesWithAssociations(sites []*Site) ([]*DHCPLease, error)
	// AssociateDHCPLeases associates DHCP leases with clients, devices, and networks using pre-fetched data.
	AssociateDHCPLeases(leases []*DHCPLease, clients []*Client, devices *Devices, networks []Network) error
	// GetWANEnrichedConfiguration returns enriched WAN configuration for all WAN interfaces.
	GetWANEnrichedConfiguration(sites []*Site) ([]*WANEnrichedConfiguration, error)
	// GetWANLoadBalancingStatus returns the current load balancing status for WAN interfaces.
	GetWANLoadBalancingStatus(sites []*Site) (*WANLoadBalancingStatus, error)
	// GetWANISPStatus returns the ISP status for WAN interfaces.
	GetWANISPStatus(sites []*Site, wanNetworkgroup string) (*WANISPStatusDetailed, error)
	// GetWANSLAs returns WAN SLA monitoring data.
	GetWANSLAs(sites []*Site) ([]*WANSLA, error)
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
	GetServerData() (*ServerStatus, error)
	// GetUsers returns a response full of clients that connected to the UDM within the provided amount of time
	// using the insight historical connection data set.
	GetUsers(sites []*Site, hours int) ([]*User, error)
	// GetClientTraffic returns a response full of clients' traffic data from the UniFi Controller for the provided time period.
	GetClientTraffic(sites []*Site, epochMillisTimePeriod *EpochMillisTimePeriod, includeUnidentified bool) ([]*ClientUsageByApp, error)
	// GetClientTrafficByMac returns a response full of clients' traffic data from the UniFi Controller for the provided time period
	// and each of the mac addressees provided.
	GetClientTrafficByMac(site *Site, epochMillisTimePeriod *EpochMillisTimePeriod, includeUnidentified bool, macs ...string) ([]*ClientUsageByApp, error)
	// GetCountryTraffic returns a response full of clients' traffic data from the UniFi Controller for the provided time period.'
	GetCountryTraffic(sites []*Site, epochMillisTimePeriod *EpochMillisTimePeriod) ([]*UsageByCountry, error)
	// GetProtectLogs returns Protect system log events.
	GetProtectLogs(req *ProtectLogRequest) ([]*ProtectLogEntry, error)
	// GetSysinfo returns controller system info and health (UniFi OS).
	GetSysinfo(sites []*Site) ([]*Sysinfo, error)
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	*Config
	*ServerStatus
	csrf                        string
	fingerprints                fingerprints
	new                         bool
	deviceTagsUnavailableOnce   sync.Once
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
	ServerVersion string   `fake:"{appversion}" json:"server_version"`
	UUID          string   `fake:"{uuid}"       json:"uuid"`
}

type FlexString struct {
	Val         string
	Arr         []string
	hintIsArray bool
}

// DeviceTag represents a device tag from the UniFi API.
// Device tags allow grouping devices for organization and filtering.
type DeviceTag struct {
	ID               string   `json:"_id"`
	Name             string   `json:"name"`
	MemberDeviceMacs []string `json:"member_device_macs"`
}

func NewFlexString(v string) *FlexString {
	return &FlexString{
		Val:         v,
		Arr:         []string{v},
		hintIsArray: false,
	}
}

func NewFlexStringArray(v []string) *FlexString {
	return &FlexString{
		Val:         strings.Join(v, ", "),
		Arr:         v,
		hintIsArray: true,
	}
}

// UnmarshalJSON converts a string or number to an integer.
// Generally, do not call this directly, it's used in the json interface.
func (f *FlexString) UnmarshalJSON(b []byte) error {
	var ust interface{}

	if err := json.Unmarshal(b, &ust); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	switch i := ust.(type) {
	case []interface{}:
		f.hintIsArray = true
		// try to cast to string
		for _, v := range i {
			if s, ok := v.(string); ok {
				f.Arr = append(f.Arr, s)
			}
		}

		f.Val = strings.Join(f.Arr, ", ")
	case []string:
		f.hintIsArray = true
		f.Val = strings.Join(i, ", ")
		f.Arr = i
	case string:
		f.Val = i
		f.Arr = []string{i}
	case nil:
		// noop, consider it empty values
	default:
		return fmt.Errorf("%v: %w", b, ErrCannotUnmarshalFlexString)
	}

	return nil
}

func (f FlexString) MarshalJSON() ([]byte, error) {
	// array case
	if f.hintIsArray {
		return json.Marshal(f.Arr)
	}

	// plain string case
	return json.Marshal(f.Val)
}

func (f FlexString) String() string {
	return f.Val
}

func (f FlexString) Fake(faker *gofakeit.Faker) interface{} {
	randValue := math.Min(math.Max(0.1, math.Abs(faker.Rand.Float64())), 120)
	s := fmt.Sprintf("fake-%0.2f", randValue)

	if faker.Rand.Intn(2) == 0 {
		// plain string value
		return FlexString{
			Val: s,
			Arr: []string{s},
		}
	}

	// array case
	s2 := fmt.Sprintf("fake-%0.2f-2", randValue)
	s3 := fmt.Sprintf("fake-%0.2f-3", randValue)
	arr := []string{s, s2, s3}

	return FlexString{
		Val: strings.Join(arr, ", "),
		Arr: arr,
	}
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
// Generally, do not call this directly, it's used in the json interface.
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

func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Val)
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
func (f FlexInt) Fake(faker *gofakeit.Faker) interface{} {
	randValue := math.Min(math.Max(1, math.Abs(faker.Rand.Float64())), 500)

	if faker.Rand.Intn(2) == 0 {
		// int-value
		return FlexInt{
			Val: float64(int64(randValue)),
			Txt: strconv.FormatInt(int64(randValue), 10),
		}
	}

	return FlexInt{
		Val: randValue,
		Txt: strconv.FormatFloat(randValue, 'f', 8, 64),
	}
}

// FlexBool provides a container and unmarshalling for fields that may be
// boolean or strings in the Unifi API.
type FlexBool struct {
	Val bool
	Txt string
}

func NewFlexBool(v bool) *FlexBool {
	textValue := "false"

	if v {
		textValue = "true"
	}

	return &FlexBool{
		Val: v,
		Txt: textValue,
	}
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

func (f FlexBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Val)
}

func (f *FlexBool) String() string {
	return f.Txt
}

func (f *FlexBool) Float64() float64 {
	if f.Val {
		return 1
	}

	return 0
}

// Fake implements gofakeit Fake interface
func (f FlexBool) Fake(faker *gofakeit.Faker) interface{} {
	opts := []bool{
		true,
		false,
	}

	v := opts[faker.Rand.Intn(2)]

	return FlexBool{
		Val: v,
		Txt: strconv.FormatBool(v),
	}
}

// FlexTemp provides a container and unmarshalling for fields that may be
// numbers or strings in the Unifi API as temperatures.
type FlexTemp struct {
	Val float64 // in Celsius
	Txt string
}

func NewFlexTemp(v float64) *FlexTemp {
	return &FlexTemp{
		Val: v,
		Txt: strconv.FormatFloat(v, 'f', -1, 64),
	}
}

// UnmarshalJSON converts a string or number to Celsius (stored as float64).
// Generally, do not call this directly, it's used in the json interface.
// Accepts float64, int, int64, string (e.g. "72 C"), or nil so that all API
// variants are stored as float64 and avoid InfluxDB field type conflicts.
func (f *FlexTemp) UnmarshalJSON(b []byte) error {
	var unk interface{}

	if err := json.Unmarshal(b, &unk); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	switch i := unk.(type) {
	case float64:
		f.Val = i
		f.Txt = strconv.FormatFloat(i, 'f', -1, 64)
	case int:
		f.Val = float64(i)
		f.Txt = strconv.FormatInt(int64(i), 10)
	case int64:
		f.Val = float64(i)
		f.Txt = strconv.FormatInt(i, 10)
	case string:
		f.Txt = i
		parts := strings.SplitN(i, " ", 2)

		if len(parts) == 2 {
			// format is: $val(int or float) $unit(C or F)
			f.Val, _ = strconv.ParseFloat(parts[0], 64)
		} else {
			// assume Celsius
			f.Val, _ = strconv.ParseFloat(i, 64)
		}
	case nil:
		f.Txt = "0"
		f.Val = 0
	default:
		return fmt.Errorf("%v: %w", b, ErrCannotUnmarshalFlexInt)
	}

	return nil
}

func (f FlexTemp) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Val)
}

// Celsius returns the temperature in Celsius as float64.
// Use this (or CelsiusSafe) for InfluxDB/metrics to avoid field type conflicts; prefer over CelsiusInt/CelsiusInt64.
func (f *FlexTemp) Celsius() float64 {
	return f.Val
}

// CelsiusSafe returns the temperature in Celsius as float64, or 0 if f is nil.
// Use for InfluxDB and metrics so the field is always float64 (avoids "field type conflict" errors).
func (f *FlexTemp) CelsiusSafe() float64 {
	if f == nil {
		return 0
	}

	return f.Val
}

func (f *FlexTemp) CelsiusInt() int {
	return int(f.Val)
}

func (f *FlexTemp) CelsiusInt64() int64 {
	return int64(f.Val)
}

func (f *FlexTemp) Fahrenheit() float64 {
	return (f.Val * (9 / 5)) + 32
}

func (f *FlexTemp) FahrenheitInt() int {
	return int(f.Fahrenheit())
}

func (f *FlexTemp) FahrenheitInt64() int64 {
	return int64(f.Fahrenheit())
}

func (f *FlexTemp) String() string {
	return f.Txt
}

func (f *FlexTemp) Add(o *FlexTemp) {
	f.Val += o.Val
	f.Txt = strconv.FormatFloat(f.Val, 'f', -1, 64)
}

func (f *FlexTemp) AddFloat64(v float64) {
	f.Val += v
	f.Txt = strconv.FormatFloat(f.Val, 'f', -1, 64)
}

// Fake implements gofakeit Fake interface
func (f FlexTemp) Fake(faker *gofakeit.Faker) interface{} {
	randValue := math.Min(math.Max(0.1, math.Abs(faker.Rand.Float64())), 120)
	if faker.Rand.Intn(2) == 0 {
		// int-value
		return FlexTemp{
			Val: float64(int64(randValue)),
			Txt: strconv.FormatInt(int64(randValue), 10) + " C",
		}
	}

	return FlexTemp{
		Val: randValue,
		Txt: strconv.FormatFloat(randValue, 'f', 8, 64) + " C",
	}
}

// DownlinkTable is part of a UXG and UDM output.
type DownlinkTable struct {
	PortIdx    FlexInt  `json:"port_idx"`
	Speed      FlexInt  `json:"speed"`
	FullDuplex FlexBool `json:"full_duplex"`
	Mac        string   `fake:"{macaddress}" json:"mac"`
}

// ConfigNetwork comes from gateways.
type ConfigNetwork struct {
	Type string `fake:"{randomstring:[wan,lan,vlan]}" json:"type"`
	IP   string `fake:"{ipv4address}"                 json:"ip"`
}

type EthernetTable struct {
	Mac     string  `fake:"{macaddress}" json:"mac"`
	NumPort FlexInt `json:"num_port"`
	Name    string  `fake:"{animal}"     json:"name"`
}

// Port is a physical connection on a USW or Gateway.
// Not every port has the same capabilities.
type Port struct {
	AggregatedBy       FlexBool   `json:"aggregated_by"`
	Autoneg            FlexBool   `json:"autoneg,omitempty"`
	BytesR             FlexInt    `json:"bytes-r"`
	DNS                []string   `fakesize:"5"                                                    json:"dns,omitempty"`
	Dot1XMode          string     `json:"dot1x_mode"`
	Dot1XStatus        string     `json:"dot1x_status"`
	Enable             FlexBool   `json:"enable"`
	FlowctrlRx         FlexBool   `json:"flowctrl_rx"`
	FlowctrlTx         FlexBool   `json:"flowctrl_tx"`
	FullDuplex         FlexBool   `json:"full_duplex"`
	IP                 string     `fake:"{ipv4address}"                                            json:"ip,omitempty"`
	Ifname             string     `fake:"{randomstring:[wlan0,wlan1,lan0,lan1,vlan1,vlan0,vlan2]}" json:"ifname,omitempty"`
	IsUplink           FlexBool   `json:"is_uplink"`
	Mac                string     `fake:"{macaddress}"                                             json:"mac,omitempty"`
	MacTable           []MacTable `fakesize:"5"                                                    json:"mac_table,omitempty"`
	Jumbo              FlexBool   `json:"jumbo,omitempty"`
	Masked             FlexBool   `json:"masked"`
	Media              string     `json:"media"`
	Name               string     `fake:"{animal}"                                                 json:"name"`
	NetworkName        string     `fake:"{animal}"                                                 json:"network_name,omitempty"`
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
	PortPoe            FlexBool   `fake:"{constFlexBool:true}"                                     json:"port_poe"`
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
	SFPFound           FlexBool   `fake:"{constFlexBool:true}"                                     json:"sfp_found"`
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
