package domain

type QueryName = string

type Metrics struct {
	QueryConfig *Config
	// Data holds the timeseries in a map
	// should look something like;
	/*
		{
			metrics_name : {
				deployment_name: []float64
				deployment_name: []float64
			},
			metrics_name : {
				deployment_name: []float64
				deployment_name: []float64
			},
		},
	*/
	Data map[QueryName]MetricsMatrix
}

type MetricsMatrix struct {
	// LabelType indicates what the string keys in Matrix represent. e.g. deployment
	LabelType  string
	Matrix     map[string][]float64
	Timestamps []int64
}
