package main

import (
	"github.com/TokenTeam/Token-Static-Center/util"
	"flag"
	"log"
	"os"
	"strconv"
	"github.com/TokenTeam/Token-Static-Center/core"
	_ "github.com/mkevac/debugcharts" // 可选，添加后可以查看几个实时图表数据
	_ "net/http/pprof" // 必须，引入 pprof 模块
	"net/http"
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {

	if os.Getenv("DEBUG") == "true" {
		go func() {
			// terminal: $ go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
			// web:
			// 1、http://localhost:8081/ui
			// 2、http://localhost:6060/debug/charts
			// 3、http://localhost:6060/debug/pprof
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	// 初始化
	imagick.Initialize()
	// 延迟执行
	defer imagick.Terminate()

	// Step 1. 从命令行参数获取配置文件路径
	configFilePath := flag.String("config", "", "配置文件路径，必须指定，无默认值！")
	flag.Parse()

	if *configFilePath == "" {
		fmt.Println("请指定配置文件位置，格式为[--config 配置文件路径.yaml！]")
		return
	}

	// Step 2. 加载配置文件
	err := util.ReadConfig(*configFilePath)

	if err != nil {
		util.ErrorLog("main", err.Error(), "main->ReadConfig")
		return
	}

	util.OperationLog("main", "配置文件加载成功，路径为" + *configFilePath, "main->ReadConfig")

	// Step 3. 从配置文件读取监听端口
	// 注意，此处defaultPort返回回来的是一个Interface
	defaultPortInterface, err := util.GetConfig("Global", "Port")

	// 将Interface转换成int
	listenPort := defaultPortInterface.(uint32)

	if err != nil {
		util.ErrorLog("main", "无法读取默认配置！" + err.Error(), "main->GetConfig")
		return
	}

	util.OperationLog("main", "使用配置文件中的端口" + strconv.Itoa(int(listenPort)), "main->listenPort")

	// 输出欢迎信息
	util.AccessLog("main", "欢迎使用 Token-Static-Center v1.10, 初始化完成，开始接受外部请求，启动服务", "main")


	// Step 4. 开始监听
	server, err := core.NewServer()

	if err != nil {
		util.ErrorLog("main", "服务器启动失败，原因：" + err.Error(), "main->NewServer")
	}

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), server).Error())
}