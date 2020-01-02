//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"net/http"
	"os"
	"path"
	"strings"
)

//IEngine 一个引擎接口
type IEngine interface {
	Init(*Combination) IEngine //初始化引擎
	Handle() error             //执行方法
	Check(interface{}) bool    //判断是否属于引擎处理请求
	GetType() string           //获取引擎结构类型
}

//DefaultFile 文件服务器默认输出文件
var DefaultFile = "index.html"

//httpExceptFile http包内的文件服务函数 会针对index.html做301
const httpExceptFile = "index.html"

//ctxEngine 请求处理引擎
//业务处理的相关处理引擎
type ctxEngine struct {
	enableFile    bool         //是否支持静态文件
	baseDir       string       //文件仓库地址
	indexFilename string       //默认文件
	method        string       //请求方法
	path          string       //请求地址
	cb            *Combination //http复合体
}

func createEngine(eFile bool, config ...string) *ctxEngine {
	e := &ctxEngine{}
	e.enableFile = eFile
	if eFile {
		if len(config) > 1 {
			e.baseDir = config[0]
			if len(config) == 1 {
				e.indexFilename = DefaultFile
			} else {
				e.indexFilename = config[1]
			}
			//判断路径是否可用
			fileInfo, err := os.Stat(e.baseDir)
			if err != nil {
				panic(err)
			}
			//判断文件路径是否是一个文件夹
			if !fileInfo.IsDir() {
				panic(e.baseDir + " is not a directory")
			}
		} else {
			e.enableFile = false
		}
	}
	return e
}

//Init 引擎初始化
//获取请求类型和请求路由 并保存
func (e *ctxEngine) Init(c *Combination) IEngine {
	method := c.GetMethod()
	path := "/" + strings.Trim(c.GetPath(), "/")

	return &ctxEngine{
		baseDir:    e.baseDir,
		enableFile: e.enableFile,
		cb:         c,
		method:     method,
		path:       path,
	}

}

//Check 在路由列表中判断 动态请求路由是否已经注册
//如果已经注册 则本次请求由动态引擎处理
//保存路由中已经注册的业务方法
func (e *ctxEngine) Check(i interface{}) bool {
	var fn HandlerFunc 
	rtg := i.(*RouteGroup)
	//先进行文件判断
	if e.enableFile {
		filename := strings.Trim(e.cb.GetPath(), "/")
		//http包内的文件服务函数 会针对index.html做301
		//在这里做一下特殊处理
		if filename == httpExceptFile {
			e.cb.Request.URL.Path = "/"
		}
		//默认index.html 如果直接访问根目录 则返回index.html页面
		if filename == "" {
			filename = e.indexFilename
		}
		filePath := path.Clean(e.baseDir + filename)
		//获取文件/文件夹信息
		file, err := os.Stat(filePath)
		//如果存在 则判断是否是文件夹
		if err == nil && !file.IsDir() {
			e.path = filePath
			fn = e.serveFile
		}
	}
	if rtg != nil {
		if fn == nil  {
			method := ""
			if e.cb.Request.Header.Get("Upgrade") == "websocket" {
				method = MethodWs
			} else {
				method = e.method
			}
			Handler, err := rtg.Get(method, e.path)
			if err == nil {
				fn = Handler
			}
		}
		if fn == nil {
			fn = rtg.route404
		}
		//加载所有路由
		for _,f := range rtg.routeMiddleware {
			e.cb.handlerChain.add(f)
		}
	}
	e.cb.handlerChain.add(fn)
	return true
}

func (e *ctxEngine) serveFile(c *Combination) error {
	//调用http包输出文件的方法
	http.ServeFile(c.ResponseWriter, c.Request, e.path)
	return nil
}

//Handle 执行已经保存的业务方法
//暂时不做错误返回处理
func (e *ctxEngine) Handle() (err error) {
	defer doRecover(&err, e.cb)
	err = e.cb.handlerChain.next(e.cb)
	return
}

//GetType 获取引擎结构类型
func (e *ctxEngine) GetType() string {
	return "ContextEngine"
}
