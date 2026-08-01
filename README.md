# kubediag

A CLI tool that automatically diagnoses Kubernetes clusters, detects common issues, explains what's wrong, and suggests fixes.

```
$ kubediag diagnose

Cluster Status

✓ API Server Healthy
⚠ Pod default/worker-7cf9: Pending
  Reason: Pod default/worker-7cf9 is Pending: insufficient CPU
  Recommendation: Increase node pool size, reduce pod CPU requests, or enable the Cluster Autoscaler.

Cluster Score: 91/100
Health: Excellent
Warnings: 3
Critical: 0
```

## Status

Early scaffold. v1.0 diagnostics (cluster health, nodes, pods) are implemented; everything else in the
[roadmap](docs/roadmap.md) is a stubbed command pointing back at this repo.

## Install

See [docs/installation.md](docs/installation.md).

## Usage

See [docs/commands.md](docs/commands.md) for the full command reference.

```
kubediag diagnose            # full cluster diagnosis
kubediag nodes                # node health only
kubediag pods                 # pod health only
kubediag report --out r.json  # export a report
kubediag version
```

## How it works

```
CLI -> Kubernetes API -> Collectors -> Analysis Engine -> Recommendation Engine -> Report Generator -> Terminal / JSON / HTML
```

See [docs/architecture.md](docs/architecture.md) for the full breakdown of diagnostic modules.

## Development

```
go mod tidy
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

## License

[MIT](LICENSE)
