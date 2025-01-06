package unifi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseClientHistory(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	validClientHistoryJSON := `[ { "blocked": false, "channel": 0, "display_name": "Device1", "fingerprint": { "has_override": false }, "first_seen": 0, "id": "4a:29:ca:c4:41:aa", "is_allowed_in_visual_programming": true, "is_guest": false, "is_mlo": false, "is_wired": false, "last_ip": "192.168.1.10", "last_radio": "na", "last_seen": 1730247156, "last_uplink_mac": "2e:83:d8:77:fe:3f", "last_uplink_name": "U6 Lite", "local_dns_record_enabled": false, "mac": "4a:29:ca:c4:41:aa", "name": "Device 1", "noted": true, "oui": "Cloud Network Device Manufacturer", "site_id": "1e1007989d059c14580b8403", "status": "online", "tags": [], "type": "WIRELESS", "unifi_device": false, "uplink_mac": "", "use_fixedip": false, "user_id": "6743b053dcec90060a653200", "usergroup_id": "", "virtual_network_override_enabled": false, "wlanconf_id": "5e0606700dd59c14580b80e4" }, { "blocked": false, "channel": 0, "display_name": "UnifiDevice1", "fingerprint": { "has_override": false }, "first_seen": 1664090137, "hostname": "UnifiDevice1", "id": "59:29:43:c4:2a:ef", "is_allowed_in_visual_programming": true, "is_guest": false, "is_mlo": false, "is_wired": false, "last_ip": "192.168.1.6", "last_seen": 1710202100, "local_dns_record_enabled": false, "mac": "59:29:43:c4:2a:ef", "noted": false, "oui": "Ubiquiti Inc.", "site_id": "1e1007989d059c14580b8403", "status": "offline", "tags": [], "type": "WIRELESS", "unifi_device": true, "unifi_device_info": { "icon_filename": "aea0b2d0-63b7-4a11-9055-d4278f7430aa", "icon_resolutions": [ [ 25, 25 ], [ 51, 51 ], [ 101, 101 ], [ 129, 129 ], [ 257, 257 ] ], "view_in_application": true }, "uplink_mac": "", "use_fixedip": false, "user_id": "6743b053dcec90060a653200", "usergroup_id": "", "virtual_network_override_enabled": false, "wlanconf_id": "5e0606700dd59c14580b80e4" } ]`

	clientHistory := []*ClientHistory{}
	err := json.Unmarshal([]byte(validClientHistoryJSON), &clientHistory)
	a.NoError(err)

	a.Len(clientHistory, 2)
	a.EqualValues(clientHistory[0].DisplayName, "Device1")
	a.EqualValues(clientHistory[0].ID, "4a:29:ca:c4:41:aa")
	a.EqualValues(clientHistory[1].WlanconfID, "5e0606700dd59c14580b80e4")
	a.EqualValues(clientHistory[1].UnifiDevice, FlexBool{Val: true, Txt: "true"})
}

func TestClientHistoryOpts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		actual   *ClientHistoryOpts
		expected *ClientHistoryOpts
		name     string
	}{
		{
			name:   "Default Options",
			actual: NewClientHistoryOpts(),
			expected: &ClientHistoryOpts{
				OnlyNonBlocked:      false,
				IncludeUnifiDevices: true,
				WithinHours:         0,
			},
		},
		{
			name:   "OnlyNonBlocked",
			actual: NewClientHistoryOpts().SetOnlyNonBlocked(true),
			expected: &ClientHistoryOpts{
				OnlyNonBlocked:      true,
				IncludeUnifiDevices: true,
				WithinHours:         0,
			},
		},
		{
			name:   "IncludeUnifiDevices",
			actual: NewClientHistoryOpts().SetIncludeUnifiDevices(false),
			expected: &ClientHistoryOpts{
				OnlyNonBlocked:      false,
				IncludeUnifiDevices: false,
				WithinHours:         0,
			},
		},
		{
			name:   "WithinHours",
			actual: NewClientHistoryOpts().SetWithinHours(24),
			expected: &ClientHistoryOpts{
				OnlyNonBlocked:      false,
				IncludeUnifiDevices: true,
				WithinHours:         24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual)
		})
	}
}
