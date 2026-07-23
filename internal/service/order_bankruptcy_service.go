package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/repository"
)

// invalidateOrderParticipationAfterBankruptcy 只失效尚未完成的轮次机会。
// 已选订单、已完成轮次和订单归属保持不变；若破产小组正好轮到选单，则继续推进当前标段。
func invalidateOrderParticipationAfterBankruptcy(
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	operatorID int64,
	operatorName string,
	operateTime time.Time,
) error {
	stateRepo := repository.NewMarketBiddingStateRepository(tx)
	sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
	poolRepo := repository.NewOrderPoolRepository(tx)

	// 订单操作统一先锁标段状态再锁顺序记录；破产联动保持同样的锁顺序，避免与选单并发死锁。
	currentSegments, err := stateRepo.ListSelectingByCurrentGroupForUpdate(ctx, yearNo, groupID)
	if err != nil {
		return fmt.Errorf("lock bankrupt group current order segments: %w", err)
	}
	if _, err := sequenceRepo.MarkUnfinishedBankrupt(ctx, groupID, yearNo, operatorName, operateTime); err != nil {
		return fmt.Errorf("invalidate bankrupt group order rounds: %w", err)
	}
	for _, segment := range currentSegments {
		if _, _, err := advanceOrderSegment(ctx, tx, stateRepo, sequenceRepo, poolRepo, segment, operatorID, operatorName, operateTime); err != nil {
			return fmt.Errorf("advance order segment after bankruptcy: %w", err)
		}
	}
	return nil
}
