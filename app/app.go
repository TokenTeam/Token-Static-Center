// Token-Static-Center
// 主业务模块
// 负责主要的业务处理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TokenTeam/Token-Static-Center/db"
	"github.com/TokenTeam/Token-Static-Center/util"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

// 全局变量：上传资源的AppCode（用于多模块共享，便于存储AppCode到数据库）
// 默认为debug_mode，因为debug模式开启的时候不会进行AppCode检查
var CurrentUploadAppCode string

// 首页
func HomePage(w http.ResponseWriter, r *http.Request) {
	//解析指定模板文件homepage.html
	page, err := template.ParseFiles("template/homepage.html")

	if err != nil {
		errorLogger(r, err, "HomePage")
		return
	}

	accessLogger(r, "HomePage")

	//输出到浏览器
	page.Execute(w, nil)
}


// 图像输出
func ImageFetchHandler(w http.ResponseWriter, r *http.Request) {
	// 记录访问请求
	accessLogger(r,"ImageFetchHandler")

	// 记录执行开始时间，用于执行速率统计
	startTime := time.Now()

	// 获取请求（/image/bla-bla-bla-bla.bla）
	requestParam := r.RequestURI

	// 以连字符拆分URI
	var (
		imageData []byte
		fileExtension string
		requestPath []string
		requestParams []string
		tempSlice []string
	)

	// 检查水印信息中是否含有非法字符
	// 非法字符列表：
	// - #
	if strings.Contains(requestParam, "#") {
			ErrorPage(w, r, 404, "ImageFetchHandler", "该请求含有非法字符，被安全模块拦截")
			return
	}

	// 检查URL信息中是否缺少必须含有的字符
	// 必须含有的字符列表
	// - .
	// - -
	// 避免存在example.com/images/blabla#blabla.bla，#号截断其后文字造成后续字符串裁切出现索引失败的状况
	if strings.Contains(requestParam, ".") != true ||
		strings.Contains(requestParam, "-") != true {
			ErrorPage(w, r, 404, "ImageFetchHandler", "该请求未包含必需字符，被安全模块拦截")
			return
	}

	// Step1: /image/bla-bla-bla-bla.bla => {"", "image", "bla-bla-bla-bla.bla"}
	requestPath = strings.Split(requestParam, "/")
	// Step2: bla-bla-bla-bla.bla => {"bla", "bla", "bla", "bla.bla"}
	requestParams = strings.Split(requestPath[2], "-")
	// Step3: bla.bla => {"bla", "bla"}
	tempSlice = strings.Split(requestParams[len(requestParams) - 1], ".")
	// Step4: {"bla", "bla", "bla", "bla.bla"} + {"bla", "bla"} => {"bla", "bla", "bla", "bla", "bla"}
	requestParams = append(requestParams[0:len(requestParams) - 1], tempSlice[0], tempSlice[1])

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

	// 缓存开关检查
	cacheConfig, err := util.GetConfig("Cache", "Status")
	// 转换Interface到String
	cacheConfig = cacheConfig.(string)
	if err != nil {
		ErrorPage(w, r, 500, "ImageFetchHandler", "缓存开关获取失败，请检查配置文件状态！" + err.Error())
		return
	}

	// 调试模式检查
	debugConfig, err := util.GetConfig("Global", "Debug")
	// 转换Interface到String
	debugConfig = debugConfig.(string)
	if err != nil {
		ErrorPage(w, r, 500, "ImageFetchHandler", "调试模式状态获取失败，请检查配置文件状态！" + err.Error())
		return
	}

	// 调试模式下提醒缓存被旁通
	if debugConfig == "on" {
		util.WarningLog("app", "调试模式已开启，缓存机制将会被旁通", "app->ImageFetchHandler")
	}

	// 读取缓存
	cacheStatus, imageData, cacheSizeByte, err := CacheFetchHandler(requestPath[2])

	if err != nil {
		ErrorPage(w, r, 500, "ImageFetchHandler", "读取缓存失败，原因：" + err.Error())
	}

	// 检查缓存状态
	// 如果缓存已开启，并且缓存存在，而且配置文件没有打开调试模式
	if cacheConfig == "on" && cacheStatus == true && debugConfig != "on" {
		// 切分出文件扩展名
		extensionName := strings.Split(requestPath[2], ".")[1]
		w.Header().Set("Content-Type", mime.TypeByExtension(extensionName))
		w.Write(imageData)
		util.OperationLog("app", "缓存成功命中，读取缓存成功，缓存大小：" + strconv.Itoa(cacheSizeByte / 1024) + "kb", "app->ImageFetchHandler")
	} else {
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
				imageData, err = ReadImage(GUID, uint(width), targetFormat, uint(quality), "", 0, 0, 0, "", "", "")
				if err != nil {
					ErrorPage(w, r, 404, "ImageFetchHandler", "图像处理模块返回错误信息：" + err.Error())
					return
				}
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
				imageData, err = ReadImage(GUID, uint(width), targetFormat, uint(quality), watermarkName, uint(watermarkPosition), uint(watermarkOpacity), uint(watermarkSize), "", "", "")
				if err != nil {
					ErrorPage(w, r, 404, "ImageFetchHandler", "图像处理模块返回错误信息：" + err.Error())
					return
				}
				break
			// 例：http://static2.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80-%40Token+Team-1-20-30-FFFFFF-regular.jpg
			// 即为请求GUID为 e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，水印文本为@Token Team，水印位置为左上角，水印透明度为20%透明，水印字体大小为30px，水印颜色为FFFFFF，水印字体样式为普通字体样式的JPG格式图片资源
			// 带文字水印获取图片资源
			case 14:
				watermarkText := requestParams[7]
				watermarkPosition, err := strconv.Atoi(requestParams[8])
				watermarkOpacity, err := strconv.Atoi(requestParams[9])
				watermarkSize, err := strconv.Atoi(requestParams[10])
				watermarkColor := requestParams[11]
				watermarkStyle := requestParams[12]
				targetFormat := requestParams[13]
				fileExtension = targetFormat
				// 校验参数有效性
				if err != nil || watermarkOpacity > 100 || watermarkOpacity < 0 {
					ErrorPage(w, r, 404, "ImageFetchHandler", "文字水印参数中存在不合法数值")
					return
				}
				imageData, err = ReadImage(GUID, uint(width), targetFormat, uint(quality), "", uint(watermarkPosition), uint(watermarkOpacity), uint(watermarkSize), watermarkText, watermarkColor, watermarkStyle)
				if err != nil {
					ErrorPage(w, r, 404, "ImageFetchHandler", "图像处理模块返回错误信息：" + err.Error())
					return
				}
				break
			// 错误捕获
			default:
				ErrorPage(w, r, 404, "ImageFetchHandler", "请求URL参数个数不正确")
				return
		}

		// 校验返回图像数据有效性，如果返回数据为空，报错
		if imageData == nil {
			ErrorPage(w, r, 404, "ImageFetchHandler", "图像读取与处理时，出现致命错误，返回空数据，原因：" + err.Error())
			return
		}

		// 输出图片
		extensionName := "." + fileExtension
		w.Header().Set("Content-Type", mime.TypeByExtension(extensionName))
		w.Write(imageData)

		// 执行缓存操作
		if debugConfig != "on" && cacheConfig == "on" {
			cacheLocation, cacheSizeByte, err := CacheWriteHandler(imageData, requestPath[2])
			if err != nil {
				util.ErrorLog("app", "执行存储缓存失败，原因：" + err.Error(), "app->ImageFetchHandler")
			}
			util.OperationLog("app", "执行存储缓存成功，缓存位置：" + cacheLocation + "，缓存大小：" + strconv.Itoa(cacheSizeByte / 1024) + "kb", "app->ImageFetchHandler")
		}

	}

	// 记录访问流量与频次
	err = db.DownloadCounter(GUID, cacheSizeByte)
	if err != nil {
		util.ErrorLog("app", "记录访问失败，原因：" + err.Error(), "app->ImageFetchHandler")
	}

	// 记录执行时间
	elapsedTime := time.Since(startTime)
	util.OperationLog("app", "执行图片获取请求成功，耗时" + strconv.FormatInt(int64(elapsedTime) / 1000000, 10) + "ms", "app->ImageFetchHandler")

}

// 图像上传
func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 记录上传请求
	accessLogger(r, "ImageUploadHandler")

	// 记录执行开始时间，用于执行速率统计
	startTime := time.Now()

	// 读取配置文件
	// 最大宽度
	maxWidthInterface, err1 := util.GetConfig("Image", "MaxWidth")
	maxWidth, err2 := strconv.Atoi(maxWidthInterface.(string))

	// 可上传文件后缀名类型
	uploadableFileTypeInterface, err3 := util.GetConfig("Image", "UploadableFileType")
	uploadableFileTypeStringSlice := uploadableFileTypeInterface.([]string)

	// 存储文件类型
	storageFileTypeInterface, err4 := util.GetConfig("Image", "StorageFileType")
	storageFileType := storageFileTypeInterface.(string)

	// 压缩等级
	compressLevelInterface, err5 := util.GetConfig("Image", "JpegCompressLevel")
	compressLevel := compressLevelInterface.(int)

	// 最大文件体积
	maxImageFileSizeInterface, err6 := util.GetConfig("Image", "MaxImageFileSize")
	maxImageFileSize := maxImageFileSizeInterface.(int)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -1, "图像上传过程中读取配置文件时出现致命错误，请检查配置文件是否完整！")
		return
	}

	// 解析表单
	uploadFile, handle, err := r.FormFile("image")
	if err != nil {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -2, "解析表单时出现错误，请检查表单是否正确！")
		return
	}
	defer uploadFile.Close()

	// 获取文件
	imageFileStream := bytes.NewBuffer(nil)
	_, err = io.Copy(imageFileStream, uploadFile)
	if err != nil {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -3, "从表单读取文件时出现错误！")
		return
	}

	// 检查文件大小是否超过限额
	fileSizeMB := handle.Size / (1024 * 1024)
	if int(fileSizeMB) > maxImageFileSize {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -4, "文件大小超过限制，最大允许" + strconv.Itoa(maxImageFileSize) + "MB，上传文件大小为" + strconv.Itoa(int(fileSizeMB)) + "MB")
		return
	}

	// 检查文件格式是否在配置文件中
	fileFormat := strings.ToLower(path.Ext(handle.Filename))
	tempSlice := strings.Split(fileFormat, ".")
	fileFormat = tempSlice[1]

	formatHitStatus := false
	for i := 0; i < len(uploadableFileTypeStringSlice); i++ {
		if uploadableFileTypeStringSlice[i] == fileFormat {
			formatHitStatus = true
			break
		}
	}

	if formatHitStatus == false {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -5, "你上传的格式" + fileFormat + "不受支持！")
		return
	}

	// 生成GUID
	guidObj := uuid.NewV4()
	guid := fmt.Sprintf("%s", guidObj)

	// 写入图像
	err = WriteImage(guid, imageFileStream.Bytes(), storageFileType, compressLevel, maxWidth)
	if err != nil {
		JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", -6, "保存图像失败：" + err.Error())
		return
	}

	//fmt.Println(guid, bytes.Count(imageFileStream.Bytes(), nil), storageFileType, compressLevel, maxWidth)

	// 返回数据
	JsonReturn(w, r, "ImageUploadHandler", "app->ImageUploadHandler", 0, guid)

	// 记录执行时间
	elapsedTime := time.Since(startTime)
	util.OperationLog("app", "执行图片上传请求成功，耗时" + strconv.FormatInt(int64(elapsedTime) / 1000000, 10) + "ms", "app->ImageFetchHandler")
}

// 错误页面
func ErrorPage(w http.ResponseWriter, r *http.Request, errorType int, errorModule string, errorMessage string) {
	page, err := template.ParseFiles("template/" + strconv.Itoa(errorType) + ".html")

	if err != nil {
		errorLogger(r, err, errorModule)
		return
	}

	// 自动记录错误日志
	switch r.Method {
		case "GET":
			errorLogger(r, errors.New("接收访问请求时出现错误，错误码：" + strconv.Itoa(errorType) + "，相关信息：" + errorMessage), errorModule)
			break
		case "POST":
			errorLogger(r, errors.New("接收上传请求时出现错误，错误码：" + strconv.Itoa(errorType) + "，相关信息：" + errorMessage), errorModule)
			break
		default:
			errorLogger(r, errors.New("接收方法为" + r.Method + "的请求时出现错误，错误码：" + strconv.Itoa(errorType) + "，相关信息：" + errorMessage), errorModule)
	}


	// 输出自定义错误码
	w.WriteHeader(errorType)

	//输出到浏览器
	page.Execute(w, nil)
}

// 返回Json数据(直接输出到浏览器)
// errNumber 错误码，如果错误码表示错误，值为负数，会自动计入错误日志中
func JsonReturn(w http.ResponseWriter, r *http.Request, module string, trace string, errNumber int, message string) {

	// 返回Json数据格式
	type Json struct{
		Errno int			`json:"error_code"`
		Message string		`json:"message"`
	}

	data := Json{errNumber, message}

	jsonDataByte, err := json.Marshal(data)

	if err != nil {
		util.ErrorLog("JsonReturn", "生成Json数据时出现错误：" + err.Error(), "app->JsonReturn")
		return
	}

	// 记录错误日志
	if errNumber < 0 {
		util.ErrorLog(module, message, trace)
	}

	w.Write(jsonDataByte)

	return
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
