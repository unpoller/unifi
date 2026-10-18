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

	a.Equal(map[string]string{"WAN": "online", "WAN2": "online"}, d.LastWanStatus)

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

// TestUDMWANOfflineParse decodes the same UDM Pro about a minute after WAN1's
// upstream went dark. The port link stayed up; the controller reports latency
// and availability as 0 and swaps the uptime_stats averages for downtime.
func TestUDMWANOfflineParse(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("endpoints_data/stat-device-udm-wan-offline.json")
	require.NoError(t, err)

	u := &Unifi{Config: &Config{URL: "http://unifi", DebugLog: discardLogs, ErrorLog: discardLogs}}
	devices := u.parseDevices([]json.RawMessage{raw}, &Site{SiteName: "Default"})

	require.Len(t, devices.UDMs, 1)

	d := devices.UDMs[0]
	a := assert.New(t)
	a.Equal(map[string]string{"WAN": "offline", "WAN2": "online"}, d.LastWanStatus)
	a.True(d.Wan1.Up.Val, "the port link stays up")
	a.Equal("0", d.Wan1.Latency.Txt)
	a.Equal("0", d.Wan1.Availability.Txt)

	wan := d.UptimeStats["WAN"]
	require.NotNil(t, wan)
	a.InDelta(66, wan.Downtime.Val, 0)
	a.Empty(wan.Availability.Txt, "averages are omitted while offline")
	a.Empty(wan.LatencyAverage.Txt)
	a.Empty(wan.TimePeriod.Txt)
	a.Empty(wan.Uptime.Txt)

	require.Len(t, wan.Monitors, 3)
	a.Equal("www.microsoft.com", wan.Monitors[0].Target)
	a.Empty(wan.Monitors[0].LatencyAverage.Txt, "a monitor can omit latency_average")
	a.Equal("0", wan.Monitors[1].LatencyAverage.Txt)
	a.InDelta(100, wan.Monitors[1].Availability.Val, 0)

	require.Len(t, wan.AlertingMonitors, 3)
	a.InDelta(0, wan.AlertingMonitors[0].Availability.Val, 0)

	wan2 := d.UptimeStats["WAN2"]
	require.NotNil(t, wan2)
	a.Empty(wan2.Downtime.Txt, "an online WAN has no downtime")
	a.NotEmpty(wan2.Uptime.Txt)
}

// TestUDMWANLinkDownParse decodes the same UDM Pro about a minute after WAN2's
// cable was unplugged at the console. The port goes down and omits latency and
// availability entirely; uptime_stats reports downtime as for any outage.
func TestUDMWANLinkDownParse(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("endpoints_data/stat-device-udm-wan-link-down.json")
	require.NoError(t, err)

	u := &Unifi{Config: &Config{URL: "http://unifi", DebugLog: discardLogs, ErrorLog: discardLogs}}
	devices := u.parseDevices([]json.RawMessage{raw}, &Site{SiteName: "Default"})

	require.Len(t, devices.UDMs, 1)

	d := devices.UDMs[0]
	a := assert.New(t)
	a.Equal(map[string]string{"WAN": "online", "WAN2": "offline"}, d.LastWanStatus)
	a.False(d.Wan2.Up.Val)
	a.Empty(d.Wan2.Latency.Txt, "a down port omits latency")
	a.Empty(d.Wan2.Availability.Txt)

	wan2 := d.UptimeStats["WAN2"]
	require.NotNil(t, wan2)
	a.InDelta(62, wan2.Downtime.Val, 0)
	a.Empty(wan2.Availability.Txt)
	a.Empty(wan2.Uptime.Txt)

	require.Len(t, wan2.Monitors, 3)

	for _, m := range append(wan2.Monitors, wan2.AlertingMonitors...) {
		a.InDelta(0, m.Availability.Val, 0, m.Target)
		a.Empty(m.LatencyAverage.Txt, m.Target)
	}
}
