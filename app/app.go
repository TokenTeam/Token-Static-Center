package app

import (
	"net/http"
	"html/template"
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
