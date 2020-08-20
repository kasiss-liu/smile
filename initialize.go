package smile

import (
	"sync"
)

//InitFunc 初始化函数类型
type InitFunc func()

var initFuncList []InitFunc

func init() {
	initFuncList = make([]InitFunc, 0, 5)
}

//是否执行过初始化操作
var initState = false

//
var initWithGoroutine = false

//InitFuncPush  将函数注入到框架启动初始化函数列表中
func InitFuncPush(f InitFunc) {
	initFuncList = append(initFuncList, f)
}

//DoCustomInit 显式提供可调用api 调用时机
func DoCustomInit() {
	if initState {
		panic(`初始化操作已经执行`)
	}
	if initWithGoroutine {
		routineInit()
	} else {
		syncInit()
	}
	initState = true
}

//使用goroutine异步执行初始化操作
func routineInit() {
	var wg sync.WaitGroup
	for _, initFunc := range initFuncList {
		wg.Add(1)
		go func(f InitFunc) {
			f()
			wg.Done()
		}(initFunc)
	}
	wg.Wait()
}

//采用线性方式执行初始化操作
func syncInit() {
	for _, initFunc := range initFuncList {
		initFunc()
	}
}

//GetInitState 获取初始化操作状态
func GetInitState() bool {
	return initState
}

//SetInitGoroutine 设置初始化时采用协程
func SetInitGoroutine() {
	initWithGoroutine = true
}

//SetInitSync 设置初始化时不采用协程
func SetInitSync() {
	initWithGoroutine = false
}
