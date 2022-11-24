package gee

import (
	"net/http"
)

// HandleFunc 定义了请求的处理方式
type HandleFunc func(c *Context)

// Engine 实现了serverHttp 接口
type Engine struct {
	router *router
}

// New 是Engine的构造方法
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 开头小写，不导出该方法
//func (engine *Engine) addRouter(method string, pattern string, handler HandleFunc) {
//	key := method + "-" + pattern
//	engine.router[key] = handler
//}

// GET 定义了添加get请求的方法
func (engine *Engine) GET(pattern string, handle HandleFunc) {
	engine.router.addRoute("GET", pattern, handle)
}

func (engine *Engine) POST(pattern string, handle HandleFunc) {
	engine.router.addRoute("POST", pattern, handle)
}

// Run 开启了httpserver服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 方法，实现该方法使得Engine结构体可以作为Handle类，启动服务器并放入Handle实例，
// 可以在指定拦截所有http请求并处理指定请求，处理方式就是ServeHTTP方法
// 这是框架的内容，框架使用者只需要写好处理方法HandleFunc放入engine，剩下的事情框架来做
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	engine.router.handle(c)
}
