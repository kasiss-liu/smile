//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"net/http"
	"os"
	"path"
	"strings"
)

//注册几种请求类型
const (
	ActTypeFile = "FILE"
	ActTypeWs   = "WS"
	ActTypeDn   = "DYNAMIC"
	ActType404  = "NOT_FOUND"
)

//IEngine 一个引擎接口
type IEngine interface {
	Init(*Combination) IEngine //初始化引擎
	Handle() error             //执行方法
	Check(interface{}) bool    //判断是否属于引擎处理请求
	GetType() string           //获取引擎结构类型
}

var _ IEngine = &FileEngine{}

//FileEngine 文件请求引擎
//响应静态文件数据请求
type FileEngine struct {
	BaseDir  string       //文件仓库地址
	FilePath string       //文件路径
	FileExt  string       //文件后缀
	cb       *Combination //请求复合结构
}

//DefaultFile 文件服务器默认输出文件
var DefaultFile = "index.html"

//httpExceptFile http包内的文件服务函数 会针对index.html做301
const httpExceptFile = "index.html"

//Init 文件引擎初始化
//储存请求复合结构
//处理请求文件地址
//处理请求文件后缀
func (f *FileEngine) Init(c *Combination) IEngine {

	filename := strings.Trim(c.GetPath(), "/")

	//http包内的文件服务函数 会针对index.html做301
	//在这里做一下特殊处理
	if filename == httpExceptFile {
		c.Request.URL.Path = "/"
	}
	//默认index.html 如果直接访问根目录 则返回index.html页面
	if filename == "" {
		filename = DefaultFile
	}
	filePath := path.Clean(f.BaseDir + filename)
	fileExt := path.Ext(filePath)

	return &FileEngine{
		cb:       c,
		FileExt:  fileExt,
		FilePath: filePath,
		BaseDir:  f.BaseDir,
	}
}

//Handle 响应处理方法
func (f *FileEngine) Handle() (err error) {
	//调用http包输出文件的方法
	http.ServeFile(f.cb.ResponseWriter, f.cb.Request, f.FilePath)
	return nil
}

//FileExist 判断请求文件是否存在
func (f *FileEngine) FileExist() bool {

	//获取文件/文件夹信息
	file, err := os.Stat(f.FilePath)
	//如果存在 则判断是否是文件夹
	if err == nil && !file.IsDir() {
		return true
	}
	//针对默认页面做校验

	//不存在或者不是文件 返回false
	return false
}

//Check 判断请求是否为文件请求
func (f *FileEngine) Check(i interface{}) bool {
	return f.FileExist()
}

//GetType 获取引擎结构类型
func (f *FileEngine) GetType() string {
	return ActTypeFile
}

//DynamicEngine 动态请求引擎
//业务处理的相关处理引擎
type DynamicEngine struct {
	cb     *Combination
	method string
	path   string
	handle HandlerFunc
}

var _ IEngine = &DynamicEngine{}

//Init 引擎初始化
//获取请求类型和请求路由 并保存
func (d *DynamicEngine) Init(c *Combination) IEngine {

	method := c.GetMethod()
	path := "/" + strings.Trim(c.GetPath(), "/")

	return &DynamicEngine{
		cb:     c,
		method: method,
		path:   path,
	}

}

//Check 在路由列表中判断 动态请求路由是否已经注册
//如果已经注册 则本次请求由动态引擎处理
//保存路由中已经注册的业务方法
func (d *DynamicEngine) Check(i interface{}) bool {
	if rtg, ok := i.(*RouteGroup); ok {
		Handler, err := rtg.Get(d.method, d.path)
		if err == nil {
			d.handle = Handler
			return true
		}
	}
	return false
}

//Handle 执行已经保存的业务方法
//暂时不做错误返回处理
func (d *DynamicEngine) Handle() (err error) {
	defer doRecover(&err, d.cb)
	err = d.handle(d.cb)
	return err
}

//GetType 获取引擎结构类型
func (d *DynamicEngine) GetType() string {
	return ActTypeDn
}

//WsEngine websocket 引擎
//主要用于处理websocket链接请求
//websocket也属于http请求
//处理逻辑于http请求类似
type WsEngine struct {
	cb     *Combination //请求复合体
	method string       //请求类型 默认赋值为WS
	path   string       //请求的路由
	handle HandlerFunc  //请求方法
}

var _ IEngine = &WsEngine{}

//Init 初始化一个websocket引擎
func (w *WsEngine) Init(c *Combination) IEngine {

	method := MethodWs                           //默认赋值为WS
	path := "/" + strings.Trim(c.GetPath(), "/") //处理路由

	return &WsEngine{
		cb:     c,
		method: method,
		path:   path,
	}
}

//Check 判断请求是否存在于websocket路由列表
//如果路由列表中已经注册，则本次请求由websocket引擎处理
//保存对应的处理方法
func (w *WsEngine) Check(i interface{}) bool {
	if rtg, ok := i.(*RouteGroup); ok {
		Handler, err := rtg.Get(w.method, w.path)
		if err == nil {
			w.handle = Handler
			return true
		}
	}
	return false
}

//Handle 请求处理方法
func (w *WsEngine) Handle() (err error) {
	defer doRecover(&err, w.cb)
	err = w.handle(w.cb)
	return err
}

//GetType 获取引擎结构类型
func (w *WsEngine) GetType() string {
	return ActTypeWs
}
