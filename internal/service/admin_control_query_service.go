package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

const (
	openNextYearBlockedReasonFinalYearReached    = "已达到最终年份，无法继续开放"
	openNextYearBlockedReasonBaselineUnsubmitted = "初始基线尚未提交"
	openNextYearBlockedReasonUnfinishedReports   = "当前仍有未完成财报的小组"
	openNextYearBlockedReasonRollbackPending     = "仍有小组处于回退补提中，不能开放下一年"
)

type AdminActionSummary struct {
	ActionCode    string `json:"actionCode"`
	OperatorName  string `json:"operatorName"`
	OperateTime   string `json:"operateTime"`
	TargetGroupID *int64 `json:"targetGroupId,omitempty"`
	TargetYearNo  *int   `json:"targetYearNo,omitempty"`
}

type AdminControlConfigResult struct {
	FinalYear                  int                 `json:"finalYear"`
	CurrentOpenYear            int                 `json:"currentOpenYear"`
	EditionCode                string              `json:"editionCode"`
	EditionName                string              `json:"editionName"`
	CanOpenNextYear            bool                `json:"canOpenNextYear"`
	NextOpenableYear           int                 `json:"nextOpenableYear"`
	OpenNextYearBlockedReason  string              `json:"openNextYearBlockedReason"`
	RuleVersion                string              `json:"ruleVersion"`
	FormulaVersion             string              `json:"formulaVersion"`
	TemplateVersion            string              `json:"templateVersion"`
	OperatingTemplateVersion   string              `json:"operatingTemplateVersion"`
	ReportTemplateVersion      string              `json:"reportTemplateVersion"`
	OrderTemplateVersion       string              `json:"orderTemplateVersion"`
	ProcessRuleVersion         string              `json:"processRuleVersion"`
	InitialBaselineSubmitted   bool                `json:"initialBaselineSubmitted"`
	InitialBaselineSubmittedAt *string             `json:"initialBaselineSubmittedAt"`
	InitialBaselineSubmitter   *string             `json:"initialBaselineSubmitterName"`
	LatestAdminAction          *AdminActionSummary `json:"latestAdminAction"`
}

type AdminControlSetupStatusResult struct {
	Initialized              bool          `json:"initialized"`
	GroupCount               int           `json:"groupCount"`
	FinalYear                int           `json:"finalYear"`
	CurrentOpenYear          int           `json:"currentOpenYear"`
	EditionCode              string        `json:"editionCode"`
	EditionName              string        `json:"editionName"`
	RuleVersion              string        `json:"ruleVersion"`
	TemplateVersion          string        `json:"templateVersion"`
	AvailableEditions        []GameEdition `json:"availableEditions"`
	InitialBaselineSubmitted bool          `json:"initialBaselineSubmitted"`
	DefaultRoute             string        `json:"defaultRoute"`
}

type InitialBaselineViewResult struct {
	Submitted         bool                    `json:"submitted"`
	Editable          bool                    `json:"editable"`
	BaselinePayload   payload.BaselinePayload `json:"baselinePayload"`
	AppliedGroupCount int                     `json:"appliedGroupCount"`
	SubmitterName     *string                 `json:"submitterName"`
	SubmittedAt       *string                 `json:"submittedAt"`
}

type AdminControlQueryService struct {
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	initialBaselineRepo *repository.InitialBaselineRepository
	accountRepo         *repository.AccountRepository
	adminActionLogRepo  *repository.AdminActionLogRepository
}

func NewAdminControlQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	initialBaselineRepo *repository.InitialBaselineRepository,
	accountRepo *repository.AccountRepository,
	adminActionLogRepo *repository.AdminActionLogRepository,
) *AdminControlQueryService {
	return &AdminControlQueryService{
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		initialBaselineRepo: initialBaselineRepo,
		accountRepo:         accountRepo,
		adminActionLogRepo:  adminActionLogRepo,
	}
}

func (s *AdminControlQueryService) GetConfig(ctx context.Context) (*AdminControlConfigResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("load groups: %w", err)
	}
	states, err := s.groupYearRepo.ListByYear(ctx, gameConfig.CurrentOpenYear)
	if err != nil {
		return nil, fmt.Errorf("load current year states: %w", err)
	}

	openStatus := evaluateOpenNextYearStatus(gameConfig, groups, states)
	var baselineSubmittedAt *string
	var baselineSubmitterName *string
	if gameConfig.InitialBaselineSubmitted {
		baselineItem, baselineErr := s.initialBaselineRepo.FindLatestSubmitted(ctx)
		if baselineErr != nil {
			return nil, fmt.Errorf("load initial baseline summary: %w", baselineErr)
		}
		baselineSubmittedAt, baselineSubmitterName, err = s.extractInitialBaselineMeta(ctx, baselineItem)
		if err != nil {
			return nil, err
		}
	}
	latestAdminAction, err := s.lookupLatestAdminAction(ctx)
	if err != nil {
		return nil, err
	}
	edition := normalizeGameConfigEdition(*gameConfig)

	return &AdminControlConfigResult{
		FinalYear:                  gameConfig.FinalYear,
		CurrentOpenYear:            gameConfig.CurrentOpenYear,
		EditionCode:                edition.EditionCode,
		EditionName:                edition.EditionName,
		CanOpenNextYear:            openStatus.CanOpenNextYear,
		NextOpenableYear:           openStatus.NextOpenableYear,
		OpenNextYearBlockedReason:  openStatus.BlockedReason,
		RuleVersion:                edition.RuleVersion,
		FormulaVersion:             edition.FormulaVersion,
		TemplateVersion:            edition.TemplateVersion,
		OperatingTemplateVersion:   edition.OperatingTemplateVersion,
		ReportTemplateVersion:      edition.ReportTemplateVersion,
		OrderTemplateVersion:       edition.OrderTemplateVersion,
		ProcessRuleVersion:         edition.ProcessRuleVersion,
		InitialBaselineSubmitted:   gameConfig.InitialBaselineSubmitted,
		InitialBaselineSubmittedAt: baselineSubmittedAt,
		InitialBaselineSubmitter:   baselineSubmitterName,
		LatestAdminAction:          latestAdminAction,
	}, nil
}

func (s *AdminControlQueryService) GetSetupStatus(ctx context.Context) (*AdminControlSetupStatusResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	groupCount, err := s.groupRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("count groups: %w", err)
	}

	initialized := groupCount > 0
	defaultRoute := "/sandbox-game/admin/setup"
	if initialized {
		defaultRoute = "/sandbox-game/admin/summary"
	}
	edition := normalizeGameConfigEdition(*gameConfig)

	return &AdminControlSetupStatusResult{
		Initialized:              initialized,
		GroupCount:               int(groupCount),
		FinalYear:                gameConfig.FinalYear,
		CurrentOpenYear:          gameConfig.CurrentOpenYear,
		EditionCode:              edition.EditionCode,
		EditionName:              edition.EditionName,
		RuleVersion:              edition.RuleVersion,
		TemplateVersion:          edition.TemplateVersion,
		AvailableEditions:        ListGameEditions(),
		InitialBaselineSubmitted: gameConfig.InitialBaselineSubmitted,
		DefaultRoute:             defaultRoute,
	}, nil
}

func (s *AdminControlQueryService) GetInitialBaseline(ctx context.Context) (*InitialBaselineViewResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	baselinePayload := payload.BaselinePayload{}
	appliedGroupCount := 0
	var submitterName *string
	var submittedAtValue *time.Time

	if gameConfig.InitialBaselineSubmitted {
		baselineItem, baselineErr := s.initialBaselineRepo.FindLatestSubmitted(ctx)
		if baselineErr != nil {
			return nil, fmt.Errorf("load submitted initial baseline: %w", baselineErr)
		}
		if err := unmarshalBaselinePayload(baselineItem.BaselinePayload, &baselinePayload); err != nil {
			return nil, err
		}
		submittedAtValue, submitterName, err = s.extractInitialBaselineMetaRaw(ctx, baselineItem)
		if err != nil {
			return nil, err
		}
		count, countErr := s.initialBaselineRepo.CountSubmitted(ctx)
		if countErr != nil {
			return nil, fmt.Errorf("count submitted initial baseline: %w", countErr)
		}
		appliedGroupCount = int(count)
	} else {
		templateItem, templateErr := s.initialBaselineRepo.FindTemplateSample(ctx)
		switch {
		case templateErr == nil:
			if err := unmarshalBaselinePayload(templateItem.BaselinePayload, &baselinePayload); err != nil {
				return nil, err
			}
		case errors.Is(templateErr, gorm.ErrRecordNotFound):
		default:
			return nil, fmt.Errorf("load initial baseline template: %w", templateErr)
		}
	}

	return buildInitialBaselineViewResult(
		gameConfig.InitialBaselineSubmitted,
		baselinePayload,
		appliedGroupCount,
		submitterName,
		submittedAtValue,
	), nil
}

type openNextYearStatus struct {
	CanOpenNextYear  bool
	NextOpenableYear int
	BlockedReason    string
}

func evaluateOpenNextYearStatus(gameConfig *entity.GameConfig, groups []entity.Group, states []entity.GroupYearState) openNextYearStatus {
	nextOpenableYear := 0
	if gameConfig.CurrentOpenYear < gameConfig.FinalYear {
		nextOpenableYear = gameConfig.CurrentOpenYear + 1
	}

	result := openNextYearStatus{
		CanOpenNextYear:  false,
		NextOpenableYear: nextOpenableYear,
		BlockedReason:    "",
	}

	if gameConfig.CurrentOpenYear >= gameConfig.FinalYear {
		result.BlockedReason = openNextYearBlockedReasonFinalYearReached
		return result
	}
	if !gameConfig.InitialBaselineSubmitted {
		result.BlockedReason = openNextYearBlockedReasonBaselineUnsubmitted
		return result
	}

	stateMap := make(map[int64]entity.GroupYearState, len(states))
	for _, item := range states {
		stateMap[item.GroupID] = item
	}
	for _, group := range groups {
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			continue
		}
		stateItem, ok := stateMap[group.ID]
		if !ok || stateItem.YearStatus != enum.YearStatusCompleted {
			result.BlockedReason = openNextYearBlockedReasonUnfinishedReports
			return result
		}
	}
	for _, item := range states {
		if item.RollbackPending {
			result.BlockedReason = openNextYearBlockedReasonRollbackPending
			return result
		}
	}

	result.CanOpenNextYear = true
	return result
}

func (s *AdminControlQueryService) extractInitialBaselineMeta(ctx context.Context, item *entity.InitialBaseline) (*string, *string, error) {
	submittedAtValue, submitterName, err := s.extractInitialBaselineMetaRaw(ctx, item)
	if err != nil {
		return nil, nil, err
	}
	return formatOptionalRFC3339(submittedAtValue), submitterName, nil
}

func (s *AdminControlQueryService) extractInitialBaselineMetaRaw(ctx context.Context, item *entity.InitialBaseline) (*time.Time, *string, error) {
	if item == nil {
		return nil, nil, nil
	}

	var submitterName *string
	if item.SubmitterID != nil {
		account, accountErr := s.accountRepo.GetByID(ctx, *item.SubmitterID)
		if accountErr != nil {
			if !errors.Is(accountErr, gorm.ErrRecordNotFound) {
				return nil, nil, fmt.Errorf("load initial baseline submitter: %w", accountErr)
			}
		} else {
			value := account.Username
			submitterName = &value
		}
	}

	return item.SubmittedAt, submitterName, nil
}

func (s *AdminControlQueryService) lookupLatestAdminAction(ctx context.Context) (*AdminActionSummary, error) {
	item, err := s.adminActionLogRepo.FindLatest(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load latest admin action: %w", err)
	}
	return buildAdminActionSummary(item), nil
}

func buildAdminActionSummary(item *entity.AdminActionLog) *AdminActionSummary {
	if item == nil {
		return nil
	}
	return &AdminActionSummary{
		ActionCode:    item.ActionCode,
		OperatorName:  item.OperatorName,
		OperateTime:   item.OperateTime.Format(time.RFC3339),
		TargetGroupID: item.TargetGroupID,
		TargetYearNo:  item.TargetYearNo,
	}
}

func buildInitialBaselineViewResult(submitted bool, baselinePayload payload.BaselinePayload, appliedGroupCount int, submitterName *string, submittedAt *time.Time) *InitialBaselineViewResult {
	return &InitialBaselineViewResult{
		Submitted:         submitted,
		Editable:          !submitted,
		BaselinePayload:   baselinePayload,
		AppliedGroupCount: appliedGroupCount,
		SubmitterName:     submitterName,
		SubmittedAt:       formatOptionalRFC3339(submittedAt),
	}
}

func unmarshalBaselinePayload(raw []byte, target *payload.BaselinePayload) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("unmarshal initial baseline: %w", err)
	}
	return nil
}

func formatOptionalRFC3339(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(time.RFC3339)
	return &formatted
}
