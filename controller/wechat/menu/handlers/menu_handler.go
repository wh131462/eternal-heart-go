package handlers

import (
	"eh_go/controller/wechat/menu/path/node"
	"errors"
	"strconv"
)

// MenuHandler 日期的处理
func MenuHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	target := context.GetNode()
	session := context.GetSession()
	// 可能是子路由的数字
	num, _ := strconv.ParseInt(context.GetRawCommand(), 10, 0)
	child := target.GetChild(int(num - 1))
	if child == nil {
		return "", errors.New("child not found")
	}
	session.UpdateId(child.ID)
	return context.GoToMenu(child.ID), nil
}
