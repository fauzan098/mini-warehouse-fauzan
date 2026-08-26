package usecase

import (
	"context"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/repository"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseProductUsecaseInterface interface {
	GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, []httpclient.ProductResponse, error)
	GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, *httpclient.ProductResponse, error)
	CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID, ProductID uint) (*model.WarehouseProduct, *httpclient.ProductResponse, error)
	UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error
	DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error
	GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error)
	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
}

type warehouseProductUsecase struct {
	warehouseProductRepo repository.WarehouseProductRepositoryInterface
	productClient        httpclient.ProductClientInterface
}

// CreateWarehouseProduct implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	result, err := w.warehouseProductRepo.GetWarehouseProductByWarehouseIDAndProductID(ctx, warehouseProduct.WarehouseID, warehouseProduct.ProductID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Errorf("[warehouseProductUsecase] CreateWarehouseProduct - 1: %v", err)
			return err	
		}

		log.Errorf("[warehouseProductUsecase] CreateWarehouseProduct - 2: %v", err)
		return err
	}

	if result != nil {
		warehouseProduct.ID = result.ID
		return w.warehouseProductRepo.UpdateWarehouseProduct(ctx, warehouseProduct)
	}

	return w.warehouseProductRepo.CreateWarehouseProduct(ctx, warehouseProduct)
}

// DeleteAllWarehouseProductByProductID implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error {
	return w.warehouseProductRepo.DeleteAllWarehouseProductByProductID(ctx, productID)
}

// DeleteWarehouseProduct implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error {
	return w.warehouseProductRepo.DeleteWarehouseProduct(ctx, warehouseProductID)
}

// GetDetailWarehouse implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, []httpclient.ProductResponse, error) {
	warehouse, err := w.warehouseProductRepo.GetDetailWarehouse(ctx, warehouseID)
	if err != nil {
		log.Errorf("[warehouseProductUsecase] GetDetailWarehouse - 1: %v", err)
		return nil, nil, err
	}

	var products []httpclient.ProductResponse
	if len(warehouse.WarehouseProducts) > 0 {
		for _, wp := range warehouse.WarehouseProducts {
			product, err := w.productClient.GetProductByID(ctx, wp.ProductID)
			if err != nil {
				log.Errorf("[warehouseProductUsecase] GetDetailWarehouse - 2: %v", err)
				return nil, nil, err
			}

			products = append(products, *product)
		}
	}

	return warehouse, products, nil
}

// GetDetailWarehouseProductByID implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, *httpclient.ProductResponse, error) {
	warehouseProduct, err := w.warehouseProductRepo.GetDetailWarehouseProductByID(ctx, warehouseProductID)
	if err != nil {
		log.Errorf("[warehouseProductUsecase] GetDetailWarehouseProductByID - 1: %v", err)
		return nil, nil, err
	}

	product, err := w.productClient.GetProductByID(ctx, warehouseProduct.ID)
	if err != nil {
		log.Errorf("[warehouseProductUsecase] GetDetailWarehouseProductByID - 2: %v", err)
		return nil, nil, err
	}

	return warehouseProduct, product, nil
}

// GetProductTotalStock implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	return w.warehouseProductRepo.GetProductTotalStock(ctx, productID)
}

// GetWarehouseProductByProductID implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error) {
	return w.warehouseProductRepo.GetWarehouseProductByProductID(ctx, productID)
}

// GetWarehouseProductByWarehouseIDAndProductID implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID uint, ProductID uint) (*model.WarehouseProduct, *httpclient.ProductResponse, error) {
	warehouseProduct, err := w.warehouseProductRepo.GetWarehouseProductByWarehouseIDAndProductID(ctx, warehouseID, ProductID)
	if err != nil {
		log.Errorf("[WarehouseProductUsecase] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", err)
		return nil, nil, err
	}

	product, err := w.productClient.GetProductByID(ctx, warehouseProduct.ProductID)
	if err != nil {
		log.Errorf("[WarehouseProductUsecase] GetWarehouseProductByWarehouseIDAndProductID - 2: %v", err)
		return nil, nil, err
	}

	return warehouseProduct, product, nil
}

// UpdateWarehouseProduct implements [WarehouseProductUsecaseInterface].
func (w *warehouseProductUsecase) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	return w.warehouseProductRepo.UpdateWarehouseProduct(ctx, warehouseProduct)
}

func NewWarehouseProductUsecase(warehouseProductRepo repository.WarehouseProductRepositoryInterface, productClient httpclient.ProductClientInterface) WarehouseProductUsecaseInterface {
	return &warehouseProductUsecase{
		warehouseProductRepo: warehouseProductRepo,
		productClient: productClient,
	}
}
