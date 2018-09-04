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
	"strings"
	"fmt"
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
	//extensionName := ".jpg"
	//
	//// 输出图片
	//w.Header().Set("Content-Type", mime.TypeByExtension(extensionName))
	//w.Write(ReadImage("e44378ac-0237-4331-aaf2-63b8818e5c34", 500, "jpg", "go", 3, 10, 30))

	ImageFetchHandler(w, r)

}


// 图像输出
func ImageFetchHandler(w http.ResponseWriter, r *http.Request) {
	var (
		imageData []byte
		fileExtension string
	)

	// 获取请求（/image/bla-bla-bla-bla.bla）
	requestParam := r.RequestURI

	//// 检查缓存状态
	//if true {
	//
	//} else {
		// 以连字符拆分URI
		var (
			requestPath []string
			requestParams []string
			tempSlice []string
		)

		// Step1: /image/bla-bla-bla-bla.bla => {"", "image", "bla-bla-bla-bla.bla"}
		requestPath = strings.Split(requestParam, "/")
		// Step2: bla-bla-bla-bla.bla => {"bla", "bla", "bla", "bla.bla"}
		requestParams = strings.Split(requestPath[2], "-")
		// Step3: bla.bla => {"bla", "bla"}
		tempSlice = strings.Split(requestParams[len(requestParams) - 1], ".")
		// Step4: {"bla", "bla", "bla", "bla.bla"} + {"bla", "bla"} => {"bla", "bla", "bla", "bla", "bla"}
		requestParams = append(requestParams[0:len(requestParams) - 1], tempSlice[0], tempSlice[1])

		fmt.Println(requestParams)

		// 解决参数过少时引起的故障
		if len(requestParams) < 7 {
			ErrorPage(w, r, 404,"ImageFetchHandler", "参数严重不足，只有" + strconv.Itoa(len(requestParams)) + "个")
			return
		}

		// 解析GUID
		GUID := ""
		// GUID总共由五个片段组成，之间使用连字符进行连接
		for i := 0; i < 4; i++ {
			GUID = GUID + requestParams[i] + "-"
		}
		GUID = GUID + requestParams[4]

		// GUID基础校验（仅校验长度）
		if len(GUID) != 36 {
			// 抛出404错误
			ErrorPage(w, r, 404, "ImageFetchHandler", "GUID校验失败")
			return
		}

		// 根据URL参数个数，筛选输出类型
		width, err := strconv.Atoi(requestParams[5])
		quality, err := strconv.Atoi(requestParams[6])
		// 校验参数有效性
		if err != nil || quality > 100 || quality < 0 {
			ErrorPage(w, r, 404, "ImageFetchHandler", "请求URL中存在不合法数值")
			return
		}
		switch len(requestParams) {
			// 例：http://static2.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80.jpg
			// 即为请求GUID为e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，不带水印的JPG格式图片资源
			// 不带水印获取图片资源
			case 8:
				targetFormat := requestParams[7]
				fileExtension = targetFormat
				imageData = ReadImage(GUID, uint(width), targetFormat, uint(quality), "", 0, 0, 0, "", "", "")
				break
			// 例：http://static2.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80-wutnews-1-30-15.jpg
			// 即为请求GUID为 e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，水印名称为wutnews，水印位置为左上角，水印透明度为30%透明，水印大小为15%宽度（相对于图片宽度）的JPG格式图片资源
			// 带图片水印获取图片资源
			case 12:
				watermarkName := requestParams[7]
				watermarkPosition, err := strconv.Atoi(requestParams[8])
				watermarkOpacity, err := strconv.Atoi(requestParams[9])
				watermarkSize, err := strconv.Atoi(requestParams[10])
				targetFormat := requestParams[11]
				fileExtension = targetFormat
				// 校验参数有效性
				if err != nil || watermarkOpacity > 100 || watermarkOpacity < 0 {
					ErrorPage(w, r, 404, "ImageFetchHandler", "图片水印参数中存在不合法数值")
					return
				}
				imageData = ReadImage(GUID, uint(width), targetFormat, uint(quality), watermarkName, uint(watermarkPosition), uint(watermarkOpacity), uint(watermarkSize), "", "", "")
				break

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
