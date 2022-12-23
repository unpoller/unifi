package unifi // nolint: testpackage

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed examples/usw.json
var uswSample []byte

func testGetControllerJSON() (string, string) {
	return `{
		"sw": {
		  "site_id": "mySite",
		  "o": "sw",
		  "oid": "00:00:00:00:00:00",
		  "sw": "00:00:00:00:00:00",
		  "time": 1577742600000,
		  "datetime": "2019-12-30T09:40:00Z",
		  "rx_packets": 321,
		  "rx_bytes": 321,
		  "rx_errors": 123,
		  "rx_dropped": 123,
		  "rx_crypts": 123,
		  "rx_frags": 123,
		  "tx_packets": 123,
		  "tx_bytes": 123,
		  "tx_errors": 0,
		  "tx_dropped": 0,
		  "tx_retries": 0,
		  "rx_multicast": 123,
		  "rx_broadcast": 123,
		  "tx_multicast": 123,
		  "tx_broadcast": 123,
		  "bytes": 123,
		  "duration": 123}}`,
		`{
		"site_id": "mySite",
		"o": "sw",
		"oid": "00:00:00:00:00:00",
		"sw": "00:00:00:00:00:00",
		"time": 1577742600000,
		"datetime": "2019-12-30T09:40:00Z",
		"rx_packets": 321,
		"rx_bytes": 321,
		"rx_errors": 123,
		"rx_dropped": 123,
		"rx_crypts": 123,
		"rx_frags": 123,
		"tx_packets": 123,
		"tx_bytes": 123,
		"tx_errors": 0,
		"tx_dropped": 0,
		"tx_retries": 0,
		"rx_multicast": 123,
		"rx_broadcast": 123,
		"tx_multicast": 123,
		"tx_broadcast": 123,
		"bytes": 123,
		"duration": 123}`
}

func TestUSWUnmarshalJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	testcontroller511, testcontroller510 := testGetControllerJSON()
	rxMulticast := float64(123)
	u := &USWStat{}
	err := u.UnmarshalJSON([]byte(testcontroller510))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(rxMulticast, u.RxMulticast.Val, "data was not properly unmarshaled")

	u = &USWStat{} // reset
	err = u.UnmarshalJSON([]byte(testcontroller511))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(rxMulticast, u.RxMulticast.Val, "data was not properly unmarshaled")

	usw := &USW{}
	err = json.Unmarshal(uswSample, usw)
	a.Nil(err, "must be no error unmarshaling sample")
	a.Equal(true, usw.Adopted.Val, "data was not properly unmarshaled")
}
