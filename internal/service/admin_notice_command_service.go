package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

const (
	noticeScopeAll                   = "ALL"
	noticeScopeGroup                 = "GROUP"
	adjustmentTypeReward             = "REWARD"
	adjustmentTypePenalty            = "PENALTY"
	noticeKindGeneral                = "GENERAL"
	adminActionCodeSendGeneralNotice = "SEND_GENERAL_NOTICE"
	adminActionCodeSendAdjustment    = "SEND_ADJUSTMENT"
	timeLayoutRFC3339                = time.RFC3339
)

var (
	ErrAdminNoticeContentRequired       = errors.New("admin notice content required")
	ErrAdminNoticeTargetScopeInvalid    = errors.New("admin notice target scope invalid")
	ErrAdminNoticeTargetGroupRequired   = errors.New("admin notice target group required")
	ErrAdminAdjustmentStageInvalid      = errors.New("admin adjustment stage invalid")
	ErrAdminAdjustmentTypeInvalid       = errors.New("admin adjustment type invalid")
	ErrAdminAdjustmentAmountInvalid     = errors.New("admin adjustment amount invalid")
	ErrAdminAdjustmentReasonRequired    = errors.New("admin adjustment reason required")
	ErrAdminAdjustmentYearNotOpen       = errors.New("admin adjustment year not open")
	ErrAdminAdjustmentStageLocked       = errors.New("admin adjustment stage locked")
	ErrAdminAdjustmentGroupNotAvailable = errors.New("admin adjustment group not available")
)

type SendGeneralNoticeCommand struct {
	TargetScope   string
	TargetGroupID *int64
	Content       string
	Pinned        bool
	OperatorID    int64
	OperatorName  string
}

type SendGeneralNoticeResult struct {
	NoticeID      int64     `json:"noticeId"`
	TargetScope   string    `json:"targetScope"`
	TargetGroupID *int64    `json:"targetGroupId"`
	Content       string    `json:"content"`
	Pinned        bool      `json:"pinned"`
	PublishedAt   time.Time `json:"publishedAt"`
}

type SendAdjustmentCommand struct {
	GroupID        int64
	YearNo         int
	StageCode      string
	AdjustmentType string
	Amount         float64
	Reason         string
	OperatorID     int64
	OperatorName   string
}

type SendAdjustmentResult struct {
	AdjustmentID   int64     `json:"adjustmentId"`
	GroupID        int64     `json:"groupId"`
	YearNo         int       `json:"yearNo"`
	StageCode      string    `json:"stageCode"`
	AdjustmentType string    `json:"adjustmentType"`
	Amount         float64   `json:"amount"`
	Reason         string    `json:"reason"`
	PublishedAt    time.Time `json:"publishedAt"`
}

type AdminNoticeCommandService struct {
	db *gorm.DB
}

func NewAdminNoticeCommandService(db *gorm.DB) *AdminNoticeCommandService {
	return &AdminNoticeCommandService{db: db}
}

func (s *AdminNoticeCommandService) SendGeneralNotice(
	ctx context.Context,
	cmd SendGeneralNoticeCommand,
) (*SendGeneralNoticeResult, error) {
	targetScope := strings.ToUpper(strings.TrimSpace(cmd.TargetScope))
	content := strings.TrimSpace(cmd.Content)
	if content == "" {
		return nil, ErrAdminNoticeContentRequired
	}
	if targetScope != noticeScopeAll && targetScope != noticeScopeGroup {
		return nil, ErrAdminNoticeTargetScopeInvalid
	}
	if targetScope == noticeScopeGroup && (cmd.TargetGroupID == nil || *cmd.TargetGroupID <= 0) {
		return nil, ErrAdminNoticeTargetGroupRequired
	}

	publishedAt := time.Now()
	item := &entity.Notice{
		TargetScope:   targetScope,
		TargetGroupID: cmd.TargetGroupID,
		Content:       content,
		Pinned:        cmd.Pinned,
		PublishedAt:   publishedAt,
		OperatorID:    cmd.OperatorID,
		OperatorName:  cmd.OperatorName,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: publishedAt,
			Updater:    cmd.OperatorName,
			UpdateTime: publishedAt,
		},
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txNoticeRepo := repository.NewNoticeRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		if targetScope == noticeScopeGroup {
			if _, err := txGroupRepo.GetByID(ctx, *cmd.TargetGroupID); err != nil {
				return err
			}
		}

		if err := txNoticeRepo.Create(ctx, item); err != nil {
			return err
		}

		payloadJSON, err := marshalJSON(map[string]any{
			"targetScope":   targetScope,
			"targetGroupId": cmd.TargetGroupID,
			"content":       content,
			"pinned":        cmd.Pinned,
			"noticeId":      item.ID,
		})
		if err != nil {
			return err
		}

		if err := txAdminActionLogRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeSendGeneralNotice,
			TargetGroupID: cmd.TargetGroupID,
			TargetYearNo:  nil,
			ActionPayload: payloadJSON,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  cmd.OperatorName,
			OperateTime:   publishedAt,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("send general notice transaction: %w", err)
	}

	return &SendGeneralNoticeResult{
		NoticeID:      item.ID,
		TargetScope:   item.TargetScope,
		TargetGroupID: item.TargetGroupID,
		Content:       item.Content,
		Pinned:        item.Pinned,
		PublishedAt:   item.PublishedAt,
	}, nil
}

func (s *AdminNoticeCommandService) SendAdjustment(
	ctx context.Context,
	cmd SendAdjustmentCommand,
) (*SendAdjustmentResult, error) {
	stageCode := strings.ToUpper(strings.TrimSpace(cmd.StageCode))
	adjustmentType := strings.ToUpper(strings.TrimSpace(cmd.AdjustmentType))
	reason := strings.TrimSpace(cmd.Reason)

	if !isQuarterStageCode(stageCode) {
		return nil, ErrAdminAdjustmentStageInvalid
	}
	if adjustmentType != adjustmentTypeReward && adjustmentType != adjustmentTypePenalty {
		return nil, ErrAdminAdjustmentTypeInvalid
	}
	if cmd.Amount <= 0 {
		return nil, ErrAdminAdjustmentAmountInvalid
	}
	if reason == "" {
		return nil, ErrAdminAdjustmentReasonRequired
	}

	publishedAt := time.Now()
	item := &entity.GroupAdjustment{
		GroupID:        cmd.GroupID,
		YearNo:         cmd.YearNo,
		StageCode:      stageCode,
		AdjustmentType: adjustmentType,
		Amount:         cmd.Amount,
		Reason:         reason,
		PublishedAt:    publishedAt,
		OperatorID:     cmd.OperatorID,
		OperatorName:   cmd.OperatorName,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: publishedAt,
			Updater:    cmd.OperatorName,
			UpdateTime: publishedAt,
		},
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGroupRepo := repository.NewGroupRepository(tx)
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txOperatingRepo := repository.NewOperatingRepository(tx)
		txAdjustmentRepo := repository.NewGroupAdjustmentRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		group, err := txGroupRepo.GetByID(ctx, cmd.GroupID)
		if err != nil {
			return err
		}
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			return ErrAdminAdjustmentGroupNotAvailable
		}

		gameConfig, err := txGameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return err
		}
		if cmd.YearNo > gameConfig.CurrentOpenYear {
			return ErrAdminAdjustmentYearNotOpen
		}

		yearState, err := txGroupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return err
		}
		if yearState.YearStatus == enum.YearStatusLocked {
			return ErrAdminAdjustmentYearNotOpen
		}

		locked, err := txOperatingRepo.HasStageSubmission(ctx, cmd.GroupID, cmd.YearNo, stageCode)
		if err != nil {
			return err
		}
		if locked {
			return ErrAdminAdjustmentStageLocked
		}

		if err := txAdjustmentRepo.Create(ctx, item); err != nil {
			return err
		}

		targetYearNo := cmd.YearNo
		payloadJSON, err := marshalJSON(map[string]any{
			"groupId":        cmd.GroupID,
			"yearNo":         cmd.YearNo,
			"stageCode":      stageCode,
			"adjustmentType": adjustmentType,
			"amount":         cmd.Amount,
			"reason":         reason,
			"adjustmentId":   item.ID,
		})
		if err != nil {
			return err
		}

		if err := txAdminActionLogRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeSendAdjustment,
			TargetGroupID: &cmd.GroupID,
			TargetYearNo:  &targetYearNo,
			ActionPayload: payloadJSON,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  cmd.OperatorName,
			OperateTime:   publishedAt,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("send adjustment transaction: %w", err)
	}

	return &SendAdjustmentResult{
		AdjustmentID:   item.ID,
		GroupID:        item.GroupID,
		YearNo:         item.YearNo,
		StageCode:      item.StageCode,
		AdjustmentType: item.AdjustmentType,
		Amount:         item.Amount,
		Reason:         item.Reason,
		PublishedAt:    item.PublishedAt,
	}, nil
}

func isQuarterStageCode(stageCode string) bool {
	switch stageCode {
	case state.StageCodeQ1, state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4:
		return true
	default:
		return false
	}
}

func marshalJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal json payload: %w", err)
	}
	return raw, nil
}
