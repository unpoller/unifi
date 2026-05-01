package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
)

func TestPortForwardStruct(t *testing.T) {
	t.Parallel()

	var p unifi.PortForward

	err := gofakeit.Struct(&p)
	require.NoError(t, err)
	require.NotEmpty(t, p.ID)
	require.NotEmpty(t, p.Name)
}
