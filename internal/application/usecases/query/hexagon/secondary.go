package hexagon

import (
	"fmt"
	"time"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

const SECONDARY_SUM_KEY = "secondary_id"

type SecondaryDurationVariant = float64

const (
	Task SecondaryDurationVariant = iota
	Call
)

// NewSecondaryCountQuery create secondary adapter invocation count query
func NewSecondaryCountQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	var countQuery query.MetricsName
	switch variant {
	case Task:
		countQuery = TaskDurationCount
	case Call:
		countQuery = CallDurationCount
	}

	return promql.NewQuery(countQuery.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
}

// NewAvgSecondaryDurationQuery create average secondary adapter Duration
func NewAvgSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	var sumQuery query.MetricsName
	var countQuery query.MetricsName
	switch variant {
	case Task:
		sumQuery = TaskDurationSum
		countQuery = TaskDurationCount
	case Call:
		sumQuery = CallDurationSum
		countQuery = CallDurationCount
	}

	sum := promql.NewQuery(sumQuery.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	count := promql.NewQuery(countQuery.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	return sum.Divide(count)
}

// NewPercentileSecondaryDurationQuery create query for given percentile for secondary duration
func NewPercentileSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateDuration time.Duration, percentile float32) *promql.Query {
	var bucketQuery query.MetricsName
	switch variant {
	case Task:
		bucketQuery = TaskDurationBucket
	case Call:
		bucketQuery = CallDurationBucket
	}
	return promql.NewQuery(bucketQuery.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY, "le"}).
		HistogramQuantile(percentile)
}

// NewThresholdSecondaryDurationQuery create query for ratio under percentile for secondary duration
func NewThresholdBucketSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateDuration time.Duration, le float32) *promql.Query {
	var bucketQuery query.MetricsName
	var countQuery query.MetricsName
	switch variant {
	case Task:
		bucketQuery = TaskDurationBucket
		countQuery = TaskDurationCount
	case Call:
		bucketQuery = CallDurationBucket
		countQuery = CallDurationCount
	}

	count := promql.NewQuery(countQuery.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	hist := promql.NewQuery(bucketQuery.AsString()).
		Filter(append(filters, promql.NewFilter("le", "=~", fmt.Sprintf("%v", le)))).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	return hist.Divide(count)
}

// NewRetryRateQuery create query for retry rate
// retry rate is the number of retry calls divide by number of all calls
func NewRetryRateQuery(filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	all := promql.NewQuery(CallDurationCount.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	retry := promql.NewQuery(CallDurationCount.AsString()).
		Filter(append(filters, promql.NewFilter("nth_attempt", "!~", "(1|0)"))).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})

	return retry.Divide(all)
}
