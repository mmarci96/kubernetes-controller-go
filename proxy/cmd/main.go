package main

import (
	"log"

	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/server"
	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/watcher"
)

func main() {
	go watcher.WatchEndpointSlices()
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	select {}
}
