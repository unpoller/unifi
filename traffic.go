package unifi

import (
	"fmt"
	"time"
)

func (u *Unifi) GetClientTraffic(sites []*Site, params *ClientTrafficParameters) ([]*ClientUsageByApp, error) {

	if !params.Period.isValid() {
		return nil, fmt.Errorf("start must be before end (%s)", params.Period.String())
	}

	data := make([]*ClientUsageByApp, 0)

	for _, site := range sites {

		var response struct {
			ClientUsageByApp []*ClientUsageByApp `json:"client_usage_by_app"`
		}

		u.DebugLog("Polling Controller, retrieving UniFi Traffic, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIClientTrafficPath,
			site.Name,
			params.Period.Start.UnixMilli(),
			params.Period.End.UnixMilli(),
			params.IncludeUnidentified)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, elem := range response.ClientUsageByApp {
			elem.SourceName = site.SourceName
			elem.SiteID = site.ID
			elem.SiteName = site.SiteName
		}

		data = append(data, response.ClientUsageByApp...)
	}

	return data, nil
}

func (u *Unifi) GetClientTrafficByMac(site *Site, mac string, params *ClientTrafficParameters) (*ClientUsageByApp, error) {

	if !params.Period.isValid() {
		return nil, fmt.Errorf("start must be before end (%s)", params.Period.String())
	}

	var response struct {
		ClientUsageByApp []*ClientUsageByApp `json:"client_usage_by_app"`
	}

	u.DebugLog("Polling Controller, retrieving UniFi Traffic, site %s and mac %s", site.SiteName, mac)

	clientPath := fmt.Sprintf(APIClientTrafficByMacPath,
		site.Name,
		mac,
		params.Period.Start.UnixMilli(),
		params.Period.End.UnixMilli(),
		params.IncludeUnidentified,
		mac)
	if err := u.GetData(clientPath, &response); err != nil {
		return nil, err
	}
	if len(response.ClientUsageByApp) == 1 {
		response.ClientUsageByApp[0].SourceName = site.SourceName
		response.ClientUsageByApp[0].SiteID = site.ID
		response.ClientUsageByApp[0].SiteName = site.SiteName
		return response.ClientUsageByApp[0], nil
	}
	return nil, fmt.Errorf("no traffic found at site %s for mac %s", site.SiteName, mac)
}

func (u *Unifi) GetCountryTraffic(sites []*Site, params *TimePeriod) ([]*UsageByCountry, error) {

	if !params.isValid() {
		return nil, fmt.Errorf("start must be before end (%s)", params.String())
	}

	data := make([]*UsageByCountry, 0)

	for _, site := range sites {

		var response struct {
			UsageByCountry []*UsageByCountry `json:"usage_by_country"`
		}

		u.DebugLog("Polling Controller, retrieving UniFi Traffic, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APICountryTrafficPath,
			site.Name,
			params.Start.UnixMilli(),
			params.End.UnixMilli())
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, elem := range response.UsageByCountry {
			elem.SourceName = site.SourceName
			elem.SiteID = site.ID
			elem.SiteName = site.SiteName
		}

		data = append(data, response.UsageByCountry...)
	}

	return data, nil
}

type ClientUsageByApp struct {
	SiteID     string       `json:"-"`
	SiteName   string       `json:"-"`
	SourceName string       `json:"-"`
	Client     ClientInfo   `json:"client"`
	UsageByApp []UsageByApp `json:"usage_by_app"`
}

// ClientInfo contains information about the network client
type ClientInfo struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Hostname    string      `json:"hostname"`
	IsWired     bool        `json:"is_wired"`
	Mac         string      `json:"mac"`
	Name        string      `json:"name"`
	Oui         string      `json:"oui"`
	WlanconfID  string      `json:"wlanconf_id,omitempty"`
}

// Fingerprint contains device fingerprinting information
type Fingerprint struct {
	ComputedDevID  int  `json:"computed_dev_id,omitempty"`
	ComputedEngine int  `json:"computed_engine,omitempty"`
	Confidence     int  `json:"confidence,omitempty"`
	DevCat         int  `json:"dev_cat,omitempty"`
	DevFamily      int  `json:"dev_family,omitempty"`
	DevID          int  `json:"dev_id,omitempty"`
	DevIDOverride  int  `json:"dev_id_override,omitempty"`
	DevVendor      int  `json:"dev_vendor,omitempty"`
	HasOverride    bool `json:"has_override"`
	OsClass        int  `json:"os_class,omitempty"`
	OsName         int  `json:"os_name,omitempty"`
}

// UsageByApp contains application usage statistics
type UsageByApp struct {
	ActivitySeconds  int   `json:"activity_seconds,omitempty"`
	Application      int   `json:"application"`
	BytesReceived    int64 `json:"bytes_received"`
	BytesTransmitted int64 `json:"bytes_transmitted"`
	Category         int   `json:"category"`
	TotalBytes       int64 `json:"total_bytes"`
}

type UsageByCountry struct {
	SiteID           string `json:"-"`
	SiteName         string `json:"-"`
	SourceName       string `json:"-"`
	BytesReceived    int64  `json:"bytes_received"`
	BytesTransmitted int64  `json:"bytes_transmitted"`
	Country          string
	TotalBytes       int64
}

type TimePeriod struct {
	Start time.Time
	End   time.Time
}

func (p *TimePeriod) isValid() bool {
	return p.Start.Before(p.End)
}

func (p *TimePeriod) String() string {
	return p.Start.String() + " to " + p.End.String()
}

type ClientTrafficParameters struct {
	Period              TimePeriod
	IncludeUnidentified bool
}
