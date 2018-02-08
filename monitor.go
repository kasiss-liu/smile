package smile

import (
	"time"
)

type RunMonitor interface {
	HandleStart(*MonitorInfo)
	HandleEnd(*MonitorInfo)
}

type MonitorInfo struct {
	CaseTime    time.Time
	Method      string
	Path        string
	Combination *Combination
}
