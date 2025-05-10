package main

import (
	"log"

	kubeclient "github.com/mmarci96/kubernetes-controller-go/proxy/pkg/kube_client"
	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/server"
)

func main() {
	go kubeclient.WatchEndpointSlices()
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	select {}
}
