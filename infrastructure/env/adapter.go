package env

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/hanapedia/rca-experiment-runner/pkg/constants"
	"github.com/hanapedia/rca-experiment-runner/pkg/domain"

	"github.com/joho/godotenv"
)

func NewExperimentConfig() (*domain.ExperimentConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INFO]: No .env file found. Proceeding to read environmental variables without it.")
	}

	name := os.Getenv("EXPERIMENT_NAME")
	if name == "" {
		return nil, errors.New("EXPERIMENT_NAME must be set")
	}

	targetNamespace := os.Getenv("TARGET_NAMESPACE")
	if targetNamespace == "" {
		return nil, errors.New("TARGET_NAMESPACE must be set")
	}

	configMapName := os.Getenv("BATCH_CONFIG_MAP_NAME")
	if configMapName == "" {
		return nil, errors.New("BATCH_CONFIG_MAP_NAME must be set")
	}

	imageTag := os.Getenv("TAG")
	if imageTag == "" {
		imageTag = constants.DefaultImageTag
	}

	normalDurationStr := os.Getenv("NORMAL_DURATION")
	normalDuration, err := time.ParseDuration(normalDurationStr)
	if err != nil {
		return nil, errors.New("NORMAL_DURATION is not a valid duration")
	}

	injectionDurationStr := os.Getenv("INJECTION_DURATION")
	injectionDuration, err := time.ParseDuration(injectionDurationStr)
	if err != nil {
		return nil, errors.New("INJECTION_DURATION is not a valid duration")
	}

	latencyStr := os.Getenv("LATENCY")
	if latencyStr == "" {
		latencyStr = constants.DefaultLatency
	}
	latencyDuration, err := time.ParseDuration(latencyStr)
	if err != nil {
		return nil, errors.New("LATENCY is not a valid duration")
	}

	jitterStr := os.Getenv("JITTER")
	if jitterStr == "" {
		jitterStr = constants.DefaultJitter
	}
	jitterDuration, err := time.ParseDuration(jitterStr)
	if err != nil {
		return nil, errors.New("JITTER is not a valid duration")
	}

	config := &domain.ExperimentConfig{
		Name:               name,
		TargetNamespace:    targetNamespace,
		BatchConfigMapName: configMapName,
		ImageTag:           imageTag,
		NormalDuration:     normalDuration,
		InjectionDuration:  injectionDuration,
		Latency:            latencyDuration,
		Jitter:             jitterDuration,
	}

	return config, nil
}
