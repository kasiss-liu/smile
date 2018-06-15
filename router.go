//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

//定一个业务执行方法
type HandlerFunc func(*Combination) error

//定义部分请求类型及其匹配式
const (
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_WS     = "WS"
	METHOD_PUT    = "PUT"
	METHOD_DELETE = "DELETE"
	regexp_Post   = "(POST)|(Post)|"
	regexp_Get    = "(GET)|(Get)|"
	regexp_Ws     = "(WS)|(Ws)|"
	regexp_Put    = "(PUT)|(Put)|"
	regexp_Det    = "(DELETE)|(Delete)|"
)

const (
	STYLE_HUMP    = "hump"
	STYLE_CONNECT = "connector"
)

//路由列表
type RouteGroup struct {
	GET       map[string]HandlerFunc
	POST      map[string]HandlerFunc
	WS        map[string]HandlerFunc
	PUT       map[string]HandlerFunc
	DELETE    map[string]HandlerFunc
	pathStyle string //自动填充路由时 方法名称转化为路径后的风格
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
	case METHOD_PUT:
		rg.PUT[path] = handler
	case METHOD_DELETE:
		rg.DELETE[path] = handler
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

func (rg *RouteGroup) SetPUT(path string, handler HandlerFunc) {
	rg.Set(METHOD_PUT, path, handler)
}

func (rg *RouteGroup) SetDEL(path string, handler HandlerFunc) {
	rg.Set(METHOD_DELETE, path, handler)
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
	case METHOD_PUT:
		if val, ok := rg.PUT[path]; ok {
			return val, nil
		}
	case METHOD_DELETE:
		if val, ok := rg.DELETE[path]; ok {
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
	r.PUT = make(map[string]HandlerFunc, 10)
	r.DELETE = make(map[string]HandlerFunc, 10)
	r.SetPathStyleConnector()
	return r
}

//填充路由基础方法
func (rg *RouteGroup) FillRoutes(method string, prefix string, c interface{}) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	l := t.NumMethod()
	for i := 0; i < l; i++ {
		fnName := t.Method(i).Name
		interf := v.Method(i).Interface()
		if fn, ok := interf.(func(*Combination) error); ok {
			fnName = rg.transFnNameToPath(fnName)
			path := strings.Trim(prefix+"/"+fnName, "/")
			switch method {
			case METHOD_GET:
				rg.SetGET(path, fn)
			case METHOD_POST:
				rg.SetPOST(path, fn)
			case METHOD_WS:
				rg.SetWS(path, fn)
			case METHOD_PUT:
				rg.SetPUT(path, fn)
			case METHOD_DELETE:
				rg.SetDEL(path, fn)
			default:
			}
		}
	}
}

//前缀匹配规则 填充路由
//暂时只支持GET、POST、WS
//将一个Controller结构下的方法按照方法名称注册到routeGroup中
func (rg *RouteGroup) PrefixFillRoutes(prefix string, c interface{}) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	l := t.NumMethod()
	reg, _ := regexp.Compile(`^(` + regexp_Post + regexp_Get + regexp_Ws + regexp_Put + regexp_Det + `).+`)
	var method string
	for i := 0; i < l; i++ {
		fnName := t.Method(i).Name
		interf := v.Method(i).Interface()
		rexSubmatch := reg.FindStringSubmatch(fnName)
		if len(rexSubmatch) > 0 {
			method = strings.ToUpper(rexSubmatch[1])
			//去掉函数名称中的方法类型
			fnName = strings.Replace(fnName, rexSubmatch[1], "", -1)
		}
		if fn, ok := interf.(func(*Combination) error); ok {
			//函数名称转化为请求路径path的全小写格式
			fnName = rg.transFnNameToPath(fnName)
			path := strings.Trim(prefix+"/"+fnName, "/")
			switch method {
			case METHOD_GET:
				rg.SetGET(path, fn)
			case METHOD_POST:
				rg.SetPOST(path, fn)
			case METHOD_WS:
				rg.SetWS(path, fn)
			case METHOD_PUT:
				rg.SetPUT(path, fn)
			case METHOD_DELETE:
				rg.SetDEL(path, fn)
			default:
			}
		}
	}
}

//设置路径风格为驼峰 即区分大小写
func (rg *RouteGroup) SetPathStyleHump() {
	rg.pathStyle = STYLE_HUMP
}

//设置路径风格为连字符格式 将驼峰转为连接符风格
func (rg *RouteGroup) SetPathStyleConnector() {
	rg.pathStyle = STYLE_CONNECT
}

//根据规则转化方法名称为路径
func (rg *RouteGroup) transFnNameToPath(fnName string) string {
	if rg.pathStyle == STYLE_CONNECT {
		reg, _ := regexp.Compile(`[A-Z]+`)
		fnName = reg.ReplaceAllStringFunc(fnName, func(str string) string {
			return "-" + strings.ToLower(str)
		})
	}
	return strings.Trim(fnName, "-")
}
