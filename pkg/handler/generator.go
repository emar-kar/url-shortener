package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/emar-kar/urlshortener"
)

// GenRequest represents structure of the API request json object for link generation.
type GenRequest struct {
	Link    string `json:"link"`
	ExpTime string `json:"expiration_time"`
}

// Time converts expiration time from the json request
// to golang time.Duration.
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

// generateLink wraps API request handler.
// Generates short URL and returns it as a link json object.
func (h *Handler) generateLink(c *gin.Context) {
	var genRequest GenRequest
	if err := c.BindJSON(&genRequest); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if genRequest.Link == "" {
		ErrorResponse(c, http.StatusBadRequest, errors.New("url is empty"))
		return
	}

	shortURL, err := h.services.GenerateShortURL(c.Request.Host)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	dur, err := genRequest.Time()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if dur < 0 {
		ErrorResponse(c, http.StatusBadRequest, errors.New("expiration time is in the past"))
		return
	}

	link := &urlshortener.Link{
		FullForm:   genRequest.Link,
		ShortForm:  shortURL,
		Expiration: dur,
	}

	if err := h.services.SetLink(link); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, link)
}
