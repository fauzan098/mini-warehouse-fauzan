package controller

import (
	"micro-warehouse/warehouse-service/controller/request"
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/conv"
	"micro-warehouse/warehouse-service/pkg/pagination"
	"micro-warehouse/warehouse-service/pkg/validator"
	"micro-warehouse/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseControllerInterface interface {
	CreateWarehouse(ctx *fiber.Ctx) error
	GetAllWarehouse(ctx *fiber.Ctx) error
	GetWarehouseByID(ctx *fiber.Ctx) error
	UpdateWarehouse(ctx *fiber.Ctx) error
	DeleteWarehouse(ctx *fiber.Ctx) error
}

type warehouseController struct {
	warehouseUsecase usecase.WarehouseUsecaseinterface
}

// CreateWarehouse implements [WarehouseControllerInterface].
func (w *warehouseController) CreateWarehouse(ctx *fiber.Ctx) error {
	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := w.warehouseUsecase.CreateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse Created successfully",
	})
}

// DeleteWarehouse implements [WarehouseControllerInterface].
func (w *warehouseController) DeleteWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	if err := w.warehouseUsecase.DeleteWarehouse(ctx.Context(), warehouseID); err != nil {
		log.Errorf("[WarehouseController] DeleteWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "warehouse delete successfully",
	})
}

// GetAllWarehouse implements [WarehouseControllerInterface].
func (w *warehouseController) GetAllWarehouse(ctx *fiber.Ctx) error {
	var req request.GetAllWarehouseRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	warehouses, total, err := w.warehouseUsecase.GetAllWarehouse(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all warehouse",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	var warehousesResponse []response.WarehouseResponse
	for _, warehouse := range warehouses {
		warehousesResponse = append(warehousesResponse, response.WarehouseResponse{
			ID:           warehouse.ID,
			Name:         warehouse.Name,
			Address:      warehouse.Address,
			Photo:        warehouse.Photo,
			Phone:        warehouse.Phone,
			CountProduct: len(warehouse.WarehouseProducts),
		})
	}

	response := response.GetAllWarehouseResponse{
		Warehouse:  warehousesResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "warehouse fetched successfully",
	})
}

// GetWarehouseByID implements [WarehouseControllerInterface].
func (w *warehouseController) GetWarehouseByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	warehouse, err := w.warehouseUsecase.GetWarehouseById(ctx.Context(), warehouseID)
	if err != nil {
		log.Errorf("[WarehouseController] GetWarehouseById - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse",
		})
	}

	resp := response.DetailWarehouseResponse{
		ID:    warehouse.ID,
		Name:  warehouse.Name,
		Photo: warehouse.Photo,
		Phone: warehouse.Phone,
		Address: warehouse.Address,
	}

	var warehouseProductResponse []response.WarehouseProductResponse
	for _, warehouseProduct := range warehouse.WarehouseProducts {
		warehouseProductResponse = append(warehouseProductResponse, response.WarehouseProductResponse{
			ID:          warehouseProduct.ID,
			WarehouseID: warehouseProduct.WarehouseID,
			ProductID:   warehouseProduct.ProductID,
			Stock:       warehouseProduct.Stock,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "warehouse fetched successfully",
	})
}

// UpdateWarehouse implements [WarehouseControllerInterface].
func (w *warehouseController) UpdateWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		ID: warehouseID,
		Name: req.Name,
		Address: req.Address,
		Phone: req.Phone,
		Photo: req.Photo,
	}

	if err := w.warehouseUsecase.UpdateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "warehouse update successfully",
	})
}

func NewWarehouseController(warehouseUsecase usecase.WarehouseUsecaseinterface) WarehouseControllerInterface {
	return &warehouseController{warehouseUsecase: warehouseUsecase}
}
