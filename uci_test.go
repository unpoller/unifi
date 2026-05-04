package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestUCI(test *testing.T) {
	uci := unifi.UCI{}
	err := gofakeit.Struct(&uci)
	require.NoError(test, err)
	require.NotEmpty(test, uci.Name)
}

// TestUCISystemStatsUnmarshal verifies that the system-stats field (hyphen, as
// returned by the UniFi controller API) correctly unmarshals into SystemStats.
func TestUCISystemStatsUnmarshal(t *testing.T) {
	data := []byte(`{
		"sys_stats":    {"loadavg_1": "8.21", "loadavg_5": "8.18", "loadavg_15": "5.44", "mem_buffer": 2572288, "mem_total": 258965504, "mem_used": 231374848},
		"system-stats": {"cpu": "15.8", "mem": "89.3", "uptime": "711050"}
	}`)

	var uci unifi.UCI
	require.NoError(t, json.Unmarshal(data, &uci))
	require.NotNil(t, uci.SysStats, "SysStats should not be nil")
	require.NotNil(t, uci.SystemStats, "SystemStats should not be nil (check JSON tag: system-stats not system_stats)")
	require.Equal(t, "15.8", uci.SystemStats.CPU.String())
	require.Equal(t, "89.3", uci.SystemStats.Mem.String())
}
