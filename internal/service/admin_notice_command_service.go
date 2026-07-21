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
)

const (
	noticeScopeAll                   = "ALL"
	noticeScopeGroup                 = "GROUP"
	adjustmentTypeReward             = "REWARD"
	adjustmentTypePenalty            = "PENALTY"
	noticeKindGeneral                = "GENERAL"
	adminActionCodeSendGeneralNotice = "SEND_GENERAL_NOTICE"
	adminActionCodeSendAdjustment    = "SEND_ADJUSTMENT"
	adminActionCodeVoidAdjustment    = "VOID_ADJUSTMENT"
	timeLayoutRFC3339                = time.RFC3339
)

var (
	ErrAdminNoticeContentRequired        = errors.New("admin notice content required")
	ErrAdminNoticeTargetScopeInvalid     = errors.New("admin notice target scope invalid")
	ErrAdminNoticeTargetGroupRequired    = errors.New("admin notice target group required")
	ErrAdminAdjustmentStageInvalid       = errors.New("admin adjustment stage invalid")
	ErrAdminAdjustmentTypeInvalid        = errors.New("admin adjustment type invalid")
	ErrAdminAdjustmentAmountInvalid      = errors.New("admin adjustment amount invalid")
	ErrAdminAdjustmentReasonRequired     = errors.New("admin adjustment reason required")
	ErrAdminAdjustmentYearNotOpen        = errors.New("admin adjustment year not open")
	ErrAdminAdjustmentStageLocked        = errors.New("admin adjustment stage locked")
	ErrAdminAdjustmentGroupNotAvailable  = errors.New("admin adjustment group not available")
	ErrAdminAdjustmentOperationInvalid   = errors.New("admin adjustment operation invalid")
	ErrAdminAdjustmentNotEffective       = errors.New("admin adjustment not effective")
	ErrAdminAdjustmentVoidReasonRequired = errors.New("admin adjustment void reason required")
	ErrAdminAdjustmentTargetMismatch     = errors.New("admin adjustment target mismatch")
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
	StageCode      string // 兼容旧调用方，服务端不会使用该值决定归属阶段。
	AdjustmentType string
	Amount         float64
	Reason         string
	OperatorID     int64
	OperatorName   string
}

type SendAdjustmentResult struct {
	AdjustmentID       int64                  `json:"adjustmentId"`
	GroupID            int64                  `json:"groupId"`
	YearNo             int                    `json:"yearNo"`
	StageCode          string                 `json:"stageCode"`
	AdjustmentType     string                 `json:"adjustmentType"`
	Amount             float64                `json:"amount"`
	Reason             string                 `json:"reason"`
	PublishedAt        time.Time              `json:"publishedAt"`
	AdjustmentRevision int64                  `json:"adjustmentRevision"`
	Impact             AdjustmentImpactResult `json:"impact"`
}

type PreviewAdjustmentCommand struct {
	Operation      string
	AdjustmentID   int64
	GroupID        int64
	YearNo         int
	AdjustmentType string
	Amount         float64
	Reason         string
}

type VoidAdjustmentCommand struct {
	AdjustmentID int64
	Reason       string
	OperatorID   int64
	OperatorName string
}

type VoidAdjustmentResult struct {
	AdjustmentID       int64                  `json:"adjustmentId"`
	GroupID            int64                  `json:"groupId"`
	YearNo             int                    `json:"yearNo"`
	AdjustmentRevision int64                  `json:"adjustmentRevision"`
	VoidedAt           time.Time              `json:"voidedAt"`
	Impact             AdjustmentImpactResult `json:"impact"`
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
	adjustmentType := strings.ToUpper(strings.TrimSpace(cmd.AdjustmentType))
	reason := strings.TrimSpace(cmd.Reason)
	if err := validateAdjustmentInput(adjustmentType, cmd.Amount, reason); err != nil {
		return nil, err
	}

	publishedAt := time.Now()
	var result *SendAdjustmentResult

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGroupRepo := repository.NewGroupRepository(tx)
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txAdjustmentRepo := repository.NewGroupAdjustmentRepository(tx)
		txRevisionRepo := repository.NewGroupAdjustmentRevisionRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		group, err := txGroupRepo.GetByIDForUpdate(ctx, cmd.GroupID)
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

		yearState, err := txGroupYearRepo.GetByGroupIDAndYearForUpdate(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return err
		}
		stageCode, err := resolveAdjustmentStage(*yearState)
		if err != nil {
			return err
		}

		item := &entity.GroupAdjustment{
			GroupID:        cmd.GroupID,
			YearNo:         cmd.YearNo,
			StageCode:      stageCode,
			AdjustmentType: adjustmentType,
			Amount:         cmd.Amount,
			Reason:         reason,
			Effective:      true,
			PublishedAt:    publishedAt,
			OperatorID:     cmd.OperatorID,
			OperatorName:   cmd.OperatorName,
			BaseEntity: entity.BaseEntity{
				Creator: cmd.OperatorName, CreateTime: publishedAt,
				Updater: cmd.OperatorName, UpdateTime: publishedAt,
			},
		}
		impact, err := newAdjustmentImpactCalculator(tx).Calculate(ctx, AdjustmentImpactRequest{
			Operation: adjustmentOperationCreate,
			GroupID:   cmd.GroupID, YearNo: cmd.YearNo, Candidate: item,
		})
		if err != nil {
			return err
		}

		if err := txAdjustmentRepo.Create(ctx, item); err != nil {
			return err
		}
		revision, err := txRevisionRepo.Increment(ctx, cmd.GroupID, cmd.YearNo, publishedAt)
		if err != nil {
			return err
		}

		if impact.WillBankrupt {
			if err := markAdjustmentBankrupt(ctx, tx, *item, *impact, cmd.OperatorID, cmd.OperatorName, publishedAt, adjustmentOperationCreate); err != nil {
				return err
			}
		}

		targetYearNo := cmd.YearNo
		payloadJSON, err := marshalJSON(map[string]any{
			"groupId":            cmd.GroupID,
			"yearNo":             cmd.YearNo,
			"stageCode":          stageCode,
			"adjustmentType":     adjustmentType,
			"amount":             cmd.Amount,
			"reason":             reason,
			"adjustmentId":       item.ID,
			"adjustmentRevision": revision,
			"impact":             impact,
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

		result = &SendAdjustmentResult{
			AdjustmentID: item.ID, GroupID: item.GroupID, YearNo: item.YearNo,
			StageCode: item.StageCode, AdjustmentType: item.AdjustmentType,
			Amount: item.Amount, Reason: item.Reason, PublishedAt: item.PublishedAt,
			AdjustmentRevision: revision, Impact: *impact,
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("send adjustment transaction: %w", err)
	}

	return result, nil
}

func (s *AdminNoticeCommandService) PreviewAdjustment(ctx context.Context, cmd PreviewAdjustmentCommand) (*AdjustmentImpactResult, error) {
	operation := strings.ToUpper(strings.TrimSpace(cmd.Operation))
	switch operation {
	case adjustmentOperationCreate:
		adjustmentType := strings.ToUpper(strings.TrimSpace(cmd.AdjustmentType))
		reason := strings.TrimSpace(cmd.Reason)
		if err := validateAdjustmentInput(adjustmentType, cmd.Amount, reason); err != nil {
			return nil, err
		}
		candidate := &entity.GroupAdjustment{
			GroupID: cmd.GroupID, YearNo: cmd.YearNo,
			AdjustmentType: adjustmentType, Amount: cmd.Amount, Reason: reason, Effective: true,
		}
		return newAdjustmentImpactCalculator(s.db).Calculate(ctx, AdjustmentImpactRequest{
			Operation: operation, GroupID: cmd.GroupID, YearNo: cmd.YearNo, Candidate: candidate,
		})
	case adjustmentOperationVoid:
		item, err := repository.NewGroupAdjustmentRepository(s.db).GetByID(ctx, cmd.AdjustmentID)
		if err != nil {
			return nil, err
		}
		if !item.Effective {
			return nil, ErrAdminAdjustmentNotEffective
		}
		return newAdjustmentImpactCalculator(s.db).Calculate(ctx, AdjustmentImpactRequest{
			Operation: operation, GroupID: item.GroupID, YearNo: item.YearNo, VoidID: item.ID,
		})
	default:
		return nil, ErrAdminAdjustmentOperationInvalid
	}
}

func (s *AdminNoticeCommandService) VoidAdjustment(ctx context.Context, cmd VoidAdjustmentCommand) (*VoidAdjustmentResult, error) {
	reason := strings.TrimSpace(cmd.Reason)
	if reason == "" {
		return nil, ErrAdminAdjustmentVoidReasonRequired
	}
	voidedAt := time.Now()
	var result *VoidAdjustmentResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		adjustmentRepo := repository.NewGroupAdjustmentRepository(tx)
		item, err := adjustmentRepo.GetByIDForUpdate(ctx, cmd.AdjustmentID)
		if err != nil {
			return err
		}
		if !item.Effective {
			return ErrAdminAdjustmentNotEffective
		}
		group, err := repository.NewGroupRepository(tx).GetByIDForUpdate(ctx, item.GroupID)
		if err != nil {
			return err
		}
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			return ErrAdminAdjustmentGroupNotAvailable
		}
		yearState, err := repository.NewGroupYearStateRepository(tx).GetByGroupIDAndYearForUpdate(ctx, item.GroupID, item.YearNo)
		if err != nil {
			return err
		}
		if _, err := resolveAdjustmentStage(*yearState); err != nil {
			return err
		}
		impact, err := newAdjustmentImpactCalculator(tx).Calculate(ctx, AdjustmentImpactRequest{
			Operation: adjustmentOperationVoid, GroupID: item.GroupID, YearNo: item.YearNo, VoidID: item.ID,
		})
		if err != nil {
			return err
		}
		updated, err := adjustmentRepo.Void(ctx, item.ID, cmd.OperatorID, cmd.OperatorName, reason, voidedAt)
		if err != nil {
			return err
		}
		if !updated {
			return ErrAdminAdjustmentNotEffective
		}
		revision, err := repository.NewGroupAdjustmentRevisionRepository(tx).Increment(ctx, item.GroupID, item.YearNo, voidedAt)
		if err != nil {
			return err
		}
		if impact.WillBankrupt {
			if err := markAdjustmentBankrupt(ctx, tx, *item, *impact, cmd.OperatorID, cmd.OperatorName, voidedAt, adjustmentOperationVoid); err != nil {
				return err
			}
		}
		payloadJSON, err := marshalJSON(map[string]any{
			"adjustmentId": item.ID, "reason": reason,
			"adjustmentRevision": revision, "impact": impact,
		})
		if err != nil {
			return err
		}
		if err := repository.NewAdminActionLogRepository(tx).Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeVoidAdjustment,
			TargetGroupID: &item.GroupID, TargetYearNo: &item.YearNo,
			ActionPayload: payloadJSON, StateBefore: []byte("{}"), StateAfter: []byte("{}"),
			OperatorID: cmd.OperatorID, OperatorName: cmd.OperatorName, OperateTime: voidedAt,
		}); err != nil {
			return err
		}
		result = &VoidAdjustmentResult{
			AdjustmentID: item.ID, GroupID: item.GroupID, YearNo: item.YearNo,
			AdjustmentRevision: revision, VoidedAt: voidedAt, Impact: *impact,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("void adjustment transaction: %w", err)
	}
	return result, nil
}

func validateAdjustmentInput(adjustmentType string, amount float64, reason string) error {
	if adjustmentType != adjustmentTypeReward && adjustmentType != adjustmentTypePenalty {
		return ErrAdminAdjustmentTypeInvalid
	}
	if amount <= 0 {
		return ErrAdminAdjustmentAmountInvalid
	}
	if !isWholeNumber(amount) {
		return ErrManualNumberNotInteger
	}
	if strings.TrimSpace(reason) == "" {
		return ErrAdminAdjustmentReasonRequired
	}
	return nil
}

func marshalJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal json payload: %w", err)
	}
	return raw, nil
}
