package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandleFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	// 实例化一个engine
	e := gee.New()
	// 加个全局中间件
	e.AddMiddleware(gee.Logger())
	v1 := e.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := e.Group("/v2")
	v2.AddMiddleware(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	e.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
	})
	e.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	e.Run(":8080")

}
