package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestRADIUSProfileMetadata(t *testing.T) {
	t.Parallel()

	var s unifi.RADIUSProfileMetadata

	require.NoError(t, gofakeit.Struct(&s))
}

func TestRADIUSProfile(t *testing.T) {
	t.Parallel()

	var s unifi.RADIUSProfile

	require.NoError(t, gofakeit.Struct(&s))
}
