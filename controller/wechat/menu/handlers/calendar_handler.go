package handlers

import (
	"eh_go/controller/wechat/menu/path/node"
	"eh_go/utils"
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"strings"
	"time"
)

// CalendarHandler 日期的处理
func CalendarHandler(context node.Context) (string, error) {
	if base, has := BaseCheck(context); has {
		return base, nil
	}
	// 在当前节点 需要判断 是1 还是指定的日期
	dateStr := context.GetRawCommand()

	if utils.NumberEqual(dateStr, 1) {
		today := time.Now()
		dateStr = today.Format("2006-01-02")
	}
	return getCalendarInfo(dateStr), nil
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

	return fmt.Sprintf("%s黄历：\n公历：%s\n农历：%s%s\n宜：%s\n忌：%s",
		datePrefix, date.Format("2006-01-02"), lunarDate, jieQiInfo, suitable, avoid)
}
