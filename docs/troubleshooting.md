# Troubleshooting

Recommendations below are surfaced automatically by `kubediag diagnose` / `kubediag nodes` /
`kubediag pods`; this page has more detail than fits in a terminal report.

## Node Not Ready

The node's kubelet is not reporting `Ready`. Check kubelet logs on the node
(`journalctl -u kubelet`) and confirm it can reach the API server. If it doesn't recover, drain
and replace the node.

## Memory Pressure

The node is close to its memory limit. Evict or reschedule non-critical pods, or increase the
node's memory capacity.

## Disk Pressure

The node is low on disk space (often container images/logs). Prune unused images
(`crictl rmi --prune`) or increase disk capacity.

## PID Pressure

The node is close to its process ID limit, usually from a container leaking processes.
Investigate the offending pod or raise the node's `pid` limit.

## Pending Pods

`kubediag` inspects the pod's `PodScheduled` condition message to classify the cause:

- **Insufficient CPU / memory** — increase node pool size, reduce the pod's requests, or enable
  the Cluster Autoscaler.
- **Unschedulable (other)** — check taints, tolerations, node selectors, and affinity rules.

## CrashLoopBackOff

The container starts and exits repeatedly. Check `kubectl logs <pod> -c <container> --previous`
for the exit reason.

## ImagePullBackOff

The image cannot be pulled. Verify the image name/tag exists and, for private registries, that
`imagePullSecrets` are configured on the pod or service account.

## OOMKilled

The container exceeded its memory limit and was killed by the kernel. Increase the container's
memory limit or investigate a memory leak.

## Failed Pods

Check `kubectl describe pod <pod>` for events and the container logs for the termination reason.
