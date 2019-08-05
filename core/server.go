// Token-Static-Center
// 服务器模块
// 路由配置&中间件配置，用于初始化服务器
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package core

import (
	"errors"
	"github.com/TokenTeam/Token-Static-Center/app"
	"github.com/TokenTeam/Token-Static-Center/security"
	"github.com/TokenTeam/Token-Static-Center/util"
	"github.com/husobee/vestigo"
	"github.com/justinas/alice"
	"net/http"
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
	var (
		imageDownloadHandler alice.Chain
		imageUploadHandler alice.Chain
	)

	// 根据Debug模式切换中间件模式
	switch debugStatus {
		case "on":
			// 获取文件路由：直接获取
			imageDownloadHandler = alice.New()
			// 上传文件路由：直接上传
			imageUploadHandler = alice.New()
			break
		case "off":
			// 获取文件路由：白名单校验->文件获取
			imageDownloadHandler = alice.New(security.WhiteListFilter)
			// 上传文件路由：白名单校验->Token校验->文件上传
			imageUploadHandler = alice.New(security.WhiteListFilter, security.SecureUploadFilter)
			break
		default:
			util.ErrorLog("server", "调试模式配置错误！请检查配置文件！", "server->debugStatus")
			return
	}

	router.Get("/", app.HomePage)

	router.Get("/image/:filename", imageDownloadHandler.ThenFunc(app.ImageFetchHandler).(http.HandlerFunc))

	router.Post("/upload/:parameter", imageUploadHandler.ThenFunc(app.ImageUploadHandler).(http.HandlerFunc))

	// 防止upload方法被误访问
	router.Get("/upload/:parameter", func(w http.ResponseWriter, r *http.Request) {
		app.ErrorPage(w, r, 403, "server->NewServer", "通过GET方法访问了只允许POST访问的上传接口")
	})

	return router, nil
}