package query

import (
	"fmt"
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateAvgServerLatencyQuery creates query for average server response time
func CreateAvgServerLatencyQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createAvgLatencyQuery(namespace, rateDuration, "inbound").
		SetName("avg_server_latency_ms")
}

// CreatePercentileServerLatencyQuery creates query for percentile server response time
func CreatePercentileServerLatencyQuery(namespace string, rateDuration time.Duration, percentile float32) *promql.Query {
	percentileInt := int(percentile * 100)
	return createPercentileLatencyQuery(namespace, rateDuration, "inbound", percentile).
		SetName(fmt.Sprintf("p%v_server_latency_ms", percentileInt))
}

// CreateAvgClientLatencyQuery creates query for average client response time
func CreateAvgClientLatencyQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createAvgLatencyQuery(namespace, rateDuration, "outbound").
		SetName("avg_client_latency_ms")
}

// CreatePercentileServerLatencyQuery creates query for percentile server response time
func CreatePercentileClientLatencyQuery(namespace string, rateDuration time.Duration, percentile float32) *promql.Query {
	percentileInt := int(percentile * 100)
	return createPercentileLatencyQuery(namespace, rateDuration, "outbound", percentile).
		SetName(fmt.Sprintf("p%v_client_latency_ms", percentileInt))
}

// createAvgLatencyQuery create query for average response latency of a deployment
func createAvgLatencyQuery(namespace string, rateDuration time.Duration, direction string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("namespace", "=", namespace),
		promql.NewFilter("direction", "=", direction),
		promql.NewFilter("target_port", "!=", "4191"),
	}
	sum := promql.NewQuery("response_latency_ms_sum").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"deployment"})

	count := promql.NewQuery("response_latency_ms_count").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"deployment"})

	return sum.Divide(count)
}

// createPercentileLatencyQuery create query for given percentile latency of a deployment
func createPercentileLatencyQuery(namespace string, rateDuration time.Duration, direction string, percentile float32) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("namespace", "=", namespace),
		promql.NewFilter("direction", "=", direction),
		promql.NewFilter("target_port", "!=", "4191"),
	}
	return promql.NewQuery("response_latency_ms_bucket").
		Filter(filters).
		Rate(rateDuration).
		HistogramQuantile(percentile).
		SumBy([]string{"deployment"})
}
