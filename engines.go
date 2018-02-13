//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//一个引擎接口
type IEngine interface {
	Init(*Combination)      //初始化引擎
	Handle() error          //执行方法
	Check(interface{}) bool //判断是否属于引擎处理请求
	Reset()                 //重置引擎数据
}

var _ IEngine = &FileEngine{}

//文件请求引擎
//响应静态文件数据请求
type FileEngine struct {
	BaseDir  string       //文件仓库地址
	FilePath string       //文件路径
	FileExt  string       //文件后缀
	cb       *Combination //请求复合结构
}

//文件引擎初始化
//储存请求复合结构
//处理请求文件地址
//处理请求文件后缀
func (f *FileEngine) Init(c *Combination) {
	f.cb = c
	f.FilePath = f.BaseDir + strings.Trim(c.GetPath(), "/")
	f.FileExt = path.Ext(f.FilePath)
}

//响应处理方法
func (f *FileEngine) Handle() (err error) {
	//根据不同的后缀 输出不同的请求头
	//目前只涉及有常见后缀文件
	switch f.FileExt {
	case ".css":
		f.cb.SetHeader("Content-Type", "text/css")
	case ".js":
		f.cb.SetHeader("Content-Type", "text/javascript")
	case ".jpg":
		fallthrough
	case ".jpeg":
		f.cb.SetHeader("Content-Type", "image/jpeg")
	case ".png":
		f.cb.SetHeader("Content-Type", "image/png")
	case ".gif":
		f.cb.SetHeader("Content-Type", "image/gif")
	case ".html":
		fallthrough
	case ".htm":
		f.cb.SetHeader("Content-Type", "text/html")
	default:
		f.cb.SetHeader("Content-Type", "text/plain")
	}
	//读取文件内容 发送到客户端
	content, err := f.readFile()
	f.cb.Write(content)
	//重置
	f.Reset()
	if err != nil {
		return err
	}
	return nil
}

//判断请求文件是否存在
func (f *FileEngine) FileExist() bool {
	//获取文件/文件夹信息
	file, err := os.Stat(f.FilePath)
	//如果存在 则判断是否是文件夹
	if err == nil && !file.IsDir() {
		return true
	}
	//不存在或者不是文件 返回false
	return false
}

//判断请求是否为文件请求
func (f *FileEngine) Check(i interface{}) bool {
	return f.FileExist()
}

//读取文件内容
func (f *FileEngine) readFile() ([]byte, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		//如果文件存在 但打开文件失败 则返回error
		return []byte{}, errors.New("[" + f.FilePath + "]" + err.Error())
	}
	//读取文件内容并返回
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

//将引擎的文件路径、后缀名、请求数据清空
func (f *FileEngine) Reset() {
	f.FilePath = ""
	f.FileExt = ""
	f.cb = nil
}

//动态请求引擎
//业务处理的相关处理引擎
type DynamicEngine struct {
	cb     *Combination
	method string
	path   string
	handle HandlerFunc
}

var _ IEngine = &DynamicEngine{}

//引擎初始化
//获取请求类型和请求路由 并保存
func (d *DynamicEngine) Init(c *Combination) {
	d.cb = c
	d.method = c.GetMethod()
	d.path = strings.Trim(c.GetPath(), "/")
}

//在路由列表中判断 动态请求路由是否已经注册
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

//执行已经保存的业务方法
//暂时不做错误返回处理
func (d *DynamicEngine) Handle() (err error) {
	d.handle(d.cb)
	d.Reset()
	return nil
}

//重置动态请求
//将请求的类型、方法、路由全部重置
func (d *DynamicEngine) Reset() {
	d.cb = nil
	d.method = ""
	d.path = ""
	d.handle = nil
}

//websocket 引擎
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

//初始化一个websocket引擎
func (w *WsEngine) Init(c *Combination) {
	w.cb = c
	w.method = METHOD_WS                    //默认赋值为WS
	w.path = strings.Trim(c.GetPath(), "/") //处理路由
}

//判断请求是否存在于websocket路由列表
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

//请求处理方法
func (w *WsEngine) Handle() (err error) {
	w.handle(w.cb)
	w.Reset()
	return nil
}

//重置引擎
func (w *WsEngine) Reset() {
	w.path = ""
	w.handle = nil
	w.cb = nil
}
