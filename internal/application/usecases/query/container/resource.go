package container

import (
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateCpuUsageQuery create query for cpu usage of a deployment
// MinBy is used instead of SumBy to account for container restarts
// When container is recreated, the metrics for old container is reported for few minutes even after killed.
// Thus, min by is used to record the newly created container's metrics
func CreateCpuUsageQuery(filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	usage := promql.NewQuery(ContainerCpuUsageSeconds.AsString()).
		Filter(filters).
		Rate(rateDuration).
		MinBy([]string{"pod"})

	limit := limitQuery(append(filters, promql.NewFilter("resource", "=", "cpu")))

	return usage.Divide(limit).SetName("cpu_usage_ratio")
}

// CreateMemoryUsageQuery create query for memory usage of a deployment
// MinBy is used instead of SumBy to account for container restarts
// When container is recreated, the metrics for old container is reported for few minutes even after killed.
// Thus, min by is used to record the newly created container's metrics
func CreateMemoryUsageQuery(filters []promql.Filter) *promql.Query {
	usage := promql.NewQuery(ContainerMemoryWorkingSetBytes.AsString()).
		Filter(filters).
		MinBy([]string{"pod"})

	limit := limitQuery(append(filters, promql.NewFilter("resource", "=", "memory")))

	return usage.Divide(limit).SetName("memory_usage_ratio")
}

// limitQuery create query for resource limit.
// type of resource should be specified via filters
func limitQuery(filters []promql.Filter) *promql.Query {
	return promql.NewQuery(KubePodContainerLimit.AsString()).
		Filter(filters).
		SumBy([]string{"pod"})
}
