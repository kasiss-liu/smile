package smile

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

func debugFunc(c *Combination) error {
	return errors.New("testing debug error catch")
}

func TestDoDebugger(t *testing.T) {
	rg.SetGET("/test_debug", debugFunc)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/test_debug", nil)
	c := InitCombination(w, r, Default())

	e := &DynamicEngine{}
	engine := e.Init(c)
	if engine.Check(rg) {
		err := engine.Handle()
		if err != nil {
			doDebug(err, c)
		} else {
			fmt.Println("success")
		}
	}

}
