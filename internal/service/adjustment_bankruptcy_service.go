package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func markAdjustmentBankrupt(
	ctx context.Context,
	tx *gorm.DB,
	adjustment entity.GroupAdjustment,
	impact AdjustmentImpactResult,
	operatorID int64,
	operatorName string,
	operateTime time.Time,
	operation string,
) error {
	reason := fmt.Sprintf("奖惩%s后所得税后现金为 %.2f，小于 0", operation, impact.CashAfter)
	if err := repository.NewGroupRepository(tx).MarkBankrupt(ctx, adjustment.GroupID, adjustment.YearNo, reason, operatorName); err != nil {
		return fmt.Errorf("mark adjustment bankruptcy: %w", err)
	}
	_, err := createGroupSnapshot(ctx, tx, CreateGroupSnapshotCommand{
		GroupID: adjustment.GroupID, YearNo: adjustment.YearNo,
		StageCode:    impact.ResolvedStageCode,
		SnapshotType: enum.SnapshotTypeAuto,
		TriggerCode:  enum.SnapshotTriggerAdjustmentBankruptcy,
		Description:  "奖惩变更触发的只读破产快照",
		OperatorID:   operatorID, OperatorName: operatorName, OperateTime: operateTime,
		UseForRestore: false,
		Metadata: map[string]any{
			"adjustmentId":    adjustment.ID,
			"operation":       operation,
			"impact":          impact,
			"operatingResult": impact.OperatingAfter,
			"reportPreview":   impact.ReportAfter,
			"bankruptAt":      operateTime,
			"bankruptReason":  reason,
		},
	})
	if err != nil {
		return fmt.Errorf("create adjustment bankruptcy snapshot: %w", err)
	}
	return nil
}
