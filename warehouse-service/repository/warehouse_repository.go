package repository

import (
	"context"
	"errors"
	"micro-warehouse/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseRepositoryInterface interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetAllWarehouse(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Warehouse, int64, error)
	GetWarehouseById(ctx context.Context, id uint) (*model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id uint) error
}

type warehouseRepository struct {
	db *gorm.DB
}

// CreateWarehouse implements [WarehouseRepositoryInterface].
func (w *warehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] CreateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouse).Error
	}
}

// DeleteWarehouse implements [WarehouseRepositoryInterface].
func (w *warehouseRepository) DeleteWarehouse(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] DeleteWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouse := model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 2: %v", err)
			return err
		}

		if len(modelWarehouse.WarehouseProducts) > 0 {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 3: %v", errors.New("warehouse has products"))
			return errors.New("warehouse has products")
		}

		return w.db.WithContext(ctx).Delete(&modelWarehouse).Error
	}
}

// GetAllWarehouse implements [WarehouseRepositoryInterface].
func (w *warehouseRepository) GetAllWarehouse(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetAllWarehouse - 1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		if page <= 0 {
			page = 1
		}

		if limit <= 0 {
			limit = 10
		}

		if sortBy == "" {
			sortBy = "created_at"
		}

		if sortOrder == "" {
			sortOrder = "desc"
		}

		offset := (page - 1) * limit

		// build query
		query := w.db.Model(&model.Warehouse{})

		if search != "" {
			query = query.Where("name ILIKE ? OR address ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		// Get total count
		var warehouse []model.Warehouse
		var total int64
		if err := query.Count(&total).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouse - 2: %v", err)
			return nil, 0, err
		}

		
		if err := query.Order(sortBy + " " + sortOrder).
			WithContext(ctx).
			Preload("WarehouseProducts").
			Offset(offset).
			Limit(limit).
			Find(&warehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouse - 3: %v", err)
			return nil, 0, err
		}

		return warehouse, total, nil
	}
}

// GetWarehouseById implements [WarehouseRepositoryInterface].
func (w *warehouseRepository) GetWarehouseById(ctx context.Context, id uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetWarehouseById - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelWarehouse := model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetWarehouseById - 2: %v", err)
			return nil, err
		}

		return &modelWarehouse, nil
	}
}

// UpdateWarehouse implements [WarehouseRepositoryInterface].
func (w *warehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] UpdateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingWarehouse := model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", warehouse.ID).First(&existingWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] UpdateWarehouse - 2: %v", err)
			return err
		}

		existingWarehouse.Name = warehouse.Name
		existingWarehouse.Address = warehouse.Address
		existingWarehouse.Photo = warehouse.Photo
		existingWarehouse.Phone = warehouse.Phone

		if err := w.db.WithContext(ctx).Save(&existingWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] UpdateWarehouse - 3: %v", err)
			return err
		}

		return nil
	}
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepositoryInterface {
	return &warehouseRepository{db: db}
}
