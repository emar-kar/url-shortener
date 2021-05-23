package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emar-kar/urlshortener/pkg/service"
)

// TODO: add custom logger

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitRoutes(ginMode string) *gin.Engine {
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if _, ok := recovered.(string); ok {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{
				"title": "500 error",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.StaticFS("web", http.Dir("web"))
	router.LoadHTMLGlob("web/*/*.html")

	router.NoRoute(h.notFound)

	router.GET("/", h.mainHandler)
	router.GET("/:url", h.redirectHandler)
	router.GET("/statistics", h.statisticsHandler)
	router.POST("/generate", h.generateHandler)

	// appApi := router.Group("api")
	// {
	// 	appApi.GET("/statistics", api.GetStatistics)
	// 	appApi.POST("/generate", api.Generate)
	// }

	return router
}
