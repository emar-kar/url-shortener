package handler

import (
	"errors"
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
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if statRequest.Link == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("url is empty"))
		return
	}

	data, err := h.services.Get(statRequest.Link)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, data)
}
