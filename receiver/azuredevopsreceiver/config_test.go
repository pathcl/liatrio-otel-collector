package azuredevopsreceiver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := &Config{}
	require.NotNil(t, cfg)
}
