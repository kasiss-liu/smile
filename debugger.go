package smile

import (
	"fmt"
)

//Debug 定义外部panic处理函数，可用于覆盖默认值
type Debugger func(*Combination, error)

// 储存一个处理panic的函数
var debugger Debugger

func init() {
	//使用默认的recover处理函数
	debugger = defaultDebugger
}

// debug钩子 handleFunc 返回error时被调用
// 调用注册的error处理函数
func doDebug(e error, c *Combination) {
	debugger(c, e)
}

//默认debug函数
func defaultDebugger(c *Combination, e error) {
	if Mode() != ModePRO {
		fmt.Printf("[debug_log] %s\n", e.Error())
	}
}

//SetDebugger 由外部注入一个处理error的函数 会替换默认函数
func SetDebugger(fnc Debugger) {
	debugger = fnc
}
