package controlplane

// gin router for control plane API
import (
	"github.com/gin-gonic/gin"
)

func GetAdminRouter() *gin.Engine {
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