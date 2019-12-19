//This software is licensed under the MIT License.
//You can get more info in license file.

package smile

//定义了一些模式和HOOK 方便调测
const (
	ModeDEBUG   = "debug"
	ModeTESTING = "testing"
	ModePRO     = "production"
)

//模式
var mode string

//SetDEBUG 开发模式
func SetDEBUG() {
	mode = ModeDEBUG
}

//SetPRODUCTION 生产模式
func SetPRODUCTION() {
	mode = ModePRO
}

//SetTESTING 测试模式
func SetTESTING() {
	mode = ModeTESTING
}

//Mode 返回当前模式
func Mode() string {
	return mode
}

//日志开关
//是否开启日志功能
var logSwitch = true

//LogON 开启日志
func LogON() {
	logSwitch = true
}

//LogOFF 关闭日志
func LogOFF() {
	logSwitch = false
}
