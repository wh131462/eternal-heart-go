package middleware

import "github.com/gin-gonic/gin"

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{"error": "服务器内部错误"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
