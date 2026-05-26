package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

const (
	yearTabStatusEnterable        = "ENTERABLE"
	yearTabStatusLocked           = "LOCKED"
	yearTabStatusCompleted        = "COMPLETED"
	yearTabStatusBankruptReadOnly = "BANKRUPT_READONLY"
)

type CurrentGameConfigResult struct {
	CurrentOpenYear          int    `json:"currentOpenYear"`
	FinalYear                int    `json:"finalYear"`
	EditionCode              string `json:"editionCode"`
	EditionName              string `json:"editionName"`
	RuleVersion              string `json:"ruleVersion"`
	FormulaVersion           string `json:"formulaVersion"`
	TemplateVersion          string `json:"templateVersion"`
	OperatingTemplateVersion string `json:"operatingTemplateVersion"`
	ReportTemplateVersion    string `json:"reportTemplateVersion"`
	OrderTemplateVersion     string `json:"orderTemplateVersion"`
	ProcessRuleVersion       string `json:"processRuleVersion"`
	DemoYearEnabled          bool   `json:"demoYearEnabled"`
}

type YearTabItem struct {
	YearNo            int    `json:"yearNo"`
	Label             string `json:"label"`
	TabStatus         string `json:"tabStatus"`
	CanEnter          bool   `json:"canEnter"`
	IsCurrentOpenYear bool   `json:"isCurrentOpenYear"`
	IsFormalYear      bool   `json:"isFormalYear"`
}

type YearTabsResult struct {
	CurrentOpenYear int           `json:"currentOpenYear"`
	FinalYear       int           `json:"finalYear"`
	DemoYearEnabled bool          `json:"demoYearEnabled"`
	Tabs            []YearTabItem `json:"tabs"`
}

type GameConfigQueryService struct {
	gameConfigRepo *repository.GameConfigRepository
	groupRepo      *repository.GroupRepository
	groupYearRepo  *repository.GroupYearStateRepository
}

func NewGameConfigQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
) *GameConfigQueryService {
	return &GameConfigQueryService{
		gameConfigRepo: gameConfigRepo,
		groupRepo:      groupRepo,
		groupYearRepo:  groupYearRepo,
	}
}

func (s *GameConfigQueryService) GetCurrent(ctx context.Context) (*CurrentGameConfigResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	edition := normalizeGameConfigEdition(*gameConfig)

	return &CurrentGameConfigResult{
		CurrentOpenYear:          gameConfig.CurrentOpenYear,
		FinalYear:                gameConfig.FinalYear,
		EditionCode:              edition.EditionCode,
		EditionName:              edition.EditionName,
		RuleVersion:              edition.RuleVersion,
		FormulaVersion:           edition.FormulaVersion,
		TemplateVersion:          edition.TemplateVersion,
		OperatingTemplateVersion: edition.OperatingTemplateVersion,
		ReportTemplateVersion:    edition.ReportTemplateVersion,
		OrderTemplateVersion:     edition.OrderTemplateVersion,
		ProcessRuleVersion:       edition.ProcessRuleVersion,
		DemoYearEnabled:          true,
	}, nil
}

func (s *GameConfigQueryService) ListGameEditions(ctx context.Context) ([]GameEdition, error) {
	return ListGameEditions(), nil
}

func (s *GameConfigQueryService) GetYearTabs(ctx context.Context, roleType string, groupID *int64) (*YearTabsResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	result := &YearTabsResult{
		CurrentOpenYear: gameConfig.CurrentOpenYear,
		FinalYear:       gameConfig.FinalYear,
		DemoYearEnabled: true,
	}

	if roleType == enum.RoleTypeGroup && groupID != nil {
		group, err := s.groupRepo.GetByID(ctx, *groupID)
		if err != nil {
			return nil, fmt.Errorf("load group: %w", err)
		}
		stateMap, err := s.loadGroupYearStateMap(ctx, *groupID, gameConfig.FinalYear)
		if err != nil {
			return nil, err
		}
		result.Tabs = buildGroupYearTabs(gameConfig.CurrentOpenYear, gameConfig.FinalYear, group.BusinessStatus, group.BankruptYearNo, stateMap)
		return result, nil
	}

	result.Tabs = buildAdminYearTabs(gameConfig.CurrentOpenYear, gameConfig.FinalYear)
	return result, nil
}

func (s *GameConfigQueryService) loadGroupYearStateMap(ctx context.Context, groupID int64, finalYear int) (map[int]entity.GroupYearState, error) {
	result := make(map[int]entity.GroupYearState, finalYear+1)
	for yearNo := 0; yearNo <= finalYear; yearNo++ {
		item, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo)
		switch {
		case err == nil:
			result[yearNo] = *item
		case errors.Is(err, gorm.ErrRecordNotFound):
			continue
		default:
			return nil, fmt.Errorf("load group year state: %w", err)
		}
	}
	return result, nil
}

func buildAdminYearTabs(currentOpenYear int, finalYear int) []YearTabItem {
	tabs := make([]YearTabItem, 0, finalYear+1)
	for yearNo := 0; yearNo <= finalYear; yearNo++ {
		tabStatus := yearTabStatusLocked
		canEnter := false
		if yearNo <= currentOpenYear {
			tabStatus = yearTabStatusEnterable
			canEnter = true
		}
		tabs = append(tabs, YearTabItem{
			YearNo:            yearNo,
			Label:             buildYearLabel(yearNo),
			TabStatus:         tabStatus,
			CanEnter:          canEnter,
			IsCurrentOpenYear: yearNo == currentOpenYear,
			IsFormalYear:      yearNo > 0,
		})
	}
	return tabs
}

func buildGroupYearTabs(currentOpenYear int, finalYear int, businessStatus string, bankruptYearNo *int, stateMap map[int]entity.GroupYearState) []YearTabItem {
	tabs := make([]YearTabItem, 0, finalYear+1)
	for yearNo := 0; yearNo <= finalYear; yearNo++ {
		tabStatus := resolveGroupYearTabStatus(yearNo, currentOpenYear, businessStatus, bankruptYearNo, stateMap)
		tabs = append(tabs, YearTabItem{
			YearNo:            yearNo,
			Label:             buildYearLabel(yearNo),
			TabStatus:         tabStatus,
			CanEnter:          tabStatus != yearTabStatusLocked,
			IsCurrentOpenYear: yearNo == currentOpenYear,
			IsFormalYear:      yearNo > 0,
		})
	}
	return tabs
}

func resolveGroupYearTabStatus(yearNo int, currentOpenYear int, businessStatus string, bankruptYearNo *int, stateMap map[int]entity.GroupYearState) string {
	if bankruptYearNo != nil {
		if yearNo == *bankruptYearNo {
			return yearTabStatusBankruptReadOnly
		}
		if yearNo > *bankruptYearNo {
			return yearTabStatusLocked
		}
	}

	if item, ok := stateMap[yearNo]; ok && item.YearStatus == enum.YearStatusCompleted {
		return yearTabStatusCompleted
	}

	if businessStatus == enum.BusinessStatusBankrupt && bankruptYearNo == nil && yearNo == currentOpenYear {
		return yearTabStatusBankruptReadOnly
	}
	if yearNo <= currentOpenYear {
		return yearTabStatusEnterable
	}
	return yearTabStatusLocked
}

func buildYearLabel(yearNo int) string {
	return fmt.Sprintf("%d年", yearNo)
}
