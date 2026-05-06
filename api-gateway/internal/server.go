package internal

import (
	controlplane "api-gateway/internal/controlplane"
	dataplane "api-gateway/internal/dataplane"
	"net/http"
)

func GetServer(hostname string) *http.ServeMux {

	// --- Root mux ---
	mux := http.NewServeMux()

	// Specific route first
	mux.Handle("/v1/admin/", controlplane.GetAdminRouter())

	// General route after

	proxy, err := dataplane.GetGatewayRouter(hostname)
	if err != nil {
		panic(err)
	}

	mux.Handle("/v1/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))

	// Finally return the mux
	return mux

}