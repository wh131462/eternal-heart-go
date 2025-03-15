package sessions

import (
	"log"
	"time"
)

// UserSession 用户会话状态
type UserSession struct {
	ID        string
	Timestamp int64
}

// 用户会话管理
var userSessions = make(map[string]*UserSession)

// GetUserSession 获取用户会话
func GetUserSession(userID string) *UserSession {
	session, exists := userSessions[userID]
	if !exists {
		session = &UserSession{
			ID:        "",
			Timestamp: time.Now().Unix(),
		}
		userSessions[userID] = session
	}
	return session
}

// UpdateId 更新id
func (session *UserSession) UpdateId(id string) {
	session.ID = id
	log.Println("更新session", id)
}

// DeleteUserSession 删除用户session
func DeleteUserSession(userID string) {
	delete(userSessions, userID)
}
