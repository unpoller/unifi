package unifi

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// IDS holds an Intrusion Prevention System Event.
type IDS struct { //nolint:revive
	AppProto              string     `json:"app_proto,omitempty"`
	Archived              FlexBool   `json:"archived"`
	Catname               FlexString `json:"catname"`
	Datetime              time.Time  `fake:"{recent_time}"            json:"datetime"`
	DestIP                string     `fake:"{ipv4address}"            json:"dest_ip"`
	DestIPGeo             IPGeo      `json:"dstipGeo"`
	DestPort              FlexInt    `json:"dest_port,omitempty"`
	DstIPASN              string     `fake:"{address}"                json:"dstipASN"`
	DstIPCountry          string     `fake:"{country}"                json:"dstipCountry"`
	DstMAC                string     `fake:"{macaddress}"             json:"dst_mac"`
	EventType             string     `json:"event_type"`
	FlowID                FlexInt    `json:"flow_id"`
	Host                  string     `json:"host"`
	ID                    string     `fake:"{uuid}"                   json:"_id"`
	InIface               string     `json:"in_iface"`
	InnerAlertAction      string     `json:"inner_alert_action"`
	InnerAlertCategory    string     `json:"inner_alert_category"`
	InnerAlertGID         FlexInt    `json:"inner_alert_gid"`
	InnerAlertRev         FlexInt    `json:"inner_alert_rev"`
	InnerAlertSeverity    FlexInt    `json:"inner_alert_severity"`
	InnerAlertSignature   string     `json:"inner_alert_signature"`
	InnerAlertSignatureID FlexInt    `json:"inner_alert_signature_id"`
	Key                   string     `fake:"{uuid}"                   json:"key"`
	Msg                   string     `fake:"{buzzword}"               json:"msg"`
	Proto                 string     `json:"proto"`
	SiteID                string     `fake:"{uuid}"                   json:"site_id"`
	SiteName              string     `json:"-"`
	SourceIPGeo           IPGeo      `json:"srcipGeo"`
	SourceName            string     `json:"-"`
	SrcIP                 string     `fake:"{ipv4address}"            json:"src_ip"`
	SrcIPASN              string     `fake:"{address}"                json:"srcipASN"`
	SrcIPCountry          string     `fake:"{country}"                json:"srcipCountry"`
	SrcMAC                string     `fake:"{macaddress}"             json:"src_mac"`
	SrcPort               FlexInt    `json:"src_port,omitempty"`
	Subsystem             string     `json:"subsystem"`
	Time                  FlexInt    `json:"time"`
	Timestamp             FlexInt    `json:"timestamp"`
	USGIP                 string     `fake:"{ipv4address}"            json:"usgip"`
	USGIPASN              string     `fake:"{address}"                json:"usgipASN"`
	USGIPCountry          string     `fake:"{country}"                json:"usgipCountry"`
	USGIPGeo              IPGeo      `json:"usgipGeo"`
	UniqueAlertID         string     `json:"unique_alertid"`
}

// GetIDS returns Intrusion Detection Systems events for a list of Sites.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetIDS(sites []*Site, timeRange ...time.Time) ([]*IDS, error) { //nolint:revive
	data := []*IDS{}

	for _, site := range sites {
		response, err := u.GetIDSSite(site, timeRange...)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetIDSSite retrieves the Intrusion Detection System Data for a single Site.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetIDSSite(site *Site, timeRange ...time.Time) ([]*IDS, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for IDS Events, site %s", site.SiteName)

	var (
		path = fmt.Sprintf(APIEventPathIDS, site.Name)
		ids  struct {
			Data idsList `json:"data"`
		}
	)

	if params, err := makeEventParams(timeRange...); err != nil {
		return ids.Data, err
	} else if err = u.GetData(path, &ids, params); err != nil {
		return ids.Data, err
	}

	for i := range ids.Data {
		// Add special SourceName value.
		ids.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		ids.Data[i].SiteName = site.SiteName
	}

	sort.Sort(ids.Data)

	return ids.Data, nil
}

func makeEventParams(timeRange ...time.Time) (string, error) {
	type eventReq struct {
		Start int64  `json:"start,omitempty"`
		End   int64  `json:"end,omitempty"`
		Limit int    `json:"_limit,omitempty"`
		Sort  string `json:"_sort"`
	}

	rp := eventReq{Limit: eventLimit, Sort: "-time"}

	switch len(timeRange) {
	case 0:
		rp.End = time.Now().Unix() * int64(time.Microsecond)
	case 1:
		rp.Start = timeRange[0].Unix() * int64(time.Microsecond)
		rp.End = time.Now().Unix() * int64(time.Microsecond)
	case 2: // nolint: gomnd
		rp.Start = timeRange[0].Unix() * int64(time.Microsecond)
		rp.End = timeRange[1].Unix() * int64(time.Microsecond)
	default:
		return "", ErrInvalidTimeRange
	}

	params, err := json.Marshal(&rp)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return string(params), nil
}

type idsList []*IDS

// Len satisfies sort.Interface.
func (e idsList) Len() int {
	return len(e)
}

// Swap satisfies sort.Interface.
func (e idsList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less satisfies sort.Interface. Sort our list by Datetime.
func (e idsList) Less(i, j int) bool {
	return e[i].Datetime.Before(e[j].Datetime)
}
