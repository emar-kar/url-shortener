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
	timeLayout = "2006-01-02"
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
	if exp == "" {
		dur, _ = time.ParseDuration("24h")
	} else {
		parsedTime, _ := time.Parse(timeLayout, exp)
		dur = time.Until(parsedTime)
	}

	shortURL, err := h.services.GenerateShortURL(c.Request.Host)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"title": "500 error",
		})
	}

	link := &urlshortener.Link{FullForm: url, ShortForm: shortURL, Expiration: dur}

	if err := h.services.Set(link); err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"title": "500 error",
		})
	}

	c.HTML(http.StatusOK, "generator.html", gin.H{
		"title":      "Generated link",
		"longLink":   url,
		"shortLink":  shortURL,
		"expiration": dur.String(),
	})
}

func (h *Handler) statisticsHandler(c *gin.Context) {
	url := c.Request.FormValue("userLink")
	if url == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
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

	// TODO: fix move then 1 redirect
	c.Redirect(303, data.FullForm)
	// c.Abort()
}
