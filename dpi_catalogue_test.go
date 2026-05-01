package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestDPIApplication(t *testing.T) {
	t.Parallel()

	var s unifi.DPIApplication

	require.NoError(t, gofakeit.Struct(&s))
}

func TestDPICategory(t *testing.T) {
	t.Parallel()

	var s unifi.DPICategory

	require.NoError(t, gofakeit.Struct(&s))
}
