package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error string `json:"error"`
}

func ErrorResponse(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, Response{Error: err.Error()})
}
