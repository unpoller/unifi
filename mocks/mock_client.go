package mocks

import (
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/unpoller/unifi"
)

type MockUnifi struct {
	faker *gofakeit.Faker
}

func fakeSeedValue() int64 {
	seedVal := os.Getenv("UNPOLLER_FAKE_GEN_SEED")
	if seedVal != "" {
		if seed, err := strconv.ParseInt(seedVal, 10, 64); err != nil {
			return seed
		}
	}
	return 0
}

func NewMockUnifi() *MockUnifi {
	faker := gofakeit.New(fakeSeedValue())
	return &MockUnifi{
		faker: faker,
	}
}

// GetAlarms returns Alarms for a list of Sites.
func (m *MockUnifi) GetAlarms(sites []*unifi.Site) ([]*unifi.Alarm, error) {
	qty := m.faker.Rand.Intn(5)
	alarms := make([]*unifi.Alarm, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Alarm
		err := m.faker.Struct(&a)
		if err != nil {
			return alarms, err
		}
		alarms[i] = &a
	}
	return alarms, nil
}

// GetAlarmsSite retreives the Alarms for a single Site.
func (m *MockUnifi) GetAlarmsSite(site *unifi.Site) ([]*unifi.Alarm, error) {
	qty := m.faker.Rand.Intn(5)
	alarms := make([]*unifi.Alarm, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Alarm
		err := m.faker.Struct(&a)
		if err != nil {
			return alarms, err
		}
		alarms[i] = &a
	}
	return alarms, nil
}

// GetAnomalies returns Anomalies for a list of Sites.
func (m *MockUnifi) GetAnomalies(sites []*unifi.Site, timeRange ...time.Time) ([]*unifi.Anomaly, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Anomaly, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Anomaly
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetAnomaliesSite retreives the Anomalies for a single Site.
func (m *MockUnifi) GetAnomaliesSite(site *unifi.Site, timeRange ...time.Time) ([]*unifi.Anomaly, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Anomaly, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Anomaly
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetClients returns a response full of clients' data from the UniFi Controller.
func (m *MockUnifi) GetClients(sites []*unifi.Site) ([]*unifi.Client, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Client, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Client
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetClientsDPI garners dpi data for clients.
func (m *MockUnifi) GetClientsDPI(sites []*unifi.Site) ([]*unifi.DPITable, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.DPITable, qty)
	for i := 0; i < qty; i++ {
		var a unifi.DPITable
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetDevices returns a response full of devices' data from the UniFi Controller.
func (m *MockUnifi) GetDevices(sites []*unifi.Site) (*unifi.Devices, error) {
	var d unifi.Devices
	err := m.faker.Struct(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetUSWs returns all switches, an error, or nil if there are no switches.
func (m *MockUnifi) GetUSWs(site *unifi.Site) ([]*unifi.USW, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.USW, qty)
	for i := 0; i < qty; i++ {
		var a unifi.USW
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetUAPs returns all access points, an error, or nil if there are no APs.
func (m *MockUnifi) GetUAPs(site *unifi.Site) ([]*unifi.UAP, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.UAP, qty)
	for i := 0; i < qty; i++ {
		var a unifi.UAP
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetUDMs returns all dream machines, an error, or nil if there are no UDMs.
func (m *MockUnifi) GetUDMs(site *unifi.Site) ([]*unifi.UDM, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.UDM, qty)
	for i := 0; i < qty; i++ {
		var a unifi.UDM
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetUXGs returns all 10Gb gateways, an error, or nil if there are no UXGs.
func (m *MockUnifi) GetUXGs(site *unifi.Site) ([]*unifi.UXG, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.UXG, qty)
	for i := 0; i < qty; i++ {
		var a unifi.UXG
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetUSGs returns all 1Gb gateways, an error, or nil if there are no USGs.
func (m *MockUnifi) GetUSGs(site *unifi.Site) ([]*unifi.USG, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.USG, qty)
	for i := 0; i < qty; i++ {
		var a unifi.USG
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetEvents returns a response full of UniFi Events for the last 1 hour from multiple sites.
func (m *MockUnifi) GetEvents(sites []*unifi.Site, hours time.Duration) ([]*unifi.Event, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Event, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Event
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetSiteEvents retrieves the last 1 hour's worth of events from a single site.
func (m *MockUnifi) GetSiteEvents(site *unifi.Site, hours time.Duration) ([]*unifi.Event, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Event, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Event
		err := m.faker.Struct(&a)
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
func (m *MockUnifi) GetIDS(sites []*unifi.Site, timeRange ...time.Time) ([]*unifi.IDS, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.IDS, qty)
	for i := 0; i < qty; i++ {
		var a unifi.IDS
		err := m.faker.Struct(&a)
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
func (m *MockUnifi) GetIDSSite(site *unifi.Site, timeRange ...time.Time) ([]*unifi.IDS, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.IDS, qty)
	for i := 0; i < qty; i++ {
		var a unifi.IDS
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetNetworks returns a response full of network data from the UniFi Controller.
func (m *MockUnifi) GetNetworks(sites []*unifi.Site) ([]unifi.Network, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]unifi.Network, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Network
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = a
	}
	return results, nil
}

// GetSites returns a list of configured sites on the UniFi controller.
func (m *MockUnifi) GetSites() ([]*unifi.Site, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.Site, qty)
	for i := 0; i < qty; i++ {
		var a unifi.Site
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetSiteDPI garners dpi data for sites.
func (m *MockUnifi) GetSiteDPI(sites []*unifi.Site) ([]*unifi.DPITable, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.DPITable, qty)
	for i := 0; i < qty; i++ {
		var a unifi.DPITable
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetRogueAPs returns RogueAPs for a list of Sites.
// Use GetRogueAPsSite if you want more control.
func (m *MockUnifi) GetRogueAPs(sites []*unifi.Site) ([]*unifi.RogueAP, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.RogueAP, qty)
	for i := 0; i < qty; i++ {
		var a unifi.RogueAP
		err := m.faker.Struct(&a)
		if err != nil {
			return results, err
		}
		results[i] = &a
	}
	return results, nil
}

// GetRogueAPsSite returns RogueAPs for a single Site.
func (m *MockUnifi) GetRogueAPsSite(site *unifi.Site) ([]*unifi.RogueAP, error) {
	qty := m.faker.Rand.Intn(5)
	results := make([]*unifi.RogueAP, qty)
	for i := 0; i < qty; i++ {
		var a unifi.RogueAP
		err := m.faker.Struct(&a)
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
func (m *MockUnifi) GetServerData() error {
	return nil
}
