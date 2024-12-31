package unifi_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unifi/v5"
)

//go:embed examples/pdu.json
var pduSample []byte

func TestPDUUnmarshalJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	powerBudgetVal := float64(1875)
	p := &unifi.PDU{}
	err := json.Unmarshal(pduSample, p)
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(powerBudgetVal, p.OutletACPowerBudget.Val, "data was not properly unmarshaled")
}
