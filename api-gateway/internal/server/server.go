package server

import (
	controlplane "api-gateway/internal/controlplane"
	dataplane "api-gateway/internal/dataplane"
	"database/sql"
	"net/http"
	_ "github.com/lib/pq"
)

type ServerParams struct {
	Env      string
	DBString string
	JWTSecret string
}

func GetServer(params ServerParams) *http.ServeMux {

	// Initialize database connection here and pass it to control plane modules
	conn, err := sql.Open("postgres", params.DBString)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// --- Root mux ---
	mux := http.NewServeMux()

	// Specific route first
	mux.Handle(
		"/v1/admin/",
		controlplane.GetAdminRouter(controlplane.AdminRouterParams{
			CONN:      conn,
			ENV:       params.Env,
			JWTSECRET: params.JWTSecret,
		}),
	)

	// General route after

	proxy, err := dataplane.GetGatewayRouter("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	mux.Handle(
		"/v1/",
		http.StripPrefix(
			"/v1",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				proxy.ServeHTTP(w, r)
			}),
		),
	)

	// Finally return the mux
	return mux

}
