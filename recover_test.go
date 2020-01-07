package smile

import (
	"net/http/httptest"
	"testing"
)

func init() {
	rg.SetGET("test_recover", recoverFunc)
}

func recoverFunc(c *Combination) error {
	testRecover()
	return nil
}
func testRecover() {
	panic("byte error")
}

func TestDoRecover(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/test_recover", nil)
	c := InitCombination(w, r, Default())

	e := createEngine(false)
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			t.Log("handleErr:", err.Error())
		} else {
			t.Log("success")
		}
	}
}
