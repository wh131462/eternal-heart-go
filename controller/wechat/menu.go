package wechat

import (
	"fmt"
	"log"
	"strings"
)

// 定义菜单项结构
type MenuItem struct {
	Key      string
	Name     string
	SubMenus map[string]*MenuItem
}

// 定义主菜单
var mainMenu = &MenuItem{
	Key:  "main",
	Name: "主菜单",
	SubMenus: map[string]*MenuItem{
		"1": {
			Key:  "profile",
			Name: "个人信息",
			SubMenus: map[string]*MenuItem{
				"1": {Key: "view_profile", Name: "查看信息"},
				"2": {Key: "edit_profile", Name: "修改信息"},
				"0": {Key: "back", Name: "返回上级"},
			},
		},
		"2": {
			Key:  "settings",
			Name: "系统设置",
			SubMenus: map[string]*MenuItem{
				"1": {Key: "notifications", Name: "通知设置"},
				"2": {Key: "privacy", Name: "隐私设置"},
				"0": {Key: "back", Name: "返回上级"},
			},
		},
		"0": {Key: "help", Name: "帮助信息"},
	},
}

// 用户会话状态
type UserSession struct {
	CurrentMenu *MenuItem
	Path        []string
}

// 用户会话管理
var userSessions = make(map[string]*UserSession)

// 获取用户会话
func getUserSession(userID string) *UserSession {
	session, exists := userSessions[userID]
	if !exists {
		session = &UserSession{
			CurrentMenu: mainMenu,
			Path:        []string{"main"},
		}
		userSessions[userID] = session
	}
	return session
}

// 生成菜单显示文本
func generateMenuText(menu *MenuItem) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s ===\n", menu.Name))

	for key, item := range menu.SubMenus {
		sb.WriteString(fmt.Sprintf("%s. %s\n", key, item.Name))
	}

	sb.WriteString("\n请输入对应的数字选择功能：")
	return sb.String()
}

// 处理菜单选择
func handleMenuSelection(userID string, selection string) string {
	session := getUserSession(userID)
	log.Println("收到消息:", userID, selection)

	// 处理返回上级菜单
	if selection == "0" && len(session.Path) > 1 {
		session.Path = session.Path[:len(session.Path)-1]
		session.CurrentMenu = getMenuByPath(session.Path)
		return generateMenuText(session.CurrentMenu)
	}

	// 查找选择的菜单项
	if item, exists := session.CurrentMenu.SubMenus[selection]; exists {
		// 如果是叶子节点（没有子菜单的选项）
		if len(item.SubMenus) == 0 {
			return handleMenuAction(item.Key)
		}

		// 更新当前菜单状态
		session.CurrentMenu = item
		session.Path = append(session.Path, item.Key)
		return generateMenuText(item)
	}

	return fmt.Sprintf("欢迎使用筑梦恒心公众号！\n\n您可以输入以下指令：\n1. 输入\"菜单\"显示功能列表\n2. 输入数字选择对应功能\n3. 输入\"0\"返回主菜单\n\n更多精彩功能等待您的发现！")
}

// 根据路径获取菜单
func getMenuByPath(path []string) *MenuItem {
	current := mainMenu
	for i := 1; i < len(path); i++ {
		for _, item := range current.SubMenus {
			if item.Key == path[i] {
				current = item
				break
			}
		}
	}
	return current
}

// 处理菜单动作
func handleMenuAction(key string) string {
	switch key {
	case "view_profile":
		return "您的个人信息：\n昵称：测试用户\n注册时间：2024-01-01"
	case "edit_profile":
		return "暂不支持修改个人信息"
	case "notifications":
		return "通知设置：已开启系统通知"
	case "privacy":
		return "隐私设置：当前账号可见"
	case "help":
		return "帮助信息：\n1. 输入数字选择对应功能\n2. 输入0返回上级菜单\n3. 如需帮助请输入'help'"
	default:
		return "功能开发中，敬请期待"
	}
}
