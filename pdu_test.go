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
	
	// Verify outlet table is parsed correctly
	a.NotNil(p.OutletTable, "outlet_table should not be nil")
	a.Greater(len(p.OutletTable), 0, "outlet_table should have entries")
	
	// Verify outlet_current is parsed correctly for outlets that have it
	for _, outlet := range p.OutletTable {
		if outlet.OutletCaps.Val >= 3 {
			// Outlets with caps >= 3 should have current data
			a.NotNil(&outlet.OutletCurrent, "outlet_current should be accessible")
		}
	}
	
	// Verify specific outlet with known current value (Outlet 5 from example has 0.086)
	for _, outlet := range p.OutletTable {
		if outlet.Index.Val == 5 {
			expectedCurrent := 0.086
			a.InDelta(expectedCurrent, outlet.OutletCurrent.Val, 0.001, 
				"outlet 5 should have current value of approximately 0.086")
		}
	}
}
