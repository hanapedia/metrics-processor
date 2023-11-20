package query

import (
	"time"

	"github.com/hanapedia/metrics-processor/internal/infrastructure/prometheus"
)

// CreateCpuUsageQuery create query for cpu usage of a deployment
func CreateCpuUsageQuery(namespace, containers string, rateDuration time.Duration) *prometheus.Query {
	filters := []prometheus.Filter{
		prometheus.NewFilter("namespace", "=", namespace),
		prometheus.NewFilter("container", "=~", containers),
	}
	usage := prometheus.NewQuery("container_cpu_usage_seconds_total").
		Filter(filters).
		Rate(rateDuration).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})

	limit := limitQuery(namespace, "server|redis", "cpu")

	return usage.Divide(limit).SetName("cpu_usage_ratio")
}

// CreateMemoryUsageQuery create query for memory usage of a deployment
func CreateMemoryUsageQuery(namespace, containers string) *prometheus.Query {
	filters := []prometheus.Filter{
		prometheus.NewFilter("namespace", "=", namespace),
		prometheus.NewFilter("container", "=~", containers),
	}
	usage := prometheus.NewQuery("container_memory_working_set_bytes").
		Filter(filters).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})

	limit := limitQuery(namespace, "server|redis", "memory")

	return usage.Divide(limit).SetName("memory_usage_ratio")
}

// limitQuery create query for resource limit of
// a resource set on specified containers in a namespace
// and sum them by owner deployemnt.
func limitQuery(namespace, containers, resource string) *prometheus.Query {
	filters := []prometheus.Filter{
		prometheus.NewFilter("namespace", "=", namespace),
		prometheus.NewFilter("container", "=~", containers),
		prometheus.NewFilter("resource", "=", resource),
	}
	return prometheus.NewQuery("kube_pod_container_resource_limits").
		Filter(filters).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})
}
