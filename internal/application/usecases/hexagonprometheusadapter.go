package usecases

import (
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query/container"
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query/hexagon"
	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/prometheus"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

func HexagonPrometheusQueryAdapter(config *domain.Config) *prometheus.PrometheusAdapter {
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
	statusOkFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "ok"))
	statusErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "!=", "ok"))
	statusTimeoutErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-ctx-timed-out|error-ctx-canceled"))
	statusCBOpenErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-cb-open"))
	containerFilter := []promql.Filter{
		promql.NewFilter("namespace", "=", config.Namespace),
		promql.NewFilter("container", "!=", ""),
	}

	// Register non-rate or non-irate queries
	queries := []*promql.Query{
		// queue length
		hexagon.NewPrimaryInProgressQuery(filters).SetName("primary_in_progress"),

		// container metrics
		container.CreateMemoryUsageQuery(containerFilter).SetName("memory_usage"),
		container.CreateContainerRestartsQuery(containerFilter).SetName("container_restarts"),

		// adaptive timeout
		hexagon.NewAdaptiveTimeoutQuery(hexagon.Call, filters).SetName("adaptive_call_timeout"), // adaptive call timeout
		hexagon.NewAdaptiveTimeoutQuery(hexagon.Task, filters).SetName("adaptive_task_timeout"), // adaptive task timeout
	}
	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	// Register rate & irate queries
	for _, rateConfig := range rateConfigs {
		queries := []*promql.Query{
			// primary adatper metrics
			// p99, p50, avg durations
			// in progress count
			// success rate
			hexagon.NewAvgPrimaryDurationQuery(filters, rateConfig).
				SetName(rateConfig.AddSuffix("avg_primary_duration")),
			hexagon.NewPercentilePrimaryDurationQuery(filters, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_primary_duration")), // p99
			hexagon.NewPrimaryDurationHistogramQuery(filters, rateConfig).
				SetName(rateConfig.AddSuffix("primary_duration_histogram")), // p99

			// secondary adatper call metrics
			hexagon.NewSecondaryCountQuery(hexagon.Call, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_all_count")), // count all
			hexagon.NewSecondaryCountQuery(hexagon.Call, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_ok_count")), // count ok
			hexagon.NewSecondaryCountQuery(hexagon.Call, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_timeout_err_count")), // count timeout
			hexagon.NewSecondaryCountQuery(hexagon.Call, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_cb_err_count")), // count cb error
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, filters, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_call_all_duration")), // call avg duration all
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_call_ok_duration")), // call avg duration ok
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_call_err_duration")), // call avg duration err
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_call_timeout_err_duration")), // call avg duration err
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_call_cb_err_duration")), // call avg duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, filters, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_call_all_duration")), // call p99 duration all
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_call_ok_duration")), // call p99 duration ok
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_call_err_duration")), // call p99 duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusTimeoutErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_call_timeout_err_duration")), // call p99 duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusCBOpenErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_call_cb_err_duration")), // call p99 duration err

			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Call, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_all_duration_histogram")), // call p99 duration all
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Call, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_ok_duration_histogram")), // call p99 duration ok
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Call, statusErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_err_duration_histogram")), // call p99 duration err
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Call, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_timeout_err_duration_histogram")), // call p99 duration err
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Call, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_call_cb_err_duration_histogram")), // call p99 duration err

			hexagon.NewThresholdBucketSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateConfig, 2.5).
				SetName(rateConfig.AddSuffix("secondary_duration_under_p99")), // calls under 2.5ms
			hexagon.NewRetryRateQuery(filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_retry_rate")), // retry rate

			// secondary adatper task metrics
			hexagon.NewSecondaryCountQuery(hexagon.Task, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_all_count")), // count all
			hexagon.NewSecondaryCountQuery(hexagon.Task, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_ok_count")), // count ok
			hexagon.NewSecondaryCountQuery(hexagon.Task, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_timeout_err_count")), // count timeout
			hexagon.NewSecondaryCountQuery(hexagon.Task, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_cb_err_count")), // count cb error
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, filters, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_task_all_duration")), // task avg duration all
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_task_ok_duration")), // task avg duration ok
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_task_err_duration")), // task avg duration err
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_task_timeout_err_duration")), // task avg duration err
			hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("avg_secondary_task_cb_err_duration")), // task avg duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, filters, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_task_all_duration")), // task p99 duration all
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusOkFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_task_ok_duration")), // task p99 duration ok
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_task_err_duration")), // task p99 duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusTimeoutErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_task_timeout_err_duration")), // task p99 duration err
			hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusCBOpenErrFilter, rateConfig, 0.99).
				SetName(rateConfig.AddSuffix("p99_secondary_task_cb_err_duration")), // task p99 duration err

			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Task, filters, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_all_duration_histogram")), // task p99 duration all
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Task, statusOkFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_ok_duration_histogram")), // task p99 duration ok
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Task, statusErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_err_duration_histogram")), // task p99 duration err
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Task, statusTimeoutErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_timeout_err_duration_histogram")), // task p99 duration err
			hexagon.NewSecondaryDurationHistogramQuery(hexagon.Task, statusCBOpenErrFilter, rateConfig).
				SetName(rateConfig.AddSuffix("secondary_task_cb_err_duration_histogram")), // task p99 duration err

			// container metrics
			container.CreateCpuUsageQuery(containerFilter, rateConfig).
				SetName(rateConfig.AddSuffix("cpu_usage")),
			container.CreateCpuThrottleQuery(containerFilter, rateConfig).
				SetName(rateConfig.AddSuffix("cpu_throttled")),

			// k6 metrics
			query.CreateK6IterationRateQuery(config.K6TestName, rateConfig).
				SetName(rateConfig.AddSuffix("k6_iterations")),
			query.CreateK6DroppedIterationRateQuery(config.K6TestName, rateConfig).
				SetName(rateConfig.AddSuffix("k6_dropped_iterations")),
		}
		for _, query := range queries {
			prometheusAdapter.RegisterQuery(query)
		}
	}

	return prometheusAdapter
}
