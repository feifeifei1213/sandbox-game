package assembler

import (
	"fmt"

	"sandbox-game/internal/model/entity"
)

func buildRollbackNotice(yearState entity.GroupYearState) string {
	if !yearState.RollbackPending || yearState.RollbackTargetYearNo == nil {
		return ""
	}
	targetStage := "财报"
	if yearState.RollbackTargetStageCode != nil && *yearState.RollbackTargetStageCode != "" {
		targetStage = *yearState.RollbackTargetStageCode
	}
	return fmt.Sprintf("管理员已将本组退回到 %d 年 %s，请核对保留草稿并从当前阶段重新提交。", *yearState.RollbackTargetYearNo, targetStage)
}
