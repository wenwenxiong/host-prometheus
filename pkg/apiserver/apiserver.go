package apiserver

import (
	"context"
	"github.com/wenwenxiong/host-prometheus/pkg/client/cache"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"github.com/wenwenxiong/host-prometheus/pkg/client/mysql"
	"github.com/wenwenxiong/host-prometheus/pkg/cron"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type APIServer struct {
	ListenPort int
	Endpoint string
	Config *Config
	// monitoring client set
	MonitoringClient monitoring.Interface
	RedisClient cache.Interface
	MysqlClient *mysql.Client
	Server *http.Server
	Cron *cron.Cron
}

func (s *APIServer) Run() (err error) {

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	go func() {
		// Block until we receive our signal.
		<-c
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		s.Server.Shutdown(ctx)
	}()

	// Schedule Task

	log.Println("Start cron Scheduler Task on %s", time.Now().String())
	if _,e := s.Cron.Cron().Every(1).Minute().Do(Task,s); e != nil{
		log.Println(e)
	}
	s.Cron.Cron().StartAsync()
	log.Println("Start listening on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return  err
}
