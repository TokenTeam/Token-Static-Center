package core

import (
	"github.com/husobee/vestigo"
	"github.com/LuRenJiasWorld/Token-Static-Center/app"
	"github.com/justinas/alice"
)

// 初始化服务器
func NewServer() (*vestigo.Router) {
	router := vestigo.NewRouter()

	imageFileHandler := alice.New()

	router.Get("/", app.HomePage)
	router.Get("/image/:filename")


	return router
}