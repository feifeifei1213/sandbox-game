package service

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func TestGroupAdjustmentRevisionGetEmptyDoesNotEmitRecordNotFound(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	groupID := int64(987654321)
	yearNo := 77
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Delete(&entity.GroupAdjustmentRevision{}).Error; err != nil {
		t.Fatalf("clear adjustment revision: %v", err)
	}

	var loggerOutput bytes.Buffer
	queryDB := tx.Session(&gorm.Session{
		Logger: gormlogger.New(
			log.New(&loggerOutput, "", 0),
			gormlogger.Config{
				LogLevel: gormlogger.Warn,
				Colorful: false,
			},
		),
	})

	revision, err := repository.NewGroupAdjustmentRevisionRepository(queryDB).Get(ctx, groupID, yearNo)
	if err != nil {
		t.Fatalf("expected empty adjustment revision to be treated as zero revision, got err=%v", err)
	}
	if revision != 0 {
		t.Fatalf("expected empty adjustment revision to return 0, got %d", revision)
	}
	if strings.Contains(loggerOutput.String(), "record not found") {
		t.Fatalf("expected empty adjustment revision query to avoid GORM record-not-found noise, got log: %s", loggerOutput.String())
	}
}
