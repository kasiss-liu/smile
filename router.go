package smile

import (
	"errors"
)

type HandlerFunc func(*Combination) error

const (
	METHOD_GET  = "GET"
	METHOD_POST = "POST"
	METHOD_WS   = "WS"
)

type RouteGroup struct {
	GET  map[string]HandlerFunc
	POST map[string]HandlerFunc
	WS   map[string]HandlerFunc
}

func (rg *RouteGroup) Set(method string, path string, handler HandlerFunc) {
	switch method {
	case METHOD_GET:
		rg.GET[path] = handler
	case METHOD_POST:
		rg.POST[path] = handler
	case METHOD_WS:
		rg.WS[path] = handler
	default:
	}
}

func (rg *RouteGroup) SetGET(path string, handler HandlerFunc) {
	rg.Set(METHOD_GET, path, handler)
}
func (rg *RouteGroup) SetPOST(path string, handler HandlerFunc) {
	rg.Set(METHOD_POST, path, handler)
}
func (rg *RouteGroup) SetWS(path string, handler HandlerFunc) {
	rg.Set(METHOD_WS, path, handler)
}

func (rg *RouteGroup) Get(method string, path string) (HandlerFunc, error) {
	switch method {
	case METHOD_GET:
		if val, ok := rg.GET[path]; ok {
			return val, nil
		}
	case METHOD_POST:
		if val, ok := rg.POST[path]; ok {
			return val, nil
		}
	case METHOD_WS:
		if val, ok := rg.WS[path]; ok {
			return val, nil
		}
	default:
	}
	return nil, errors.New("METHOD:" + method + " PATH:" + path + " DID NOT REGISTER YET")
}

func NewRouteGroup() *RouteGroup {
	r := &RouteGroup{}
	r.GET = make(map[string]HandlerFunc, 10)
	r.POST = make(map[string]HandlerFunc, 10)
	r.WS = make(map[string]HandlerFunc, 10)
	return r
}
