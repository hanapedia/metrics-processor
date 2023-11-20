package domain

import "time"

type Config struct {
	MetricsQueryEndpoing string
	EndTime              time.Time
	Duration             time.Duration
	Step                 time.Duration
	AWSRegion            string
	S3Bucket             string
	TestName             string
	Namespace            string
	WorkloadContainers   string
}
