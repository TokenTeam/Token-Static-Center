package security

import (
	"net/http"
	"fmt"
)

// 过滤白名单以外的IP访问
func WhiteListFilter(r *http.Request) {
	fmt.Println()
}
