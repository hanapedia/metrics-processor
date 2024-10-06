package domain

import (
	"github.com/prometheus/common/model"
)

type QueryName = string

type MetricsMatrix struct {
	Name   string                        `json:"name"`
	Matrix map[string][]model.SamplePair `json:"matrix"`
	End    float64                       `json:"end"`
}
