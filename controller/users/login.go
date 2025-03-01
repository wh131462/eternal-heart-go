package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (uc *UsersController) Login(c *gin.Context) {
	// 登录逻辑实现
	type LoginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 认证逻辑
	c.JSON(http.StatusOK, gin.H{"token": "generated_jwt_token"})
}
