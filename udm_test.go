package unifi // nolint: testpackage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUDMWANHealthParse decodes per-WAN latency, availability and uptime_stats
// from a real UDM Pro (firmware 5.1.33, Network 10.6.106) stat/device entry.
func TestUDMWANHealthParse(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("endpoints_data/stat-device-udm.json")
	require.NoError(t, err)

	u := &Unifi{Config: &Config{URL: "http://unifi", DebugLog: discardLogs, ErrorLog: discardLogs}}
	devices := u.parseDevices([]json.RawMessage{raw}, &Site{SiteName: "Default"})

	require.Len(t, devices.UDMs, 1)

	d := devices.UDMs[0]
	a := assert.New(t)
	a.Equal("eth9", d.Wan1.Ifname)
	a.InDelta(30, d.Wan1.Latency.Val, 0)
	a.InDelta(100, d.Wan1.Availability.Val, 0)
	a.Equal("eth8", d.Wan2.Ifname)
	a.InDelta(21, d.Wan2.Latency.Val, 0)
	a.InDelta(100, d.Wan2.Availability.Val, 0)

	require.Len(t, d.UptimeStats, 2)
	require.Contains(t, d.UptimeStats, "WAN")
	require.Contains(t, d.UptimeStats, "WAN2")

	wan := d.UptimeStats["WAN"]
	a.InDelta(99.9306, wan.Availability.Val, 0.0001, "availability is fractional")
	a.InDelta(38, wan.LatencyAverage.Val, 0)
	a.InDelta(86400, wan.TimePeriod.Val, 0)
	a.InDelta(510587, wan.Uptime.Val, 0)

	require.Len(t, wan.Monitors, 3)
	a.Equal(UptimeMonitor{
		Availability:   *NewFlexInt(100),
		LatencyAverage: *NewFlexInt(34),
		Target:         "www.microsoft.com",
		Type:           "icmp",
	}, wan.Monitors[0])

	require.Len(t, wan.AlertingMonitors, 3)
	a.Equal("1.1.1.1", wan.AlertingMonitors[1].Target)
	a.Equal("dns", wan.AlertingMonitors[1].Type)
	a.InDelta(36, wan.AlertingMonitors[1].LatencyAverage.Val, 0)

	wan2 := d.UptimeStats["WAN2"]
	a.InDelta(100, wan2.Availability.Val, 0)
	a.InDelta(22, wan2.LatencyAverage.Val, 0)
	a.InDelta(1425562, wan2.Uptime.Val, 0)
	a.Len(wan2.Monitors, 3)
	a.Len(wan2.AlertingMonitors, 3)
}
