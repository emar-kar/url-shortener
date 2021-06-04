package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emar-kar/urlshortener"
	"github.com/gin-gonic/gin"
)

const (
	TimeLayout = "2006-01-02"
)

func (h *Handler) mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main",
	})
}

func (h *Handler) notFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"title": "404 error",
	})
}

func (h *Handler) generateHandler(c *gin.Context) {
	url := c.PostForm("userLink")
	if url == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Main",
			"warning": "URL was not set! ",
		})
		return
	}
	exp := c.PostForm("expirationDate")

	var dur time.Duration
	var err error
	if exp == "" {
		dur, err = time.ParseDuration("24h")
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{
				"title": "500 error",
			})
			return
		}
	} else {
		parsedTime, err := time.Parse(TimeLayout, exp)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{
				"title": "500 error",
			})
			return
		}
		dur = time.Until(parsedTime)
	}

	if dur < 0 {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":   "Main",
			"warning": "Expiration time should be in the future! ",
		})
		return
	}

	shortURL, err := h.services.GenerateShortURL(c.Request.Host)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"title": "500 error",
		})
		return
	}

	link := &urlshortener.Link{FullForm: url, ShortForm: shortURL, Expiration: dur}

	if err := h.services.Set(link); err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"title": "500 error",
		})
		return
	}

	c.HTML(http.StatusCreated, "generator.html", gin.H{
		"title":      "Generated link",
		"longLink":   url,
		"shortLink":  shortURL,
		"expiration": dur.String(),
	})
}

func (h *Handler) statisticsHandler(c *gin.Context) {
	url := c.Request.FormValue("userLink")
	if url == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":   "Main",
			"warning": "URL was not set! ",
		})
		return
	}
	data, err := h.services.Get(url)
	if err != nil {
		log.Println(err)
		h.notFound(c)
		return
	}
	c.HTML(http.StatusOK, "statistics.html", gin.H{
		"title":      "Link statistics",
		"longLink":   data.FullForm,
		"shortLink":  data.ShortForm,
		"expiration": data.Expiration.String(),
		"redirects":  data.Redirects,
	})
}

func (h *Handler) redirectHandler(c *gin.Context) {
	url := fmt.Sprint(c.Request.Host + c.Request.URL.Path)
	data, err := h.services.Get(url)
	if err != nil {
		log.Println(err)
		h.notFound(c)
		return
	}

	if err := h.services.Redirect(data.FullForm); err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"title": "500 error",
		})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, data.FullForm)
	c.Abort()
}
