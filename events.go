package unifi

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

var (
	ErrNoSiteProvided   = fmt.Errorf("site must not be nil or empty")
	ErrInvalidTimeRange = fmt.Errorf("only 0, 1 or 2 times may be provided to timeRange")
)

const (
	eventLimit = 50000
)

// GetEvents returns a response full of UniFi Events for the last 1 hour from multiple sites.
func (u *Unifi) GetEvents(sites []*Site, hours time.Duration) ([]*Event, error) {
	data := make([]*Event, 0)

	for _, site := range sites {
		response, err := u.GetSiteEvents(site, hours)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetSiteEvents retrieves the last 1 hour's worth of events from a single site.
func (u *Unifi) GetSiteEvents(site *Site, hours time.Duration) ([]*Event, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	if hours < time.Hour {
		hours = time.Hour
	}

	u.DebugLog("Polling Controller, retreiving UniFi Events, site %s", site.SiteName)

	var (
		path   = fmt.Sprintf(APIEventPath, site.Name)
		params = fmt.Sprintf(`{"_limit":%d,"within":%d,"_sort":"-time"}`,
			eventLimit, int(hours.Round(time.Hour).Hours()))
		event struct {
			Data events `json:"data"`
		}
	)

	if err := u.GetData(path, &event, params); err != nil {
		return event.Data, err
	}

	for i := range event.Data {
		// Add special SourceName value.
		event.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		event.Data[i].SiteName = site.SiteName
	}

	sort.Sort(event.Data)

	return event.Data, nil
}

// Event describes a UniFi Event.
// API Path: /api/s/default/stat/event.
type Event struct {
	Admin                 string     `json:"admin"`
	Ap                    string     `json:"ap"`
	ApFrom                string     `json:"ap_from"`
	ApName                string     `json:"ap_name"`
	ApTo                  string     `json:"ap_to"`
	AppProto              string     `json:"app_proto"`
	Bytes                 FlexInt    `json:"bytes"`
	Catname               FlexString `json:"catname"`
	Channel               FlexInt    `json:"channel"`
	ChannelFrom           FlexInt    `json:"channel_from"`
	ChannelTo             FlexInt    `json:"channel_to"`
	Datetime              time.Time  `fake:"{recent_time}"            json:"datetime"`
	DestIP                string     `fake:"{ipv4address}"            json:"dest_ip"`
	DestIPGeo             IPGeo      `json:"dstipGeo"`
	DestPort              int        `fake:"{port}"                   json:"dest_port"`
	DstMAC                string     `fake:"{macaddress}"             json:"dst_mac"`
	Duration              FlexInt    `json:"duration"`
	EventType             string     `json:"event_type"`
	FlowID                FlexInt    `json:"flow_id"`
	Guest                 string     `json:"guest"`
	Gw                    string     `json:"gw"`
	GwName                string     `json:"gw_name"`
	Host                  string     `json:"host"`
	Hostname              string     `json:"hostname"`
	ID                    string     `fake:"{uuid}"                   json:"_id"`
	IP                    string     `fake:"{ipv4address}"            json:"ip"`
	InIface               string     `json:"in_iface"`
	InnerAlertAction      string     `json:"inner_alert_action"`
	InnerAlertCategory    string     `json:"inner_alert_category"`
	InnerAlertGID         FlexInt    `json:"inner_alert_gid"`
	InnerAlertRev         FlexInt    `json:"inner_alert_rev"`
	InnerAlertSeverity    FlexInt    `json:"inner_alert_severity"`
	InnerAlertSignature   string     `json:"inner_alert_signature"`
	InnerAlertSignatureID FlexInt    `json:"inner_alert_signature_id"`
	IsAdmin               FlexBool   `json:"is_admin"`
	Key                   string     `fake:"{uuid}"                   json:"key"`
	Msg                   string     `fake:"{buzzword}"               json:"msg"`
	Network               string     `json:"network"`
	Proto                 string     `json:"proto"`
	Radio                 string     `json:"radio"`
	RadioFrom             string     `json:"radio_from"`
	RadioTo               string     `json:"radio_to"`
	SSID                  string     `fake:"{macaddress}"             json:"ssid"`
	SiteID                string     `fake:"{}"                       json:"site_id"`
	SiteName              string     `json:"-"`
	SourceIPGeo           IPGeo      `json:"srcipGeo"`
	SourceName            string     `json:"-"`
	SrcIP                 string     `fake:"{ipv4address}"            json:"src_ip"`
	SrcIPASN              string     `fake:"{address}"                json:"srcipASN"`
	SrcIPCountry          string     `fake:"{country}"                json:"srcipCountry"`
	SrcMAC                string     `fake:"{macaddress}"             json:"src_mac"`
	SrcPort               int        `fake:"{port}"                   json:"src_port"`
	Subsystem             string     `json:"subsystem"`
	Sw                    string     `json:"sw"`
	SwName                string     `json:"sw_name"`
	Time                  int64      `fake:"{timestamp}"              json:"time"`
	Timestamp             int64      `fake:"{timestamp}"              json:"timestamp"`
	USGIP                 string     `fake:"{ipv4address}"            json:"usgip"`
	USGIPASN              string     `fake:"{address}"                json:"usgipASN"`
	USGIPCountry          string     `fake:"{country}"                json:"usgipCountry"`
	USGIPGeo              IPGeo      `json:"usgipGeo"`
	UniqueAlertID         string     `json:"unique_alertid"`
	User                  string     `json:"user"`
}

// IPGeo is part of the UniFi Event data. Each event may have up to three of these.
// One for source, one for dest and one for the USG location.
type IPGeo struct {
	Asn           int64   `json:"asn"`
	City          string  `fake:"{city}"         json:"city"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `fake:"{countryabr}"   json:"country_code"`
	CountryName   string  `fake:"{country}"      json:"country_name"`
	Latitude      float64 `fake:"{latitude}"     json:"latitude"`
	Longitude     float64 `fake:"{longitude}"    json:"longitude"`
	Organization  string  `fake:"{company}"      json:"organization"`
}

// Events satisfied the sort.Interface.
type events []*Event

// Len satisfies sort.Interface.
func (e events) Len() int {
	return len(e)
}

// Swap satisfies sort.Interface.
func (e events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less satisfies sort.Interface. Sort our list by date/time.
func (e events) Less(i, j int) bool {
	return e[i].Datetime.Before(e[j].Datetime)
}

// UnmarshalJSON is required because sometimes the unifi api returns
// an empty array instead of a struct filled with data.
func (v *IPGeo) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		return nil // it's empty
	}

	g := struct {
		Asn           int64   `json:"asn"`
		City          string  `json:"city"`
		ContinentCode string  `json:"continent_code"`
		CountryCode   string  `json:"country_code"`
		CountryName   string  `json:"country_name"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Organization  string  `json:"organization"`
	}{}

	err := json.Unmarshal(data, &g)
	v.Asn = g.Asn
	v.Latitude = g.Latitude
	v.Longitude = g.Longitude
	v.City = g.City
	v.ContinentCode = g.ContinentCode
	v.CountryCode = g.CountryCode
	v.CountryName = g.CountryName
	v.Organization = g.Organization

	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}
