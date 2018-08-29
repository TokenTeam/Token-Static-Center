package app

import (
	"net/http"
	"html/template"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"strings"
	"github.com/LuRenJiasWorld/Token-Static-Center/security"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	//解析指定模板文件index.html
	page, err := template.ParseFiles("template/homepage.html")

	if err != nil {
		errorLogger(r, err)
		return
	}

	accessLogger(r, "HomePage")

	//输出到浏览器
	page.Execute(w, nil)
}


// 获取所请求的链接地址
func getRequestURI(r *http.Request) (uri string) {
	// 该业务只提供HTTP服务，但是因为反向代理对外是HTTPS，因此使用HTTPS的scheme
	scheme := "https://"
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

// 记录访问信息
func accessLogger(r *http.Request, module string) {
	util.AccessLog("app", "访问请求：" + getRequestURI(r) + "，客户端IP：" + security.GetRequestIP(r), "app->" + module)
}


// 记录错误信息
func errorLogger(r *http.Request, err error) {
	util.ErrorLog("app", "访问请求 " + getRequestURI(r) + " 时页面渲染失败：" + err.Error() + "，客户端IP：" + security.GetRequestIP(r), "app->HomePage")
}
