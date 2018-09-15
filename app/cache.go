// Token-Static-Center
// 缓存模块
// 负责缓存的处理与管理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"bytes"
	"errors"
	"github.com/LuRenJiasWorld/Token-Static-Center/db"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"os"
	"path/filepath"
	"strconv"
)

// 缓存获取
func CacheFetchHandler(cacheFileName string) (cacheStatus bool, cacheFileStream []byte, fileSizeByte int, err error) {
	// 垃圾回收
	CacheGCHandler()

	// 获取缓存路径
	cachePath, err := getCachePath(cacheFileName)
	if err != nil {
		return false, nil, -1, errors.New("尝试获取缓存时遭遇致命错误：" + err.Error())

	}

	// 尝试读取文件
	cacheFileStream, err = readFile(cachePath)
	// 处理错误
	// 文件不存在，或者文件为文件夹
	if err != nil {
		return false, nil, -1, nil
	}

	// 获取文件体积
	cacheFileSize := bytes.Count(cacheFileStream, nil)

	return true, cacheFileStream, cacheFileSize, nil
}

// 缓存写入
func CacheWriteHandler(cacheStream []byte, cacheFileName string) (cacheLocation string, fileSizeByte int, err error) {
	// 垃圾回收
	CacheGCHandler()

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
	// 获取上次缓存垃圾回收时间（秒）
	intervalTimeSecond, err := db.ReadGC()
	if err != nil {
		util.ErrorLog("CacheGCHandler", "获取上次缓存垃圾回收时间失败，原因：" + err.Error(), "app->CacheGCHandler")
	}

	// 获取缓存配置
	cacheIntervalTimeHourInterface, err := util.GetConfig("Cache", "GCInterval")
	if err != nil {
		util.ErrorLog("CacheGCHandler", "获取配置文件中的缓存垃圾回收时间间隔时出现错误，原因：" + err.Error(), "app->CacheGCHandler")
	}
	cacheIntervalTimeHour := cacheIntervalTimeHourInterface.(int)

	cacheThresholdInterface, err := util.GetConfig("Cache", "GCThreshold")
	if err != nil {
		util.ErrorLog("CacheGCHandler", "获取配置文件中的缓存垃圾回收数量阈值时出现错误，原因：" + err.Error(), "app->CacheGCHandler")
	}
	cacheThreshold := cacheThresholdInterface.(int)

	// 获取缓存文件路径
	cachePath, err := getCachePath("")

	// 获取缓存文件列表
	cacheFileList, err := filepath.Glob(cachePath + "*")

	// 缓存文件数量
	cacheFileCount := len(cacheFileList)

	// 判断条件进行缓存清除
	if cacheFileCount >= cacheThreshold || intervalTimeSecond >= cacheIntervalTimeHour * 3600 {
		for i := range cacheFileList {
			err = os.Remove(cacheFileList[i])
			if err != nil {
				util.WarningLog("CacheGCHandler", "缓存文件垃圾回收时出现错误：文件" + cacheFileList[i] + "无法删除，错误信息为：" + err.Error(), "app->CacheGCHandler")
			}
		}

		// 日志记录
		util.OperationLog("CacheGCHandler", "缓存已清除，总共清除缓存数量：" + strconv.Itoa(cacheFileCount), "app->CacheGCHandler")

		// 更新垃圾回收数据库
		db.UpdateGC(cacheFileCount)
	}
}

// 获取缓存路径
// 将缓存目录与缓存文件名进行拼接获取缓存路径
func getCachePath(cacheFileName string) (path string, err error) {
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