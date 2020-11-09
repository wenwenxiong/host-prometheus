package monitoring

import "time"

type Interface interface {
	GetMetric(expr string, time time.Time) Metric
	GetMetricOverTime(expr string, start, end time.Time, step time.Duration) Metric
}
