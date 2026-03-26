package operating

import (
	"fmt"
	"strings"

	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
	"sandbox-game/internal/state"
)

// ValidationIssue 表示规则层发现的业务问题。
type ValidationIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult 表示规则校验结果。
type ValidationResult struct {
	Passed bool              `json:"passed"`
	Issues []ValidationIssue `json:"issues"`
}

type requiredScope struct {
	path       string
	label      string
	value      any
	allowEmpty bool
}

// Validator 是经营页校验入口。
type Validator struct {
	guard *state.TransitionGuard
}

// NewValidator 创建经营校验器。
func NewValidator() *Validator {
	return &Validator{
		guard: state.NewTransitionGuard(),
	}
}

// ValidateStageSubmit 校验经营阶段提交的最小前置条件。
func (v *Validator) ValidateStageSubmit(ctx calcctx.CalculationContext, stageCode string) ValidationResult {
	issues := make([]ValidationIssue, 0)

	if err := ctx.Validate(); err != nil {
		issues = append(issues, ValidationIssue{
			Field:   "context",
			Message: err.Error(),
		})
		return failedResult(issues)
	}

	if !ctx.HasCarryForwardSource() {
		issues = append(issues, ValidationIssue{
			Field:   "carryForwardSource",
			Message: "缺少跨年承接来源，无法进入经营计算",
		})
	}

	if ctx.OperatingPayload == nil {
		issues = append(issues, ValidationIssue{
			Field:   "operatingPayload",
			Message: "经营页负载不能为空",
		})
	}

	if !v.guard.CanSubmitStage(ctx.State, stageCode) {
		issues = append(issues, ValidationIssue{
			Field:   "stageCode",
			Message: "当前阶段不允许提交该阶段动作",
		})
	}

	if len(issues) > 0 {
		return failedResult(issues)
	}

	issues = append(issues, validateRequiredScopes(ctx.OperatingPayload.Normalize(), stageCode)...)
	if len(issues) > 0 {
		return failedResult(issues)
	}

	return ValidationResult{
		Passed: true,
		Issues: []ValidationIssue{},
	}
}

func validateRequiredScopes(operatingPayload payload.OperatingPayload, stageCode string) []ValidationIssue {
	scopes := buildRequiredScopes(operatingPayload, stageCode)
	issues := make([]ValidationIssue, 0)
	for _, scope := range scopes {
		issues = append(issues, validateScope(scope)...)
	}
	return issues
}

func buildRequiredScopes(operatingPayload payload.OperatingPayload, stageCode string) []requiredScope {
	stageCode = strings.ToUpper(strings.TrimSpace(stageCode))

	switch stageCode {
	case state.StageCodeQ1:
		return buildQuarterStageScopes(operatingPayload, "q1", true)
	case state.StageCodeQ2:
		return buildQuarterStageScopes(operatingPayload, "q2", false)
	case state.StageCodeQ3:
		return buildQuarterStageScopes(operatingPayload, "q3", false)
	case state.StageCodeQ4:
		return buildQuarterStageScopes(operatingPayload, "q4", false)
	case state.StageCodeYearEnd:
		return []requiredScope{
			{
				path:  "yearEnd.longTermLoan",
				label: "年末长期贷款区",
				value: operatingPayload.YearEnd.LongTermLoan,
			},
			{
				path:  "yearEnd.assetAdjustment",
				label: "年末资产调整区",
				value: operatingPayload.YearEnd.AssetAdjustment,
			},
		}
	default:
		return []requiredScope{}
	}
}

func buildQuarterStageScopes(operatingPayload payload.OperatingPayload, quarterKey string, includeBeginning bool) []requiredScope {
	quarterLabel := strings.ToUpper(quarterKey)
	scopes := make([]requiredScope, 0, 11)

	if includeBeginning {
		scopes = append(scopes,
			requiredScope{
				path:  "beginning.taxAndPlanning",
				label: "年初规划区",
				value: operatingPayload.Beginning.TaxAndPlanning,
			},
			requiredScope{
				path:  "beginning.marketBid",
				label: "年初市场竞标区",
				value: operatingPayload.Beginning.MarketBid,
			},
		)
	}

	scopes = append(scopes,
		requiredScope{
			path:  fmt.Sprintf("quarter.shortTermLoan.%s", quarterKey),
			label: quarterLabel + "短期贷款区",
			value: findQuarterValue(operatingPayload.Quarter.ShortTermLoan, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.materialPayment.%s", quarterKey),
			label: quarterLabel + "材料费区",
			value: findQuarterValue(operatingPayload.Quarter.MaterialPayment, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.productionLineAdjustment.%s", quarterKey),
			label: quarterLabel + "生产线调整区",
			value: findQuarterValue(operatingPayload.Quarter.ProductionLineAdjust, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.humanResource.%s", quarterKey),
			label: quarterLabel + "人力资源区",
			value: findQuarterValue(operatingPayload.Quarter.HumanResource, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.salaryAndProduction.%s", quarterKey),
			label: quarterLabel + "工资与生产区",
			value: findQuarterValue(operatingPayload.Quarter.SalaryAndProduction, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.researchAndManagement.%s", quarterKey),
			label: quarterLabel + "研发与管理区",
			value: findQuarterValue(operatingPayload.Quarter.ResearchAndManagement, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.receivableUpdate.%s", quarterKey),
			label: quarterLabel + "应收更新区",
			value: findQuarterValue(operatingPayload.Quarter.ReceivableUpdate, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("quarter.deliverySettlement.%s", quarterKey),
			label: quarterLabel + "交货结算区",
			value: findQuarterValue(operatingPayload.Quarter.DeliverySettlement, quarterKey),
		},
		requiredScope{
			path:  fmt.Sprintf("extra.incomeAndPenalty.%s", quarterKey),
			label: quarterLabel + "其他收支区",
			value: findQuarterValue(operatingPayload.Extra.IncomeAndPenalty, quarterKey),
		},
	)

	return scopes
}

func findQuarterValue(source payload.OperatingQuarterMap, quarterKey string) any {
	if len(source) == 0 {
		return nil
	}
	target := normalizeQuarterKey(quarterKey)
	for key, value := range source {
		if normalizeQuarterKey(key) == target {
			return value
		}
	}
	return nil
}

func normalizeQuarterKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validateScope(scope requiredScope) []ValidationIssue {
	if isMissingContainer(scope.value) {
		return []ValidationIssue{{
			Field:   scope.path,
			Message: scope.label + "未填写",
		}}
	}

	issues := make([]ValidationIssue, 0)
	walkMissingValues(scope.path, scope.label, scope.value, &issues)
	return issues
}

func walkMissingValues(path string, label string, value any, issues *[]ValidationIssue) {
	switch typed := value.(type) {
	case map[string]any:
		if len(typed) == 0 {
			*issues = append(*issues, ValidationIssue{
				Field:   path,
				Message: label + "未填写",
			})
			return
		}
		for key, child := range typed {
			childPath := path + "." + key
			if isMissingLeaf(child) {
				*issues = append(*issues, ValidationIssue{
					Field:   childPath,
					Message: childPath + " 未填写",
				})
				continue
			}
			walkMissingValues(childPath, label, child, issues)
		}
	case []map[string]any:
		if len(typed) == 0 {
			*issues = append(*issues, ValidationIssue{
				Field:   path,
				Message: label + "未填写",
			})
			return
		}
		for index, item := range typed {
			childPath := fmt.Sprintf("%s[%d]", path, index)
			if len(item) == 0 {
				*issues = append(*issues, ValidationIssue{
					Field:   childPath,
					Message: childPath + " 未填写",
				})
				continue
			}
			walkMissingValues(childPath, label, item, issues)
		}
	case []any:
		if len(typed) == 0 {
			*issues = append(*issues, ValidationIssue{
				Field:   path,
				Message: label + "未填写",
			})
			return
		}
		for index, item := range typed {
			childPath := fmt.Sprintf("%s[%d]", path, index)
			if isMissingLeaf(item) {
				*issues = append(*issues, ValidationIssue{
					Field:   childPath,
					Message: childPath + " 未填写",
				})
				continue
			}
			walkMissingValues(childPath, label, item, issues)
		}
	default:
		if isMissingLeaf(typed) {
			*issues = append(*issues, ValidationIssue{
				Field:   path,
				Message: path + " 未填写",
			})
		}
	}
}

func isMissingContainer(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case map[string]any:
		return len(typed) == 0
	case []map[string]any:
		return len(typed) == 0
	case []any:
		return len(typed) == 0
	default:
		return false
	}
}

func isMissingLeaf(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	default:
		return false
	}
}

func failedResult(issues []ValidationIssue) ValidationResult {
	return ValidationResult{
		Passed: false,
		Issues: issues,
	}
}
