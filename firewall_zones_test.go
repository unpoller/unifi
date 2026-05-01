package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestFirewallZone(t *testing.T) {
	t.Parallel()

	var z unifi.FirewallZone

	require.NoError(t, gofakeit.Struct(&z))
}

func TestFirewallZoneMetadata(t *testing.T) {
	t.Parallel()

	var m unifi.FirewallZoneMetadata

	require.NoError(t, gofakeit.Struct(&m))
}
