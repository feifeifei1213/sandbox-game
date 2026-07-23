package service

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func TestGetYearSummaryAndFinalRankingUseFormalSnapshotsConsistently(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 2, 2, true)
	markAllExistingGroupsBankrupt(t, ctx, tx, 2)
	withdrawExistingSummarySnapshotsForIntegrationTest(t, ctx, tx, 2)

	bankruptYearNo := 2
	groupOneID := createAdminIntegrationGroup(t, ctx, tx, "汇总测试-第1组", enum.BusinessStatusNormal, nil)
	_ = createAdminIntegrationGroup(t, ctx, tx, "汇总测试-第2组", enum.BusinessStatusBankrupt, &bankruptYearNo)
	groupThreeID := createAdminIntegrationGroup(t, ctx, tx, "汇总测试-第3组", enum.BusinessStatusNormal, nil)
	groupFourID := createAdminIntegrationGroup(t, ctx, tx, "汇总测试-第4组", enum.BusinessStatusBankrupt, &bankruptYearNo)

	createDetailedGroupYearState(t, ctx, tx, groupOneID, 2, enum.YearTypeFormal, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, true, 5, 1)
	createDetailedGroupYearState(t, ctx, tx, groupThreeID, 2, enum.YearTypeFormal, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, true, 5, 1)
	createDetailedGroupYearState(t, ctx, tx, groupFourID, 2, enum.YearTypeFormal, enum.YearStatusLocked, enum.StageStatusQ1Open, enum.ReportStatusLocked, true, 0, 0)

	createSummarySnapshotFixture(t, ctx, tx, groupOneID, 2, 80, 10, 68, enum.BusinessStatusNormal, 68, true, 1)
	createSummarySnapshotFixture(t, ctx, tx, groupThreeID, 2, 92, 14, 82, enum.BusinessStatusNormal, 82, true, 1)
	// 汇总快照保留历史生成时的状态；管理员汇总页行级经营状态必须显示小组当前状态。
	createSummarySnapshotFixture(t, ctx, tx, groupFourID, 2, 60, 5, 50, enum.BusinessStatusNormal, 50, true, 1)

	service := buildAdminSummaryQueryService(tx)

	yearSummary, err := service.GetYearSummary(ctx, 2)
	if err != nil {
		t.Fatalf("get year summary: %v", err)
	}
	if yearSummary.YearNo != 2 {
		t.Fatalf("expected year summary for year 2, got %d", yearSummary.YearNo)
	}
	if len(yearSummary.List) != 3 {
		t.Fatalf("expected 3 effective summary items, got %d", len(yearSummary.List))
	}
	viewByGroupID := make(map[int64]YearSummaryItem, len(yearSummary.List))
	for _, item := range yearSummary.List {
		viewByGroupID[item.GroupID] = item
	}
	if viewByGroupID[groupOneID].Ranking == nil || *viewByGroupID[groupOneID].Ranking != 2 {
		t.Fatalf("expected group one ranking 2 in year summary, got %#v", viewByGroupID[groupOneID])
	}
	if viewByGroupID[groupThreeID].Ranking == nil || *viewByGroupID[groupThreeID].Ranking != 1 {
		t.Fatalf("expected group three ranking 1 in year summary, got %#v", viewByGroupID[groupThreeID])
	}
	if viewByGroupID[groupFourID].BusinessStatus != enum.BusinessStatusBankrupt {
		t.Fatalf("expected year summary to expose current BANKRUPT status instead of historical snapshot status, got %#v", viewByGroupID[groupFourID])
	}

	finalRanking, err := service.GetFinalRanking(ctx)
	if err != nil {
		t.Fatalf("get final ranking: %v", err)
	}
	if finalRanking.FinalYear != 2 {
		t.Fatalf("expected final year 2, got %d", finalRanking.FinalYear)
	}
	if len(finalRanking.List) != 3 {
		t.Fatalf("expected 3 final ranking items, got %d", len(finalRanking.List))
	}
	if finalRanking.List[0].GroupID != groupThreeID || finalRanking.List[0].Ranking != 1 {
		t.Fatalf("expected group three ranking 1, got %#v", finalRanking.List[0])
	}
	if finalRanking.List[1].GroupID != groupOneID || finalRanking.List[1].Ranking != 2 {
		t.Fatalf("expected group one ranking 2, got %#v", finalRanking.List[1])
	}
	if finalRanking.List[2].GroupID != groupFourID || finalRanking.List[2].BusinessStatus != enum.BusinessStatusBankrupt {
		t.Fatalf("expected final ranking to expose current BANKRUPT status for group four, got %#v", finalRanking.List[2])
	}
}

func buildAdminSummaryQueryService(db *gorm.DB) *AdminSummaryQueryService {
	return NewAdminSummaryQueryService(
		repository.NewGameConfigRepository(db),
		repository.NewGroupRepository(db),
		repository.NewGroupYearStateRepository(db),
		repository.NewSummarySnapshotRepository(db),
	)
}

func withdrawExistingSummarySnapshotsForIntegrationTest(t *testing.T, ctx context.Context, tx *gorm.DB, yearNo int) {
	t.Helper()

	if err := tx.WithContext(ctx).
		Model(&entity.GroupSummarySnapshot{}).
		Where("year_no = ?", yearNo).
		Update("summary_effective", false).Error; err != nil {
		t.Fatalf("withdraw existing summary snapshots: %v", err)
	}
}
