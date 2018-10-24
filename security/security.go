// Token-Static-Center
// 安全模块
// 用于防盗链、防恶意访问等功能，通常作为中间件
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package security

import (
	"fmt"
	"github.com/TokenTeam/Token-Static-Center/app"
	"github.com/TokenTeam/Token-Static-Center/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 过滤白名单以外的域名访问图片资源
func WhiteListFilter(next http.Handler) http.Handler {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		// 获取防盗链配置
		antiLeechStatus, err := util.GetConfig("Security", "AntiLeech", "Status")

		if err != nil {
			util.ErrorLog("security", "读取防盗链配置失败，错误信息为：" + err.Error(), "security->WhiteListFilter")
			return
		}

		// 如果防盗链开启
		if antiLeechStatus.(string) == "on" {
			// 获取防盗链白名单
			whiteListInterface, err := util.GetConfig("Security", "WhiteList")
			if err != nil {
				util.ErrorLog("security", "无法获取防盗链白名单，请检查防盗链配置！，错误信息：" + err.Error(), "security->WhiteListFilter")
				return
			}

			// 转换Interface到StringSlice
			whiteListArray := whiteListInterface.([]string)

			// 获取当前的HTTP referrer
			requestReferer := r.Referer()

			// httpReferer是否命中白名单
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
			if hitFlag == false || requestReferer == "" {
				switch antiLeechWarning {
				case "on":
					app.AntiLeechImage(w, r)
					return
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
func SecureUploadFilter(next http.Handler) http.Handler {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) {
		// 获取安全上传配置
		secureUploadStatus, err := util.GetConfig("Security", "Token")

		if err != nil {
			util.ErrorLog("security", "读取安全上传配置失败，错误信息为：" + err.Error(), "security->SecureUploadFilter")
			return
		}

		// 如果安全上传为启动
		if secureUploadStatus.(string) == "on" {
			// 获取安全配置-AppCode
			appCodeInterface, err := util.GetConfig("Security", "AppCode")

			if err != nil {
				util.ErrorLog("security", "读取安全上传AppCode失败，错误信息为：" + err.Error(), "security->SecureUploadFilter")
				return
			}

			// 转换Interface到StringSlice
			appCodeStringSlice := appCodeInterface.([]string)

			// 解析请求
			// Step1. 获取形如/upload/blablabla-blabla.bla的URI
			requestParam := r.RequestURI
			// Step2. /upload/blablabla-blabla.bla => {"", "upload", "blablabla-blabla.bla"}
			requestParams := strings.Split(requestParam, "/")
			// Step3. blablabla-blabla.bla => {"blablabla-blabla", "bla"}
			requestParams = strings.Split(requestParams[2], ".")
			// Step4. blablabla-blabla => {"blablabla", "blabla"}
			requestParams = strings.Split(requestParams[0], "-")

			// appCode匹配状态
			appCodeHitStatus := false

			// 获取AccessToken
			accessToken := requestParams[0]
			// 获取Nonce
			nonce := requestParams[1]

			// 轮询已有AppCode，进行匹配计算
			for i := 0; i < len(appCodeStringSlice); i++ {
				// 转换字符为rune
				appCodeRune := []rune(appCodeStringSlice[i])

				// 检测是否满足64位长度，不满足，跳转到下一个AppCode
				if strings.Count(appCodeStringSlice[i], "") != (64 + 1) {
					continue
				}

				// 本地计算
				// 获取AppCode前32位
				appCodePrefix := string(appCodeRune[0:32])
				// 获取AppCode后32位
				appCodePostfix := string(appCodeRune[32:64])
				// 获取时间戳（去掉后四位）
				currentTimeStamp := time.Now().Unix()
				currentTimeStampRune := []rune(strconv.FormatInt(currentTimeStamp, 10))
				currentTimeStampString := string(currentTimeStampRune[0:6])
				// 获取Salt
				saltStringInterface, err := util.GetConfig("Security", "TokenSalt")
				if err != nil {
					util.ErrorLog("security", "读取Token验证SaltString失败，错误信息为：" + err.Error(), "security->SecureUploadFilter")
					return
				}
				saltString := saltStringInterface.(string)

				totalString := fmt.Sprintf("%s%s%s%s%s", appCodePrefix, currentTimeStampString, appCodePostfix, nonce, saltString)

				totalStringByte := []byte(totalString)

				currentAccessToken := util.GetMD5Hash(totalStringByte)

				if currentAccessToken == accessToken {
					// 输出当前AppCode，用于存储到app模块，便于存储当前业务信息到数据库
					app.CurrentUploadAppCode = appCodeStringSlice[i]
					appCodeHitStatus = true
					break
				}
			}

			// 如果都不匹配
			if appCodeHitStatus == false {
				app.ErrorPage(w, r, 403, "SecureUploadFilter", "存在非法授权，AccessToken校验失败，所发送的AccessToken为"+accessToken+"，nonce为"+nonce)
				return
			}
		}

		// 传递到下一个中间件
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunction)
}
