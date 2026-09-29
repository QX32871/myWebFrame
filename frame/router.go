package frame

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots获取路由地址	handlers获取路由对应的handler
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)

			// 如果为通配符，就跳过循环
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 重写添加路由逻辑，从map改为使用trie树
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 先解析路由
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 如果获取不到请求方法（GET这种）
	_, ok := r.roots[method]

	// 传一个空的node
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchPaths := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 查找path
	n := root.search(searchPaths, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {

			// 处理两个通配符
			if part[0] == ':' {
				params[part[1:]] = searchPaths[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPaths[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 重写handle方法，使用trie树
func (r *router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
