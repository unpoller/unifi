package unifi_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

func TestACLRule(t *testing.T) {
	t.Parallel()

	var r unifi.ACLRule

	require.NoError(t, gofakeit.Struct(&r))
}
