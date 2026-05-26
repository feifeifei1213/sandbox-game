package service

import (
	"strings"

	"sandbox-game/internal/model/entity"
)

const (
	GameEditionVIPServiceV1              = "VIP_SERVICE_V1"
	GameEditionVIPServiceV1Name          = "贵宾服务版 V1"
	FormulaVersionCommonV1               = "COMMON_FORMULA_V1"
	ProcessRuleVersionCommonV1           = "COMMON_PROCESS_V1"
	TemplateVersionVIPServiceV1          = "VIP_SERVICE_V1"
	OperatingTemplateVersionVIPServiceV1 = "VIP_OPERATING_TEMPLATE_V1"
	ReportTemplateVersionVIPServiceV1    = "VIP_REPORT_TEMPLATE_V1"
	OrderTemplateVersionVIPServiceV1     = "VIP_ORDER_TEMPLATE_V1"
)

type GameEdition struct {
	EditionCode              string `json:"editionCode"`
	EditionName              string `json:"editionName"`
	Description              string `json:"description"`
	DefaultEdition           bool   `json:"defaultEdition"`
	RuleVersion              string `json:"ruleVersion"`
	FormulaVersion           string `json:"formulaVersion"`
	TemplateVersion          string `json:"templateVersion"`
	OperatingTemplateVersion string `json:"operatingTemplateVersion"`
	ReportTemplateVersion    string `json:"reportTemplateVersion"`
	OrderTemplateVersion     string `json:"orderTemplateVersion"`
	ProcessRuleVersion       string `json:"processRuleVersion"`
}

var builtInGameEditions = []GameEdition{
	{
		EditionCode:              GameEditionVIPServiceV1,
		EditionName:              GameEditionVIPServiceV1Name,
		Description:              "贵宾服务版字段模板，公式规则和流程规则沿用当前通用版本。",
		DefaultEdition:           true,
		RuleVersion:              FormulaVersionCommonV1,
		FormulaVersion:           FormulaVersionCommonV1,
		TemplateVersion:          TemplateVersionVIPServiceV1,
		OperatingTemplateVersion: OperatingTemplateVersionVIPServiceV1,
		ReportTemplateVersion:    ReportTemplateVersionVIPServiceV1,
		OrderTemplateVersion:     OrderTemplateVersionVIPServiceV1,
		ProcessRuleVersion:       ProcessRuleVersionCommonV1,
	},
}

func ListGameEditions() []GameEdition {
	result := make([]GameEdition, len(builtInGameEditions))
	copy(result, builtInGameEditions)
	return result
}

func DefaultGameEdition() GameEdition {
	for _, edition := range builtInGameEditions {
		if edition.DefaultEdition {
			return edition
		}
	}
	return builtInGameEditions[0]
}

func FindGameEdition(code string) (GameEdition, bool) {
	normalized := normalizeEditionCode(code)
	for _, edition := range builtInGameEditions {
		if edition.EditionCode == normalized {
			return edition, true
		}
	}
	return GameEdition{}, false
}

func ResolveGameEdition(code string) GameEdition {
	if edition, ok := FindGameEdition(code); ok {
		return edition
	}
	return DefaultGameEdition()
}

func normalizeGameConfigEdition(config entity.GameConfig) GameEdition {
	edition := ResolveGameEdition(config.EditionCode)
	if strings.TrimSpace(config.EditionCode) != "" {
		edition.EditionCode = strings.TrimSpace(config.EditionCode)
	}
	if strings.TrimSpace(config.EditionName) != "" {
		edition.EditionName = strings.TrimSpace(config.EditionName)
	}
	if strings.TrimSpace(config.RuleVersion) != "" {
		edition.RuleVersion = strings.TrimSpace(config.RuleVersion)
		edition.FormulaVersion = strings.TrimSpace(config.RuleVersion)
	}
	if strings.TrimSpace(config.TemplateVersion) != "" {
		edition.TemplateVersion = strings.TrimSpace(config.TemplateVersion)
	}
	if strings.TrimSpace(config.OperatingTemplateVersion) != "" {
		edition.OperatingTemplateVersion = strings.TrimSpace(config.OperatingTemplateVersion)
	}
	if strings.TrimSpace(config.ReportTemplateVersion) != "" {
		edition.ReportTemplateVersion = strings.TrimSpace(config.ReportTemplateVersion)
	}
	if strings.TrimSpace(config.OrderTemplateVersion) != "" {
		edition.OrderTemplateVersion = strings.TrimSpace(config.OrderTemplateVersion)
	}
	if strings.TrimSpace(config.ProcessRuleVersion) != "" {
		edition.ProcessRuleVersion = strings.TrimSpace(config.ProcessRuleVersion)
	}
	return edition
}

func normalizeEditionCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
