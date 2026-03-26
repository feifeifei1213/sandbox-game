package payload

// StateSnapshot 用于记录提交前后最小状态快照。
type StateSnapshot struct {
	YearStatus       string `json:"yearStatus"`
	StageStatus      string `json:"stageStatus"`
	ReportStatus     string `json:"reportStatus"`
	BusinessStatus   string `json:"businessStatus"`
	SummaryEffective bool   `json:"summaryEffective"`
}
