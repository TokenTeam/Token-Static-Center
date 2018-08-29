// Token-Static-Center
// 客户端信息获取模块
// 获取客户端的IP地址、请求URL等信息
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"net/http"
	"strings"
)

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

// 获取所请求的链接地址
func GetRequestURI(r *http.Request) (uri string) {
	// 该业务只提供HTTP服务，但是因为反向代理对外是HTTPS，因此使用HTTPS的scheme
	scheme := "https://"
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}
