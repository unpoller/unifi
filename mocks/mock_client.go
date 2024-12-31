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
