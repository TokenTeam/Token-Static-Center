// Token-Static-Center
// 数据库模块
// 封装静态资源引擎常用的数据库操作
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package db

import (
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"errors"
	"time"
	"strconv"
	"fmt"
)

// 数据库结构（以MySQL为例，SQLite具体参考sqlite.go->checkDBStructureSQLite方法）
// token_static_center数据库（具体数据库名称由配置文件决定）SQLite数据库情况下无索引
// |- image_info数据表
// |  |- guid varchar 128 primary		图片资源的GUID
// |  |- year bigint index				图片资源存放的年份（同时也是存储的一级目录名）
// |  |- month bigint index				图片资源存放的月份（同时也是存储的二级目录名）
// |  |- file_size_byte bigint			图片资源大小（byte）
// |  |- file_storage_format varchar 32	图片存储格式（存储在服务器上的格式，由当前配置文件决定）
// |  |- upload_time datetime			图片存储时间
// |  |- app_code varchar 64 index		所属业务AppCode
// |  |- md5 varchar 32 index			该文件MD5（用于数据校验）
// |  |- download_count bigint			总下载次数
// |- gc_log数据表
// |  |- id int ai primary				垃圾收集ID
// |  |- collection_time timestamp index上次收集时间(时间戳)
// |  |- garbage_count bigint 			此次收集垃圾数量
// |- image_statistics数据表
// |  |- date varchar 10 primary		统计时间（yyyy-mm-dd）
// |  |- upload_count bigint			当日上传计数
// |  |- download_count bigint			当日下载计数
// |  |- upload_size_byte bigint		当日上传大小（Byte）
// |  |- download_size_byte bigint		当日下载大小（Byte）


// 写入图片数据
func WriteImageDB(GUID string, fileSize uint64, fileFormat string, AppCode string, md5 string) (err error) {
	t := time.Now()
	year := t.Year()
	month := int(t.Month())

	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return errors.New("写入图片数据过程中校验数据库类型失败，原因：" + err.Error())
	}

	switch dbType {
		case "mysql":
			writeData := []string{GUID, strconv.Itoa(year), strconv.Itoa(month), strconv.FormatUint(fileSize, 10), fileFormat, "", AppCode, md5, "0"}
			err = insertMySQL("image_info", writeData)
			break
		case "sqlite":
			writeData := []string{GUID, strconv.Itoa(year), strconv.Itoa(month), strconv.FormatUint(fileSize, 10), fileFormat, time.Now().Format("2006-01-02 15:04:05"), AppCode, md5, "0"}
			err = insertSQLite("image_info", writeData)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err != nil {
		return errors.New("写入图片数据过程中写入到数据库失败，原因：" + err.Error())
	}

	return nil
}

// 读取图片数据
func ReadImageDB(GUID string) (year int, month int, md5 string, format string, err error) {
	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return -1, -1, "", "",  errors.New("读取图片数据过程中校验数据库类型失败，原因：" + err.Error())
	}

	selectData := map[string]string{"guid":GUID}

	// 返回回来的数据集
	imageData := map[int]map[string]string{}

	switch dbType {
		case "mysql":
			imageData, err = selectMySQL("image_info", selectData)
			break
		case "sqlite":
			imageData, err = selectSQLite("image_info", selectData)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err != nil {
		return -1, -1, "", "", errors.New("读取图片数据过程中读取数据集失败，原因：" + err.Error())
	}

	// 此处只需要一条（唯一一条）记录
	year, _ = strconv.Atoi(imageData[0]["year"])
	month, _ = strconv.Atoi(imageData[0]["month"])
	format = imageData[0]["file_storage_format"]
	md5 = imageData[0]["md5"]

	return year, month, md5, format, nil
}

// 写入上次GC时间
func UpdateGC(garbageCount int) (err error) {
	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return errors.New("写入垃圾收集记录中校验数据库类型失败，原因：" + err.Error())
	}

	switch dbType {
		case "mysql":
			writeData := []string{"", "", strconv.Itoa(garbageCount)}
			err = insertMySQL("gc_log", writeData)
			break
		case "sqlite":
			writeData := []string{"", time.Now().Format("2006-01-02 15:04:05"), strconv.Itoa(garbageCount)}
			err = insertSQLite("gc_log", writeData)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err != nil {
		return errors.New("写入垃圾收集记录中写入到数据库失败，原因：" + err.Error())
	}

	return nil
}

// 读取上次GC时间间隔（秒）
func ReadGC() (intervalTime int, err error) {
	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return -1, errors.New("读取垃圾收集记录中校验数据库类型失败，原因：" + err.Error())
	}

	selectData := map[string]string{}

	// 返回回来的数据集
	gcData := map[int]map[string]string{}

	switch dbType {
		case "mysql":
			gcData, err = selectMySQL("gc_log", selectData)
			break
		case "sqlite":
			gcData, err = selectSQLite("gc_log", selectData)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err != nil {
		return -1, errors.New("读取垃圾收集记录中读取数据集失败，原因：" + err.Error())
	}

	// 返回所有GC日志
	// 获取最新一条的GC日志
	lastIndex := len(gcData) - 1
	lastGCTime := gcData[lastIndex]["collection_time"]

	// 获取上一次垃圾回收时间戳
	// time.ParseInLocation的Layout参数是魔法数字，不要修改
	lastGCTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", lastGCTime, time.Local)
	lastGCTimeStampInt64 := lastGCTimeStamp.Unix()

	// 获取当前时间戳
	currentTime := time.Now()
	currentTimeStampInt64 := currentTime.Unix()

	timeInterval := int(currentTimeStampInt64 - lastGCTimeStampInt64)
	return timeInterval, nil
}

// 上传次数增加计数
func UploadCounter(upload_size_byte int) (err error) {
	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return errors.New("上传计数自增中校验数据库类型失败，原因：" + err.Error())
	}

	// t.Format的layout参数是魔法数字，不要修改
	t := time.Now()
	date := t.Format("2006-01-02")

	queryData := map[string]string{"date":date}

	var (
		err1 error
		err2 error
	)

	switch dbType {
		case "mysql":
			// 先检查是否存在当日记录（可能会产生错误，无视之，有错误说明该记录已存在，否则写入新纪录）
			execMySQL("INSERT INTO image_statistics VALUES ('" + date + "', 0, 0, 0, 0)")
			err1 = updateMySQL("image_statistics", "upload_count", 1, queryData)
			err2 = updateMySQL("image_statistics", "upload_size_byte", upload_size_byte, queryData)
			break
		case "sqlite":
			execSQLite("INSERT INTO image_statistics VALUES ('" + date + "', 0, 0, 0, 0)")
			err1 = updateSQLite("image_statistics", "upload_count", 1, queryData)
			err2 = updateSQLite("image_statistics", "upload_size_byte", upload_size_byte, queryData)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err1 != nil || err2 != nil {
		return errors.New("上传计数自增中更新数据失败，原因：" + err1.Error() + err2.Error())
	}

	return nil
}

// 下载次数增加计数
func DownloadCounter(guid string, download_size_byte int) (err error) {
	// 获取数据库类型
	dbType, err := getDBType()

	if err != nil {
		return errors.New("上传计数自增中校验数据库类型失败，原因：" + err.Error())
	}

	// t.Format的layout参数是魔法数字，不要修改
	t := time.Now()
	date := t.Format("2006-01-02")

	queryData := map[string]string{"date":date}

	var (
		err1 error		// 当日总下载计数
		err2 error		// 当日总下载数据量
		err3 error		// 对应文件下载数据量
	)

	fileQuery := map[string]string{"guid":guid}

	switch dbType {
		case "mysql":
			// 先检查是否存在当日记录（可能会产生错误，无视之，有错误说明该记录已存在，否则写入新纪录）
			execMySQL("INSERT INTO image_statistics VALUES ('" + date + "', 0, 0, 0, 0)")
			err1 = updateMySQL("image_statistics", "download_count", 1, queryData)
			err2 = updateMySQL("image_statistics", "download_size_byte", download_size_byte, queryData)
			err3 = updateMySQL("image_info", "download_count", 1, fileQuery)
			break
		case "sqlite":
			execSQLite("INSERT INTO image_statistics VALUES ('" + date + "', 0, 0, 0, 0)")
			err2 = updateSQLite("image_statistics", "download_count", 1, queryData)
			err2 = updateSQLite("image_statistics", "download_size_byte", download_size_byte, queryData)
			err3 = updateSQLite("image_info", "download_count", 1, fileQuery)
			break
		// Default类型此前已经判断过，不需要重复判断
	}

	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("上传计数自增中更新数据失败，原因：" + err1.Error() + err2.Error())
	}

	return nil
}

// 获取数据库类型
func getDBType() (dbType string, err error) {
	// 获取数据库类型配置文件
	dbTypeInterface, err := util.GetConfig("Global", "Db", "DbType")

	if err != nil {
		return "", errors.New("数据库类型获取失败，原因：" + err.Error())
	}

	// 转换Interface到String
	dbTypeString := dbTypeInterface.(string)

	switch dbTypeString {
		case "mysql":
			return "mysql", nil
		case "sqlite":
			return "sqlite", nil
		default:
			return "", errors.New("获取数据库类型时出现致命错误：不受支持的数据库类型：" + dbTypeString)
	}
}

func Test() {
	//insert := []string{"2018-09-01", "23", "312231", "3123313", "12231324543243"}
	//fmt.Println(insertMySQL("image_statistics", insert))

	//update := map[string]string{"date":"2018-09-01"}
	//fmt.Println(updateMySQL("image_statistics", "upload_count", 10000, update))

	//query := map[string]string{"date":"2018-09-01"}
	//fmt.Println(selectMySQL("image_statistics", query))

	//fmt.Println(execMySQL("SELECT * FROM image_statistics"))
}