package server_menu

import (
	"eh_go/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// 全局配置
var (
	accessToken    string
	accessTokenExp int64
	tokenMutex     sync.Mutex
)

// GetAccessToken 获取微信接口访问令牌
func GetAccessToken() (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// 检查令牌是否过期
	if accessToken != "" && time.Now().Unix() < accessTokenExp-60 {
		return accessToken, nil
	}

	// 请求新的访问令牌
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		config.WechatConfig.AppID, config.WechatConfig.AppSecret)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %d %s", result.ErrCode, result.ErrMsg)
	}

	// 更新令牌和过期时间
	accessToken = result.AccessToken
	accessTokenExp = time.Now().Unix() + int64(result.ExpiresIn)
	log.Printf("获取新的访问令牌成功，有效期至: %v", time.Unix(accessTokenExp, 0))

	return accessToken, nil
}
