package usecase

import (
	"context"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/repository"
)

type WarehouseUsecaseinterface interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetAllWarehouse(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Warehouse, int64, error)
	GetWarehouseById(ctx context.Context, id uint) (*model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id uint) error
}

type warehouseUsecase struct {
	warehouseRepo repository.WarehouseRepositoryInterface
}

// CreateWarehouse implements [WarehouseUsecaseinterface].
func (w *warehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return w.warehouseRepo.CreateWarehouse(ctx, warehouse)
}

// DeleteWarehouse implements [WarehouseUsecaseinterface].
func (w *warehouseUsecase) DeleteWarehouse(ctx context.Context, id uint) error {
	return w.warehouseRepo.DeleteWarehouse(ctx, id)
}

// GetAllWarehouse implements [WarehouseUsecaseinterface].
func (w *warehouseUsecase) GetAllWarehouse(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error) {
	return w.warehouseRepo.GetAllWarehouse(ctx, page, limit, search, sortBy, sortOrder)
}

// GetWarehouseById implements [WarehouseUsecaseinterface].
func (w *warehouseUsecase) GetWarehouseById(ctx context.Context, id uint) (*model.Warehouse, error) {
	return w.warehouseRepo.GetWarehouseById(ctx, id)
}

// UpdateWarehouse implements [WarehouseUsecaseinterface].
func (w *warehouseUsecase) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return w.warehouseRepo.UpdateWarehouse(ctx, warehouse)
}

func NewWarehouseUsecase(warehouseRepo repository.WarehouseRepositoryInterface) WarehouseUsecaseinterface {
	return &warehouseUsecase{warehouseRepo: warehouseRepo}
}
