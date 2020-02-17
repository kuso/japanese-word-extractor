package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/gzip"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Serve frontend static files
	//router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	router.Static("/", "./frontend/build")

	/*
	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	 */

	// Start and run the server
	router.Run(":8080")

}
