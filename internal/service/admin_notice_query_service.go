package service

import (
	"context"
	"fmt"

	"sandbox-game/internal/repository"
)

type AdminGeneralNoticeRecord struct {
	ID              int64   `json:"id"`
	TargetScope     string  `json:"targetScope"`
	TargetGroupID   *int64  `json:"targetGroupId"`
	TargetGroupName *string `json:"targetGroupName"`
	Content         string  `json:"content"`
	Pinned          bool    `json:"pinned"`
	PublishedAt     string  `json:"publishedAt"`
	OperatorName    string  `json:"operatorName"`
}

type AdminAdjustmentRecord struct {
	ID             int64   `json:"id"`
	GroupID        int64   `json:"groupId"`
	GroupNo        int     `json:"groupNo"`
	GroupName      string  `json:"groupName"`
	YearNo         int     `json:"yearNo"`
	StageCode      string  `json:"stageCode"`
	AdjustmentType string  `json:"adjustmentType"`
	Amount         float64 `json:"amount"`
	Reason         string  `json:"reason"`
	PublishedAt    string  `json:"publishedAt"`
	OperatorName   string  `json:"operatorName"`
}

type AdminNoticeRecordsResult struct {
	GeneralNotices []AdminGeneralNoticeRecord `json:"generalNotices"`
	Adjustments    []AdminAdjustmentRecord    `json:"adjustments"`
}

type AdminNoticeQueryService struct {
	noticeRepo     *repository.NoticeRepository
	adjustmentRepo *repository.GroupAdjustmentRepository
	groupRepo      *repository.GroupRepository
}

func NewAdminNoticeQueryService(
	noticeRepo *repository.NoticeRepository,
	adjustmentRepo *repository.GroupAdjustmentRepository,
	groupRepo *repository.GroupRepository,
) *AdminNoticeQueryService {
	return &AdminNoticeQueryService{
		noticeRepo:     noticeRepo,
		adjustmentRepo: adjustmentRepo,
		groupRepo:      groupRepo,
	}
}

func (s *AdminNoticeQueryService) GetRecords(ctx context.Context, limit int) (*AdminNoticeRecordsResult, error) {
	notices, err := s.noticeRepo.ListRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent notices: %w", err)
	}

	adjustments, err := s.adjustmentRepo.ListRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent adjustments: %w", err)
	}

	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list groups for notice records: %w", err)
	}

	groupNameMap := map[int64]string{}
	groupNoMap := map[int64]int{}
	for _, item := range groups {
		groupNameMap[item.ID] = item.GroupName
		groupNoMap[item.ID] = item.GroupNo
	}

	noticeRecords := make([]AdminGeneralNoticeRecord, 0, len(notices))
	for _, item := range notices {
		var targetGroupName *string
		if item.TargetGroupID != nil {
			if name, ok := groupNameMap[*item.TargetGroupID]; ok {
				copied := name
				targetGroupName = &copied
			}
		}

		noticeRecords = append(noticeRecords, AdminGeneralNoticeRecord{
			ID:              item.ID,
			TargetScope:     item.TargetScope,
			TargetGroupID:   item.TargetGroupID,
			TargetGroupName: targetGroupName,
			Content:         item.Content,
			Pinned:          item.Pinned,
			PublishedAt:     item.PublishedAt.Format(timeLayoutRFC3339),
			OperatorName:    item.OperatorName,
		})
	}

	adjustmentRecords := make([]AdminAdjustmentRecord, 0, len(adjustments))
	for _, item := range adjustments {
		adjustmentRecords = append(adjustmentRecords, AdminAdjustmentRecord{
			ID:             item.ID,
			GroupID:        item.GroupID,
			GroupNo:        groupNoMap[item.GroupID],
			GroupName:      groupNameMap[item.GroupID],
			YearNo:         item.YearNo,
			StageCode:      item.StageCode,
			AdjustmentType: item.AdjustmentType,
			Amount:         item.Amount,
			Reason:         item.Reason,
			PublishedAt:    item.PublishedAt.Format(timeLayoutRFC3339),
			OperatorName:   item.OperatorName,
		})
	}

	return &AdminNoticeRecordsResult{
		GeneralNotices: noticeRecords,
		Adjustments:    adjustmentRecords,
	}, nil
}
