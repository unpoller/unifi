package mocks_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5/mocks"
	"testing"
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
}
