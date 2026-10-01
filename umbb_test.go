package unifi // nolint: testpackage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUMBB(t *testing.T) {
	t.Parallel()

	umbb := UMBB{}
	require.NoError(t, gofakeit.Struct(&umbb))
	require.NotEmpty(t, umbb.Name)
}

// TestUMBBParse decodes a real U5G Max (UMBBE630) device entry from stat/device.
func TestUMBBParse(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("endpoints_data/stat-device-umbb.json")
	require.NoError(t, err)

	u := &Unifi{Config: &Config{URL: "http://unifi", DebugLog: discardLogs, ErrorLog: discardLogs}}
	devices := u.parseDevices([]json.RawMessage{raw}, &Site{SiteName: "Default"})

	require.Len(t, devices.UMBBs, 1)

	d := devices.UMBBs[0]
	a := assert.New(t)
	a.Equal("U5G Max", d.Name)
	a.Equal("UMBBE630", d.Model)
	a.Equal("Default", d.SiteName)
	a.Equal("failover", d.Mbb.Mode)
	a.Equal("ready", d.Mbb.State)
	a.InDelta(1, d.MbbOverrides.PrimarySlot.Val, 0)

	r := d.Mbb.Radio
	a.True(r.SA5GMode.Val)
	a.Equal("5G", r.Rat)
	a.Equal("n78", r.Band)
	a.Equal("EE", r.NetworkOperator)
	a.InDelta(-96, r.Rsrp.Val, 0)
	a.InDelta(-13, r.Rsrq.Val, 0)
	a.InDelta(17.6, r.Snr.Val, 0.001)
	a.InDelta(95, r.SignalPercent.Val, 0)
	require.Len(t, r.CaNr, 2)
	a.True(r.CaNr[0].Primary.Val)
	a.InDelta(40, r.CaNr[0].DlBwMhz.Val, 0)
	a.InDelta(0, r.CaNr[1].UlBwMhz.Val, 0, "null uplink bandwidth decodes as zero")

	require.Len(t, d.Mbb.Sim, 2)
	a.Equal("no-sim", d.Mbb.Sim[0].DisplayState)

	sim := d.Mbb.Sim[1]
	a.True(sim.Active.Val)
	a.Equal("operational", sim.DisplayState)
	a.Equal("Carrier rejected data", sim.NetworkRejectText)
	a.InDelta(3013150278176, sim.RxBytes.Val, 0, "byte counters arrive as strings")
	a.InDelta(725650860811, sim.TxBytes.Val, 0)
}
