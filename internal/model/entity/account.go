package entity

import "time"

type Account struct {
	ID            int64      `gorm:"column:id;primaryKey"`
	Username      string     `gorm:"column:username"`
	PasswordHash  string     `gorm:"column:password_hash"`
	RoleType      string     `gorm:"column:role_type"`
	GroupID       *int64     `gorm:"column:group_id"`
	Status        string     `gorm:"column:status"`
	LastLoginTime *time.Time `gorm:"column:last_login_time"`
	BaseEntity
}

func (Account) TableName() string {
	return "sg_account"
}
