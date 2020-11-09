package apiserver

import(
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring/prometheus"
	"github.com/wenwenxiong/host-prometheus/pkg/constant"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type APIServer struct {
	ListenPort int
	Endpoint string
	// monitoring client set
	MonitoringClient monitoring.Interface
	Server *http.Server
}

func NewApiServer() (*APIServer, error){
	s := &APIServer{
		ListenPort: constant.ListenPort,
		Endpoint: constant.Endpoint,
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

	router := mux.NewRouter()
	RegisterRoutes(router)
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

func (s *APIServer) Run() (err error) {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	go func() {
		// Block until we receive our signal.
		<-c
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		s.Server.Shutdown(ctx)
	}()

	log.Println("Start listening on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return  err
}