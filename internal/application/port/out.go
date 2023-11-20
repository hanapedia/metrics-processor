package port

import (
	"github.com/hanapedia/metrics-processor/internal/domain"
)

// MetricsQueryPort represents port for querying metrics from arbitrary backend
type MetricsQueryPort interface {
	// Query run the registered queries
	Query() (*domain.Metrics, error)
}

// MetricsStoragePort represents port for storing metrics to arbitrary backend
type MetricsStoragePort interface {
	Save(*domain.Metrics) error
}
