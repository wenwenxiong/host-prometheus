package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"log"
	"net/http"
)

type handler struct{
	client monitoring.Interface
}

func newHandler( m monitoring.Interface) *handler {
	return &handler{ client: m}
}

func (h handler) ShowVisitorInfo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	country := vars["country"]
	_, err := fmt.Fprintf(writer, "This guy named %s, was coming from %s .", name, country)
	if err != nil {
		log.Println(err)
	}
}

func (h handler) CpuHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "(100-avg(irate(node_cpu_seconds_total{instance="+params.host+",mode=\"idle\"}[30m]))*100)"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) DealWithResponse( params reqParams, opt queryOptions, err error,res monitoring.Metric, writer http.ResponseWriter){

	if err != nil {
		if err.Error() == ErrNoHit {
			WriteAsJson(res, writer)
			return
		}

		HandleBadRequest(writer, nil, err)
		return
	}

	if opt.isRangeQuery() {
		res = h.client.GetMetricOverTime(params.expression, opt.start, opt.end, opt.step)
	} else {
		res = h.client.GetMetric(params.expression, opt.time)
	}

	if err != nil {
		HandleBadRequest(writer, nil, err)
	} else {
		WriteAsJson(res, writer)
	}

}

func WriteAsJson(res interface{}, writer http.ResponseWriter){

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(js)
	if err != nil {
		log.Println(err)
	}
}

func (h handler) Load1Handler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "node_load1{instance=~"+params.host+"}"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) Load5Handler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "node_load5{instance=~"+params.host+"}"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) Load15Handler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "node_load15{instance=~"+params.host+"}"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) MemHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "{__name__=~\"node_memory_MemFree_bytes|node_memory_Buffers_bytes|node_memory_Cached_bytes|node_memory_Slab_bytes|node_memory_MemTotal_bytes\",instance="+params.host+"}"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) DiskSizeHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "1-(node_filesystem_free_bytes{instance=~"+params.host+",fstype=~\"ext4|xfs\"} / node_filesystem_size_bytes{instance=~"+params.host+",fstype=~\"ext4|xfs\"})"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) DiskInodeHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "1-(node_filesystem_files_free{instance=~"+params.host+",fstype=~\"ext4|xfs\"} / node_filesystem_files{instance=~"+params.host+",fstype=~\"ext4|xfs\"})"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) IopsReadHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_disk_reads_completed_total{instance=~"+params.host+"}[30m])"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) IopsWriteHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_disk_writes_completed_total{instance=~"+params.host+"}[30m])"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) ThrReadHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_disk_read_bytes_total{instance=~"+params.host+"}[30m])"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) ThrWriteHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_disk_written_bytes_total{instance=~"+params.host+"}[30m])"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) NetworkRevHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_network_receive_bytes_total{instance=~"+params.host+",device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}

func (h handler) NetworkTransHandler(writer http.ResponseWriter, request *http.Request){
	var res monitoring.Metric
	params := parseRequestParams(request)
	params.host = BuildInstanceStr(params.host)
	params.expression = "irate(node_network_transmit_bytes_total{instance=~"+params.host+",device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8"
	opt, err := makeQueryOptions(params)
	h.DealWithResponse(params, opt, err, res, writer)
}