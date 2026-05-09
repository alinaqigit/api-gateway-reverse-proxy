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
	dbString := os.Getenv("API_GATEWAY_DATABASE_URL")
	env := os.Getenv("ENV")
	JWTSecret := os.Getenv("API_GATEWAY_JWT_SECRET")

	fmt.Println("Reading Environment Variables...")
	if hostname == "" || listenPort == "" || env == "" || dbString == "" || JWTSecret == "" {
		panic("API_GATEWAY_HOSTNAME, API_GATEWAY_LISTEN_PORT, API_GATEWAY_DATABASE_URL, API_GATEWAY_JWT_SECRET, and ENV environment variables must be set")
	}

	fmt.Printf("Hostname: %s, Listen Port: %s, Environment: %s\n", hostname, listenPort, env)

	// Get the main server mux
	server := internal.GetServer(internal.ServerParams{
		Env:      env,
		DBString: dbString,
		JWTSecret: JWTSecret,
	})
	fmt.Println("Server Mux initialized...")

	// Start the server
	http.ListenAndServe(hostname + ":" + listenPort, server)
}