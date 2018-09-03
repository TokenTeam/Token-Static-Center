// Token-Static-Center
// 哈希校验模块
// 校验文件（文件流）的哈希值
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"crypto"
	"encoding/hex"
)

// 计算文件流的md5哈希值
func GetMD5Hash(fileStream []byte) (md5Hash string) {
	fileHashObj := crypto.MD5.New()
	fileHashObj.Write(fileStream)
	fileHash := hex.EncodeToString(fileHashObj.Sum(nil))
	return fileHash
}