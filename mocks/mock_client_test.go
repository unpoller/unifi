package mocks_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5/mocks"
)

func TestMockUnifiClient(t *testing.T) {
	m := mocks.NewMockUnifi()
	devices, err := m.GetDevices(nil)
	require.NoError(t, err)
	require.NotEmpty(t, devices.UAPs)
	require.NotEmpty(t, devices.USGs)
	require.NotEmpty(t, devices.UDMs)
	require.NotEmpty(t, devices.USWs)
	require.NotEmpty(t, devices.PDUs)
	require.NotEmpty(t, devices.UCIs)
	require.NotEmpty(t, devices.UXGs)
	require.NotEmpty(t, devices.UBBs)

	clients, err := m.GetClientTraffic(nil, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, clients)

	mac_clients, err := m.GetClientTrafficByMac(nil, nil, true, "00:00:00:00:00:00")
	require.NoError(t, err)
	require.NotEmpty(t, mac_clients)

	countries, err := m.GetCountryTraffic(nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, countries)
}
