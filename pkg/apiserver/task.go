package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/wenwenxiong/host-prometheus/pkg/client/cache"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"github.com/wenwenxiong/host-prometheus/pkg/client/mysql"
	"log"
	"time"
)

type Host struct {
	HostName string `json:"hostName" db:"name"`
	HostIp string `json:"hostIp" db:"ip"`
}

type HostMetrics struct {
	Host Host `json:"host"`
	CpuUsage float64 `json:"cpuUsage"`
	Load1 float64 `json:"load1"`
	Load5 float64 `json:"load5"`
	Load15 float64 `json:"load15"`
	MemUsage float64 `json:"memUsage"`
	/*DiskSizeUsage float64 `json:"diskSizeUsage"`
	DiskInodeUsage float64 `json:"diskInodeUsage"`
	DiskReadIops float64 `json:"diskReadIops"`
	DiskWriteIops float64 `json:"diskWriteIops"`
	DiskReadThroughput float64 `json:"diskReadThroughput"`
	DiskWriteThroughput float64 `json:"diskWriteThroughput"`
	NetBytesReceived float64 `json:"netBytesReceived"`
	NetBytesTransmitted float64 `json:"netBytesTransmitted"`*/
}


func Task(a *APIServer){
	hosts := QueryHost(a.MysqlClient)
	hms := GetherMetrics(hosts, a.MonitoringClient)
	saveCache(hms, a.RedisClient)
}

func QueryHost(s *mysql.Client) ([]Host) {
	  var hosts []Host
	_,err := s.Database().Select("ip","name").From("kubeops_api_host").Where(mysql.Eq("npg", 0)).Load(&hosts)
	if err != nil {
		log.Println(err)
	}
	return  hosts
}

func GetherMetrics(host [] Host, s monitoring.Interface) ([]HostMetrics){
	var hostMetrics []HostMetrics
	for _, h := range host {
		var ho HostMetrics
		ho.Host = h
		instance := "\""+h.HostIp+":9100\""
		//for cpu usage
		expression := "(100-avg(irate(node_cpu_seconds_total{instance="+instance+",mode=\"idle\"}[30m]))*100)"
		m := s.GetMetric(expression, time.Now())
		if len(m.MetricData.MetricValues)>0{
			ho.CpuUsage = m.MetricData.MetricValues[0].Sample[1]
		}else {
			ho.CpuUsage = 0.0
		}
		//for load usage
		expression = "node_load1{instance=~"+instance+"}"
		m = s.GetMetric(expression, time.Now())
		if len(m.MetricData.MetricValues)>0{
			ho.Load1 = m.MetricData.MetricValues[0].Sample[1]
		}else {
			ho.Load1 = 0.0
		}
		expression = "node_load5{instance=~"+instance+"}"
		m = s.GetMetric(expression, time.Now())
		if len(m.MetricData.MetricValues)>0{
			ho.Load5 = m.MetricData.MetricValues[0].Sample[1]
		}else {
			ho.Load5 = 0.0
		}
		expression = "node_load15{instance=~"+instance+"}"
		m = s.GetMetric(expression, time.Now())
		if len(m.MetricData.MetricValues)>0{
			ho.Load15 = m.MetricData.MetricValues[0].Sample[1]
		}else {
			ho.Load15 = 0.0
		}
		//for mem usage
		expression ="{__name__=~\"node_memory_MemFree_bytes|node_memory_Buffers_bytes|node_memory_Cached_bytes|node_memory_Slab_bytes|node_memory_MemTotal_bytes\",instance="+instance+"}"
		m = s.GetMetric(expression, time.Now())
		if len(m.MetricData.MetricValues)>0{
			var  memFree, memCache, memBuffer, memSlab,memTotal float64
			for _, mValue := range m.MetricData.MetricValues{
				if mValue.Metadata["__name__"] == "node_memory_MemFree_bytes"{
					memFree = mValue.Sample[1]/1000000;
				} else if mValue.Metadata["__name__"] == "node_memory_Cached_bytes" {
					memCache = mValue.Sample[1]/1000000;
				} else if mValue.Metadata["__name__"] == "node_memory_Buffers_bytes" {
					memBuffer = mValue.Sample[1]/1000000;
				} else if mValue.Metadata["__name__"] == "node_memory_Slab_bytes" {
					memSlab = mValue.Sample[1]/1000000;
				} else if mValue.Metadata["__name__"] == "node_memory_MemTotal_bytes" {
					memTotal = mValue.Sample[1]/1000000;
				}
			}
			memUsage := memFree+memCache + memBuffer + memSlab;
			ho.MemUsage = (1-(memUsage/memTotal))*100;
		}else {
			ho.MemUsage = 0.0
		}

		hostMetrics = append(hostMetrics,ho)
	}
	return  hostMetrics
}

func saveCache (hostMetrics []HostMetrics, redis cache.Interface){
	for _,hm := range hostMetrics{
		json, err := json.Marshal(hm)
		if err != nil {
			log.Println(err)
		}
		err = redis.Set(hm.Host.HostIp, string(json[:]), 0)
		if err != nil {
			log.Println(err)
		}

	}

}

func getCache (hostIp string, redis cache.Interface)(hm HostMetrics){
	val, err := redis.Get(hostIp)
	if err != nil {
		fmt.Println(err)
	}
	data := HostMetrics{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}