package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestIntegrationDeviceStats(t *testing.T) {
	t.Parallel()

	var s unifi.IntegrationDeviceStats

	require.NoError(t, gofakeit.Struct(&s))
}

func TestIntegrationDeviceRadioStats(t *testing.T) {
	t.Parallel()

	var s unifi.IntegrationDeviceRadioStats

	require.NoError(t, gofakeit.Struct(&s))
}

func TestIntegrationDeviceUplinkStats(t *testing.T) {
	t.Parallel()

	var s unifi.IntegrationDeviceUplinkStats

	require.NoError(t, gofakeit.Struct(&s))
}
