package prometheus

import (
	"context"
	/* "encoding/json" */
	"log/slog"
	/* "os" */
	"testing"
	"time"

	"github.com/hanapedia/metrics-processor/internal/domain"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func TestNewPrometheusAdapter(t *testing.T) {
	// Test for successful client creation
	cfg := &domain.Config{
		MetricsQueryEndpoint: "http://localhost:9091",
		EndTime:              time.Now(),
		Duration:             30 * time.Minute,
		Step:                 15 * time.Second,
	}

	adapter, err := NewPrometheusAdapter(cfg)
	if err != nil {
		slog.Error("failed to initialize prometheus adapter")
		t.Fail()
	}

	query := "sum by (deployment)(label_replace(rate(container_cpu_usage_seconds_total{namespace=\"emulation\",container=~\"server|redis\"}[1m0s]),\"deployment\",\"$1\",\"pod\",\"(.*)-[^-]+-[^-]+\")) / sum by (deployment)(label_replace(kube_pod_container_resource_limits{namespace=\"emulation\",container=~\"server|redis\",resource=\"cpu\"},\"deployment\",\"$1\",\"pod\",\"(.*)-[^-]+-[^-]+\"))"

	result, warnings, err := adapter.client.QueryRange(
		context.Background(),
		query,
		adapter.queryRange,
		v1.WithTimeout(5*time.Second),
	)
	if err != nil {
		slog.Error("Query failed")
		t.Fail()
	}

	for _, warning := range warnings {
		slog.Warn(warning)
	}

	if matrix, ok := result.(model.Matrix); ok {
		_ = adapter.handleMatrixResult("test-query", &matrix)
		/* jsonData, err := json.Marshal(metricsMatrix) */
		/* if err != nil { */
		/* 	slog.Error("Failed to encode to json", "err", err, "name", metricsMatrix.Name) */
		/* 	t.Fail() */
		/* } */
		/* // Write data to file */
		/* err = os.WriteFile("query_result.json", jsonData, 0644) */
		/* if err != nil { */
		/* 	slog.Error("Failed to write json") */
		/* 	t.Fail() */
		/* } */
	} else {
		slog.Warn("Query did not return matrix. Skipping.")
	}
}
