package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestCountry(t *testing.T) {
	t.Parallel()

	var s unifi.Country

	require.NoError(t, gofakeit.Struct(&s))
}
