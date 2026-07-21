package entity

import "time"

// GroupAdjustmentRevision 维护单个小组单个年份的奖惩增量同步版本。
type GroupAdjustmentRevision struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	GroupID   int64     `gorm:"column:group_id;uniqueIndex:uk_adjustment_revision_group_year,priority:1"`
	YearNo    int       `gorm:"column:year_no;uniqueIndex:uk_adjustment_revision_group_year,priority:2"`
	Revision  int64     `gorm:"column:revision"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (GroupAdjustmentRevision) TableName() string {
	return "sg_group_adjustment_revision"
}
