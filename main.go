package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// set the router as the default one shipped with Gin
	router := gin.Default()

	// serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views/js", true)))

	// setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// start and run the server
	router.Run(":3000")
}
