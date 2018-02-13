//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"errors"
)

//定一个业务执行方法
type HandlerFunc func(*Combination) error

//定义部分请求类型
const (
	METHOD_GET  = "GET"
	METHOD_POST = "POST"
	METHOD_WS   = "WS"
)

//路由列表
type RouteGroup struct {
	GET  map[string]HandlerFunc
	POST map[string]HandlerFunc
	WS   map[string]HandlerFunc
}

//注册一个路由
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

//注册一个GET方法请求到的路由
func (rg *RouteGroup) SetGET(path string, handler HandlerFunc) {
	rg.Set(METHOD_GET, path, handler)
}

//注册一个POST方法可用的路由
func (rg *RouteGroup) SetPOST(path string, handler HandlerFunc) {
	rg.Set(METHOD_POST, path, handler)
}

//注册一个websocket路由
func (rg *RouteGroup) SetWS(path string, handler HandlerFunc) {
	rg.Set(METHOD_WS, path, handler)
}

//根据请求方法 获取一个注册的方法
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

//生成一个新的路由列表
func NewRouteGroup() *RouteGroup {
	r := &RouteGroup{}
	r.GET = make(map[string]HandlerFunc, 10)
	r.POST = make(map[string]HandlerFunc, 10)
	r.WS = make(map[string]HandlerFunc, 10)
	return r
}
