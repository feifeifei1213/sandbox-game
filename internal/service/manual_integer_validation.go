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
var ErrSupplyChainOrderQuantityInvalid = fmt.Errorf("supply chain order quantity invalid")

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
	collectManualIntegerIssues("quarter.supplyChainOrderRecord", normalized.Quarter.SupplyChainOrderRecord, &issues)
	collectManualIntegerIssues("quarter.receivableUpdate", normalized.Quarter.ReceivableUpdate, &issues)
	collectManualIntegerIssues("quarter.deliverySettlement", normalized.Quarter.DeliverySettlement, &issues)
	collectManualIntegerIssues("yearEnd.longTermLoan", normalized.YearEnd.LongTermLoan, &issues)
	collectManualIntegerIssues("yearEnd.assetAdjustment", normalized.YearEnd.AssetAdjustment, &issues)
	collectManualIntegerIssues("yearEnd.projectProgressUpdate", normalized.YearEnd.ProjectProgressUpdate.Items, &issues)
	collectManualIntegerIssues("yearEnd.marketCultivation", normalized.YearEnd.MarketCultivation, &issues)
	collectManualIntegerIssues("extra.incomeAndPenalty", normalized.Extra.IncomeAndPenalty, &issues)

	if err := buildManualIntegerError("经营页手工数字", issues); err != nil {
		return err
	}
	if err := validateSupplyChainOrderQuantities(normalized.Quarter.SupplyChainOrderRecord); err != nil {
		return err
	}
	return validateOperatingFeatureInputs(normalized)
}

func validateSupplyChainOrderQuantities(value payload.OperatingQuarterMap) error {
	issues := make([]ManualIntegerIssue, 0)
	for quarterKey, quarterValue := range value {
		for _, fieldKey := range []string{"basicProduct", "standardProduct", "precisionProduct", "intelligentProduct"} {
			raw, exists := quarterValue[fieldKey]
			if !exists || raw == nil || (reflect.ValueOf(raw).Kind() == reflect.String && strings.TrimSpace(fmt.Sprint(raw)) == "") {
				continue
			}
			parsed, ok := manualNumericValue(raw)
			if ok && parsed >= 0 && isWholeNumber(parsed) {
				continue
			}
			issues = append(issues, ManualIntegerIssue{
				Path:  fmt.Sprintf("quarter.supplyChainOrderRecord.%s.%s", quarterKey, fieldKey),
				Value: strings.TrimSpace(fmt.Sprint(raw)),
			})
		}
	}
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, minInt(len(issues), manualIntegerPreviewLimit))
	for index, issue := range issues {
		if index >= manualIntegerPreviewLimit {
			break
		}
		parts = append(parts, fmt.Sprintf("%s=%s", issue.Path, issue.Value))
	}
	return fmt.Errorf("%w: 订单数量必须为非负整数，发现：%s", ErrSupplyChainOrderQuantityInvalid, strings.Join(parts, "、"))
}

func manualNumericValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	case json.Number:
		parsed, err := strconv.ParseFloat(typed.String(), 64)
		return parsed, err == nil
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed, err == nil
	default:
		return 0, false
	}
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
	case payload.OperatingMarketCultivationPayload:
		collectManualIntegerIssues(path+".regional", typed.Regional, issues)
		collectManualIntegerIssues(path+".national", typed.National, issues)
		collectManualIntegerIssues(path+".global", typed.Global, issues)
	case payload.OperatingMarketCultivationItem:
		collectManualIntegerIssues(path+".annualInvestment", typed.AnnualInvestment, issues)
	case payload.OperatingQualificationCertification:
		return
	case payload.OperatingQualificationItem:
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
