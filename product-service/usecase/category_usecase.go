package usecase

import (
	"context"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/repository"
)

type CategoryUsecaseInterface interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetAllCategory(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepositoryInterface
}

// CreateCategory implements [CategoryUseCaseInterface].
func (c *categoryUsecase) CreateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryRepo.CreateCategory(ctx, category)
}

// DeleteCategory implements [CategoryUseCaseInterface].
func (c *categoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	return c.categoryRepo.DeleteCategory(ctx, id)
}

// GetAllCategory implements [CategoryUseCaseInterface].
func (c *categoryUsecase) GetAllCategory(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Category, int64, error) {
	return c.categoryRepo.GetAllCategory(ctx, page, limit, search, sortBy, sortOrder)
}

// GetCategoryByID implements [CategoryUseCaseInterface].
func (c *categoryUsecase) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	return c.categoryRepo.GetCategoryByID(ctx, id)
}

// UpdateCategory implements [CategoryUseCaseInterface].
func (c *categoryUsecase) UpdateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryRepo.UpdateCategory(ctx, category)
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepositoryInterface) CategoryUsecaseInterface {
	return &categoryUsecase{categoryRepo: categoryRepo}
}
