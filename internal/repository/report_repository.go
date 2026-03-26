package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) FindByGroupIDAndYear(ctx context.Context, groupID int64, yearNo int) (*entity.GroupReport, error) {
	var item entity.GroupReport
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

type UpsertReportDraftCommand struct {
	GroupID             int64
	YearNo              int
	ReportManualPayload payload.ReportManualPayload
	LastAutoSavedAt     time.Time
	OperatorName        string
}

// UpsertDraft 保存组 + 年维度的财报草稿，采用最后一次写入覆盖当前草稿。
func (r *ReportRepository) UpsertDraft(ctx context.Context, cmd UpsertReportDraftCommand) error {
	rawPayload, err := json.Marshal(cmd.ReportManualPayload)
	if err != nil {
		return fmt.Errorf("marshal report manual payload: %w", err)
	}

	item := entity.GroupReport{
		GroupID:             cmd.GroupID,
		YearNo:              cmd.YearNo,
		ReportManualPayload: rawPayload,
		// 首次保存财报草稿时，自动计算结果尚未正式落地，这里用空 JSON 对象占位。
		ReportComputedPayload: []byte("{}"),
		LastAutoSavedAt:       &cmd.LastAutoSavedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.LastAutoSavedAt,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.LastAutoSavedAt,
		},
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "group_id"},
				{Name: "year_no"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"report_manual_payload_json": item.ReportManualPayload,
				"last_auto_saved_at":         item.LastAutoSavedAt,
				"updater":                    item.Updater,
				"update_time":                item.UpdateTime,
			}),
		}).
		Create(&item).Error
}

type UpsertCurrentReportCommand struct {
	GroupID               int64
	YearNo                int
	ReportManualPayload   payload.ReportManualPayload
	ReportComputedPayload payload.ReportComputedPayload
	BalanceCheckPassed    bool
	LastAutoSavedAt       time.Time
	SubmittedAt           *time.Time
	OperatorName          string
}

// UpsertCurrentReport 保存当前组 + 年维度的财报最新内容。
func (r *ReportRepository) UpsertCurrentReport(ctx context.Context, cmd UpsertCurrentReportCommand) error {
	rawManualPayload, err := json.Marshal(cmd.ReportManualPayload)
	if err != nil {
		return fmt.Errorf("marshal report manual payload: %w", err)
	}
	rawComputedPayload, err := json.Marshal(cmd.ReportComputedPayload)
	if err != nil {
		return fmt.Errorf("marshal report computed payload: %w", err)
	}

	item := entity.GroupReport{
		GroupID:               cmd.GroupID,
		YearNo:                cmd.YearNo,
		ReportManualPayload:   rawManualPayload,
		ReportComputedPayload: rawComputedPayload,
		BalanceCheckPassed:    cmd.BalanceCheckPassed,
		LastAutoSavedAt:       &cmd.LastAutoSavedAt,
		SubmittedAt:           cmd.SubmittedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.LastAutoSavedAt,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.LastAutoSavedAt,
		},
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "group_id"},
				{Name: "year_no"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"report_manual_payload_json":   item.ReportManualPayload,
				"report_computed_payload_json": item.ReportComputedPayload,
				"balance_check_passed":         item.BalanceCheckPassed,
				"last_auto_saved_at":           item.LastAutoSavedAt,
				"submitted_at":                 item.SubmittedAt,
				"updater":                      item.Updater,
				"update_time":                  item.UpdateTime,
			}),
		}).
		Create(&item).Error
}

type CreateReportSubmissionCommand struct {
	GroupID                int64
	YearNo                 int
	SubmitVersion          int
	ReportManualSnapshot   payload.ReportManualPayload
	ReportComputedSnapshot payload.ReportComputedPayload
	BalanceCheckPassed     bool
	StateBeforeJSON        []byte
	StateAfterJSON         []byte
	SubmitterID            int64
	SubmitTime             time.Time
}

// CreateSubmission 写入财报正式提交流水。
func (r *ReportRepository) CreateSubmission(ctx context.Context, cmd CreateReportSubmissionCommand) error {
	rawManualPayload, err := json.Marshal(cmd.ReportManualSnapshot)
	if err != nil {
		return fmt.Errorf("marshal report manual snapshot: %w", err)
	}
	rawComputedPayload, err := json.Marshal(cmd.ReportComputedSnapshot)
	if err != nil {
		return fmt.Errorf("marshal report computed snapshot: %w", err)
	}

	item := entity.GroupReportSubmission{
		GroupID:                cmd.GroupID,
		YearNo:                 cmd.YearNo,
		SubmitVersion:          cmd.SubmitVersion,
		ReportManualSnapshot:   rawManualPayload,
		ReportComputedSnapshot: rawComputedPayload,
		BalanceCheckPassed:     cmd.BalanceCheckPassed,
		StateBefore:            cmd.StateBeforeJSON,
		StateAfter:             cmd.StateAfterJSON,
		SubmitterID:            cmd.SubmitterID,
		SubmitTime:             cmd.SubmitTime,
	}

	return r.db.WithContext(ctx).Create(&item).Error
}
