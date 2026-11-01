package unifi // nolint: testpackage

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed examples/uap.json
var uapSample []byte

func TestUAPUnmarshalJSON(t *testing.T) {
	testcontroller511 := `{
	  "ap": {
	    "site_id": "mySite",
	    "o": "ap",
	    "oid": "00:00:00:00:00:00",
	    "ap": "00:00:00:00:00:00",
			"time": 1577742600000,
      "datetime": "2019-12-30T09:50:00Z",
	    "user-wifi1-rx_packets": 6596670,
	    "user-wifi0-rx_packets": 42659527,
	    "user-rx_packets": 49294197,
	    "guest-rx_packets": 0,
	    "wifi0-rx_packets": 42639527,
	    "wifi1-rx_packets": 6591670,
	    "rx_packets": 49299197}}`

	testcontroller510 := `{
		"site_id": "mySite",
		"o": "ap",
		"oid": "00:00:00:00:00:00",
		"ap": "00:00:00:00:00:00",
		"time": 1577742600000,
		"datetime": "2019-12-30T09:50:00Z",
		"user-wifi1-rx_packets": 6596670,
		"user-wifi0-rx_packets": 42659527,
		"user-rx_packets": 49294197,
		"guest-rx_packets": 0,
		"wifi0-rx_packets": 42639527,
		"wifi1-rx_packets": 6591670,
		"rx_packets": 49299197}`

	t.Parallel()
	a := assert.New(t)
	rxPakcets := 49299197
	u := &UAPStat{}
	err := u.UnmarshalJSON([]byte(testcontroller510))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(rxPakcets), u.RxPackets.Val, "data was not properly unmarshaled")

	u = &UAPStat{} // reset
	err = u.UnmarshalJSON([]byte(testcontroller511))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(rxPakcets), u.RxPackets.Val, "data was not properly unmarshaled")

	uap := &UAP{}
	err = json.Unmarshal(uapSample, uap)
	a.Nil(err, "must be no error unmarshaling uap sample")
	a.Equal(true, uap.Adopted.Val, "data was not properaly unmarshaled")
}

func TestUAPLteState(t *testing.T) {
	t.Parallel()

	var ap UAP
	require.NoError(t, json.Unmarshal([]byte(`{"lte_state":"ready"}`), &ap))
	assert.Equal(t, "ready", ap.LteState)

	var other UAP
	require.NoError(t, json.Unmarshal([]byte(`{"lte_state":"connecting"}`), &other))
	assert.Equal(t, "connecting", other.LteState)
}

func TestUAPWirelessUplinkUnmarshalJSON(t *testing.T) {
	t.Parallel()

	// Trimmed stat/device uplink of a U6 Pro meshed to a wired parent AP.
	data := `{"uplink": {"ap_mac": "60:22:32:00:00:19", "bssid": "76:22:32:00:00:1b",
	  "channel": 44, "radio": "na", "radio_name": "wifi1", "rssi": 12, "signal": -83,
	  "tx_rate": 51600, "rx_rate": 68800, "type": "wireless", "up": true, "port_idx": -1,
	  "uplink_device_name": "Front", "uplink_mac": "60:22:32:00:00:19",
	  "uplink_source": "lldp_uplink"}}`

	var u UAP

	require.NoError(t, json.Unmarshal([]byte(data), &u))
	assert.Equal(t, "wireless", u.Uplink.Type)
	assert.Equal(t, "60:22:32:00:00:19", u.Uplink.UplinkMac)
	assert.Equal(t, "Front", u.Uplink.UplinkDeviceName)
	assert.Equal(t, "60:22:32:00:00:19", u.Uplink.ApMac)
	assert.Equal(t, "na", u.Uplink.Radio)
	assert.Equal(t, "wifi1", u.Uplink.RadioName)
	assert.InDelta(t, 44.0, u.Uplink.Channel.Val, 0)
	assert.InDelta(t, 12.0, u.Uplink.Rssi.Val, 0)
	assert.InDelta(t, -83.0, u.Uplink.Signal.Val, 0)
	assert.InDelta(t, 51600.0, u.Uplink.TxRate.Val, 0)
	assert.InDelta(t, 68800.0, u.Uplink.RxRate.Val, 0)
}
