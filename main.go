package main

import (
	"k8s-platform/config"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/router"
	"k8s-platform/service"
	"k8s-platform/service/cicd"
	"net/http"

	_ "k8s-platform/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化gin对象
	r := gin.Default()
	//初始化k8s client
	// service.K8s.Init()
	//初始化Jenkins
	cicd.Jenkins.Init()
	//数据库初始化
	db.Init()
	//注册中间件
	r.Use(middle.Cors())
	// r.Use(middle.JWTAuth())
	//初始化路由规则
	router.Router.InitApiRouter(r)
	//终端websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
	}()
	//gin程序启动
	r.Run(config.ListenAddr)
}
