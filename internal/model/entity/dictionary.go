package entity

import "time"

type DictionaryScheme struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	EditionCode string `gorm:"column:edition_code"`
	SchemeName  string `gorm:"column:scheme_name"`
	Description string `gorm:"column:description"`
	BuiltIn     bool   `gorm:"column:built_in"`
	BaseEntity
}

func (DictionaryScheme) TableName() string {
	return "sg_dictionary_scheme"
}

type DictionarySchemeItem struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	SchemeID       int64  `gorm:"column:scheme_id"`
	EditionCode    string `gorm:"column:edition_code"`
	ItemCode       string `gorm:"column:item_code"`
	ItemCategory   string `gorm:"column:item_category"`
	DefaultName    string `gorm:"column:default_name"`
	DisplayName    string `gorm:"column:display_name"`
	DisplayOrder   int    `gorm:"column:display_order"`
	Editable       bool   `gorm:"column:editable"`
	RelatedPayload []byte `gorm:"column:related_payload"`
	BaseEntity
}

func (DictionarySchemeItem) TableName() string {
	return "sg_dictionary_scheme_item"
}

type CurrentDictionaryItem struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	EditionCode    string `gorm:"column:edition_code"`
	ItemCode       string `gorm:"column:item_code"`
	ItemCategory   string `gorm:"column:item_category"`
	DefaultName    string `gorm:"column:default_name"`
	DisplayName    string `gorm:"column:display_name"`
	DisplayOrder   int    `gorm:"column:display_order"`
	Editable       bool   `gorm:"column:editable"`
	RelatedPayload []byte `gorm:"column:related_payload"`
	BaseEntity
}

func (CurrentDictionaryItem) TableName() string {
	return "sg_current_dictionary_item"
}

type DictionaryChangeLog struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	EditionCode  string    `gorm:"column:edition_code"`
	ChangeType   string    `gorm:"column:change_type"`
	SchemeID     *int64    `gorm:"column:scheme_id"`
	Reason       string    `gorm:"column:reason"`
	BeforeJSON   []byte    `gorm:"column:before_json"`
	AfterJSON    []byte    `gorm:"column:after_json"`
	Revision     int       `gorm:"column:revision"`
	OperatorID   int64     `gorm:"column:operator_id"`
	OperatorName string    `gorm:"column:operator_name"`
	OperateTime  time.Time `gorm:"column:operate_time"`
}

func (DictionaryChangeLog) TableName() string {
	return "sg_dictionary_change_log"
}
