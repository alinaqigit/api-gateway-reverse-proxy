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
	env := os.Getenv("ENV")

	fmt.Println("Reading Environment Variables...")
	if hostname == "" || listenPort == "" || env == "" {
		panic("API_GATEWAY_HOSTNAME, API_GATEWAY_LISTEN_PORT, and ENV environment variables must be set")
	}

	fmt.Printf("Hostname: %s, Listen Port: %s, Environment: %s\n", hostname, listenPort, env)

	// Get the main server mux
	server := internal.GetServer(env)
	fmt.Println("Server Mux initialized...")

	// Start the server
	http.ListenAndServe(hostname + ":" + listenPort, server)
}