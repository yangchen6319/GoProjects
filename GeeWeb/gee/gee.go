package gee

import (
	"log"
	"net/http"
)

// HandleFunc 定义了请求的处理方式
type HandleFunc func(c *Context)

// RouterGroup 实现路由分组
type RouterGroup struct {
	// 同组路由前缀
	prefix string
	// 中间件
	middlewares []HandleFunc
	// 父路由组
	parent *RouterGroup
	// 所有的路由分组共享同一个engin实例
	engine *Engine
}

// Engine 实现了serverHttp 接口
type Engine struct {
	// 这里存在一个嵌入类型
	// 将RouterGroup作为嵌入类型放入Engin
	// 意味着engin实例可以随意调用RouterGroup的属性
	*RouterGroup
	// 全部采用指针形式意味着engin可以调用任意的RouterGroup实例
	groups []*RouterGroup
	router *router
}

// New 是Engine的构造方法
func New() *Engine {
	engine := &Engine{router: newRouter()}
	// 实例化一个engine, 将engine放入RouterGroup中，这样engine可以访问routergroup，routergroup也可以访问engine
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 开头小写，不导出该方法
//func (engine *Engine) addRouter(method string, pattern string, handler HandleFunc) {
//	key := method + "-" + pattern
//	engine.router[key] = handler
//}

// 创建新的routerGroup并放入engine

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// 所有group共享同一个engine的目的在这
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		// 这里说明，所有的group都是由一个父group来的
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 以下的方法是通过路由分组的方式添加路由信息

func (group *RouterGroup) addRouter(method string, comp string, handle HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handle)
}

func (group *RouterGroup) GET(pattern string, handle HandleFunc) {
	group.addRouter("GET", pattern, handle)
}

func (group *RouterGroup) POST(pattern string, handle HandleFunc) {
	group.addRouter("POST", pattern, handle)
}

// 这里的方式是通过engine直接添加路由信息

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
