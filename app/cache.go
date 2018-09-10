// Token-Static-Center
// 缓存模块
// 负责缓存的处理与管理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"errors"
	"bytes"
)

// 缓存获取
func CacheFetchHandler(cacheFileName string) (cacheStatus bool, cacheFileStream []byte, fileSizeByte int, err error) {
	// 获取缓存路径
	cachePath, err := getCachePath(cacheFileName)
	if err != nil {
		return false, nil, -1, errors.New("尝试获取缓存时遭遇致命错误：" + err.Error())

	}

	// 尝试读取文件
	cacheFileStream, err = readFile(cachePath)
	// 处理错误
	// 文件不存在，或者文件为文件夹（不可能为文件夹，此前已经过滤过了）
	if err != nil {
		return false, nil, -1, nil
	}

	// 获取文件体积
	cacheFileSize := bytes.Count(cacheFileStream, nil)

	return true, cacheFileStream, cacheFileSize, nil
}

// 缓存写入
func CacheWriteHandler(cacheStream []byte, cacheFileName string) (cacheLocation string, fileSizeByte int, err error) {
	// 获取缓存路径
	cachePath, err := getCachePath(cacheFileName)
	if err != nil {
		return "", -1, errors.New("写入缓存时遭遇致命错误：" + err.Error())
	}

	// 写入文件
	err = writeFile(cachePath, cacheStream)
	if err != nil {
		return "", -1, errors.New("图片资源缓存时存储文件过程中遭遇致命错误：" + err.Error())
	}

	// 计算文件大小
	fileSizeByte = bytes.Count(cacheStream, nil)

	return cachePath, fileSizeByte, nil
}

// 缓存垃圾回收
func CacheGCHandler() () {

}

// 获取缓存路径
// 将缓存目录与缓存文件名进行拼接获取缓存路径
func getCachePath(cacheFileName string) (path string, err error) {
	// 检查文件名是否为空
	if cacheFileName == "" {
		return "", errors.New("图片资源缓存时获取缓存目录配置过程中遭遇致命错误：所请求的文件名为空！")
	}

	// 获取缓存目录配置文件
	cachePathInterface, _ := util.GetConfig("Cache", "CacheDir")

	// 转换Interface到String
	cachePathString := cachePathInterface.(string)
	if err != nil {
		return "", errors.New("图片资源缓存时获取缓存目录配置过程中遭遇致命错误：" + err.Error())
	}

	// 拼接最终缓存目录
	cachePath := cachePathString + cacheFileName

	return cachePath, nil
}