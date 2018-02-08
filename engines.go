package smile

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type IEngine interface {
	Init(*Combination)
	Handle() error
	Check(interface{}) bool
	Reset()
}

var _ IEngine = &FileEngine{}

type FileEngine struct {
	BaseDir  string
	FilePath string
	FileExt  string
	cb       *Combination
}

func (f *FileEngine) Init(c *Combination) {
	f.cb = c
	f.FilePath = f.BaseDir + strings.Trim(c.GetPath(), "/")
	f.FileExt = path.Ext(f.FilePath)
}

func (f *FileEngine) Handle() (err error) {
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
	content, err := f.readFile()
	f.cb.Write(content)
	f.Reset()
	if err != nil {
		return err
	}
	return nil
}
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

func (f *FileEngine) Check(i interface{}) bool {
	return f.FileExist()
}

func (f *FileEngine) readFile() ([]byte, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return []byte{}, errors.New("[" + f.FilePath + "]" + err.Error())
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
func (f *FileEngine) Reset() {
	f.FilePath = ""
	f.cb = nil
}

type DynamicEngine struct {
	cb     *Combination
	method string
	path   string
	handle HandlerFunc
}

var _ IEngine = &DynamicEngine{}

func (d *DynamicEngine) Init(c *Combination) {
	d.cb = c
	d.method = c.GetMethod()
	d.path = strings.Trim(c.GetPath(), "/")
}
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

func (d *DynamicEngine) Handle() (err error) {
	d.handle(d.cb)
	d.Reset()
	return nil
}

func (d *DynamicEngine) Reset() {
	d.cb = nil
	d.method = ""
	d.path = ""
	d.handle = nil
}

type WsEngine struct {
	cb     *Combination
	method string
	path   string
	handle HandlerFunc
}

var _ IEngine = &WsEngine{}

func (w *WsEngine) Init(c *Combination) {
	w.cb = c
	w.method = METHOD_WS
	w.path = strings.Trim(c.GetPath(), "/")
}

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

func (w *WsEngine) Handle() (err error) {
	w.handle(w.cb)
	w.Reset()
	return nil
}

func (w *WsEngine) Reset() {
	w.path = ""
	w.handle = nil
	w.cb = nil
}
