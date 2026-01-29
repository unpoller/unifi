package unifi

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// GetActiveDHCPLeases returns active DHCP leases for the given sites.
func (u *Unifi) GetActiveDHCPLeases(sites []*Site) ([]*DHCPLease, error) {
	leases := make([]*DHCPLease, 0)

	for _, site := range sites {
		var response struct {
			DHCPLeaseInfo []json.RawMessage `json:"dhcp_lease_info"`
		}

		leasePath := fmt.Sprintf(APIActiveDHCPLeasesPath, site.Name)
		if err := u.GetData(leasePath, &response); err != nil {
			return nil, fmt.Errorf("failed to fetch DHCP leases for site %s: %w", site.SiteName, err)
		}

		for _, data := range response.DHCPLeaseInfo {
			lease, err := u.parseDHCPLease(data, site)
			if err != nil {
				return leases, fmt.Errorf("failed to parse DHCP lease: %w", err)
			}

			leases = append(leases, lease)
		}
	}

	return leases, nil
}

// GetActiveDHCPLeasesWithAssociations returns active DHCP leases enriched with client, device, and network associations.
// This method fetches leases, clients, devices, and networks, then matches them appropriately.
func (u *Unifi) GetActiveDHCPLeasesWithAssociations(sites []*Site) ([]*DHCPLease, error) {
	leases, err := u.GetActiveDHCPLeases(sites)
	if err != nil {
		return nil, err
	}

	// Fetch clients, devices, and networks for association
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

	networks, err := u.GetNetworks(sites)
	if err != nil {
		// Log error but continue without network associations
		u.ErrorLog("failed to fetch networks for DHCP lease association: %v", err)
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

	// Create network lookup map by ID and name
	networkByID := make(map[string]*Network)
	networkByName := make(map[string]*Network)

	for i := range networks {
		network := &networks[i]
		if network.ID != "" {
			networkByID[network.ID] = network
		}

		if network.Name != "" {
			networkByName[network.Name] = network
		}
	}

	// Also get network info from device NetworkTable (has DhcpdStart/DhcpdStop)
	networkTableByID := make(map[string]*NetworkTableEntry)

	for _, device := range devices.UDMs {
		for i := range device.NetworkTable {
			nt := device.NetworkTable[i]
			if nt.ID != "" {
				networkTableByID[nt.ID] = &NetworkTableEntry{
					ID:                   nt.ID,
					Name:                 nt.Name,
					DhcpdEnabled:         nt.DhcpdEnabled,
					DhcpdStart:           nt.DhcpdStart,
					DhcpdStop:            nt.DhcpdStop,
					ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
					DhcpdLeasetime:       nt.DhcpdLeasetime,
					IPSubnet:             nt.IPSubnet,
				}
			}
		}
	}

	for _, device := range devices.UXGs {
		for i := range device.NetworkTable {
			nt := device.NetworkTable[i]
			if nt.ID != "" {
				networkTableByID[nt.ID] = &NetworkTableEntry{
					ID:                   nt.ID,
					Name:                 nt.Name,
					DhcpdEnabled:         nt.DhcpdEnabled,
					DhcpdStart:           nt.DhcpdStart,
					DhcpdStop:            nt.DhcpdStop,
					ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
					DhcpdLeasetime:       nt.DhcpdLeasetime,
					IPSubnet:             nt.IPSubnet,
				}
			}
		}
	}

	for _, device := range devices.USGs {
		for i := range device.NetworkTable {
			nt := device.NetworkTable[i]
			if nt.ID != "" {
				networkTableByID[nt.ID] = &NetworkTableEntry{
					ID:                   nt.ID,
					Name:                 nt.Name,
					DhcpdEnabled:         nt.DhcpdEnabled,
					DhcpdStart:           nt.DhcpdStart,
					DhcpdStop:            nt.DhcpdStop,
					ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
					DhcpdLeasetime:       nt.DhcpdLeasetime,
					IPSubnet:             nt.IPSubnet,
				}
			}
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

		// Match network by NetworkID first, then by Network name
		// The API response may not include the readable network name, so we populate it
		// from the associated network data for better metric labeling.
		if lease.NetworkID != "" {
			if network, found := networkByID[lease.NetworkID]; found {
				lease.AssociatedNetwork = network
				// Populate Network name from association (e.g., "LAN", "Guest") for readable metrics
				// This ensures the network label in Prometheus metrics shows a human-readable name
				// instead of an empty string, even when the API doesn't provide it directly.
				if lease.Network == "" && network.Name != "" {
					lease.Network = network.Name
				}
			}

			// Also check NetworkTable for DHCP pool info (DhcpdStart/DhcpdStop range)
			if ntEntry, found := networkTableByID[lease.NetworkID]; found {
				lease.NetworkTableEntry = ntEntry
				// Populate Network name from NetworkTable if not already set from AssociatedNetwork
				// NetworkTable entries also contain the network name and are available from device data
				if lease.Network == "" && ntEntry.Name != "" {
					lease.Network = ntEntry.Name
				}
			}
		} else if lease.Network != "" {
			if network, found := networkByName[lease.Network]; found {
				lease.AssociatedNetwork = network
			}
		}
	}

	return leases, nil
}

// AssociateDHCPLeases associates DHCP leases with clients, devices, and networks using pre-fetched data.
// This method is used to avoid redundant API calls when the data is already available.
func (u *Unifi) AssociateDHCPLeases(leases []*DHCPLease, clients []*Client, devices *Devices, networks []Network) error {
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
	if devices != nil {
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
	}

	// Create network lookup map by ID and name
	networkByID := make(map[string]*Network)
	networkByName := make(map[string]*Network)

	for i := range networks {
		network := &networks[i]
		if network.ID != "" {
			networkByID[network.ID] = network
		}

		if network.Name != "" {
			networkByName[network.Name] = network
		}
	}

	// Also get network info from device NetworkTable (has DhcpdStart/DhcpdStop)
	networkTableByID := make(map[string]*NetworkTableEntry)

	if devices != nil {
		for _, device := range devices.UDMs {
			for i := range device.NetworkTable {
				nt := device.NetworkTable[i]
				if nt.ID != "" {
					networkTableByID[nt.ID] = &NetworkTableEntry{
						ID:                   nt.ID,
						Name:                 nt.Name,
						DhcpdEnabled:         nt.DhcpdEnabled,
						DhcpdStart:           nt.DhcpdStart,
						DhcpdStop:            nt.DhcpdStop,
						ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
						DhcpdLeasetime:       nt.DhcpdLeasetime,
						IPSubnet:             nt.IPSubnet,
					}
				}
			}
		}

		for _, device := range devices.USGs {
			for i := range device.NetworkTable {
				nt := device.NetworkTable[i]
				if nt.ID != "" {
					networkTableByID[nt.ID] = &NetworkTableEntry{
						ID:                   nt.ID,
						Name:                 nt.Name,
						DhcpdEnabled:         nt.DhcpdEnabled,
						DhcpdStart:           nt.DhcpdStart,
						DhcpdStop:            nt.DhcpdStop,
						ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
						DhcpdLeasetime:       nt.DhcpdLeasetime,
						IPSubnet:             nt.IPSubnet,
					}
				}
			}
		}

		for _, device := range devices.UXGs {
			for i := range device.NetworkTable {
				nt := device.NetworkTable[i]
				if nt.ID != "" {
					networkTableByID[nt.ID] = &NetworkTableEntry{
						ID:                   nt.ID,
						Name:                 nt.Name,
						DhcpdEnabled:         nt.DhcpdEnabled,
						DhcpdStart:           nt.DhcpdStart,
						DhcpdStop:            nt.DhcpdStop,
						ActiveDhcpLeaseCount: nt.ActiveDhcpLeaseCount,
						DhcpdLeasetime:       nt.DhcpdLeasetime,
						IPSubnet:             nt.IPSubnet,
					}
				}
			}
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

		// Match network by NetworkID first, then by Network name
		if lease.NetworkID != "" {
			if network, found := networkByID[lease.NetworkID]; found {
				lease.AssociatedNetwork = network
				if lease.Network == "" && network.Name != "" {
					lease.Network = network.Name
				}
			}

			// Also check NetworkTable for DHCP pool info
			if ntEntry, found := networkTableByID[lease.NetworkID]; found {
				lease.NetworkTableEntry = ntEntry
				if lease.Network == "" && ntEntry.Name != "" {
					lease.Network = ntEntry.Name
				}
			}
		} else if lease.Network != "" {
			if network, found := networkByName[lease.Network]; found {
				lease.AssociatedNetwork = network
			}
		}
	}

	return nil
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

	// Set ClientName from display_name if hostname is empty
	if lease.ClientName != "" && lease.Hostname == "" {
		lease.Hostname = lease.ClientName
	}

	// Calculate lease_start if we have lease_end (expiration_time) and can estimate duration
	// Note: We don't have lease duration from the API, so lease_start will remain 0
	// The lease_time will also remain 0 unless we can get it from network config

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
	Network   string `json:"network"`    // Network name this lease belongs to (populated from association)
	NetworkID string `json:"network_id"` // Network ID this lease belongs to

	// Lease timing information
	LeaseStart FlexInt `json:"lease_start"`           // Lease start timestamp (calculated from expiration)
	LeaseEnd   FlexInt `json:"lease_expiration_time"` // Lease expiration timestamp
	LeaseTime  FlexInt `json:"lease_time"`            // Lease duration in seconds (calculated)

	// Additional metadata
	ClientName string   `json:"display_name"` // Display name from API
	IsStatic   FlexBool `json:"use_fixedip"`  // Whether this is a static/reserved lease
	Status     string   `json:"status"`       // online/offline status
	ClientType string   `json:"client_type"`  // WIRELESS, etc.
	DeviceType string   `json:"device_type"`  // usw, etc. (for UniFi devices)

	// Site and source tracking
	SiteName   string `json:"-"`
	SourceName string `json:"-"`

	// Associations (populated by GetActiveDHCPLeasesWithAssociations)
	AssociatedClient  *Client            `json:"-"` // Associated client if found
	AssociatedDevice  interface{}        `json:"-"` // Associated device if found (UAP, USW, USG, UDM, UXG, PDU, UBB, or UCI)
	AssociatedNetwork *Network           `json:"-"` // Associated network if found
	NetworkTableEntry *NetworkTableEntry `json:"-"` // Network table entry from device (contains DHCP pool range)
}

// NetworkTableEntry represents a network entry from a device's NetworkTable.
// This contains DHCP pool information (DhcpdStart/DhcpdStop) not available in the Network struct.
type NetworkTableEntry struct {
	ID                   string   `json:"_id"`
	Name                 string   `json:"name"`
	DhcpdEnabled         FlexBool `json:"dhcpd_enabled"`
	DhcpdStart           string   `json:"dhcpd_start"`
	DhcpdStop            string   `json:"dhcpd_stop"`
	ActiveDhcpLeaseCount FlexInt  `json:"active_dhcp_lease_count"`
	DhcpdLeasetime       FlexInt  `json:"dhcpd_leasetime"`
	IPSubnet             string   `json:"ip_subnet"`
}

// GetPoolSize calculates the DHCP pool size (number of available IPs) from DhcpdStart and DhcpdStop.
// Returns 0 if the range cannot be calculated or if DHCP is not enabled.
func (l *DHCPLease) GetPoolSize() int {
	if l.NetworkTableEntry == nil {
		return 0
	}

	if !l.NetworkTableEntry.DhcpdEnabled.Val {
		return 0
	}

	if l.NetworkTableEntry.DhcpdStart == "" || l.NetworkTableEntry.DhcpdStop == "" {
		return 0
	}

	return calculateIPRangeSize(l.NetworkTableEntry.DhcpdStart, l.NetworkTableEntry.DhcpdStop)
}

// GetActiveLeaseCount returns the number of active DHCP leases for this network.
// Returns 0 if network table entry is not available.
func (l *DHCPLease) GetActiveLeaseCount() int {
	if l.NetworkTableEntry == nil {
		return 0
	}

	return int(l.NetworkTableEntry.ActiveDhcpLeaseCount.Val)
}

// GetUtilizationPercentage calculates the DHCP pool utilization percentage.
// Returns 0 if pool size cannot be determined.
func (l *DHCPLease) GetUtilizationPercentage() float64 {
	poolSize := l.GetPoolSize()
	if poolSize == 0 {
		return 0
	}

	activeCount := l.GetActiveLeaseCount()
	if activeCount == 0 {
		return 0
	}

	return (float64(activeCount) / float64(poolSize)) * 100.0
}

// GetAvailableIPs returns the number of available IPs in the DHCP pool.
func (l *DHCPLease) GetAvailableIPs() int {
	poolSize := l.GetPoolSize()
	if poolSize == 0 {
		return 0
	}

	activeCount := l.GetActiveLeaseCount()
	available := poolSize - activeCount

	if available < 0 {
		return 0
	}

	return available
}

// calculateIPRangeSize calculates the number of IPs in a range from start to stop (inclusive).
func calculateIPRangeSize(startIP, stopIP string) int {
	start := net.ParseIP(startIP)
	stop := net.ParseIP(stopIP)

	if start == nil || stop == nil {
		return 0
	}

	// Convert to IPv4 addresses
	start4 := start.To4()
	stop4 := stop.To4()

	if start4 == nil || stop4 == nil {
		return 0
	}

	// Convert to uint32 for calculation
	startInt := uint32(start4[0])<<24 | uint32(start4[1])<<16 | uint32(start4[2])<<8 | uint32(start4[3])
	stopInt := uint32(stop4[0])<<24 | uint32(stop4[1])<<16 | uint32(stop4[2])<<8 | uint32(stop4[3])

	if stopInt < startInt {
		return 0
	}

	// Add 1 because the range is inclusive
	return int(stopInt-startInt) + 1
}
