package controller

import (
	"micro-warehouse/product-service/controller/request"
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/pkg/conv"
	"micro-warehouse/product-service/pkg/pagination"
	"micro-warehouse/product-service/pkg/validator"
	"micro-warehouse/product-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CategoryControllerInterface interface {
	CreateCategory(ctx *fiber.Ctx) error
	GetAllCategory(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error
}

type categoryController struct {
	categoryUsecase usecase.CategoryUsecaseInterface
}

// CreateCategory implements [CategoryControllerInterface].
func (c *categoryController) CreateCategory(ctx *fiber.Ctx) error {
	var req request.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[categoryController] CreateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[categoryController] CreateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Category{
		Name:    req.Name,
		Tagline: req.Tagline,
		Photo:   req.Photo,
	}

	if err := c.categoryUsecase.CreateCategory(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[categoryController] CreateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to create category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category create successfully",
	})
}

// DeleteCategory implements [CategoryControllerInterface].
func (c *categoryController) DeleteCategory(ctx *fiber.Ctx) error {
	context := ctx.Context()

	idCategory := ctx.Params("id")
	if idCategory == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID is required",
		})
	}

	id := conv.StringToUint(idCategory)

	if err := c.categoryUsecase.DeleteCategory(context, id); err != nil {
		log.Errorf("[categoryController] DeleteCategory - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// GetAllCategory implements [CategoryControllerInterface].
func (c *categoryController) GetAllCategory(ctx *fiber.Ctx) error {
	var req request.GetAllCategoryRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[categoryController] GetAllCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	categories, total, err := c.categoryUsecase.GetAllCategory(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[categoryController] GetAllCategory - 2: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all categories",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	var categoriesResponse []response.CategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, response.CategoryResponse{
			ID:      category.ID,
			Name:    category.Name,
			Tagline: category.Tagline,
			Photo:   category.Photo,
		})
	}

	response := response.GetAllCategoryResponse{
		Categories: categoriesResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data":    response,
	})
}

// GetCategoryByID implements [CategoryControllerInterface].
func (c *categoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idToUint := conv.StringToUint(id)

	category, err := c.categoryUsecase.GetCategoryByID(ctx.Context(), idToUint)
	if err != nil {
		log.Errorf("[categoryController] GetCategoryByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get category by id",
		})
	}

	response := response.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Tagline:      category.Tagline,
		Photo:        category.Photo,
		CountProduct: len(category.Products),
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data":    response,
	})
}

// UpdateCategory implements [CategoryControllerInterface].
func (c *categoryController) UpdateCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idToUint := conv.StringToUint(id)

	var req request.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[categoryController] UpdateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[categoryController] UpdateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Category{
		ID:      idToUint,
		Name:    req.Name,
		Tagline: req.Tagline,
		Photo:   req.Photo,
	}

	if err := c.categoryUsecase.UpdateCategory(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[categoryController] UpdateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category update successfully",
	})
}

func NewCategoryController(categoryUsecase usecase.CategoryUsecaseInterface) CategoryControllerInterface {
	return &categoryController{categoryUsecase: categoryUsecase}
}
