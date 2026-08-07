package unifi // nolint: testpackage

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed examples/usp.json
var upsTwoUProSample []byte

func TestParseDevicesUPS2UPro(t *testing.T) {
	t.Parallel()

	client := &Unifi{Config: &Config{
		URL:      "https://synthetic-controller.example",
		ErrorLog: discardLogs,
		DebugLog: discardLogs,
	}}
	site := &Site{SiteName: "Synthetic Site"}
	devices := client.parseDevices([]json.RawMessage{upsTwoUProSample}, site)

	require.Empty(t, devices.USWs)
	require.Len(t, devices.PDUs, 1)

	pdu := devices.PDUs[0]
	require.Equal(t, "pdu", pdu.Type)
	require.Equal(t, "USPDA2B", pdu.Model)
	require.Equal(t, "Synthetic UPS-2U-Pro", pdu.Name)
	require.Equal(t, "02:00:00:00:00:01", pdu.Mac)
	require.Equal(t, "Synthetic Site", pdu.SiteName)
	require.Len(t, pdu.OutletTable, 8)

	expectedMetrics := []struct {
		voltage     float64
		current     float64
		power       float64
		powerFactor float64
	}{
		{voltage: 120.4, current: 0.125, power: 15.05, powerFactor: 0.99},
		{voltage: 120.3, current: 0.25, power: 29.8, powerFactor: 0.98},
		{voltage: 120.5, current: 0, power: 0, powerFactor: 0},
		{voltage: 120.2, current: 0.875, power: 103.75, powerFactor: 0.97},
		{voltage: 120.1, current: 1.125, power: 132.4, powerFactor: 0.98},
		{voltage: 120, current: 0.055, power: 5.9, powerFactor: 0.89},
		{voltage: 119.9, current: 0.41, power: 47.2, powerFactor: 0.96},
		{voltage: 120.6, current: 2.015, power: 238.7, powerFactor: 0.99},
	}

	for index, expected := range expectedMetrics {
		outlet := pdu.OutletTable[index]

		require.Equal(t, float64(index+1), outlet.Index.Val)
		require.InDelta(t, expected.voltage, outlet.OutletVoltage.Val, 0.001)
		require.InDelta(t, expected.current, outlet.OutletCurrent.Val, 0.001)
		require.InDelta(t, expected.power, outlet.OutletPower.Val, 0.001)
		require.InDelta(t, expected.powerFactor, outlet.OutletPowerFactor.Val, 0.001)
	}
}
