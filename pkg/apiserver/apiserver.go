package apiserver

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring"
	"github.com/wenwenxiong/host-prometheus/pkg/client/monitoring/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

func NewApiServer(listenPort string,endpoint string) (*APIServer, error){
	port, err := strconv.Atoi(listenPort)
	if err != nil {
		return nil, fmt.Errorf("listenPort must be right int format, error: %v", err)
	}
	s := &APIServer{
		ListenPort: port,
		Endpoint: endpoint,
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
	RegisterRoutes(router, s.MonitoringClient)
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

	log.Println("Start listening on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return  err
}