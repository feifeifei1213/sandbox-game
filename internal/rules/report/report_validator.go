package report

import (
	"math"

	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
	"sandbox-game/internal/rules/operating"
	"sandbox-game/internal/state"
)

const balanceTolerance = 0.000001

// BalanceCheckResult 表示资产负债平衡校验结果。
type BalanceCheckResult struct {
	Passed bool    `json:"passed"`
	Gap    float64 `json:"gap"`
}

// Validator 是财报校验入口。
type Validator struct {
	guard *state.TransitionGuard
}

// NewValidator 创建财报校验器。
func NewValidator() *Validator {
	return &Validator{
		guard: state.NewTransitionGuard(),
	}
}

// ValidateDraft 校验财报草稿保存时的最小约束。
// 首版只校验已填写值的合法性，不要求手工项完整。
func (v *Validator) ValidateDraft(manual *payload.ReportManualPayload) operating.ValidationResult {
	if manual == nil {
		return operating.ValidationResult{
			Passed: true,
			Issues: []operating.ValidationIssue{},
		}
	}

	issues := validateManualPayloadForDraft(manual)
	if len(issues) > 0 {
		return operating.ValidationResult{
			Passed: false,
			Issues: issues,
		}
	}

	return operating.ValidationResult{
		Passed: true,
		Issues: []operating.ValidationIssue{},
	}
}

// ValidateSubmit 校验财报提交前提、手工项完整性、税率合法性和平衡关系。
func (v *Validator) ValidateSubmit(ctx calcctx.CalculationContext, computed payload.ReportComputedPayload) operating.ValidationResult {
	issues := make([]operating.ValidationIssue, 0)

	if err := ctx.Validate(); err != nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "context",
			Message: err.Error(),
		})
		return operating.ValidationResult{Passed: false, Issues: issues}
	}

	if !v.guard.CanSubmitReport(ctx.State) {
		issues = append(issues, operating.ValidationIssue{
			Field:   "reportStatus",
			Message: "当前状态不允许提交财报",
		})
	}

	issues = append(issues, validateManualPayloadForSubmit(ctx.ReportManualPayload)...)

	balance := v.CheckBalance(computed)
	if !balance.Passed {
		issues = append(issues, operating.ValidationIssue{
			Field:   "balance",
			Message: "总资产必须等于总负债和权益",
		})
	}

	if len(issues) > 0 {
		return operating.ValidationResult{
			Passed: false,
			Issues: issues,
		}
	}

	return operating.ValidationResult{
		Passed: true,
		Issues: []operating.ValidationIssue{},
	}
}

// CheckBalance 校验财报平衡关系。
func (v *Validator) CheckBalance(computed payload.ReportComputedPayload) BalanceCheckResult {
	gap := computed.BalanceGap()
	return BalanceCheckResult{
		Passed: math.Abs(gap) <= balanceTolerance,
		Gap:    gap,
	}
}

func validateManualPayloadForSubmit(manual *payload.ReportManualPayload) []operating.ValidationIssue {
	if manual == nil {
		return []operating.ValidationIssue{{
			Field:   "reportManualPayload",
			Message: "财报手工项不能为空",
		}}
	}

	issues := validateManualPayloadForDraft(manual)
	if manual.WorkInProgress == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "workInProgress",
			Message: "在制品不能为空",
		})
	}
	if manual.FinishedGoods == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "finishedGoods",
			Message: "成品不能为空",
		})
	}
	if manual.RawMaterials == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "rawMaterials",
			Message: "材料不能为空",
		})
	}
	if manual.IncomeTaxRate == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "incomeTaxRate",
			Message: "所得税税率不能为空",
		})
	}
	if manual.EnterpriseCertificationScore == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "enterpriseCertificationScore",
			Message: "企业认证得分不能为空",
		})
	}
	if manual.ProductionHumanScore == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "productionHumanScore",
			Message: "最佳生产人力总监得分不能为空",
		})
	}
	if manual.ClosingSpeedScore == nil {
		issues = append(issues, operating.ValidationIssue{
			Field:   "closingSpeedScore",
			Message: "关账速度得分不能为空",
		})
	}

	return issues
}

func validateManualPayloadForDraft(manual *payload.ReportManualPayload) []operating.ValidationIssue {
	issues := make([]operating.ValidationIssue, 0)
	if manual.IncomeTaxRate != nil && !isAllowedIncomeTaxRate(*manual.IncomeTaxRate) {
		issues = append(issues, operating.ValidationIssue{
			Field:   "incomeTaxRate",
			Message: "所得税税率不在允许值范围内",
		})
	}
	return issues
}

func isAllowedIncomeTaxRate(value float64) bool {
	for _, allowed := range payload.AllowedIncomeTaxRates() {
		if math.Abs(value-allowed) <= balanceTolerance {
			return true
		}
	}
	return false
}
