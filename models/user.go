package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model // 内嵌包含 ID、CreatedAt、UpdatedAt、DeletedAt 字段

	// 基础身份信息
	Name     string `gorm:"size:100;not null;comment:用户昵称" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null;comment:邮箱(唯一登录标识)" json:"email"`
	Password string `gorm:"size:255;not null;comment:加密后的密码" json:"-"` // 敏感字段禁止序列化

	// 多因素认证支持
	Phone         string `gorm:"size:20;uniqueIndex;comment:手机号" json:"phone"`
	PhoneVerified bool   `gorm:"default:false;comment:手机号是否已验证" json:"-"`

	// 状态控制
	IsActive  bool      `gorm:"default:true;comment:账户是否激活" json:"isActive"`
	LastLogin time.Time `gorm:"comment:最后登录时间" json:"lastLogin"`

	// 社交扩展
	AvatarURL string `gorm:"size:512;comment:头像地址" json:"avatarUrl"`
	OAuthBind string `gorm:"size:32;comment:第三方登录绑定" json:"-"` // 例如 github,wechat

	// 权限控制
	Role string `gorm:"size:50;default:user;comment:用户角色" json:"role"` // 例如 admin,user
}
