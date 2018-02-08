package smile

import (
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type Combination struct {
	ResponseWriter
	Request *http.Request
}

const (
	MaxFileSize = 1 << 20
)

func InitCombination(w http.ResponseWriter, r *http.Request) *Combination {

	writer := &responseWriter{}
	writer.Init(w)

	r.ParseForm()
	r.ParseMultipartForm(MaxFileSize)
	return &Combination{writer, r}
}

func (c *Combination) GetURL() string {
	return c.Request.URL.String()
}
func (c *Combination) GetPath() string {
	return c.Request.URL.Path
}
func (c *Combination) GetScheme() string {
	return c.Request.URL.Scheme
}
func (c *Combination) GetQueryString() string {
	return c.Request.URL.RawQuery
}
func (c *Combination) GetUserAgent() string {
	return c.Request.UserAgent()
}
func (c *Combination) GetMethod() string {
	return c.Request.Method
}
func (c *Combination) GetProto() string {
	return c.Request.Proto
}
func (c *Combination) GetHost() string {
	return c.Request.Host
}
func (c *Combination) GetHeader() http.Header {
	return c.Request.Header
}
func (c *Combination) GetClientIP() string {

	clientIP := c.GetHeader().Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if clientIP != "" {
		return clientIP
	}
	clientIP = strings.TrimSpace(c.GetHeader().Get("X-Real-Ip"))
	if clientIP != "" {
		return clientIP
	}

	if addr := c.GetHeader().Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (c *Combination) GetRawBody() string {
	byte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	return string(byte)
}
func (c *Combination) GetQueryParam(key string) string {
	return c.Request.Form.Get(key)
}

func (c *Combination) GetPostParam(key string) string {
	return c.Request.PostFormValue(key)
}

func (c *Combination) GetMultipartFormParam(key string) []string {
	return c.Request.MultipartForm.Value[key]
}
func (c *Combination) GetMultipartFormFile(key string) []*multipart.FileHeader {
	return c.Request.MultipartForm.File[key]
}
func (c *Combination) GetCookie(key string) (*http.Cookie, error) {
	return c.Request.Cookie(key)
}
func (c *Combination) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

func (c *Combination) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}
