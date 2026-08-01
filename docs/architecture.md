# Architecture

## Overall flow

```
CLI -> Kubernetes API -> Collectors -> Analysis Engine -> Recommendation Engine -> Report Generator -> Terminal / JSON / HTML
```

## Package layout

- `main.go` — entry point, delegates to `cmd.Execute()`.
- `cmd/` — one package per CLI command, each exposing `NewCommand() *cobra.Command`. `cmd/root.go`
  wires them onto the root command and defines global (persistent) flags.
- `internal/config` — global flag values (`--kubeconfig`, `--context`, `-o/--output`), read by
  every subcommand.
- `internal/cluster` — loads kubeconfig and builds a `client-go` clientset; verifies API server
  connectivity.
- `internal/<module>` (`nodes`, `pods`, ...) — one package per diagnostic module. Each exposes a
  `Check(ctx, clientset) ([]reporter.Issue, error)` function that collects cluster state and
  evaluates it against known-bad conditions.
- `internal/recommendations` — maps a root-cause `Code` to human-readable advice and a docs link.
  Diagnostic modules look up recommendations here instead of hardcoding advice inline, so the
  mapping is centralized and easy to extend.
- `internal/reporter` — the `Issue`/`Check`/`Report` types shared by every module, cluster score
  calculation, and terminal/JSON/HTML rendering.
- `internal/diagnose` — orchestrates a full run: connect, run every v1.0 module, aggregate into a
  `reporter.Report`. Used by both `kubediag diagnose` and `kubediag report`.

## Diagnostic modules

| Module | Status | Checks |
|---|---|---|
| Cluster Health | v1.0 (implemented) | API server reachability, version |
| Nodes | v1.0 (implemented) | Ready, MemoryPressure, DiskPressure, PIDPressure |
| Pods | v1.0 (implemented) | CrashLoopBackOff, ImagePullBackOff, Pending (root cause), Failed, OOMKilled |
| Deployments | planned | Replica mismatch, failed rollout, stale ReplicaSets, restart count |
| Networking | v2.0 | DNS/CoreDNS, Services, endpoints, CNI |
| Storage | v2.0 | Pending PVC, failed PV, mount failures |
| Ingress | v2.0 | Missing TLS/backend/service, invalid host |
| Resource Usage | planned | CPU/memory requests, limits, overcommit (metrics API) |
| Autoscaling | v3.0 | Missing/misconfigured HPA, scaling failures |
| Security | v3.0 | RBAC, privileged pods, root containers, missing NetworkPolicy, secret exposure |

## Recommendation engine

```
Problem -> Root Cause Detection -> Best Practice Mapping -> Suggested Fix -> Documentation Link
```

Each diagnostic module classifies *why* a problem occurred (e.g. a Pending pod's `PodScheduled`
condition message is parsed to distinguish "insufficient CPU" from "insufficient memory" from a
generic scheduling constraint) and looks up the matching `internal/recommendations.Code`.

## Adding a new diagnostic module

1. Create `internal/<module>` with a `Check(ctx, clientset) ([]reporter.Issue, error)`.
2. Add any new root causes to `internal/recommendations`.
3. Wire it into `internal/diagnose.Run` (for `kubediag diagnose`/`report`) and, if it should be
   runnable standalone, add a `cmd/<module>` package and register it in `cmd/root.go`.
