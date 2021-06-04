package handler

import (
	"log"
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
	} else {
		parsedTime, err := time.Parse(TimeLayout, gr.ExpTime)
		if err != nil {
			return 0, err
		}
		return time.Until(parsedTime), nil
	}
}

func (h *Handler) generate(c *gin.Context) {
	var genRequest GenRequest
	if err := c.BindJSON(&genRequest); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if genRequest.Link == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	shortURL, err := h.services.GenerateShortURL(c.Request.Host)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dur, err := genRequest.Time()
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	link := &urlshortener.Link{
		FullForm:   genRequest.Link,
		ShortForm:  shortURL,
		Expiration: dur,
	}

	if err := h.services.Set(link); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, link)
}
