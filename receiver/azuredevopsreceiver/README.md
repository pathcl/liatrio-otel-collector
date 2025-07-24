# Azure DevOps Receiver

This receiver ingests Azure DevOps Service Hook events and converts them into OpenTelemetry traces. It is modeled after the GitHub Actions receiver, but is tailored for Azure DevOps event payloads and semantics.

## Supported Event Types

- Build completed (`build.complete`)
- Release events (e.g., `release.created`, `release.deployment.completed`)
- Pipeline run state changed (`run.state.changed`)
- Pull request events (e.g., `pullrequest.created`, `pullrequest.updated`)
- Work item events (e.g., `workitem.created`, `workitem.updated`)

See [Azure DevOps Service Hook Events](https://learn.microsoft.com/en-us/azure/devops/service-hooks/events?view=azure-devops) for a full list of event types and payloads.

## Usage

Configure the receiver to accept webhook payloads from Azure DevOps. You can use processors to filter and transform events as needed.

## Example Configuration

```yaml
receivers:
  azuredevops:
    endpoint: 0.0.0.0:8080
    path: /azuredevops/webhook
```

## Development

This receiver is based on the structure and logic of the GitHub Actions receiver. See the code and tests for implementation details.

## Scaffolded Files

- `config.go`, `config_test.go`: Receiver configuration
- `factory.go`, `factory_test.go`: Factory registration and tests
- `trace_event_handling.go`, `trace_event_handling_test.go`: Event parsing and model
- `trace_receiver.go`: HTTP handler and receiver skeleton
- `metadata.yaml`, `Makefile`, `go.mod`, `doc.go`: Metadata and build files

---

*This is a work in progress. Contributions and feedback are welcome!* 