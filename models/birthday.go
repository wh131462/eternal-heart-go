package models

import (
	"gorm.io/gorm"
	"time"
)

type RemindType string

const (
	RemindByDay  RemindType = "day"  // 提前 N 天提醒
	RemindByTime RemindType = "time" // 指定具体时间提醒
)

type Birthday struct {
	gorm.Model

	// 基础信息
	UserID      uint   `gorm:"not null;index;comment:关联用户ID" json:"userId"`
	ContactName string `gorm:"size:100;not null;comment:联系人姓名" json:"contactName"`

	// 日期配置
	SolarDate  *time.Time `gorm:"comment:阳历生日(二选一)" json:"solarDate"`   // 格式: 2006-01-02
	LunarMonth int        `gorm:"comment:农历月份(1-12)" json:"lunarMonth"` // 1-12
	LunarDay   int        `gorm:"comment:农历日期(1-30)" json:"lunarDay"`   // 1-30
	LunarLeap  bool       `gorm:"comment:是否为闰月" json:"lunarLeap"`       // 处理闰月场景

	// 提醒设置
	RemindType string `gorm:"size:20;default:day;comment:提醒类型(day/time)" json:"remindType"` // day:按天数提前, time:指定时间
	RemindDays int    `gorm:"default:1;comment:提前提醒天数" json:"remindDays"`                   // 仅当RemindType=day时生效
	RemindTime string `gorm:"size:5;default:'09:00';comment:提醒时间(HH:mm)" json:"remindTime"` // 格式: 15:04
	Enabled    bool   `gorm:"default:true;comment:是否启用提醒" json:"enabled"`

	// 扩展字段
	Relation string `gorm:"size:50;comment:关系描述(如父亲/同事)" json:"relation"`
	Notes    string `gorm:"type:text;comment:备注信息" json:"notes"`

	// 关联关系
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
