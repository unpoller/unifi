package unifi

import (
	"encoding/json"
	"fmt"
)

// GetPortAnomalies returns port anomalies for all provided sites.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/ports/port-anomalies
// An empty slice is returned when no anomalies are detected (healthy network).
func (u *Unifi) GetPortAnomalies(sites []*Site) ([]*PortAnomaly, error) {
	anomalies := make([]*PortAnomaly, 0)

	for _, site := range sites {
		siteAnomalies, err := u.GetPortAnomaliesSite(site)
		if err != nil {
			return anomalies, err
		}

		anomalies = append(anomalies, siteAnomalies...)
	}

	return anomalies, nil
}

// GetPortAnomaliesSite returns port anomalies for a single site.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/ports/port-anomalies
// An empty slice is returned when no anomalies are detected (healthy network).
func (u *Unifi) GetPortAnomaliesSite(site *Site) ([]*PortAnomaly, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for Port Anomalies, site %s", site.SiteName)

	path := fmt.Sprintf(APIPortAnomaliesPath, site.Name)

	body, err := u.GetJSON(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch port anomalies for site %s: %w", site.SiteName, err)
	}

	var raw []*PortAnomaly
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse port anomalies for site %s: %w", site.SiteName, err)
	}

	anomalies := make([]*PortAnomaly, 0, len(raw))

	for _, a := range raw {
		if a == nil {
			continue
		}

		a.SiteName = site.SiteName
		a.SourceName = u.URL
		anomalies = append(anomalies, a)
	}

	return anomalies, nil
}

// PortAnomaly represents a port anomaly event from the UniFi controller.
// The API returns an empty array when the network is healthy (no anomalies).
type PortAnomaly struct {
	DeviceMAC  string  `fake:"{macaddress}" json:"device_mac"`
	PortIdx    FlexInt `fake:"{number:1,48}"  json:"port_idx"`
	AnomalyType string `fake:"{randomstring:[CRC_ERROR,COLLISION,PORT_FLAPPING,RX_ERROR,TX_ERROR,PACKET_DROP]}" json:"anomaly_type"`
	Count      FlexInt `fake:"{number:1,1000}" json:"count"`
	LastSeen   FlexInt `fake:"{unixtime}"      json:"last_seen"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}
