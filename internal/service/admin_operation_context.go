package service

import (
	"context"
	"errors"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

type AdminGroupOperationContextResult struct {
	GroupID                 int64   `json:"groupId"`
	GroupNo                 int     `json:"groupNo"`
	GroupName               string  `json:"groupName"`
	BusinessStatus          string  `json:"businessStatus"`
	CurrentOpenYear         int     `json:"currentOpenYear"`
	OperationYearNo         int     `json:"operationYearNo"`
	YearStatus              string  `json:"yearStatus"`
	StageStatus             string  `json:"stageStatus"`
	ReportStatus            string  `json:"reportStatus"`
	CurrentStageCode        string  `json:"currentStageCode"`
	RollbackPending         bool    `json:"rollbackPending"`
	RollbackTargetYearNo    *int    `json:"rollbackTargetYearNo"`
	RollbackTargetStageCode *string `json:"rollbackTargetStageCode"`
	CanUnlockRetry          bool    `json:"canUnlockRetry"`
	CanAdjust               bool    `json:"canAdjust"`
	AdjustmentStageCode     string  `json:"adjustmentStageCode"`
	BlockedReason           string  `json:"blockedReason"`
}

type adminGroupOperationContext struct {
	Group           entity.Group
	GameConfig      entity.GameConfig
	YearState       entity.GroupYearState
	OperationYearNo int
}

func loadAdminGroupOperationContext(
	ctx context.Context,
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	groupID int64,
	forUpdate bool,
) (*adminGroupOperationContext, error) {
	var group *entity.Group
	var err error
	if forUpdate {
		group, err = groupRepo.GetByIDForUpdate(ctx, groupID)
	} else {
		group, err = groupRepo.GetByID(ctx, groupID)
	}
	if err != nil {
		return nil, err
	}

	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, err
	}

	yearStates, err := groupYearRepo.ListByGroupIDFromYear(ctx, groupID, 0)
	if err != nil {
		return nil, err
	}
	operationYearNo := resolveAdminOperationYearNo(yearStates, gameConfig.CurrentOpenYear)

	var yearState *entity.GroupYearState
	if forUpdate {
		yearState, err = groupYearRepo.GetByGroupIDAndYearForUpdate(ctx, groupID, operationYearNo)
	} else {
		yearState, err = groupYearRepo.GetByGroupIDAndYear(ctx, groupID, operationYearNo)
	}
	if err != nil {
		return nil, err
	}

	return &adminGroupOperationContext{
		Group:           *group,
		GameConfig:      *gameConfig,
		YearState:       *yearState,
		OperationYearNo: operationYearNo,
	}, nil
}

func resolveAdminOperationYearNo(yearStates []entity.GroupYearState, currentOpenYear int) int {
	for _, item := range yearStates {
		if item.RollbackPending {
			return item.YearNo
		}
	}
	return currentOpenYear
}

func buildAdminGroupOperationContextResult(ctx adminGroupOperationContext) AdminGroupOperationContextResult {
	currentStageCode := state.CurrentStageCode(ctx.YearState.StageStatus)
	adjustmentStageCode := ""
	canAdjust := false
	blockedReason := ""

	if ctx.Group.BusinessStatus == enum.BusinessStatusBankrupt {
		blockedReason = "目标小组已破产，不能继续操作"
	} else if stageCode, err := resolveAdjustmentStage(ctx.YearState); err == nil {
		adjustmentStageCode = stageCode
		canAdjust = true
	} else {
		blockedReason = resolveAdminOperationBlockedReason(err)
	}

	nextYearAlreadyOpened := ctx.GameConfig.CurrentOpenYear > ctx.OperationYearNo && !ctx.YearState.RollbackPending
	canUnlockRetry := ctx.Group.BusinessStatus != enum.BusinessStatusBankrupt &&
		ctx.YearState.YearStatus != enum.YearStatusLocked &&
		!nextYearAlreadyOpened

	return AdminGroupOperationContextResult{
		GroupID:                 ctx.Group.ID,
		GroupNo:                 ctx.Group.GroupNo,
		GroupName:               ctx.Group.GroupName,
		BusinessStatus:          ctx.Group.BusinessStatus,
		CurrentOpenYear:         ctx.GameConfig.CurrentOpenYear,
		OperationYearNo:         ctx.OperationYearNo,
		YearStatus:              ctx.YearState.YearStatus,
		StageStatus:             ctx.YearState.StageStatus,
		ReportStatus:            ctx.YearState.ReportStatus,
		CurrentStageCode:        currentStageCode,
		RollbackPending:         ctx.YearState.RollbackPending,
		RollbackTargetYearNo:    ctx.YearState.RollbackTargetYearNo,
		RollbackTargetStageCode: ctx.YearState.RollbackTargetStageCode,
		CanUnlockRetry:          canUnlockRetry,
		CanAdjust:               canAdjust,
		AdjustmentStageCode:     adjustmentStageCode,
		BlockedReason:           blockedReason,
	}
}

func resolveAdminOperationBlockedReason(err error) string {
	switch {
	case errors.Is(err, ErrAdminAdjustmentYearNotOpen):
		return "目标小组当前年份未开放"
	case errors.Is(err, ErrAdminAdjustmentStageLocked):
		return "目标小组当前阶段已锁定，请先退回到可编辑状态"
	case errors.Is(err, ErrAdminAdjustmentStageInvalid):
		return "目标小组当前阶段无法自动归属奖惩"
	default:
		return "目标小组当前状态暂不可操作"
	}
}
