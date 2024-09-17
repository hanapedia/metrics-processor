package query

import "time"

type MetricsName string

func (m MetricsName) AsString() string {
	return string(m)
}

const SCRAPE_ITERVAL = 15 * time.Second
