package unifi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetActiveDHCPLeases returns active DHCP leases for the given sites.
func (u *Unifi) GetActiveDHCPLeases(sites []*Site) ([]*DHCPLease, error) {
	leases := make([]*DHCPLease, 0)

	for _, site := range sites {
		var response struct {
			Data []json.RawMessage `json:"data"`
		}

		leasePath := fmt.Sprintf(APIActiveDHCPLeasesPath, site.Name)
		if err := u.GetData(leasePath, &response); err != nil {
			return nil, fmt.Errorf("failed to fetch DHCP leases for site %s: %w", site.SiteName, err)
		}

		for _, data := range response.Data {
			lease, err := u.parseDHCPLease(data, site)
			if err != nil {
				return leases, fmt.Errorf("failed to parse DHCP lease: %w", err)
			}

			leases = append(leases, lease)
		}
	}

	return leases, nil
}

// GetActiveDHCPLeasesWithAssociations returns active DHCP leases enriched with client and device associations.
// This method fetches leases, clients, and devices, then matches them by MAC address and IP address.
func (u *Unifi) GetActiveDHCPLeasesWithAssociations(sites []*Site) ([]*DHCPLease, error) {
	leases, err := u.GetActiveDHCPLeases(sites)
	if err != nil {
		return nil, err
	}

	// Fetch clients and devices for association
	clients, err := u.GetClients(sites)
	if err != nil {
		// Log error but continue without client associations
		u.ErrorLog("failed to fetch clients for DHCP lease association: %v", err)
	}

	devices, err := u.GetDevices(sites)
	if err != nil {
		// Log error but continue without device associations
		u.ErrorLog("failed to fetch devices for DHCP lease association: %v", err)
	}

	// Create lookup maps for efficient matching
	clientByMAC := make(map[string]*Client)
	clientByIP := make(map[string]*Client)

	for _, client := range clients {
		if client.Mac != "" {
			clientByMAC[normalizeMAC(client.Mac)] = client
		}

		if client.IP != "" {
			clientByIP[client.IP] = client
		}
	}

	deviceByMAC := make(map[string]interface{})

	// Add all device types to the map
	for _, device := range devices.UAPs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.USWs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.USGs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.UDMs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.UXGs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.PDUs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.UBBs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	for _, device := range devices.UCIs {
		if device.Mac != "" {
			deviceByMAC[normalizeMAC(device.Mac)] = device
		}
	}

	// Enrich leases with associations
	for _, lease := range leases {
		// Match by MAC address (most reliable)
		if lease.Mac != "" {
			normalizedMAC := normalizeMAC(lease.Mac)

			if client, found := clientByMAC[normalizedMAC]; found {
				lease.AssociatedClient = client
			}

			if device, found := deviceByMAC[normalizedMAC]; found {
				lease.AssociatedDevice = device
			}
		}

		// Also try matching by IP address as fallback
		if lease.AssociatedClient == nil && lease.IP != "" {
			if client, found := clientByIP[lease.IP]; found {
				lease.AssociatedClient = client
			}
		}
	}

	return leases, nil
}

// parseDHCPLease parses the raw JSON from the UniFi Controller into a DHCP lease structure.
func (u *Unifi) parseDHCPLease(data json.RawMessage, site *Site) (*DHCPLease, error) {
	lease := new(DHCPLease)

	if err := u.unmarshalDevice(site.SiteName, data, lease); err != nil {
		return nil, err
	}

	lease.SiteName = site.SiteName
	lease.SourceName = u.URL

	// Normalize MAC address for consistent matching
	if lease.Mac != "" {
		lease.Mac = normalizeMAC(lease.Mac)
	}

	return lease, nil
}

// normalizeMAC normalizes a MAC address to lowercase with colons for consistent matching.
// Handles MAC addresses in various formats: "aa:bb:cc:dd:ee:ff", "aa-bb-cc-dd-ee-ff", etc.
func normalizeMAC(mac string) string {
	// Convert to lowercase
	mac = strings.ToLower(mac)

	// Replace dashes with colons
	mac = strings.ReplaceAll(mac, "-", ":")

	// Remove any spaces
	mac = strings.ReplaceAll(mac, " ", "")

	return mac
}

// DHCPLease represents an active DHCP lease from the UniFi controller.
type DHCPLease struct {
	// Basic lease information
	IP        string `json:"ip"`         // Assigned IP address
	Mac       string `json:"mac"`        // MAC address of the client
	Hostname  string `json:"hostname"`   // Hostname of the client
	Network   string `json:"network"`    // Network name this lease belongs to
	NetworkID string `json:"network_id"` // Network ID this lease belongs to

	// Lease timing information
	LeaseStart FlexInt `json:"lease_start"` // Lease start timestamp
	LeaseEnd   FlexInt `json:"lease_end"`   // Lease expiration timestamp
	LeaseTime  FlexInt `json:"lease_time"`  // Lease duration in seconds

	// Additional metadata
	ClientName string   `json:"client_name"` // Client name if available
	IsStatic   FlexBool `json:"is_static"`   // Whether this is a static/reserved lease

	// Site and source tracking
	SiteName   string `json:"-"`
	SourceName string `json:"-"`

	// Associations (populated by GetActiveDHCPLeasesWithAssociations)
	AssociatedClient *Client     `json:"-"` // Associated client if found
	AssociatedDevice interface{} `json:"-"` // Associated device if found (UAP, USW, USG, UDM, UXG, PDU, UBB, or UCI)
}
