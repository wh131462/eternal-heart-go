package handlers

import "eh_go/controller/wechat/menu/path/node"

func RepeaterHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	mes := context.GetRawCommand()
	if mes == "我不玩了!!!" {
		return "OK,那好吧~\n" + context.GoToMenu("main"), nil
	}
	return mes, nil
}
