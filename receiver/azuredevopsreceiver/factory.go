package azuredevopsreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver"
)

const (
	TypeStr = "azuredevops"
)

// NewFactory creates a new receiver factory for Azure DevOps.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType(TypeStr),
		createDefaultConfig,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ServerConfig: confighttp.ServerConfig{
			Endpoint: "0.0.0.0:8080",
		},
		Path: "/azuredevops/webhook",
	}
}
