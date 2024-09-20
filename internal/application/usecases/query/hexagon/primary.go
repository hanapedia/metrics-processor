package hexagon

import (
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

const PRIMARY_SUM_KEY = "primary_id"

// NewAvgPrimaryDurationQuery create average primary adapter Duration
func NewAvgPrimaryDurationQuery(filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	sum := promql.NewQuery(PrimaryDurationSum.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY})

	count := promql.NewQuery(PrimaryDurationCount.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY})

	return sum.Divide(count)
}

// NewPercentilePrimaryDurationQuery create query for given percentile for primary duration
func NewPercentilePrimaryDurationQuery(filters []promql.Filter, rateDuration time.Duration, percentile float32) *promql.Query {
	return promql.NewQuery(PrimaryDurationBucket.AsString()).
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{PRIMARY_SUM_KEY, "le"}).
		HistogramQuantile(percentile)
}

// NewPrimaryInProgressQuery create query for primary adapter in progress
func NewPrimaryInProgressQuery(filters []promql.Filter, rateDuration time.Duration) *promql.Query {
	return promql.NewQuery(PrimaryAdapterInProgress.AsString()).
		Filter(filters).
		SumBy([]string{PRIMARY_SUM_KEY})
}
