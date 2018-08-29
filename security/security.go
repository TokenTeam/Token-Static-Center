package security

import (
	"net/http"
	"strings"
)

// 过滤白名单以外的IP访问


// 获取当前请求者的IP地址
func GetRequestIP(r *http.Request) (requestIP string) {
	// 兼容IPv6（八段地址+一个端口）
	requestAddrArray := make([]string, 9)

	// 按冒号分割（如果是IPv6则分割出九个元素）
	requestAddrArray = strings.Split(r.RemoteAddr, ":")

	// 最后返回的地址字符串
	requestAddrString := ""

	// 将最后一个端口号去除
	for i := 0; i < len(requestAddrArray) - 1; i++ {
		requestAddrString += requestAddrArray[i]
	}

	return requestAddrString
}