package manager

import (
	"eh_go/controller/wechat/menu/path/builder"
	"eh_go/controller/wechat/menu/path/node"
	"eh_go/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var Slogan = "\t𝑩𝒖𝒊𝒍𝒅 𝒚𝒐𝒖𝒓 𝒅𝒓𝒆𝒂𝒎 𝒃𝒚 𝒚𝒐𝒖𝒓 𝒉𝒆𝒂𝒓𝒕."
var Player = `𝗗𝗲𝗹𝗮𝗰𝗲𝘆 - 𝗗𝗿𝗲𝗮𝗺 𝗜𝘁 𝗣𝗼𝘀𝘀𝗶𝗯𝗹𝗲
0:59━━●━━━━━──3:20
  ⇆       ◁      ❚❚      ▷       ↻`
var DivideTag = "─━─━─━─━─━─━─" //`━━━━━━━━━━━━━━`
var Power = "\t𝑃𝑜𝑤𝑒𝑟𝑒𝑑 𝑏𝑦 𝐸𝑡𝑒𝑟𝑛𝑎𝑙𝐻𝑒𝑎𝑟𝑡."

// GetMenuText 获取菜单的text 只通过id即可 传入user创建的时间戳实现动态音乐播放器(模拟)
func GetMenuText(id string, timestamp int64) string {
	target := builder.Menu.FindByID(id)
	breadcrumbs := BuildBreadcrumbs(target)
	subMenu := BuildSubMenuText(target)
	totalTime := 200 // 秒 3分20秒
	curTime := int(time.Now().Unix()-timestamp) % totalTime
	player := GeneratePlayer(curTime, totalTime)
	return BuildMenuText(target.Content, breadcrumbs, subMenu, player)
}

// BuildBreadcrumbs 生成面包屑
func BuildBreadcrumbs(tree *node.PathNode) string {
	paths := utils.Map(tree.Path(), func(u *node.PathNode, i int) string {
		return u.Name
	})
	return "➤ 当前位置:" + strings.Join(paths, " > ")
}

// BuildSubMenuText 构建子菜单
func BuildSubMenuText(tree *node.PathNode) string {
	if len(tree.Children) < 1 {
		return "输入<返回>回到上一级菜单"
	}
	subMenus := utils.Map(tree.Children, func(u *node.PathNode, i int) string {
		return strconv.Itoa(i+1) + "." + u.Name
	})
	return "输入对应序号查看:\n" + strings.Join(subMenus, "\n")
}

// BuildMenuText 生成标准菜单项目文本
func BuildMenuText(content string, breadcrumbs string, subMenu string, player string) string {
	textList := []string{DivideTag + "\n" + breadcrumbs, content, subMenu, Slogan, player, Power}
	return strings.Join(textList, "\n"+DivideTag+"\n")
}

// GeneratePlayer 生成音乐播放器进度条
// currentTime - 当前播放时间（秒）
// totalDuration - 总时长（秒）
func GeneratePlayer(currentTime int, totalDuration int) string {
	// 时间格式转换 (秒 -> 分:秒)
	formatTime := func(seconds int) string {
		return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
	}

	// 进度条参数
	const (
		totalBlocks   = 10  // 进度条总块数
		progressChar  = "━" // 已播放部分字符
		remainingChar = "─" // 未播放部分字符
		indicatorChar = "●" // 进度指示符
	)

	// 计算进度比例 (0.0 ~ 1.0)
	progress := math.Min(1.0, math.Max(0.0, float64(currentTime)/float64(totalDuration)))

	// 生成动态进度条
	progressPos := int(math.Round(float64(totalBlocks) * progress))
	progressBar := strings.Repeat(progressChar, progressPos) +
		indicatorChar +
		strings.Repeat(remainingChar, totalBlocks-progressPos)

	// 替换模板内容
	return fmt.Sprintf(`𝗗𝗲𝗹𝗮𝗰𝗲𝘆 - 𝗗𝗿𝗲𝗮𝗺 𝗜𝘁 𝗣𝗼𝘀𝘀𝗶𝗯𝗹𝗲
%s%s%s
  ⇆       ◁      ❚❚      ▷       ↻`,
		formatTime(currentTime),   // 当前时间
		progressBar,               // 进度条
		formatTime(totalDuration)) // 总时长
}
