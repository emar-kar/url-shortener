package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatRequest struct {
	Link string `json:"link"`
}

func (h *Handler) getStatistics(c *gin.Context) {
	log.Println(c.Request)
	var statRequest StatRequest
	if err := c.BindJSON(&statRequest); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if statRequest.Link == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := h.services.Get(statRequest.Link)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, data)
}
