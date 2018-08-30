// Token-Static-Center
// 日志模板
// 便于输出更加规范的业务日志，增强代码复用能力
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"net/http"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
)

// 记录访问信息
func accessLogger(r *http.Request, module string) {
	util.AccessLog("app", "处理请求：" + util.GetRequestURI(r) + "，客户端IP：" + util.GetRequestIP(r), "app->" + module)
}


// 记录错误信息
func errorLogger(r *http.Request, err error) {
	util.ErrorLog("app", "处理请求 " + util.GetRequestURI(r) + " 时页面渲染失败：" + err.Error() + "，客户端IP：" + util.GetRequestIP(r), "app->HomePage")
}


// 记录文件读取信息
func readFileLogger(filePath string, module string) {
	util.OperationLog("app", "读取文件 " + filePath, "app->" + module)
}

// 记录文件写入信息
func writeFileLogger(storagePath string, module string) {
	util.OperationLog("app", "写入文件成功！存储位置：" + storagePath, "app->" + module)
}

