// Token-Static-Center
// 配置文件转换模块
// 便于util.GetConfig实现配置文件的嵌套读取
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"errors"
	"reflect"
)

// 将结构体转换为Map
func Struct2Map(struc interface{}) (map[string]interface{}, error) {

	returnMap := make(map[string]interface{})

	sType := getStructType(struc)

	if sType.Kind() != reflect.Struct {
		return returnMap, errors.New("传入变量既不是结构体也不是结构体指针！")
	}

	for i := 0; i < sType.NumField(); i++ {
		structFieldName := sType.Field(i).Name
		structVal := reflect.ValueOf(struc)
		returnMap[structFieldName] = structVal.FieldByName(structFieldName).Interface()
	}

	return returnMap, nil
}

// 检查所输入变量的类型
func getStructType(struc interface{}) (reflect.Type) {
	sType := reflect.TypeOf(struc)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	return sType
}
