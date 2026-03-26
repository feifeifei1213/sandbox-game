package summary

import (
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

// SummaryResult 表示汇总页所需的最小结果结构。
type SummaryResult struct {
	Revenue          float64 `json:"revenue"`
	Profit           float64 `json:"profit"`
	Equity           float64 `json:"equity"`
	BusinessStatus   string  `json:"businessStatus"`
	RankingValue     float64 `json:"rankingValue"`
	SummaryEffective bool    `json:"summaryEffective"`
}

// Builder 是汇总构建统一入口。
type Builder struct{}

// NewBuilder 创建汇总构建器。
func NewBuilder() *Builder {
	return &Builder{}
}

// Build 根据财报结果生成汇总快照载体。
func (b *Builder) Build(ctx calcctx.CalculationContext, computed payload.ReportComputedPayload) (SummaryResult, error) {
	if err := ctx.Validate(); err != nil {
		return SummaryResult{}, err
	}

	result := SummaryResult{
		Revenue:        computed.ReportSalesRevenue,
		Profit:         computed.ReportNetProfit,
		Equity:         computed.ReportTotalEquity,
		BusinessStatus: ctx.State.BusinessStatus,
		RankingValue:   computed.ReportTotalEquity,
	}

	result.SummaryEffective = ctx.State.YearStatus == enum.YearStatusCompleted && !ctx.IsDemoYear()

	return result, nil
}
