package builder

import (
	"eh_go/controller/wechat/menu/handlers"
	"eh_go/controller/wechat/menu/path/node"
	"eh_go/utils"
)

var Menu = BuildMenuTree()

// BuildMenuTree 构建菜单树
func BuildMenuTree() *node.PathNode {
	root := node.NewRoot("main", "主页", "欢迎访问『筑梦恒心』!\n愿以恒心筑就你的梦想~\n在此我们为你提供了很多有用的功能~")
	root.AddHandler(handlers.MenuHandler)
	addCalendar(root)
	addBirth(root)
	addAIBinding(root)
	addRepeater(root)
	addHelper(root)
	return root
}

// 添加日历节点
func addCalendar(root *node.PathNode) {
	funcListText := utils.CreateListText([]string{
		"输入<1>查看今日黄历,也可以输入指定日期<yyyy-MM-dd>或者<yyyy年MM月dd日>来查看指定日期的黄历,例:2012-12-31",
	})
	calendar := &node.PathNode{
		ID:          "calendar",
		Name:        "日历服务",
		Content:     "这里有一些有用的日历服务:\n" + funcListText,
		HandlerFunc: handlers.CalendarHandler,
	}
	root.AddChild(calendar)
}

func addBirth(root *node.PathNode) {
	funcListText := utils.CreateListText([]string{
		"更新/新增:<更新:张三:1999-12-31>记录公历生日或<更新:张三:甲子年正月初三>记录农历生日,也支持指定年份<更新:张三:1999年正月初三>",
		"查看:<查看:张三>查看指定人员或<查看>可以查看全部记录",
		"删除:<删除:张三>删除记录",
	})
	birth := &node.PathNode{
		ID:          "birth",
		Name:        "生日服务",
		Content:     "管理重要的人的生日,可以设置到期提醒哦~(因为一些原因,暂时只能通过邮件)\n按格式输入指令可以进行管理哦,注意名字是唯一的,如果处理同名最好加一些提示区分哦~" + funcListText,
		HandlerFunc: handlers.BirthdayHandler,
	}
	root.AddChild(birth)
}

func addAIBinding(root *node.PathNode) {
	funcListText := utils.CreateListText([]string{
		"绑定api-key:<绑定:[name]:[key]>绑定一个api-key,并指定一个名字,可以用来管理api-key",
		"解绑api-key:<解绑:[name]>,name为指定的api名称",
		"查看绑定:<查看:[name]>,name为指定的api名称,或<查看>可以查看当前全部绑定",
		"刷新上下文:<刷新上下文>刷新当前上下文,注意此操作不可逆",
	})
	aiBinding := &node.PathNode{
		ID:          "ai-binding",
		Name:        "人工智能服务",
		Content:     "可以在此绑定自己的api-key来实现接口消费token,以免纠纷,绑定即视为同意本公众号代理.(每个人的api-key只会绑定自己的微信uuid,所以请放心不会产生误用token!)\n绑定后直接对话即可!\n" + funcListText,
		HandlerFunc: handlers.AIBindingHandler,
	}
	root.AddChild(aiBinding)
}

// 复读机
func addRepeater(root *node.PathNode) {
	birth := &node.PathNode{
		ID:          "repeater",
		Name:        "我是复读机",
		Content:     "我是复读机,现在开始你说什么话我都会复读!(停止复读请说:<我不玩了!!!>)",
		HandlerFunc: handlers.RepeaterHandler,
	}
	root.AddChild(birth)
}

// 添加帮助信息
func addHelper(root *node.PathNode) {
	// 帮助路径
	helpNode := &node.PathNode{
		ID:   "help",
		Name: "帮助",
		Content: utils.CreateListText([]string{
			"当前位置标识你所在的菜单具体位置",
			"在任意位置输入<帮助>都可以返回提示信息",
			"在任意菜单输入<返回>都可以返回上一级目录",
			"输入<主页>可以返回主菜单",
			`如果没有特别说明,在"<>"中的内容都为指令的内容`,
		}),
		HandlerFunc: handlers.BaseHandler,
	}
	root.AddChild(helpNode)
}

// utils
