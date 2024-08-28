package config

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/hanapedia/metrics-processor/internal/domain"
)

func NewConfigFromEnv() *domain.Config {
	var endTime time.Time
	if GetEnvs().END_TIME == "" {
		endTime = time.Now()
	} else {
		endTime = parseStringUnixMilliSecTimestamp(GetEnvs().END_TIME)
	}

	duration, err := time.ParseDuration(GetEnvs().DURATION)
	if err != nil {
		slog.Warn("Failed to parse DURATION. Using 30m", "err", err)
		duration = 30 * time.Minute
	}

	step, err := time.ParseDuration(GetEnvs().STEP)
	if err != nil {
		slog.Warn("Failed to parse STEP", "err", err)
		step = 15 * time.Second
	}

	return &domain.Config{
		MetricsQueryEndpoint: GetEnvs().METRICS_QUERY_ENDPOINT,
		EndTime:              endTime,
		Duration:             duration,
		Step:                 step,
		AWSRegion:            GetEnvs().AWS_REGION,
		S3Bucket:             GetEnvs().S3_BUCKET,
		S3BucketDir:          GetEnvs().S3_BUCKET_DIR,
		K6TestName:           GetEnvs().K6_TEST_NAME,
		Namespace:            GetEnvs().NAMESPACE,
		WorkloadContainers:   GetEnvs().WORKLOAD_CONTAINERS,
	}
}

func parseStringUnixMilliSecTimestamp(timestamp string) time.Time {
	unixTimeStr, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		slog.Warn("Failed to parse END_TIME. Using time.Now()", "err", err)
		return time.Now()
	}
	if len(timestamp) == 13 {
		return time.UnixMilli(unixTimeStr)
	}
	return time.Unix(unixTimeStr, 0)
}
