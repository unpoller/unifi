package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestDNSPolicy(t *testing.T) {
	t.Parallel()

	var p unifi.DNSPolicy

	require.NoError(t, gofakeit.Struct(&p))
}
