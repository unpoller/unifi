package unifi

import (
	"fmt"
	"net/url"
	"strings"
)

// GetClients returns a response full of clients' data from the UniFi Controller.
func (u *Unifi) GetClients(sites []*Site) ([]*Client, error) {
	data := make([]*Client, 0)

	for _, site := range sites {
		var response struct {
			Data []*Client `json:"data"`
		}

		u.DebugLog("Polling Controller, retreiving UniFi Clients, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIClientPath, site.Name)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for i, d := range response.Data {
			// Add special SourceName value.
			response.Data[i].SourceName = u.URL
			// Add the special "Site Name" to each client. This becomes a Grafana filter somewhere.
			response.Data[i].SiteName = site.SiteName
			// Fix name and hostname fields. Sometimes one or the other is blank.
			response.Data[i].Hostname = strings.TrimSpace(pick(d.Hostname, d.Name, d.Mac))
			response.Data[i].Name = strings.TrimSpace(pick(d.Name, d.Hostname))
		}

		data = append(data, response.Data...)
	}

	return data, nil
}

// GetClientsDPI garners dpi data for clients.
func (u *Unifi) GetClientsDPI(sites []*Site) ([]*DPITable, error) {
	var data []*DPITable

	for _, site := range sites {
		u.DebugLog("Polling Controller, retreiving Client DPI data, site %s", site.SiteName)

		var response struct {
			Data []*DPITable `json:"data"`
		}

		clientDPIpath := fmt.Sprintf(APIClientDPI, site.Name)
		if err := u.GetData(clientDPIpath, &response, `{"type":"by_app"}`); err != nil {
			return nil, err
		}

		for _, d := range response.Data {
			d.SourceName = site.SourceName
			d.SiteName = site.SiteName
			data = append(data, d)
		}
	}

	return data, nil
}

// GetClientHistory pulls client history data from the controller.
func (u *Unifi) GetClientHistory(sites []*Site, opts *ClientHistoryOpts) ([]ClientHistory, error) {
	// TODO: could use external library to convert struct to url query params ?
	if opts == nil {
		opts = NewClientHistoryOpts()
	}
	params := url.Values{}
	params.Add("onlyNonBlocked", fmt.Sprint(opts.OnlyNonBlocked))
	params.Add("includeUnifiDevices", fmt.Sprint(opts.IncludeUnifiDevices))
	params.Add("withinHours", fmt.Sprint(opts.WithinHours))

	data := make([]ClientHistory, 0)

	for _, site := range sites {
		response := []ClientHistory{}

		u.DebugLog("Polling Controller, retreiving UniFi Client History, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIClientHistoryPath, site.Name)
		if err := u.GetData(clientPath, &response, params.Encode()); err != nil {
			return nil, err
		}

		for i, d := range response {
			// Add special SourceName value.
			response[i].SourceName = u.URL
			// Add the special "Site Name" to each client. This becomes a Grafana filter somewhere.
			response[i].SiteName = site.SiteName
			// Fix name and hostname fields. Sometimes one or the other is blank.
			response[i].Hostname = strings.TrimSpace(pick(d.Hostname, d.Name, d.Mac))
			response[i].Name = strings.TrimSpace(pick(d.Name, d.Hostname))
		}

		data = append(data, response...)
	}

	return data, nil
}

// Client defines all the data a connected-network client contains.
type Client struct {
	Anomalies        FlexInt  `json:"anomalies,omitempty"`
	ApMac            string   `fake:"{macaddress}"           json:"ap_mac"`
	ApName           string   `json:"-"`
	AssocTime        FlexInt  `json:"assoc_time"`
	Blocked          bool     `json:"blocked,omitempty"`
	Bssid            string   `fake:"{macaddress}"           json:"bssid"`
	BytesR           FlexInt  `json:"bytes-r"`
	Ccq              FlexInt  `json:"ccq"`
	Channel          FlexInt  `json:"channel"`
	DevCat           FlexInt  `json:"dev_cat"`
	DevFamily        FlexInt  `json:"dev_family"`
	DevID            FlexInt  `json:"dev_id"`
	DevVendor        FlexInt  `json:"dev_vendor,omitempty"`
	DhcpendTime      FlexInt  `json:"dhcpend_time,omitempty"`
	Essid            string   `fake:"{macaddress}"           json:"essid"`
	FirstSeen        FlexInt  `json:"first_seen"`
	FixedIP          string   `fake:"{ipv4address}"          json:"fixed_ip"`
	GwMac            string   `fake:"{macaddress}"           json:"gw_mac"`
	GwName           string   `json:"-"`
	Hostname         string   `json:"hostname"`
	ID               string   `fake:"{uuid}"                 json:"_id"`
	IP               string   `fake:"{ipv4address}"          json:"ip"`
	IdleTime         FlexInt  `json:"idle_time"`
	Is11R            FlexBool `json:"is_11r"`
	IsGuest          FlexBool `json:"is_guest"`
	IsGuestByUAP     FlexBool `json:"_is_guest_by_uap"`
	IsGuestByUGW     FlexBool `json:"_is_guest_by_ugw"`
	IsGuestByUSW     FlexBool `json:"_is_guest_by_usw"`
	IsWired          FlexBool `json:"is_wired"`
	LastSeen         FlexInt  `json:"last_seen"`
	LastSeenByUAP    FlexInt  `json:"_last_seen_by_uap"`
	LastSeenByUGW    FlexInt  `json:"_last_seen_by_ugw"`
	LastSeenByUSW    FlexInt  `json:"_last_seen_by_usw"`
	LatestAssocTime  FlexInt  `json:"latest_assoc_time"`
	Mac              string   `fake:"{macaddress}"           json:"mac"`
	Name             string   `fake:"{randomstring:[client-1,client-2,client-3]}" json:"name"`
	Network          string   `json:"network"`
	NetworkID        string   `fake:"{uuid}"                 json:"network_id"`
	Noise            FlexInt  `json:"noise"`
	Note             string   `fake:"{sentence 20}"          json:"note"`
	Noted            FlexBool `json:"noted"`
	OsClass          FlexInt  `json:"os_class"`
	OsName           FlexInt  `json:"os_name"`
	Oui              string   `json:"oui"`
	PowersaveEnabled FlexBool `json:"powersave_enabled"`
	QosPolicyApplied FlexBool `json:"qos_policy_applied"`
	Radio            string   `json:"radio"`
	RadioDescription string   `json:"-"`
	RadioName        string   `json:"radio_name"`
	RadioProto       string   `json:"radio_proto"`
	RoamCount        FlexInt  `json:"roam_count"`
	Rssi             FlexInt  `json:"rssi"`
	RxBytes          FlexInt  `json:"rx_bytes"`
	RxBytesR         FlexInt  `json:"rx_bytes-r"`
	RxPackets        FlexInt  `json:"rx_packets"`
	RxRate           FlexInt  `json:"rx_rate"`
	Satisfaction     FlexInt  `json:"satisfaction,omitempty"`
	Signal           FlexInt  `json:"signal"`
	SiteID           string   `fake:"{uuid}"                 json:"site_id"`
	SiteName         string   `json:"-"`
	SourceName       string   `json:"-"`
	SwDepth          int      `json:"sw_depth"`
	SwMac            string   `fake:"{macaddress}"           json:"sw_mac"`
	SwName           string   `json:"-"`
	SwPort           FlexInt  `json:"sw_port"`
	TxBytes          FlexInt  `json:"tx_bytes"`
	TxBytesR         FlexInt  `json:"tx_bytes-r"`
	TxPackets        FlexInt  `json:"tx_packets"`
	TxPower          FlexInt  `json:"tx_power"`
	TxRate           FlexInt  `json:"tx_rate"`
	TxRetries        FlexInt  `json:"tx_retries"`
	Uptime           FlexInt  `json:"uptime"`
	UptimeByUAP      FlexInt  `json:"_uptime_by_uap"`
	UptimeByUGW      FlexInt  `json:"_uptime_by_ugw"`
	UptimeByUSW      FlexInt  `json:"_uptime_by_usw"`
	UseFixedIP       FlexBool `json:"use_fixedip"`
	UserGroupID      string   `fake:"{uuid}"                 json:"usergroup_id"`
	UserID           string   `fake:"{uuid}"                 json:"user_id"`
	Vlan             FlexInt  `json:"vlan"`
	WifiTxAttempts   FlexInt  `json:"wifi_tx_attempts"`
	WiredRxBytes     FlexInt  `json:"wired-rx_bytes"`
	WiredRxBytesR    FlexInt  `json:"wired-rx_bytes-r"`
	WiredRxPackets   FlexInt  `json:"wired-rx_packets"`
	WiredTxBytes     FlexInt  `json:"wired-tx_bytes"`
	WiredTxBytesR    FlexInt  `json:"wired-tx_bytes-r"`
	WiredTxPackets   FlexInt  `json:"wired-tx_packets"`
}

type ClientHistory struct {
	Blocked     FlexBool `json:"blocked"`
	Channel     FlexInt  `json:"channel"`
	DisplayName string   `json:"display_name"`
	Fingerprint struct {
		HasOverride FlexBool `json:"has_override"`
	} `json:"fingerprint"`
	FirstSeen                     FlexInt      `json:"first_seen"`
	Hostname                      string       `json:"hostname,omitempty"`
	ID                            string       `json:"id"`
	IsAllowedInVisualProgramming  FlexBool     `json:"is_allowed_in_visual_programming"`
	IsGuest                       FlexBool     `json:"is_guest"`
	IsMlo                         FlexBool     `json:"is_mlo"`
	IsWired                       FlexBool     `json:"is_wired"`
	LastIP                        string       `json:"last_ip"`
	LastRadio                     string       `json:"last_radio"`
	LastSeen                      FlexInt      `json:"last_seen"`
	LastUplinkMac                 string       `json:"last_uplink_mac"`
	LastUplinkName                string       `json:"last_uplink_name"`
	LocalDNSRecordEnabled         FlexBool     `json:"local_dns_record_enabled"`
	Mac                           string       `json:"mac"`
	Noted                         FlexBool     `json:"noted"`
	Oui                           string       `json:"oui"`
	SiteID                        string       `json:"site_id"`
	Status                        string       `json:"status"`
	Tags                          []FlexString `json:"tags"`
	Type                          string       `json:"type"`
	UnifiDevice                   FlexBool     `json:"unifi_device"`
	UplinkMac                     string       `json:"uplink_mac"`
	UseFixedip                    FlexBool     `json:"use_fixedip"`
	UserID                        string       `json:"user_id"`
	UsergroupID                   string       `json:"usergroup_id"`
	VirtualNetworkOverrideEnabled FlexBool     `json:"virtual_network_override_enabled"`
	WlanconfID                    string       `json:"wlanconf_id"`
	Name                          string       `json:"name,omitempty"`
	SiteName                      string       `json:"-"`
	SourceName                    string       `json:"-"`
}

type ClientHistoryOpts struct {
	OnlyNonBlocked      bool
	IncludeUnifiDevices bool
	WithinHours         int // 0 for all time
}

// NewClientHistoryOpts returns a new ClientHistoryOpts with values set to include all device history
func NewClientHistoryOpts() *ClientHistoryOpts {
	return &ClientHistoryOpts{
		OnlyNonBlocked:      false,
		IncludeUnifiDevices: true,
		WithinHours:         0,
	}
}

func (c *ClientHistoryOpts) SetOnlyNonBlocked(onlyNonBlocked bool) *ClientHistoryOpts {
	c.OnlyNonBlocked = onlyNonBlocked
	return c
}

func (c *ClientHistoryOpts) SetIncludeUnifiDevices(includeUnifiDevices bool) *ClientHistoryOpts {
	c.IncludeUnifiDevices = includeUnifiDevices
	return c
}
func (c *ClientHistoryOpts) SetWithinHours(withinHours int) *ClientHistoryOpts {
	c.WithinHours = withinHours
	return c
}
