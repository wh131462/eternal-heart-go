package handlers

import (
	"eh_go/controller/wechat/menu/path/node"
)

// BirthdayHandler 生日的处理
func BirthdayHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	// todo 处理生日
	//command := context.GetRawCommand()
	return "暂未开放,敬请期待~", nil
}
