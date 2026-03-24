package unifi

import (
	"encoding/json"
	"fmt"
)

// GetTopology returns network topology data for all provided sites.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/topology
func (u *Unifi) GetTopology(sites []*Site) ([]*Topology, error) {
	topologies := make([]*Topology, 0)

	for _, site := range sites {
		path := fmt.Sprintf(APITopologyPath, site.Name)

		body, err := u.GetJSON(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch topology for site %s: %w", site.SiteName, err)
		}

		var topo Topology
		if err := json.Unmarshal(body, &topo); err != nil {
			return nil, fmt.Errorf("failed to parse topology for site %s: %w", site.SiteName, err)
		}

		topo.SiteName = site.SiteName
		topo.SourceName = u.URL
		topologies = append(topologies, &topo)
	}

	return topologies, nil
}

// TopologyVertex represents a node in the network topology (device or client).
type TopologyVertex struct {
	AllowedInVisualProgramming bool            `json:"allowedInVisualProgramming"`
	Default                    bool            `json:"default"`
	Mac                        string          `json:"mac"`
	Model                      string          `json:"model"`
	Name                       string          `json:"name"`
	State                      FlexInt         `json:"state"`
	Type                       string          `json:"type"` // DEVICE or CLIENT
	UnifiDevice                bool            `json:"unifiDevice"`
	WifiRadios                 []TopologyRadio `json:"wifiRadios"`
}

// TopologyRadio represents a WiFi radio on a device vertex.
type TopologyRadio struct {
	Channel   FlexInt `json:"channel"`
	Protocol  string  `json:"protocol"`  // ax, ac, n, g
	RadioBand string  `json:"radioBand"` // ng (2.4GHz), na (5GHz), 6e (6GHz)
}

// TopologyEdge represents a connection between two vertices in the topology.
type TopologyEdge struct {
	Channel            FlexInt `json:"channel"`
	DownlinkMac        string  `json:"downlinkMac"`
	DownlinkPortNumber FlexInt `json:"downlinkPortNumber"`
	Duplex             string  `json:"duplex"` // FULL_DUPLEX, HALF_DUPLEX
	Essid              string  `json:"essid"`
	ExperienceScore    FlexInt `json:"experienceScore"`
	NetworkID          string  `json:"networkId"`
	Protocol           string  `json:"protocol"`  // ax, ac, n, g
	RadioBand          string  `json:"radioBand"` // ng, na, 6e
	RateMbps           FlexInt `json:"rateMbps"`
	Type               string  `json:"type"` // WIRED or WIRELESS
	UplinkMac          string  `json:"uplinkMac"`
	UplinkPortNumber   FlexInt `json:"uplinkPortNumber"`
}

// Topology represents the full network topology for a site.
type Topology struct {
	Edges            []TopologyEdge   `json:"edges"`
	HasUnknownSwitch bool             `json:"has_unknown_switch"`
	Vertices         []TopologyVertex `json:"vertices"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}
