package service

import (
	"context"
	"fmt"
	"time"

	"sandbox-game/internal/enum"
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
	ID                 int64   `json:"id"`
	GroupID            int64   `json:"groupId"`
	GroupNo            int     `json:"groupNo"`
	GroupName          string  `json:"groupName"`
	YearNo             int     `json:"yearNo"`
	StageCode          string  `json:"stageCode"`
	AdjustmentType     string  `json:"adjustmentType"`
	Amount             float64 `json:"amount"`
	Reason             string  `json:"reason"`
	PublishedAt        string  `json:"publishedAt"`
	OperatorName       string  `json:"operatorName"`
	Status             string  `json:"status"`
	CanVoid            bool    `json:"canVoid"`
	AdjustmentRevision int64   `json:"adjustmentRevision"`
	VoidedByName       *string `json:"voidedByName"`
	VoidReason         *string `json:"voidReason"`
	VoidedAt           *string `json:"voidedAt"`
}

type AdminNoticeRecordsResult struct {
	GeneralNotices []AdminGeneralNoticeRecord `json:"generalNotices"`
	Adjustments    []AdminAdjustmentRecord    `json:"adjustments"`
}

type AdminNoticeQueryService struct {
	noticeRepo     *repository.NoticeRepository
	adjustmentRepo *repository.GroupAdjustmentRepository
	groupRepo      *repository.GroupRepository
	revisionRepo   *repository.GroupAdjustmentRevisionRepository
}

func NewAdminNoticeQueryService(
	noticeRepo *repository.NoticeRepository,
	adjustmentRepo *repository.GroupAdjustmentRepository,
	groupRepo *repository.GroupRepository,
	revisionRepo *repository.GroupAdjustmentRevisionRepository,
) *AdminNoticeQueryService {
	return &AdminNoticeQueryService{
		noticeRepo:     noticeRepo,
		adjustmentRepo: adjustmentRepo,
		groupRepo:      groupRepo,
		revisionRepo:   revisionRepo,
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
	groupBusinessStatusMap := map[int64]string{}
	for _, item := range groups {
		groupNameMap[item.ID] = item.GroupName
		groupNoMap[item.ID] = item.GroupNo
		groupBusinessStatusMap[item.ID] = item.BusinessStatus
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
		revision, revisionErr := s.revisionRepo.Get(ctx, item.GroupID, item.YearNo)
		if revisionErr != nil {
			return nil, fmt.Errorf("load adjustment revision: %w", revisionErr)
		}
		status := "SNAPSHOT_INACTIVE"
		if item.Effective {
			status = "EFFECTIVE"
		} else if item.VoidedAt != nil {
			status = "VOIDED"
		}
		adjustmentRecords = append(adjustmentRecords, AdminAdjustmentRecord{
			ID:                 item.ID,
			GroupID:            item.GroupID,
			GroupNo:            groupNoMap[item.GroupID],
			GroupName:          groupNameMap[item.GroupID],
			YearNo:             item.YearNo,
			StageCode:          item.StageCode,
			AdjustmentType:     item.AdjustmentType,
			Amount:             item.Amount,
			Reason:             item.Reason,
			PublishedAt:        item.PublishedAt.Format(timeLayoutRFC3339),
			OperatorName:       item.OperatorName,
			Status:             status,
			CanVoid:            item.Effective && groupBusinessStatusMap[item.GroupID] != enum.BusinessStatusBankrupt,
			AdjustmentRevision: revision,
			VoidedByName:       item.VoidedByName,
			VoidReason:         item.VoidReason,
			VoidedAt:           formatAdjustmentTime(item.VoidedAt),
		})
	}

	return &AdminNoticeRecordsResult{
		GeneralNotices: noticeRecords,
		Adjustments:    adjustmentRecords,
	}, nil
}

func formatAdjustmentTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(timeLayoutRFC3339)
	return &formatted
}
