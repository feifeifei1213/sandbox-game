package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

type PlayerAdjustmentSyncResult struct {
	NotModified           bool                           `json:"notModified"`
	AdjustmentRevision    int64                          `json:"adjustmentRevision"`
	IncomeAndPenalty      payload.OperatingQuarterMap    `json:"incomeAndPenalty,omitempty"`
	DerivedValues         map[string]float64             `json:"derivedValues,omitempty"`
	QuarterCashChecks     map[string]float64             `json:"quarterCashChecks,omitempty"`
	PeriodEndCash         float64                        `json:"periodEndCash,omitempty"`
	ReportComputedPayload *payload.ReportComputedPayload `json:"reportComputedPayload,omitempty"`
	BalanceGap            float64                        `json:"balanceGap,omitempty"`
	BalanceCheckPassed    bool                           `json:"balanceCheckPassed,omitempty"`
	NoticeBoard           *assembler.PlayerNoticeBoard   `json:"noticeBoard,omitempty"`
	BusinessStatus        string                         `json:"businessStatus,omitempty"`
	Bankrupt              bool                           `json:"bankrupt,omitempty"`
}

type PlayerAdjustmentSyncService struct {
	db                  *gorm.DB
	playerNoticeService *PlayerNoticeService
}

func NewPlayerAdjustmentSyncService(db *gorm.DB, playerNoticeService *PlayerNoticeService) *PlayerAdjustmentSyncService {
	return &PlayerAdjustmentSyncService{db: db, playerNoticeService: playerNoticeService}
}

func (s *PlayerAdjustmentSyncService) GetSync(ctx context.Context, groupID int64, yearNo int, knownRevision int64) (*PlayerAdjustmentSyncResult, error) {
	revision, err := repository.NewGroupAdjustmentRevisionRepository(s.db).Get(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("load adjustment revision: %w", err)
	}
	if revision == knownRevision {
		return &PlayerAdjustmentSyncResult{NotModified: true, AdjustmentRevision: revision}, nil
	}

	impact, err := newAdjustmentImpactCalculator(s.db).Calculate(ctx, AdjustmentImpactRequest{
		Operation: adjustmentOperationCurrent, GroupID: groupID, YearNo: yearNo,
	})
	if err != nil {
		return nil, err
	}
	board, err := s.playerNoticeService.BuildBoard(ctx, groupID, yearNo)
	if err != nil {
		return nil, err
	}
	group, err := repository.NewGroupRepository(s.db).GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	report := impact.ReportAfter
	return &PlayerAdjustmentSyncResult{
		NotModified: false, AdjustmentRevision: revision,
		IncomeAndPenalty:      impact.OperatingAfter.OperatingPayload.Extra.IncomeAndPenalty,
		DerivedValues:         impact.OperatingAfter.DerivedValues,
		QuarterCashChecks:     impact.OperatingAfter.QuarterCashChecks,
		PeriodEndCash:         impact.OperatingAfter.PeriodEndCash,
		ReportComputedPayload: &report,
		BalanceGap:            report.BalanceGap(), BalanceCheckPassed: reportBalancePassed(report),
		NoticeBoard: board, BusinessStatus: group.BusinessStatus,
		Bankrupt: group.BusinessStatus == enum.BusinessStatusBankrupt,
	}, nil
}
