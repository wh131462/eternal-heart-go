package dispatch

import (
	"eh_go/controller/wechat/menu/path/builder"
	"eh_go/controller/wechat/menu/path/context"
	"eh_go/controller/wechat/menu/sessions"
	"fmt"
	"log"
)

// Dispatch 处理菜单选择 从此识别指令去分发事件 支持 直接响应的指令 不记录 和 菜单路径 完全保持路径
func Dispatch(userID string, rawCommand string) string {
	session := sessions.GetUserSession(userID)
	menuText := ""
	log.Println("收到消息:", userID, rawCommand)
	ctx := &context.PathContext{
		UserID:     userID,
		RawCommand: rawCommand,
		Session:    session,
	}
	if session.ID == "" {
		ctx.Node = builder.Menu
		menuText = ctx.GoToHome()
	} else {
		target := builder.Menu.FindByID(session.ID)
		ctx.Node = target
		if target.HandlerFunc != nil {
			returnVal, err := target.HandlerFunc(ctx)
			if err != nil {
				log.Println(err)
				menuText = "无效的输入,请按要求输入符合要求的指令,如需帮助,请输入<帮助>查看帮助信息~"
			} else {
				menuText = returnVal
			}
		}
	}

	return fmt.Sprintf(menuText)
}
