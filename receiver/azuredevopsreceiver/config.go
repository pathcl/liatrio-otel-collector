package azuredevopsreceiver

import (
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for the Azure DevOps receiver.
type Config struct {
	confighttp.ServerConfig `mapstructure:",squash"`
	Path                    string `mapstructure:"path"`
}
