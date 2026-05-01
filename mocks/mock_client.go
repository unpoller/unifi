package mocks

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/unpoller/unifi/v5"
)

type MockUnifi struct {
	*unifi.Config
}

// ensure MockUnifi implements the interface fully, this will fail to compile otherwise
var _ unifi.UnifiClient = &MockUnifi{}

func NewMockUnifi() *MockUnifi {
	return &MockUnifi{}
}

const numItemsMocked = 5

// GetAlarms returns Alarms for a list of Sites.
func (m *MockUnifi) GetAlarms(_ []*unifi.Site) ([]*unifi.Alarm, error) {
	alarms := make([]*unifi.Alarm, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Alarm

		err := gofakeit.Struct(&a)
		if err != nil {
			return alarms, err
		}

		alarms[i] = &a
	}

	return alarms, nil
}

// GetAlarmsSite retreives the Alarms for a single Site.
func (m *MockUnifi) GetAlarmsSite(_ *unifi.Site) ([]*unifi.Alarm, error) {
	alarms := make([]*unifi.Alarm, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Alarm

		err := gofakeit.Struct(&a)
		if err != nil {
			return alarms, err
		}

		alarms[i] = &a
	}

	return alarms, nil
}

// GetAnomalies returns Anomalies for a list of Sites.
func (m *MockUnifi) GetAnomalies(_ []*unifi.Site, _ ...time.Time) ([]*unifi.Anomaly, error) {
	results := make([]*unifi.Anomaly, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Anomaly

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetAnomaliesSite retreives the Anomalies for a single Site.
func (m *MockUnifi) GetAnomaliesSite(_ *unifi.Site, _ ...time.Time) ([]*unifi.Anomaly, error) {
	results := make([]*unifi.Anomaly, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Anomaly

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetClients returns a response full of clients' data from the UniFi Controller.
func (m *MockUnifi) GetClients(_ []*unifi.Site) ([]*unifi.Client, error) {
	results := make([]*unifi.Client, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Client

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetClientsDPI garners dpi data for clients.
func (m *MockUnifi) GetClientsDPI(_ []*unifi.Site) ([]*unifi.DPITable, error) {
	results := make([]*unifi.DPITable, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.DPITable

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetClientHistory returns a response full of client history from the UniFi Controller
func (m *MockUnifi) GetClientHistory(_ []*unifi.Site, _ *unifi.ClientHistoryOpts) ([]*unifi.ClientHistory, error) {
	// TODO : add logic to generate data based on ClientHistoryOpts
	results := make([]*unifi.ClientHistory, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.ClientHistory

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetDevices returns a response full of devices' data from the UniFi Controller.
func (m *MockUnifi) GetDevices(_ []*unifi.Site) (*unifi.Devices, error) {
	var d unifi.Devices

	err := gofakeit.Struct(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// GetUSWs returns all switches, an error, or nil if there are no switches.
func (m *MockUnifi) GetUSWs(_ *unifi.Site) ([]*unifi.USW, error) {
	results := make([]*unifi.USW, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.USW

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUAPs returns all access points, an error, or nil if there are no APs.
func (m *MockUnifi) GetUAPs(_ *unifi.Site) ([]*unifi.UAP, error) {
	results := make([]*unifi.UAP, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UAP

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUDMs returns all dream machines, an error, or nil if there are no UDMs.
func (m *MockUnifi) GetUDMs(_ *unifi.Site) ([]*unifi.UDM, error) {
	results := make([]*unifi.UDM, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UDM

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUXGs returns all 10Gb gateways, an error, or nil if there are no UXGs.
func (m *MockUnifi) GetUXGs(_ *unifi.Site) ([]*unifi.UXG, error) {
	results := make([]*unifi.UXG, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UXG

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUBBs returns all UBB devices, an error, or nil if there are no UBBs.
func (m *MockUnifi) GetUBBs(_ *unifi.Site) ([]*unifi.UBB, error) {
	results := make([]*unifi.UBB, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UBB

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUCIs returns all UCI devices, an error, or nil if there are no UCIs.
func (m *MockUnifi) GetUCIs(_ *unifi.Site) ([]*unifi.UCI, error) {
	results := make([]*unifi.UCI, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UCI

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetPDUs returns all PDU devices, an error, or nil if there are no PDUs.
func (m *MockUnifi) GetPDUs(_ *unifi.Site) ([]*unifi.PDU, error) {
	results := make([]*unifi.PDU, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.PDU

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUDBs returns all UDB devices, an error, or nil if there are no UDBs.
func (m *MockUnifi) GetUDBs(_ *unifi.Site) ([]*unifi.UDB, error) {
	results := make([]*unifi.UDB, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UDB

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetUSGs returns all 1Gb gateways, an error, or nil if there are no USGs.
func (m *MockUnifi) GetUSGs(_ *unifi.Site) ([]*unifi.USG, error) {
	results := make([]*unifi.USG, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.USG

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetEvents returns a response full of UniFi Events for the last 1 hour from multiple sites.
func (m *MockUnifi) GetEvents(_ []*unifi.Site, _ time.Duration) ([]*unifi.Event, error) {
	results := make([]*unifi.Event, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Event

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSiteEvents retrieves the last 1 hour's worth of events from a single site.
func (m *MockUnifi) GetSiteEvents(_ *unifi.Site, _ time.Duration) ([]*unifi.Event, error) {
	results := make([]*unifi.Event, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Event

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetIDS returns Intrusion Detection Systems events for a list of Sites.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (m *MockUnifi) GetIDS(_ []*unifi.Site, _ ...time.Time) ([]*unifi.IDS, error) { //nolint:revive
	results := make([]*unifi.IDS, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IDS

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetIDSSite retrieves the Intrusion Detection System Data for a single Site.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (m *MockUnifi) GetIDSSite(_ *unifi.Site, _ ...time.Time) ([]*unifi.IDS, error) {
	results := make([]*unifi.IDS, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IDS

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetNetworks returns a response full of network data from the UniFi Controller.
func (m *MockUnifi) GetNetworks(_ []*unifi.Site) ([]unifi.Network, error) {
	results := make([]unifi.Network, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Network

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = a
	}

	return results, nil
}

// GetActiveDHCPLeases returns active DHCP leases for the given sites.
func (m *MockUnifi) GetActiveDHCPLeases(_ []*unifi.Site) ([]*unifi.DHCPLease, error) {
	results := make([]*unifi.DHCPLease, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var l unifi.DHCPLease

		err := gofakeit.Struct(&l)
		if err != nil {
			return results, err
		}

		results[i] = &l
	}

	return results, nil
}

// GetActiveDHCPLeasesWithAssociations returns active DHCP leases enriched with client and device associations.
func (m *MockUnifi) GetActiveDHCPLeasesWithAssociations(_ []*unifi.Site) ([]*unifi.DHCPLease, error) {
	results := make([]*unifi.DHCPLease, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var l unifi.DHCPLease

		err := gofakeit.Struct(&l)
		if err != nil {
			return results, err
		}

		results[i] = &l
	}

	return results, nil
}

// AssociateDHCPLeases associates DHCP leases with clients, devices, and networks using pre-fetched data.
func (m *MockUnifi) AssociateDHCPLeases(_ []*unifi.DHCPLease, _ []*unifi.Client, _ *unifi.Devices, _ []unifi.Network) error {
	return nil
}

// GetWANEnrichedConfiguration returns enriched WAN configuration for all WAN interfaces.
func (m *MockUnifi) GetWANEnrichedConfiguration(_ []*unifi.Site) ([]*unifi.WANEnrichedConfiguration, error) {
	results := make([]*unifi.WANEnrichedConfiguration, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var w unifi.WANEnrichedConfiguration

		err := gofakeit.Struct(&w)
		if err != nil {
			return results, err
		}

		results[i] = &w
	}

	return results, nil
}

// GetWANLoadBalancingStatus returns the current load balancing status for WAN interfaces.
func (m *MockUnifi) GetWANLoadBalancingStatus(_ []*unifi.Site) (*unifi.WANLoadBalancingStatus, error) {
	var w unifi.WANLoadBalancingStatus

	err := gofakeit.Struct(&w)
	if err != nil {
		return &w, err
	}

	return &w, nil
}

// GetWANISPStatus returns the ISP status for WAN interfaces.
func (m *MockUnifi) GetWANISPStatus(_ []*unifi.Site, _ string) (*unifi.WANISPStatusDetailed, error) {
	var w unifi.WANISPStatusDetailed

	err := gofakeit.Struct(&w)
	if err != nil {
		return &w, err
	}

	return &w, nil
}

// GetWANSLAs returns WAN SLA monitoring data.
func (m *MockUnifi) GetWANSLAs(_ []*unifi.Site) ([]*unifi.WANSLA, error) {
	results := make([]*unifi.WANSLA, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var w unifi.WANSLA

		err := gofakeit.Struct(&w)
		if err != nil {
			return results, err
		}

		results[i] = &w
	}

	return results, nil
}

// GetSites returns a list of configured sites on the UniFi controller.
func (m *MockUnifi) GetSites() ([]*unifi.Site, error) {
	results := make([]*unifi.Site, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Site

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSiteDPI garners dpi data for sites.
func (m *MockUnifi) GetSiteDPI(_ []*unifi.Site) ([]*unifi.DPITable, error) {
	// we only should ever return 1 of these, regardless, unless there is an error, return 0
	var a unifi.DPITable

	err := gofakeit.Struct(&a)
	if err != nil {
		return []*unifi.DPITable{}, err
	}

	return []*unifi.DPITable{&a}, nil
}

// GetRogueAPs returns RogueAPs for a list of Sites.
// Use GetRogueAPsSite if you want more control.
func (m *MockUnifi) GetRogueAPs(_ []*unifi.Site) ([]*unifi.RogueAP, error) {
	results := make([]*unifi.RogueAP, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.RogueAP

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetRogueAPsSite returns RogueAPs for a single Site.
func (m *MockUnifi) GetRogueAPsSite(_ *unifi.Site) ([]*unifi.RogueAP, error) {
	results := make([]*unifi.RogueAP, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.RogueAP

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// Login is a helper method. It can be called to grab a new authentication cookie.
func (m *MockUnifi) Login() error {
	return nil
}

// Logout closes the current session.
func (m *MockUnifi) Logout() error {
	return nil
}

// GetServerData sets the controller's version and UUID. Only call this if you
// previously called Login and suspect the controller version has changed.
func (m *MockUnifi) GetServerData() (*unifi.ServerStatus, error) {
	var response unifi.ServerStatus

	err := gofakeit.Struct(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUsers returns a response full of clients that connected to the UDM within the provided amount of time
// using the insight historical connection data set.
func (m *MockUnifi) GetUsers(_ []*unifi.Site, _ int) ([]*unifi.User, error) {
	results := make([]*unifi.User, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.User

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

func (m *MockUnifi) GetClientTraffic(_ []*unifi.Site, _ *unifi.EpochMillisTimePeriod, _ bool) ([]*unifi.ClientUsageByApp, error) {
	results := make([]*unifi.ClientUsageByApp, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.ClientUsageByApp

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

func (m *MockUnifi) GetClientTrafficByMac(_ *unifi.Site, _ *unifi.EpochMillisTimePeriod, _ bool, macs ...string) ([]*unifi.ClientUsageByApp, error) {
	results := make([]*unifi.ClientUsageByApp, len(macs))

	for i := 0; i < len(macs); i++ {
		var a unifi.ClientUsageByApp

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

func (m *MockUnifi) GetCountryTraffic(_ []*unifi.Site, _ *unifi.EpochMillisTimePeriod) ([]*unifi.UsageByCountry, error) {
	results := make([]*unifi.UsageByCountry, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UsageByCountry

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetProtectLogs returns Protect system log events.
func (m *MockUnifi) GetProtectLogs(_ *unifi.ProtectLogRequest) ([]*unifi.ProtectLogEntry, error) {
	results := make([]*unifi.ProtectLogEntry, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.ProtectLogEntry

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSysinfo returns controller system info and health (UniFi OS).
func (m *MockUnifi) GetSysinfo(_ []*unifi.Site) ([]*unifi.Sysinfo, error) {
	results := make([]*unifi.Sysinfo, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var s unifi.Sysinfo

		err := gofakeit.Struct(&s)
		if err != nil {
			return results, err
		}

		results[i] = &s
	}

	return results, nil
}

// GetPortAnomalies returns port anomalies for a list of Sites.
func (m *MockUnifi) GetPortAnomalies(_ []*unifi.Site) ([]*unifi.PortAnomaly, error) {
	results := make([]*unifi.PortAnomaly, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.PortAnomaly

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetPortAnomaliesSite returns port anomalies for a single Site.
func (m *MockUnifi) GetPortAnomaliesSite(_ *unifi.Site) ([]*unifi.PortAnomaly, error) {
	results := make([]*unifi.PortAnomaly, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.PortAnomaly

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetMagicSiteToSiteVPN returns Site Magic VPN mesh configurations for a list of Sites.
func (m *MockUnifi) GetMagicSiteToSiteVPN(_ []*unifi.Site) ([]*unifi.MagicSiteToSiteVPN, error) {
	results := make([]*unifi.MagicSiteToSiteVPN, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var v unifi.MagicSiteToSiteVPN

		err := gofakeit.Struct(&v)
		if err != nil {
			return results, err
		}

		results[i] = &v
	}

	return results, nil
}

// GetMagicSiteToSiteVPNSite returns Site Magic VPN mesh configurations for a single Site.
func (m *MockUnifi) GetMagicSiteToSiteVPNSite(_ *unifi.Site) ([]*unifi.MagicSiteToSiteVPN, error) {
	results := make([]*unifi.MagicSiteToSiteVPN, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var v unifi.MagicSiteToSiteVPN

		err := gofakeit.Struct(&v)
		if err != nil {
			return results, err
		}

		results[i] = &v
	}

	return results, nil
}

// GetIntegrationSites returns all sites from the Integration/v1 API.
func (m *MockUnifi) GetIntegrationSites() ([]*unifi.IntegrationSite, error) {
	results := make([]*unifi.IntegrationSite, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IntegrationSite

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetIntegrationInfo returns application version info from the Integration/v1 API.
func (m *MockUnifi) GetIntegrationInfo() (*unifi.IntegrationInfo, error) {
	var a unifi.IntegrationInfo

	err := gofakeit.Struct(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetIntegrationDeviceStats returns statistics for a single device.
func (m *MockUnifi) GetIntegrationDeviceStats(_ *unifi.IntegrationSite, _ string) (*unifi.IntegrationDeviceStats, error) {
	var s unifi.IntegrationDeviceStats

	err := gofakeit.Struct(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetAllIntegrationDeviceStats returns statistics for all devices in a site.
func (m *MockUnifi) GetAllIntegrationDeviceStats(_ *unifi.IntegrationSite) ([]*unifi.IntegrationDeviceStats, error) {
	results := make([]*unifi.IntegrationDeviceStats, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IntegrationDeviceStats

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetWifiBroadcasts returns WiFi broadcast configurations for a site.
func (m *MockUnifi) GetWifiBroadcasts(_ *unifi.IntegrationSite) ([]*unifi.WifiBroadcast, error) {
	results := make([]*unifi.WifiBroadcast, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.WifiBroadcast

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetFirewallZones returns firewall zones for a site.
func (m *MockUnifi) GetFirewallZones(_ *unifi.IntegrationSite) ([]*unifi.FirewallZone, error) {
	results := make([]*unifi.FirewallZone, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.FirewallZone

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetACLRules returns access control rules for a site.
func (m *MockUnifi) GetACLRules(_ *unifi.IntegrationSite) ([]*unifi.ACLRule, error) {
	results := make([]*unifi.ACLRule, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.ACLRule

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetIntegrationNetworks returns networks for a site from the Integration/v1 API.
func (m *MockUnifi) GetIntegrationNetworks(_ *unifi.IntegrationSite) ([]*unifi.IntegrationNetwork, error) {
	results := make([]*unifi.IntegrationNetwork, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IntegrationNetwork

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetIntegrationWANs returns WAN interface identifiers for a site.
func (m *MockUnifi) GetIntegrationWANs(_ *unifi.IntegrationSite) ([]*unifi.IntegrationWAN, error) {
	results := make([]*unifi.IntegrationWAN, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.IntegrationWAN

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetVPNServers returns VPN server configurations for a site.
func (m *MockUnifi) GetVPNServers(_ *unifi.IntegrationSite) ([]*unifi.VPNServer, error) {
	results := make([]*unifi.VPNServer, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.VPNServer

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSiteToSiteTunnels returns site-to-site VPN tunnel configurations for a site.
func (m *MockUnifi) GetSiteToSiteTunnels(_ *unifi.IntegrationSite) ([]*unifi.SiteToSiteTunnel, error) {
	results := make([]*unifi.SiteToSiteTunnel, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.SiteToSiteTunnel

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetLAGs returns link aggregation group configurations for a site.
func (m *MockUnifi) GetLAGs(_ *unifi.IntegrationSite) ([]*unifi.LAG, error) {
	results := make([]*unifi.LAG, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.LAG

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetMCLAGDomains returns multi-chassis LAG domain configurations for a site.
func (m *MockUnifi) GetMCLAGDomains(_ *unifi.IntegrationSite) ([]*unifi.MCLAGDomain, error) {
	results := make([]*unifi.MCLAGDomain, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.MCLAGDomain

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSwitchStacks returns switch stack configurations for a site.
func (m *MockUnifi) GetSwitchStacks(_ *unifi.IntegrationSite) ([]*unifi.SwitchStack, error) {
	results := make([]*unifi.SwitchStack, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.SwitchStack

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetDNSPolicies returns DNS policies for a site.
func (m *MockUnifi) GetDNSPolicies(_ *unifi.IntegrationSite) ([]*unifi.DNSPolicy, error) {
	results := make([]*unifi.DNSPolicy, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.DNSPolicy

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetRADIUSProfiles returns RADIUS authentication profiles for a site.
func (m *MockUnifi) GetRADIUSProfiles(_ *unifi.IntegrationSite) ([]*unifi.RADIUSProfile, error) {
	results := make([]*unifi.RADIUSProfile, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.RADIUSProfile

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetTrafficMatchingLists returns traffic matching lists for a site.
func (m *MockUnifi) GetTrafficMatchingLists(_ *unifi.IntegrationSite) ([]*unifi.TrafficMatchingList, error) {
	results := make([]*unifi.TrafficMatchingList, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.TrafficMatchingList

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetHotspotVouchers returns guest portal vouchers for a site.
func (m *MockUnifi) GetHotspotVouchers(_ *unifi.IntegrationSite) ([]*unifi.HotspotVoucher, error) {
	results := make([]*unifi.HotspotVoucher, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.HotspotVoucher

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetDPIApplications returns the DPI application reference catalogue (global, no site).
func (m *MockUnifi) GetDPIApplications() ([]*unifi.DPIApplication, error) {
	results := make([]*unifi.DPIApplication, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.DPIApplication

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetDPICategories returns the DPI category reference catalogue (global, no site).
func (m *MockUnifi) GetDPICategories() ([]*unifi.DPICategory, error) {
	results := make([]*unifi.DPICategory, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.DPICategory

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetPendingDevices returns devices waiting to be adopted (global, no site).
func (m *MockUnifi) GetPendingDevices() ([]*unifi.PendingDevice, error) {
	results := make([]*unifi.PendingDevice, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.PendingDevice

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetCountries returns the list of countries for geo-based firewall filters (global, no site).
func (m *MockUnifi) GetCountries() ([]*unifi.Country, error) {
	results := make([]*unifi.Country, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.Country

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetWANStatus returns WAN interface status for a site.
func (m *MockUnifi) GetWANStatus(_ *unifi.Site) (*unifi.WANStatus, error) {
	var w unifi.WANStatus

	err := gofakeit.Struct(&w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

// GetUPSDeviceList returns UPS/PDU device selectors for a site.
func (m *MockUnifi) GetUPSDeviceList(_ *unifi.Site) ([]*unifi.UPSDeviceSelector, error) {
	results := make([]*unifi.UPSDeviceSelector, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.UPSDeviceSelector

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetPortForwards returns port forwarding rules for a site.
func (m *MockUnifi) GetPortForwards(_ *unifi.Site) ([]*unifi.PortForward, error) {
	results := make([]*unifi.PortForward, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var a unifi.PortForward

		err := gofakeit.Struct(&a)
		if err != nil {
			return results, err
		}

		results[i] = &a
	}

	return results, nil
}

// GetSSLCertificate returns SSL certificate information for a site.
func (m *MockUnifi) GetSSLCertificate(_ *unifi.Site) (*unifi.SSLCertificate, error) {
	var a unifi.SSLCertificate

	err := gofakeit.Struct(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
