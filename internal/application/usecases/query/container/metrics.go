package container

type MetricsName string

const (
	ContainerCpuUsageSeconds          MetricsName = "container_cpu_usage_seconds_total"
	ContainerCpuThrottledPeriodsTotal MetricsName = "container_cpu_cfs_throttled_periods_total"
	ContainerMemoryWorkingSetBytes    MetricsName = "container_memory_working_set_bytes"
	KubePodContainerLimit             MetricsName = "kube_pod_container_resource_limits"
	KubePodContainerRestarts          MetricsName = "kube_pod_container_status_restarts_total"
)

func (m MetricsName) AsString() string {
	return string(m)
}
