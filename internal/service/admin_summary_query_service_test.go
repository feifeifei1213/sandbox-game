package service

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func TestBuildRankingMapAssignsDescendingOrderRanking(t *testing.T) {
	items := []repository.SummarySnapshotWithGroup{
		{GroupID: 3, RankingValue: 680},
		{GroupID: 1, RankingValue: 520},
		{GroupID: 2, RankingValue: 410},
	}

	rankingMap := buildRankingMap(items)
	if rankingMap[3] != 1 {
		t.Fatalf("expected group 3 ranking to be 1, got %d", rankingMap[3])
	}
	if rankingMap[1] != 2 {
		t.Fatalf("expected group 1 ranking to be 2, got %d", rankingMap[1])
	}
	if rankingMap[2] != 3 {
		t.Fatalf("expected group 2 ranking to be 3, got %d", rankingMap[2])
	}
}

func TestBuildYearSummaryItemsUsesRankingMapInsteadOfListOrder(t *testing.T) {
	items := []repository.SummarySnapshotWithGroup{
		{GroupID: 1, GroupNo: 1, GroupName: "第1组", Revenue: 100, Profit: 10, Equity: 200, BusinessStatus: enum.BusinessStatusNormal},
		{GroupID: 2, GroupNo: 2, GroupName: "第2组", Revenue: 120, Profit: 20, Equity: 260, BusinessStatus: enum.BusinessStatusNormal},
	}

	result := buildYearSummaryItems(items, map[int64]int{1: 2, 2: 1})
	if result[0].Ranking == nil || *result[0].Ranking != 2 {
		t.Fatalf("expected group 1 ranking to be 2")
	}
	if result[1].Ranking == nil || *result[1].Ranking != 1 {
		t.Fatalf("expected group 2 ranking to be 1")
	}
}

func TestIsFinalRankingReadyRequiresAllNonBankruptGroupsCompleted(t *testing.T) {
	groups := []entity.Group{
		{ID: 1, BusinessStatus: enum.BusinessStatusNormal},
		{ID: 2, BusinessStatus: enum.BusinessStatusNormal},
	}
	states := []entity.GroupYearState{
		{GroupID: 1, YearStatus: enum.YearStatusCompleted},
		{GroupID: 2, YearStatus: enum.YearStatusOperating},
	}

	if isFinalRankingReady(groups, states) {
		t.Fatalf("expected final ranking to be not ready")
	}
}

func TestIsFinalRankingReadyAllowsBankruptGroupsToStopBlocking(t *testing.T) {
	groups := []entity.Group{
		{ID: 1, BusinessStatus: enum.BusinessStatusNormal},
		{ID: 2, BusinessStatus: enum.BusinessStatusBankrupt},
	}
	states := []entity.GroupYearState{
		{GroupID: 1, YearStatus: enum.YearStatusCompleted},
		{GroupID: 2, YearStatus: enum.YearStatusLocked},
	}

	if !isFinalRankingReady(groups, states) {
		t.Fatalf("expected final ranking to be ready")
	}
}
