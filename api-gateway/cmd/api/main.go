package main

import (
	"api-gateway/internal"
	"fmt"

	"net/http"
	"os"
	// Import the control plane API router from the control-plane package
)

func main() {

	fmt.Println("Starting API Gateway...")

	// Take the configuration from environment variables
	hostname := os.Getenv("API_GATEWAY_HOSTNAME")
	listenPort := os.Getenv("API_GATEWAY_LISTEN_PORT")

	if hostname == "" || listenPort == "" {
		panic("API_GATEWAY_HOSTNAME and API_GATEWAY_LISTEN_PORT environment variables must be set")
	}

	// Get the main server mux
	mux := internal.GetServer(hostname)
	
	
	// Start the server
	fmt.Printf("API Gateway listening on https://%s:%s", hostname, listenPort)
	http.ListenAndServe(":" + listenPort, mux)
}