package hexagon

import "github.com/hanapedia/metrics-processor/internal/application/usecases/query"

const (
	PrimaryDurationBucket    query.MetricsName = "primary_adapter_duration_ms_bucket"
	PrimaryDurationCount     query.MetricsName = "primary_adapter_duration_ms_count"
	PrimaryDurationSum       query.MetricsName = "primary_adapter_duration_ms_sum"
	PrimaryAdapterInProgress query.MetricsName = "primary_adapter_in_progress"
	CallDurationBucket       query.MetricsName = "secondary_adapter_call_duration_ms_bucket"
	CallDurationCount        query.MetricsName = "secondary_adapter_call_duration_ms_count"
	CallDurationSum          query.MetricsName = "secondary_adapter_call_duration_ms_sum"
	CallAdaptiveTimeout      query.MetricsName = "adaptive_call_timeout_duration"
	TaskDurationBucket       query.MetricsName = "secondary_adapter_task_duration_ms_bucket"
	TaskDurationCount        query.MetricsName = "secondary_adapter_task_duration_ms_count"
	TaskDurationSum          query.MetricsName = "secondary_adapter_task_duration_ms_sum"
	TaskAdaptiveTimeout      query.MetricsName = "adaptive_task_timeout_duration"
)
