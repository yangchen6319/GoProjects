package main

import "github.com/gin-gonic/gin"

func Gin() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"jack": "danny",
			"rose": "mary",
		})
	})
	r.Run(":8080")
}
