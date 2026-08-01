// Package recommendations maps known root causes to actionable fixes and
// documentation links. It is the "best practice mapping" step of the
// recommendation engine described in docs/architecture.md.
package recommendations

// Recommendation pairs human-readable advice with a docs link for a known
// root cause.
type Recommendation struct {
	Text     string
	DocsLink string
}

// Code identifies a known root cause. Diagnostic modules look these up
// instead of hardcoding advice inline, so the mapping can be tuned in one
// place.
type Code string

const (
	NodeNotReady                 Code = "node.not_ready"
	NodeMemoryPressure           Code = "node.memory_pressure"
	NodeDiskPressure             Code = "node.disk_pressure"
	NodePIDPressure              Code = "node.pid_pressure"
	PodPendingInsufficientCPU    Code = "pod.pending.insufficient_cpu"
	PodPendingInsufficientMemory Code = "pod.pending.insufficient_memory"
	PodPendingUnschedulable      Code = "pod.pending.unschedulable"
	PodCrashLoopBackOff          Code = "pod.crash_loop_backoff"
	PodImagePullBackOff          Code = "pod.image_pull_backoff"
	PodOOMKilled                 Code = "pod.oom_killed"
	PodFailed                    Code = "pod.failed"
)

var catalog = map[Code]Recommendation{
	NodeNotReady: {
		Text:     "Check kubelet logs on the node and verify it can reach the API server. Consider draining and replacing the node if it does not recover.",
		DocsLink: "docs/troubleshooting.md#node-not-ready",
	},
	NodeMemoryPressure: {
		Text:     "Evict or reschedule non-critical pods, or increase node memory capacity.",
		DocsLink: "docs/troubleshooting.md#memory-pressure",
	},
	NodeDiskPressure: {
		Text:     "Clean up unused images/containers on the node or increase disk capacity.",
		DocsLink: "docs/troubleshooting.md#disk-pressure",
	},
	NodePIDPressure: {
		Text:     "Investigate processes/containers leaking PIDs on the node, or raise the node's PID limit.",
		DocsLink: "docs/troubleshooting.md#pid-pressure",
	},
	PodPendingInsufficientCPU: {
		Text:     "Increase node pool size, reduce pod CPU requests, or enable the Cluster Autoscaler.",
		DocsLink: "docs/troubleshooting.md#pending-pods",
	},
	PodPendingInsufficientMemory: {
		Text:     "Increase node pool size, reduce pod memory requests, or enable the Cluster Autoscaler.",
		DocsLink: "docs/troubleshooting.md#pending-pods",
	},
	PodPendingUnschedulable: {
		Text:     "Check taints, tolerations, node selectors, and affinity rules preventing this pod from scheduling.",
		DocsLink: "docs/troubleshooting.md#pending-pods",
	},
	PodCrashLoopBackOff: {
		Text:     "Inspect the container logs and previous termination reason; the application is exiting shortly after start.",
		DocsLink: "docs/troubleshooting.md#crashloopbackoff",
	},
	PodImagePullBackOff: {
		Text:     "Verify the image name/tag exists and that imagePullSecrets are configured correctly for private registries.",
		DocsLink: "docs/troubleshooting.md#imagepullbackoff",
	},
	PodOOMKilled: {
		Text:     "Increase the container's memory limit or fix a memory leak in the application.",
		DocsLink: "docs/troubleshooting.md#oomkilled",
	},
	PodFailed: {
		Text:     "Inspect the pod's events and container logs to determine why it terminated in a Failed state.",
		DocsLink: "docs/troubleshooting.md#failed-pods",
	},
}

// For returns the recommendation for a known code. If the code is not in
// the catalog, ok is false.
func For(code Code) (Recommendation, bool) {
	rec, ok := catalog[code]
	return rec, ok
}
