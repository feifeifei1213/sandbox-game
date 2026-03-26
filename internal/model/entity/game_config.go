package entity

type GameConfig struct {
	ID                       int64  `gorm:"column:id;primaryKey"`
	FinalYear                int    `gorm:"column:final_year"`
	CurrentOpenYear          int    `gorm:"column:current_open_year"`
	RuleVersion              string `gorm:"column:rule_version"`
	TemplateVersion          string `gorm:"column:template_version"`
	InitialBaselineSubmitted bool   `gorm:"column:initial_baseline_submitted"`
	BaseEntity
}

func (GameConfig) TableName() string {
	return "sg_game_config"
}
