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
			// 获取文件路由：文件获取
			imageFileHandler = alice.New()
			// 上传文件路由：直接上传
			break
		case "off":
			imageFileHandler = alice.New(security.WhiteListFilter)
			// 上传文件路由：Token校验->文件上传
			// accessToken校验：检查 md5(AppCode前32位+时间戳去掉最后三位+AppCode后32位+Nonce+"token123")是否通过校验
			// 校验规则：客户端传递accessToken、Nonce给服务器，服务器对所有AppCode进行检查->根据服务器端时间戳计算accessToken->与客户端accessToken进行对比
			break
		default:
			util.ErrorLog("server", "调试模式配置错误！请检查配置文件！", "server->debugStatus")
	}

	router.Get("/", app.HomePage)

	router.Get("/image/:filename", imageFileHandler.ThenFunc(app.ImageFetchHandler).(http.HandlerFunc))

	router.Post("/upload/:parameter", imageFileHandler.ThenFunc(app.ImageUploadHandler).(http.HandlerFunc))

	// 防止upload方法被误访问
	router.Get("/upload/:parameter", func(w http.ResponseWriter, r *http.Request) {
		app.ErrorPage(w, r, 403, "server->NewServer", "通过GET方法访问了只允许POST访问的上传接口")
	})

	return router, nil
}