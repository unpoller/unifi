package unifi

import (
	"fmt"
	"strings"
)

var ErrDPIDataBug = fmt.Errorf("dpi data table contains more than 1 item; please open a bug report")

// GetSites returns a list of configured sites on the UniFi controller.
func (u *Unifi) GetSites() ([]*Site, error) {
	var response struct {
		Data []*Site `json:"data"`
	}

	if err := u.GetData(APISiteList, &response); err != nil {
		return nil, err
	}

	sites := []string{} // used for debug log only

	for i, d := range response.Data {
		// Add the unifi struct to the site.
		response.Data[i].controller = u
		// Add special SourceName value.
		response.Data[i].SourceName = u.URL
		// If the human name is missing (description), set it to the cryptic name.
		response.Data[i].Desc = strings.TrimSpace(pick(d.Desc, d.Name))
		// Add the custom site name to each site. used as a Grafana filter somewhere.
		response.Data[i].SiteName = d.Desc + " (" + d.Name + ")"
		sites = append(sites, d.Name) // used for debug log only
	}

	u.DebugLog("Found %d site(s): %s", len(sites), strings.Join(sites, ","))

	return response.Data, nil
}

// GetSiteDPI garners dpi data for sites.
func (u *Unifi) GetSiteDPI(sites []*Site) ([]*DPITable, error) {
	data := []*DPITable{}

	for _, site := range sites {
		u.DebugLog("Polling Controller, retreiving Site DPI data, site %s", site.SiteName)

		var response struct {
			Data []*DPITable `json:"data"`
		}

		siteDPIpath := fmt.Sprintf(APISiteDPI, site.Name)
		if err := u.GetData(siteDPIpath, &response, `{"type":"by_app"}`); err != nil {
			return nil, err
		}

		if l := len(response.Data); l > 1 {
			return nil, ErrDPIDataBug
		} else if l == 0 {
			u.DebugLog("Site DPI data missing! Is DPI enabled in UniFi controller? Site %s", site.SiteName)

			continue
		}

		response.Data[0].SourceName = site.SourceName
		response.Data[0].SiteName = site.SiteName
		data = append(data, response.Data[0])
	}

	return data, nil
}

// Site represents a site's data.
type Site struct {
	AttrHiddenID string   `json:"attr_hidden_id"`
	AttrNoDelete FlexBool `json:"attr_no_delete"`
	controller   *Unifi
	Desc         string `fake:"{buzzword}" json:"desc"`
	Health       []struct {
		Drops         FlexInt  `json:"drops,omitempty"`
		Gateways      []string `fakesize:"5"             json:"gateways,omitempty"`
		GwMac         string   `fake:"{macaddress}"      json:"gw_mac,omitempty"`
		GwName        string   `json:"gw_name,omitempty"`
		GwSystemStats struct {
			CPU    FlexInt `json:"cpu"`
			Mem    FlexInt `json:"mem"`
			Uptime FlexInt `json:"uptime"`
		} `json:"gw_system-stats,omitempty"`
		GwVersion             string   `fake:"{appversion}"                       json:"gw_version,omitempty"`
		LanIP                 string   `json:"lan_ip,omitempty"`
		Latency               FlexInt  `json:"latency,omitempty"`
		Nameservers           []string `fakesize:"5"                              json:"nameservers,omitempty"`
		Netmask               string   `json:"netmask,omitempty"`
		NumAdopted            FlexInt  `json:"num_adopted,omitempty"`
		NumAp                 FlexInt  `json:"num_ap,omitempty"`
		NumDisabled           FlexInt  `json:"num_disabled,omitempty"`
		NumDisconnected       FlexInt  `json:"num_disconnected,omitempty"`
		NumGuest              FlexInt  `json:"num_guest,omitempty"`
		NumGw                 FlexInt  `json:"num_gw,omitempty"`
		NumIot                FlexInt  `json:"num_iot,omitempty"`
		NumPending            FlexInt  `json:"num_pending,omitempty"`
		NumSta                FlexInt  `json:"num_sta,omitempty"`
		NumSw                 FlexInt  `json:"num_sw,omitempty"`
		NumUser               FlexInt  `json:"num_user,omitempty"`
		RemoteUserEnabled     FlexBool `json:"remote_user_enabled,omitempty"`
		RemoteUserNumActive   FlexInt  `json:"remote_user_num_active,omitempty"`
		RemoteUserNumInactive FlexInt  `json:"remote_user_num_inactive,omitempty"`
		RemoteUserRxBytes     FlexInt  `json:"remote_user_rx_bytes,omitempty"`
		RemoteUserRxPackets   FlexInt  `json:"remote_user_rx_packets,omitempty"`
		RemoteUserTxBytes     FlexInt  `json:"remote_user_tx_bytes,omitempty"`
		RemoteUserTxPackets   FlexInt  `json:"remote_user_tx_packets,omitempty"`
		RxBytesR              FlexInt  `json:"rx_bytes-r,omitempty"`
		SiteToSiteEnabled     FlexBool `json:"site_to_site_enabled,omitempty"`
		SpeedtestLastrun      FlexInt  `json:"speedtest_lastrun,omitempty"`
		SpeedtestPing         FlexInt  `json:"speedtest_ping,omitempty"`
		SpeedtestStatus       string   `json:"speedtest_status,omitempty"`
		Status                string   `json:"status"`
		Subsystem             string   `json:"subsystem"`
		TxBytesR              FlexInt  `json:"tx_bytes-r,omitempty"`
		Uptime                FlexInt  `json:"uptime,omitempty"`
		WanIP                 string   `fake:"{ipv4address}"                      json:"wan_ip,omitempty"`
		XputDown              FlexInt  `json:"xput_down,omitempty"`
		XputUp                FlexInt  `json:"xput_up,omitempty"`
	} `fakesize:"5"                          json:"health"`
	ID           string  `fake:"{uuid}"                         json:"_id"`
	Name         string  `fake:"{randomstring:[site-1,site-2]}" json:"name"`
	NumNewAlarms FlexInt `json:"num_new_alarms"`
	SiteName     string  `json:"-"`
	SourceName   string  `json:"-"`
}
