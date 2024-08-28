package query

import (
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateCpuUsageQuery create query for cpu usage of a deployment
func CreateCpuUsageQuery(filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	usage := promql.NewQuery("container_cpu_usage_seconds_total").
		Filter(filters).
		Rate(rateDuration).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})

	limit := limitQuery(append(filters, promql.NewFilter("resource", "=", "cpu")))

	return usage.Divide(limit).SetName("cpu_usage_ratio")
}

// CreateMemoryUsageQuery create query for memory usage of a deployment
func CreateMemoryUsageQuery(filters []promql.Filter) *promql.Query {
	usage := promql.NewQuery("container_memory_working_set_bytes").
		Filter(filters).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})

	limit := limitQuery(append(filters, promql.NewFilter("resource", "=", "memory")))

	return usage.Divide(limit).SetName("memory_usage_ratio")
}

// limitQuery create query for resource limit of
// a resource set on specified containers in a namespace
// and sum them by owner deployemnt.
func limitQuery(filters []promql.Filter) *promql.Query {
	return promql.NewQuery("kube_pod_container_resource_limits").
		Filter(filters).
		LabelReplace("deployment", "pod", "(.*)-[^-]+-[^-]+").
		SumBy([]string{"deployment"})
}
