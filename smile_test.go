package smile

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)

type smileController struct{}

func (t *smileController) PostFunc(c *Combination) error {
	c.WriteString("hello post")
	return nil
}

func (t *smileController) GetFunc(c *Combination) error {
	c.WriteString("hello get")
	return nil
}

func (t *smileController) WsFunc(c *Combination) error {
	c.WriteString("hello ws")
	return nil
}

var sfc = &smileController{}

func TestSmile(t *testing.T) {
	var startChan = make(chan int)
	e := NewEngine("./examples/websocket/")

	go func() {
		go e.Run(":9999")
		go e.RunTLS(":443", "cert.crt", "key.key")
		startChan <- 1
	}()

	<-startChan

	e.GzipOff()
	e.GzipOn()

	e.SetLoger(&Logger{os.Stdout, true})
	e.SetRouteGroup(&RouteGroup{
		GET:  map[string]HandlerFunc{"func": sfc.GetFunc},
		POST: map[string]HandlerFunc{"func": sfc.PostFunc},
		WS:   map[string]HandlerFunc{"func": sfc.WsFunc},
	})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:9999/func", nil)
	e.ServeHTTP(resp, req)
	req = httptest.NewRequest("POST", "http://localhost:9999/func", nil)
	e.ServeHTTP(resp, req)
	LogOFF()
	fmt.Println("------log Off -----")
	req = httptest.NewRequest("WS", "http://localhost:9999/func", nil)
	e.ServeHTTP(resp, req)
	LogON()
	fmt.Println("------log On -----")
	req = httptest.NewRequest("PUT", "http://localhost:9999/test", nil)
	e.ServeHTTP(resp, req)

}

func TestMode(t *testing.T) {
	LogON()
	if logSwitch != true {
		t.Error("logOn failed")
	}
	LogOFF()
	if logSwitch != false {
		t.Error("logOff failed")
	}
	SetTESTING()
	if mode != ModeTESTING {
		t.Error("SetTESTING failed")
	}
	SetDEBUG()
	if mode != ModeDEBUG {
		t.Error("SetDEBUG failed")
	}
	SetPRODUCTION()
	if mode != ModePRO {
		t.Error("SetPRODUCTION failed")
	}
}
