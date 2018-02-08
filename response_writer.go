package smile

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

const (
	isWritten       = -1
	defaultDataSize = 0
	defaultStatus   = 200
)

type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	DataSize() int
	Status() int
	WriteString(string) (int, error)
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) Init(writer http.ResponseWriter) {
	w.size = defaultDataSize
	w.status = isWritten
	w.ResponseWriter = writer
}

func (w *responseWriter) DataSize() int {
	return w.size
}
func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) WriteHeader(code int) {
	if code != w.status && code > 0 {
		if !w.isWritten() {
			w.status = code
		}

	}
}

func (w *responseWriter) isWritten() bool {
	return w.status != isWritten
}
func (w *responseWriter) WriteHeaderAtOnce() {
	if !w.isWritten() {
		w.size = 0
		w.status = defaultStatus
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderAtOnce()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *responseWriter) WriteString(data string) (n int, err error) {
	w.WriteHeaderAtOnce()
	n, err = io.WriteString(w.ResponseWriter, data)
	w.size += n
	return
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
func (w *responseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
