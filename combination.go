//This software is licensed under the MIT License.
//You can get more info in license file.
package smile

import (
	"compress/gzip"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

//一个复合结构，将writer和 request保存到一起，方便被调用
//实现了一些便捷方法 从而缩短获取数据的路径长度
type Combination struct {
	ResponseWriter
	Request *http.Request
}

//文件上传大小限制
const (
	MaxFileSize = 1 << 20
)

//初始化一个*Combination
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
	//解析传参数据
	r.ParseForm()
	r.ParseMultipartForm(MaxFileSize)
	return &Combination{writer, r}
}

//获取请求的URL
func (c *Combination) GetURL() string {
	return c.Request.URL.String()
}

//获取请求的Path
func (c *Combination) GetPath() string {
	return c.Request.URL.Path
}

//获取请求Scheme
func (c *Combination) GetScheme() string {
	return c.Request.URL.Scheme
}

//获取请求的url参数
func (c *Combination) GetQueryString() string {
	return c.Request.URL.RawQuery
}

//获取请求的代理头 user-agent
func (c *Combination) GetUserAgent() string {
	return c.Request.UserAgent()
}

//获取请求的方法 GET/POST
func (c *Combination) GetMethod() string {
	return c.Request.Method
}

//获取请求的传输协议 HTTP1.1 / HTTP2
func (c *Combination) GetProto() string {
	return c.Request.Proto
}

//获取请求host
func (c *Combination) GetHost() string {
	return c.Request.Host
}

//获取请求的header
//返回http.Header
func (c *Combination) GetHeader() http.Header {
	return c.Request.Header
}

//获取请求的IP地址
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

//获取请求中的body体 当传输类型为 urlencoded时 可用
func (c *Combination) GetRawBody() string {
	byte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	return string(byte)
}

//根据键名从url参数中取值
func (c *Combination) GetQueryParam(key string) string {
	return c.Request.Form.Get(key)
}

//根据键名从post表单中取值
func (c *Combination) GetPostParam(key string) string {
	return c.Request.PostFormValue(key)
}

//根据键名从form-data类型中取值
func (c *Combination) GetMultipartFormParam(key string) []string {
	if c.Request.MultipartForm == nil {
		return nil
	}
	return c.Request.MultipartForm.Value[key]
}

//根据键名冲form-data类型中取得上传文件
func (c *Combination) GetMultipartFormFile(key string) []*multipart.FileHeader {
	if c.Request.MultipartForm == nil {
		return nil
	}
	return c.Request.MultipartForm.File[key]
}

//从请求携带的cookie中取值
func (c *Combination) GetCookie(key string) (*http.Cookie, error) {
	return c.Request.Cookie(key)
}

//设置cookie
func (c *Combination) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

//设置header
func (c *Combination) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

//请求响应结束后的一些操作
func (c *Combination) Close() {
	//如果本次请求使用gzip压缩 则关闭资源
	if c.ResponseWriter.(*responseWriter).Gz() {
		if c.ResponseWriter.(*responseWriter).Writer.(*gzip.Writer) != nil {
			c.ResponseWriter.(*responseWriter).Writer.(*gzip.Writer).Close()
		}
	}
}
