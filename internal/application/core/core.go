package core

import (
	"log/slog"

	"github.com/hanapedia/metrics-processor/internal/application/port"
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

func (ms *MetricsProcessor) Process() error {
	metrics, err := ms.query.Query()
	if err != nil {
		return err
	}
	slog.Info("Metrics queried")

	if err = ms.storage.Save(metrics); err != nil {
		return err
	}
	slog.Info("Metrics saved")

	return nil
}
