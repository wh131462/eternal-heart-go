package main

import (
	"eh_go/routers"
	"log"
	"net/http"
)

func main() {
	// 初始化路由
	routers.SetupRouter()
	// 初始化数据库
	//config.InitDB()
	// 启动服务器
	log.Println("Starting server on :2341")
	if err := http.ListenAndServe(":2341", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
