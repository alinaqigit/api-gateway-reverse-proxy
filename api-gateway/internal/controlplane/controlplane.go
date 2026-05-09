package controlplane

// gin router for control plane API
import (

	"github.com/gin-gonic/gin"
)

func GetAdminRouter(env string) *gin.Engine {

	switch env {
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

	// mount other control plane API routes here

	

	return router
}