package unifi

import (
	"fmt"
	"strings"
)

// GetUsers returns a response full of clients that connected to the UDM within the provided amount of time
// using the insight historical connection data set.
func (u *Unifi) GetUsers(sites []*Site, hours int) ([]*User, error) {
	data := make([]*User, 0)

	for _, site := range sites {
		var (
			response struct {
				Data []*User `json:"data"`
			}
			params = fmt.Sprintf(`{ "type": "all:", "conn": "all", "within":%d }`, hours)
		)

		u.DebugLog("Polling Controller, retrieving UniFi Users, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIAllUserPath, site.Name)
		if err := u.GetData(clientPath, &response, params); err != nil {
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

// User defines the metadata available for previously connected clients.
type User struct {
	Blocked             FlexBool `json:"blocked,omitempty"`
	DevIDOverride       FlexInt  `json:"dev_id_override,omitempty"`
	Duration            FlexInt  `json:"duration,omitempty"`
	FingerprintOverride FlexBool `json:"fingerprint_override,omitempty"`
	FirstSeen           FlexInt  `json:"first_seen,omitempty"`
	FixedIp             string   `fake:"{ipv4address}"                  json:"fixed_ip,omitempty"` //nolint:revive
	Hostname            string   `json:"hostname,omitempty"`
	ID                  string   `fake:"{uuid}"                         json:"_id"`
	IsGuest             bool     `json:"is_guest"`
	IsWired             bool     `json:"is_wired,omitempty"`
	LastSeen            FlexInt  `json:"last_seen,omitempty"`
	Mac                 string   `fake:"{macaddress}"                   json:"mac"`
	Name                string   `fake:"{animal}"                       json:"name,omitempty"`
	Note                string   `fake:"{buzzword}"                     json:"note,omitempty"`
	Noted               FlexBool `json:"noted,omitempty"`
	Oui                 string   `json:"oui,omitempty"`
	RxBytes             FlexInt  `json:"rx_bytes,omitempty"`
	RxPackets           FlexInt  `json:"rx_packets,omitempty"`
	SiteID              string   `fake:"{uuid}"                         json:"site_id"`
	SiteName            string   `json:"-"`
	SourceName          string   `json:"-"`
	TxBytes             FlexInt  `json:"tx_bytes,omitempty"`
	TxPackets           FlexInt  `json:"tx_packets,omitempty"`
	TxRetries           FlexInt  `json:"tx_retries,omitempty"`
	UseFixedIp          FlexBool `json:"use_fixedip,omitempty"` //nolint:revive
	UsergroupID         string   `json:"usergroup_id,omitempty"`
	WifiTxAttempts      FlexInt  `json:"wifi_tx_attempts,omitempty"`
}
