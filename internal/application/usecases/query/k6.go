package query

import (
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateK6IterationRateQuery create query for iteration per second (rps)
func CreateK6IterationRateQuery(testName string, rateDuration time.Duration) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_iteration_duration_avg").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"scenario"}).
		SetName("lg_iteration_rate")
}

// CreateAvgK6IterationDurationQuery create query for average duration for each request
func CreateAvgK6IterationDurationQuery(testName string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_http_req_duration_avg").
		Filter(filters).
		SumBy([]string{"scenario"}).
		MultiplyByConstant(1000).
		SetName("avg_lg_request_duration_ms")
}

// CreateP95K6IterationDurationQuery create query for p95 duration for each request
func CreateP95K6IterationDurationQuery(testName string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_http_req_duration_p95").
		Filter(filters).
		SumBy([]string{"scenario"}).
		MultiplyByConstant(1000).
		SetName("p95_lg_request_duration_ms")
}

// CreateP99K6IterationDurationQuery create query for p99 duration for each request
func CreateP99K6IterationDurationQuery(testName string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_http_req_duration_p99").
		Filter(filters).
		SumBy([]string{"scenario"}).
		MultiplyByConstant(1000).
		SetName("p99_lg_request_duration_ms")
}

// CreateK6BytesSentQuery create query for bytes sent by loadgenerator
func CreateK6BytesSentQuery(testName string, rateDuration time.Duration) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_data_sent_total").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"scenario"}).
		SetName("lg_bytes_sent")
}

// CreateK6BytesReceivedQuery create query for bytes received by loadgenerator
func CreateK6BytesReceivedQuery(testName string, rateDuration time.Duration) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("name", "=", testName),
	}
	return promql.NewQuery("k6_data_received_total").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"scenario"}).
		SetName("lg_bytes_received")
}
