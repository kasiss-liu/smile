//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"time"
)

//注册几种请求类型
const (
	ACT_TYPE_FILE = "FILE"
	ACT_TYPE_WS   = "WS"
	ACT_TYPE_GET  = "GET"
	ACT_TYPE_POST = "POST"
)

//一个服务器引擎
type Engine struct {
	RouteGroup    *RouteGroup
	actEngine     IEngine
	actType       string
	fileEngine    IEngine
	dynamicEngine IEngine
	wsEngine      IEngine
	RunHandle     RunMonitor
	Logger        ILogger
	Gzip          bool
	//debug
}

//生成一个默认配置的服务器
//有动态引擎和websocket引擎
func Default() *Engine {
	return &Engine{
		dynamicEngine: &DynamicEngine{},
		wsEngine:      &WsEngine{},
		Logger:        &Logger{os.Stdout, true},
		Gzip:          false,
	}
}

//获取一个具有全部处理引擎的服务器
func NewEngine(fileDir string) *Engine {
	e := Default()
	//判断路径是否可用
	fileInfo, err := os.Stat(fileDir)
	if err != nil {
		panic(err)
	}
	//判断文件路径是否是一个文件夹
	if !fileInfo.IsDir() {
		panic(fileDir + " is not a directory")
	}
	e.fileEngine = &FileEngine{BaseDir: fileDir}
	return e
}

//有请求的时候 把请求处理以后 储存到结构中
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//默认生成一个*gzip.Writer
	//请求结束后关闭
	gz := gzip.NewWriter(w)
	defer gz.Close()
	//初始化一个请求复合 包含了本次请求及响应的数据
	combine := InitCombination(w, r, e, gz)
	//初始化使用引擎
	if e.initActEngine(combine) {
		//如果监控开关打开 并且注册了方法
		//则在请求业务方法前后调用注册的start 和end 方法
		if monitorSwitch && e.RunHandle != nil {
			e.RunHandle.HandleStart(&MonitorInfo{time.Now(), e.actType, combine.GetPath(), combine})
		}
		//业务处理方法调用
		err := e.actEngine.Handle()
		if err != nil {
			//debug
			fmt.Println("[debug] " + err.Error())
		}
		//监控结束方法
		if monitorSwitch && e.RunHandle != nil {
			e.RunHandle.HandleEnd(&MonitorInfo{time.Now(), e.actType, combine.GetPath(), combine})
		}
	} else {
		//未初始化到路由的请求 写入一个404头
		//debug
		combine.WriteHeader(404)

	}
	//如果已经注册了 并且日志开关开启
	//则进行日志打印
	if e.Logger != nil && logSwitch {
		e.Logger.Log(combine)
	}
}

//匹配本次请求的处理引擎
func (e *Engine) initActEngine(c *Combination) bool {

	if e.fileEngine != nil {
		e.fileEngine.Init(c)
		if e.fileEngine.Check(nil) {
			e.actEngine = e.fileEngine
			e.actType = ACT_TYPE_FILE
			return true
		}
		e.fileEngine.Reset()
	}
	if e.wsEngine != nil {
		e.wsEngine.Init(c)
		if e.wsEngine.Check(e.RouteGroup) {
			e.actEngine = e.wsEngine
			e.actType = ACT_TYPE_WS
			return true
		}
		e.wsEngine.Reset()
	}

	if e.dynamicEngine != nil {
		e.dynamicEngine.Init(c)
		if e.dynamicEngine.Check(e.RouteGroup) {
			e.actEngine = e.dynamicEngine
			e.actType = c.GetMethod()
			return true
		}
		e.dynamicEngine.Reset()
	}
	return false
}

//注册一个监控器
func (e *Engine) SetMonitor(m RunMonitor) {
	e.RunHandle = m
}

//注册一个logger
func (e *Engine) SetLoger(l ILogger) {
	e.Logger = l
}

//注册一个路由组
func (e *Engine) SetRouteGroup(r *RouteGroup) {
	e.RouteGroup = r
}

//开启Gzip
func (e *Engine) GzipOn() {
	e.Gzip = true
}

//启动一个HttpServer
func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

//启动一个HttpsServer
func (e *Engine) RunTLS(port, cert, key string) {
	http.ListenAndServeTLS(port, cert, key, e)
}
