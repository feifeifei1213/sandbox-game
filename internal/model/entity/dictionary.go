package entity

import "time"

type DictionaryScheme struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	EditionCode string `gorm:"column:edition_code;type:varchar(64);not null"`
	SchemeName  string `gorm:"column:scheme_name;type:varchar(64);not null"`
	Description string `gorm:"column:description;type:varchar(255);not null;default:''"`
	BuiltIn     bool   `gorm:"column:built_in"`
	BaseEntity
}

func (DictionaryScheme) TableName() string {
	return "sg_dictionary_scheme"
}

type DictionarySchemeItem struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	SchemeID       int64  `gorm:"column:scheme_id"`
	EditionCode    string `gorm:"column:edition_code;type:varchar(64);not null"`
	ItemCode       string `gorm:"column:item_code;type:varchar(128);not null"`
	ItemCategory   string `gorm:"column:item_category;type:varchar(32);not null"`
	DefaultName    string `gorm:"column:default_name;type:varchar(128);not null"`
	DisplayName    string `gorm:"column:display_name;type:varchar(128);not null"`
	DisplayOrder   int    `gorm:"column:display_order"`
	Editable       bool   `gorm:"column:editable"`
	RelatedPayload []byte `gorm:"column:related_payload;type:json"`
	BaseEntity
}

func (DictionarySchemeItem) TableName() string {
	return "sg_dictionary_scheme_item"
}

type CurrentDictionaryItem struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	EditionCode    string `gorm:"column:edition_code;type:varchar(64);not null"`
	ItemCode       string `gorm:"column:item_code;type:varchar(128);not null"`
	ItemCategory   string `gorm:"column:item_category;type:varchar(32);not null"`
	DefaultName    string `gorm:"column:default_name;type:varchar(128);not null"`
	DisplayName    string `gorm:"column:display_name;type:varchar(128);not null"`
	DisplayOrder   int    `gorm:"column:display_order"`
	Editable       bool   `gorm:"column:editable"`
	RelatedPayload []byte `gorm:"column:related_payload;type:json"`
	BaseEntity
}

func (CurrentDictionaryItem) TableName() string {
	return "sg_current_dictionary_item"
}

type DictionaryChangeLog struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	EditionCode  string    `gorm:"column:edition_code;type:varchar(64);not null"`
	ChangeType   string    `gorm:"column:change_type;type:varchar(32);not null"`
	SchemeID     *int64    `gorm:"column:scheme_id"`
	Reason       string    `gorm:"column:reason;type:varchar(255);not null;default:''"`
	BeforeJSON   []byte    `gorm:"column:before_json;type:json;not null"`
	AfterJSON    []byte    `gorm:"column:after_json;type:json;not null"`
	Revision     int       `gorm:"column:revision"`
	OperatorID   int64     `gorm:"column:operator_id"`
	OperatorName string    `gorm:"column:operator_name;type:varchar(64);not null"`
	OperateTime  time.Time `gorm:"column:operate_time;type:datetime;not null"`
}

func (DictionaryChangeLog) TableName() string {
	return "sg_dictionary_change_log"
}
