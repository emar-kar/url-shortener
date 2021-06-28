package handler

import (
	"github.com/gin-gonic/gin"
)

// Response represents error message response on API calls.
type Response struct {
	Error string `json:"error"`
}

// ErrorResponse simplifies error with a message delivery.
func ErrorResponse(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, Response{Error: err.Error()})
}
