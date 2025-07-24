package azuredevopsreceiver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	require.NotNil(t, cfg)
	_, ok := cfg.(*Config)
	require.True(t, ok)
}
