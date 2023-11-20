package prometheus

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusAdapter struct {
	client     v1.API
	queryRange v1.Range
	queries    []*Query
	metrics    *domain.Metrics
}

func NewPrometheusAdapter(config *domain.Config) (*PrometheusAdapter, error) {
	client, err := api.NewClient(api.Config{
		Address: config.MetricsQueryEndpoing,
	})
	if err != nil {
		return nil, err
	}

	return &PrometheusAdapter{
		client: v1.NewAPI(client),
		queryRange: v1.Range{
			Start: config.EndTime.Add(-1 * config.Duration),
			End:   config.EndTime,
			Step:  config.Step,
		},
		metrics: &domain.Metrics{
			QueryConfig: config,
			Data:        make(map[string]domain.MetricsMatrix),
		},
	}, nil
}

func (pa *PrometheusAdapter) RegisterQuery(query *Query) {
	pa.queries = append(pa.queries, query)
}

func (pa *PrometheusAdapter) Query() (*domain.Metrics, error) {
	for _, query := range pa.queries {
		result, warnings, err := pa.client.QueryRange(
			context.Background(),
			query.AsString(),
			pa.queryRange,
			v1.WithTimeout(5*time.Second),
		)
		if err != nil {
			return nil, err
		}

		for _, warning := range warnings {
			slog.Warn(warning)
		}

		if matrix, ok := result.(model.Matrix); ok {
			pa.handleMatrixResult(query.name, &matrix)
		} else {
			slog.Warn("Query did not return matrix. Skipping.", "name", query.name, "query", query.q)
			continue
		}
	}
	return pa.metrics, nil
}

func (pa *PrometheusAdapter) handleMatrixResult(name string, matrix *model.Matrix) {
	metricsMatrix := domain.MetricsMatrix{
		LabelType:  "",
		Matrix:     make(map[string][]float64),
		Timestamps: []int64{},
	}
	for i, sampleStream := range *matrix {
		if i == 0 {
			metricsMatrix.LabelType = extractLabelType(sampleStream.Metric)
			metricsMatrix.Timestamps = extractTimestamps(sampleStream.Values)
		}
		label := extractLabelValue(sampleStream.Metric)
		metricsMatrix.Matrix[label] = extractSampleValues(sampleStream.Values)
	}
	pa.metrics.Data[name] = metricsMatrix
}

func extractLabelType(metric model.Metric) string {
	keys := []string{}
	for k := range metric {
		keys = append(keys, string(k))
	}
	return strings.Join(keys, "_")
}

func extractLabelValue(metric model.Metric) string {
	values := []string{}
	for _, v := range metric {
		values = append(values, string(v))
	}
	return strings.Join(values, "_")
}

func extractTimestamps(samples []model.SamplePair) []int64 {
	timestamps := []int64{}
	for _, sample := range samples {
		timestamps = append(timestamps, sample.Timestamp.Unix())
	}
	return timestamps
}

func extractSampleValues(samples []model.SamplePair) []float64 {
	values := []float64{}
	for _, sample := range samples {
		values = append(values, float64(sample.Value))
	}
	return values
}
