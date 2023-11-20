package usecases

import (
	"log/slog"
	"os"

	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/s3"
)

func NewS3Adapter(config *domain.Config) *s3.S3Adapter {
	adapter, err := s3.NewS3Adapter(config)
	if err != nil {
		slog.Error("Failed to create new S3 adapter", "err", err)
		os.Exit(1)
	}
	return adapter
}
