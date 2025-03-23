package server_menu

import (
	"eh_go/config"
	"log"
)

// InitMenu 初始化微信公众号菜单
func InitMenu() error {
	// 定义菜单结构
	menu := Menu{
		Buttons: []Button{
			{
				Name:     "EH-TOOLS",
				Type:     ButtonTypeMiniProgram,
				AppID:    config.WechatConfig.MiniAppID,
				PagePath: config.WechatConfig.MiniAppPath,
			},
			{
				Name: "关于我",
				SubButtons: []Button{
					{
						Type: ButtonTypeView,
						Name: "GITHUB",
						URL:  config.WechatConfig.GithubUrl,
					},
					{
						Type: ButtonTypeView,
						Name: "个人网站",
						URL:  config.WechatConfig.WebsiteUrl,
					},
				},
			},
		},
	}

	// 创建菜单
	if err := CreateMenu(menu); err != nil {
		log.Printf("创建菜单失败: %v", err)
		return err
	}

	log.Println("微信公众号菜单创建成功")
	return nil
}
