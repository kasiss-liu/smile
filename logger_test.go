package smile

import (
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {

	logger := &Logger{
		os.Stdout,
		true,
	}

	logger.TermOff()
	logger.TermOn()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/websocket/index.html", nil)
	c := InitCombination(w, r, Default())

	c.WriteString("hello world")

	logger.Log(c)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/websocket/index.html", nil)
	c = InitCombination(w, r, Default())
	c.WriteHeader(301)

	logger.Log(c)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("WS", "/websocket/index.html", nil)
	c = InitCombination(w, r, Default())
	c.WriteHeader(404)

	logger.Log(c)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("WS", "/websocket/index.html", nil)
	c = InitCombination(w, r, Default())
	c.WriteHeader(500)

	logger.Log(c)
}
