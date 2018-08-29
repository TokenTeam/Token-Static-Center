// Token-Static-Center
// 主业务模块
// 负责主要的业务处理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"net/http"
	"html/template"
)

// 首页
func HomePage(w http.ResponseWriter, r *http.Request) {
	//解析指定模板文件homepage.html
	page, err := template.ParseFiles("template/homepage.html")

	if err != nil {
		errorLogger(r, err)
		return
	}

	accessLogger(r, "HomePage")

	//输出到浏览器
	page.Execute(w, nil)
}
