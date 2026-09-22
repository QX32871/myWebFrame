package frame

import (
	"fmt"
	"net/http"
)

// HandlerFunc define the request handler used by cos-frame
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine impl the  interface of ServerHttp
type Engine struct {
	router map[string]HandlerFunc
}

// New constructor of frame.Engine
func New() *Engine {
	// 创建Engine之后返回它的值
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由方法
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET func to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

// POST func to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// Run func to start http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(writer, req)
	} else {
		fmt.Fprintf(writer, "404 NOT FOUND: %s\n", req.URL)
	}
}
