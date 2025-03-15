package handlers

import (
	"eh_go/controller/wechat/menu/path/node"
	"errors"
)

// BaseCheck 最基础的Handler
func BaseCheck(context node.Context) (string, bool) {
	command := context.GetRawCommand()
	switch command {
	case "返回":
		return context.Back(), true
	case "主页":
		return context.GoToMenu("main"), true
	case "帮助":
		return "帮助信息:\n" + context.GetNode().FindByID("help").Content, true
	default:
		return "", false
	}
}

func BaseHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	return "", errors.New("BaseHandler 无效参数处理")
}
