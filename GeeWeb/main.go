package main

import (
	"fmt"
	"gee"
	"html/template"
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

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	// 实例化一个engine
	e := gee.New()
	// 设置模板
	e.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./static")
	// 加个全局中间件
	e.AddMiddleware(gee.Logger())

	stu1 := &student{"zhang", 23}
	stu2 := &student{"wang", 16}

	v1 := e.Group("/v1")
	{

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
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	e.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	e.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	e.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	e.Run(":8080")

}
