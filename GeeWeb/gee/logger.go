package gee

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		// 记录开始时间
		n := time.Now()
		// 使用next函数使指定part按照顺序运行
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(n))
	}
}
