package unifi_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v5"
	"testing"
)

func TestUCI(test *testing.T) {
	uci := unifi.UCI{}
	err := gofakeit.Struct(&uci)
	require.NoError(test, err)
	require.NotEmpty(test, uci.Name)
}
