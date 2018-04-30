//Testing File for Routes
package smile

import (
	"testing"
)

//测试需要的Controller类

type testFillController struct{}

func (t *testFillController) Func(c *Combination) error {
	return nil
}

var tfc *testFillController = &testFillController{}

func TestFillRoutes(t *testing.T) {
	rg := NewRouteGroup()
	rg.FillRoutes(METHOD_GET, "", tfc)
	rg.FillRoutes(METHOD_POST, "", tfc)
	rg.FillRoutes(METHOD_WS, "", tfc)
	fn, err := rg.Get(METHOD_GET, "func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}

}

//测试自动匹配前缀
type testController struct{}

func (t *testController) PostFunc(c *Combination) error {
	return nil
}

func (t *testController) GetFunc(c *Combination) error {
	return nil
}

func (t *testController) WsFunc(c *Combination) error {
	return nil
}

var tc *testController = &testController{}

func TestFillPrefixRoutes(t *testing.T) {
	rg := NewRouteGroup()
	rg.PrefixFillRoutes("", tc)
	fn, err := rg.Get(METHOD_GET, "func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}
	fn, err = rg.Get(METHOD_POST, "func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}

	fn, err = rg.Get(METHOD_WS, "func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}
	fn, err = rg.Get(METHOD_GET, "nilfunc")
	if err == nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Log("success")
	}

}