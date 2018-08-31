// Token-Static-Center
// 图片处理模块
// 负责图片的获取、存储、处理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"os"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"io/ioutil"
	"errors"
	"strings"
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
func WriteRawImage(filePath string, fileStream []byte) (err error) {
	err = writeFile(filePath, fileStream)

	if err != nil {
		util.ErrorLog("app", err.Error(), "app->WriteRawImage")
		return
	}

	writeFileLogger(filePath, "WriteRawImage")
	return
}

// 读取图片文件
func ReadImage(GUID string, width uint32, format string, watermarkName string) (data []byte) {
	// 根据GUID从数据库获取文件所在路径（相对）

}

// 写入图片文件
func WriteImage(GUID string, fileStream []byte) (err error) {

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
func getStaticRoot() (rootPath string) {

}