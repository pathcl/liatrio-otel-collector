package azuredevopsreceiver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEventToTraces_AllEventTypes(t *testing.T) {
	tests := []struct {
		desc      string
		file      string
		eventType AzureDevOpsEventType
		spanName  string
	}{
		{
			desc:      "Build Completed",
			file:      "testdata/build-complete.json",
			eventType: EventBuildCompleted,
			spanName:  "AzureDevOps Build Completed",
		},
		{
			desc:      "Release Created",
			file:      "testdata/release-created.json",
			eventType: EventReleaseCreated,
			spanName:  "AzureDevOps Release Created",
		},
		{
			desc:      "Release Deployment Completed",
			file:      "testdata/release-deployment-completed.json",
			eventType: EventReleaseDeploymentCompleted,
			spanName:  "AzureDevOps Release Deployment Completed",
		},
		{
			desc:      "Pipeline Run State Changed",
			file:      "testdata/pipeline-run-state-changed.json",
			eventType: EventPipelineRunStateChanged,
			spanName:  "AzureDevOps Pipeline Run State Changed",
		},
		{
			desc:      "Pull Request Created",
			file:      "testdata/pullrequest-created.json",
			eventType: EventPullRequestCreated,
			spanName:  "AzureDevOps Pull Request",
		},
		{
			desc:      "Work Item Created",
			file:      "testdata/workitem-created.json",
			eventType: EventWorkItemCreated,
			spanName:  "AzureDevOps Work Item",
		},
		{
			desc:      "Advanced Security Alert State Changed",
			file:      "testdata/advsec-alert-state-changed.json",
			eventType: EventAdvancedSecurityAlertStateChanged,
			spanName:  "AzureDevOps Security Alert",
		},
	}

	logger := zap.NewNop()

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			payload, err := os.ReadFile(tc.file)
			require.NoError(t, err)

			eventType, resource, err := ParseAzureDevOpsEvent(payload)
			require.NoError(t, err)
			require.Equal(t, tc.eventType, eventType)

			traces, err := eventToTraces(eventType, resource, logger)
			require.NoError(t, err)
			rs := traces.ResourceSpans().At(0)
			ss := rs.ScopeSpans().At(0)
			require.Equal(t, 1, ss.Spans().Len())
			span := ss.Spans().At(0)
			require.Equal(t, tc.spanName, span.Name())
		})
	}
}
