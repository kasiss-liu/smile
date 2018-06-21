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

func testHandle(c *Combination) error {

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
	c := InitCombination(w, r, Default())

	e := &FileEngine{BaseDir: "./examples/"}
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			fmt.Println("fileErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	} else {
		fmt.Println("check:", engine.Check(rg))
	}
}

func TestDynamicEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	c := InitCombination(w, r, Default())

	e := &DynamicEngine{}
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			fmt.Println("handleErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	}
}

func TestWsEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("WS", "/test2", nil)
	c := InitCombination(w, r, Default())

	e := &WsEngine{}
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			fmt.Println("handleErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	} else {
		fmt.Println("checkERR:", engine.Check(rg))
	}
}
