package hexagon

type MetricsName string

const (
	PrimaryDurationBucket    MetricsName = "primary_adapter_duration_ms_bucket"
	PrimaryDurationCount     MetricsName = "primary_adapter_duration_ms_count"
	PrimaryDurationSum       MetricsName = "primary_adapter_duration_ms_sum"
	PrimaryAdapterInProgress MetricsName = "primary_adapter_in_progress"
	CallDurationBucket       MetricsName = "secondary_adapter_call_duration_ms_bucket"
	CallDurationCount        MetricsName = "secondary_adapter_call_duration_ms_count"
	CallDurationSum          MetricsName = "secondary_adapter_call_duration_ms_sum"
	TaskDurationBucket       MetricsName = "secondary_adapter_task_duration_ms_bucket"
	TaskDurationCount        MetricsName = "secondary_adapter_task_duration_ms_count"
	TaskDurationSum          MetricsName = "secondary_adapter_task_duration_ms_sum"
)

func (m MetricsName) AsString() string {
	return string(m)
}
