//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"net/http"
	"os"
)

//Engine 一个服务器引擎
type Engine struct {
	RouteGroup 		*RouteGroup
	engine     		IEngine
	Logger     		ILogger
	Gzip       		bool
	//debug
	Errors 			[]error
}

//Default 生成一个默认配置的服务器
//有动态引擎和websocket引擎
func Default() *Engine {
	return &Engine{
		engine:     createEngine(false),
		Logger:     &Logger{os.Stdout, true},
		Gzip:       true,
		RouteGroup: new(RouteGroup),
	}
}

//NewEngine 获取一个具有全部处理引擎的服务器
func NewEngine(config ...string) *Engine {
	e := Default()
	if len(config) > 1 {
		e.engine = createEngine(true, config...)
	}
	return e
}

//ServeHTTP 有请求的时候 把请求处理以后 储存到结构中
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//初始化一个请求复合 包含了本次请求及响应的数据
	cb := initContext(w, r, e)
	defer cb.Close()

	//初始化使用引擎
	engine := e.engine.Init(cb)

	var err error
	//路由校验以及404处理
	engine.Check(e.RouteGroup)

	//如果已经注册了 并且日志开关开启
	//则进行日志打印
	if e.Logger != nil && logSwitch {
		cb.handlerChain.add(e.Logger.Log)
	}
	err = engine.Handle()
	
	if err != nil {
		//debug
		doDebug(err, cb)
	}

}

//SetLoger 注册一个logger
func (e *Engine) SetLoger(l ILogger) {
	e.Logger = l
}

//SetRouteGroup 注册一个路由组
func (e *Engine) SetRouteGroup(r *RouteGroup) {
	e.RouteGroup = r
}

//GzipOn 开启Gzip
func (e *Engine) GzipOn() {
	e.Gzip = true
}

//GzipOff 关闭Gzip
func (e *Engine) GzipOff() {
	e.Gzip = false
}

func (e *Engine) prepareRun() {
	if !GetInitState() {
		DoCustomInit()
	}
	doPrintRoutes(e.RouteGroup.FormatRoutes())
}

//Run 启动一个HttpServer
func (e *Engine) Run(port string) (err error) {

	defer doRecover(&err, nil)

	e.prepareRun()
	
	err = http.ListenAndServe(port, e)
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
	return
}

//RunTLS 启动一个HttpsServer
func (e *Engine) RunTLS(port, cert, key string) (err error) {

	defer doRecover(&err, nil)

	e.prepareRun()

	err = http.ListenAndServeTLS(port, cert, key, e)
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
	return
}

//GetErrors 获取引擎中的错误
func (e *Engine) GetErrors() []error {
	return e.Errors
}
