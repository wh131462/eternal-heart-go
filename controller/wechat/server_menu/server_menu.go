package server_menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateMenu 创建微信公众号自定义菜单
func CreateMenu(menu Menu) error {
	// 获取访问令牌
	token, err := GetAccessToken()
	if err != nil {
		return fmt.Errorf("获取访问令牌失败: %v", err)
	}

	// 准备请求数据
	menuData, err := json.Marshal(menu)
	if err != nil {
		return fmt.Errorf("菜单数据序列化失败: %v", err)
	}

	// 发送创建菜单请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", token)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(menuData))
	if err != nil {
		return fmt.Errorf("创建菜单请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var result CreateMenuResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("创建菜单失败: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// GetMenu 查询当前菜单配置
func GetMenu() (*Menu, error) {
	// 获取访问令牌
	token, err := GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	// 发送查询菜单请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("查询菜单请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result struct {
		Menu    *Menu  `json:"menu"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 && result.ErrCode != 46003 { // 46003是菜单不存在的错误码
		return nil, fmt.Errorf("查询菜单失败: %d %s", result.ErrCode, result.ErrMsg)
	}

	return result.Menu, nil
}

// DeleteMenu 删除当前菜单
func DeleteMenu() error {
	// 获取访问令牌
	token, err := GetAccessToken()
	if err != nil {
		return fmt.Errorf("获取访问令牌失败: %v", err)
	}

	// 发送删除菜单请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("删除菜单请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("删除菜单失败: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}
