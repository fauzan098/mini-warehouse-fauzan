package controller

import (
	"micro-warehouse/product-service/controller/request"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/pkg/conv"
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
	panic("unimplemented")
}

// GetCategoryByID implements [CategoryControllerInterface].
func (c *categoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateCategory implements [CategoryControllerInterface].
func (c *categoryController) UpdateCategory(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryController(categoryUsecase usecase.CategoryUsecaseInterface) CategoryControllerInterface {
	return &categoryController{categoryUsecase: categoryUsecase}
}
