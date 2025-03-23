package server_menu

// 微信菜单按钮类型
const (
	ButtonTypeClick       = "click"       // 点击推事件
	ButtonTypeView        = "view"        // 跳转URL
	ButtonTypeMiniProgram = "miniprogram" // 小程序
	ButtonTypeSubMenu     = "sub_button"  // 二级菜单
)

// Button 微信菜单按钮
type Button struct {
	Type       string   `json:"type,omitempty"`
	Name       string   `json:"name"`
	Key        string   `json:"key,omitempty"`
	URL        string   `json:"url,omitempty"`
	SubButtons []Button `json:"sub_button,omitempty"`
	AppID      string   `json:"appid,omitempty"`    // 小程序appid
	PagePath   string   `json:"pagepath,omitempty"` // 小程序页面路径
}

// Menu 微信自定义菜单
type Menu struct {
	Buttons []Button `json:"button"`
}

// CreateMenuResponse 创建菜单响应
type CreateMenuResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
