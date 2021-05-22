package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) statisticsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "statistics.html", nil)
}
