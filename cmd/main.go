package main

import (
	"github.com/hanapedia/metrics-processor/internal/application/core"
	"github.com/hanapedia/metrics-processor/internal/application/usecases"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/config"
)

func main() {
	config := config.NewConfigFromEnv()
	prometheusAdapter := usecases.PrometheusQueryAdapter(config)
	s3Adapter := usecases.NewS3Adapter(config)

	processor := core.NewMetricsProcessor(prometheusAdapter, s3Adapter)
	processor.Process()
}
