package routers

import (
	"eh_go/controller/users"
	"eh_go/controller/wechat"
	"eh_go/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 设置为调试模式
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	// 添加日志中间件
	router.Use(gin.Logger())
	// 在路由中启用
	router.Use(middleware.Recovery())
	// 添加一个根路由，用于测试
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to API"})
	})
	// 健康接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "1.0.0"})
	})
	// API 版本前缀
	v1 := router.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", (&users.UsersController{}).Register)
			userGroup.POST("/login", (&users.UsersController{}).Login)
		}
		// 初始化微信的路由
		wechat.SetupWechatRoutes(router)
	}
	return router
}
