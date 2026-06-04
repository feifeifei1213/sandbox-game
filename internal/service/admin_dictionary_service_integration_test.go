package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func TestAdminDictionaryLifecycleInitializesSnapshotAndRevision(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 3, 0, false)
	clearInitializationFixtures(t, ctx, tx)
	clearDictionaryFixturesForIntegrationTest(t, ctx, tx)

	dictionaryService := NewAdminDictionaryService(tx)
	scheme, err := dictionaryService.SaveScheme(ctx, SaveDictionarySchemeCommand{
		EditionCode:  GameEditionVIPServiceV1,
		SchemeName:   "集成测试显示名称方案",
		Description:  "用于验证初始化字典快照",
		Items:        []DictionaryItemInput{{ItemCode: "market.LOCAL", DisplayName: "本地市场测试名"}},
		OperatorID:   91001,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("save dictionary scheme: %v", err)
	}
	if scheme.ID <= 0 {
		t.Fatalf("expected scheme id to be generated, got %#v", scheme)
	}
	detail, err := dictionaryService.GetSchemeDetail(ctx, scheme.ID)
	if err != nil {
		t.Fatalf("get dictionary scheme detail: %v", err)
	}
	if displayNameOf(detail.Items, "market.LOCAL") != "本地市场测试名" {
		t.Fatalf("expected saved scheme display name, got %q", displayNameOf(detail.Items, "market.LOCAL"))
	}
	scheme, err = dictionaryService.SaveScheme(ctx, SaveDictionarySchemeCommand{
		SchemeID:     &scheme.ID,
		EditionCode:  GameEditionVIPServiceV1,
		SchemeName:   "集成测试显示名称方案",
		Description:  "用于验证初始化字典快照",
		Items:        []DictionaryItemInput{{ItemCode: "market.LOCAL", DisplayName: "本地市场方案更新"}},
		OperatorID:   91001,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("update dictionary scheme: %v", err)
	}
	if displayNameOf(scheme.Items, "market.LOCAL") != "本地市场方案更新" {
		t.Fatalf("expected updated scheme display name, got %q", displayNameOf(scheme.Items, "market.LOCAL"))
	}

	controlService := NewAdminControlCommandService(tx)
	initResult, err := controlService.InitializeGame(ctx, InitializeGameCommand{
		GroupCount:         2,
		EditionCode:        GameEditionVIPServiceV1,
		DictionarySchemeID: &scheme.ID,
		OperatorID:         91002,
		OperatorName:       "integration-admin",
	})
	if err != nil {
		t.Fatalf("initialize game with dictionary scheme: %v", err)
	}
	if initResult.DictionaryRevision != 1 {
		t.Fatalf("expected dictionary revision 1 after initialization, got %d", initResult.DictionaryRevision)
	}

	current, err := dictionaryService.GetCurrent(ctx, "")
	if err != nil {
		t.Fatalf("get current dictionary after initialization: %v", err)
	}
	if !current.Initialized || current.DictionaryRevision != 1 || current.EditionCode != GameEditionVIPServiceV1 {
		t.Fatalf("unexpected current dictionary state: %#v", current)
	}
	if displayNameOf(current.Items, "market.LOCAL") != "本地市场方案更新" {
		t.Fatalf("expected scheme display name to be copied into current dictionary, got %q", displayNameOf(current.Items, "market.LOCAL"))
	}

	updated, err := dictionaryService.UpdateCurrent(ctx, UpdateCurrentDictionaryCommand{
		Items:        []DictionaryItemInput{{ItemCode: "market.LOCAL", DisplayName: "本地市场修改后"}},
		Reason:       "验证当前比赛字典修改",
		OperatorID:   91003,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("update current dictionary: %v", err)
	}
	if updated.DictionaryRevision != 2 {
		t.Fatalf("expected dictionary revision 2 after update, got %d", updated.DictionaryRevision)
	}
	if displayNameOf(updated.Items, "market.LOCAL") != "本地市场修改后" {
		t.Fatalf("expected updated display name, got %q", displayNameOf(updated.Items, "market.LOCAL"))
	}

	applied, err := dictionaryService.ApplySchemeToCurrent(ctx, ApplyDictionarySchemeCommand{
		SchemeID:     scheme.ID,
		OperatorID:   91007,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("apply scheme to current: %v", err)
	}
	if applied.DictionaryRevision != 3 {
		t.Fatalf("expected dictionary revision 3 after apply scheme, got %d", applied.DictionaryRevision)
	}
	if displayNameOf(applied.Items, "market.LOCAL") != "本地市场方案更新" {
		t.Fatalf("expected applied scheme display name, got %q", displayNameOf(applied.Items, "market.LOCAL"))
	}

	restored, err := dictionaryService.RestoreCurrentDefault(ctx, RestoreCurrentDictionaryDefaultCommand{
		OperatorID:   91008,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("restore current default: %v", err)
	}
	if restored.DictionaryRevision != 4 {
		t.Fatalf("expected dictionary revision 4 after restore default, got %d", restored.DictionaryRevision)
	}
	if displayNameOf(restored.Items, "market.LOCAL") != "本地市场" {
		t.Fatalf("expected default display name after restore, got %q", displayNameOf(restored.Items, "market.LOCAL"))
	}

	revision, err := dictionaryService.GetRevision(ctx)
	if err != nil {
		t.Fatalf("get dictionary revision: %v", err)
	}
	if revision.DictionaryRevision != 4 || revision.EditionCode != GameEditionVIPServiceV1 {
		t.Fatalf("unexpected revision result: %#v", revision)
	}

	if err := dictionaryService.DeleteScheme(ctx, DeleteDictionarySchemeCommand{
		SchemeID:     scheme.ID,
		OperatorID:   91009,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("delete dictionary scheme: %v", err)
	}
	schemeList, err := dictionaryService.ListSchemes(ctx, GameEditionVIPServiceV1)
	if err != nil {
		t.Fatalf("list dictionary schemes after delete: %v", err)
	}
	if hasScheme(schemeList.List, scheme.ID) {
		t.Fatalf("expected scheme to be deleted, got %#v", schemeList.List)
	}

	logs, err := dictionaryService.PageChangeLogs(ctx, 1, 20)
	if err != nil {
		t.Fatalf("page dictionary change logs: %v", err)
	}
	if logs.Total < 7 {
		t.Fatalf("expected save, update scheme, initialize, update current, apply, restore and delete logs, got %#v", logs)
	}
}

func TestAdminDictionaryRejectsCrossEditionSchemeForCurrentGame(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 3, 0, false)
	clearInitializationFixtures(t, ctx, tx)
	clearDictionaryFixturesForIntegrationTest(t, ctx, tx)

	dictionaryService := NewAdminDictionaryService(tx)
	productionScheme, err := dictionaryService.SaveScheme(ctx, SaveDictionarySchemeCommand{
		EditionCode:  GameEditionProductionV1,
		SchemeName:   "生产版显示名称方案",
		Items:        []DictionaryItemInput{{ItemCode: "operating.materialGroupTitle", DisplayName: "生产版材料标题"}},
		OperatorID:   91004,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("save production dictionary scheme: %v", err)
	}

	controlService := NewAdminControlCommandService(tx)
	if _, err := controlService.InitializeGame(ctx, InitializeGameCommand{
		GroupCount:   1,
		EditionCode:  GameEditionVIPServiceV1,
		OperatorID:   91005,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("initialize vip game: %v", err)
	}

	_, err = dictionaryService.ApplySchemeToCurrent(ctx, ApplyDictionarySchemeCommand{
		SchemeID:     productionScheme.ID,
		OperatorID:   91006,
		OperatorName: "integration-admin",
	})
	if !errors.Is(err, ErrDictionarySchemeCrossEdition) {
		t.Fatalf("expected cross edition scheme error, got %v", err)
	}
}

func clearDictionaryFixturesForIntegrationTest(t *testing.T, ctx context.Context, tx *gorm.DB) {
	t.Helper()

	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.DictionaryChangeLog{}).Error; err != nil {
		t.Fatalf("clear dictionary change logs: %v", err)
	}
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.CurrentDictionaryItem{}).Error; err != nil {
		t.Fatalf("clear current dictionary items: %v", err)
	}
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.DictionarySchemeItem{}).Error; err != nil {
		t.Fatalf("clear dictionary scheme items: %v", err)
	}
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.DictionaryScheme{}).Error; err != nil {
		t.Fatalf("clear dictionary schemes: %v", err)
	}
	gameConfig, err := repository.NewGameConfigRepository(tx).GetCurrent(ctx)
	if err != nil {
		t.Fatalf("load game config for dictionary revision reset: %v", err)
	}
	if err := repository.NewGameConfigRepository(tx).UpdateDictionaryRevision(ctx, gameConfig.ID, 0, "integration-test", time.Now()); err != nil {
		t.Fatalf("reset dictionary revision: %v", err)
	}
}

func displayNameOf(items []DictionaryItemResult, itemCode string) string {
	for _, item := range items {
		if item.ItemCode == itemCode {
			return item.DisplayName
		}
	}
	return ""
}

func hasScheme(items []DictionarySchemeSummary, schemeID int64) bool {
	for _, item := range items {
		if item.ID == schemeID {
			return true
		}
	}
	return false
}
