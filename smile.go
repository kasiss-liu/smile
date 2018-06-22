//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

//一个服务器引擎
type Engine struct {
	RouteGroup    *RouteGroup
	fileEngine    IEngine
	dynamicEngine IEngine
	wsEngine      IEngine
	RunHandle     RunMonitor
	Logger        ILogger
	Gzip          bool
	Rout404       HandlerFunc //注册时 404调用
	//debug
}

//生成一个默认配置的服务器
//有动态引擎和websocket引擎
func Default() *Engine {
	return &Engine{
		dynamicEngine: &DynamicEngine{},
		wsEngine:      &WsEngine{},
		Logger:        &Logger{os.Stdout, true},
		Gzip:          true,
		RouteGroup:    new(RouteGroup),
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
	//初始化一个请求复合 包含了本次请求及响应的数据
	combine := InitCombination(w, r, e)

	//初始化使用引擎
	engine := e.initActEngine(combine)

	var actType string = "NOT_FOUND"

	if engine != nil {
		actType = engine.GetType()
	}

	//如果监控开关打开 并且注册了方法
	//则在请求业务方法前后调用注册的start 和end 方法
	if monitorSwitch && e.RunHandle != nil {
		e.RunHandle.HandleStart(&MonitorInfo{
			time.Now(),
			actType,
			combine.GetPath(),
			combine,
		})
	}

	var err error

	if engine != nil {

		err = engine.Handle()

	} else {
		//当请求的路由不在注册列表中时
		//如果注册了Route404修复方法 则调用Route404
		if e.Rout404 != nil {
			err = e.Rout404(combine)
		}
	}

	if err != nil {
		//debug
		fmt.Println("[debug] " + err.Error())
	}

	//监控结束方法
	if monitorSwitch && e.RunHandle != nil {
		e.RunHandle.HandleEnd(&MonitorInfo{
			time.Now(),
			actType,
			combine.GetPath(),
			combine,
		})
	}

	//如果已经注册了 并且日志开关开启
	//则进行日志打印
	if e.Logger != nil && logSwitch {
		e.Logger.Log(combine)
	}
	combine.Close()

}

//匹配本次请求的处理引擎
func (e *Engine) initActEngine(c *Combination) (threadEngine IEngine) {

	if e.fileEngine != nil {
		threadEngine = e.fileEngine.Init(c)
		if threadEngine.Check(nil) {
			return threadEngine
		}
	}
	if e.wsEngine != nil {
		threadEngine = e.wsEngine.Init(c)
		if threadEngine.Check(e.RouteGroup) {
			return
		}
	}

	if e.dynamicEngine != nil {
		threadEngine = e.dynamicEngine.Init(c)
		if threadEngine.Check(e.RouteGroup) {
			return
		}
	}
	return nil
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

//关闭Gzip
func (e *Engine) GzipOff() {
	e.Gzip = false
}

//注册404回调方法
func (e *Engine) SetRout404(fn HandlerFunc) {
	e.Rout404 = fn
}

//启动一个HttpServer
func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

//启动一个HttpsServer
func (e *Engine) RunTLS(port, cert, key string) {
	http.ListenAndServeTLS(port, cert, key, e)
}
