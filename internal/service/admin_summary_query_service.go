package service

import (
	"context"
	"errors"
	"fmt"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

var (
	ErrAdminSummaryYearInvalid = errors.New("admin summary year invalid")
	ErrFinalRankingNotReady    = errors.New("final ranking not ready")
)

type YearSummaryItem struct {
	GroupID        int64   `json:"groupId"`
	GroupNo        int     `json:"groupNo"`
	GroupName      string  `json:"groupName"`
	Revenue        float64 `json:"revenue"`
	Profit         float64 `json:"profit"`
	Equity         float64 `json:"equity"`
	BusinessStatus string  `json:"businessStatus"`
	Ranking        *int    `json:"ranking,omitempty"`
}

type GetAdminYearSummaryResult struct {
	YearNo int               `json:"yearNo"`
	List   []YearSummaryItem `json:"list"`
}

type FinalRankingItem struct {
	Ranking        int     `json:"ranking"`
	GroupID        int64   `json:"groupId"`
	GroupNo        int     `json:"groupNo"`
	GroupName      string  `json:"groupName"`
	Equity         float64 `json:"equity"`
	Revenue        float64 `json:"revenue"`
	Profit         float64 `json:"profit"`
	BusinessStatus string  `json:"businessStatus"`
}

type GetFinalRankingResult struct {
	FinalYear int                `json:"finalYear"`
	List      []FinalRankingItem `json:"list"`
}

type AdminSummaryQueryService struct {
	gameConfigRepo *repository.GameConfigRepository
	groupRepo      *repository.GroupRepository
	groupYearRepo  *repository.GroupYearStateRepository
	summaryRepo    *repository.SummarySnapshotRepository
}

func NewAdminSummaryQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	summaryRepo *repository.SummarySnapshotRepository,
) *AdminSummaryQueryService {
	return &AdminSummaryQueryService{
		gameConfigRepo: gameConfigRepo,
		groupRepo:      groupRepo,
		groupYearRepo:  groupYearRepo,
		summaryRepo:    summaryRepo,
	}
}

func (s *AdminSummaryQueryService) GetYearSummary(ctx context.Context, yearNo int) (*GetAdminYearSummaryResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	if !isFormalSummaryYear(yearNo, gameConfig.FinalYear) {
		return nil, ErrAdminSummaryYearInvalid
	}

	summaryItems, err := s.summaryRepo.ListEffectiveByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("load summary snapshots: %w", err)
	}

	rankingMap := map[int64]int{}
	if yearNo == gameConfig.FinalYear {
		rankingItems, rankingErr := s.summaryRepo.ListEffectiveRankingByYear(ctx, yearNo)
		if rankingErr != nil {
			return nil, fmt.Errorf("load final year ranking snapshots: %w", rankingErr)
		}
		rankingMap = buildRankingMap(rankingItems)
	}

	return &GetAdminYearSummaryResult{
		YearNo: yearNo,
		List:   buildYearSummaryItems(summaryItems, rankingMap),
	}, nil
}

func (s *AdminSummaryQueryService) GetFinalRanking(ctx context.Context) (*GetFinalRankingResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("load groups: %w", err)
	}
	states, err := s.groupYearRepo.ListByYear(ctx, gameConfig.FinalYear)
	if err != nil {
		return nil, fmt.Errorf("load final year states: %w", err)
	}
	if !isFinalRankingReady(groups, states) {
		return nil, ErrFinalRankingNotReady
	}

	items, err := s.summaryRepo.ListEffectiveRankingByYear(ctx, gameConfig.FinalYear)
	if err != nil {
		return nil, fmt.Errorf("load final ranking snapshots: %w", err)
	}

	return &GetFinalRankingResult{
		FinalYear: gameConfig.FinalYear,
		List:      buildFinalRankingItems(items),
	}, nil
}

func buildYearSummaryItems(items []repository.SummarySnapshotWithGroup, rankingMap map[int64]int) []YearSummaryItem {
	result := make([]YearSummaryItem, 0, len(items))
	for _, item := range items {
		viewItem := YearSummaryItem{
			GroupID:        item.GroupID,
			GroupNo:        item.GroupNo,
			GroupName:      item.GroupName,
			Revenue:        item.Revenue,
			Profit:         item.Profit,
			Equity:         item.Equity,
			BusinessStatus: item.BusinessStatus,
		}
		if ranking, ok := rankingMap[item.GroupID]; ok {
			rankingValue := ranking
			viewItem.Ranking = &rankingValue
		}
		result = append(result, viewItem)
	}
	return result
}

func buildRankingMap(items []repository.SummarySnapshotWithGroup) map[int64]int {
	result := make(map[int64]int, len(items))
	for index, item := range items {
		result[item.GroupID] = index + 1
	}
	return result
}

func buildFinalRankingItems(items []repository.SummarySnapshotWithGroup) []FinalRankingItem {
	result := make([]FinalRankingItem, 0, len(items))
	for index, item := range items {
		result = append(result, FinalRankingItem{
			Ranking:        index + 1,
			GroupID:        item.GroupID,
			GroupNo:        item.GroupNo,
			GroupName:      item.GroupName,
			Equity:         item.Equity,
			Revenue:        item.Revenue,
			Profit:         item.Profit,
			BusinessStatus: item.BusinessStatus,
		})
	}
	return result
}

func isFinalRankingReady(groups []entity.Group, states []entity.GroupYearState) bool {
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
			return false
		}
	}
	return true
}

func isFormalSummaryYear(yearNo int, finalYear int) bool {
	return yearNo > 0 && yearNo <= finalYear
}
