package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 封装每次http请求和处理过程中产生的数据
type Context struct {
	// 这里记住首字母大写供包外访问
	// 原数据
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	// 请求参数
	Params map[string]string
	// 响应信息
	StatusCode int
	// 新增中间件部分
	handlers []HandleFunc
	index    int
	engine   *Engine // 使Context使用Engine中的模板
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

// PostForm 返回post请求中key对应的value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 返回get请求中key对应的value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Context-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Context-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err)
	}
}

func (c *Context) Fail(code int, err error) {
	c.Status(code)
}

// Next 函数用于执行context中的handles函数列表
func (c *Context) Next() {
	c.index++
	size := len(c.handlers)
	// 这里可以配合写好的中间件达到按照想要的顺序执行各部分的目的
	for ; c.index < size; c.index++ {
		c.handlers[c.index](c)
	}
}
