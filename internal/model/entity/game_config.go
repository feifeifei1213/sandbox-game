package entity

type GameConfig struct {
	ID                       int64  `gorm:"column:id;primaryKey"`
	FinalYear                int    `gorm:"column:final_year"`
	CurrentOpenYear          int    `gorm:"column:current_open_year"`
	EditionCode              string `gorm:"column:edition_code"`
	EditionName              string `gorm:"column:edition_name"`
	RuleVersion              string `gorm:"column:rule_version"`
	TemplateVersion          string `gorm:"column:template_version"`
	OperatingTemplateVersion string `gorm:"column:operating_template_version"`
	ReportTemplateVersion    string `gorm:"column:report_template_version"`
	OrderTemplateVersion     string `gorm:"column:order_template_version"`
	ProcessRuleVersion       string `gorm:"column:process_rule_version"`
	DictionaryRevision       int    `gorm:"column:dictionary_revision"`
	InitialBaselineSubmitted bool   `gorm:"column:initial_baseline_submitted"`
	BaseEntity
}

func (GameConfig) TableName() string {
	return "sg_game_config"
}
