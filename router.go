//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//HandlerFunc 定一个业务执行方法
type HandlerFunc func(*Combination) error

//定义部分请求类型及其匹配式
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodWs     = "WS"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	regexpPost   = "(POST)|(Post)|"
	regexpGet    = "(GET)|(Get)|"
	regexpWs     = "(WS)|(Ws)|"
	regexpPut    = "(PUT)|(Put)|"
	regexpDet    = "(DELETE)|(Delete)|"
)

//定义自动生成路由的风格
const (
	StyleHump    = "hump"
	StyleConnect = "connector"
)

//path的打印基准长度
const pathLen = 10

//RouteGroup 路由列表
type RouteGroup struct {
	GET            map[string]HandlerFunc
	POST           map[string]HandlerFunc
	WS             map[string]HandlerFunc
	PUT            map[string]HandlerFunc
	DELETE         map[string]HandlerFunc
	pathStyle      string                       //自动填充路由时 方法名称转化为路径后的风格
	routeFnameList map[string]map[string]string //路由->方法名列表
}

func (rg *RouteGroup) initRouteFnameList() {
	list := make(map[string]map[string]string, 5)
	list[MethodGet] = make(map[string]string, 10)
	list[MethodPost] = make(map[string]string, 10)
	list[MethodPut] = make(map[string]string, 10)
	list[MethodWs] = make(map[string]string, 10)
	list[MethodDelete] = make(map[string]string, 10)
	rg.routeFnameList = list
}

func (rg *RouteGroup) setRouteFnameList(method, path string, handler interface{}) {
	switch v := handler.(type) {
	case string:
		rg.routeFnameList[method][path] = v
	case HandlerFunc:
		rg.routeFnameList[method][path] = getFuncName(handler)
	}
}

//Set 注册一个路由
func (rg *RouteGroup) Set(method string, path string, handler HandlerFunc) {
	path = trimPath(path)
	set := true
	switch method {
	case MethodGet:
		rg.GET[path] = handler
	case MethodPost:
		rg.POST[path] = handler
	case MethodWs:
		rg.WS[path] = handler
	case MethodPut:
		rg.PUT[path] = handler
	case MethodDelete:
		rg.DELETE[path] = handler
	default:
		set = false
	}
	if set {
		rg.setRouteFnameList(method, path, handler)
	}

}

//SetGET 注册一个GET方法请求到的路由
func (rg *RouteGroup) SetGET(path string, handler HandlerFunc) {
	rg.Set(MethodGet, path, handler)
}

//SetPOST 注册一个POST方法可用的路由
func (rg *RouteGroup) SetPOST(path string, handler HandlerFunc) {
	rg.Set(MethodPost, path, handler)
}

//SetWS 注册一个websocket路由
func (rg *RouteGroup) SetWS(path string, handler HandlerFunc) {
	rg.Set(MethodWs, path, handler)
}

//SetPUT 注册一个PUT方法可用的路由
func (rg *RouteGroup) SetPUT(path string, handler HandlerFunc) {
	rg.Set(MethodPut, path, handler)
}

//SetDEL 注册一个PUT方法可用的路由
func (rg *RouteGroup) SetDEL(path string, handler HandlerFunc) {
	rg.Set(MethodDelete, path, handler)
}

//Get 根据请求方法 获取一个注册的方法
func (rg *RouteGroup) Get(method string, path string) (HandlerFunc, error) {
	switch method {
	case MethodGet:
		if val, ok := rg.GET[path]; ok {
			return val, nil
		}
	case MethodPost:
		if val, ok := rg.POST[path]; ok {
			return val, nil
		}
	case MethodWs:
		if val, ok := rg.WS[path]; ok {
			return val, nil
		}
	case MethodPut:
		if val, ok := rg.PUT[path]; ok {
			return val, nil
		}
	case MethodDelete:
		if val, ok := rg.DELETE[path]; ok {
			return val, nil
		}
	default:
	}
	return nil, errors.New("METHOD:" + method + " PATH:" + path + " DID NOT REGISTER YET")
}

//NewRouteGroup 生成一个新的路由列表
func NewRouteGroup() *RouteGroup {
	r := &RouteGroup{}
	r.GET = make(map[string]HandlerFunc, 10)
	r.POST = make(map[string]HandlerFunc, 10)
	r.WS = make(map[string]HandlerFunc, 10)
	r.PUT = make(map[string]HandlerFunc, 10)
	r.DELETE = make(map[string]HandlerFunc, 10)
	r.SetPathStyleConnector()
	r.initRouteFnameList()
	return r
}

//FillRoutes 填充路由基础方法
func (rg *RouteGroup) FillRoutes(method string, prefix string, c interface{}) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	l := t.NumMethod()
	for i := 0; i < l; i++ {
		fnName := t.Method(i).Name
		interf := v.Method(i).Interface()
		set := true
		if fn, ok := interf.(func(*Combination) error); ok {
			fnName = rg.transFnNameToPath(fnName)
			path := strings.Trim(prefix+"/"+fnName, "/")
			switch method {
			case MethodGet:
				rg.SetGET(path, fn)
			case MethodPost:
				rg.SetPOST(path, fn)
			case MethodWs:
				rg.SetWS(path, fn)
			case MethodPut:
				rg.SetPUT(path, fn)
			case MethodDelete:
				rg.SetDEL(path, fn)
			default:
				set = false
			}
			if set {
				rg.setRouteFnameList(method, trimPath(path), t.String()+"."+t.Method(i).Name)
			}
		}
	}
}

//PrefixFillRoutes 前缀匹配规则 填充路由
//暂时只支持GET、POST、WS
//将一个Controller结构下的方法按照方法名称注册到routeGroup中
func (rg *RouteGroup) PrefixFillRoutes(prefix string, c interface{}) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	l := t.NumMethod()
	reg, _ := regexp.Compile(`^(` + regexpPost + regexpGet + regexpWs + regexpPut + regexpDet + `).+`)
	var method string
	for i := 0; i < l; i++ {
		fnName := t.Method(i).Name
		interf := v.Method(i).Interface()
		rexSubmatch := reg.FindStringSubmatch(fnName)
		if len(rexSubmatch) > 0 {
			method = strings.ToUpper(rexSubmatch[1])
			//去掉函数名称中的方法类型
			if len(rexSubmatch[1]) < 1 {
				fnName = "/"
				method = strings.ToUpper(rexSubmatch[0])
			} else {
				fnName = strings.Replace(fnName, rexSubmatch[1], "", -1)
			}
		}
		if fn, ok := interf.(func(*Combination) error); ok {
			//函数名称转化为请求路径path的全小写格式
			fnName = rg.transFnNameToPath(fnName)
			path := strings.Trim(prefix+"/"+fnName, "/")
			set := true
			switch method {
			case MethodGet:
				rg.SetGET(path, fn)
			case MethodPost:
				rg.SetPOST(path, fn)
			case MethodWs:
				rg.SetWS(path, fn)
			case MethodPut:
				rg.SetPUT(path, fn)
			case MethodDelete:
				rg.SetDEL(path, fn)
			default:
				set = false
			}
			if set {
				rg.setRouteFnameList(method, trimPath(path), t.String()+"."+t.Method(i).Name)
			}

		}
	}
}

//SetPathStyleHump 设置路径风格为驼峰 即区分大小写
func (rg *RouteGroup) SetPathStyleHump() {
	rg.pathStyle = StyleHump
}

//SetPathStyleConnector 设置路径风格为连字符格式 将驼峰转为连接符风格
func (rg *RouteGroup) SetPathStyleConnector() {
	rg.pathStyle = StyleConnect
}

//根据规则转化方法名称为路径
func (rg *RouteGroup) transFnNameToPath(fnName string) string {
	if rg.pathStyle == StyleConnect {
		reg, _ := regexp.Compile(`[A-Z]+`)
		fnName = reg.ReplaceAllStringFunc(fnName, func(str string) string {
			return "-" + strings.ToLower(str)
		})
	}
	return strings.Trim(fnName, "-")
}

//FormatRoutes 返回格式化的路由信息 每个路由信息为一个string
func (rg *RouteGroup) FormatRoutes() []string {
	rm := make(map[string]string, 10)
	rt := make(map[string][]byte)
	baseLen := pathLen
	for method, rgroup := range rg.routeFnameList {
		for path, fnName := range rgroup {
			rt[path] = []byte(method)
			rm[path] = fnName
		}
	}

	//重组每条路由信息数据
	rs := make([]string, 0, 30)
	for path, fnName := range rm {
		s := fmt.Sprintf("%-6s  %-"+strconv.Itoa(baseLen)+"s --> %s",
			rt[path],
			path,
			fnName,
		)

		rs = append(rs, s)
	}
	return rs
}
