package unifi_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
	"testing"
)

// TestUBBIssue821 tests unmarshaling of UBB device with peer_ubb and radio_table_stats.
// This test reproduces the issue reported in https://github.com/unpoller/unpoller/issues/821
// where UBB-XG devices with radio information would fail to unmarshal due to RadioTableStats
// being incorrectly defined as []RadioTableStats (slice of slices) instead of RadioTableStats.
func TestUBBIssue821(t *testing.T) {
	// Simplified JSON from issue #821 - focusing on the radio_table_stats difference
	jsonData := `{
		"type":"ubb",
		"name":"EBES-UBB-XG",
		"mac":"0c:ea:14:3d:a3:25",
		"model":"UBBXG",
		"adopted":true,
		"radio_table_stats":[
			{"name":"wifi0","channel":149,"radio":"na","gain":12,"satisfaction":-1,"state":"RUN","extchannel":0,"tx_power":16,"num_sta":0,"user-num_sta":0},
			{"name":"terra2","channel":4,"radio":"ad","gain":19,"satisfaction":-1,"state":"RUN","extchannel":0,"tx_power":5,"num_sta":0,"user-num_sta":0}
		],
		"peer_ubb":{
			"type":"ubb",
			"name":"EBES-UBB-XG-PEER",
			"mac":"0c:ea:14:3d:a3:8e",
			"model":"UBBXG",
			"adopted":true,
			"radio_table_stats":[
				{"name":"wifi0","channel":149,"radio":"na","gain":12,"satisfaction":-1,"state":"RUN","extchannel":0,"tx_power":16},
				{"name":"terra2","channel":4,"radio":"ad","gain":19,"satisfaction":-1,"state":"RUN","extchannel":0,"tx_power":7}
			]
		}
	}`

	var ubb unifi.UBB

	err := json.Unmarshal([]byte(jsonData), &ubb)
	if err != nil {
		t.Logf("Unmarshal error: %v", err)
	}

	require.NoError(t, err, "Should be able to unmarshal UBB device with peer_ubb")
	require.Equal(t, "EBES-UBB-XG", ubb.Name)
	require.NotNil(t, ubb.PeerUbb, "peer_ubb should be present")
	require.Equal(t, "EBES-UBB-XG-PEER", ubb.PeerUbb.Name)
}
