package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwf/go-speedtest/config"
	"github.com/lwf/go-speedtest/servers"
	"github.com/lwf/go-speedtest/transfer"
	"github.com/lwf/go-speedtest/transfer/downloader"
	"github.com/lwf/go-speedtest/transfer/uploader"
	"github.com/lwf/go-speedtest/unit"
)

func main() {
	var bytes bool
	var id int
	flag.BoolVar(&bytes, "bytes", false, "Report in bytes rather than bits")
	flag.IntVar(&id, "id", 0, "Id of server to use")
	flag.Parse()

	log.Println("Retrieving configuration from speedtest.net...")
	config := config.GetConfig()

	var url string
	s := servers.GetServers()
	if id == 0 {
		log.Println("Retrieving server list from speedtest.net...")
		servers.SortByDistance(config, s.Servers)
		s2 := s.Servers[0:5]
		servers.SortByLatency(config, s2)
		url = s2[0].Url
		log.Printf("Selected remote %s based on distance and latency", url)
	} else {
		for _, server := range s.Servers {
			if server.Id == id {
				url = server.Url
				break
			}
		}
		if url == "" {
			log.Fatal("Could not find server for id", id)
		}
		log.Println("Using remote", url)
	}

	manager := transfer.NewManager(url)
	hookSignals(manager.Close)
	log.Printf("Performing test of download/upload speed")
	var labels unit.Labels
	var m float64
	if bytes == true {
		labels = unit.Bytes
		m = 1
	} else {
		labels = unit.Bits
		m = 8
	}
	workers := []transfer.Worker{
		&downloader.Downloader{},
		&uploader.Uploader{},
	}
	for _, worker := range workers {
		v := measure(manager, worker, func() {
			fmt.Printf(".")
		})
		if v > 0 {
			speed, unit := unit.Unit(v*m, labels)
			fmt.Printf("\n%s speed: %.2f %s/s\n", worker, speed, unit)
		}
	}
}

func measure(manager *transfer.Manager, worker transfer.Worker, progress func()) float64 {
	t1 := time.Now()
	v := manager.Run(worker, 4, progress)
	t2 := time.Since(t1).Seconds()
	return float64(v) / t2
}

func hookSignals(h func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		h()
	}()
}
