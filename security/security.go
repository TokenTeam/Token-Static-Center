// Token-Static-Center
// 安全模块
// 用于防盗链、防恶意访问等功能，通常作为中间件
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package security

import (
	"net/http"
	"fmt"
)

// 过滤白名单以外的IP访问
func WhiteListFilter(next http.Handler) (http.Handler) {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Referer())
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunction)
}
