package saber

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	//方法到trie树的映射
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{roots: make(map[string]*node), handlers: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	_, exist := r.roots[method]
	if !exist {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, exist := r.roots[method]
	if !exist {
		return nil, nil
	}
	for _, res := range searchParts {
		fmt.Println(res)
	}
	resNode := root.search(searchParts, 0)
	fmt.Println(resNode)
	if resNode != nil {
		parts := parsePattern(resNode.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i][1:]
			}
			if part[0] == '*' && len(part) > 1 {
				//将字符串切片连接成一个单独的字符串，之间用第二个参数分隔。
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return resNode, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	fmt.Println(c.Path)
	resNode, params := r.getRoute(c.Method, c.Path)
	fmt.Println(resNode)
	if resNode != nil {
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[c.Method+"-"+resNode.pattern])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
