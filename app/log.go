package app

import (
	"net/http"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
)

// 记录访问信息
func accessLogger(r *http.Request, module string) {
	util.AccessLog("app", "访问请求：" + util.GetRequestURI(r) + "，客户端IP：" + util.GetRequestIP(r), "app->" + module)
}


// 记录错误信息
func errorLogger(r *http.Request, err error) {
	util.ErrorLog("app", "访问请求 " + util.GetRequestURI(r) + " 时页面渲染失败：" + err.Error() + "，客户端IP：" + util.GetRequestIP(r), "app->HomePage")
}

