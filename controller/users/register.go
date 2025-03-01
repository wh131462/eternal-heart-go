package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (uc *UsersController) Register(c *gin.Context) {
	// 注册逻辑实现
	type RegisterReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用业务层处理
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
