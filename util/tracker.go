// Token-Static-Center
// 客户端信息获取模块
// 获取客户端的IP地址、请求URL等信息
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package util

import (
	"net"
	"net/http"
	"strings"
)

// 获取当前请求者的IP地址
func GetRequestIP(r *http.Request) (requestIP string) {
	ipType, err := GetConfig("Global", "Log", "IPType")

	if err != nil {
		return "配置文件读取失败"
	}

	var ipPort string

	switch ipType {
	case "native":
		ipPort, _, _ = net.SplitHostPort(r.RemoteAddr)
		break
	case "real-ip":
		ipPort = r.Header.Get("X-Real-IP")
		break
	case "x-forwarded-for":
		ipPort = r.Header.Get("X-Forwarded-For")
		break
	case "cloud-flare":
		ipPort = r.Header.Get("CF-Connecting-IP")
		break
	default:
		return "配置文件中日志采集IP类别错误"
	}

	// Fallback
	if ipPort == "" {
		ipPort, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return ipPort
}

// 获取所请求的链接地址
func GetRequestURI(r *http.Request) (uri string) {
	// 该业务只提供HTTP服务，但是因为反向代理对外是HTTPS，因此使用HTTPS的scheme
	scheme := "https://"
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}
