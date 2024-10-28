package usecases

import (
	"log/slog"
	"os"
	"time"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query/container"
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
	containerFilter := []promql.Filter{
		promql.NewFilter("namespace", "=", config.Namespace),
		promql.NewFilter("container", "!=", ""),
	}

	// Register non-rate or non-irate queries
	queries := []*promql.Query{
		// container metrics
		container.CreateMemoryUsageQuery(containerFilter).SetName("memory_usage"),
		container.CreateContainerRestartsQuery(containerFilter).SetName("container_restarts"),
	}
	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	// Register rate & irate queries
	for _, rateConfig := range rateConfigs {
		queries := []*promql.Query{

			// container metrics
			container.CreateCpuUsageQuery(containerFilter, rateConfig).
				SetName(rateConfig.AddSuffix("cpu_usage")),
			container.CreateCpuThrottleQuery(containerFilter, rateConfig).
				SetName(rateConfig.AddSuffix("cpu_throttled")),
		}
		for _, query := range queries {
			prometheusAdapter.RegisterQuery(query)
		}
	}

	return prometheusAdapter
}
