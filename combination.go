//Package smile This software is licensed under the MIT License.
//You can get more info in license file.
package smile

import (
	"compress/gzip"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

//handlerChain 请求方法调用链
type handlerChain struct {
	current int 
	handlerList []HandlerFunc
	aborted  bool
}

func (hc *handlerChain) add(fn HandlerFunc) {
	hc.handlerList = append(hc.handlerList,fn)
}

func (hc *handlerChain) next(cb *Combination) (err error) {

	defer func(){
		if hc.current >= len(hc.handlerList) {
			hc.abort()
		}
	}()

	err = errors.New("HanlderChain reached end")

	if hc.current < len(hc.handlerList) {	
		if fn := hc.handlerList[hc.current];fn != nil  {
			hc.current++
			if !hc.aborted {
				err = fn(cb)
				if err == nil {
					hc.next(cb)
				}
			}else {
				err = nil;
			}
			return
		}
	}
	return
}

func (hc *handlerChain) reset() {
	hc.aborted = false
	hc.current = 0
}
func (hc *handlerChain) isAorted() bool {
	return hc.aborted
}

func (hc *handlerChain) abort() {
	hc.aborted = true
}

func newHanlderChain() *handlerChain {
	return &handlerChain{
		0,
		make([]HandlerFunc,0,5),
		false,
	}
}



//Combination 一个复合结构，将writer和 request保存到一起，方便被调用
//实现了一些便捷方法 从而缩短获取数据的路径长度
type Combination struct {
	handlerChain *handlerChain
	ResponseWriter
	Request *http.Request
}

//默认文件上传大小限制
const (
	MaxFileSize = 5 << 20
)

//用户自定上传大小限制
var (
	CustomFileSize int64
)

//InitCombination 初始化一个*Combination
//解析url传参 解析form-data
func InitCombination(w http.ResponseWriter, r *http.Request, e *Engine) *Combination {

	writer := &responseWriter{}

	//初始化http.ResponseWriter
	writer.Init(w)
	//如果开启了Gzip 则设置响应的headers 并将writer的io.writer调整为gzwriter
	if e.Gzip && strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") {
		//默认生成一个*gzip.Writer
		//请求结束后关闭
		gz := gzip.NewWriter(w)
		//设置响应头 告知浏览器本次请求为gzip压缩
		writer.Header().Set("Content-Encoding", "gzip")
		//非常重要 如果不设置此头 浏览器将不解析gzip
		writer.Header().Set("Transfer-Encoding", "chunked")
		writer.GzOn(gz)
	}
	var FileSize int64
	if CustomFileSize > 0 {
		FileSize = CustomFileSize
	} else {
		FileSize = MaxFileSize
	}
	//解析传参数据
	r.ParseForm()
	r.ParseMultipartForm(FileSize)
	cb := &Combination{ResponseWriter:writer, Request: r,handlerChain: newHanlderChain()}
	//初始化中间件

	return cb
}

//GetURL 获取请求的URL
func (c *Combination) GetURL() string {
	return c.Request.URL.String()
}

//GetPath 获取请求的Path
func (c *Combination) GetPath() string {
	return c.Request.URL.Path
}

//GetScheme 获取请求Scheme
func (c *Combination) GetScheme() string {
	return c.Request.URL.Scheme
}

//GetQueryString 获取请求的url参数
func (c *Combination) GetQueryString() string {
	return c.Request.URL.RawQuery
}

//GetUserAgent 获取请求的代理头 user-agent
func (c *Combination) GetUserAgent() string {
	return c.Request.UserAgent()
}

//GetMethod 获取请求的方法 GET/POST
func (c *Combination) GetMethod() string {
	return c.Request.Method
}

//GetProto 获取请求的传输协议 HTTP1.1 / HTTP2
func (c *Combination) GetProto() string {
	return c.Request.Proto
}

//GetHost 获取请求host
func (c *Combination) GetHost() string {
	return c.Request.Host
}

//GetHeader 获取请求的header
//返回http.Header
func (c *Combination) GetHeader() http.Header {
	return c.Request.Header
}

//GetClientIP 获取请求的IP地址
//从请求头中截取
func (c *Combination) GetClientIP() string {
	//如果请求头中含有 X-Forwarded-For 则首先取用该值
	clientIP := c.GetHeader().Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if clientIP != "" {
		return clientIP
	}
	//如果 X-Forwarded-For 值无效 则使用 X-Real-Ip 的值
	clientIP = strings.TrimSpace(c.GetHeader().Get("X-Real-Ip"))
	if clientIP != "" {
		return clientIP
	}
	//针对app的请求头获取到ip
	if addr := c.GetHeader().Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}
	//以上手段均失败时 尝试从请求地址截取ip
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

//GetRawBody 获取请求中的body体 当传输类型为 urlencoded时 可用
func (c *Combination) GetRawBody() string {
	byte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	return string(byte)
}

//GetQueryParam 根据键名从url参数中取值
func (c *Combination) GetQueryParam(key string) string {
	return c.Request.Form.Get(key)
}

//GetPostParam 根据键名从post表单中取值
func (c *Combination) GetPostParam(key string) string {
	return c.Request.PostFormValue(key)
}

//GetMultipartFormParam 根据键名从form-data类型中取值
func (c *Combination) GetMultipartFormParam(key string) []string {
	if c.Request.MultipartForm == nil {
		return nil
	}
	return c.Request.MultipartForm.Value[key]
}

//GetMultipartFormFile 根据键名冲form-data类型中取得上传文件头信息
func (c *Combination) GetMultipartFormFile(key string) []*multipart.FileHeader {
	if c.Request.MultipartForm == nil {
		return nil
	}
	return c.Request.MultipartForm.File[key]
}

//GetFormFile 根据键名获取上传文件
func (c *Combination) GetFormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

//GetCookie 从请求携带的cookie中取值
func (c *Combination) GetCookie(key string) (*http.Cookie, error) {
	return c.Request.Cookie(key)
}

//SetCookie 设置cookie
func (c *Combination) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

//SetHeader 设置header
func (c *Combination) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

//Close 请求响应结束后的一些操作
func (c *Combination) Close() {
	//如果本次请求使用gzip压缩 则关闭资源
	if c.ResponseWriter.(*responseWriter).Gz() {
		if c.ResponseWriter.(*responseWriter).Writer.(*gzip.Writer) != nil {
			c.ResponseWriter.(*responseWriter).Writer.(*gzip.Writer).Close()
		}
	}
}

//Redirect 302跳转到指定地址
func (c *Combination) Redirect(url string) {
	c.WriteHeader(http.StatusFound)
	c.Header().Set("Location", url)
	c.Done()
}

//Next 执行注册请求
func (c *Combination) Next() error {
	if !c.handlerChain.isAorted() {
		return c.handlerChain.next(c)
	}
	return errors.New("Combination is aborted")
}
//Abort 调用中断执行 后续注册函数将不再执行
func (c *Combination) Abort() bool {
	c.handlerChain.abort()
	return c.handlerChain.aborted
}