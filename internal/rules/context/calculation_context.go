package context

import (
	"fmt"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/state"
)

// CalculationContext 统一承接规则层计算和校验所需上下文。
type CalculationContext struct {
	Group               entity.Group
	YearState           entity.GroupYearState
	GameConfig          entity.GameConfig
	State               state.RuntimeState
	InitialBaseline     *payload.BaselinePayload
	PreviousReport      *payload.ReportComputedPayload
	OperatingPayload    *payload.OperatingPayload
	ReportManualPayload *payload.ReportManualPayload
}

// NewCalculationContext 根据实体构建规则层上下文。
func NewCalculationContext(group entity.Group, yearState entity.GroupYearState, gameConfig entity.GameConfig) CalculationContext {
	return CalculationContext{
		Group:      group,
		YearState:  yearState,
		GameConfig: gameConfig,
		State:      state.NewRuntimeStateFromEntities(group, yearState),
	}
}

// Validate 负责校验上下文基础结构是否可进入规则层。
func (c CalculationContext) Validate() error {
	if c.Group.ID <= 0 {
		return fmt.Errorf("invalid calculation context: group id is required")
	}
	if c.YearState.GroupID <= 0 {
		return fmt.Errorf("invalid calculation context: yearState groupId is required")
	}
	if c.Group.ID != c.YearState.GroupID {
		return fmt.Errorf("invalid calculation context: group/yearState mismatch")
	}
	if c.YearState.YearNo != c.State.YearNo {
		return fmt.Errorf("invalid calculation context: state yearNo mismatch")
	}
	if c.GameConfig.FinalYear < 0 {
		return fmt.Errorf("invalid calculation context: finalYear cannot be negative")
	}
	if err := state.NewStateMachine().Validate(c.State); err != nil {
		return fmt.Errorf("invalid calculation context: %w", err)
	}

	return nil
}

// IsDemoYear 判断当前是否为 0 年引导年。
func (c CalculationContext) IsDemoYear() bool {
	return c.YearState.YearNo == 0
}

// HasCarryForwardSource 判断当前是否具备跨年承接所需来源。
func (c CalculationContext) HasCarryForwardSource() bool {
	if c.IsDemoYear() {
		return c.InitialBaseline != nil
	}
	return c.PreviousReport != nil
}

// WithInitialBaseline 返回带初始基线的上下文副本。
func (c CalculationContext) WithInitialBaseline(baseline *payload.BaselinePayload) CalculationContext {
	c.InitialBaseline = baseline
	return c
}

// WithPreviousReport 返回带上年财报结果的上下文副本。
func (c CalculationContext) WithPreviousReport(previous *payload.ReportComputedPayload) CalculationContext {
	c.PreviousReport = previous
	return c
}

// WithOperatingPayload 返回带经营页负载的上下文副本。
func (c CalculationContext) WithOperatingPayload(operating *payload.OperatingPayload) CalculationContext {
	c.OperatingPayload = operating
	return c
}

// WithReportManualPayload 返回带财报手工项的上下文副本。
func (c CalculationContext) WithReportManualPayload(reportManual *payload.ReportManualPayload) CalculationContext {
	c.ReportManualPayload = reportManual
	return c
}
