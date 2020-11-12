package options

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wenwenxiong/host-prometheus/pkg/apiserver"
	"github.com/wenwenxiong/host-prometheus/pkg/client/cache"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring/prometheus"
	"github.com/wenwenxiong/host-prometheus/pkg/client/mysql"
	"github.com/wenwenxiong/host-prometheus/pkg/cron"
	cliflag "k8s.io/component-base/cli/flag"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultListenPort = "9111"
)

type ServerRunOptions struct {
	ConfigFile              string
	ListenPort				string
	*apiserver.Config
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		ListenPort: DefaultListenPort,
		Config:      apiserver.New(),
	}

	return s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	s.RedisOptions.AddFlags(fss.FlagSet("redis"), s.RedisOptions)
	s.MonitoringOptions.AddFlags(fss.FlagSet("monitoring"), s.MonitoringOptions)
	s.MysqlOptions.AddFlags(fss.FlagSet("mysql"), s.MysqlOptions)
	return fss
}

func (s *ServerRunOptions) Validate() []error {
	var errors []error
	errors = append(errors, s.MonitoringOptions.Validate()...)
	return errors
}

func NewApiServer(listenPort string,endpoint string, sr *ServerRunOptions, stopCh <-chan struct{}) (*apiserver.APIServer, error){
	port, err := strconv.Atoi(listenPort)
	if err != nil {
		return nil, fmt.Errorf("listenPort must be right int format, error: %v", err)
	}
	s := &apiserver.APIServer{
		ListenPort: port,
		Endpoint: endpoint,
		Config: sr.Config,
	}
	if (strings.TrimSpace(s.Endpoint)) == "" {
		return nil, fmt.Errorf("moinitoring service address MUST not be empty, please check config endpoint")
	} else {
		monitoringClient, err := prometheus.NewPrometheus(s.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to prometheus, please check prometheus status, error: %v", err)
		}
		s.MonitoringClient = monitoringClient
	}

	var cacheClient cache.Interface
	if sr.RedisOptions == nil || len(sr.RedisOptions.Host) == 0 {
		return nil, fmt.Errorf("redis service address MUST not be empty, please check config file")
	} else {
		cacheClient, err = cache.NewRedisClient(sr.RedisOptions, stopCh)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to redis service, please check redis status, error: %v", err)
		}
		s.RedisClient = cacheClient
	}

	if sr.MysqlOptions == nil || len(sr.MysqlOptions.Host) == 0 {
		return nil, fmt.Errorf("mysql service address MUST not be empty, please check config file")
	} else {
		mysqlClient, err := mysql.NewMySQLClient(sr.MysqlOptions, stopCh)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mysql service, please check mysql status, error: %v", err)
		}
		s.MysqlClient = mysqlClient
	}

	cronSchedule, err := cron.NewCron()
	if err != nil {
		return nil, fmt.Errorf("failed to create cronSchedule, error: %v", err)
	}
	s.Cron = cronSchedule

	router := mux.NewRouter()
	apiserver.RegisterRoutes(router, s.MonitoringClient, s.RedisClient)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.ListenPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router, // Pass our instance of gorilla/mux in.
	}
	s.Server = srv

	return s,nil
}