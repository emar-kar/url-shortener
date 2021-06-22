package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emar-kar/urlshortener/pkg/service"
)

// Handler represents an object to communicate with services and handle
// project functionality.
type Handler struct {
	services *service.Service
}

// NewHandler creates handler object with the given services.
func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

// InitRoutes initialize gin.Engine to handle URI routes.
func (h *Handler) InitRoutes(ginMode string) *gin.Engine {
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Use standart gin logger.
	router.Use(gin.Logger())

	// Override panic. Recovery will transfer user to the 500 error page.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if _, ok := recovered.(string); ok {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{
				"title": "500 error",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// Integrate static files.
	router.StaticFS("web", http.Dir("web"))
	router.LoadHTMLGlob("web/*/*.html")

	// If gin.Engine cannot determine URI path it will transfer user to the 404 page.
	router.NoRoute(h.notFound)

	router.GET("/", h.mainHandler)
	router.GET("/:url", h.redirectHandler)
	router.GET("/statistics", h.statisticsHandler)
	router.POST("/generate", h.generateHandler)

	// API group.
	appAPI := router.Group("api")
	{
		appAPI.POST("/generate", h.generateLink)
		appAPI.GET("/statistics", h.getStatistics)
	}

	return router
}
