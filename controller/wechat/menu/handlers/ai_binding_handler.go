package handlers

import (
	"eh_go/controller/wechat/menu/path/node"
)

// AIBindingHandler AI绑定API功能
func AIBindingHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	// todo 处理AI的绑定
	//command := context.GetRawCommand()
	return "暂未开放,敬请期待~", nil
}
