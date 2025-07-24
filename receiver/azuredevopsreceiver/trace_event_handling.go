package azuredevopsreceiver

import (
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// AzureDevOpsEventType represents the type of Azure DevOps event.
type AzureDevOpsEventType string

// AzureDevOpsWebhookPayload is a generic struct for Azure DevOps webhook events.
type AzureDevOpsWebhookPayload struct {
	EventType string          `json:"eventType"`
	Resource  json.RawMessage `json:"resource"`
	// Add more fields as needed from the top-level event
}

// BuildCompletedResource represents the resource for build.complete events.
type BuildCompletedResource struct {
	ID          int    `json:"id"`
	BuildNumber string `json:"buildNumber"`
	Status      string `json:"status"`
	Result      string `json:"result"`
	StartTime   string `json:"startTime"`
	FinishTime  string `json:"finishTime"`
	Definition  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"definition"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

// ReleaseCreatedResource represents the resource for release.created events.
type ReleaseCreatedResource struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CreatedOn  string `json:"createdOn"`
	ModifiedOn string `json:"modifiedOn"`
	Project    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

// ReleaseDeploymentCompletedResource represents the resource for release.deployment.completed events.
type ReleaseDeploymentCompletedResource struct {
	DeploymentID int `json:"deploymentId"`
	Release      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"release"`
	Status      string `json:"status"`
	StartedOn   string `json:"startedOn"`
	CompletedOn string `json:"completedOn"`
	Project     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

// PipelineRunStateChangedResource represents the resource for run.state.changed events.
type PipelineRunStateChangedResource struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Result     string `json:"result"`
	StartedOn  string `json:"startedOn"`
	FinishedOn string `json:"finishedOn"`
	Project    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

// PullRequestResource represents the resource for pullrequest.created/updated events.
type PullRequestResource struct {
	PullRequestID int    `json:"pullRequestId"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	CreatedBy     struct {
		ID   string `json:"id"`
		Name string `json:"displayName"`
	} `json:"createdBy"`
	CreationDate string `json:"creationDate"`
	ClosedDate   string `json:"closedDate"`
	Repository   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"repository"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

// WorkItemResource represents the resource for workitem.created/updated events.
type WorkItemResource struct {
	ID     int                    `json:"id"`
	Rev    int                    `json:"rev"`
	Fields map[string]interface{} `json:"fields"`
	URL    string                 `json:"url"`
}

// AdvancedSecurityAlertResource represents the resource for advanced security alert events.
type AdvancedSecurityAlertResource struct {
	AlertID       int    `json:"alertId"`
	Severity      string `json:"severity"`
	Title         string `json:"title"`
	AlertType     string `json:"alertType"`
	State         string `json:"state"`
	RepositoryUrl string `json:"repositoryUrl"`
	GitRef        string `json:"gitRef"`
}

// Update event type constants
const (
	EventBuildCompleted                    AzureDevOpsEventType = "build.complete"
	EventReleaseCreated                    AzureDevOpsEventType = "release.created"
	EventReleaseDeploymentCompleted        AzureDevOpsEventType = "release.deployment.completed"
	EventPipelineRunStateChanged           AzureDevOpsEventType = "run.state.changed"
	EventPullRequestCreated                AzureDevOpsEventType = "pullrequest.created"
	EventPullRequestUpdated                AzureDevOpsEventType = "pullrequest.updated"
	EventWorkItemCreated                   AzureDevOpsEventType = "workitem.created"
	EventWorkItemUpdated                   AzureDevOpsEventType = "workitem.updated"
	EventAdvancedSecurityAlertStateChanged AzureDevOpsEventType = "ms.vss-alerts.alert-state-changed-event"
	EventAdvancedSecurityAlertUpdated      AzureDevOpsEventType = "ms.vss-alerts.alert-updated-event"
)

// Update the parser to handle all event types
func ParseAzureDevOpsEvent(data []byte) (AzureDevOpsEventType, interface{}, error) {
	var payload AzureDevOpsWebhookPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	switch AzureDevOpsEventType(payload.EventType) {
	case EventBuildCompleted:
		var resource BuildCompletedResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal build resource: %w", err)
		}
		return EventBuildCompleted, resource, nil
	case EventReleaseCreated:
		var resource ReleaseCreatedResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal release created resource: %w", err)
		}
		return EventReleaseCreated, resource, nil
	case EventReleaseDeploymentCompleted:
		var resource ReleaseDeploymentCompletedResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal release deployment completed resource: %w", err)
		}
		return EventReleaseDeploymentCompleted, resource, nil
	case EventPipelineRunStateChanged:
		var resource PipelineRunStateChangedResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal pipeline run state changed resource: %w", err)
		}
		return EventPipelineRunStateChanged, resource, nil
	case EventPullRequestCreated, EventPullRequestUpdated:
		var resource PullRequestResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal pull request resource: %w", err)
		}
		return AzureDevOpsEventType(payload.EventType), resource, nil
	case EventWorkItemCreated, EventWorkItemUpdated:
		var resource WorkItemResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal work item resource: %w", err)
		}
		return AzureDevOpsEventType(payload.EventType), resource, nil
	case EventAdvancedSecurityAlertStateChanged, EventAdvancedSecurityAlertUpdated:
		var resource AdvancedSecurityAlertResource
		if err := json.Unmarshal(payload.Resource, &resource); err != nil {
			return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("failed to unmarshal advanced security alert resource: %w", err)
		}
		return AzureDevOpsEventType(payload.EventType), resource, nil
	default:
		return AzureDevOpsEventType(payload.EventType), nil, fmt.Errorf("unsupported event type: %s", payload.EventType)
	}
}

// eventToTraces generates OpenTelemetry traces for Azure DevOps events.
func eventToTraces(eventType AzureDevOpsEventType, resource interface{}, logger *zap.Logger) (ptrace.Traces, error) {
	traces := ptrace.NewTraces()
	resourceSpans := traces.ResourceSpans().AppendEmpty()
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()

	switch eventType {
	case EventBuildCompleted:
		build, ok := resource.(BuildCompletedResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Build Completed")
		span.Attributes().PutStr("azuredevops.build.number", build.BuildNumber)
		span.Attributes().PutStr("azuredevops.build.status", build.Status)
		span.Attributes().PutStr("azuredevops.build.result", build.Result)
		span.Attributes().PutStr("azuredevops.build.definition", build.Definition.Name)
		span.Attributes().PutStr("azuredevops.project.name", build.Project.Name)
		start, _ := time.Parse(time.RFC3339, build.StartTime)
		finish, _ := time.Parse(time.RFC3339, build.FinishTime)
		span.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
		span.SetEndTimestamp(pcommon.NewTimestampFromTime(finish))
		if build.Result == "succeeded" {
			span.Status().SetCode(ptrace.StatusCodeOk)
		} else {
			span.Status().SetCode(ptrace.StatusCodeError)
		}
	case EventReleaseCreated:
		release, ok := resource.(ReleaseCreatedResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Release Created")
		span.Attributes().PutStr("azuredevops.release.name", release.Name)
		span.Attributes().PutStr("azuredevops.release.status", release.Status)
		span.Attributes().PutStr("azuredevops.project.name", release.Project.Name)
		created, _ := time.Parse(time.RFC3339, release.CreatedOn)
		modified, _ := time.Parse(time.RFC3339, release.ModifiedOn)
		span.SetStartTimestamp(pcommon.NewTimestampFromTime(created))
		span.SetEndTimestamp(pcommon.NewTimestampFromTime(modified))
		span.Status().SetCode(ptrace.StatusCodeUnset)
	case EventReleaseDeploymentCompleted:
		deploy, ok := resource.(ReleaseDeploymentCompletedResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Release Deployment Completed")
		span.Attributes().PutStr("azuredevops.release.name", deploy.Release.Name)
		span.Attributes().PutStr("azuredevops.deployment.status", deploy.Status)
		span.Attributes().PutStr("azuredevops.project.name", deploy.Project.Name)
		start, _ := time.Parse(time.RFC3339, deploy.StartedOn)
		finish, _ := time.Parse(time.RFC3339, deploy.CompletedOn)
		span.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
		span.SetEndTimestamp(pcommon.NewTimestampFromTime(finish))
		if deploy.Status == "succeeded" {
			span.Status().SetCode(ptrace.StatusCodeOk)
		} else {
			span.Status().SetCode(ptrace.StatusCodeError)
		}
	case EventPipelineRunStateChanged:
		run, ok := resource.(PipelineRunStateChangedResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Pipeline Run State Changed")
		span.Attributes().PutStr("azuredevops.pipeline.name", run.Name)
		span.Attributes().PutStr("azuredevops.pipeline.state", run.State)
		span.Attributes().PutStr("azuredevops.pipeline.result", run.Result)
		span.Attributes().PutStr("azuredevops.project.name", run.Project.Name)
		start, _ := time.Parse(time.RFC3339, run.StartedOn)
		finish, _ := time.Parse(time.RFC3339, run.FinishedOn)
		span.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
		span.SetEndTimestamp(pcommon.NewTimestampFromTime(finish))
		if run.Result == "succeeded" {
			span.Status().SetCode(ptrace.StatusCodeOk)
		} else {
			span.Status().SetCode(ptrace.StatusCodeError)
		}
	case EventPullRequestCreated, EventPullRequestUpdated:
		pr, ok := resource.(PullRequestResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Pull Request")
		span.Attributes().PutStr("azuredevops.pullrequest.title", pr.Title)
		span.Attributes().PutStr("azuredevops.pullrequest.status", pr.Status)
		span.Attributes().PutStr("azuredevops.pullrequest.created_by", pr.CreatedBy.Name)
		span.Attributes().PutStr("azuredevops.project.name", pr.Project.Name)
		created, _ := time.Parse(time.RFC3339, pr.CreationDate)
		closed, _ := time.Parse(time.RFC3339, pr.ClosedDate)
		span.SetStartTimestamp(pcommon.NewTimestampFromTime(created))
		if !closed.IsZero() {
			span.SetEndTimestamp(pcommon.NewTimestampFromTime(closed))
		}
		if pr.Status == "completed" {
			span.Status().SetCode(ptrace.StatusCodeOk)
		} else {
			span.Status().SetCode(ptrace.StatusCodeUnset)
		}
	case EventWorkItemCreated, EventWorkItemUpdated:
		wi, ok := resource.(WorkItemResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Work Item")
		span.Attributes().PutInt("azuredevops.workitem.id", int64(wi.ID))
		span.Attributes().PutInt("azuredevops.workitem.rev", int64(wi.Rev))
		span.Attributes().PutStr("azuredevops.workitem.url", wi.URL)
		span.Status().SetCode(ptrace.StatusCodeUnset)
	case EventAdvancedSecurityAlertStateChanged, EventAdvancedSecurityAlertUpdated:
		alert, ok := resource.(AdvancedSecurityAlertResource)
		if !ok {
			return traces, nil
		}
		span := scopeSpans.Spans().AppendEmpty()
		span.SetName("AzureDevOps Security Alert")
		span.Attributes().PutStr("azuredevops.alert.title", alert.Title)
		span.Attributes().PutStr("azuredevops.alert.severity", alert.Severity)
		span.Attributes().PutStr("azuredevops.alert.type", alert.AlertType)
		span.Attributes().PutStr("azuredevops.alert.state", alert.State)
		span.Attributes().PutStr("azuredevops.alert.repository_url", alert.RepositoryUrl)
		span.Attributes().PutStr("azuredevops.alert.git_ref", alert.GitRef)
		span.Status().SetCode(ptrace.StatusCodeUnset)
	default:
		logger.Error("unknown Azure DevOps event type", zap.String("eventType", string(eventType)))
		return traces, nil
	}

	return traces, nil
}

// TODO: Implement event-to-trace logic for supported event types.
