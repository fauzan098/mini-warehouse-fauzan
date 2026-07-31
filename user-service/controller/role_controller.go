package controller

import (
	"fmt"
	"micro-warehouse/user-service/controller/request"
	"micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/pkg/validator"
	"micro-warehouse/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type RoleControllerInterface interface {
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
	GetRoleById(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
}

type roleController struct {
	roleUseCase usecase.RoleUsecaseInterface
}

// CreateRole implements [RoleControllerInterface].
func (r *roleController) CreateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreatRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[RoleController] CreateRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[RoleController] CreateRole - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Role{
		Name: req.Name,
	}

	if err := r.roleUseCase.CreateRole(ctx, reqModel); err != nil {
		log.Errorf("[RoleController] CreateRole - 3: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Role create successfully",
	})
}

// DeleteRole implements [RoleControllerInterface].
func (r *roleController) DeleteRole(c *fiber.Ctx) error {
	ctx := c.Context()

	roleID := c.Params("id")
	if roleID != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID is required",
		})
	}

	id := conv.StringToUint(roleID)

	if err := r.roleUseCase.DeleteRole(ctx, id); err != nil {
		log.Errorf("[RoleController] DeleteRole - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}

// GetAllRoles implements [RoleControllerInterface].
func (r *roleController) GetAllRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	roles, err := r.roleUseCase.GetAllRoles(ctx)
	if err != nil {
		log.Errorf("[RoleController] GetAllRoles - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.RoleResponse{}
	for _, role := range roles {
		resp = append(resp, response.RoleResponse{
			ID: role.ID,
			Name: role.Name,
			CountUsers: int64(len(role.Users)),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Roles fetched successfully",
		"data": resp,
	})
}

// GetRoleById implements [RoleControllerInterface].
func (r *roleController) GetRoleById(c *fiber.Ctx) error {
	ctx := c.Context()

	roleID := c.Params("id")
	fmt.Println(roleID)
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID is required",
		})
	}

	id := conv.StringToUint(roleID)
	role, err := r.roleUseCase.GetRoleById(ctx, id)
	if err != nil {
		log.Errorf("[RoleController] GetRoleById - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Roles fetched successfully",
		"data": role,
	})
}

// UpdateRole implements [RoleControllerInterface].
func (r *roleController) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreatRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[RoleController] UpdateRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[RoleController] UpdateRole - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Role{
		ID: conv.StringToUint(c.Params("id")),
		Name: req.Name,
	}

	if  err := r.roleUseCase.UpdateRole(ctx, reqModel); err != nil {
		log.Errorf("[RoleController] UpdateRole - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role Update Successfully",
	})
}

func NewRoleController(roleUsecase usecase.RoleUsecaseInterface) RoleControllerInterface {
	return &roleController{roleUseCase: roleUsecase}
}
