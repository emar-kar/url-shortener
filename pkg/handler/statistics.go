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
	var statRequest StatRequest
	if err := c.BindJSON(&statRequest); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if statRequest.Link == "" {
		ErrorResponse(c, http.StatusBadRequest, errors.New("url is empty"))
		return
	}

	data, err := h.services.GetLink(statRequest.Link)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, data)
}
