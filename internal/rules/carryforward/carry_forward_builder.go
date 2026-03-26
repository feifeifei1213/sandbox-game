package carryforward

import (
	"fmt"

	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

// CarryForwardResult 表示从上年财报或初始基线承接到下一年经营的关键字段。
type CarryForwardResult struct {
	PreviousIncomeTax         float64 `json:"previousIncomeTax"`
	PreviousShortTermLoan     float64 `json:"previousShortTermLoan"`
	PreviousLongTermLoan      float64 `json:"previousLongTermLoan"`
	PreviousEquipmentResidual float64 `json:"previousEquipmentResidual"`
	PreviousDepreciableAsset  float64 `json:"previousDepreciableAsset"`
	PreviousCash              float64 `json:"previousCash"`
	PreviousReceivable        float64 `json:"previousReceivable"`
	ShareholderCapital        float64 `json:"shareholderCapital"`
	RetainedEarnings          float64 `json:"retainedEarnings"`
}

// Builder 是跨年承接统一入口。
type Builder struct{}

// NewBuilder 创建跨年承接构建器。
func NewBuilder() *Builder {
	return &Builder{}
}

// Build 根据当前年份上下文构建跨年承接结果。
func (b *Builder) Build(ctx calcctx.CalculationContext) (CarryForwardResult, error) {
	if err := ctx.Validate(); err != nil {
		return CarryForwardResult{}, err
	}

	if ctx.IsDemoYear() {
		if ctx.InitialBaseline == nil {
			return CarryForwardResult{}, fmt.Errorf("missing initial baseline for demo year")
		}

		return CarryForwardResult{
			PreviousIncomeTax:         ctx.InitialBaseline.BaselineIncomeTax,
			PreviousShortTermLoan:     ctx.InitialBaseline.BaselineShortTermLoan,
			PreviousLongTermLoan:      ctx.InitialBaseline.BaselineLongTermLoan,
			PreviousEquipmentResidual: ctx.InitialBaseline.BaselineLineResidual,
			PreviousDepreciableAsset:  ctx.InitialBaseline.BaselineDepreciableAsset,
			PreviousCash:              ctx.InitialBaseline.BaselineCash,
			PreviousReceivable:        ctx.InitialBaseline.BaselineReceivable,
			ShareholderCapital:        ctx.InitialBaseline.BaselineShareCapital,
			RetainedEarnings:          ctx.InitialBaseline.BaselineRetainedEarnings + calculateBaselineNetProfit(*ctx.InitialBaseline),
		}, nil
	}

	if ctx.PreviousReport == nil {
		return CarryForwardResult{}, fmt.Errorf("missing previous report for formal year")
	}

	return CarryForwardResult{
		PreviousIncomeTax:         ctx.PreviousReport.ReportIncomeTax,
		PreviousShortTermLoan:     ctx.PreviousReport.ReportShortTermLiability,
		PreviousLongTermLoan:      ctx.PreviousReport.ReportLongTermLiability,
		PreviousEquipmentResidual: ctx.PreviousReport.ReportLineResidual,
		PreviousDepreciableAsset:  ctx.PreviousReport.ReportDepreciableAsset,
		PreviousCash:              ctx.PreviousReport.ReportCash,
		PreviousReceivable:        ctx.PreviousReport.ReportReceivable,
		ShareholderCapital:        ctx.PreviousReport.ReportShareCapital,
		RetainedEarnings:          ctx.PreviousReport.ReportRetainedEarnings + ctx.PreviousReport.ReportNetProfit,
	}, nil
}

func calculateBaselineNetProfit(baseline payload.BaselinePayload) float64 {
	grossProfit := baseline.BaselineSalesRevenue - baseline.BaselineDirectCost
	operatingProfit := grossProfit - baseline.BaselineComprehensiveCost - baseline.BaselineDepreciation
	preTaxProfit := operatingProfit - baseline.BaselineFinanceIncomeExpense + baseline.BaselineExtraIncomeExpense
	return preTaxProfit - baseline.BaselineIncomeTax
}
