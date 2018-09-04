// Token-Static-Center
// 安全模块
// 用于防盗链、防恶意访问等功能，通常作为中间件
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package security

import (
	"net/http"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"github.com/LuRenJiasWorld/Token-Static-Center/app"
)

// 过滤白名单以外的域名访问图片资源
func WhiteListFilter(next http.Handler) (http.Handler) {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		// 获取防盗链配置
		antiLeechStatus, err := util.GetConfig("Security", "AntiLeech", "Status")

		if err != nil {
			util.ErrorLog("security", "读取防盗链配置失败，错误信息为：" + err.Error(), "security->WhiteListFilter")
			return
		}

		// 转换Interface到String
		antiLeechStatus = antiLeechStatus.(string)

		// 如果防盗链开启
		if antiLeechStatus == "on" {
			// 获取防盗链白名单
			whiteListInterface, err := util.GetConfig("Security", "WhiteList")
			if err != nil {
				util.ErrorLog("security", "无法获取防盗链白名单，请检查防盗链配置！，错误信息：" + err.Error(), "security->WhiteListFilter")
				return
			}

			// 转换Interface到StringSlice
			whiteListArray := whiteListInterface.([]string)

			// 获取当前的HTTP referrer
			requestReferrer := r.Referer()

			// httpReferrer是否命中白名单
			hitFlag := false
			for i := 0; i < len(whiteListArray); i++ {
				if whiteListArray[i] == requestReferrer {
					hitFlag = true
					break
				}
			}

			// 获取防盗链警告配置
			antiLeechWarning, err := util.GetConfig("Security", "AntiLeech", "ShowWarning")

			if err != nil {
				util.ErrorLog("security", "读取防盗链警告配置失败，错误信息为：" + err.Error(), "security->WhiteListFilter")
			}

			// 转换Interface到String
			antiLeechWarning = antiLeechWarning.(string)

			// 如果未命中白名单，或者根本没有来源头，返回错误
			if hitFlag == false || requestReferrer == "" {
				switch antiLeechWarning {
					case "on":
						app.AntiLeechImage(w, r)
						break
					case "off":
						app.ErrorPage(w, r, 403, "WhiteListFilter", "触发反盗链机制")
						return
				}
				return
			}

		}

		// 传递到下一个中间件
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunction)
}


// 安全上传
func SecureUploadFilter() {

}
