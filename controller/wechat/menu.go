package wechat

import (
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/6tail/lunar-go/calendar"
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
			Key:  "calendar",
			Name: "📅 万年历中心",
			SubMenus: map[string]*MenuItem{
				"1": {Key: "today_calendar", Name: "今日黄历"},
				"2": {Key: "calendar_wiki", Name: "黄历百科"},
				"3": {Key: "year_calendar", Name: "年度日历"},
				"0": {Key: "back", Name: "返回上级"},
			},
		},
		"2": {
			Key:  "birthday",
			Name: "🎂 生日管家",
			SubMenus: map[string]*MenuItem{
				"1": {Key: "add_birthday", Name: "添加生日"},
				"2": {Key: "birthday_list", Name: "我的生日簿"},
				"3": {Key: "reminder_settings", Name: "提醒设置"},
				"0": {Key: "back", Name: "返回上级"},
			},
		},
		"3": {
			Key:  "service",
			Name: "⚙️ 我的服务",
			SubMenus: map[string]*MenuItem{
				"1": {Key: "user_guide", Name: "使用指南"},
				"2": {Key: "data_manage", Name: "数据管理"},
				"3": {Key: "feedback", Name: "意见反馈"},
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
	CurrentDir  string // 缓存当前所在目录
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
			CurrentDir:  "main", // 初始化当前目录
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

	// 处理特殊指令
	if selection == "菜单" {
		session.CurrentMenu = mainMenu
		session.Path = []string{"main"}
		session.CurrentDir = "main"
		return generateMenuText(mainMenu)
	}

	// 处理返回上级菜单
	if selection == "0" {
		// 只有在当前不是主菜单时才允许返回
		if session.CurrentDir != "main" {
			// 移除当前路径的最后一个元素
			if len(session.Path) > 1 {
				session.Path = session.Path[:len(session.Path)-1]
				// 更新当前菜单到父级菜单
				session.CurrentMenu = getMenuByPath(session.Path)
				// 更新当前目录为新的路径最后一个元素
				session.CurrentDir = session.Path[len(session.Path)-1]
				return generateMenuText(session.CurrentMenu)
			}
		}
		// 如果已经在主菜单，则显示主菜单
		session.CurrentMenu = mainMenu
		session.Path = []string{"main"}
		session.CurrentDir = "main"
		return generateMenuText(mainMenu)
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
		session.CurrentDir = item.Key
		return generateMenuText(item)
	}

	// 处理日期查询
	if isDateFormat(selection) {
		return getCalendarInfo(selection)
	}

	// 处理生日添加指令
	if strings.HasPrefix(selection, "生日+") {
		parts := strings.Split(selection, "+")
		if len(parts) == 3 {
			return fmt.Sprintf("收到添加生日请求：\n姓名：%s\n日期：%s\n\n功能开发中，敬请期待！", parts[1], parts[2])
		}
	}

	// 如果当前在子菜单中，显示当前菜单
	if session.CurrentDir != "main" {
		return generateMenuText(session.CurrentMenu)
	}

	// 其他情况显示欢迎信息
	return fmt.Sprintf("欢迎使用筑梦恒心公众号！\n\n您可以：\n1. 输入\"菜单\"显示功能列表\n2. 输入数字选择对应功能\n3. 直接输入日期(如：2024-01-19)查询黄历\n4. 输入\"生日+姓名+日期\"添加生日提醒\n\n更多精彩功能等待您的发现！")
}

// 日期格式验证和转换
func isDateFormat(text string) bool {
	// 支持 2024-01-19 或 2024年1月19日 格式
	return strings.Contains(text, "-") || (strings.Contains(text, "年") && strings.Contains(text, "月") && strings.Contains(text, "日"))
}

// 获取日历信息
func getCalendarInfo(dateStr string) string {
	var date time.Time
	var err error

	// 处理不同格式的日期
	if strings.Contains(dateStr, "-") {
		date, err = time.Parse("2006-01-02", dateStr)
	} else {
		// 处理中文日期格式（简化处理）
		dateStr = strings.ReplaceAll(dateStr, "年", "-")
		dateStr = strings.ReplaceAll(dateStr, "月", "-")
		dateStr = strings.ReplaceAll(dateStr, "日", "")
		date, err = time.Parse("2006-1-2", dateStr)
	}

	if err != nil {
		return "日期格式不正确，请使用正确的格式（如：2024-01-19 或 2024年1月19日）"
	}

	// 判断是否为今日
	today := time.Now()
	isToday := date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()

	// 使用农历库获取详细信息
	lunar := calendar.NewLunarFromDate(date)
	
	// 获取农历日期
	lunarDate := fmt.Sprintf("%s年%s月%s", lunar.GetYearInChinese(), lunar.GetMonthInChinese(), lunar.GetDayInChinese())
	
	// 获取节气信息
	jieQi := lunar.GetJieQi()
	jieQiInfo := ""
	if jieQi != "" {
		jieQiInfo = fmt.Sprintf("\n节气：%s", jieQi)
	}

	// 获取宜忌
	var suitableItems []string
	yiList := lunar.GetDayYi()
	for e := yiList.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			suitableItems = append(suitableItems, str)
		}
	}

	var avoidItems []string
	jiList := lunar.GetDayJi()
	for e := jiList.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			avoidItems = append(avoidItems, str)
		}
	}

	suitable := strings.Join(suitableItems, "、")
	avoid := strings.Join(avoidItems, "、")

	datePrefix := "指定日期"
	if isToday {
		datePrefix = "今日"
	}

	return fmt.Sprintf("%s黄历：\n公历：%s\n农历：%s%s\n宜：%s\n忌：%s\n\n回复【日期】如2024-08-15可查看指定日期黄历",
		datePrefix, date.Format("2006-01-02"), lunarDate, jieQiInfo, suitable, avoid)
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
	case "today_calendar":
		return getCalendarInfo(time.Now().Format("2006-01-02"))
	case "year_calendar":
		return "年度日历功能开发中，敬请期待！"
	case "calendar_wiki":
		return "黄历百科：\n1. 节气：二十四节气说明\n2. 吉凶宜忌详解\n3. 传统节日习俗\n4. 择日指南"
	case "add_birthday":
		return "请按以下格式添加生日：\n生日+姓名+日期\n例如：生日+张三+2000-01-01\n或：生日+李四+农历2000正月初一"
	case "birthday_list":
		return "暂无保存的生日信息，请先添加生日。"
	case "reminder_settings":
		return "提醒设置：\n1. 提醒时间：提前1天\n2. 提醒方式：文字提醒\n3. 农历转换：自动转换"
	case "user_guide":
		return "使用指南：\n1. 回复数字选择功能\n2. 输入日期查询黄历\n3. 生日+姓名+日期添加提醒\n4. 输入0返回上级菜单"
	case "data_manage":
		return "数据管理功能开发中，敬请期待！"
	case "feedback":
		return "请直接输入您的意见或建议，我们会认真听取并改进！"
	case "help":
		return "帮助信息：\n1. 输入数字选择对应功能\n2. 输入0返回上级菜单\n3. 直接输入日期查询黄历\n4. 如需帮助请输入'help'"
	default:
		return "功能开发中，敬请期待"
	}
}
