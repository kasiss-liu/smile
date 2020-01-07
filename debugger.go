//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"fmt"
	"runtime/debug"
)

//Debugger 定义外部panic处理函数，可用于覆盖默认值
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
func defaultDebugger(cb *Combination, e error) {
	stack := debug.Stack()
	fmt.Printf("[debug_log] error: %s\n", e.Error())
	fmt.Printf("[debug_log] stacks: %s\n", string(stack))
	cb.Header().Add("Content-Type", "text/html;charset=utf-8;")
	cb.WriteHeader(cb.Status())
	str := fmt.Sprintf("{\"path\":\"%s\",\"status\":\"500\",\"message\":\"Internal Server Error \r\n %s\"}", cb.Request.URL.Path,e.Error())
	cb.WriteString(str)
}

//SetDebugger 由外部注入一个处理error的函数 会替换默认函数
func SetDebugger(fnc Debugger) {
	debugger = fnc
}
