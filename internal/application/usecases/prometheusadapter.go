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

	rateDuration := config.Step * 4

	queries := []*promql.Query{
		// server metrics
		query.CreateAvgServerLatencyQuery(config.Namespace, rateDuration),
		// query.CreatePercentileServerLatencyQuery(config.Namespace, rateDuration, 0.95),
		// query.CreatePercentileServerLatencyQuery(config.Namespace, rateDuration, 0.99),
		// query.CreateServerReadBytesQuery(config.Namespace, rateDuration),
		// query.CreateServerWriteBytesQuery(config.Namespace, rateDuration),

		// client metrics
		// query.CreateAvgClientLatencyQuery(config.Namespace, rateDuration),
		// query.CreatePercentileClientLatencyQuery(config.Namespace, rateDuration, 0.95),
		// query.CreatePercentileClientLatencyQuery(config.Namespace, rateDuration, 0.99),
		// query.CreateClientReadBytesQuery(config.Namespace, rateDuration),
		// query.CreateClientWriteBytesQuery(config.Namespace, rateDuration),

		// resource metrics
		// query.CreateCpuUsageQuery(config.Namespace, config.WorkloadContainers, rateDuration),
		// query.CreateMemoryUsageQuery(config.Namespace, config.WorkloadContainers),

		// k6 metrics
		query.CreateK6IterationRateQuery(config.TestName, rateDuration),
		// query.CreateK6BytesReceivedQuery(config.TestName, rateDuration),
		// query.CreateK6BytesSentQuery(config.TestName, rateDuration),
		// query.CreateAvgK6IterationDurationQuery(config.TestName),
		// query.CreateP95K6IterationDurationQuery(config.TestName),
		// query.CreateP99K6IterationDurationQuery(config.TestName),
	}

	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	return prometheusAdapter
}
