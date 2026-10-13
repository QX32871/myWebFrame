package frame

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc define the request handler used by cos-frame
type HandlerFunc func(c *Context)

type H map[string]interface{}

// Engine impl the  interface of ServerHttp
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 顶层Engine保存所有的分组
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件支持
	parent      *RouterGroup  // 支持子分组
	engine      *Engine       // 所有的分组都复用同一个路由引擎实例
}

// New constructor of frame.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由方法，兼容使用分组添加路由和普通方式添加路由两种添加方法
func (group *RouterGroup) addRouter(method string, uriStr string, handler HandlerFunc) {
	pattern := group.prefix + uriStr
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET func to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

// POST func to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

// Run func to start http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(writer, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
