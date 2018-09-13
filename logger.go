//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"fmt"
	"io"
	"time"
)

//ILogger 日志处理
type ILogger interface {
	Write(io.Writer, string)
	Log(*Combination)
}

//注册基本的terminal颜色
var (
	green       = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white       = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow      = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red         = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue        = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta     = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan        = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset       = string([]byte{27, 91, 48, 109})
	enableColor = true
)

//预定义几种日志前缀
var (
	infoprefix    = "[INFO]"
	warningprefix = "[WARNING]"
	errorprefix   = "[ERROR]"
	logprefix     = "[LOG]"
)

//Logger 实现一个日志结构
type Logger struct {
	Writer io.Writer
	isTerm bool
}

var _ ILogger = &Logger{}

//logger写方法 向一个io内写入数据
func (l *Logger) Write(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}

//TermOn 开启终端输出（数据染色）
func (l *Logger) TermOn() {
	l.isTerm = true
}

//TermOff 关闭终端输出（数据染色取消）
func (l *Logger) TermOff() {
	l.isTerm = false
}

//Log 方法
//整理需要log的数据拼接后进行输出
func (l *Logger) Log(c *Combination) {

	statusCode := c.ResponseWriter.Status()   //请求响应状态
	prefix := prefixForStatus(statusCode)     //日志前缀
	statusColor := colorForStatus(statusCode) //状态打印颜色
	method := c.GetMethod()                   //请求方法
	methodColor := colorForMethod(method)     //方法打印颜色
	clientIP := c.GetClientIP()               //客户端ip
	path := c.GetURL()                        //请求路由
	//如果不是终端 则不输出颜色
	if !l.isTerm {
		statusColor = ""
		methodColor = ""
		reset = ""
	}
	//将准备的数据进行拼接
	s := l.joinLog(prefix, path, statusCode, statusColor, clientIP, method, methodColor)
	//数据写入
	l.Write(l.Writer, s)
}

//拼接日志字符串
//[LOG]2018/02/07-18:19:20 GET /test Status 200 IP 127.0.0.1
//[%s] %v |%s %3d %s| %15s |%s %-7s %s %s,
func (l *Logger) joinLog(prefix string, path string, statusCode int, statusColor string, clientIP string, method string, methodColor string) string {
	s := fmt.Sprintf("[SMILE %s]%v |%s %-7s %s| %s | code %s %3d %s|  ClientIP %s",
		prefix,
		time.Now().Format("2006/01/02 15:04:05"),
		methodColor, method, reset,
		path,
		statusColor, statusCode, reset,
		clientIP,
	)
	return s
}

//根据不同的状态获取日志前缀
func prefixForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return logprefix
	case status >= 300 && status < 400:
		return infoprefix
	case status >= 400 && status < 500:
		return warningprefix
	default:
		return errorprefix

	}
}

//根据不同的状态获取颜色
func colorForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return cyan
	case status >= 400 && status < 500:
		return blue
	default:
		return red

	}
}

//根据不同的方法获取颜色
func colorForMethod(method string) string {
	switch method {
	case MethodPost:
		return magenta
	case MethodGet:
		return yellow
	case MethodPut:
		return blue
	case MethodDelete:
		return red
	default:
		return white
	}
}
