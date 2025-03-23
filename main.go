package main

import (
	"eh_go/controller/wechat/server_menu"
	"eh_go/routers"
	"log"
)

func main() {
	// 初始化路由
	router := routers.SetupRouter()
	// 初始化数据库
	//config.InitDB()
	// 初始化微信菜单
	if err := server_menu.InitMenu(); err != nil {
		log.Printf("初始化微信菜单失败: %v", err)
	}
	// 启动服务器
	log.Println("Starting server on :9999")
	if err := router.Run(":9999"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
