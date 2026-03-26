package entity

import "time"

// BaseEntity 承载首版各业务主表共享的审计字段。
type BaseEntity struct {
	Creator    string    `gorm:"column:creator"`
	CreateTime time.Time `gorm:"column:create_time"`
	Updater    string    `gorm:"column:updater"`
	UpdateTime time.Time `gorm:"column:update_time"`
}
