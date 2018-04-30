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

//常见的文件后缀对应的输出类型
var CONTENTTYPE = map[string]string{
	"hqx":     "application/mac-binhex40",
	"cpt":     "application/mac-compactpro",
	"doc":     "application/msword",
	"bin":     "application/octet-stream",
	"dms":     "application/octet-stream",
	"lha":     "application/octet-stream",
	"lzh":     "application/octet-stream",
	"exe":     "application/octet-stream",
	"class":   "application/octet-stream",
	"so":      "application/octet-stream",
	"dll":     "application/octet-stream",
	"oda":     "application/oda",
	"pdf":     "application/pdf",
	"ai":      "application/postscript",
	"eps":     "application/postscript",
	"ps":      "application/postscript",
	"smi":     "application/smil",
	"smil":    "application/smil",
	"mif":     "application/vnd.mif",
	"xls":     "application/vnd.ms-excel",
	"ppt":     "application/vnd.ms-powerpoint",
	"wbxml":   "application/vnd.wap.wbxml",
	"wmlc":    "application/vnd.wap.wmlc",
	"wmlsc":   "application/vnd.wap.wmlscriptc",
	"bcpio":   "application/x-bcpio",
	"vcd":     "application/x-cdlink",
	"pgn":     "application/x-chess-pgn",
	"cpio":    "application/x-cpio",
	"csh":     "application/x-csh",
	"dcr":     "application/x-director",
	"dir":     "application/x-director",
	"dxr":     "application/x-director",
	"dvi":     "application/x-dvi",
	"spl":     "application/x-futuresplash",
	"gtar":    "application/x-gtar",
	"hdf":     "application/x-hdf",
	"js":      "application/x-javascript",
	"skp":     "application/x-koan",
	"skd":     "application/x-koan",
	"skt":     "application/x-koan",
	"skm":     "application/x-koan",
	"latex":   "application/x-latex",
	"nc":      "application/x-netcdf",
	"cdf":     "application/x-netcdf",
	"sh":      "application/x-sh",
	"shar":    "application/x-shar",
	"swf":     "application/x-shockwave-flash",
	"sit":     "application/x-stuffit",
	"sv4cpio": "application/x-sv4cpio",
	"sv4crc":  "application/x-sv4crc",
	"tar":     "application/x-tar",
	"tcl":     "application/x-tcl",
	"tex":     "application/x-tex",
	"texinfo": "application/x-texinfo",
	"texi":    "application/x-texinfo",
	"t":       "application/x-troff",
	"tr":      "application/x-troff",
	"roff":    "application/x-troff",
	"man":     "application/x-troff-man",
	"me":      "application/x-troff-me",
	"ms":      "application/x-troff-ms",
	"ustar":   "application/x-ustar",
	"src":     "application/x-wais-source",
	"xhtml":   "application/xhtml+xml",
	"xht":     "application/xhtml+xml",
	"zip":     "application/zip",
	"au":      "audio/basic",
	"snd":     "audio/basic",
	"mid":     "audio/midi",
	"midi":    "audio/midi",
	"kar":     "audio/midi",
	"mpga":    "audio/mpeg",
	"mp2":     "audio/mpeg",
	"mp3":     "audio/mpeg",
	"aif":     "audio/x-aiff",
	"aiff":    "audio/x-aiff",
	"aifc":    "audio/x-aiff",
	"m3u":     "audio/x-mpegurl",
	"ram":     "audio/x-pn-realaudio",
	"rm":      "audio/x-pn-realaudio",
	"rpm":     "audio/x-pn-realaudio-plugin",
	"ra":      "audio/x-realaudio",
	"wav":     "audio/x-wav",
	"pdb":     "chemical/x-pdb",
	"xyz":     "chemical/x-xyz",
	"bmp":     "image/bmp",
	"gif":     "image/gif",
	"ico":     "image/x-icon",
	"ief":     "image/ief",
	"jpeg":    "image/jpeg",
	"jpg":     "image/jpeg",
	"jpe":     "image/jpeg",
	"png":     "image/png",
	"tiff":    "image/tiff",
	"tif":     "image/tiff",
	"djvu":    "image/vnd.djvu",
	"djv":     "image/vnd.djvu",
	"wbmp":    "image/vnd.wap.wbmp",
	"ras":     "image/x-cmu-raster",
	"pnm":     "image/x-portable-anymap",
	"pbm":     "image/x-portable-bitmap",
	"pgm":     "image/x-portable-graymap",
	"ppm":     "image/x-portable-pixmap",
	"rgb":     "image/x-rgb",
	"xbm":     "image/x-xbitmap",
	"xpm":     "image/x-xpixmap",
	"xwd":     "image/x-xwindowdump",
	"igs":     "model/iges",
	"iges":    "model/iges",
	"msh":     "model/mesh",
	"mesh":    "model/mesh",
	"silo":    "model/mesh",
	"wrl":     "model/vrml",
	"vrml":    "model/vrml",
	"css":     "text/css",
	"html":    "text/html",
	"htm":     "text/html",
	"asc":     "text/plain",
	"txt":     "text/plain",
	"rtx":     "text/richtext",
	"rtf":     "text/rtf",
	"sgml":    "text/sgml",
	"sgm":     "text/sgml",
	"tsv":     "text/tab-separated-values",
	"wml":     "text/vnd.wap.wml",
	"wmls":    "text/vnd.wap.wmlscript",
	"etx":     "text/x-setext",
	"xsl":     "text/xml",
	"xml":     "text/xml",
	"mpeg":    "video/mpeg",
	"mpg":     "video/mpeg",
	"mpe":     "video/mpeg",
	"qt":      "video/quicktime",
	"mov":     "video/quicktime",
	"mxu":     "video/vnd.mpegurl",
	"avi":     "video/x-msvideo",
	"movie":   "video/x-sgi-movie",
	"ice":     "x-conference/x-cooltalk",
}

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

const DEFAULT_FILE = "index.html"

//文件引擎初始化
//储存请求复合结构
//处理请求文件地址
//处理请求文件后缀
func (f *FileEngine) Init(c *Combination) {
	f.cb = c
	filename := strings.Trim(c.GetPath(), "/")
	if filename == "" {
		filename = DEFAULT_FILE
	}
	f.FilePath = f.BaseDir + filename
	f.FileExt = path.Ext(f.FilePath)
}

//响应处理方法
func (f *FileEngine) Handle() (err error) {
	//根据不同的后缀 输出不同的请求头
	//目前只涉及有常见后缀文件
	ext := strings.Trim(f.FileExt, ".")
	if contentType, ok := CONTENTTYPE[ext]; ok {
		f.cb.SetHeader("Content-Type", contentType)
	} else {
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
	//针对默认页面做校验

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
