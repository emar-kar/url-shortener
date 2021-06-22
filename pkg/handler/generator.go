package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/emar-kar/urlshortener"
)

type GenRequest struct {
	Link    string `json:"link"`
	ExpTime string `json:"expiration_time"`
}

func (gr *GenRequest) Time() (time.Duration, error) {
	if gr.ExpTime == "" {
		duration, err := time.ParseDuration("24h")
		if err != nil {
			return 0, err
		}
		return duration, nil
	}
	parsedTime, err := time.Parse(TimeLayout, gr.ExpTime)
	if err != nil {
		return 0, err
	}
	return time.Until(parsedTime), nil
}

func (h *Handler) generate(c *gin.Context) {
	var genRequest GenRequest
	if err := c.BindJSON(&genRequest); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if genRequest.Link == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("url is empty"))
		return
	}

	shortURL, err := h.services.GenerateShortURL(c.Request.Host)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	dur, err := genRequest.Time()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if dur < 0 {
		errorResponse(c, http.StatusBadRequest, errors.New("expiration time is in the past"))
		return
	}

	link := &urlshortener.Link{
		FullForm:   genRequest.Link,
		ShortForm:  shortURL,
		Expiration: dur,
	}

	if err := h.services.Set(link); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, link)
}
