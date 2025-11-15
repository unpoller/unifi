package unifi

import (
	"fmt"
)

func (u *Unifi) GetClientTraffic(sites []*Site, epochMillisTimePeriod *EpochMillisTimePeriod, includeUnidentified bool) ([]*ClientUsageByApp, error) {
	_, err := epochMillisTimePeriod.isValid()
	if err != nil {
		return nil, err
	}

	data := make([]*ClientUsageByApp, 0)

	for _, site := range sites {
		var response struct {
			ClientUsageByApp []*ClientUsageByApp `json:"client_usage_by_app"`
		}

		trafficSite := &TrafficSite{
			SiteID:     site.ID,
			SiteName:   site.SiteName,
			SourceName: site.SourceName,
		}

		u.DebugLog("Polling Controller, retrieving UniFi Client Traffic, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIClientTrafficPath,
			site.Name,
			epochMillisTimePeriod.StartEpochMillis,
			epochMillisTimePeriod.EndEpochMillis,
			includeUnidentified)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, elem := range response.ClientUsageByApp {
			elem.TrafficSite = trafficSite
		}

		data = append(data, response.ClientUsageByApp...)
	}

	return data, nil
}

func (u *Unifi) GetClientTrafficByMac(site *Site, epochMillisTimePeriod *EpochMillisTimePeriod, includeUnidentified bool, macs ...string) ([]*ClientUsageByApp, error) {
	_, err := epochMillisTimePeriod.isValid()
	if err != nil {
		return nil, err
	}

	data := make([]*ClientUsageByApp, 0)

	var response struct {
		ClientUsageByApp []*ClientUsageByApp `json:"client_usage_by_app"`
	}

	trafficSite := &TrafficSite{
		SiteID:     site.ID,
		SiteName:   site.SiteName,
		SourceName: site.SourceName,
	}

	for _, mac := range macs {
		u.DebugLog("Polling Controller, retrieving UniFi Client Traffic By MAC address, site %s and mac %s", site.SiteName, mac)

		clientPath := fmt.Sprintf(APIClientTrafficByMacPath,
			site.Name,
			mac,
			epochMillisTimePeriod.StartEpochMillis,
			epochMillisTimePeriod.EndEpochMillis,
			includeUnidentified,
			mac)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, elem := range response.ClientUsageByApp {
			elem.TrafficSite = trafficSite
		}

		data = append(data, response.ClientUsageByApp...)
	}

	return data, nil
}

func (u *Unifi) GetCountryTraffic(sites []*Site, epochMillisTimePeriod *EpochMillisTimePeriod) ([]*UsageByCountry, error) {
	_, err := epochMillisTimePeriod.isValid()
	if err != nil {
		return nil, err
	}

	data := make([]*UsageByCountry, 0)

	for _, site := range sites {
		var response struct {
			UsageByCountry []*UsageByCountry `json:"usage_by_country"`
		}

		trafficSite := &TrafficSite{
			SiteID:     site.ID,
			SiteName:   site.SiteName,
			SourceName: site.SourceName,
		}

		u.DebugLog("Polling Controller, retrieving UniFi Country Traffic, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APICountryTrafficPath,
			site.Name,
			epochMillisTimePeriod.StartEpochMillis,
			epochMillisTimePeriod.EndEpochMillis)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, elem := range response.UsageByCountry {
			elem.TrafficSite = trafficSite
		}

		data = append(data, response.UsageByCountry...)
	}

	return data, nil
}

type ClientUsageByApp struct {
	TrafficSite *TrafficSite `json:"site"`
	Client      ClientInfo   `json:"client"`
	UsageByApp  []UsageByApp `json:"usage_by_app"`
}

// ClientInfo contains information about the network client
type ClientInfo struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Hostname    string      `fake:"{domainname}"          json:"hostname"`
	IsWired     bool        `json:"is_wired"`
	Mac         string      `fake:"{macaddress}"          json:"mac"`
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
	TrafficSite      *TrafficSite `json:"site"`
	BytesReceived    int64        `json:"bytes_received"`
	BytesTransmitted int64        `json:"bytes_transmitted"`
	Country          string       `fake:"{countryabr}"      json:"country"`
	TotalBytes       int64        `json:"total_bytes"`
}

type TrafficSite struct {
	SiteID     string `fake:"{uuid}" json:"site_id"`
	SiteName   string `json:"name"`
	SourceName string `json:"source"`
}

// Parameters

type EpochMillisTimePeriod struct {
	StartEpochMillis int64
	EndEpochMillis   int64
}

func (p *EpochMillisTimePeriod) isValid() (bool, error) {
	if p.StartEpochMillis < p.EndEpochMillis {
		return true, nil
	}

	return false, fmt.Errorf("start must be before end (%s)", p.String())
}

func (p *EpochMillisTimePeriod) String() string {
	return fmt.Sprintf("%d to %d", p.StartEpochMillis, p.EndEpochMillis)
}
