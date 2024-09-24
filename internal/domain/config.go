package domain

import "time"

type Config struct {
	MetricsQueryEndpoint string
	EndTime              time.Time
	Duration             time.Duration
	RateDuration         time.Duration
	Step                 time.Duration
	AWSRegion            string
	S3Bucket             string
	S3BucketDir          string
	K6TestName           string
	Namespace            string
	WorkloadContainers   string
}
