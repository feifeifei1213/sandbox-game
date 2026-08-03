package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"

	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

func applyOperatingFeatureState(
	ctx context.Context,
	operatingRepo *repository.OperatingRepository,
	groupID int64,
	yearNo int,
	source payload.OperatingPayload,
	previousYearEffective bool,
	currentYearEndEffective bool,
	preserveLockedAnnual bool,
) (payload.OperatingPayload, error) {
	normalized := source.Normalize()
	previousMarketState := payload.MarketCultivationCarryState{}
	previousQualificationState := payload.QualificationCarryState{}

	if previousYearEffective && yearNo > 0 && operatingRepo != nil {
		previousPayload, err := loadPreviousEffectiveOperatingPayload(ctx, operatingRepo, groupID, yearNo)
		if err != nil {
			return payload.OperatingPayload{}, err
		}
		if previousPayload != nil {
			previousMarketState = payload.ExtractMarketCultivationCarryState(*previousPayload)
			previousQualificationState = payload.ExtractQualificationCarryState(*previousPayload)
		}
	}

	normalized.YearEnd.MarketCultivation = payload.ApplyMarketCultivationState(
		normalized.YearEnd.MarketCultivation,
		previousMarketState,
		currentYearEndEffective,
		preserveLockedAnnual,
	)
	normalized.YearEnd.QualificationCertification = payload.ApplyQualificationState(
		normalized.YearEnd.QualificationCertification,
		previousQualificationState,
		currentYearEndEffective,
	)
	return normalized, nil
}

func loadPreviousEffectiveOperatingPayload(
	ctx context.Context,
	operatingRepo *repository.OperatingRepository,
	groupID int64,
	yearNo int,
) (*payload.OperatingPayload, error) {
	draft, err := operatingRepo.FindDraft(ctx, groupID, yearNo-1)
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, nil
	default:
		return nil, fmt.Errorf("load previous operating payload: %w", err)
	}
	if draft == nil || len(draft.OperatingPayload) == 0 {
		return nil, nil
	}

	previous := payload.NewOperatingPayload()
	if err := json.Unmarshal(draft.OperatingPayload, &previous); err != nil {
		return nil, fmt.Errorf("unmarshal previous operating payload: %w", err)
	}
	previous = previous.Normalize()
	return &previous, nil
}

func validateOperatingFeatureInputs(value payload.OperatingPayload) error {
	normalized := value.Normalize()
	if err := validateProjectProgressInputs(normalized.YearEnd.ProjectProgressUpdate); err != nil {
		return err
	}
	if err := validateMarketCultivationInputs(normalized.YearEnd.MarketCultivation); err != nil {
		return err
	}
	if err := validateQualificationCertificationInputs(normalized.YearEnd.QualificationCertification); err != nil {
		return err
	}
	return nil
}

func validateProjectProgressInputs(value payload.OperatingProjectProgressPayload) error {
	lineTypeOptions := map[string]bool{
		"":    true,
		"人工":  true,
		"半自动": true,
		"自动":  true,
		"智能":  true,
	}
	issues := make([]string, 0)
	for _, item := range value.Items {
		projectName := strings.TrimSpace(fmt.Sprint(item["projectName"]))
		if projectName == "" {
			projectName = "项目进度"
		}

		lineType := ""
		if raw, exists := item["lineType"]; exists && raw != nil {
			lineType = strings.TrimSpace(fmt.Sprint(raw))
		}
		if !lineTypeOptions[lineType] {
			issues = append(issues, fmt.Sprintf("%s产线=%s", projectName, lineType))
		}

		rawProgress, exists := item["progress"]
		if !exists || rawProgress == nil || strings.TrimSpace(fmt.Sprint(rawProgress)) == "" {
			continue
		}
		progress, ok := payload.ParseManualNumber(rawProgress)
		if !ok || !isWholeNumber(progress) || progress < 0 || progress > 4 {
			issues = append(issues, fmt.Sprintf("%s进度=%s", projectName, strings.TrimSpace(fmt.Sprint(rawProgress))))
		}
	}

	if len(issues) == 0 {
		return nil
	}
	return fmt.Errorf("项目进度更新的产线只能填写“人工 / 半自动 / 自动 / 智能”，进度只能填写 0~4：%s", strings.Join(issues, "、"))
}

func validateMarketCultivationInputs(value payload.OperatingMarketCultivationPayload) error {
	type marketItem struct {
		label string
		item  payload.OperatingMarketCultivationItem
	}

	issues := make([]string, 0)
	for _, current := range []marketItem{
		{label: "区域", item: value.Regional},
		{label: "全国", item: value.National},
		{label: "全球", item: value.Global},
	} {
		raw := current.item.AnnualInvestment
		if raw == nil || strings.TrimSpace(fmt.Sprint(raw)) == "" {
			continue
		}
		number, ok := payload.ParseManualNumber(raw)
		if !ok || !isWholeNumber(number) || (number != 0 && number != 1) {
			issues = append(issues, fmt.Sprintf("%s=%s", current.label, strings.TrimSpace(fmt.Sprint(raw))))
			continue
		}
		if current.item.LockedByPrevious && math.Abs(number) > 0.000001 {
			issues = append(issues, fmt.Sprintf("%s已解锁，不允许继续投入%sM", current.label, formatValidationAmount(number)))
		}
	}

	if len(issues) == 0 {
		return nil
	}
	return fmt.Errorf("新市场培育每个市场每年只能填写 0 或 1，且已解锁市场不能继续投入：%s", strings.Join(issues, "、"))
}

func validateQualificationCertificationInputs(value payload.OperatingQualificationCertification) error {
	type qualificationItem struct {
		label string
		item  payload.OperatingQualificationItem
	}

	issues := make([]string, 0)
	for _, current := range []qualificationItem{
		{label: "质量、环境健康体系认证企业", item: value.QualityEnvironmentalHealth},
		{label: "高新技术企业", item: value.HighTechEnterprise},
		{label: "专精特新小巨人", item: value.SpecializedInnovation},
		{label: "上市企业", item: value.ListedCompany},
	} {
		status := strings.TrimSpace(fmt.Sprint(current.item.Status))
		if status == "" {
			continue
		}
		if status != "未解锁" && status != "解锁" {
			issues = append(issues, fmt.Sprintf("%s=%s", current.label, status))
		}
		if current.item.LockedByPrevious && status != "" && status != "解锁" {
			issues = append(issues, fmt.Sprintf("%s已继承解锁，不允许改回%s", current.label, status))
		}
	}

	if len(issues) == 0 {
		return nil
	}
	return fmt.Errorf("资质认证状态只能填写“未解锁 / 解锁”：%s", strings.Join(issues, "、"))
}
