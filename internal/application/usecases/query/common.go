package query

import (
	"fmt"
	"time"
)

type MetricsName string

func (m MetricsName) AsString() string {
	return string(m)
}

const SCRAPE_ITERVAL = 15 * time.Second

type RateConfig struct {
	Name     string
	Duration time.Duration
	// Specify whether rate or irate
	IsInstant bool
}

func (rc RateConfig) AddSuffix(name string) string {
	if rc.IsInstant {
		return fmt.Sprintf("%s_irate_%s", name, rc.Duration.String())
	}
	return fmt.Sprintf("%s_rate_%s", name, rc.Duration.String())
}
