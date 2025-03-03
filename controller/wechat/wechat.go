package wechat

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Token = "wh131462"
)

// VerifySignatureMiddleware 中间件：微信签名验证
func VerifySignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")

		// 验证签名
		if !verifySignature(signature, timestamp, nonce, Token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}

		// 处理验证请求
		if c.Request.Method == "GET" {
			c.String(http.StatusOK, echostr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetupWechatRoutes 路由处理函数
func SetupWechatRoutes(router *gin.Engine) {
	wechatGroup := router.Group("/api/v1/wx")
	{
		wechatGroup.Use(VerifySignatureMiddleware())
		wechatGroup.POST("", handleWxMessage)
		wechatGroup.GET("", func(c *gin.Context) {}) // 由中间件处理验证请求
	}
}

// 核心处理逻辑
func handleWxMessage(c *gin.Context) {
	var msg WxMessage
	if err := c.ShouldBindXML(&msg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid xml format"})
		return
	}
	response := dispatchMessage(msg)
	c.XML(http.StatusOK, response)
}

// 签名验证实现
func verifySignature(signature, timestamp, nonce, token string) bool {
	params := []string{token, timestamp, nonce}
	sort.Strings(params)

	hash := sha1.New()
	io.WriteString(hash, strings.Join(params, ""))
	return fmt.Sprintf("%x", hash.Sum(nil)) == signature
}

// 消息分发逻辑（保持不变）
func dispatchMessage(msg WxMessage) WxResponse {
	switch msg.MsgType {
	case "text":
		return handleTextMessage(msg)
	case "event":
		return handleEventMessage(msg)
	default:
		return WxResponse{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "不支持此类消息,请检查消息!",
		}
	}
}

// 文本消息处理
func handleTextMessage(msg WxMessage) WxResponse {
	// 处理菜单选择
	responseContent := handleMenuSelection(msg.FromUserName, msg.Content)
	log.Printf("用户 %s 发送消息: %s, 响应: %s", msg.FromUserName, msg.Content, responseContent)
	
	return WxResponse{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      responseContent,
	}
}

// 事件消息处理
func handleEventMessage(msg WxMessage) WxResponse {
	var content string
	switch msg.Event {
	case "subscribe":
		// 发送欢迎消息和主菜单
		session := getUserSession(msg.FromUserName)
		content = fmt.Sprintf("【欢迎关注筑梦恒心】\n\n回复数字获取服务：\n1️⃣ 查今日黄历\n2️⃣ 添加生日提醒\n3️⃣ 查看帮助手册\n\n🌟 输入任意日期（如\"2024-10-1\"）立即查询\n🎂 输入\"生日+名字+日期\"快速添加提醒\n\n%s", generateMenuText(session.CurrentMenu))
	case "unsubscribe":
		// 清理用户会话
		delete(userSessions, msg.FromUserName)
		content = "感谢使用，期待再次相见！"
	default:
		content = "暂不支持的事件类型"
	}
	return WxResponse{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}
}
