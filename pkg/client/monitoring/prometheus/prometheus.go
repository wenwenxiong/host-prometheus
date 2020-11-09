package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"time"
)
// prometheus implements monitoring interface backed by Prometheus
type prometheus struct {
	client apiv1.API
}

func NewPrometheus(endpoint string) (monitoring.Interface, error) {
	cfg := api.Config{
		Address: endpoint,
	}

	client, err := api.NewClient(cfg)
	return prometheus{client: apiv1.NewAPI(client)}, err
}

func (p prometheus) GetMetric(expr string, ts time.Time) monitoring.Metric {
	var parsedResp monitoring.Metric

	value,_, err := p.client.Query(context.Background(), expr, ts)
	if err != nil {
		parsedResp.Error = err.Error()
	} else {
		parsedResp.MetricData = parseQueryResp(value)
	}

	return parsedResp
}

func (p prometheus) GetMetricOverTime(expr string, start, end time.Time, step time.Duration) monitoring.Metric {
	timeRange := apiv1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	value,_, err := p.client.QueryRange(context.Background(), expr, timeRange)

	var parsedResp monitoring.Metric
	if err != nil {
		parsedResp.Error = err.Error()
	} else {
		parsedResp.MetricData = parseQueryRangeResp(value)
	}
	return parsedResp
}

func parseQueryRangeResp(value model.Value) monitoring.MetricData {
	res := monitoring.MetricData{MetricType: monitoring.MetricTypeMatrix}

	data, _ := value.(model.Matrix)

	for _, v := range data {
		mv := monitoring.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		for _, k := range v.Values {
			mv.Series = append(mv.Series, monitoring.Point{float64(k.Timestamp) / 1000, float64(k.Value)})
		}

		res.MetricValues = append(res.MetricValues, mv)
	}

	return res
}

func parseQueryResp(value model.Value) monitoring.MetricData {
	res := monitoring.MetricData{MetricType: monitoring.MetricTypeVector}

	data, _ := value.(model.Vector)

	for _, v := range data {
		mv := monitoring.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		mv.Sample = &monitoring.Point{float64(v.Timestamp) / 1000, float64(v.Value)}

		res.MetricValues = append(res.MetricValues, mv)
	}

	return res
}
