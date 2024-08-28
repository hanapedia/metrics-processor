package prometheus

import (
	"context"
	"log/slog"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/pkg/promql"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusAdapter struct {
	client     v1.API
	queryRange v1.Range
	queries    []*promql.Query
}

func NewPrometheusAdapter(config *domain.Config) (*PrometheusAdapter, error) {
	client, err := api.NewClient(api.Config{
		Address: config.MetricsQueryEndpoint,
	})
	if err != nil {
		return nil, err
	}

	start := config.EndTime.Add(-1 * config.Duration)

	slog.Info("Query Range set.", "start", start, "end", config.EndTime)

	return &PrometheusAdapter{
		client: v1.NewAPI(client),
		queryRange: v1.Range{
			Start: start,
			End:   config.EndTime,
			Step:  config.Step,
		},
	}, nil
}

func (pa *PrometheusAdapter) RegisterQuery(query *promql.Query) {
	pa.queries = append(pa.queries, query)
}

func (pa *PrometheusAdapter) Len() int {
	return len(pa.queries)
}

func (pa *PrometheusAdapter) Query(metricsChan chan<- *domain.MetricsMatrix) {
	var wg sync.WaitGroup

	for _, query := range pa.queries {
		wg.Add(1)
		go func(q *promql.Query) {
			defer wg.Done()
			pa.runQuery(q, metricsChan)
		}(query)
	}

	go func() {
		wg.Wait()
		close(metricsChan)
	}()
}

func (pa *PrometheusAdapter) runQuery(query *promql.Query, metricsChan chan<- *domain.MetricsMatrix) {
	slog.Info("Running Query.", "name", query.Name, "query", query.AsString())
	result, warnings, err := pa.client.QueryRange(
		context.Background(),
		query.AsString(),
		pa.queryRange,
		v1.WithTimeout(5*time.Second),
	)
	if err != nil {
		slog.Error("Query failed", "name", query.Name, "error", err, "query", query.AsString())
		return
	}

	for _, warning := range warnings {
		slog.Warn(warning)
	}

	if matrix, ok := result.(model.Matrix); ok {
		metricsChan <- pa.handleMatrixResult(query.Name, &matrix)
	} else {
		slog.Warn("Query did not return matrix. Skipping.", "name", query.Name, "query", query.AsString())
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

	// replace NaN

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
	var lastValue float64 = 0
	for _, sample := range samples {
		thisValue := float64(sample.Value)
		// if NaN, use last non NaN value
		if math.IsNaN(thisValue) {
			thisValue = lastValue
		} else {
			lastValue = thisValue
		}

		values = append(values, thisValue)
	}
	return values
}
