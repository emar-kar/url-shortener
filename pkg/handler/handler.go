package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emar-kar/url-shortener/api"
)

type Handler struct{}

func NewHandler() *Handler {
	return nil
}

func (h *Handler) InitRoutes(ginMode string) *gin.Engine {
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.StaticFS("web", http.Dir("web"))
	router.LoadHTMLGlob("web/*/*.html")

	router.GET("", h.mainHandler)
	router.GET("/generator", h.generatorHandler)
	router.GET("/statistic", h.statisticsHandler)

	appApi := router.Group("api")
	{
		appApi.GET("/graph", api.GetStatistics)
		appApi.POST("/generate", api.Generate)
	}

	return router
}

func (h *Handler) mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
