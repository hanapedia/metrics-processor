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

	rateDuration := 5 * time.Minute
	defaultSvc := "service-.*"
	filters := []promql.Filter{
		promql.NewFilter("experiment", "=~", config.K6TestName),
		promql.NewFilter("service", "=~", defaultSvc),
	}
	statusOkFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "ok"))
	statusErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "!=", "ok"))
	statusTimeoutErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-ctx-timed-out|error-ctx-canceled"))
	statusCBOpenErrFilter := append(slices.Clone(filters), promql.NewFilter("status", "=~", "error-cb-open"))

	queries := []*promql.Query{
		// primary adatper metrics
		// p95, p99, p50, avg durations
		// in progress count
		// success rate
		hexagon.NewAvgPrimaryDurationQuery(filters, rateDuration).SetName("avg_primary_duration"),
		hexagon.NewPercentilePrimaryDurationQuery(filters, rateDuration, 0.99).SetName("p99_primary_duration"), // p99
		hexagon.NewPercentilePrimaryDurationQuery(filters, rateDuration, 0.95).SetName("p95_primary_duration"), // p95
		hexagon.NewPrimaryInProgressQuery(filters, rateDuration).SetName("primary_in_progress"),

		// secondary adatper call metrics
		hexagon.NewSecondaryCountQuery(hexagon.Call, filters, rateDuration).SetName("secondary_call_all_count"),                                                  // count all
		hexagon.NewSecondaryCountQuery(hexagon.Call, statusOkFilter, rateDuration).SetName("secondary_call_ok_count"),                                            // count ok
		hexagon.NewSecondaryCountQuery(hexagon.Call, statusTimeoutErrFilter, rateDuration).SetName("secondary_call_timeout_err_count"),                           // count timeout
		hexagon.NewSecondaryCountQuery(hexagon.Call, statusCBOpenErrFilter, rateDuration).SetName("secondary_call_cb_err_count"),                                 // count cb error
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, filters, rateDuration).SetName("avg_secondary_call_all_duration"),                                     // call avg duration all
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateDuration).SetName("avg_secondary_call_ok_duration"),                               // call avg duration ok
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusErrFilter, rateDuration).SetName("avg_secondary_call_err_duration"),                             // call avg duration err
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusTimeoutErrFilter, rateDuration).SetName("avg_secondary_call_timeout_err_duration"),              // call avg duration err
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Call, statusCBOpenErrFilter, rateDuration).SetName("avg_secondary_call_cb_err_duration"),                    // call avg duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, filters, rateDuration, 0.99).SetName("p99_secondary_call_all_duration"),                        // call p99 duration all
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateDuration, 0.99).SetName("p99_secondary_call_ok_duration"),                  // call p99 duration ok
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusErrFilter, rateDuration, 0.99).SetName("p99_secondary_call_err_duration"),                // call p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusTimeoutErrFilter, rateDuration, 0.99).SetName("p99_secondary_call_timeout_err_duration"), // call p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusCBOpenErrFilter, rateDuration, 0.99).SetName("p99_secondary_call_cb_err_duration"),       // call p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, filters, rateDuration, 0.95).SetName("p95_secondary_call_all_duration"),                        // call p95 duration all
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateDuration, 0.95).SetName("p95_secondary_call_ok_duration"),                  // call p95 duration ok
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusErrFilter, rateDuration, 0.95).SetName("p95_secondary_call_err_duration"),                // call p95 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusTimeoutErrFilter, rateDuration, 0.95).SetName("p95_secondary_call_timeout_err_duration"), // call p95 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Call, statusCBOpenErrFilter, rateDuration, 0.95).SetName("p95_secondary_call_cb_err_duration"),       // call p95 duration err
		hexagon.NewThresholdBucketSecondaryDurationQuery(hexagon.Call, statusOkFilter, rateDuration, 2.5).SetName("secondary_duration_under_p99"),                // calls under 2.5ms
		hexagon.NewRetryRateQuery(filters, rateDuration).SetName("secondary_retry_rate"),                                                                         // retry rate

		// secondary adatper task metrics
		hexagon.NewSecondaryCountQuery(hexagon.Task, filters, rateDuration).SetName("secondary_task_all_count"),                                                  // count all
		hexagon.NewSecondaryCountQuery(hexagon.Task, statusOkFilter, rateDuration).SetName("secondary_task_ok_count"),                                            // count ok
		hexagon.NewSecondaryCountQuery(hexagon.Task, statusTimeoutErrFilter, rateDuration).SetName("secondary_task_timeout_err_count"),                           // count timeout
		hexagon.NewSecondaryCountQuery(hexagon.Task, statusCBOpenErrFilter, rateDuration).SetName("secondary_task_cb_err_count"),                                 // count cb error
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, filters, rateDuration).SetName("avg_secondary_task_all_duration"),                                     // task avg duration all
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusOkFilter, rateDuration).SetName("avg_secondary_task_ok_duration"),                               // task avg duration ok
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusErrFilter, rateDuration).SetName("avg_secondary_task_err_duration"),                             // task avg duration err
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusTimeoutErrFilter, rateDuration).SetName("avg_secondary_task_timeout_err_duration"),              // task avg duration err
		hexagon.NewAvgSecondaryDurationQuery(hexagon.Task, statusCBOpenErrFilter, rateDuration).SetName("avg_secondary_task_cb_err_duration"),                    // task avg duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, filters, rateDuration, 0.99).SetName("p99_secondary_task_all_duration"),                        // task p99 duration all
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusOkFilter, rateDuration, 0.99).SetName("p99_secondary_task_ok_duration"),                  // task p99 duration ok
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusErrFilter, rateDuration, 0.99).SetName("p99_secondary_task_err_duration"),                // task p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusTimeoutErrFilter, rateDuration, 0.99).SetName("p99_secondary_task_timeout_err_duration"), // task p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusCBOpenErrFilter, rateDuration, 0.99).SetName("p99_secondary_task_cb_err_duration"),       // task p99 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, filters, rateDuration, 0.95).SetName("p95_secondary_task_all_duration"),                        // task p95 duration all
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusOkFilter, rateDuration, 0.95).SetName("p95_secondary_task_ok_duration"),                  // task p95 duration ok
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusErrFilter, rateDuration, 0.95).SetName("p95_secondary_task_err_duration"),                // task p95 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusTimeoutErrFilter, rateDuration, 0.95).SetName("p95_secondary_task_timeout_err_duration"), // task p95 duration err
		hexagon.NewPercentileSecondaryDurationQuery(hexagon.Task, statusCBOpenErrFilter, rateDuration, 0.95).SetName("p95_secondary_task_cb_err_duration"),       // task p95 duration err

		// container metrics
		container.CreateCpuUsageQuery(
			[]promql.Filter{
				promql.NewFilter("namespace", "=", config.Namespace),
				promql.NewFilter("container", "!=", ""),
			},
			rateDuration,
		).SetName("cpu_usage"),
		container.CreateMemoryUsageQuery(
			[]promql.Filter{
				promql.NewFilter("namespace", "=", config.Namespace),
				promql.NewFilter("container", "!=", ""),
			},
		).SetName("memory_usage"),
		container.CreateContainerRestartsQuery(
			[]promql.Filter{
				promql.NewFilter("namespace", "=", config.Namespace),
				promql.NewFilter("container", "!=", ""),
			},
		).SetName("container_restarts"),

		// k6 metrics
		query.CreateK6IterationRateQuery(config.K6TestName, rateDuration).SetName("k6_iterations"),
		query.CreateK6DroppedIterationRateQuery(config.K6TestName, rateDuration).SetName("k6_dropped_iterations"),
	}

	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	return prometheusAdapter
}
