package usecases

import (
	"log/slog"
	"os"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/prometheus"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

func PrometheusQueryAdapter(config *domain.Config) *prometheus.PrometheusAdapter {
	prometheusAdapter, err := prometheus.NewPrometheusAdapter(config)
	if err != nil {
		slog.Error("Failed to create new Prometheus adapter", "err", err)
		os.Exit(1)
	}

	rateDuration := config.Step * 16

	queries := []*promql.Query{
		// server metrics
		query.CreateAvgServerLatencyQuery(config.Namespace, rateDuration),
		query.CreatePercentileServerLatencyQuery(config.Namespace, rateDuration, 0.95),
		query.CreatePercentileServerLatencyQuery(config.Namespace, rateDuration, 0.99),
		query.CreateServerReadBytesQuery(config.Namespace, rateDuration),
		query.CreateServerWriteBytesQuery(config.Namespace, rateDuration),

		// server latency from client
		query.CreateAvgServerLatencyFromClientQuery(config.Namespace, rateDuration),
		query.CreatePercentileServerLatencyFromClientQuery(config.Namespace, rateDuration, 0.95),
		query.CreatePercentileServerLatencyFromClientQuery(config.Namespace, rateDuration, 0.99),

		// client metrics
		query.CreateAvgClientLatencyQuery(config.Namespace, rateDuration),
		query.CreatePercentileClientLatencyQuery(config.Namespace, rateDuration, 0.95),
		query.CreatePercentileClientLatencyQuery(config.Namespace, rateDuration, 0.99),
		query.CreateClientReadBytesQuery(config.Namespace, rateDuration),
		query.CreateClientWriteBytesQuery(config.Namespace, rateDuration),

		// resource metrics
		query.CreateCpuUsageQuery([]promql.Filter{
			promql.NewFilter("namespace", "=", config.Namespace),
			promql.NewFilter("container", "=", config.WorkloadContainers),
		},
			rateDuration),
		query.CreateMemoryUsageQuery([]promql.Filter{
			promql.NewFilter("namespace", "=", config.Namespace),
			promql.NewFilter("container", "=", config.WorkloadContainers),
		}),

		// k6 metrics
		query.CreateK6IterationRateQuery(config.K6TestName, rateDuration),
		query.CreateK6BytesReceivedQuery(config.K6TestName, rateDuration),
		query.CreateK6BytesSentQuery(config.K6TestName, rateDuration),
		query.CreateAvgK6IterationDurationQuery(config.K6TestName),
		query.CreateP95K6IterationDurationQuery(config.K6TestName),
		query.CreateP99K6IterationDurationQuery(config.K6TestName),
	}

	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	return prometheusAdapter
}
