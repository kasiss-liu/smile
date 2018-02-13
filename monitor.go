//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

import (
	"time"
)

//monitor接口 可以注册任何实现该接口的结构到引擎中
type RunMonitor interface {
	HandleStart(*MonitorInfo)
	HandleEnd(*MonitorInfo)
}

//接口方法的注入参数
type MonitorInfo struct {
	CaseTime    time.Time
	Method      string
	Path        string
	Combination *Combination
}
