package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

const (
	dictionaryChangeTypeInitialize     = "INITIALIZE"
	dictionaryChangeTypeUpdateCurrent  = "UPDATE_CURRENT"
	dictionaryChangeTypeApplyScheme    = "APPLY_SCHEME"
	dictionaryChangeTypeRestoreDefault = "RESTORE_DEFAULT"
	dictionaryChangeTypeSaveScheme     = "SAVE_SCHEME"
	dictionaryChangeTypeDeleteScheme   = "DELETE_SCHEME"

	dictionaryDefaultSchemeName = "版本默认名称"
)

var (
	ErrDictionaryEditionRequired     = errors.New("dictionary edition required")
	ErrDictionaryEditionInvalid      = errors.New("dictionary edition invalid")
	ErrDictionarySchemeNameRequired  = errors.New("dictionary scheme name required")
	ErrDictionarySchemeInvalid       = errors.New("dictionary scheme invalid")
	ErrDictionarySchemeCrossEdition  = errors.New("dictionary scheme cross edition")
	ErrDictionaryDisplayNameRequired = errors.New("dictionary display name required")
	ErrDictionaryItemInvalid         = errors.New("dictionary item invalid")
	ErrDictionaryNotInitialized      = errors.New("dictionary game not initialized")
)

type DictionaryItemInput struct {
	ItemCode    string `json:"itemCode"`
	DisplayName string `json:"displayName"`
}

type DictionaryItemResult struct {
	ItemCode       string `json:"itemCode"`
	ItemCategory   string `json:"itemCategory"`
	DefaultName    string `json:"defaultName"`
	DisplayName    string `json:"displayName"`
	DisplayOrder   int    `json:"displayOrder"`
	Editable       bool   `json:"editable"`
	RelatedPayload any    `json:"relatedPayload,omitempty"`
}

type DictionaryCurrentResult struct {
	Initialized        bool                   `json:"initialized"`
	EditionCode        string                 `json:"editionCode"`
	EditionName        string                 `json:"editionName"`
	DictionaryRevision int                    `json:"dictionaryRevision"`
	CanUpdateCurrent   bool                   `json:"canUpdateCurrent"`
	Items              []DictionaryItemResult `json:"items"`
}

type DictionarySchemeSummary struct {
	ID          int64  `json:"id"`
	EditionCode string `json:"editionCode"`
	SchemeName  string `json:"schemeName"`
	Description string `json:"description"`
	BuiltIn     bool   `json:"builtIn"`
	ItemCount   int    `json:"itemCount"`
	UpdatedAt   string `json:"updatedAt"`
	UpdatedBy   string `json:"updatedBy"`
}

type DictionarySchemeListResult struct {
	EditionCode string                    `json:"editionCode"`
	List        []DictionarySchemeSummary `json:"list"`
}

type DictionarySchemeDetailResult struct {
	ID          int64                  `json:"id"`
	EditionCode string                 `json:"editionCode"`
	SchemeName  string                 `json:"schemeName"`
	Description string                 `json:"description"`
	BuiltIn     bool                   `json:"builtIn"`
	Items       []DictionaryItemResult `json:"items"`
	UpdatedAt   string                 `json:"updatedAt"`
	UpdatedBy   string                 `json:"updatedBy"`
}

type SaveDictionarySchemeCommand struct {
	SchemeID     *int64
	EditionCode  string
	SchemeName   string
	Description  string
	Items        []DictionaryItemInput
	OperatorID   int64
	OperatorName string
}

type DeleteDictionarySchemeCommand struct {
	SchemeID     int64
	OperatorID   int64
	OperatorName string
}

type UpdateCurrentDictionaryCommand struct {
	Items        []DictionaryItemInput
	Reason       string
	OperatorID   int64
	OperatorName string
}

type ApplyDictionarySchemeCommand struct {
	SchemeID     int64
	OperatorID   int64
	OperatorName string
}

type RestoreCurrentDictionaryDefaultCommand struct {
	OperatorID   int64
	OperatorName string
}

type DictionaryRevisionResult struct {
	EditionCode        string `json:"editionCode"`
	DictionaryRevision int    `json:"dictionaryRevision"`
}

type DictionaryChangeLogItem struct {
	ID             int64  `json:"id"`
	EditionCode    string `json:"editionCode"`
	ChangeType     string `json:"changeType"`
	SchemeID       *int64 `json:"schemeId"`
	Reason         string `json:"reason"`
	Revision       int    `json:"revision"`
	OperatorName   string `json:"operatorName"`
	OperateTime    string `json:"operateTime"`
	ChangedSummary string `json:"changedSummary"`
}

type DictionaryChangeLogPageResult struct {
	List     []DictionaryChangeLogItem `json:"list"`
	PageNo   int                       `json:"pageNo"`
	PageSize int                       `json:"pageSize"`
	Total    int64                     `json:"total"`
}

type AdminDictionaryService struct {
	db *gorm.DB
}

func NewAdminDictionaryService(db *gorm.DB) *AdminDictionaryService {
	return &AdminDictionaryService{db: db}
}

func (s *AdminDictionaryService) GetCurrent(ctx context.Context, editionCode string) (*DictionaryCurrentResult, error) {
	gameConfigRepo := repository.NewGameConfigRepository(s.db)
	groupRepo := repository.NewGroupRepository(s.db)
	dictionaryRepo := repository.NewDictionaryRepository(s.db)

	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	groupCount, err := groupRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("count groups: %w", err)
	}
	initialized := groupCount > 0

	edition := normalizeGameConfigEdition(*gameConfig)
	if !initialized {
		edition, err = validateDictionaryEditionOrDefault(editionCode)
		if err != nil {
			return nil, err
		}
	}

	items := []DictionaryItemResult{}
	if initialized {
		currentItems, loadErr := dictionaryRepo.ListCurrentItems(ctx)
		if loadErr != nil {
			return nil, fmt.Errorf("load current dictionary: %w", loadErr)
		}
		if len(currentItems) > 0 {
			items = mergeDictionaryItemsWithDefinitions(edition.EditionCode, buildDictionaryItemResultsFromCurrent(currentItems))
		}
	}
	if len(items) == 0 {
		items = buildDictionaryItemResultsFromDefinitions(BuiltInDictionaryDefinitions(edition.EditionCode), nil)
	}

	return &DictionaryCurrentResult{
		Initialized:        initialized,
		EditionCode:        edition.EditionCode,
		EditionName:        edition.EditionName,
		DictionaryRevision: gameConfig.DictionaryRevision,
		CanUpdateCurrent:   initialized,
		Items:              items,
	}, nil
}

func (s *AdminDictionaryService) ListSchemes(ctx context.Context, editionCode string) (*DictionarySchemeListResult, error) {
	edition, err := validateDictionaryEditionOrDefault(editionCode)
	if err != nil {
		return nil, err
	}
	repo := repository.NewDictionaryRepository(s.db)
	schemes, err := repo.ListSchemesByEdition(ctx, edition.EditionCode)
	if err != nil {
		return nil, fmt.Errorf("load dictionary schemes: %w", err)
	}
	list := make([]DictionarySchemeSummary, 0, len(schemes)+1)
	list = append(list, DictionarySchemeSummary{
		ID:          0,
		EditionCode: edition.EditionCode,
		SchemeName:  dictionaryDefaultSchemeName,
		Description: "系统内置默认显示名称，不需要保存即可使用。",
		BuiltIn:     true,
		ItemCount:   len(BuiltInDictionaryDefinitions(edition.EditionCode)),
		UpdatedAt:   "",
		UpdatedBy:   "system",
	})
	for _, scheme := range schemes {
		count, countErr := repo.CountSchemeItems(ctx, scheme.ID)
		if countErr != nil {
			return nil, fmt.Errorf("count dictionary scheme items: %w", countErr)
		}
		list = append(list, DictionarySchemeSummary{
			ID:          scheme.ID,
			EditionCode: scheme.EditionCode,
			SchemeName:  scheme.SchemeName,
			Description: scheme.Description,
			BuiltIn:     scheme.BuiltIn,
			ItemCount:   int(count),
			UpdatedAt:   scheme.UpdateTime.Format(time.RFC3339),
			UpdatedBy:   scheme.Updater,
		})
	}
	return &DictionarySchemeListResult{EditionCode: edition.EditionCode, List: list}, nil
}

func (s *AdminDictionaryService) GetSchemeDetail(ctx context.Context, schemeID int64) (*DictionarySchemeDetailResult, error) {
	repo := repository.NewDictionaryRepository(s.db)
	scheme, err := repo.GetSchemeByID(ctx, schemeID)
	if err != nil {
		return nil, fmt.Errorf("load dictionary scheme: %w", err)
	}
	items, err := repo.ListSchemeItems(ctx, scheme.ID)
	if err != nil {
		return nil, fmt.Errorf("load dictionary scheme items: %w", err)
	}
	return buildDictionarySchemeDetailResult(*scheme, mergeDictionaryItemsWithDefinitions(scheme.EditionCode, buildDictionaryItemResultsFromSchemeItems(items))), nil
}

func (s *AdminDictionaryService) SaveScheme(ctx context.Context, cmd SaveDictionarySchemeCommand) (*DictionarySchemeDetailResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	edition, err := validateDictionaryEditionOrDefault(cmd.EditionCode)
	if err != nil {
		return nil, err
	}
	schemeName := strings.TrimSpace(cmd.SchemeName)
	if schemeName == "" {
		return nil, ErrDictionarySchemeNameRequired
	}

	var result *DictionarySchemeDetailResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := repository.NewDictionaryRepository(tx)
		now := time.Now()
		items, normalizeErr := buildDictionaryItemsForEdition(edition.EditionCode, cmd.Items)
		if normalizeErr != nil {
			return normalizeErr
		}

		var scheme entity.DictionaryScheme
		if cmd.SchemeID != nil && *cmd.SchemeID > 0 {
			loaded, loadErr := repo.GetSchemeByID(ctx, *cmd.SchemeID)
			if loadErr != nil {
				return fmt.Errorf("load dictionary scheme: %w", loadErr)
			}
			if loaded.EditionCode != edition.EditionCode {
				return ErrDictionarySchemeCrossEdition
			}
			if loaded.BuiltIn {
				return ErrDictionarySchemeInvalid
			}
			if updateErr := repo.UpdateScheme(ctx, loaded.ID, schemeName, strings.TrimSpace(cmd.Description), operatorName, now); updateErr != nil {
				return fmt.Errorf("update dictionary scheme: %w", updateErr)
			}
			scheme = *loaded
			scheme.SchemeName = schemeName
			scheme.Description = strings.TrimSpace(cmd.Description)
			scheme.Updater = operatorName
			scheme.UpdateTime = now
		} else {
			scheme = entity.DictionaryScheme{
				EditionCode: edition.EditionCode,
				SchemeName:  schemeName,
				Description: strings.TrimSpace(cmd.Description),
				BuiltIn:     false,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			}
			if createErr := repo.CreateScheme(ctx, &scheme); createErr != nil {
				return fmt.Errorf("create dictionary scheme: %w", createErr)
			}
		}

		schemeItems := make([]entity.DictionarySchemeItem, 0, len(items))
		for _, item := range items {
			schemeItems = append(schemeItems, entity.DictionarySchemeItem{
				SchemeID:       scheme.ID,
				EditionCode:    edition.EditionCode,
				ItemCode:       item.ItemCode,
				ItemCategory:   item.ItemCategory,
				DefaultName:    item.DefaultName,
				DisplayName:    item.DisplayName,
				DisplayOrder:   item.DisplayOrder,
				Editable:       item.Editable,
				RelatedPayload: marshalAny(item.RelatedPayload),
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
		if replaceErr := repo.ReplaceSchemeItems(ctx, scheme.ID, schemeItems); replaceErr != nil {
			return fmt.Errorf("replace dictionary scheme items: %w", replaceErr)
		}
		if err := createDictionaryChangeLog(ctx, repo, edition.EditionCode, dictionaryChangeTypeSaveScheme, &scheme.ID, "", nil, items, 0, cmd.OperatorID, operatorName, now); err != nil {
			return err
		}
		result = buildDictionarySchemeDetailResult(scheme, mergeDictionaryItemsWithDefinitions(edition.EditionCode, buildDictionaryItemResultsFromSchemeItems(schemeItems)))
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AdminDictionaryService) DeleteScheme(ctx context.Context, cmd DeleteDictionarySchemeCommand) error {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := repository.NewDictionaryRepository(tx)
		scheme, err := repo.GetSchemeByID(ctx, cmd.SchemeID)
		if err != nil {
			return fmt.Errorf("load dictionary scheme: %w", err)
		}
		if scheme.BuiltIn {
			return ErrDictionarySchemeInvalid
		}
		items, err := repo.ListSchemeItems(ctx, scheme.ID)
		if err != nil {
			return fmt.Errorf("load dictionary scheme items: %w", err)
		}
		if err := repo.ReplaceSchemeItems(ctx, scheme.ID, nil); err != nil {
			return fmt.Errorf("delete dictionary scheme items: %w", err)
		}
		if err := repo.DeleteScheme(ctx, scheme.ID); err != nil {
			return fmt.Errorf("delete dictionary scheme: %w", err)
		}
		return createDictionaryChangeLog(
			ctx,
			repo,
			scheme.EditionCode,
			dictionaryChangeTypeDeleteScheme,
			&scheme.ID,
			"",
			buildDictionaryItemResultsFromSchemeItems(items),
			nil,
			0,
			cmd.OperatorID,
			operatorName,
			time.Now(),
		)
	})
}

func (s *AdminDictionaryService) UpdateCurrent(ctx context.Context, cmd UpdateCurrentDictionaryCommand) (*DictionaryCurrentResult, error) {
	return s.replaceCurrent(ctx, replaceCurrentDictionaryCommand{
		Items:        cmd.Items,
		Reason:       cmd.Reason,
		ChangeType:   dictionaryChangeTypeUpdateCurrent,
		OperatorID:   cmd.OperatorID,
		OperatorName: cmd.OperatorName,
	})
}

func (s *AdminDictionaryService) ApplySchemeToCurrent(ctx context.Context, cmd ApplyDictionarySchemeCommand) (*DictionaryCurrentResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	var result *DictionaryCurrentResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		repo := repository.NewDictionaryRepository(tx)
		gameConfig, edition, err := loadInitializedDictionaryGameForUpdate(ctx, gameConfigRepo, groupRepo)
		if err != nil {
			return err
		}
		scheme, err := repo.GetSchemeByID(ctx, cmd.SchemeID)
		if err != nil {
			return fmt.Errorf("load dictionary scheme: %w", err)
		}
		if scheme.EditionCode != edition.EditionCode {
			return ErrDictionarySchemeCrossEdition
		}
		schemeItems, err := repo.ListSchemeItems(ctx, scheme.ID)
		if err != nil {
			return fmt.Errorf("load dictionary scheme items: %w", err)
		}
		items := mergeDictionaryItemsWithDefinitions(edition.EditionCode, buildDictionaryItemResultsFromSchemeItems(schemeItems))
		currentBefore, err := loadCurrentDictionaryOrDefault(ctx, repo, edition.EditionCode)
		if err != nil {
			return err
		}
		now := time.Now()
		revision := gameConfig.DictionaryRevision + 1
		currentItems := currentDictionaryEntitiesFromResults(edition.EditionCode, items, operatorName, now)
		if err := repo.ReplaceCurrentItems(ctx, currentItems); err != nil {
			return fmt.Errorf("replace current dictionary: %w", err)
		}
		if err := gameConfigRepo.UpdateDictionaryRevision(ctx, gameConfig.ID, revision, operatorName, now); err != nil {
			return fmt.Errorf("update dictionary revision: %w", err)
		}
		if err := createDictionaryChangeLog(ctx, repo, edition.EditionCode, dictionaryChangeTypeApplyScheme, &scheme.ID, scheme.SchemeName, currentBefore, items, revision, cmd.OperatorID, operatorName, now); err != nil {
			return err
		}
		result = &DictionaryCurrentResult{
			Initialized:        true,
			EditionCode:        edition.EditionCode,
			EditionName:        edition.EditionName,
			DictionaryRevision: revision,
			CanUpdateCurrent:   true,
			Items:              items,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AdminDictionaryService) RestoreCurrentDefault(ctx context.Context, cmd RestoreCurrentDictionaryDefaultCommand) (*DictionaryCurrentResult, error) {
	return s.replaceCurrent(ctx, replaceCurrentDictionaryCommand{
		Items:        nil,
		Reason:       "恢复版本默认显示名称",
		ChangeType:   dictionaryChangeTypeRestoreDefault,
		OperatorID:   cmd.OperatorID,
		OperatorName: cmd.OperatorName,
	})
}

func (s *AdminDictionaryService) GetRevision(ctx context.Context) (*DictionaryRevisionResult, error) {
	repo := repository.NewGameConfigRepository(s.db)
	gameConfig, err := repo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	edition := normalizeGameConfigEdition(*gameConfig)
	return &DictionaryRevisionResult{
		EditionCode:        edition.EditionCode,
		DictionaryRevision: gameConfig.DictionaryRevision,
	}, nil
}

func (s *AdminDictionaryService) PageChangeLogs(ctx context.Context, pageNo int, pageSize int) (*DictionaryChangeLogPageResult, error) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	repo := repository.NewDictionaryRepository(s.db)
	items, total, err := repo.PageChangeLogs(ctx, pageNo, pageSize)
	if err != nil {
		return nil, fmt.Errorf("load dictionary logs: %w", err)
	}
	result := make([]DictionaryChangeLogItem, 0, len(items))
	for _, item := range items {
		result = append(result, DictionaryChangeLogItem{
			ID:             item.ID,
			EditionCode:    item.EditionCode,
			ChangeType:     item.ChangeType,
			SchemeID:       item.SchemeID,
			Reason:         item.Reason,
			Revision:       item.Revision,
			OperatorName:   item.OperatorName,
			OperateTime:    item.OperateTime.Format(time.RFC3339),
			ChangedSummary: buildDictionaryChangeSummary(item.BeforeJSON, item.AfterJSON),
		})
	}
	return &DictionaryChangeLogPageResult{List: result, PageNo: pageNo, PageSize: pageSize, Total: total}, nil
}

type replaceCurrentDictionaryCommand struct {
	Items        []DictionaryItemInput
	Reason       string
	ChangeType   string
	OperatorID   int64
	OperatorName string
}

func (s *AdminDictionaryService) replaceCurrent(ctx context.Context, cmd replaceCurrentDictionaryCommand) (*DictionaryCurrentResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	var result *DictionaryCurrentResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		repo := repository.NewDictionaryRepository(tx)
		gameConfig, edition, err := loadInitializedDictionaryGameForUpdate(ctx, gameConfigRepo, groupRepo)
		if err != nil {
			return err
		}
		currentBefore, err := loadCurrentDictionaryOrDefault(ctx, repo, edition.EditionCode)
		if err != nil {
			return err
		}
		items, err := buildDictionaryItemsForEdition(edition.EditionCode, cmd.Items)
		if err != nil {
			return err
		}
		now := time.Now()
		revision := gameConfig.DictionaryRevision + 1
		if err := repo.ReplaceCurrentItems(ctx, currentDictionaryEntitiesFromResults(edition.EditionCode, items, operatorName, now)); err != nil {
			return fmt.Errorf("replace current dictionary: %w", err)
		}
		if err := gameConfigRepo.UpdateDictionaryRevision(ctx, gameConfig.ID, revision, operatorName, now); err != nil {
			return fmt.Errorf("update dictionary revision: %w", err)
		}
		if err := createDictionaryChangeLog(ctx, repo, edition.EditionCode, cmd.ChangeType, nil, strings.TrimSpace(cmd.Reason), currentBefore, items, revision, cmd.OperatorID, operatorName, now); err != nil {
			return err
		}
		result = &DictionaryCurrentResult{
			Initialized:        true,
			EditionCode:        edition.EditionCode,
			EditionName:        edition.EditionName,
			DictionaryRevision: revision,
			CanUpdateCurrent:   true,
			Items:              items,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func loadInitializedDictionaryGameForUpdate(ctx context.Context, gameConfigRepo *repository.GameConfigRepository, groupRepo *repository.GroupRepository) (*entity.GameConfig, GameEdition, error) {
	gameConfig, err := gameConfigRepo.GetCurrentForUpdate(ctx)
	if err != nil {
		return nil, GameEdition{}, fmt.Errorf("load game config: %w", err)
	}
	groupCount, err := groupRepo.CountAll(ctx)
	if err != nil {
		return nil, GameEdition{}, fmt.Errorf("count groups: %w", err)
	}
	if groupCount <= 0 {
		return nil, GameEdition{}, ErrDictionaryNotInitialized
	}
	edition := normalizeGameConfigEdition(*gameConfig)
	return gameConfig, edition, nil
}

func validateDictionaryEditionOrDefault(editionCode string) (GameEdition, error) {
	if strings.TrimSpace(editionCode) == "" {
		return DefaultGameEdition(), nil
	}
	edition, ok := FindGameEdition(editionCode)
	if !ok {
		return GameEdition{}, ErrDictionaryEditionInvalid
	}
	return edition, nil
}

func buildDictionaryItemsForEdition(editionCode string, inputs []DictionaryItemInput) ([]DictionaryItemResult, error) {
	definitions := BuiltInDictionaryDefinitions(editionCode)
	overrides := make(map[string]string, len(inputs))
	definitionMap := make(map[string]DictionaryDefinition, len(definitions))
	for _, definition := range definitions {
		definitionMap[definition.ItemCode] = definition
	}
	for _, input := range inputs {
		itemCode := strings.TrimSpace(input.ItemCode)
		if itemCode == "" {
			return nil, ErrDictionaryItemInvalid
		}
		if _, ok := definitionMap[itemCode]; !ok {
			return nil, ErrDictionaryItemInvalid
		}
		displayName := strings.TrimSpace(input.DisplayName)
		if displayName == "" {
			return nil, ErrDictionaryDisplayNameRequired
		}
		overrides[itemCode] = displayName
	}
	return buildDictionaryItemResultsFromDefinitions(definitions, overrides), nil
}

func buildDictionaryItemResultsFromDefinitions(definitions []DictionaryDefinition, overrides map[string]string) []DictionaryItemResult {
	items := make([]DictionaryItemResult, 0, len(definitions))
	for _, definition := range definitions {
		displayName := definition.DefaultName
		if value, ok := overrides[definition.ItemCode]; ok {
			displayName = value
		}
		items = append(items, DictionaryItemResult{
			ItemCode:       definition.ItemCode,
			ItemCategory:   definition.ItemCategory,
			DefaultName:    definition.DefaultName,
			DisplayName:    displayName,
			DisplayOrder:   definition.DisplayOrder,
			Editable:       definition.Editable,
			RelatedPayload: definition.RelatedPayload,
		})
	}
	return items
}

func buildDictionaryItemResultsFromCurrent(items []entity.CurrentDictionaryItem) []DictionaryItemResult {
	result := make([]DictionaryItemResult, 0, len(items))
	for _, item := range items {
		result = append(result, DictionaryItemResult{
			ItemCode:       item.ItemCode,
			ItemCategory:   item.ItemCategory,
			DefaultName:    item.DefaultName,
			DisplayName:    item.DisplayName,
			DisplayOrder:   item.DisplayOrder,
			Editable:       item.Editable,
			RelatedPayload: unmarshalAny(item.RelatedPayload),
		})
	}
	return result
}

func buildDictionaryItemResultsFromSchemeItems(items []entity.DictionarySchemeItem) []DictionaryItemResult {
	result := make([]DictionaryItemResult, 0, len(items))
	for _, item := range items {
		result = append(result, DictionaryItemResult{
			ItemCode:       item.ItemCode,
			ItemCategory:   item.ItemCategory,
			DefaultName:    item.DefaultName,
			DisplayName:    item.DisplayName,
			DisplayOrder:   item.DisplayOrder,
			Editable:       item.Editable,
			RelatedPayload: unmarshalAny(item.RelatedPayload),
		})
	}
	return result
}

func currentDictionaryEntitiesFromResults(editionCode string, items []DictionaryItemResult, operatorName string, now time.Time) []entity.CurrentDictionaryItem {
	result := make([]entity.CurrentDictionaryItem, 0, len(items))
	for _, item := range items {
		result = append(result, entity.CurrentDictionaryItem{
			EditionCode:    editionCode,
			ItemCode:       item.ItemCode,
			ItemCategory:   item.ItemCategory,
			DefaultName:    item.DefaultName,
			DisplayName:    item.DisplayName,
			DisplayOrder:   item.DisplayOrder,
			Editable:       item.Editable,
			RelatedPayload: marshalAny(item.RelatedPayload),
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: now,
				Updater:    operatorName,
				UpdateTime: now,
			},
		})
	}
	return result
}

func loadCurrentDictionaryOrDefault(ctx context.Context, repo *repository.DictionaryRepository, editionCode string) ([]DictionaryItemResult, error) {
	currentItems, err := repo.ListCurrentItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("load current dictionary: %w", err)
	}
	if len(currentItems) > 0 {
		return mergeDictionaryItemsWithDefinitions(editionCode, buildDictionaryItemResultsFromCurrent(currentItems)), nil
	}
	return buildDictionaryItemResultsFromDefinitions(BuiltInDictionaryDefinitions(editionCode), nil), nil
}

func mergeDictionaryItemsWithDefinitions(editionCode string, items []DictionaryItemResult) []DictionaryItemResult {
	definitions := BuiltInDictionaryDefinitions(editionCode)
	itemMap := make(map[string]DictionaryItemResult, len(items))
	for _, item := range items {
		itemMap[item.ItemCode] = item
	}
	result := make([]DictionaryItemResult, 0, len(definitions))
	for _, definition := range definitions {
		if item, ok := itemMap[definition.ItemCode]; ok {
			item.ItemCategory = definition.ItemCategory
			item.DefaultName = definition.DefaultName
			item.DisplayOrder = definition.DisplayOrder
			item.Editable = definition.Editable
			item.RelatedPayload = definition.RelatedPayload
			result = append(result, item)
			continue
		}
		result = append(result, DictionaryItemResult{
			ItemCode:       definition.ItemCode,
			ItemCategory:   definition.ItemCategory,
			DefaultName:    definition.DefaultName,
			DisplayName:    definition.DefaultName,
			DisplayOrder:   definition.DisplayOrder,
			Editable:       definition.Editable,
			RelatedPayload: definition.RelatedPayload,
		})
	}
	return result
}

func buildDictionarySchemeDetailResult(scheme entity.DictionaryScheme, items []DictionaryItemResult) *DictionarySchemeDetailResult {
	return &DictionarySchemeDetailResult{
		ID:          scheme.ID,
		EditionCode: scheme.EditionCode,
		SchemeName:  scheme.SchemeName,
		Description: scheme.Description,
		BuiltIn:     scheme.BuiltIn,
		Items:       items,
		UpdatedAt:   scheme.UpdateTime.Format(time.RFC3339),
		UpdatedBy:   scheme.Updater,
	}
}

func createDictionaryChangeLog(
	ctx context.Context,
	repo *repository.DictionaryRepository,
	editionCode string,
	changeType string,
	schemeID *int64,
	reason string,
	before any,
	after any,
	revision int,
	operatorID int64,
	operatorName string,
	operateTime time.Time,
) error {
	beforeJSON := marshalAny(before)
	if len(beforeJSON) == 0 {
		beforeJSON = []byte("[]")
	}
	afterJSON := marshalAny(after)
	if len(afterJSON) == 0 {
		afterJSON = []byte("[]")
	}
	item := &entity.DictionaryChangeLog{
		EditionCode:  editionCode,
		ChangeType:   changeType,
		SchemeID:     schemeID,
		Reason:       strings.TrimSpace(reason),
		BeforeJSON:   beforeJSON,
		AfterJSON:    afterJSON,
		Revision:     revision,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		OperateTime:  operateTime,
	}
	if err := repo.CreateChangeLog(ctx, item); err != nil {
		return fmt.Errorf("create dictionary change log: %w", err)
	}
	return nil
}

func marshalAny(value any) []byte {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return raw
}

func unmarshalAny(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}
	return value
}

func buildDictionaryChangeSummary(beforeJSON []byte, afterJSON []byte) string {
	beforeItems := decodeDictionaryLogItems(beforeJSON)
	afterItems := decodeDictionaryLogItems(afterJSON)
	if len(afterItems) == 0 {
		return "删除或清空字典方案"
	}
	beforeMap := make(map[string]string, len(beforeItems))
	for _, item := range beforeItems {
		beforeMap[item.ItemCode] = item.DisplayName
	}
	changed := 0
	for _, item := range afterItems {
		if beforeMap[item.ItemCode] != item.DisplayName {
			changed++
		}
	}
	if changed == 0 {
		return "显示名称无变化"
	}
	return fmt.Sprintf("修改 %d 个显示名称", changed)
}

func decodeDictionaryLogItems(raw []byte) []DictionaryItemResult {
	var items []DictionaryItemResult
	if len(raw) == 0 {
		return items
	}
	_ = json.Unmarshal(raw, &items)
	return items
}
