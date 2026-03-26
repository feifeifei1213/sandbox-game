package entity

type Group struct {
	ID             int64   `gorm:"column:id;primaryKey"`
	GroupNo        int     `gorm:"column:group_no"`
	GroupCode      string  `gorm:"column:group_code"`
	GroupName      string  `gorm:"column:group_name"`
	BusinessStatus string  `gorm:"column:business_status"`
	BankruptYearNo *int    `gorm:"column:bankrupt_year_no"`
	BankruptReason *string `gorm:"column:bankrupt_reason"`
	BaseEntity
}

func (Group) TableName() string {
	return "sg_group"
}
