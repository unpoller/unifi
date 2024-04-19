package unifi

import (
	"fmt"
	"sort"
	"time"
)

type Alarm struct {
	AppProto              string     `json:"app_proto,omitempty"`
	Archived              FlexBool   `json:"archived"`
	Catname               FlexString `json:"catname"`
	Datetime              time.Time  `fake:"{recent_time}"                             json:"datetime"`
	DestIP                string     `fake:"{ipv4address}"                             json:"dest_ip"`
	DestIPGeo             IPGeo      `json:"dstipGeo"`
	DestPort              int        `fake:"{port}"                                    json:"dest_port"`
	DstIPASN              string     `json:"dstipASN,omitempty"`
	DstIPCountry          string     `json:"dstipCountry,omitempty"`
	DstMAC                string     `fake:"{macaddress}"                              json:"dst_mac"`
	EventType             string     `json:"event_type"`
	FlowID                int64      `json:"flow_id"`
	HandledAdminID        string     `json:"handled_admin_id,omitempty"`
	HandledTime           time.Time  `json:"handled_time,omitempty"`
	Host                  string     `json:"host"`
	ID                    string     `fake:"{uuid}"                                    json:"_id"`
	InIface               string     `fake:"{randomstring:[eth0,eth1,lan1,wan1,wan2]}" json:"in_iface"`
	InnerAlertAction      string     `json:"inner_alert_action"`
	InnerAlertCategory    string     `json:"inner_alert_category"`
	InnerAlertGID         int64      `json:"inner_alert_gid"`
	InnerAlertRev         int64      `json:"inner_alert_rev"`
	InnerAlertSeverity    int64      `json:"inner_alert_severity"`
	InnerAlertSignature   string     `json:"inner_alert_signature"`
	InnerAlertSignatureID int64      `json:"inner_alert_signature_id"`
	Key                   string     `json:"key"`
	Msg                   string     `fake:"{sentence:5}"                              json:"msg"`
	Proto                 string     `json:"proto"`
	SiteID                string     `fake:"{uuid}"                                    json:"site_id"`
	SiteName              string     `json:"-"`
	SourceIPGeo           IPGeo      `json:"usgipGeo"`
	SourceName            string     `json:"-"`
	SrcIP                 string     `fake:"{ipv4address}"                             json:"src_ip"`
	SrcIPASN              string     `json:"srcipASN,omitempty"`
	SrcIPCountry          string     `json:"srcipCountry,omitempty"`
	SrcMAC                string     `fake:"{macaddress}"                              json:"src_mac"`
	SrcPort               int        `fake:"{port}"                                    json:"src_port"`
	Subsystem             string     `json:"subsystem"`
	Time                  int64      `fake:"{timestamp}"                               json:"time"`
	Timestamp             int64      `fake:"{timestamp}"                               json:"timestamp"`
	TxID                  FlexInt    `json:"tx_id,omitempty"`
	USGIP                 string     `fake:"{ipv4address}"                             json:"usgip"`
	USGIPASN              string     `json:"usgipASN"`
	USGIPCountry          string     `json:"usgipCountry"`
	USGIPGeo              IPGeo      `json:"srcipGeo,omitempty"`
	UniqueAlertID         string     `json:"unique_alertid"`
}

// GetAlarms returns Alarms for a list of Sites.
func (u *Unifi) GetAlarms(sites []*Site) ([]*Alarm, error) {
	data := []*Alarm{}

	for _, site := range sites {
		response, err := u.GetAlarmsSite(site)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetAlarmsSite retreives the Alarms for a single Site.
func (u *Unifi) GetAlarmsSite(site *Site) ([]*Alarm, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for Alarms, site %s", site.SiteName)

	var (
		path   = fmt.Sprintf(APIEventPathAlarms, site.Name)
		alarms struct {
			Data alarms `json:"data"`
		}
	)

	if err := u.GetData(path, &alarms, ""); err != nil {
		return alarms.Data, err
	}

	for i := range alarms.Data {
		// Add special SourceName value.
		alarms.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		alarms.Data[i].SiteName = site.SiteName
	}

	sort.Sort(alarms.Data)

	return alarms.Data, nil
}

// alarms satisfies the sort.Sort Interface.
type alarms []*Alarm

// Len satisfies sort.Interface.
func (a alarms) Len() int {
	return len(a)
}

// Swap satisfies sort.Interface.
func (a alarms) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less satisfies sort.Interface. Sort our list by Datetime.
func (a alarms) Less(i, j int) bool {
	return a[i].Datetime.Before(a[j].Datetime)
}
