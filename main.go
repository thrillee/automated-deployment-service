package main

import (
	"log"
	"os"

	"github.com/thrillee/automated-deployment-service/internals/server"
)

func main() {
	portString := os.Getenv("HTTP_PORT")
	if portString == "" {
		log.Fatal("HTTP_PORT is required")
	}

	httpServer := server.HttpAPIServer{
		ListenAddr: ":" + portString,
	}
	httpServer.Run()
}
