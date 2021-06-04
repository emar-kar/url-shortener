package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error string `json:"error"`
}

func errorResponse(c *gin.Context, status int, err error) {
	log.Println(err)
	c.AbortWithStatusJSON(status, Response{Error: err.Error()})
}
