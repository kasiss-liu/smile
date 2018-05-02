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
	c := InitCombination(w, r, Default(), nil)

	e := &FileEngine{BaseDir: "./examples/"}
	e.Init(c)
	if e.Check(rg) {
		err := e.Handle()
		if err != nil {
			fmt.Println("fileErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	} else {
		fmt.Println("check:", e.Check(rg))
	}
	e.Reset()
}

func TestDynamicEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	c := InitCombination(w, r, Default(), nil)

	e := &DynamicEngine{}
	e.Init(c)
	if e.Check(rg) {
		err := e.Handle()
		if err != nil {
			fmt.Println("handleErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	}

	e.Reset()

}

func TestWsEngine(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("WS", "/test2", nil)
	c := InitCombination(w, r, Default(), nil)

	e := &WsEngine{}
	e.Init(c)
	if e.Check(rg) {
		err := e.Handle()
		if err != nil {
			fmt.Println("handleErr:", err.Error())
		} else {
			fmt.Println("success")
		}
	} else {
		fmt.Println("checkERR:", e.Check(rg))
	}
	e.Reset()

}
