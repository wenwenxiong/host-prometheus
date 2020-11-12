package apiserver

import (
	"github.com/pkg/errors"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultStep   = 10 * time.Minute
	DefaultFilter = ".*"
	DefaultOrder  = OrderDescending
	DefaultPage   = 1
	DefaultLimit  = 5
	OrderAscending  = "asc"
	OrderDescending = "desc"

	ErrNoHit           = "'end' or 'time' must be after the namespace creation time."
	ErrParamConflict   = "'time' and the combination of 'start' and 'end' are mutually exclusive."
	ErrInvalidStartEnd = "'start' must be before 'end'."
	ErrInvalidPage     = "Invalid parameter 'page'."
	ErrInvalidLimit    = "Invalid parameter 'limit'."
)

type reqParams struct {
	time             string
	start            string
	end              string
	step             string
	target           string
	order            string
	page             string
	limit            string
	metricFilter     string
	resourceFilter   string
	host         string
	expression       string
	metric           string
}

type queryOptions struct {
	metricFilter string
	namedMetrics []string

	start time.Time
	end   time.Time
	time  time.Time
	step  time.Duration

	target     string
	identifier string
	order      string
	page       int
	limit      int

	option monitoring.QueryOption
}

func (q queryOptions) isRangeQuery() bool {
	return q.time.IsZero()
}

func parseRequestParams(req *http.Request) reqParams {
	var r reqParams
	r.time = req.URL.Query().Get("time")
	r.start = req.URL.Query().Get("start")
	r.end = req.URL.Query().Get("end")
	r.step = req.URL.Query().Get("step")
	r.target = req.URL.Query().Get("sort_metric")
	r.order = req.URL.Query().Get("sort_type")
	r.page = req.URL.Query().Get("page")
	r.limit = req.URL.Query().Get("limit")
	r.metricFilter = req.URL.Query().Get("metrics_filter")
	r.resourceFilter = req.URL.Query().Get("resources_filter")
	r.host = req.URL.Query().Get("host")
	r.metric = req.URL.Query().Get("metric")
	return r
}

func makeQueryOptions(r reqParams) (q queryOptions, err error) {
	if r.resourceFilter == "" {
		r.resourceFilter = DefaultFilter
	}

	q.metricFilter = r.metricFilter
	if r.metricFilter == "" {
		q.metricFilter = DefaultFilter
	}

	q.namedMetrics = NodeMetrics
	q.option = monitoring.HostOption{
		ResourceFilter: r.resourceFilter,
		HostName:       r.host,
	}

	// Parse time params
	if r.start != "" && r.end != "" {
		startInt, err := strconv.ParseInt(r.start, 10, 64)
		if err != nil {
			return q, err
		}
		q.start = time.Unix(startInt, 0)

		endInt, err := strconv.ParseInt(r.end, 10, 64)
		if err != nil {
			return q, err
		}
		q.end = time.Unix(endInt, 0)

		if r.step == "" {
			q.step = DefaultStep
		} else {
			q.step, err = time.ParseDuration(r.step)
			if err != nil {
				return q, err
			}
		}

		if q.start.After(q.end) {
			return q, errors.New(ErrInvalidStartEnd)
		}
	} else if r.start == "" && r.end == "" {
		if r.time == "" {
			q.time = time.Now()
		} else {
			timeInt, err := strconv.ParseInt(r.time, 10, 64)
			if err != nil {
				return q, err
			}
			q.time = time.Unix(timeInt, 0)
		}
	} else {
		return q, errors.Errorf(ErrParamConflict)
	}

	// Parse sorting and paging params
	if r.target != "" {
		q.target = r.target
		q.page = DefaultPage
		q.limit = DefaultLimit
		q.order = r.order
		if r.order != OrderAscending {
			q.order = DefaultOrder
		}
		if r.page != "" {
			q.page, err = strconv.Atoi(r.page)
			if err != nil || q.page <= 0 {
				return q, errors.New(ErrInvalidPage)
			}
		}
		if r.limit != "" {
			q.limit, err = strconv.Atoi(r.limit)
			if err != nil || q.limit <= 0 {
				return q, errors.New(ErrInvalidLimit)
			}
		}
	}

	return q, nil
}

func BuildInstanceStr(hostIp string)(host string){
	var res string
	res = "\""+hostIp+":9100\""
	return res
}