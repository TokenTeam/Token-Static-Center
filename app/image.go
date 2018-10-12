// Token-Static-Center
// 图片读写模块
// 负责图片的获取、存储
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"bytes"
	"errors"
	"github.com/TokenTeam/Token-Static-Center/db"
	"github.com/TokenTeam/Token-Static-Center/util"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 直接读取图片，不经过处理
// 主要与缓存模块&&反盗链模块&&图片处理模块对接，不暴露给用户，不存在安全隐患
func ReadRawImage(filePath string) (fileStream []byte) {
	file, err := readFile(filePath)

	if err != nil {
		util.ErrorLog("app", err.Error(), "app->ReadRawImage")
		return
	}

	readFileLogger(filePath, "ReadRawImage")
	return file
}

// 直接写入图片，不经过处理
// 主要与缓存模块对接，不暴露给用户，不存在安全隐患
func WriteRawImage(filePath string, fileStream []byte) {
	err := writeFile(filePath, fileStream)

	if err != nil {
		util.ErrorLog("app", err.Error(), "app->WriteRawImage")
		return
	}

	writeFileLogger(filePath, "WriteRawImage")
}

// 读取图片文件
// 对参数有效性的校验已经由app.go完成，因此此处不校验参数有效性
func ReadImage(GUID string, width uint, targetFormat string, quality uint,
	watermarkName string, watermarkPosition uint, watermarkOpacity uint, watermarkSize uint,
	watermarkText string, watermarkColor string, watermarkStyle string) (data []byte, err error) {
	// 根据GUID从数据库获取文件所在路径（相对路径）
	filePath := ""
	year, month, md5, format, err := db.ReadImageDB(GUID)

	// 如果不存在该图片，返回空白
	if year == -2 || month == -2 {
		return nil, errors.New("图片数据不存在")
	}

	// 解决诸如月份09变成9的bug（存储的时候是按照整数进行存储的）
	monthString := strconv.Itoa(month)
	if len(strconv.Itoa(month)) == 1 {
		monthString = "0" + monthString
	}

	// 捕获致命错误
	if err != nil {
		return nil, errors.New("从数据库获取信息失败" + err.Error())
	}
	rootPath, err := getStorageRoot()
	if err != nil {
		return nil, errors.New("获取资源存储目录失败" + err.Error())
	}

	// 拼接文件的真实目录
	filePath = filePath + rootPath + "image/" + strconv.Itoa(year) + "/" + monthString + "/" + GUID + "." + format

	// 如果文件不存在
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, errors.New("图片状态检查时出现致命错误：" + err.Error())
	}

	// 读取文件
	fileData, err := readFile(filePath)
	if err != nil {
		return nil, errors.New("图片读取时出现致命错误：" + err.Error())
	}

	// 读取文件日志
	readFileLogger(filePath, "ReadImage")

	// 检查校验码
	fileHash := util.GetMD5Hash(fileData)

	if fileHash != md5 {
		return nil, errors.New("校验数据时出现错误，请检查文件是否已损毁")
	}

	// 如果转换格式不在支持范围内
	accessableFileTypeInterface, err := util.GetConfig("Image", "AccessableFileType")

	accessableFileTypeSlice := accessableFileTypeInterface.([]string)

	flag := false
	for each := range accessableFileTypeSlice {
		if accessableFileTypeSlice[each] == targetFormat {
			flag = true
			break
		}
	}
	if flag != true {
		return nil, errors.New("格式" + targetFormat + "不受支持")
	}

	// 图片处理（默认按比例缩放）
	fileData, err = ImageResize(fileData, int(width), true)

	if err != nil {
		return nil, errors.New("缩放图片时出现错误：" + err.Error())
	}

	// 如果选项为图片水印
	if watermarkName != "" && watermarkText == "" {
		// 检测水印是否存在
		watermarkPath := rootPath + "watermark/" + watermarkName + ".png"
		_, err := os.Stat(watermarkPath)

		// 如果水印不存在，不添加水印
		if err != nil {
			util.WarningLog("app", "处理图片时发生普通错误：所请求的水印" + watermarkPath + "不存在", "app->ReadRawImage")
		} else {
			// 添加水印
			watermarkFile, _ := ioutil.ReadFile(watermarkPath)
			fileData, err = ImageWatermark(fileData, watermarkFile, watermarkPosition, watermarkOpacity, watermarkSize)
			if err != nil {
				return nil, errors.New("图片水印处理时遭遇致命错误：" + err.Error())
			}
		}
	}

	// 如果选项为文字水印
	if watermarkText != "" && watermarkName == "" {
		// 解析URI Encode之后的字符为可读文字
		watermarkText, err = url.QueryUnescape(watermarkText)
		if err != nil {
			return nil, errors.New("解析URL中文字水印的文字时遭遇致命错误：" + err.Error())
		}
		fileData, err = TextWatermark(fileData, watermarkPosition, watermarkOpacity, watermarkSize, watermarkColor, watermarkText, watermarkStyle)
		if err != nil {
			return nil, errors.New("文字水印处理时遭遇致命错误：" + err.Error())
		}
	}

	// 压缩图片
	fileData, err = ImageCompress(fileData, quality)
	if err != nil {
		return nil, errors.New("压缩图片时遭遇致命错误：" + err.Error())
	}

	// 转换图片格式
	fileData, err = ImageReformat(fileData, targetFormat)
	if err != nil {
		return nil, errors.New("图片格式转换时遭遇致命错误：" + err.Error())
	}

	// 返回文件数据
	return fileData, nil
}

// 写入图片文件
func WriteImage(GUID string, fileStream []byte, defaultFormat string, preCompressLevel int, maxWidth int) (err error, fileSizeByte int) {
	// 在数据库里面检查GUID是否有重复
	year, month, md5, _, err := db.ReadImageDB(GUID)
	// err不为空的情况有可能是数据本身就不存在（正好满足写入图片的需求）
	// 因此需要结合year和month进行判断
	if err != nil && year > 0 && month > 0 {
		return errors.New("写入图片过程中检查GUID是否重复时出现致命错误：" + err.Error()), -1
	}

	// 如果有重复，报错，交给业务引擎返回重试信息（基本不可能，但是还是要检测一下）
	if md5 != "" {
		return errors.New("写入图片过程中检查GUID是否重复时出现致命错误：GUID重复，无法保存数据，请使用不重复的GUID替代之！"), -1
	}

	// 图片尺寸预压缩（按比例预压缩）
	if maxWidth != 0 {
		fileStream, err = ImageResize(fileStream, maxWidth, true)
		if err != nil {
			return errors.New("写入图片过程中对图片尺寸进行预压缩时出现致命错误：" + err.Error()), -1
		}
	}

	// 图片质量预压缩
	if preCompressLevel != 0 {
		compressValue := (10 - preCompressLevel) * 10
		ImageCompress(fileStream, uint(compressValue))
	}

	// 转换成配置文件中的默认格式
	fileStream, err = ImageReformat(fileStream, defaultFormat)
	if err != nil {
		return errors.New("写入图片过程中转换格式时出现致命错误：" + err.Error()), -1
	}

	// 生成待存储文件路径
	storagePath, err := getStorageRoot()
	if err != nil {
		return errors.New("写入图片过程中获取存储根目录时出现致命错误：" + err.Error()), -1
	}
	// 拼接方式：根目录 + image + 年 + 月 + GUID + . + 格式
	t := time.Now()
	year = t.Year()
	month = int(t.Month())
	// 解决诸如月份09变成9的bug（存储的时候是按照整数进行存储的）
	monthString := strconv.Itoa(month)
	if len(strconv.Itoa(month)) == 1 {
		monthString = "0" + monthString
	}
	filePath := storagePath + "image" + "/" + strconv.Itoa(year) + "/" + monthString + "/" + GUID + "." + defaultFormat

	// 生成校验码
	md5 = util.GetMD5Hash(fileStream)

	// 生成文件大小
	finalFileSizeByte := bytes.Count(fileStream, nil)

	// 存储文件
	err = writeFile(filePath, fileStream)
	if err != nil {
		return errors.New("写入图片过程中存储文件时出现致命错误：" + err.Error()), -1
	}

	// 写入文件日志
	writeFileLogger(filePath, "WriteImage")

	// 写入数据库
	// 解决debug模式下AppCode为空，导致数据库插入失败的BUG
	if CurrentUploadAppCode == "" {
		CurrentUploadAppCode = "debug_mode"
	}
	err = db.WriteImageDB(GUID, uint64(finalFileSizeByte), defaultFormat, CurrentUploadAppCode, md5)
	if err != nil {
		return errors.New("写入图片过程中写入数据库时出现致命错误：" + err.Error()), -1
	}

	// 报告保存成功状态
	util.OperationLog("WriteImage", "保存图片成功，保存路径为" + filePath, "app->WriteImage")

	return nil, finalFileSizeByte
}


// 直接读取文件
// filePath必须为绝对路径
func readFile(filePath string) (fileStream []byte, err error) {
	// 检查路径是否合法
	err = DeformityDirectoryFilter(filePath)

	if err != nil {
		return nil, err
	}

	// 检查文件是否存在
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New("文件 " + filePath + " 不存在")
	}
	if file.IsDir() {
		return nil, errors.New("文件 " + filePath + " 是文件夹！")
	}

	// 此前已经检查该文件状态，不需要再获取错误信息
	rawFileStream, _ := ioutil.ReadFile(filePath)

	return rawFileStream, nil
}

// 直接写入文件
// filePath必须为绝对路径
func writeFile(filePath string, fileStream []byte) (err error) {
	// 检查路径是否合法
	err = DeformityDirectoryFilter(filePath)

	if err != nil {
		return err
	}

	// 切割文件名与目录名（目录最高不超过32层）
	filePathArray := make([]string, 32)
	filePathArray = strings.Split(filePath, "/")

	pathName := ""

	for i := 0; i < len(filePathArray) - 1; i++ {
		// 如果已经读取到文件名位置（还没有读入）
		if filePathArray[i + 1] == "" {
			break
		}

		pathName = pathName + filePathArray[i] + "/"
	}


	// 检查所在目录是否存在
	_, err = os.Stat(pathName)

	// 目录不存在，则新建目录
	if err != nil {
		// 递归创建目录
		os.MkdirAll(pathName, 0755)
	}

	// 检查目录下是否已经有该文件
	_, err = os.Stat(filePath)

	// 如果已经存在该文件
	if err == nil {
		return errors.New("文件 " + filePath + " 已存在！")
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, fileStream, 0755)

	if err != nil {
		return errors.New("无法创建文件！请检查权限配置！")
	}

	return
}

// 畸形目录名检测，防止攻击者通过畸形目录名操作不被允许操作的信息
// 目录名中不允许存在的内容：
// - ~
// - ../
// - ./
// - *
// - ?
// - \
// - .db
// - .conf
// - .yaml
func DeformityDirectoryFilter(path string) (err error) {
	if strings.Contains(path, "~") ||
		strings.Contains(path, "../") ||
		strings.Contains(path, "./") ||
		strings.Contains(path, "*") ||
		strings.Contains(path, "?") ||
		strings.Contains(path, "\\") ||
		strings.Contains(path, ".db") ||
		strings.Contains(path, ".conf") ||
		strings.Contains(path, ".yaml") {
		return errors.New("目录" + path + "为畸形目录！")
	}
	return
}

// 获取当前静态资源根目录
// 便于读写文件的时候调用
func getStorageRoot() (rootPath string, err error) {
	staticRootInterface, err := util.GetConfig("Global", "StorageDir")

	if err != nil {
		return "", errors.New("读取静态资源根目录配置时，出现错误：" + err.Error())
	}

	// 转换Interface到String
	staticRootString := staticRootInterface.(string)

	return staticRootString, nil
}