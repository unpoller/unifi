package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestWANStatusStruct(t *testing.T) {
	t.Parallel()

	var s unifi.WANStatus

	err := gofakeit.Struct(&s)
	require.NoError(t, err)
}

func TestWANStatusInterfaceStruct(t *testing.T) {
	t.Parallel()

	var iface unifi.WANStatusInterface

	err := gofakeit.Struct(&iface)
	require.NoError(t, err)
	require.NotEmpty(t, iface.Name)
}
