package core

import (
	"log/slog"

	"github.com/hanapedia/metrics-processor/internal/application/port"
	"github.com/hanapedia/metrics-processor/internal/domain"
)

type MetricsProcessor struct {
	query port.MetricsQueryPort
	storage port.MetricsStoragePort
}

func NewMetricsProcessor(query port.MetricsQueryPort, storage port.MetricsStoragePort) *MetricsProcessor {
	return &MetricsProcessor{
		query: query,
		storage: storage,
	}
}

func (ms *MetricsProcessor) Process() {
	metricsChan := make(chan *domain.MetricsMatrix, ms.query.Len())
	ms.query.Query(metricsChan)
	slog.Info("Metrics queried")

	ms.storage.Save(metricsChan)
	slog.Info("Metrics saved")
}
