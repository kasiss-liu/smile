package smile

import (
	"fmt"
	"io"
	"time"
)

type ILogger interface {
	Write(io.Writer, string)
	Log(*Combination)
}

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

var (
	infoprefix    = "[INFO]"
	warningprefix = "[WARNING]"
	errorprefix   = "[ERROR]"
)

type Logger struct {
	Writer io.Writer
	isTerm bool
}

var _ ILogger = &Logger{}

func (l *Logger) Write(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}
func (l *Logger) Log(c *Combination) {
	prefix := "LOG"
	statusCode := c.ResponseWriter.Status()
	statusColor := colorForStatus(statusCode)
	method := c.GetMethod()
	methodColor := colorForMethod(method)
	clientIP := c.GetClientIP()
	path := c.GetPath()
	//如果不是终端 则不输出颜色
	if !l.isTerm {
		statusColor = ""
		methodColor = ""
	}
	s := l.joinLog(prefix, path, statusCode, statusColor, clientIP, method, methodColor)
	l.Write(l.Writer, s)
}

//[LOG]2018/02/07-18:19:20 GET /test Status 200 IP 127.0.0.1
//[%s] %v |%s %3d %s| %15s |%s %-7s %s %s\n,
func (l *Logger) joinLog(prefix string, path string, statusCode int, statusColor string, clientIP string, method string, methodColor string) string {
	s := fmt.Sprintf("[SMILE %s]%v |%s %-7s %s| %s | code %s %3d %s|  ClientIP %s\n",
		prefix,
		time.Now().Format("2006/01/02 15:04:05"),
		methodColor, method, reset,
		path,
		statusColor, statusCode, reset,
		clientIP,
	)
	return s
}

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

func colorForMethod(method string) string {
	switch method {
	case "POST":
		return magenta
	case "GET":
		return yellow
	default:
		return white
	}
}
