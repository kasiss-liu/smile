package smile

import (
	"fmt"
	"testing"
)

func init1() {
	fmt.Println(`init1`)
}

func init2() {
	fmt.Println(`init2`)
}

func init3() {
	fmt.Println(`init3`)
}

func TestInitilization(t *testing.T) {
	initState = false

	InitFuncPush(init1)
	InitFuncPush(init2)
	InitFuncPush(init3)
	t.Log(initWithGoroutine)
	SetInitGoroutine()
	t.Log(initWithGoroutine)
	DoCustomInit()
	t.Log(GetInitState())
}

func TestSyncInitilization(t *testing.T) {
	initState = false

	InitFuncPush(init1)
	InitFuncPush(init2)
	InitFuncPush(init3)
	t.Log(initWithGoroutine)
	SetInitSync()
	t.Log(initWithGoroutine)
	DoCustomInit()
}
