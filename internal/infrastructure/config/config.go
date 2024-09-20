package config

import (
	"log/slog"
	"strconv"
	"strings"
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
	// Try to parse the input as a float for potential sub-second precision
	unixTimeFloat, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		// Log a warning and return the current time if parsing fails
		slog.Warn("Failed to parse END_TIME. Using time.Now()", "err", err)
		return time.Now()
	}

	// Check if the timestamp has sub-second precision (i.e., contains a dot)
	if strings.Contains(timestamp, ".") {
		// Separate the integer seconds and the fractional milliseconds
		seconds := int64(unixTimeFloat)
		nanoSeconds := int64((unixTimeFloat - float64(seconds)) * 1e9)
		return time.Unix(seconds, nanoSeconds)
	}

	// If it's an integer, check if it's in milliseconds (13 digits)
	if len(timestamp) == 13 {
		return time.UnixMilli(int64(unixTimeFloat))
	}

	// Fallback for a standard Unix timestamp in seconds
	return time.Unix(int64(unixTimeFloat), 0)
}
