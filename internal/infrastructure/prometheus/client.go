package prometheus

import (
	"context"
	"log/slog"
	"strings"
	"sync"
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
	}, nil
}

func (pa *PrometheusAdapter) RegisterQuery(query *Query) {
	pa.queries = append(pa.queries, query)
}

func (pa *PrometheusAdapter) Len() int {
	return len(pa.queries)
}

func (pa *PrometheusAdapter) Query(metricsChan chan<- *domain.MetricsMatrix) {
	var wg sync.WaitGroup

	for _, query := range pa.queries {
		wg.Add(1)
		go func(q *Query) {
			defer wg.Done()
			pa.runQuery(q, metricsChan)
		}(query)
	}

	go func() {
		wg.Wait()
		close(metricsChan)
	}()
}

func (pa *PrometheusAdapter) runQuery(query *Query, metricsChan chan<- *domain.MetricsMatrix) {
	result, warnings, err := pa.client.QueryRange(
		context.Background(),
		query.AsString(),
		pa.queryRange,
		v1.WithTimeout(5*time.Second),
	)
	if err != nil {
		slog.Error("Query failed", "name", query.name, "error", err, "query", query.AsString())
		return
	}

	for _, warning := range warnings {
		slog.Warn(warning)
	}

	if matrix, ok := result.(model.Matrix); ok {
		metricsChan <- pa.handleMatrixResult(query.name, &matrix)
	} else {
		slog.Warn("Query did not return matrix. Skipping.", "name", query.name, "query", query.q)
	}
}

func (pa *PrometheusAdapter) handleMatrixResult(name string, matrix *model.Matrix) *domain.MetricsMatrix {
	metricsMatrix := domain.MetricsMatrix{
		Name:       name,
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
	return &metricsMatrix
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
