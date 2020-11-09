package main

import (
	"github.com/wenwenxiong/host-prometheus/cmd/app"
	"log"
)

func main() {

	cmd := app.NewPrometheusCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
