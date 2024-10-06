package hexagon

import (
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

const PRIMARY_SUM_KEY = "primary_id"

// NewAvgPrimaryDurationQuery create average primary adapter Duration
func NewAvgPrimaryDurationQuery(filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	sum := promql.NewQuery(PrimaryDurationSum.AsString()).Filter(filters)
	count := promql.NewQuery(PrimaryDurationCount.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		sum.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY})
		count.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY})
	} else {
		sum.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY})
		count.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY})
	}

	return sum.Divide(count)
}

// NewPercentilePrimaryDurationQuery create query for given percentile for primary duration
func NewPercentilePrimaryDurationQuery(filters []promql.Filter, rateConfig query.RateConfig, percentile float32) *promql.Query {
	query := promql.NewQuery(PrimaryDurationBucket.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).
			SumBy([]string{PRIMARY_SUM_KEY, "le"}).
			HistogramQuantile(percentile)
	}
	return query.Rate(rateConfig.Duration).
		SumBy([]string{PRIMARY_SUM_KEY, "le"}).
		HistogramQuantile(percentile)
}

// NewPrimaryDurationHistogramQuery create query for histogram for primary duration
func NewPrimaryDurationHistogramQuery(filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	query := promql.NewQuery(PrimaryDurationBucket.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).
			SumBy([]string{PRIMARY_SUM_KEY, "le"})
	}
	return query.Rate(rateConfig.Duration).
		SumBy([]string{PRIMARY_SUM_KEY, "le"})
}

// NewPrimaryInProgressQuery create query for primary adapter in progress
func NewPrimaryInProgressQuery(filters []promql.Filter) *promql.Query {
	return promql.NewQuery(PrimaryAdapterInProgress.AsString()).
		Filter(filters).
		SumBy([]string{PRIMARY_SUM_KEY})
}
