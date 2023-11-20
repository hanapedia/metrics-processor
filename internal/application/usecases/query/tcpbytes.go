package query

import (
	"time"

	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// CreateServerReadBytesQuery create query for bytes read by server
func CreateServerReadBytesQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createReadBytesQuery(namespace, rateDuration, "inbound", "dst").
		SetName("server_read_bytes")
}

// CreateServerWriteBytesQuery create query for bytes written by server
func CreateServerWriteBytesQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createWriteBytesQuery(namespace, rateDuration, "inbound", "dst").
		SetName("server_write_bytes")
}

// CreateClientReadBytesQuery create query for bytes read by client
func CreateClientReadBytesQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createReadBytesQuery(namespace, rateDuration, "outbound", "src").
		SetName("client_read_bytes")
}

// CreateClientWriteBytesQuery create query for bytes written by client
func CreateClientWriteBytesQuery(namespace string, rateDuration time.Duration) *promql.Query {
	return createWriteBytesQuery(namespace, rateDuration, "outbound", "src").
		SetName("client_write_bytes")
}

// createReadBytesQuery create query for bytes received
func createReadBytesQuery(namespace string, rateDuration time.Duration, direction, peer string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("namespace", "=", namespace),
		promql.NewFilter("direction", "=", direction),
		promql.NewFilter("peer", "=", peer),
	}
	// write metrics is used because it is recorded by the proxy and not the application
	return promql.NewQuery("tcp_write_bytes_total").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"deployment"})
}

// createWriteBytesQuery create query for bytes sent
func createWriteBytesQuery(namespace string, rateDuration time.Duration, direction, peer string) *promql.Query {
	filters := []promql.Filter{
		promql.NewFilter("namespace", "=", namespace),
		promql.NewFilter("direction", "=", direction),
		promql.NewFilter("peer", "=", peer),
	}
	// read metrics is used because it is recorded by the proxy and not the application
	return promql.NewQuery("tcp_read_bytes_total").
		Filter(filters).
		Rate(rateDuration).
		SumBy([]string{"deployment"})
}
