package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"youOrHell/src/storage"
)

func New(schema string, host string, storage storage.Service) *gin.Engine {
	router := gin.Default()
	router.POST("/links/", storeLink(schema, host, storage))
	router.GET("/links/all", listAll(schema, host, storage))

	router.GET("/links/stats/:shortLink", getStatistics(storage))
	router.GET("/r/:shortLink", redirect(storage))

	return router
}

func listAll(schema string, host string, storage storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		lst := storage.GetAll()

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"links":  lst,
		})

	}
}

func serveFile(c *gin.Context) {

}

func redirect(storage storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortLink, _ := c.Params.Get("shortLink")

		url, err := storage.Load(shortLink)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
			return
		}
		c.Redirect(301, url)

	}
}

func getStatistics(storage storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortLink, _ := c.Params.Get("shortLink")
		statistics, err := storage.LoadInfo(shortLink)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statistics": statistics,
		})

	}
}

func storeLink(schema string, host string, storage storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, _ := c.GetRawData()
		id, err := storage.Save(string(content))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err,
			})
			return

		}
		u := url.URL{
			Scheme: schema,
			Host:   host,
			Path:   id}

		fmt.Printf("Generated link: %v \n", u.String())
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"url":    u.String(),
		})
	}
}
