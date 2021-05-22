package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) generatorHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "generator.html", nil)
}
