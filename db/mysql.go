// Token-Static-Center
// 数据库模块-MySQL兼容层
// 针对MySQL进行兼容
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package db

import (
	"database/sql"
	"errors"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

// 插入数据
func insertMySQL(table string, data []string) (err error) {
	// 获取数据库连接句柄
	dbHandle, err := connectMySQL()
	defer dbHandle.Close()
	if err != nil {
		return errors.New("插入数据时，连接数据库失败，原因：" + err.Error())
	}

	// 检查数据库结构
	err = checkDBStructureMySQL(dbHandle)
	if err != nil {
		return errors.New("插入数据时，检查数据库结构失败，原因：" + err.Error())
	}

	// 处理插入数据
	insertValueString := "("
	for i := range data {
		// 利用Atoi的报错机制，将字符串和数字区分出来
		_, err := strconv.Atoi(data[i])
		if err != nil {
			if len(data[i]) == 0 {
				insertValueString = insertValueString + "null"
			} else {
				insertValueString = insertValueString + "\"" + data[i] + "\""
			}
		} else {
			insertValueString = insertValueString + data[i]
		}


		// 如果i不是最后一个索引，尾部添加逗号
		if i != len(data) - 1 {
			insertValueString = insertValueString + ","
		}
	}
	insertValueString = insertValueString + ")"

	// 预格式化插入操作
	insertObj, err := dbHandle.Prepare("INSERT INTO " + table + " VALUES "+ insertValueString)
	if err != nil {
		return errors.New("插入数据时，预格式化语句失败，原因：" + err.Error())
	}

	// 执行格式化语句
	insertResult, err := insertObj.Exec()

	if err != nil {
		return errors.New("插入数据失败，原因：" + err.Error())
	}

	// 检查是否已插入
	insertRows, err := insertResult.RowsAffected()
	if err != nil || insertRows != 1 {
		return errors.New("插入数据后，校验失败，相关原因：" + err.Error())
	}

	return nil
}

// 读取数据
// - 支持AND，不支持OR
// - 只支持=，不支持其他类别
// 如果要使用更高级的操作，建议执行execMySQL方法
func selectMySQL(table string, query map[string]string) (queryResults map[int]map[string]string, err error) {
	// 获取数据库连接句柄
	dbHandle, err := connectMySQL()
	defer dbHandle.Close()
	if err != nil {
		return nil, errors.New("执行查询语句时，连接数据库失败，原因：" + err.Error())
	}

	// 检查数据库结构
	err = checkDBStructureMySQL(dbHandle)
	if err != nil {
		return nil, errors.New("插入数据时，检查数据库结构失败，原因：" + err.Error())
	}

	// 处理查询语句
	// 将map转换为map[1]key=map[1]value AND ...... AND true 格式
	queryWhereString := ""
	for key, value := range query {
		queryWhereString = queryWhereString + key + "=" + "'" + value + "'"
		queryWhereString = queryWhereString + " AND "
	}
	// 让最后一个AND有意义，不引起语法错误
	queryWhereString = queryWhereString + "true"

	// 查询
	queryRows, err := dbHandle.Query("SELECT * FROM " + table + " WHERE " + queryWhereString)

	if err != nil {
		return nil, errors.New("执行查询语句时，执行失败，原因：" + err.Error())
	}

	// 存储查询结果
	// 此处使用『骚操作』，解决了Scan方法必须为定长的限制
	columns, _ := queryRows.Columns()
	tempValues := make([][]byte, len(columns))
	tempValuesInterface := make([]interface{}, len(columns))
	for i := range tempValues {
		tempValuesInterface[i] = &tempValues[i]
	}
	queryResults = make(map[int]map[string]string)
	i := 0
	for queryRows.Next() {
		// 将当前指针所在行的数据写入可变长度的Interface中
		err := queryRows.Scan(tempValuesInterface...)
		if err != nil {
			return nil, errors.New("执行查询语句时，抓取第" + strconv.Itoa(i) + "项的值出现错误：" + err.Error())
		}
		// 当前行数据
		rowData := make(map[string]string)
		for key, value := range tempValues {
			keyName := columns[key]
			rowData[keyName] = string(value)
		}
		queryResults[i] = rowData
		i++
	}
	return queryResults, nil
}

// 更新数据
// - 查询条件支持AND，不支持OR
// - 只能够对数值进行自增处理
// 如果要使用更高级的操作，建议执行execMySQL方法
func updateMySQL(table string, updateItem string, updateValue int, query map[string]string) (err error) {
	// 获取数据库连接句柄
	dbHandle, err := connectMySQL()
	defer dbHandle.Close()
	if err != nil {
		return errors.New("执行更新语句时，连接数据库失败，原因：" + err.Error())
	}

	// 检查数据库结构
	err = checkDBStructureMySQL(dbHandle)
	if err != nil {
		return errors.New("插入数据时，检查数据库结构失败，原因：" + err.Error())
	}

	// 处理查询语句
	// 将map转换为map[1]key=map[1]value AND ...... AND true 格式
	queryWhereString := ""
	for key, value := range query {
		queryWhereString = queryWhereString + key + "=" + "'" + value + "'"
		queryWhereString = queryWhereString + " AND "
	}
	// 让最后一个AND有意义，不引起语法错误
	queryWhereString = queryWhereString + "true"

	// 预格式化更新操作
	updateObj, err := dbHandle.Prepare("UPDATE "+ table + " SET " + updateItem + " = " + updateItem + " + " + strconv.Itoa(updateValue) + " WHERE " + queryWhereString)

	if err != nil {
		return errors.New("插入数据时，预格式化语句失败，原因：" + err.Error())
	}

	updateResult, err := updateObj.Exec()

	updateRows, err := updateResult.RowsAffected()

	if err != nil || updateRows != 1 {
		return errors.New("更新数据后，校验失败，相关原因：" + err.Error())
	}

	return nil
}

// 执行任意SQL语句
// 部分代码类似于selectMySQL
func execMySQL(query string) (execResults map[int]map[string]string, err error) {
	// 获取数据库连接句柄
	dbHandle, err := connectMySQL()
	defer dbHandle.Close()
	if err != nil {
		return nil, errors.New("执行语句时，连接数据库失败，原因：" + err.Error())
	}

	// 检查数据库结构
	err = checkDBStructureMySQL(dbHandle)
	if err != nil {
		return nil, errors.New("执行语句时，检查数据库结构失败，原因：" + err.Error())
	}

	// 查询
	queryRows, err := dbHandle.Query(query)

	if err != nil {
		return nil, errors.New("执行语句时，执行失败，原因：" + err.Error())
	}

	// 存储查询结果
	// 此处使用『骚操作』，解决了Scan方法必须为定长的限制
	columns, _ := queryRows.Columns()
	tempValues := make([][]byte, len(columns))
	tempValuesInterface := make([]interface{}, len(columns))
	for i := range tempValues {
		tempValuesInterface[i] = &tempValues[i]
	}
	queryResults := make(map[int]map[string]string)
	i := 0
	for queryRows.Next() {
		// 将当前指针所在行的数据写入可变长度的Interface中
		err := queryRows.Scan(tempValuesInterface...)
		if err != nil {
			return nil, errors.New("执行语句时，抓取第" + strconv.Itoa(i) + "项的值出现错误：" + err.Error())
		}
		// 当前行数据
		rowData := make(map[string]string)
		for key, value := range tempValues {
			keyName := columns[key]
			rowData[keyName] = string(value)
		}
		queryResults[i] = rowData
		i++
	}
	return queryResults, nil
}

// 检查数据库结构，如果结构不存在，新建数据表
func checkDBStructureMySQL(dbHandle *sql.DB) (err error) {
	// 新建image_info表
	_, err1 := dbHandle.Exec("CREATE TABLE IF NOT EXISTS `image_info` ( `guid` VARCHAR(128) NOT NULL , `year` BIGINT NOT NULL , `month` BIGINT NOT NULL , `file_size_byte` BIGINT NOT NULL , `file_storage_format` VARCHAR(32) NOT NULL , `upload_time` TIMESTAMP NOT NULL , `app_code` VARCHAR(64) NOT NULL , `md5` VARCHAR(32) NOT NULL , `download_count` BIGINT NOT NULL , PRIMARY KEY (`guid`), INDEX (`year`), INDEX (`month`), INDEX (`app_code`), INDEX (`md5`)) ENGINE = MyISAM;")

	// 新建gc_log表
	_, err2 := dbHandle.Exec("CREATE TABLE IF NOT EXISTS `gc_log` ( `id` INT NOT NULL AUTO_INCREMENT , `collection_time` TIMESTAMP NOT NULL , `garbage_count` BIGINT NOT NULL , PRIMARY KEY (`id`), INDEX (`collection_time`)) ENGINE = MyISAM;")

	// 新建image_statistics表
	_, err3 := dbHandle.Exec("CREATE TABLE IF NOT EXISTS `image_statistics` ( `date` VARCHAR(10) NOT NULL , `upload_count` BIGINT NOT NULL , `download_count` BIGINT NOT NULL , `upload_size_byte` BIGINT NOT NULL , `download_size_byte` BIGINT NOT NULL , PRIMARY KEY (`date`)) ENGINE = MyISAM;")

	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("新建数据表时出现错误！" + err1.Error() + err2.Error() + err3.Error())
	}

	return nil
}

// 连接数据库
func connectMySQL() (db *sql.DB, err error) {
	mysqlConfig, err := getConfigMySQL()
	if err != nil {
		return nil, errors.New("数据库连接失败，原因：" + err.Error())
	}

	// 启动连接
	db, err = sql.Open("mysql", mysqlConfig)
	if err != nil {
		return nil, errors.New("数据库连接失败，原因：" + err.Error())
	}

	return db, nil
}

// 读取数据库配置
func getConfigMySQL() (config string, err error) {
	// 读取配置文件
	dbConfigInterface, err := util.GetConfig("Global", "Db", "DbResource")

	if err != nil {
		return "", errors.New("数据库配置读取失败，原因：" + err.Error())
	}

	// 转换Interface到String
	dbConfigString := dbConfigInterface.(string)

	return dbConfigString, nil
}