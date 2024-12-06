package hexagon

import (
	"fmt"

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
func NewSecondaryCountQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	var countQuery query.MetricsName
	switch variant {
	case Task:
		countQuery = TaskDurationCount
	case Call:
		countQuery = CallDurationCount
	}

	query := promql.NewQuery(countQuery.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	}
	return query.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
}

// NewSecondaryRatioQuery create query to take ratios of two secondary adapter count queries
func NewSecondaryRatioQuery(variant SecondaryDurationVariant, numeFilter, denoFilter []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	var countQuery query.MetricsName
	switch variant {
	case Task:
		countQuery = TaskDurationCount
	case Call:
		countQuery = CallDurationCount
	}

	numeQuery := promql.NewQuery(countQuery.AsString()).Filter(numeFilter)
	denoQuery := promql.NewQuery(countQuery.AsString()).Filter(denoFilter)
	if rateConfig.IsInstant {
		numeQuery.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		denoQuery.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	} else {
		numeQuery.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		denoQuery.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	}
	return numeQuery.Divide(denoQuery)
}

// NewAvgSecondaryDurationQuery create average secondary adapter Duration
func NewAvgSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
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

	sum := promql.NewQuery(sumQuery.AsString()).Filter(filters)
	count := promql.NewQuery(countQuery.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		sum.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		count.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	} else {
		sum.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		count.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	}

	return sum.Divide(count)
}

// NewPercentileSecondaryDurationQuery create query for given percentile for secondary duration
func NewPercentileSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateConfig query.RateConfig, percentile float32) *promql.Query {
	var bucketQuery query.MetricsName
	switch variant {
	case Task:
		bucketQuery = TaskDurationBucket
	case Call:
		bucketQuery = CallDurationBucket
	}
	query := promql.NewQuery(bucketQuery.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).
			SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY, "le"}).
			HistogramQuantile(percentile)
	}
	return query.Rate(rateConfig.Duration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY, "le"}).
		HistogramQuantile(percentile)
}

// NewSecondaryDurationHistogramQuery create query for histogram of secondary duration
func NewSecondaryDurationHistogramQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	var bucketQuery query.MetricsName
	switch variant {
	case Task:
		bucketQuery = TaskDurationBucket
	case Call:
		bucketQuery = CallDurationBucket
	}
	query := promql.NewQuery(bucketQuery.AsString()).Filter(filters)
	if rateConfig.IsInstant {
		return query.IRate(rateConfig.Duration).
			SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY, "le"})
	}
	return query.Rate(rateConfig.Duration).
		SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY, "le"})
}

// NewThresholdSecondaryDurationQuery create query for ratio under percentile for secondary duration
func NewThresholdBucketSecondaryDurationQuery(variant SecondaryDurationVariant, filters []promql.Filter, rateConfig query.RateConfig, le float32) *promql.Query {
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

	count := promql.NewQuery(countQuery.AsString()).Filter(filters)
	hist := promql.NewQuery(bucketQuery.AsString()).Filter(append(filters, promql.NewFilter("le", "=~", fmt.Sprintf("%v", le))))
	if rateConfig.IsInstant {
		count.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		hist.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	} else {
		count.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		hist.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	}

	return hist.Divide(count)
}

// NewRetryRateQuery create query for retry rate
// retry rate is the number of retry calls divide by number of all calls
func NewRetryRateQuery(filters []promql.Filter, rateConfig query.RateConfig) *promql.Query {
	all := promql.NewQuery(CallDurationCount.AsString()).Filter(filters)
	retry := promql.NewQuery(CallDurationCount.AsString()).Filter(append(filters, promql.NewFilter("nth_attempt", "!~", "(1|0)")))
	if rateConfig.IsInstant {
		all.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		retry.IRate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	} else {
		all.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
		retry.Rate(rateConfig.Duration).SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
	}

	return retry.Divide(all)
}

// NewAdaptiveTimeoutQuery create secondary adapter invocation count query
func NewAdaptiveTimeoutQuery(variant SecondaryDurationVariant, filters []promql.Filter) *promql.Query {
	var countQuery query.MetricsName
	switch variant {
	case Task:
		countQuery = TaskAdaptiveTimeout
	case Call:
		countQuery = CallAdaptiveTimeout
	}

	query := promql.NewQuery(countQuery.AsString()).Filter(filters)
	return query.SumBy([]string{PRIMARY_SUM_KEY, SECONDARY_SUM_KEY})
}
