package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/enum"
	appdb "sandbox-game/internal/infra/db"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestSubmitOperatingStageAdvancesEditableScopeFromQ1ToQ2(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	groupID := createIntegrationOperatingFixtures(t, ctx, tx)
	commandService, queryService := buildPlayerOperatingServices(tx)

	result, err := commandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupID,
		YearNo:           0,
		StageCode:        state.StageCodeQ1,
		OperatingPayload: buildValidQ1OperatingPayload(),
		SubmitterID:      90001,
		OperatorName:     "integration-test",
	})
	if err != nil {
		t.Fatalf("submit stage: %v", err)
	}

	if result.StageCode != state.StageCodeQ1 {
		t.Fatalf("expected stage code Q1, got %s", result.StageCode)
	}
	if result.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected year status OPERATING, got %s", result.YearStatus)
	}
	if result.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected stage status Q2_OPEN, got %s", result.StageStatus)
	}
	if result.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected report status REPORT_LOCKED, got %s", result.ReportStatus)
	}
	if result.BusinessStatus != enum.BusinessStatusNormal {
		t.Fatalf("expected business status NORMAL, got %s", result.BusinessStatus)
	}
	if result.LatestStageSubmitVersion != 1 {
		t.Fatalf("expected latest stage submit version 1, got %d", result.LatestStageSubmitVersion)
	}
	if result.PeriodEndCash <= 0 {
		t.Fatalf("expected positive period end cash after Q1 submit, got %f", result.PeriodEndCash)
	}
	if result.SubmittedAt.IsZero() {
		t.Fatalf("expected submittedAt to be filled")
	}

	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("reload year view: %v", err)
	}

	if view.CurrentStageCode != state.StageCodeQ2 {
		t.Fatalf("expected current stage code Q2, got %s", view.CurrentStageCode)
	}
	if view.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected reloaded stage status Q2_OPEN, got %s", view.StageStatus)
	}
	if !view.CanEdit || !view.CanSubmit {
		t.Fatalf("expected reloaded view to remain editable/submittable")
	}
	if !slices.Equal(view.EditableScopes, []string{state.OperatingScopeQ2}) {
		t.Fatalf("expected editable scopes [Q2], got %#v", view.EditableScopes)
	}
	if !slices.Contains(view.ReadonlyScopes, state.OperatingScopeYearStart) || !slices.Contains(view.ReadonlyScopes, state.OperatingScopeQ1) {
		t.Fatalf("expected YEAR_START and Q1 to become readonly, got %#v", view.ReadonlyScopes)
	}
	if len(view.StageSubmitHistory) != 1 {
		t.Fatalf("expected one stage submission history item, got %d", len(view.StageSubmitHistory))
	}
	if view.StageSubmitHistory[0].StageCode != state.StageCodeQ1 || view.StageSubmitHistory[0].SubmitVersion != 1 {
		t.Fatalf("expected first history item to be Q1 version 1, got %#v", view.StageSubmitHistory[0])
	}
	if view.LastDraftSavedAt == nil || view.LastDraftSavedAt.IsZero() {
		t.Fatalf("expected last draft save time to be written on submit")
	}
}

func TestSubmitOperatingStageAfterRollbackUsesNextHistoricalVersionAndRetrySnapshot(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	groupID := createIntegrationOperatingFixtures(t, ctx, tx)
	now := time.Now()
	historicalPayload := buildValidQ1OperatingPayload()
	historicalPayload.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"] = 3
	if err := repository.NewOperatingRepository(tx).CreateStageSubmission(ctx, repository.CreateStageSubmissionCommand{
		GroupID:                  groupID,
		YearNo:                   0,
		StageCode:                state.StageCodeQ1,
		SubmitVersion:            1,
		PeriodEndCash:            30,
		OperatingPayloadSnapshot: historicalPayload,
		StateBeforeJSON:          []byte("{}"),
		StateAfterJSON:           []byte("{}"),
		SubmitterID:              90000,
		SubmitTime:               now.Add(-time.Hour),
	}); err != nil {
		t.Fatalf("create historical stage submission: %v", err)
	}
	if err := repository.NewGroupYearStateRepository(tx).MarkRollbackPending(ctx, groupID, 0, 0, state.StageCodeQ1, 10001, "integration-admin"); err != nil {
		t.Fatalf("mark rollback pending: %v", err)
	}

	commandService, _ := buildPlayerOperatingServices(tx)
	retryPayload := buildValidQ1OperatingPayload()
	retryPayload.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"] = 5
	result, err := commandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupID,
		YearNo:           0,
		StageCode:        state.StageCodeQ1,
		OperatingPayload: retryPayload,
		SubmitterID:      90001,
		OperatorName:     "integration-test",
	})
	if err != nil {
		t.Fatalf("submit rollback retry stage: %v", err)
	}
	if result.LatestStageSubmitVersion != 2 {
		t.Fatalf("expected retry stage submit version 2, got %d", result.LatestStageSubmitVersion)
	}

	var latestSubmission entity.GroupStageSubmission
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND stage_code = ?", groupID, 0, state.StageCodeQ1).
		Order("submit_version DESC").
		First(&latestSubmission).Error; err != nil {
		t.Fatalf("load latest stage submission: %v", err)
	}
	if latestSubmission.SubmitVersion != 2 {
		t.Fatalf("expected latest persisted stage version 2, got %d", latestSubmission.SubmitVersion)
	}
	var latestPayload payload.OperatingPayload
	if err := json.Unmarshal(latestSubmission.OperatingPayloadSnapshot, &latestPayload); err != nil {
		t.Fatalf("unmarshal latest operating payload snapshot: %v", err)
	}
	if got := latestPayload.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"]; got != float64(5) {
		t.Fatalf("expected rollback retry snapshot to keep edited supply chain order quantity 5, got %#v", got)
	}
	var historicalSubmission entity.GroupStageSubmission
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND stage_code = ? AND submit_version = ?", groupID, 0, state.StageCodeQ1, 1).
		First(&historicalSubmission).Error; err != nil {
		t.Fatalf("load historical stage submission: %v", err)
	}
	var historicalSnapshot payload.OperatingPayload
	if err := json.Unmarshal(historicalSubmission.OperatingPayloadSnapshot, &historicalSnapshot); err != nil {
		t.Fatalf("unmarshal historical operating payload snapshot: %v", err)
	}
	if got := historicalSnapshot.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"]; got != float64(3) {
		t.Fatalf("expected historical supply chain order quantity 3 to be preserved, got %#v", got)
	}

	snapshot := loadLatestSnapshotForTest(t, ctx, tx, groupID, 0, state.StageCodeQ1)
	if snapshot.TriggerCode != enum.SnapshotTriggerRollbackStageRetry {
		t.Fatalf("expected retry stage snapshot trigger, got %s", snapshot.TriggerCode)
	}
	if snapshot.Description == nil || *snapshot.Description != "回退后重新提交经营阶段自动快照" {
		t.Fatalf("unexpected retry stage snapshot description: %#v", snapshot.Description)
	}
}

func openIntegrationMySQL(t *testing.T) *gorm.DB {
	t.Helper()

	cfg, err := appconfig.Load(resolveProjectPath(t, "configs", "local.yaml"))
	if err != nil {
		t.Skipf("skip integration test: load local config failed: %v", err)
	}

	db, err := appdb.NewMySQL(cfg)
	if err != nil {
		t.Skipf("skip integration test: open mysql failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("extract sql db: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("skip integration test: ping mysql failed: %v", err)
	}

	ensureIntegrationGameConfigEditionColumns(t, db)
	ensureIntegrationNoticeTables(t, db)
	ensureIntegrationOrderTables(t, db)
	ensureIntegrationRollbackTables(t, db)
	ensureIntegrationDictionaryTables(t, db)
	return db
}

func ensureIntegrationGameConfigEditionColumns(t *testing.T, db *gorm.DB) {
	t.Helper()

	columns := []string{
		"EditionCode",
		"EditionName",
		"OperatingTemplateVersion",
		"ReportTemplateVersion",
		"OrderTemplateVersion",
		"ProcessRuleVersion",
		"DictionaryRevision",
	}
	for _, column := range columns {
		if db.Migrator().HasColumn(&entity.GameConfig{}, column) {
			continue
		}
		if err := db.Migrator().AddColumn(&entity.GameConfig{}, column); err != nil {
			t.Fatalf("ensure game config edition column %s: %v", column, err)
		}
	}
}

func ensureIntegrationDictionaryTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.AutoMigrate(
		&entity.DictionaryScheme{},
		&entity.DictionarySchemeItem{},
		&entity.CurrentDictionaryItem{},
		&entity.DictionaryChangeLog{},
	); err != nil {
		t.Fatalf("ensure dictionary tables: %v", err)
	}
}

func ensureIntegrationNoticeTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS sg_notice (
			id BIGINT NOT NULL AUTO_INCREMENT,
			target_scope VARCHAR(16) NOT NULL,
			target_group_id BIGINT NULL,
			content VARCHAR(1000) NOT NULL,
			pinned TINYINT(1) NOT NULL DEFAULT 0,
			published_at DATETIME NOT NULL,
			operator_id BIGINT NOT NULL,
			operator_name VARCHAR(64) NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_target_scope_group (target_scope, target_group_id),
			KEY idx_published_at (published_at),
			KEY idx_pinned (pinned)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_group_adjustment (
			id BIGINT NOT NULL AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			year_no INT NOT NULL,
			stage_code VARCHAR(16) NOT NULL,
			adjustment_type VARCHAR(16) NOT NULL,
			amount DECIMAL(18,2) NOT NULL DEFAULT 0,
			reason VARCHAR(500) NOT NULL,
			published_at DATETIME NOT NULL,
			operator_id BIGINT NOT NULL,
			operator_name VARCHAR(64) NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_group_year_stage (group_id, year_no, stage_code),
			KEY idx_published_at (published_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("ensure I2 tables: %v", err)
		}
	}
	for _, column := range []string{"VoidedByID", "VoidedByName", "VoidReason", "VoidedAt"} {
		if db.Migrator().HasColumn(&entity.GroupAdjustment{}, column) {
			continue
		}
		if err := db.Migrator().AddColumn(&entity.GroupAdjustment{}, column); err != nil {
			t.Fatalf("ensure adjustment lifecycle column %s: %v", column, err)
		}
	}
	if err := db.AutoMigrate(&entity.GroupAdjustmentRevision{}); err != nil {
		t.Fatalf("ensure adjustment revision table: %v", err)
	}
}

func ensureIntegrationOrderTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS sg_order_import_batch (
			id BIGINT NOT NULL AUTO_INCREMENT,
			original_file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL DEFAULT 0,
			parse_status VARCHAR(32) NOT NULL,
			parsed_order_count INT NOT NULL DEFAULT 0,
			parse_error_json JSON NOT NULL,
			parsed_payload_json JSON NOT NULL,
			uploader_id BIGINT NOT NULL,
			uploader_name VARCHAR(64) NOT NULL,
			uploaded_at DATETIME NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_order_import_uploaded_at (uploaded_at),
			KEY idx_order_import_status (parse_status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_generation_config (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			order_count INT NOT NULL DEFAULT 0,
			release_sequence_no INT NOT NULL,
			generation_batch_id BIGINT NULL,
			source_batch_id BIGINT NULL,
			config_status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_order_generation_config (year_no, market_code, order_type),
			KEY idx_order_release_sequence (year_no, release_sequence_no)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_forecast_control (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			forecast_stage_code VARCHAR(32) NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			order_count INT NOT NULL DEFAULT 0,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_order_forecast_control (year_no, market_code, order_type),
			KEY idx_order_forecast_stage (forecast_stage_code, market_code)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_market_forecast (
			id BIGINT NOT NULL AUTO_INCREMENT,
			forecast_stage_code VARCHAR(32) NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			forecast_data_json JSON NOT NULL,
			narrative TEXT NULL,
			formula_version VARCHAR(64) NOT NULL,
			random_seed VARCHAR(64) NOT NULL DEFAULT '',
			control_snapshot_json JSON NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_order_market_forecast (forecast_stage_code, market_code),
			KEY idx_order_market_forecast_stage (forecast_stage_code)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_generation_batch (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			batch_status VARCHAR(32) NOT NULL,
			formula_version VARCHAR(32) NOT NULL,
			random_seed VARCHAR(64) NOT NULL,
			control_snapshot_json JSON NOT NULL,
			forecast_snapshot_json JSON NULL,
			parameter_snapshot_json JSON NOT NULL,
			order_detail_json JSON NOT NULL,
			generated_order_count INT NOT NULL DEFAULT 0,
			generated_by_id BIGINT NOT NULL,
			generated_by_name VARCHAR(64) NOT NULL,
			generated_at DATETIME NOT NULL,
			confirmed_by_id BIGINT NULL,
			confirmed_by_name VARCHAR(64) NULL,
			confirmed_at DATETIME NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_order_generation_batch_year (year_no, batch_status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_market_config (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			market_enabled TINYINT(1) NOT NULL DEFAULT 0,
			market_investment_limit DECIMAL(18,2) NULL,
			config_status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
			locked_batch_id BIGINT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_order_market_config (year_no, market_code),
			KEY idx_order_market_enabled (year_no, market_enabled)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_order_pool (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			segment_code VARCHAR(64) NULL,
			card_sequence_no INT NULL,
			business_order_no VARCHAR(64) NULL,
			order_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
			order_quantity DECIMAL(18,2) NOT NULL DEFAULT 0,
			unit_price DECIMAL(18,2) NOT NULL DEFAULT 0,
			account_term INT NOT NULL DEFAULT 0,
			pool_status VARCHAR(32) NOT NULL DEFAULT 'AVAILABLE',
			selected_group_id BIGINT NULL,
			selected_at DATETIME NULL,
			generation_batch_id BIGINT NULL,
			source_batch_id BIGINT NULL,
			source_sheet_name VARCHAR(64) NULL,
			source_cell VARCHAR(16) NULL,
			source_row_index INT NULL,
			source_row_key VARCHAR(128) NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_order_pool_year_market (year_no, market_code, order_type),
			KEY idx_order_pool_status (pool_status),
			KEY idx_order_pool_selected_group (selected_group_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_market_bidding_state (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			segment_code VARCHAR(64) NOT NULL,
			release_sequence_no INT NOT NULL,
			segment_status VARCHAR(32) NOT NULL,
			leader_group_id BIGINT NULL,
			leader_rule_json JSON NULL,
			random_seed VARCHAR(64) NULL,
			current_group_id BIGINT NULL,
			opened_at DATETIME NULL,
			closed_at DATETIME NULL,
			released_at DATETIME NULL,
			completed_at DATETIME NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_market_bidding_state (year_no, market_code, order_type),
			KEY idx_segment_release_sequence (year_no, release_sequence_no),
			KEY idx_segment_status (segment_status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_group_market_bid (
			id BIGINT NOT NULL AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL DEFAULT 'AGENCY_INSPECTION',
			market_investment DECIMAL(18,2) NOT NULL DEFAULT 0,
			bid_status VARCHAR(32) NOT NULL DEFAULT 'SUBMITTED',
			submitted_at DATETIME NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_group_year_segment_bid (group_id, year_no, market_code, order_type),
			KEY idx_year_segment_bid (year_no, market_code, order_type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_market_selection_order (
			id BIGINT NOT NULL AUTO_INCREMENT,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			sequence_no INT NOT NULL,
			group_id BIGINT NOT NULL,
			market_investment DECIMAL(18,2) NOT NULL DEFAULT 0,
			previous_market_order_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
			is_market_leader TINYINT(1) NOT NULL DEFAULT 0,
			rank_basis_json JSON NOT NULL,
			selection_status VARCHAR(32) NOT NULL DEFAULT 'WAITING',
			selected_order_id BIGINT NULL,
			selected_at DATETIME NULL,
			skipped_by_admin_id BIGINT NULL,
			skipped_reason VARCHAR(255) NULL,
			skipped_at DATETIME NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_market_sequence (year_no, market_code, order_type, sequence_no),
			UNIQUE KEY uk_group_market_sequence (year_no, market_code, order_type, group_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
		`CREATE TABLE IF NOT EXISTS sg_group_order_selection (
			id BIGINT NOT NULL AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			year_no INT NOT NULL,
			market_code VARCHAR(32) NOT NULL,
			order_type VARCHAR(32) NOT NULL,
			order_id BIGINT NOT NULL,
			selection_status VARCHAR(32) NOT NULL DEFAULT 'SELECTED',
			delivery_status VARCHAR(32) NOT NULL DEFAULT 'SELECTED',
			delivered_stage_code VARCHAR(16) NULL,
			delivered_at DATETIME NULL,
			selected_at DATETIME NOT NULL,
			creator VARCHAR(64) NOT NULL DEFAULT 'system',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updater VARCHAR(64) NOT NULL DEFAULT 'system',
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY uk_group_year_segment_selection (group_id, year_no, market_code, order_type),
			UNIQUE KEY uk_order_selected (order_id),
			KEY idx_group_year_order_selection (group_id, year_no)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "duplicate column name") &&
				!strings.Contains(err.Error(), "check that column/key exists") &&
				!strings.Contains(errText, "duplicate key name") {
				t.Fatalf("ensure order tables: %v", err)
			}
		}
	}

	if !db.Migrator().HasColumn(&entity.OrderGenerationBatch{}, "forecast_snapshot_json") {
		if err := db.Migrator().AddColumn(&entity.OrderGenerationBatch{}, "ForecastSnapshot"); err != nil {
			t.Fatalf("ensure order batch forecast snapshot column: %v", err)
		}
	}
	if !db.Migrator().HasColumn(&entity.OrderMarketConfig{}, "market_investment_limit") {
		if err := db.Migrator().AddColumn(&entity.OrderMarketConfig{}, "MarketInvestmentLimit"); err != nil {
			t.Fatalf("ensure order market config investment limit column: %v", err)
		}
	}
	ensureIntegrationOrderTemplateColumns(t, db)
}

func ensureIntegrationOrderTemplateColumns(t *testing.T, db *gorm.DB) {
	t.Helper()

	entitiesWithColumns := map[any][]string{
		&entity.OrderGenerationConfig{}: {"OrderTemplateVersion"},
		&entity.OrderForecastControl{}:  {"OrderTemplateVersion"},
		&entity.OrderMarketForecast{}:   {"OrderTemplateVersion"},
		&entity.OrderMarketConfig{}:     {"OrderTemplateVersion"},
		&entity.OrderGenerationBatch{}:  {"OrderTemplateVersion"},
		&entity.OrderPool{}:             {"OrderTemplateVersion", "OrderPayloadJSON"},
		&entity.MarketBiddingState{}:    {"OrderTemplateVersion"},
		&entity.GroupMarketBid{}:        {"OrderTemplateVersion"},
		&entity.MarketSelectionOrder{}:  {"OrderTemplateVersion"},
		&entity.GroupOrderSelection{}:   {"OrderTemplateVersion"},
	}

	for item, columns := range entitiesWithColumns {
		for _, column := range columns {
			if db.Migrator().HasColumn(item, column) {
				continue
			}
			if err := db.Migrator().AddColumn(item, column); err != nil {
				t.Fatalf("ensure order template column %s on %T: %v", column, item, err)
			}
		}
	}
}

func ensureIntegrationRollbackTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	entitiesWithColumns := map[any][]string{
		&entity.GroupYearState{}: {
			"RollbackPending",
			"RollbackTargetYearNo",
			"RollbackTargetStageCode",
			"RollbackLogID",
		},
		&entity.GroupSummarySnapshot{}: {
			"InvalidatedByRollbackID",
			"InvalidatedAt",
		},
		&entity.GroupAdjustment{}: {
			"Effective",
			"InvalidatedByRollbackID",
			"InvalidReason",
			"InvalidatedAt",
		},
		&entity.AdminUnlockLog{}: {
			"UnlockTargetType",
			"TargetStageCode",
			"SafetySnapshotID",
		},
		&entity.GroupOrderSelection{}: {
			"DeliveryEffective",
			"InvalidatedByRollbackID",
			"InvalidatedAt",
		},
	}

	for item, columns := range entitiesWithColumns {
		for _, column := range columns {
			if db.Migrator().HasColumn(item, column) {
				continue
			}
			if err := db.Migrator().AddColumn(item, column); err != nil {
				t.Fatalf("ensure rollback column %s on %T: %v", column, item, err)
			}
		}
	}

	if err := db.AutoMigrate(
		&entity.StateSnapshot{},
		&entity.StateSnapshotPayload{},
		&entity.RollbackLog{},
	); err != nil {
		t.Fatalf("ensure rollback tables: %v", err)
	}

	if err := db.Exec("UPDATE sg_group_adjustment SET effective = 1 WHERE effective IS NULL").Error; err != nil {
		t.Fatalf("backfill adjustment effective: %v", err)
	}
	if err := db.Exec("UPDATE sg_group_order_selection SET delivery_effective = 1 WHERE delivery_effective IS NULL").Error; err != nil {
		t.Fatalf("backfill order delivery effective: %v", err)
	}
}

func buildPlayerOperatingServices(db *gorm.DB) (*PlayerOperatingCommandService, *PlayerOperatingQueryService) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	groupYearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	initialBaselineRepo := repository.NewInitialBaselineRepository(db)
	reportRepo := repository.NewReportRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(db)
	playerNoticeService := NewPlayerNoticeService(noticeRepo, adjustmentRepo)

	commandService := NewPlayerOperatingCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		playerNoticeService,
		nil,
	)
	queryService := NewPlayerOperatingQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		assembler.NewPlayerOperatingAssembler(),
		playerNoticeService,
		nil,
	)

	return commandService, queryService
}

func ensureGameConfigExists(t *testing.T, ctx context.Context, tx *gorm.DB) {
	t.Helper()

	var count int64
	if err := tx.WithContext(ctx).Model(&entity.GameConfig{}).Count(&count).Error; err != nil {
		t.Fatalf("count game config: %v", err)
	}
	if count > 0 {
		return
	}

	now := time.Now()
	item := entity.GameConfig{
		FinalYear:                8,
		CurrentOpenYear:          0,
		EditionCode:              GameEditionVIPServiceV1,
		EditionName:              GameEditionVIPServiceV1Name,
		RuleVersion:              FormulaVersionCommonV1,
		TemplateVersion:          TemplateVersionVIPServiceV1,
		OperatingTemplateVersion: OperatingTemplateVersionVIPServiceV1,
		ReportTemplateVersion:    ReportTemplateVersionVIPServiceV1,
		OrderTemplateVersion:     OrderTemplateVersionVIPServiceV1,
		ProcessRuleVersion:       ProcessRuleVersionCommonV1,
		InitialBaselineSubmitted: true,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create fallback game config: %v", err)
	}
}

func createIntegrationOperatingFixtures(t *testing.T, ctx context.Context, tx *gorm.DB) int64 {
	t.Helper()

	now := time.Now()
	uniqueSeed := nextIntegrationUniqueSeed()
	group := entity.Group{
		GroupNo:        integrationGroupNoFromSeed(uniqueSeed),
		GroupCode:      fmt.Sprintf("IT_OPERATING_%d", uniqueSeed),
		GroupName:      "integration operating fixture group",
		BusinessStatus: enum.BusinessStatusNormal,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&group).Error; err != nil {
		t.Fatalf("create group fixture: %v", err)
	}

	yearState := entity.GroupYearState{
		GroupID:                   group.ID,
		YearNo:                    0,
		YearType:                  enum.YearTypeDemo,
		YearStatus:                enum.YearStatusOperating,
		StageStatus:               enum.StageStatusQ1Open,
		ReportStatus:              enum.ReportStatusLocked,
		SummaryEffective:          false,
		LatestStageSubmitVersion:  0,
		LatestReportSubmitVersion: 0,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&yearState).Error; err != nil {
		t.Fatalf("create year state fixture: %v", err)
	}

	baselineJSON, err := json.Marshal(buildIntegrationBaselinePayload())
	if err != nil {
		t.Fatalf("marshal baseline payload: %v", err)
	}

	submittedAt := now
	baseline := entity.InitialBaseline{
		GroupID:         group.ID,
		BaselinePayload: baselineJSON,
		Submitted:       true,
		SubmitterID:     int64Ptr(1),
		SubmittedAt:     &submittedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&baseline).Error; err != nil {
		t.Fatalf("create baseline fixture: %v", err)
	}

	return group.ID
}

func buildIntegrationBaselinePayload() payload.BaselinePayload {
	return payload.BaselinePayload{
		BaselineSalesRevenue:         32,
		BaselineDirectCost:           15,
		BaselineComprehensiveCost:    13,
		BaselineDepreciation:         1,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   2,
		BaselineIncomeTax:            1,
		BaselineWorkInConstruction:   0,
		BaselineFactoryAsset:         40,
		BaselineLineResidual:         3,
		BaselineDepreciableAsset:     0,
		BaselineCash:                 36,
		BaselineReceivable:           0,
		BaselineWorkInProgress:       6,
		BaselineFinishedGoods:        4,
		BaselineRawMaterials:         1,
		BaselineShortTermLoan:        20,
		BaselineLongTermLoan:         0,
		BaselineShareCapital:         50,
		BaselineRetainedEarnings:     17,
	}
}

func buildValidQ1OperatingPayload() payload.OperatingPayload {
	p := payload.NewOperatingPayload()

	p.Beginning.TaxAndPlanning = map[string]any{
		"taxRate":       0,
		"marketBidCost": 0,
	}
	p.Beginning.MarketBid = []map[string]any{
		{
			"market":           "demo",
			"marketInvestment": 0,
			"orderAmount":      0,
		},
	}

	p.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"q1": {
			"shortTermRepayment": 0,
			"shortTermInterest":  0,
			"newShortTermLoan":   0,
		},
	}
	p.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"q1": {
			"materialPayment": 0,
		},
	}
	p.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"q1": {
			"changeProductCost":        0,
			"lineDismantleCost":        0,
			"lineSaleValue":            0,
			"newLineInstall":           0,
			"constructionToFixed":      0,
			"depreciableAssetIncrease": 0,
		},
	}
	p.Quarter.HumanResource = payload.OperatingQuarterMap{
		"q1": {
			"humanResourceCost": 0,
		},
	}
	p.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"q1": {
			"salaryAndProductionCost": 0,
		},
	}
	p.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"q1": {
			"researchCost":         0,
			"managementSystemCost": 0,
		},
	}
	p.Quarter.SupplyChainOrderRecord = payload.OperatingQuarterMap{
		"q1": {
			"basicProduct":       0,
			"standardProduct":    0,
			"precisionProduct":   0,
			"intelligentProduct": 0,
		},
	}
	p.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"q1": {
			"receivableRecovered": 0,
		},
	}
	p.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q1": {
			"salesRevenue":      0,
			"directCost":        0,
			"managementSalary":  0,
			"deliveryQuantity":  0,
			"receivableBalance": 0,
		},
	}
	p.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"q1": {
			"discountExpense":     0,
			"extraExpensePenalty": 0,
			"extraIncomeReward":   0,
		},
	}

	return p.Normalize()
}

func resolveProjectPath(t *testing.T, paths ...string) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("resolve project path: runtime caller not available")
	}

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	fullPath := filepath.Join(append([]string{root}, paths...)...)
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("resolve project path %s: %v", fullPath, err)
	}
	return fullPath
}

func int64Ptr(value int64) *int64 {
	return &value
}
