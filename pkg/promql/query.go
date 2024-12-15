package promql

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Query struct {
	Name string
	q    string
}

type Filter struct {
	label    string
	operator string
	value    string
}

func NewQuery(q string) *Query {
	return &Query{
		q: q,
	}
}

func (q *Query) SetName(name string) *Query {
	q.Name = name
	return q
}

func NewFilter(label, operator, value string) Filter {
	return Filter{
		label:    label,
		operator: operator,
		value:    value,
	}
}

func (q *Query) Filter(filters []Filter) *Query {
	q.q = fmt.Sprintf("%s{%s}", q.q, flattenFilters(filters))
	return q
}

func (f *Filter) AsString() string {
	return fmt.Sprintf(
		"%s%s\"%s\"",
		f.label,
		f.operator,
		f.value,
	)
}

func flattenFilters(filters []Filter) string {
	var filterStrs []string
	for _, filter := range filters {
		filterStrs = append(filterStrs, filter.AsString())
	}
	return strings.Join(filterStrs, ",")
}

func (q *Query) AsString() string {
	return q.q
}

func (q *Query) Group() *Query {
	q.q = fmt.Sprintf("(%s)", q.q)
	return q
}

func (q *Query) Rate(duration time.Duration) *Query {
	q.q = fmt.Sprintf("rate(%s[%s])", q.q, duration)
	return q
}

func (q *Query) IRate(duration time.Duration) *Query {
	q.q = fmt.Sprintf("irate(%s[%s])", q.q, duration)
	return q
}

func (q *Query) SumBy(byStrs []string) *Query {
	q.q = fmt.Sprintf("sum by (%s)(%s)", strings.Join(byStrs, ","), q.q)
	return q
}

func (q *Query) MinBy(byStrs []string) *Query {
	q.q = fmt.Sprintf("min by (%s)(%s)", strings.Join(byStrs, ","), q.q)
	return q
}

func (q *Query) HistogramQuantile(quantile float32) *Query {
	if quantile >= 1 || quantile <= 0 {
		slog.Warn("Invalid quantile, defaulting to 0.5")
		quantile = 0.5
	}
	q.q = fmt.Sprintf("histogram_quantile(%v,%s)", quantile, q.q)
	return q
}

func (q *Query) Subtract(aq *Query) *Query {
	q.q = fmt.Sprintf("%s - %s", q.q, aq.q)
	return q
}

func (q *Query) Divide(aq *Query) *Query {
	q.q = fmt.Sprintf("%s / %s", q.q, aq.q)
	return q
}

func (q *Query) MultiplyByConstant(c int) *Query {
	q.q = fmt.Sprintf("%s * %v", q.q, c)
	return q
}

func (q *Query) LabelReplace(target, source, pattern string) *Query {
	q.q = fmt.Sprintf(`label_replace(%s,"%s","$1","%s","%s")`, q.q, target, source, pattern)
	return q
}

func (q *Query) Offset(duration time.Duration) *Query {
	q.q = fmt.Sprintf("%s offset %s", q.q, duration)
	return q
}
