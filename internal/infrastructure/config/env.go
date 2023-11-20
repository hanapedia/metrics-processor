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
	STEP                   string
	S3_BUCKET              string
	TEST_NAME              string
	NAMESPACE              string
	WORKLOAD_CONTAINERS    string
}

var defaults = EnvVars{
	METRICS_QUERY_ENDPOINT: "http://localhost:9090",
	END_TIME:               "",
	DURATION:               "30m",
	STEP:                   "15s",
	S3_BUCKET:              "test",
	TEST_NAME:              "test",
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
		STEP:                   readEnv("STEP", defaults.STEP),
		S3_BUCKET:              readEnv("S3_BUCKET", defaults.S3_BUCKET),
		TEST_NAME:              readEnv("TEST_NAME", defaults.TEST_NAME),
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
