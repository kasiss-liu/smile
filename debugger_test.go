package smile

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func debugFunc(c *Context) error {
	return errors.New("testing debug error catch")
}

func TestDoDebugger(t *testing.T) {
	rg.SetGET("/test_debug", debugFunc)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/test_debug", nil)
	c := initContext(w, r, Default())

	e := createEngine(false)
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			doDebug(err, c)
		} else {
			t.Log("success")
		}
	}

}
