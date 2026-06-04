package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type DictionaryRepository struct {
	db *gorm.DB
}

func NewDictionaryRepository(db *gorm.DB) *DictionaryRepository {
	return &DictionaryRepository{db: db}
}

func (r *DictionaryRepository) ListCurrentItems(ctx context.Context) ([]entity.CurrentDictionaryItem, error) {
	var items []entity.CurrentDictionaryItem
	if err := r.db.WithContext(ctx).
		Order("display_order ASC").
		Order("item_code ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DictionaryRepository) ReplaceCurrentItems(ctx context.Context, items []entity.CurrentDictionaryItem) error {
	if err := r.db.WithContext(ctx).Where("1 = 1").Delete(&entity.CurrentDictionaryItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 200).Error
}

func (r *DictionaryRepository) ListSchemesByEdition(ctx context.Context, editionCode string) ([]entity.DictionaryScheme, error) {
	var items []entity.DictionaryScheme
	if err := r.db.WithContext(ctx).
		Where("edition_code = ?", editionCode).
		Order("built_in DESC").
		Order("update_time DESC").
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DictionaryRepository) GetSchemeByID(ctx context.Context, id int64) (*entity.DictionaryScheme, error) {
	var item entity.DictionaryScheme
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DictionaryRepository) CreateScheme(ctx context.Context, item *entity.DictionaryScheme) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DictionaryRepository) UpdateScheme(ctx context.Context, id int64, schemeName string, description string, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.DictionaryScheme{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"scheme_name": schemeName,
			"description": description,
			"updater":     operatorName,
			"update_time": operateTime,
		}).Error
}

func (r *DictionaryRepository) DeleteScheme(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.DictionaryScheme{}).Error
}

func (r *DictionaryRepository) CountSchemeItems(ctx context.Context, schemeID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.DictionarySchemeItem{}).
		Where("scheme_id = ?", schemeID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DictionaryRepository) ListSchemeItems(ctx context.Context, schemeID int64) ([]entity.DictionarySchemeItem, error) {
	var items []entity.DictionarySchemeItem
	if err := r.db.WithContext(ctx).
		Where("scheme_id = ?", schemeID).
		Order("display_order ASC").
		Order("item_code ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DictionaryRepository) ReplaceSchemeItems(ctx context.Context, schemeID int64, items []entity.DictionarySchemeItem) error {
	if err := r.db.WithContext(ctx).Where("scheme_id = ?", schemeID).Delete(&entity.DictionarySchemeItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 200).Error
}

func (r *DictionaryRepository) CreateChangeLog(ctx context.Context, item *entity.DictionaryChangeLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DictionaryRepository) PageChangeLogs(ctx context.Context, pageNo int, pageSize int) ([]entity.DictionaryChangeLog, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.DictionaryChangeLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.DictionaryChangeLog
	offset := (pageNo - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Order("operate_time DESC").
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
