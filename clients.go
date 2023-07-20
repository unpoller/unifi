package unifi

import (
	"fmt"
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

// Client defines all the data a connected-network client contains.
type Client struct {
	SourceName       string   `json:"-"`
	Anomalies        FlexInt  `json:"anomalies,omitempty"`
	ApMac            string   `json:"ap_mac" fake:"{macaddress}"`
	ApName           string   `json:"-"`
	AssocTime        FlexInt  `json:"assoc_time"`
	Blocked          bool     `json:"blocked,omitempty"`
	Bssid            string   `json:"bssid" fake:"{macaddress}"`
	BytesR           FlexInt  `json:"bytes-r"`
	Ccq              FlexInt  `json:"ccq"`
	Channel          FlexInt  `json:"channel"`
	DevCat           FlexInt  `json:"dev_cat"`
	DevFamily        FlexInt  `json:"dev_family"`
	DevID            FlexInt  `json:"dev_id"`
	DevVendor        FlexInt  `json:"dev_vendor,omitempty"`
	DhcpendTime      FlexInt  `json:"dhcpend_time,omitempty"`
	Satisfaction     FlexInt  `json:"satisfaction,omitempty"`
	Essid            string   `json:"essid" fake:"{macaddress}"`
	FirstSeen        FlexInt  `json:"first_seen"`
	FixedIP          string   `json:"fixed_ip" fake:"{ipv4address}"`
	GwMac            string   `json:"gw_mac" fake:"{macaddress}"`
	GwName           string   `json:"-"`
	Hostname         string   `json:"hostname"`
	ID               string   `json:"_id" fake:"{uuid}"`
	IP               string   `json:"ip" fake:"{ipv4address}"`
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
	Mac              string   `json:"mac" fake:"{macaddress}"`
	Name             string   `json:"name"`
	Network          string   `json:"network"`
	NetworkID        string   `json:"network_id" fake:"{uuid}"`
	Noise            FlexInt  `json:"noise"`
	Note             string   `json:"note" fake:"{sentence 20}"`
	Noted            FlexBool `json:"noted"`
	OsClass          FlexInt  `json:"os_class"`
	OsName           FlexInt  `json:"os_name"`
	Oui              string   `json:"oui"`
	PowersaveEnabled FlexBool `json:"powersave_enabled"`
	QosPolicyApplied FlexBool `json:"qos_policy_applied"`
	Radio            string   `json:"radio"`
	RadioName        string   `json:"radio_name"`
	RadioProto       string   `json:"radio_proto"`
	RadioDescription string   `json:"-"`
	RoamCount        FlexInt  `json:"roam_count"`
	Rssi             FlexInt  `json:"rssi"`
	RxBytes          FlexInt  `json:"rx_bytes"`
	RxBytesR         FlexInt  `json:"rx_bytes-r"`
	RxPackets        FlexInt  `json:"rx_packets"`
	RxRate           FlexInt  `json:"rx_rate"`
	Signal           FlexInt  `json:"signal"`
	SiteID           string   `json:"site_id" fake:"{uuid}"`
	SiteName         string   `json:"-"`
	SwDepth          int      `json:"sw_depth"`
	SwMac            string   `json:"sw_mac" fake:"{macaddress}"`
	SwName           string   `json:"-"`
	SwPort           FlexInt  `json:"sw_port"`
	TxBytes          FlexInt  `json:"tx_bytes"`
	TxBytesR         FlexInt  `json:"tx_bytes-r"`
	TxPackets        FlexInt  `json:"tx_packets"`
	TxRetries        FlexInt  `json:"tx_retries"`
	TxPower          FlexInt  `json:"tx_power"`
	TxRate           FlexInt  `json:"tx_rate"`
	Uptime           FlexInt  `json:"uptime"`
	UptimeByUAP      FlexInt  `json:"_uptime_by_uap"`
	UptimeByUGW      FlexInt  `json:"_uptime_by_ugw"`
	UptimeByUSW      FlexInt  `json:"_uptime_by_usw"`
	UseFixedIP       FlexBool `json:"use_fixedip"`
	UserGroupID      string   `json:"usergroup_id" fake:"{uuid}"`
	UserID           string   `json:"user_id" fake:"{uuid}"`
	Vlan             FlexInt  `json:"vlan"`
	WifiTxAttempts   FlexInt  `json:"wifi_tx_attempts"`
	WiredRxBytes     FlexInt  `json:"wired-rx_bytes"`
	WiredRxBytesR    FlexInt  `json:"wired-rx_bytes-r"`
	WiredRxPackets   FlexInt  `json:"wired-rx_packets"`
	WiredTxBytes     FlexInt  `json:"wired-tx_bytes"`
	WiredTxBytesR    FlexInt  `json:"wired-tx_bytes-r"`
	WiredTxPackets   FlexInt  `json:"wired-tx_packets"`
}
