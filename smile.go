package smile

import (
	"net/http"
	"os"
	"time"
)

/*
logger
router
server
*/
const (
	ACT_TYPE_FILE = "FILE"
	ACT_TYPE_WS   = "WS"
	ACT_TYPE_GET  = "GET"
	ACT_TYPE_POST = "POST"
)

type Engine struct {
	RouteGroup    *RouteGroup
	actEngine     IEngine
	actType       string
	fileEngine    IEngine
	dynamicEngine IEngine
	wsEngine      IEngine
	RunHandle     RunMonitor
	Logger        ILogger
	//logger debug
}

func Default() *Engine {
	return &Engine{
		dynamicEngine: &DynamicEngine{},
		wsEngine:      &WsEngine{},
		Logger:        &Logger{os.Stdout, true},
	}
}
func NewEngine(fileDir string) *Engine {
	e := Default()

	fileInfo, err := os.Stat(fileDir)
	if err != nil {
		panic(err)
	}
	if !fileInfo.IsDir() {
		panic(fileDir + " is not a directory")
	}
	e.fileEngine = &FileEngine{BaseDir: fileDir}
	return e
}

//有请求的时候 把请求处理以后 储存到结构中
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	combine := InitCombination(w, r)
	if e.initactEngine(combine) {
		if monitorSwitch && e.RunHandle != nil {
			e.RunHandle.HandleStart(&MonitorInfo{time.Now(), e.actType, combine.GetPath(), combine})
		}
		err := e.actEngine.Handle()
		if err != nil {
			//debug
		}
		if monitorSwitch && e.RunHandle != nil {
			e.RunHandle.HandleEnd(&MonitorInfo{time.Now(), e.actType, combine.GetPath(), combine})
		}
	} else {
		//debug
		combine.WriteHeader(404)

	}
	if e.Logger != nil {
		e.Logger.Log(combine)
	}
}
func (e *Engine) initactEngine(c *Combination) bool {

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

func (e *Engine) SetMonitor(m RunMonitor) {
	e.RunHandle = m
}
func (e *Engine) SetLoger(l ILogger) {
	e.Logger = l
}
func (e *Engine) SetRouteGroup(r *RouteGroup) {
	e.RouteGroup = r
}

func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

func (e *Engine) RunTLS(port, cert, key string) {
	http.ListenAndServeTLS(port, cert, key, e)
}
