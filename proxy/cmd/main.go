package main

import (
	"log"

	"github.com/mmarci96/kubernetes-controller-go/proxy/pkg/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	select {}
}
