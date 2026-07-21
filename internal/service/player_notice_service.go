package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"sandbox-game/internal/assembler"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

const playerNoticeRecentLimit = 8

type PlayerNoticeService struct {
	noticeRepo     *repository.NoticeRepository
	adjustmentRepo *repository.GroupAdjustmentRepository
	revisionRepo   *repository.GroupAdjustmentRevisionRepository
}

func NewPlayerNoticeService(
	noticeRepo *repository.NoticeRepository,
	adjustmentRepo *repository.GroupAdjustmentRepository,
	revisionRepos ...*repository.GroupAdjustmentRevisionRepository,
) *PlayerNoticeService {
	var revisionRepo *repository.GroupAdjustmentRevisionRepository
	if len(revisionRepos) > 0 {
		revisionRepo = revisionRepos[0]
	}
	return &PlayerNoticeService{
		noticeRepo:     noticeRepo,
		adjustmentRepo: adjustmentRepo,
		revisionRepo:   revisionRepo,
	}
}

func (s *PlayerNoticeService) GetRevision(ctx context.Context, groupID int64, yearNo int) (int64, error) {
	if s.revisionRepo == nil {
		return 0, nil
	}
	return s.revisionRepo.Get(ctx, groupID, yearNo)
}

func (s *PlayerNoticeService) OverlayAdjustments(
	ctx context.Context,
	groupID int64,
	yearNo int,
	operatingPayload payload.OperatingPayload,
) (payload.OperatingPayload, error) {
	adjustments, err := s.adjustmentRepo.ListByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return payload.OperatingPayload{}, fmt.Errorf("list group adjustments: %w", err)
	}

	return overlayAdjustmentValues(operatingPayload, adjustments), nil
}

func (s *PlayerNoticeService) BuildBoard(ctx context.Context, groupID int64, yearNo int) (*assembler.PlayerNoticeBoard, error) {
	notices, err := s.noticeRepo.ListForGroup(ctx, groupID, playerNoticeRecentLimit)
	if err != nil {
		return nil, fmt.Errorf("list notices for group: %w", err)
	}

	adjustments, err := s.adjustmentRepo.ListAllByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list adjustments for board: %w", err)
	}

	return buildPlayerNoticeBoard(notices, adjustments), nil
}

func overlayAdjustmentValues(
	operatingPayload payload.OperatingPayload,
	adjustments []entity.GroupAdjustment,
) payload.OperatingPayload {
	normalized := operatingPayload.Normalize()
	rewardByQuarter := map[string]float64{
		"q1": 0, "q2": 0, "q3": 0, "q4": 0, "year_end": 0,
	}
	penaltyByQuarter := map[string]float64{
		"q1": 0, "q2": 0, "q3": 0, "q4": 0, "year_end": 0,
	}

	for _, item := range adjustments {
		quarterKey := strings.ToLower(strings.TrimSpace(item.StageCode))
		switch item.AdjustmentType {
		case adjustmentTypeReward:
			rewardByQuarter[quarterKey] += item.Amount
		case adjustmentTypePenalty:
			penaltyByQuarter[quarterKey] += item.Amount
		}
	}

	for _, quarterKey := range []string{"q1", "q2", "q3", "q4", "year_end"} {
		normalized.Extra.IncomeAndPenalty[quarterKey] = ensureQuarterValueMap(normalized.Extra.IncomeAndPenalty[quarterKey])
		normalized.Extra.IncomeAndPenalty[quarterKey]["extraIncomeReward"] = rewardByQuarter[quarterKey]
		normalized.Extra.IncomeAndPenalty[quarterKey]["extraExpensePenalty"] = penaltyByQuarter[quarterKey]
	}

	return normalized
}

func buildPlayerNoticeBoard(
	notices []entity.Notice,
	adjustments []entity.GroupAdjustment,
) *assembler.PlayerNoticeBoard {
	if len(notices) == 0 && len(adjustments) == 0 {
		return &assembler.PlayerNoticeBoard{
			PinnedNotice: nil,
			RecentList:   []assembler.PlayerNoticeItem{},
		}
	}

	var pinnedNotice *assembler.PlayerNoticeItem
	items := make([]assembler.PlayerNoticeItem, 0, len(notices)+len(adjustments))

	for _, item := range notices {
		noticeItem := assembler.PlayerNoticeItem{
			ID:          item.ID,
			Kind:        noticeKindGeneral,
			Title:       "系统通知",
			Content:     item.Content,
			Pinned:      item.Pinned,
			PublishedAt: item.PublishedAt,
		}
		if item.Pinned && pinnedNotice == nil {
			copied := noticeItem
			pinnedNotice = &copied
			continue
		}
		items = append(items, noticeItem)
	}

	for _, item := range adjustments {
		items = append(items, buildAdjustmentNoticeItem(item))
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].PublishedAt.Equal(items[j].PublishedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].PublishedAt.After(items[j].PublishedAt)
	})

	if len(items) > playerNoticeRecentLimit {
		items = items[:playerNoticeRecentLimit]
	}

	return &assembler.PlayerNoticeBoard{
		PinnedNotice: pinnedNotice,
		RecentList:   items,
	}
}

func buildAdjustmentNoticeItem(item entity.GroupAdjustment) assembler.PlayerNoticeItem {
	stageCode := item.StageCode
	yearNo := item.YearNo
	amount := item.Amount
	status := adjustmentRecordStatus(item)
	title := fmt.Sprintf("%d年 %s %s", item.YearNo, item.StageCode, translateAdjustmentType(item.AdjustmentType))

	return assembler.PlayerNoticeItem{
		ID:          item.ID,
		Kind:        item.AdjustmentType,
		Title:       title,
		Content:     item.Reason,
		Pinned:      false,
		PublishedAt: item.PublishedAt,
		YearNo:      &yearNo,
		StageCode:   &stageCode,
		Amount:      &amount,
		Status:      &status,
	}
}

func adjustmentRecordStatus(item entity.GroupAdjustment) string {
	if item.Effective {
		return "EFFECTIVE"
	}
	if item.VoidedAt != nil {
		return "VOIDED"
	}
	return "SNAPSHOT_INACTIVE"
}

func ensureQuarterValueMap(source map[string]any) map[string]any {
	if source == nil {
		return map[string]any{}
	}
	return source
}

func translateAdjustmentType(value string) string {
	switch value {
	case adjustmentTypeReward:
		return "奖励"
	case adjustmentTypePenalty:
		return "罚款"
	default:
		return value
	}
}
