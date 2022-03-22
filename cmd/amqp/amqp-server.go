package main

import (
	"github.com/valinurovam/garagemq/config"
	"github.com/valinurovam/garagemq/metrics"
	"github.com/valinurovam/garagemq/server"
	_ "net/http/pprof"
	"time"
)

func main() {
	cfg, _ := config.CreateDefault()
	metrics.NewTrackRegistry(15, time.Second, false)
	srv := server.NewServer(cfg.TCP.IP, `5555`, cfg.Proto, cfg)
	srv.Start()
}
