package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"sandbox-game/internal/model/payload"
)

const manualIntegerPreviewLimit = 3

var ErrManualNumberNotInteger = fmt.Errorf("manual number not integer")

type ManualIntegerIssue struct {
	Path  string
	Value string
}

type ManualIntegerValidationError struct {
	Scope  string
	Issues []ManualIntegerIssue
}

func (e *ManualIntegerValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return ErrManualNumberNotInteger.Error()
	}

	parts := make([]string, 0, minInt(len(e.Issues), manualIntegerPreviewLimit))
	for index, issue := range e.Issues {
		if index >= manualIntegerPreviewLimit {
			break
		}
		parts = append(parts, fmt.Sprintf("%s=%s", issue.Path, issue.Value))
	}
	if len(e.Issues) > manualIntegerPreviewLimit {
		parts = append(parts, fmt.Sprintf("等%d项", len(e.Issues)))
	}

	scope := strings.TrimSpace(e.Scope)
	if scope == "" {
		scope = "手工数字"
	}
	return fmt.Sprintf("%s必须填写整数，发现小数：%s", scope, strings.Join(parts, "、"))
}

func (e *ManualIntegerValidationError) Unwrap() error {
	return ErrManualNumberNotInteger
}

func validateOperatingManualIntegers(value payload.OperatingPayload) error {
	normalized := value.Normalize().WithoutDerivedValues()
	issues := make([]ManualIntegerIssue, 0)

	collectManualIntegerIssues("beginning.taxAndPlanning", normalized.Beginning.TaxAndPlanning, &issues)
	collectManualIntegerIssues("beginning.marketBid", normalized.Beginning.MarketBid, &issues)
	collectManualIntegerIssues("quarter.shortTermLoan", normalized.Quarter.ShortTermLoan, &issues)
	collectManualIntegerIssues("quarter.materialPayment", normalized.Quarter.MaterialPayment, &issues)
	collectManualIntegerIssues("quarter.productionLineAdjustment", normalized.Quarter.ProductionLineAdjust, &issues)
	collectManualIntegerIssues("quarter.humanResource", normalized.Quarter.HumanResource, &issues)
	collectManualIntegerIssues("quarter.salaryAndProduction", normalized.Quarter.SalaryAndProduction, &issues)
	collectManualIntegerIssues("quarter.researchAndManagement", normalized.Quarter.ResearchAndManagement, &issues)
	collectManualIntegerIssues("quarter.receivableUpdate", normalized.Quarter.ReceivableUpdate, &issues)
	collectManualIntegerIssues("quarter.deliverySettlement", normalized.Quarter.DeliverySettlement, &issues)
	collectManualIntegerIssues("yearEnd.longTermLoan", normalized.YearEnd.LongTermLoan, &issues)
	collectManualIntegerIssues("yearEnd.assetAdjustment", normalized.YearEnd.AssetAdjustment, &issues)
	collectManualIntegerIssues("extra.incomeAndPenalty", normalized.Extra.IncomeAndPenalty, &issues)

	return buildManualIntegerError("经营页手工数字", issues)
}

func validateReportManualIntegers(value payload.ReportManualPayload) error {
	issues := make([]ManualIntegerIssue, 0)
	collectNullableFloatIntegerIssue("workInProgress", value.WorkInProgress, &issues)
	collectNullableFloatIntegerIssue("finishedGoods", value.FinishedGoods, &issues)
	collectNullableFloatIntegerIssue("rawMaterials", value.RawMaterials, &issues)
	collectNullableFloatIntegerIssue("enterpriseCertificationScore", value.EnterpriseCertificationScore, &issues)
	collectNullableFloatIntegerIssue("productionHumanScore", value.ProductionHumanScore, &issues)
	collectNullableFloatIntegerIssue("closingSpeedScore", value.ClosingSpeedScore, &issues)
	return buildManualIntegerError("财报手工数字", issues)
}

func validateBaselineManualIntegers(value payload.BaselinePayload) error {
	issues := make([]ManualIntegerIssue, 0)
	typed := reflect.ValueOf(value)
	valueType := typed.Type()
	for index := 0; index < typed.NumField(); index++ {
		field := typed.Field(index)
		if field.Kind() != reflect.Float64 {
			continue
		}
		path := jsonFieldName(valueType.Field(index))
		collectFloatIntegerIssue(path, field.Float(), &issues)
	}
	return buildManualIntegerError("初始基线手工数字", issues)
}

func collectManualIntegerIssues(path string, value any, issues *[]ManualIntegerIssue) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			collectManualIntegerIssues(path+"."+key, child, issues)
		}
	case map[string]float64:
		for key, child := range typed {
			collectFloatIntegerIssue(path+"."+key, child, issues)
		}
	case map[string]int:
		return
	case payload.OperatingQuarterMap:
		for key, child := range typed {
			collectManualIntegerIssues(path+"."+key, child, issues)
		}
	case []map[string]any:
		for index, child := range typed {
			collectManualIntegerIssues(fmt.Sprintf("%s[%d]", path, index), child, issues)
		}
	case []any:
		for index, child := range typed {
			collectManualIntegerIssues(fmt.Sprintf("%s[%d]", path, index), child, issues)
		}
	case float64:
		collectFloatIntegerIssue(path, typed, issues)
	case float32:
		collectFloatIntegerIssue(path, float64(typed), issues)
	case json.Number:
		collectNumericStringIntegerIssue(path, typed.String(), issues)
	case string:
		collectNumericStringIntegerIssue(path, typed, issues)
	default:
		return
	}
}

func collectNullableFloatIntegerIssue(path string, value *float64, issues *[]ManualIntegerIssue) {
	if value == nil {
		return
	}
	collectFloatIntegerIssue(path, *value, issues)
}

func collectFloatIntegerIssue(path string, value float64, issues *[]ManualIntegerIssue) {
	if isWholeNumber(value) {
		return
	}
	*issues = append(*issues, ManualIntegerIssue{
		Path:  path,
		Value: formatValidationAmount(value),
	})
}

func collectNumericStringIntegerIssue(path string, raw string, issues *[]ManualIntegerIssue) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	if isWholeNumber(parsed) {
		return
	}
	*issues = append(*issues, ManualIntegerIssue{
		Path:  path,
		Value: value,
	})
}

func buildManualIntegerError(scope string, issues []ManualIntegerIssue) error {
	if len(issues) == 0 {
		return nil
	}
	return &ManualIntegerValidationError{
		Scope:  scope,
		Issues: issues,
	}
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return field.Name
	}
	if commaIndex := strings.Index(tag, ","); commaIndex >= 0 {
		tag = tag[:commaIndex]
	}
	if tag == "" {
		return field.Name
	}
	return tag
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}
