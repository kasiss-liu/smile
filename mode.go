package smile

var monitorSwitch bool = false

func MonitorOn() {
	monitorSwitch = true
}

func MonitorOff() {
	monitorSwitch = false
}

func MonitorStatus() bool {
	return monitorSwitch
}

var mode string

func SetDEBUG() {
	mode = "debug"
}
func SetPRODUCTION() {
	mode = "production"
}
func SetTESTING() {
	mode = "testing"
}
func Mode() string {
	return mode
}

var logSwitch bool = true

func LogON() {
	logSwitch = true
}

func LogOFF() {
	logSwitch = false
}
