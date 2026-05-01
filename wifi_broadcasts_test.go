package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestWifiBroadcast(t *testing.T) {
	t.Parallel()

	var s unifi.WifiBroadcast

	require.NoError(t, gofakeit.Struct(&s))
}

func TestWifiBroadcastSecurityConfiguration(t *testing.T) {
	t.Parallel()

	var s unifi.WifiBroadcastSecurityConfiguration

	require.NoError(t, gofakeit.Struct(&s))
}
