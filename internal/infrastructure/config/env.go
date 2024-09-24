package config

import (
	"os"
	"strconv"
	"sync"
)

type EnvVars struct {
	METRICS_QUERY_ENDPOINT string
	END_TIME               string
	DURATION               string
	RATE_DURATION          string
	STEP                   string
	AWS_REGION             string
	S3_BUCKET              string
	S3_BUCKET_DIR          string
	K6_TEST_NAME           string
	NAMESPACE              string
	WORKLOAD_CONTAINERS    string
}

var defaults = EnvVars{
	METRICS_QUERY_ENDPOINT: "http://localhost:9090",
	END_TIME:               "",
	DURATION:               "30m",
	RATE_DURATION:          "5m",
	STEP:                   "15s",
	AWS_REGION:             "ap-northeast-1",
	S3_BUCKET:              "test",
	S3_BUCKET_DIR:          "test",
	K6_TEST_NAME:           "test",
	NAMESPACE:              "emulation",
	WORKLOAD_CONTAINERS:    "server|redis",
}

var envVars *EnvVars
var once sync.Once

func GetEnvs() *EnvVars {
	once.Do(func() {
		envVars = loadEnvVariables()
	})
	return envVars
}

func loadEnvVariables() *EnvVars {
	return &EnvVars{
		METRICS_QUERY_ENDPOINT: readEnv("METRICS_QUERY_ENDPOINT", defaults.METRICS_QUERY_ENDPOINT),
		END_TIME:               readEnv("END_TIME", defaults.END_TIME),
		DURATION:               readEnv("DURATION", defaults.DURATION),
		RATE_DURATION:          readEnv("RATE_DURATION", defaults.RATE_DURATION),
		STEP:                   readEnv("STEP", defaults.STEP),
		AWS_REGION:             readEnv("AWS_REGION", defaults.AWS_REGION),
		S3_BUCKET:              readEnv("S3_BUCKET", defaults.S3_BUCKET),
		S3_BUCKET_DIR:          readEnv("S3_BUCKET_DIR", defaults.S3_BUCKET_DIR),
		K6_TEST_NAME:           readEnv("K6_TEST_NAME", defaults.K6_TEST_NAME),
		NAMESPACE:              readEnv("NAMESPACE", defaults.NAMESPACE),
		WORKLOAD_CONTAINERS:    readEnv("WORKLOAD_CONTAINERS", defaults.WORKLOAD_CONTAINERS),
	}
}

func readEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func readBoolEnv(key string, defaultValue bool) bool {
	boolValue := defaultValue
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return boolValue
		}
		return parsed
	}
	return boolValue
}
