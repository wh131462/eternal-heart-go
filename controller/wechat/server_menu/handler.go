package server_menu

import (
	"log"
)

// HandleServerMenuClick 处理服务菜单的点击
func HandleServerMenuClick(userID, eventKey string) string {
	log.Printf("用户 %s 点击菜单: %s", userID, eventKey)
	switch eventKey {
	case "ANY":
		return "Anything you want will be true."
	default:
		return "未知的菜单命令，请重试"
	}
}
