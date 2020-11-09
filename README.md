***  获取`prometheus-node-exporter, prometheus`的`host`主要指标并提供查询接口 ***

提供以下`url`查询对应指标

```
r.HandleFunc("/", ShowVisitorInfo)
r.HandleFunc("/host_cpu_usage", CpuHandler)
r.HandleFunc("/host_load1", Load1Handler)
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
r.HandleFunc("/host_net_bytes_transmitted", NetworkTransHandler)
```

指标释义如下

```
host:node-export
//host
http://192.168.50.71:9090/api/v1/query?
参数： 起始时间，结束时间，步长、时间范围[30m]、node
cpu:
cpu使用率(host_cpu_usage)
query=(100-avg(irate(node_cpu_seconds_total{instance="'+itemIp+':9100",mode="idle"}[30m]))*100)
cpu负载
1分钟负载(host_load1)： node_load1{instance=~"$node"}
5分钟负载(host_load5)： node_load5{instance=~"$node"}
15分钟负载(host_load15)： node_load15{instance=~"$node"}
内存：
内存使用率(host_memory_usage)
query={__name__=~"node_memory_MemFree_bytes|node_memory_Buffers_bytes|node_memory_Cached_bytes|node_memory_Slab_bytes|node_memory_MemTotal_bytes",instance="$node"}
磁盘
磁盘使用率
size(host_disk_size_usage): 1-(node_filesystem_free_bytes{instance=~'$node',fstype=~"ext4|xfs"} / node_filesystem_size_bytes{instance=~'$node',fstype=~"ext4|xfs"})
inodes(host_disk_inode_usage)：1-(node_filesystem_free_bytes{instance=~'$node',fstype=~"ext4|xfs"} / node_filesystem_size_bytes{instance=~'$node',fstype=~"ext4|xfs"})
磁盘iops
读iops(host_disk_read_iops)： irate(node_disk_reads_completed_total{instance=~"$node"}[30m])
写iops(host_disk_write_iops)： irate(node_disk_writes_completed_total{instance=~"$node"}[30m])
磁盘吞吐量
读(host_disk_read_throughput)：irate(node_disk_read_bytes_total{instance=~"$node"}[30m])
写(host_disk_write_throughput)：irate(node_disk_written_bytes_total{instance=~"$node"}[30m])
网络：
网络读写带宽
读（接收）(host_net_bytes_received)：irate(node_network_receive_bytes_total{instance=~'$node',device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8
写（发送）(host_net_bytes_transmitted)：irate(node_network_transmit_bytes_total{instance=~'$node',device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8
```







