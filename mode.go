//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

//定义了一些模式和HOOK 方便调测
const (
	MODE_DEBUG   = "debug"
	MODE_TESTING = "testing"
	MODE_PRO     = "production"
)

//Hook开关
//开启后 将执行注册在引擎内的monitor方法
var monitorSwitch bool = false

//开启monitor
func MonitorOn() {
	monitorSwitch = true
}

//关闭monitor
func MonitorOff() {
	monitorSwitch = false
}

//获取monitor状态
func MonitorStatus() bool {
	return monitorSwitch
}

//模式
var mode string

//开发模式
func SetDEBUG() {
	mode = MODE_DEBUG
}

//生产模式
func SetPRODUCTION() {
	mode = MODE_PRO
}

//测试模式
func SetTESTING() {
	mode = MODE_TESTING
}

//返回当前模式
func Mode() string {
	return mode
}

//日志开关
//是否开启日志功能
var logSwitch bool = true

//开启日志
func LogON() {
	logSwitch = true
}

//关闭日志
func LogOFF() {
	logSwitch = false
}
