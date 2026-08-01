# Commands

Global flags (available on every command):

| Flag           | Default              | Description                                  |
|----------------|-----------------------|-----------------------------------------------|
| `--kubeconfig` | `$KUBECONFIG` or `~/.kube/config` | Path to kubeconfig file          |
| `--context`    | current context      | kubeconfig context to use                     |
| `-o, --output` | `table`               | Output format: `table`, `json`, `html`        |

## `kubediag diagnose`

Runs every implemented diagnostic module (cluster health, nodes, pods) and prints a scored report.

```
kubediag diagnose
kubediag diagnose -o json
```

## `kubediag nodes`

Runs node checks only: Ready, MemoryPressure, DiskPressure, PIDPressure.

## `kubediag pods`

Runs pod checks only: CrashLoopBackOff, ImagePullBackOff, Pending (with root cause), Failed, OOMKilled.

## `kubediag report`

Runs a full diagnosis and writes it to a file.

```
kubediag report --out report.json
kubediag report --out report.html -o html
```

## `kubediag version`

Prints the CLI version, commit, and build date.

## Planned (see [roadmap.md](roadmap.md))

- `kubediag ingress` — Ingress/TLS/backend diagnostics (v2.0)
- `kubediag storage` — PVC/PV/StorageClass diagnostics (v2.0)
- `kubediag network` — DNS/CoreDNS/Service/NetworkPolicy diagnostics (v2.0)
- `kubediag security` — RBAC/pod security/secret diagnostics (v3.0)
- `kubediag fix <target>` — preview-confirm-apply-verify auto-remediation (v4.0)

These commands exist today as stubs that explain their status; they do not modify the cluster.
