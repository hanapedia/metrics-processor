package container

import (
	"github.com/hanapedia/metrics-processor/internal/application/usecases/query"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateContainerRestartsQuery create query for container restarts
// it obtains gauge metric for how many times containers were restarted between observation points.
func CreateContainerRestartsQuery(filters []promql.Filter) *promql.Query {
	cur := promql.NewQuery(KubePodContainerRestarts.AsString()).
		Filter(filters).
		SumBy([]string{"pod"})

	prev := promql.NewQuery(KubePodContainerRestarts.AsString()).
		Filter(filters).
		Offset(query.SCRAPE_ITERVAL).
		SumBy([]string{"pod"})

	return cur.SetName("container_restarts").Subtract(prev)
}
