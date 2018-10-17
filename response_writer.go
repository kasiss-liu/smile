//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"
)

const (
	defaultDataSize = 0
	defaultStatus   = 200
)

//ResponseWriter 定义一个writer接口
//该接口可用于http、weibsocket的响应操作
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	DataSize() int
	Status() int
	WriteString(string) (int, error)
	Done()
}

//实现一个ResoponsWriter接口
type responseWriter struct {
	http.ResponseWriter
	io.Writer
	gz      bool //是否开启gz
	status  int  //响应状态
	size    int  //响应字节长度
	written bool
}

//初始化http.ResponseWriter 响应状态 响应数据长度
func (w *responseWriter) Init(writer http.ResponseWriter) {
	w.size = defaultDataSize
	w.status = defaultStatus
	w.ResponseWriter = writer
	w.written = false
}

//开启gz开关
//注册一个新的 *gzip.Writer
//对于本次请求响应将进行gzip压缩
func (w *responseWriter) GzOn(gz *gzip.Writer) {
	w.Writer = gz
	w.gz = true
}

//获取响应数据字节长度
func (w *responseWriter) DataSize() int {
	return w.size
}

//获取响应状态
func (w *responseWriter) Status() int {
	return w.status
}

//写入响应状态到header
func (w *responseWriter) WriteHeader(code int) {
	if code != w.status && code > 0 {
		w.status = code
	}
}

//判断是否已经设置过响应状态
func (w *responseWriter) isWritten() bool {
	return w.written
}

//如果在响应头没有设置的情况下
//立即写入状态至响应头
//并将writer的数据重置
func (w *responseWriter) WriteHeaderAtOnce() {
	if !w.isWritten() {
		w.written = true
		w.ResponseWriter.WriteHeader(w.status)
	}

}

func (w *responseWriter) Gz() bool {
	return w.gz
}

//响应结构的写方法
//如果本次请求为gzip压缩
//则使用注册为*gzip.Writer的io.writer进行写操作
//否则是用http.ResponseWriter进行写操作
func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderAtOnce()
	if w.gz {
		n, err = w.Writer.Write(data)
	} else {
		n, err = w.ResponseWriter.Write(data)
	}
	w.size += n
	return
}

//Done 向请求中写入空字节 以结束请求
//用来执行跳转 或者单纯的header设置
func (w *responseWriter) Done() {
	w.size = 0
	w.Write([]byte(""))
}

//直接写入字符串
func (w *responseWriter) WriteString(data string) (n int, err error) {
	w.Write([]byte(data))
	w.size += n
	return
}

//继承http.Hijacker的Hijack()方法
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()

}

//继承http.flusher的Flush()方法
func (w *responseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

//继承http.CloseNotifier的CloseNotify()方法
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
