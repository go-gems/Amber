package handler

import (
	"amber/src/storage"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
)

func New(schema string, host string, storage storage.Service) *gin.Engine {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.POST("/links", storeLink(schema, host, storage))
	router.GET("/r/:shortLink", redirect(storage))

	adminPwd := os.Getenv("ADMIN_PASSWORD")
	if adminPwd == "" {
		fmt.Println("WARNING: ADMIN_PASSWORD environment variable is not set. Admin route will be inaccessible or insecure.")
	}

	admin := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": adminPwd,
	}))
	admin.GET("/", func(c *gin.Context) {
		c.File("./public/admin.html")
	})
	admin.GET("/links/all", listAll(schema, host, storage))
	admin.GET("/links/stats/:shortLink", getStatistics(storage))
	admin.DELETE("/delete/:shortLink", deleteLink(storage))

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

func deleteLink(storage storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortLink, _ := c.Params.Get("shortLink")
		err := storage.Delete(shortLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
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
			Path:   "/r/" + id}

		fmt.Printf("Generated link: %v \n", u.String())
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"url":    u.String(),
		})
	}
}
