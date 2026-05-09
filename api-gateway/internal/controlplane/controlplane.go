package controlplane

// gin router for control plane API
import (
	"api-gateway/internal/controlplane/admin"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type AdminRouterParams struct {
	CONN *sql.DB
	ENV string
	JWTSECRET string
}

func GetAdminRouter(params AdminRouterParams) *gin.Engine {

	switch params.ENV {
		case "dev":
			gin.SetMode(gin.DebugMode) // Suppress GIN debug output even in dev
		case "prod":
			gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Define your control plane API routes here
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Admin routes
	adminModule := admin.NewModule(params.CONN, params.JWTSECRET)
	adminGroup := router.Group("/v1/admin");
	adminModule.RegisterRoutes(adminGroup);

	// mount other control plane API routes here

	

	return router
}