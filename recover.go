package smile

import (
	"errors"
	"fmt"
	"net/http"
)

const recoverPrefix = "Rocover panic: "

// Recovery 定义外部panic处理函数，可用于覆盖默认值
type Recovery func(*Combination, interface{}) error

// 储存一个处理panic的函数
var recovery Recovery

func init() {
	//使用默认的recover处理函数
	recovery = defaultRecover
}

// recover钩子 在engine调用handler时 会被defer触发
// 调用注册的panic处理函数
func doRecover(e *error, c *Combination) error {
	if r := recover(); r != nil {
		*e = recovery(c, r)
	}
	return *e
}

//默认panic处理函数
func defaultRecover(c *Combination, rec interface{}) error {
	c.WriteHeader(http.StatusInternalServerError)
	s := fmt.Sprintf("%#v", rec)
	return errors.New(recoverPrefix + s)
}

//SetRecovery 由外部注入一个recover函数 会替换默认函数
func SetRecovery(fnc Recovery) {
	recovery = fnc
}
