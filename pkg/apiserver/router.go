package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
)

func RegisterRoutes(r *mux.Router, cli monitoring.Interface) {

	handler := newHandler(cli)
	r.HandleFunc("/", handler.ShowVisitorInfo)
	r.HandleFunc("/host_cpu_usage", handler.CpuHandler)
	r.HandleFunc("/host_load1", handler.Load1Handler)
	r.HandleFunc("/host_load5", handler.Load5Handler)
	r.HandleFunc("/host_load15", handler.Load15Handler)
	r.HandleFunc("/host_memory_usage", handler.MemHandler)
	r.HandleFunc("/host_disk_size_usage", handler.DiskSizeHandler)
	r.HandleFunc("/host_disk_inode_usage", handler.DiskInodeHandler)
	r.HandleFunc("/host_disk_read_iops", handler.IopsReadHandler)
	r.HandleFunc("/host_disk_write_iops", handler.IopsWriteHandler)
	r.HandleFunc("/host_disk_read_throughput", handler.ThrReadHandler)
	r.HandleFunc("/host_disk_write_throughput", handler.ThrWriteHandler)
	r.HandleFunc("/host_net_bytes_received", handler.NetworkRevHandler)
	r.HandleFunc("/host_net_bytes_transmitted", handler.NetworkTransHandler)
}