package routers

import (
	"eh_go/controller/wechat"
	"eh_go/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 健康接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "1.0.0"})
	})
	// 用户相关路由
	//userGroup := router.Group("/api/v1/users")
	//{
	//	userGroup.POST("/register", users.UsersController{}.Register)
	//	userGroup.POST("/login", users.UsersController{}.Login)
	//}
	// 初始化微信的路由
	wechat.SetupWechatRoutes(router)
	// 在路由中启用
	router.Use(middleware.Recovery())
	return router
}
