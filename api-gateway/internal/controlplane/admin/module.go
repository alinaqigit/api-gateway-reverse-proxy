package admin

import (
	"api-gateway/internal/db"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Module struct {
	controller *Controller
	jwtManager *JWTManager
}

func NewModule(conn *sql.DB, jwtSecret string) *Module {
	repository := NewRepository(db.New(conn)) // Pass actual db.Queries instance here
	jwtManager, err := NewJWTManagerFromSecret(jwtSecret)
	if err != nil {
		panic(err)
	}

	service := NewService(repository, jwtManager)
	controller := NewController(service)

	return &Module{
		controller: controller,
		jwtManager: jwtManager,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Public route for login
	router.POST("/login", m.controller.LoginAdmin)

	// Protect all other admin routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware(m.jwtManager))

	// Elevated-only routes
	elevated := protected.Group("/")
	elevated.Use(RequireElevated())
	elevated.POST("/", m.controller.CreateAdmin)
	elevated.GET("/:id", m.controller.GetAdmin)
	elevated.DELETE("/:id", m.controller.DeleteAdmin)

	// Authenticated routes
	protected.PUT("/:id", m.controller.UpdateAdmin)
}
