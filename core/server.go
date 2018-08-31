// Token-Static-Center
// 服务器模块
// 路由配置&中间件配置，用于初始化服务器
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package core

import (
	"github.com/husobee/vestigo"
	"github.com/LuRenJiasWorld/Token-Static-Center/app"
	"github.com/justinas/alice"
	"net/http"
	"github.com/LuRenJiasWorld/Token-Static-Center/util"
	"errors"
	"github.com/LuRenJiasWorld/Token-Static-Center/security"
)

// 初始化服务器
func NewServer() (vestigoRouter *vestigo.Router, err error) {
	router := vestigo.NewRouter()

	// 获取Debug模式状态
	debugStatus, err := util.GetConfig("Global", "Debug")

	if err != nil {
		util.ErrorLog("server", "获取Debug模式状态失败，请检查配置文件！", "server->NewServer")
		return nil, errors.New("服务器初始化失败！")
	}

	// 转换Interface到String
	debugStatus = debugStatus.(string)

	// 如果为Debug模式，则不检查安全性、不缓存
	var imageFileHandler alice.Chain

	// 根据Debug模式切换中间件模式
	switch debugStatus {
		case "on":
			imageFileHandler = alice.New()
			break
		case "off":
			imageFileHandler = alice.New(security.WhiteListFilter)
			break
		default:
			util.ErrorLog("server", "调试模式配置错误！请检查配置文件！", "server->debugStatus")
	}

	router.Get("/", app.HomePage)
	router.Get("/image/:filename", imageFileHandler.ThenFunc(app.HomePage).(http.HandlerFunc))

	return router, nil
}