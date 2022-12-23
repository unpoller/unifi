package unifi

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed examples/pdu.json
var pduSample []byte

func TestPDUUnmarshalJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	powerBudgetVal := float64(1875)
	p := &PDU{}
	err := json.Unmarshal(pduSample, p)
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(powerBudgetVal, p.OutletACPowerBudget.Val, "data was not properly unmarshaled")
}
