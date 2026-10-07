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
	a.InDelta(637334, r.CaNr[0].DlArfcn.Val, 0)
	a.InDelta(40, r.CaNr[0].DlBwMhz.Val, 0)
	a.InDelta(0, r.CaNr[1].UlBwMhz.Val, 0, "null uplink bandwidth decodes as zero")
	a.Empty(r.CaLte)

	require.Len(t, d.Mbb.Sim, 2)
	a.Equal("no-sim", d.Mbb.Sim[0].DisplayState)

	sim := d.Mbb.Sim[1]
	a.True(sim.Active.Val)
	a.Equal("operational", sim.DisplayState)
	a.Equal("Carrier rejected data", sim.NetworkRejectText)
	a.InDelta(3013150278176, sim.RxBytes.Val, 0, "byte counters arrive as strings")
	a.InDelta(725650860811, sim.TxBytes.Val, 0)
	a.Empty(sim.NetworkReject, "this firmware uses the split reject fields")
}

func TestMBBRadioLTECarrierAndSNR(t *testing.T) {
	t.Parallel()

	const raw = `{
		"ca_lte": [{"band": 20, "dl_bw_mhz": 10, "dl_earfcn": 6400, "primary": true, "ul_bw_mhz": 10, "ul_earfcn": 18400}],
		"ca_nr": [{"band": 78, "dl_arfcn": 637334, "dl_bw_mhz": 40, "primary": true, "ul_arfcn": 637334, "ul_bw_mhz": 40}],
		"snr": 4.400000095367432,
		"snr_nr": 17.6
	}`

	var radio MBBRadio
	require.NoError(t, json.Unmarshal([]byte(raw), &radio))

	a := assert.New(t)
	require.Len(t, radio.CaLte, 1)
	a.InDelta(20, radio.CaLte[0].Band.Val, 0)
	a.InDelta(6400, radio.CaLte[0].DlEarfcn.Val, 0)
	a.InDelta(18400, radio.CaLte[0].UlEarfcn.Val, 0)
	a.InDelta(10, radio.CaLte[0].DlBwMhz.Val, 0)
	a.True(radio.CaLte[0].Primary.Val)
	a.InDelta(4.4, radio.Snr.Val, 0.001)
	a.InDelta(17.6, radio.SnrNr.Val, 0.001)

	require.Len(t, radio.CaNr, 1)
	a.InDelta(637334, radio.CaNr[0].DlArfcn.Val, 0)
}

func TestMBBOverridesSimPlan(t *testing.T) {
	t.Parallel()

	const raw = `{
		"primary_slot": 1,
		"sim": [{
			"slot": 1,
			"data_limit_enabled": true,
			"data_soft_limit_bytes": 25000000000,
			"data_soft_limit_display_unit": "GB",
			"data_warning_threshold": 90,
			"reset_policy": "month",
			"reset_date": 1
		}]
	}`

	var overrides MBBOverrides
	require.NoError(t, json.Unmarshal([]byte(raw), &overrides))

	a := assert.New(t)
	a.InDelta(1, overrides.PrimarySlot.Val, 0)
	require.Len(t, overrides.Sim, 1)

	sim := overrides.Sim[0]
	a.True(sim.DataLimitEnabled.Val)
	a.InDelta(25000000000, sim.DataSoftLimitBytes.Val, 0)
	a.Equal("GB", sim.DataSoftLimitDisplayUnit)
	a.InDelta(90, sim.DataWarningThreshold.Val, 0)
	a.Equal("month", sim.ResetPolicy)
	a.InDelta(1, sim.ResetDate.Val, 0)
	a.InDelta(1, sim.Slot.Val, 0)
}

func TestMBBSimNetworkRejectString(t *testing.T) {
	t.Parallel()

	var sim MBBSim
	require.NoError(t, json.Unmarshal([]byte(`{"network_reject":"ip","slot":1}`), &sim))

	a := assert.New(t)
	a.Equal("ip", sim.NetworkReject)
	a.Empty(sim.NetworkRejectType)
	a.Empty(sim.NetworkRejectText)
	a.InDelta(0, sim.NetworkRejectAge.Val, 0)
}
