package repository

import (
	"context"
	"errors"

	"micro-warehouse/product-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetAllCategory(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements [CategoryRepositoryInterface].
func (c *categoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	select {
	case <-ctx.Done():
		log.Errorf("[categoryRepository] CreateCategory - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return c.db.WithContext(ctx).Create(category).Error
	}
}

// DeleteCategory implements [CategoryRepositoryInterface].
func (c *categoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[categoryRepository] DeleteCategory - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelCategory := model.Category{}
		if err := c.db.WithContext(ctx).Where("id = ?", id).Preload("Products").First(&modelCategory).Error; err != nil {
			log.Errorf("[categoryRepository] DeleteCategory - 2: %v", err)
			return nil
		}
		if len(modelCategory.Products) > 0 {
			log.Errorf("[categoryRepository] DeleteCategory - 3: %v", errors.New("category has products"))
			return errors.New("category has products")
		}

		return c.db.WithContext(ctx).Delete(&modelCategory).Error
	}
}

// GetAllCategory implements [CategoryRepositoryInterface].
func (c *categoryRepository) GetAllCategory(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Category, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[categoryRepository] GetAllCategory - 1: %v", ctx.Err())
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
		query := c.db.Model(&model.Category{})

		if search != "" {
			query = query.Where("name ILIKE ? OR tagline ILIKE ?", "%"+search+"%", "%"+search+"%")
		}

		// Get total count
		var totalRecords int64
		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[categoryRepository] GetAllCategory - 2: %v", err)
			return nil, 0, err
		}

		var modelCategories []model.Category
		if err := query.Order(sortBy + " " + sortOrder).
			WithContext(ctx).
			Preload("Products").
			Offset(offset).
			Limit(limit).
			Find(&modelCategories).Error; err != nil {
			log.Errorf("[categoryRepository] GetAllCategory - 3: %v", err)
			return nil, 0, err
		}

		return modelCategories, totalRecords, nil
	}
}

// GetCategoryByID implements [CategoryRepositoryInterface].
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[categoryRepository] GetCategoryByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelCategory := model.Category{}
		if err := c.db.WithContext(ctx).Where("id = ?", id).Preload("Products").First(&modelCategory).Error; err != nil {
			log.Errorf("[categoryRepository] GetCategoryByID - 2: %v", err)
			return nil, err
		}
		return &modelCategory, nil
	}
}

// UpdateCategory implements [CategoryRepositoryInterface].
func (c *categoryRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	select {
	case <-ctx.Done():
		log.Errorf("[categoryRepository] UpdateCategory - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelCategory := model.Category{}
		if err := c.db.WithContext(ctx).Where("id = ?", category.ID).First(&modelCategory).Error; err != nil {
			log.Errorf("[categoryRepository] UpdateCategory - 2: %v", err)
			return err
		}
		modelCategory.Name = category.Name
		modelCategory.Tagline = category.Tagline
		modelCategory.Photo = category.Photo
		return c.db.WithContext(ctx).Save(&modelCategory).Error
	}
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepository{db: db}
}
