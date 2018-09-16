// Token-Static-Center
// 日志记录模块
// 管理日志记录、输出与存储
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"fmt"
	"time"
)

// 写入到文件的日志缓存
var logCache [10]string

// 日志结构
type log struct {
	log_type string				// 日志类别（access、operation、warning、error）
	log_time time.Time			// 日志时间
	log_module string			// 日志所记录的模块
	log_trace string			// 日志相关操作涉及到的路径
	log_content string			// 日志所返回相关信息
}

// Web访问日志
func AccessLog(logModule string, logContent string, logTrace string) {
	// 检查配置文件是否关闭了该项日志的记录

	logLine := packLog("access", logModule, logContent, logTrace)
	consoleOutput(logLine)
	fileOutput()
}

// 操作日志
func OperationLog (logModule string, logContent string, logTrace string) {
	logLine := packLog("operation", logModule, logContent, logTrace)
	consoleOutput(logLine)
	fileOutput()
}


// 告警日志
func WarningLog(logModule string, logContent string, logTrace string) {
	logLine := packLog("warning", logModule, logContent, logTrace)
	consoleOutput(logLine)
	fileOutput()
}


// 错误日志
func ErrorLog(logModule string, logContent string, logTrace string) {
	logLine := packLog("error", logModule, logContent, logTrace)
	consoleOutput(logLine)
	fileOutput()
}

// 输出日志到控制台（颜色可配置）
// 对应颜色
// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见
func consoleOutput(logLine log) {
	var (
		// 前景色
		foreGroundColor int
		// 背景色
		backGroundColor int
		// 字符特效
		consoleEffect int
	)

	switch logLine.log_type {
		case "access":
			foreGroundColor = 37
			backGroundColor = 0
			consoleEffect = 1
			break
		case "operation":
			foreGroundColor = 34
			backGroundColor = 0
			consoleEffect = 1
			break
		case "warning":
			foreGroundColor = 33
			backGroundColor = 0
			consoleEffect = 1
			break
		case "error":
			foreGroundColor = 31
			backGroundColor = 0
			consoleEffect = 1
			break
		default:
			// 不接受其他类别的输出
			return
	}

	// 检查调试模式是否开启
	// 只有在调试模式下才会输出信息到控制台
	if IsDebug() == true {
		fmt.Printf("%c[%d;%d;%dm [%s] @ %s (%s) %s (相关路径: %s) %c[0m\n", 0x1B, consoleEffect, backGroundColor, foreGroundColor, logLine.log_type, logLine.log_module, logLine.log_time, logLine.log_content, logLine.log_trace, 0x1B)
	}

	return
}

// 记录日志
func fileOutput() {
	// TODO: 记录日志到文件
}

// 将传入的参数打包成一个log对象
// logType类别（注意大小写）
// - access
// - operation
// - warning
// - error
// 其余类别不会被输出
func packLog(logType string, logModule string, logContent string, logTrace string) (logObj log) {
	logLine := log{}
	t := time.Now()

	logLine.log_type = logType
	// 时间格式类似于2018-08-28 22:36:38 +0800 CST
	logLine.log_time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),0,time.Local)
	logLine.log_module = logModule
	logLine.log_trace = logTrace
	logLine.log_content = logContent

	return logLine
}

// 检查是否开启调试模式
func IsDebug() (status bool) {
	debugStatus, err := GetConfig("Global", "Debug")

	if err != nil {
		fmt.Println("发生致命错误： ", err, " ，无法获取调试模式状态，默认开启调试模式以便排查错误")
		return true
	}

	// 将Interface转换成String
	debugStatus = debugStatus.(string)

	switch debugStatus {
		case "on":
			return true
		case "off":
			return false
		default:
			// 为了安全考虑，除非出现配置文件读取错误，否则默认关闭调试输出
			return false
	}
}

// 获取日志状态
func logStatus() (accessLogStatusResult bool, operationLogStatusResult bool, warningLogStatusResult bool, errorLogStatusResult bool) {
	// 读取全部配置文件
	accessLogStatus, err1 := GetConfig("Global", "Log", "LogAccess")
	operationLogStatus, err2 := GetConfig("Global", "Log", "LogOperation")
	warningLogStatus, err3 := GetConfig("Global", "Log", "LogWarning")
	errorLogStatus, err4 := GetConfig("Global", "Log", "LogError")

	// 此处错误捕获不完全原因：日志模块本身出现错误，逐层抛出错误并记录日志也无法实现
	// 但是会直接返回true，也就是默认记录当前项的日志
	// 避免由于配置文件错误导致日志未能记录的风险
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return true, true, true, true
	}

	// 转换Interface到String
	accessLogStatusString := accessLogStatus.(string)
	operationLogStatusString := operationLogStatus.(string)
	warningLogStatusString := warningLogStatus.(string)
	errorLogStatusString := errorLogStatus.(string)

	// 默认均为关闭状态
	accessLogStatusResult = false
	operationLogStatusResult = false
	warningLogStatusResult = false
	errorLogStatusResult = false

	// 判断配置状态，进行对应项的开启
	if accessLogStatusString == "on" {
		accessLogStatusResult = true
	}

	if operationLogStatusString == "on" {
		operationLogStatusResult = true
	}

	if warningLogStatusString == "on" {
		warningLogStatusResult = true
	}

	if errorLogStatusString == "on" {
		errorLogStatusResult = true
	}

	return accessLogStatusResult, operationLogStatusResult, warningLogStatusResult, errorLogStatusResult
}