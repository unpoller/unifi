package unifi

import (
	"encoding/json"
	"fmt"
)

// GetMagicSiteToSiteVPN returns Site Magic site-to-site VPN mesh configurations for all provided sites.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/magicsitetositevpn/configs
func (u *Unifi) GetMagicSiteToSiteVPN(sites []*Site) ([]*MagicSiteToSiteVPN, error) {
	meshes := make([]*MagicSiteToSiteVPN, 0)

	for _, site := range sites {
		siteMeshes, err := u.GetMagicSiteToSiteVPNSite(site)
		if err != nil {
			return meshes, err
		}

		meshes = append(meshes, siteMeshes...)
	}

	return meshes, nil
}

// GetMagicSiteToSiteVPNSite returns Site Magic site-to-site VPN mesh configurations for a single site.
// Uses the v2 API endpoint: GET /proxy/network/v2/api/site/{site}/magicsitetositevpn/configs
func (u *Unifi) GetMagicSiteToSiteVPNSite(site *Site) ([]*MagicSiteToSiteVPN, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for Site Magic VPN data, site %s", site.SiteName)

	path := fmt.Sprintf(APIMagicSiteToSiteVPNPath, site.Name)

	body, err := u.GetJSON(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch site magic VPN for site %s: %w", site.SiteName, err)
	}

	var raw []*MagicSiteToSiteVPN
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse site magic VPN for site %s: %w", site.SiteName, err)
	}

	meshes := make([]*MagicSiteToSiteVPN, 0, len(raw))

	for _, m := range raw {
		if m == nil {
			continue
		}

		m.SiteName = site.SiteName
		m.SourceName = u.URL
		meshes = append(meshes, m)
	}

	return meshes, nil
}

// MagicSiteToSiteVPN represents a Site Magic mesh VPN configuration from the UniFi controller.
type MagicSiteToSiteVPN struct {
	ID           string                `fake:"{uuid}"      json:"id"`
	Name         string                `fake:"{buzzword}"  json:"name"`
	Pause        FlexBool              `json:"pause"`
	TunnelSubnet string                `fake:"{ipv4cidr}"  json:"tunnelSubnet"`
	Connections  []MagicVPNConnection  `json:"connections"`
	Devices      []MagicVPNDevice      `json:"devices"`
	Status       []MagicVPNStatusEntry `json:"status"`

	SiteName   string `json:"-"`
	SourceName string `json:"-"`
}

// MagicVPNConnection represents a connection entry within a Site Magic mesh.
type MagicVPNConnection struct {
	ConnectionID string                   `fake:"{uuid}" json:"connectionId"`
	Sites        []MagicVPNConnectionSite `json:"sites"`
}

// MagicVPNConnectionSite represents a site endpoint within a mesh connection.
type MagicVPNConnectionSite struct {
	DeviceID string  `fake:"{uuid}"              json:"deviceId"`
	Port     FlexInt `fake:"{number:1024,65535}" json:"port"`
	SiteID   string  `fake:"{uuid}"              json:"siteId"`
	TunnelIP string  `fake:"{ipv4cidr}"          json:"tunnelIp"`
	WANIP    string  `fake:"{ipv4address}"       json:"wanIp"`
}

// MagicVPNDevice represents a device participating in a Site Magic mesh.
type MagicVPNDevice struct {
	DeviceID   string               `fake:"{uuid}"     json:"deviceId"`
	DeviceName string               `fake:"{username}" json:"deviceName"`
	Sites      []MagicVPNDeviceSite `json:"sites"`
}

// MagicVPNDeviceSite represents a site on a device within a mesh.
type MagicVPNDeviceSite struct {
	SiteID   string                      `fake:"{uuid}"     json:"siteId"`
	SiteName string                      `fake:"{buzzword}" json:"siteName"`
	Networks []MagicVPNDeviceSiteNetwork `json:"networks"`
}

// MagicVPNDeviceSiteNetwork represents a network exported through a mesh VPN site.
type MagicVPNDeviceSiteNetwork struct {
	NetworkID   string `fake:"{uuid}"        json:"networkId"`
	NetworkName string `fake:"{buzzword}"    json:"networkName"`
	RouterIP    string `fake:"{ipv4address}" json:"routerIp"`
	Subnet      string `fake:"{ipv4cidr}"    json:"subnet"`
}

// MagicVPNStatusEntry represents the status of a site within a mesh VPN.
type MagicVPNStatusEntry struct {
	SiteID      string                     `fake:"{uuid}"      json:"siteId"`
	Connections []MagicVPNStatusConnection `json:"connections"`
	Errors      []string                   `json:"errors"`
	Warnings    []string                   `json:"warnings"`
}

// MagicVPNStatusConnection represents the status of a single connection within a mesh VPN.
type MagicVPNStatusConnection struct {
	ConnectionID    string   `fake:"{uuid}"     json:"connectionId"`
	Connected       FlexBool `json:"connected"`
	AssociationTime FlexInt  `fake:"{unixtime}" json:"associationTime"`
	Errors          []string `json:"errors"`
}
