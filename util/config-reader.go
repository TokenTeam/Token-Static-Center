// Token-Static-Center
// 配置文件模块
// 解析配置文件&读取配置文件
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"errors"
)


// 全局配置数据（固化存储）
var configData = Config{}
var configReadStatus = false

// 配置文件模板
type Config struct {
	Global struct {
		Debug string				`yaml:"debug"`
		Port uint32					`yaml:"port"`
		StorageDir string			`yaml:"storage-dir"`
		Log struct {
			LogDir string			`yaml:"log-dir"`
			LogAccess string		`yaml:"log-access"`
			LogOperation string		`yaml:"log-operation"`
			LogWarning string		`yaml:"log-warning"`
			LogError string			`yaml:"log-error"`
		}							`yaml:"log"`
		Db struct {
			DbType string			`yaml:"db-type"`
			DbResource string		`yaml:"db-resource"`
		}							`yaml:"db"`
	} 								`yaml:"global"`

	Image struct {
		MaxWidth string				`yaml:"max-width"`
		MaxHeight string			`yaml:"max-height"`
		UploadableFileType []string	`yaml:"uploadable-file-type"`
		AccessableFileType []string `yaml:"accessable-file-type"`
		StorageFileType string		`yaml:"static-file-type"`
		JpegCompressLevel uint32	`yaml:"jpeg-compress-level"`
		MaxImageFileSize uint32		`yaml:"max-image-file-size"`
	}

	Security struct {
		WhiteList []string			`yaml:"white-list"`
		AppCode []string			`yaml:"app-code"`
		Token string				`yaml:"token"`
		TokenSalt string			`yaml:"token-salt"`
		AntiLeech struct {
			Status string			`yaml:"status"`
			ShowWarning string		`yaml:"show-warning"`
		}							`yaml:"anti-leech"`
	}								`yaml:"security"`

	Cache struct {
		Status string				`yaml:"status"`
		CacheDir string				`yaml:"cache-dir"`
		GCInterval string			`yaml:"gc-interval"`
		GCThreshold uint32			`yaml:"gc-threshold"`
	}
}


// 初始化配置项
func ReadConfig(configFile string) (err error) {
	file, err := ioutil.ReadFile(configFile)

	if err != nil || configFile == "" {
		return errors.New("配置文件打开错误！请检查文件是否存在、是否通过--config方式进行调用、以及文件权限是否正确！")
	}

	configData = Config{}
	err = yaml.Unmarshal(file, &configData)

	if err != nil {
		return errors.New("配置文件解析错误！请参考示例配置文件进行配置！")
	}

	configReadStatus = true

	return nil
}


// 获取配置项
func GetConfig(config ...string) (result interface{}, err error) {
	// 检查是否已初始化配置项
	if configReadStatus == false {
		return nil, errors.New("配置文件尚未初始化，无法获取配置项！")
	}

	if len(config) == 0 {
		return Struct2Map(configData)
	} else {
		configItem, err := Struct2Map(configData)
		if err != nil {
			return nil, errors.New("配置项转换失败！无法调用配置！具体原因：" + err.Error())
		}
		for i := 0; i < len(config); i ++ {
			result = configItem[config[i]]
			if i < len(config) - 1 {
				configItem, err = Struct2Map(result)
				if err != nil {
					return nil, errors.New("配置项转换失败！无法调用配置！具体原因：" + err.Error())
				}
			} else {
				break
			}
		}
	}

	return result, nil
}
