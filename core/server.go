// Token-Static-Center
// 服务器模块
// 路由配置&中间件配置，用于初始化服务器
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package core

import (
	"github.com/husobee/vestigo"
	"github.com/LuRenJiasWorld/Token-Static-Center/app"
	"github.com/justinas/alice"
	"github.com/LuRenJiasWorld/Token-Static-Center/security"
	"net/http"
)

// 初始化服务器
func NewServer() (*vestigo.Router) {
	router := vestigo.NewRouter()

	imageFileHandler := alice.New(security.WhiteListFilter)

	router.Get("/", app.HomePage)
	router.Get("/image/:filename", imageFileHandler.ThenFunc(app.HomePage).(http.HandlerFunc))


	return router
}