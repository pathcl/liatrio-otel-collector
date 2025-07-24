package azuredevopsreceiver

import (
	"context"
	"net/http"
)

// azureDevOpsReceiver is the main receiver struct.
type azureDevOpsReceiver struct {
	// TODO: Add fields for configuration, logger, etc.
}

// ServeHTTP handles incoming Azure DevOps webhook events.
func (r *azureDevOpsReceiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: Parse and handle Azure DevOps events
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Azure DevOps receiver not implemented yet"))
}

// Start starts the receiver.
func (r *azureDevOpsReceiver) Start(ctx context.Context, host interface{}) error {
	// TODO: Start HTTP server or register handler
	return nil
}

// Shutdown stops the receiver.
func (r *azureDevOpsReceiver) Shutdown(ctx context.Context) error {
	// TODO: Clean up resources
	return nil
}
