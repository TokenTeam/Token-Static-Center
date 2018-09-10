// Token-Static-Center
// 缓存模块
// 负责缓存的处理与管理
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

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
func CacheWriteHandler() () {

}

// 缓存垃圾回收
func CacheGCHandler() () {

}