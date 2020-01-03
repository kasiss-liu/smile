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

var tfc = &testFillController{}

func TestFillRoutes(t *testing.T) {
	rg := NewRouteGroup()
	rg.SetPathStyleHump()
	rg.FillRoutes(MethodGet, "", tfc)
	rg.FillRoutes(MethodPost, "", tfc)
	rg.FillRoutes(MethodWs, "", tfc)
	fn, err := rg.Get(MethodGet, "/Func")
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
	c.WriteString("hello post")
	return nil
}

func (t *testController) GetFuncTest(c *Combination) error {
	c.WriteString("hello get")
	return nil
}

func (t *testController) WsFunc(c *Combination) error {
	c.WriteString("hello ws")
	return nil
}

var tc = &testController{}

func TestFillPrefixRoutes(t *testing.T) {
	rg := NewRouteGroup()
	rg.SetPathStyleConnector()
	rg.PrefixFillRoutes("", tc)
	fn, err := rg.Get(MethodGet, "/func-test")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}
	fn, err = rg.Get(MethodPost, "/func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}

	fn, err = rg.Get(MethodWs, "/func")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}
	fn, err = rg.Get(MethodGet, "nilfunc")
	if err == nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Log("success")
	}

	rg.SetPathStyleHump()
	rg.PrefixFillRoutes("", tc)
	fn, err = rg.Get(MethodGet, "/FuncTest")
	if err != nil {
		t.Errorf("%#v\n", rg)
		t.Errorf(err.Error())
	} else {
		t.Logf("%#v\n", fn)
	}
}

func TestFormatRouteGroup(t *testing.T) {
	rg := NewRouteGroup()
	rg.SetPathStyleConnector()
	rg.PrefixFillRoutes("", tc)
	rg.SetGET("/gettest", tc.GetFuncTest)
	rsAssign, rsAuto := rg.FormatRoutes()
	t.Logf("%#v\n%#v\n", rsAssign, rsAuto)
	doPrintRoutes(rsAssign, rsAuto)
}

func TestSetRouteFuncs(t *testing.T) {
	rg := NewRouteGroup()
	rg.SetRoute404(func(cb *Combination) error {
		t.Log("testing route404")
		return nil
	})
	rg.SetMiddleware(func(cb *Combination) error {
		t.Log("testing middleware")
		return nil
	})
}
