package smile

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

var rg *RouteGroup

func init() {
	rg = NewRouteGroup()
	rg.SetGET("test", testHandle)
	rg.SetWS("test2", testHandle)
}

func testHandle(c *Context) error {

	c.WriteHeader(200)
	c.WriteString("hello world")
	status := c.Status()
	fmt.Println("status:", status)
	datasize := c.DataSize()
	fmt.Println("dataSize:", datasize)

	return nil
}

func TestFileEngine(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/websocket/index.html", nil)
	c := initContext(w, r, Default())

	e := createEngine(true, "./examples/","index.html")
	engine := e.Init(c)
	t.Log(engine.GetType())
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			t.Log("fileErr:", err.Error())
		} else {
			t.Log("success")
		}
	} else {
		t.Log("check:", engine.Check(rg))
	}
}

func TestDynamicEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	c := initContext(w, r, Default())

	e := createEngine(false)
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			t.Error("handleErr:", err.Error())
		} else {
			t.Log("success")
		}
	}
}

func TestWsEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("WS", "/test2", nil)
	c := initContext(w, r, Default())

	e := createEngine(false)
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			t.Error("handleErr:", err.Error())
		} else {
			t.Log("success")
		}
	} else {
		t.Log("checkERR:", engine.Check(rg))
	}
}
