// Token-Static-Center
// 主业务模块
// 负责主要的业务处理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"net/http"
	"html/template"
	"strconv"
	"errors"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"mime"
)

// 首页
func HomePage(w http.ResponseWriter, r *http.Request) {
	////解析指定模板文件homepage.html
	//page, err := template.ParseFiles("template/homepage.html")
	//
	//if err != nil {
	//	errorLogger(r, err)
	//	return
	//}
	//
	//
	//accessLogger(r, "HomePage")
	//
	////输出到浏览器
	//page.Execute(w, nil)


	// 配置后缀名
	extensionName := ".jpg"

	// 输出图片
	w.Header().Set("Content-Type", mime.TypeByExtension(extensionName))
	w.Write(ReadImage("e44378ac-0237-4331-aaf2-63b8818e5c34", 23, "jpg", ""))



}


// 图像输出
func ImageFetchHandler(w http.ResponseWriter, r *http.Request) {

}

// 图像上传
func ImageUploadHandler() {

}

// 错误页面
func ErrorPage(w http.ResponseWriter, r *http.Request, errorType int) {
	page, err := template.ParseFiles("template/" + strconv.Itoa(errorType) + ".html")

	if err != nil {
		errorLogger(r, err)
		return
	}

	errorLogger(r, errors.New("接收访问请求时出现错误，错误码：" + strconv.Itoa(errorType)))

	// 输出自定义错误码
	w.WriteHeader(errorType)

	//输出到浏览器
	page.Execute(w, nil)
}

// 直接返回防盗链图片
func AntiLeechImage(w http.ResponseWriter, r *http.Request) {
	// 获取静态资源
	staticFilePathInterface, err := util.GetConfig("Global", "StorageDir")

	if err != nil {
		util.ErrorLog("app", "无法获取图片路径配置！请检查配置文件！", "app->AntiLeechImage")
		return
	}

	// 转换Interface到String
	staticFilePathString := staticFilePathInterface.(string)

	// 拼接警告图片文件链接
	AntiLeechImageFilePath := staticFilePathString + "others/anti-leech.jpg"

	// 读取图片
	fileStream := ReadRawImage(AntiLeechImageFilePath)

	// 配置后缀名
	extensionName := ".jpg"

	// 输出图片
	w.Header().Set("Content-Type", mime.TypeByExtension(extensionName))
	w.Write(fileStream)
}
