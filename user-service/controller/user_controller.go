package controller

import (
	"fmt"
	"micro-warehouse/user-service/controller/request"
	"micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/pkg/pagination"
	"micro-warehouse/user-service/pkg/validator"
	"micro-warehouse/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserControllerInterface interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error

	AssignUserToRole(c *fiber.Ctx) error
	EditAssignUserToRole(c *fiber.Ctx) error
	GetUserRoleById(c *fiber.Ctx) error
	GetUserRoleByRoleName(c *fiber.Ctx) error
	GetAllUserRoles(c *fiber.Ctx) error
}

type userController struct {
	userUsecase usecase.UserUsecaseInterface
}

// AssignUserToRole implements [UserControllerInterface].
func (u *userController) AssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.AssignUserToRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] AssignUserRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] AssignUserRole - 2 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := u.userUsecase.AssignUserToRole(ctx, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] AssignUserRole - 3 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user assign to role successfully",
	})
}

// CreateUser implements [UserControllerInterface].
func (u *userController) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.CreateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] CreateUser - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] CreateUser - 2 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userModel := model.User{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
		Phone: req.Phone,
		Photo: req.Photo,
	}

	if err := u.userUsecase.CreateUser(ctx, userModel); err != nil {
		log.Errorf("[UserController] CreateUser - 3 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user create successfully",
	})
}

// DeleteUser implements [UserControllerInterface].
func (u *userController) DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	userID := conv.StringToUint(id)

	if err := u.userUsecase.DeleteUser(ctx, userID); err != nil {
		log.Errorf("[UserController] DeleteUser - 1 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User delete successfully",
	})
}

// EditAssignUserToRole implements [UserControllerInterface].
func (u *userController) EditAssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.AssignUserToRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] EditAssignUserRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] EditAssignUserRole - 2 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userRoleIDStr := c.Params("userRoleID")
	userRoleID := conv.StringToUint(userRoleIDStr)

	if err := u.userUsecase.EditAssignUserToRole(ctx, userRoleID, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] EditAssignUserRole - 3 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user role update successfully",
	})
}

// GetAllUser implements [UserControllerInterface].
func (u *userController) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUserRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUser - 1 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUser - 2 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUser(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUser - 3 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserResponse{}
	for _, user := range users {
		roleName := ""
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Photo: user.Photo,
			RoleName: roleName,
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	response := response.GetAllUsersResponse{
		Users: resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
		"message": "users fetched succesfully",
	})
}

// GetAllUserRoles implements [UserControllerInterface].
func (u *userController) GetAllUserRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUserRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUser - 1 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUser - 2 %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUserRoles(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUser - 3 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserRoleResponse{}
	for _, user := range users {
		resp = append(resp, response.UserRoleResponse{
			ID: user.ID,
			UserID: user.UserID,
			RoleID: user.RoleID,
			User: response.UserResponse{
				ID: user.User.ID,
				Email: user.User.Email,
				Name: user.User.Name,
				Phone: user.User.Phone,
				Photo: user.User.Photo,
				RoleName: user.Role.Name,
			},
			Role: response.RoleResponse{
				ID: user.Role.ID,
				Name: user.Role.Name,
			},
		})
		
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	response := response.GetAllUsersRolesResponse{
		UsersRoles: resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
		"message": "users roles fetched succesfully",
	})
}

// GetUserById implements [UserControllerInterface].
func (u *userController) GetUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	userID := conv.StringToUint(id)

	user, err := u.userUsecase.GetUserById(ctx, userID)
	if err != nil {
		log.Errorf("[UserController] GetUserById - 1 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	roleName := ""
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	response := response.UserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Photo: user.Photo,
		RoleName: roleName,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (u *userController) GetUserRoleByRoleName(c *fiber.Ctx) error {
	ctx := c.Context()
	roleName := c.Params("roleName")

	users, err := u.userUsecase.GetUserByRoleName(ctx, roleName)
	if err != nil {
		log.Errorf("[UserController] GetUserRoleByRoleName - 1 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserResponse{}
	for _, user := range users {
		roleName := ""
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Photo: user.Photo,
			RoleName: roleName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
		"message": "users fetched successfully",
	})
}

// GetUserRoleById implements [UserControllerInterface].
func (u *userController) GetUserRoleById(c *fiber.Ctx) error {
	ctx := c.Context()
	userRoleIDStr := c.Params("UserRoleID")

	userRoleID := conv.StringToUint(userRoleIDStr)
	fmt.Println(userRoleID, "controller")
	userRole, err := u.userUsecase.GetUserRoleById(ctx, userRoleID)
	if err != nil {
		log.Errorf("[UserController] GetUserRoleById - 1 %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": userRole,
		"message": "users role retrieved successfully",
	})
}

// UpdateUser implements [UserControllerInterface].
func (u *userController) UpdateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	req := request.UpdateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] UpdateUser - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] UpdateUser - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID := conv.StringToUint(id)
	userModel := model.User{
		ID: userID,
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
		Phone: req.Phone,
		Photo: req.Photo,
	}

	if req.Password != "" {
		hashedPassword, err := conv.HashPassword(req.Password)
		if err != nil {
			log.Errorf("[UserControllere] UpdateUser 3: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		userModel.Password = hashedPassword
	}

	if  err := u.userUsecase.UpdateUser(ctx, userModel); err != nil {
		log.Errorf("[UserController] UpdateUser - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User update successfully",
	})
}

func NewUserController(userUsecase usecase.UserUsecaseInterface) UserControllerInterface {
	return &userController{
		userUsecase: userUsecase,
	}
}
