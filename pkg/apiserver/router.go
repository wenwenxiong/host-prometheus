package apiserver

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", ShowVisitorInfo)
	r.HandleFunc("/host_cpu_usage", CpuHandler)
	/*r.HandleFunc("/host_load1", Load1Handler)
	r.HandleFunc("/host_load5", Load5Handler)
	r.HandleFunc("/host_load15", Load15Handler)
	r.HandleFunc("/host_memory_usage", MemHandler)
	r.HandleFunc("/host_disk_size_usage", DiskSizeHandler)
	r.HandleFunc("/host_disk_inode_usage", DiskInodeHandler)
	r.HandleFunc("/host_disk_read_iops", IopsReadHandler)
	r.HandleFunc("/host_disk_write_iops", IopsWriteHandler)
	r.HandleFunc("/host_disk_read_throughput", ThrReadHandler)
	r.HandleFunc("/host_disk_write_throughput", ThrWriteHandler)
	r.HandleFunc("/host_net_bytes_received", NetworkRevHandler)
	r.HandleFunc("/host_net_bytes_transmitted", NetworkTransHandler)*/
}