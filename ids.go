package unifi

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// IDS holds an Intrusion Prevention System Event.
type IDS struct {
	Archived              FlexBool  `json:"archived"`
	DestPort              int       `json:"dest_port,omitempty" fake:"{port}"`
	SrcPort               int       `json:"src_port,omitempty" fake:"{port}"`
	FlowID                int64     `json:"flow_id" fake:"{uuid}"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	InnerAlertGID         int64     `json:"inner_alert_gid"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	Time                  int64     `json:"time" fake:"{timestamp}"`
	Timestamp             int64     `json:"timestamp" fake:"{timestamp}"`
	Datetime              time.Time `json:"datetime"`
	AppProto              string    `json:"app_proto,omitempty"`
	Catname               string    `json:"catname"`
	DestIP                string    `json:"dest_ip" fake:"{ipv4address}"`
	DstMAC                string    `json:"dst_mac" fake:"{macaddress}"`
	DstIPASN              string    `json:"dstipASN" fake:"{address}"`
	DstIPCountry          string    `json:"dstipCountry" fake:"{country}"`
	EventType             string    `json:"event_type"`
	Host                  string    `json:"host"`
	ID                    string    `json:"_id" fake:"{uuid}"`
	InIface               string    `json:"in_iface"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	Key                   string    `json:"key" fake:"{uuid}"`
	Msg                   string    `json:"msg" fake:"{sentence:20}"`
	Proto                 string    `json:"proto"`
	SiteID                string    `json:"site_id" fake:"{uuid}"`
	SiteName              string    `json:"-"`
	SourceName            string    `json:"-"`
	SrcIP                 string    `json:"src_ip" fake:"{ipv4address}"`
	SrcIPASN              string    `json:"srcipASN" fake:"{address}"`
	SrcIPCountry          string    `json:"srcipCountry" fake:"{country}"`
	SrcMAC                string    `json:"src_mac" fake:"{macaddress}"`
	Subsystem             string    `json:"subsystem"`
	UniqueAlertID         string    `json:"unique_alertid" fake:"{uuid}"`
	USGIP                 string    `json:"usgip" fake:"{ipv4address}"`
	USGIPASN              string    `json:"usgipASN" fake:"{address}"`
	USGIPCountry          string    `json:"usgipCountry" fake:"{country}"`
	DestIPGeo             IPGeo     `json:"dstipGeo"`
	SourceIPGeo           IPGeo     `json:"srcipGeo"`
	USGIPGeo              IPGeo     `json:"usgipGeo"`
}

// GetIDS returns Intrusion Detection Systems events for a list of Sites.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetIDS(sites []*Site, timeRange ...time.Time) ([]*IDS, error) {
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
