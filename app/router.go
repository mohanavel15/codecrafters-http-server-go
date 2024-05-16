package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type HanderFunc func(req Request, res *Response)

type Router struct {
	Routes map[string]HanderFunc
	mx     sync.RWMutex
}

func NewRouter() Router {
	return Router{
		Routes: map[string]HanderFunc{},
	}
}

func (r *Router) AddRoute(method, path string, hander HanderFunc) {
	r.mx.Lock()
	path = strings.ReplaceAll(path, "{any}", ".+")
	path = fmt.Sprintf("^%s$", path)
	path = fmt.Sprintf("%s\x00%s", method, path)
	r.Routes[path] = hander
	r.mx.Unlock()
}

func (r *Router) Handle(req Request, res *Response) {
	r.mx.RLock()
	for key, handlerf := range r.Routes {
		skey := strings.Split(key, "\x00")
		match, err := regexp.MatchString(skey[1], req.Path)
		if err != nil {
			continue
		}

		if match && skey[0] == req.Method {
			res.StatusCode = 200
			res.StatusStr = "OK"
			handlerf(req, res)
		}
	}
	r.mx.RUnlock()
}
