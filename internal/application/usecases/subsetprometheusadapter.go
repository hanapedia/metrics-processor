package usecases

import (
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query/hexagon"
	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/prometheus"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// SubsetPrometheusQueryAdapter creates proemtheusAdapter with subset of queries
// use this adapter when partial requery is needed
func SubsetPrometheusQueryAdapter(config *domain.Config) *prometheus.PrometheusAdapter {
	prometheusAdapter, err := prometheus.NewPrometheusAdapter(config)
	if err != nil {
		slog.Error("Failed to create new Prometheus adapter", "err", err)
		os.Exit(1)
	}

	rateConfigs := []query.RateConfig{
		{Name: "5m", Duration: 5 * time.Minute, IsInstant: false},
		{Name: "1m", Duration: 1 * time.Minute, IsInstant: false},
		{Name: "1m", Duration: 1 * time.Minute, IsInstant: true},
	}
	defaultSvc := "service-.*"
	filters := []promql.Filter{
		promql.NewFilter("experiment", "=~", config.K6TestName),
		promql.NewFilter("service", "=~", defaultSvc),
		promql.NewFilter("namespace", "=", config.Namespace),
	}
	/* statusOkFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "ok")) */
	statusErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "!=", "ok"))
	statusTimeoutErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-ctx-timed-out|error-ctx-canceled"))
	statusCBOpenErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-cb-open"))

	// Register non-rate or non-irate queries
	/* queries := []*promql.Query{ */
	/* 	// adaptive timeout */
	/* 	hexagon.NewAdaptiveTimeoutQuery(hexagon.Call, filters).SetName("adaptive_call_timeout"), // adaptive call timeout */
	/* 	hexagon.NewAdaptiveTimeoutQuery(hexagon.Task, filters).SetName("adaptive_task_timeout"), // adaptive task timeout */
	/* } */
	/* for _, query := range queries { */
	/* 	prometheusAdapter.RegisterQuery(query) */
	/* } */

	// Register rate & irate queries
	for _, rateConfig := range rateConfigs {
		queries := []*promql.Query{
			hexagon.NewSecondaryRatioQuery(hexagon.Call, statusErrFilter, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_err_rate")), // failure rate
			hexagon.NewSecondaryRatioQuery(hexagon.Call, statusTimeoutErrFilter, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_timeout_err_rate")), // timeout rate
			hexagon.NewSecondaryRatioQuery(hexagon.Call, statusCBOpenErrFilter, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_cb_err_rate")), // cb error rate
		}
		for _, query := range queries {
			prometheusAdapter.RegisterQuery(query)
		}
	}

	return prometheusAdapter
}
