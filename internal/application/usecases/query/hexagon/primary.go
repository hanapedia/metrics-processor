package hexagon

import (
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

const PRIMARY_SUM_KEY = "primary_id"
const SERVICE_SUM_KEY = "service"

func NewPrimaryCountQuery(filters []promql.Filter, rateConfig query.RateConfig, sumBy string) *promql.Query {
	count := promql.NewQuery(PrimaryDurationCount.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return count.IRate(rateConfig.Duration).SumBy([]string{sumBy})
	}
	return count.Rate(rateConfig.Duration).SumBy([]string{sumBy})
}

// NewPrimaryRatioQuery create query to take ratios of two primary adapter count queries
func NewPrimaryRatioQuery(numeFilter, denoFilter []promql.Filter, rateConfig query.RateConfig, sumBy string) *promql.Query {
	numeQuery := promql.NewQuery(PrimaryDurationCount.AsString()).Filter(numeFilter)
	denoQuery := promql.NewQuery(PrimaryDurationCount.AsString()).Filter(denoFilter)
	if rateConfig.IsInstant {
		numeQuery.IRate(rateConfig.Duration).SumBy([]string{sumBy})
		denoQuery.IRate(rateConfig.Duration).SumBy([]string{sumBy})
	} else {
		numeQuery.Rate(rateConfig.Duration).SumBy([]string{sumBy})
		denoQuery.Rate(rateConfig.Duration).SumBy([]string{sumBy})
	}
	return numeQuery.Divide(denoQuery)
}

// NewAvgPrimaryDurationQuery create average primary adapter Duration
func NewAvgPrimaryDurationQuery(filters []promql.Filter, rateConfig query.RateConfig, sumBy string) *promql.Query {
	sum := promql.NewQuery(PrimaryDurationSum.AsString()).Filter(filters)
	count := promql.NewQuery(PrimaryDurationCount.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		sum.IRate(rateConfig.Duration).SumBy([]string{sumBy})
		count.IRate(rateConfig.Duration).SumBy([]string{sumBy})
	} else {
		sum.Rate(rateConfig.Duration).SumBy([]string{sumBy})
		count.Rate(rateConfig.Duration).SumBy([]string{sumBy})
	}

	return sum.Divide(count)
}

// NewPercentilePrimaryDurationQuery create query for given percentile for primary duration
func NewPercentilePrimaryDurationQuery(filters []promql.Filter, rateConfig query.RateConfig, sumBy string, percentile float32) *promql.Query {
	query := promql.NewQuery(PrimaryDurationBucket.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).
			SumBy([]string{sumBy, "le"}).
			HistogramQuantile(percentile)
	}
	return query.Rate(rateConfig.Duration).
		SumBy([]string{sumBy, "le"}).
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
