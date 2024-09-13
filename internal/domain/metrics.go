package domain

type QueryName = string

type MetricsMatrix struct {
	Name string `json:"name"`
	// LabelType indicates what the string keys in Matrix represent. e.g. deployment
	Matrix      map[string][]float64         `json:"matrix,omitempty"`
	Timestamps  []int64                      `json:"timestamps"`
}
