package usecase

import (
	"context"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/repository"
	"micro-warehouse/user-service/service"

	"github.com/gofiber/fiber/v2/log"
)

type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, user model.User) error
	GetAllUser(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.User, int64, error)
	GetUserById(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id uint) error

	GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error)

	AssignUserToRole(ctx context.Context, UserId uint, RoleId uint) error
	EditAssignUserToRole(ctx context.Context, assignRoleId uint, UserId uint, RoleId uint) error
	GetUserRoleById(ctx context.Context, assignRoleId uint) (*model.UserRole, error)
	GetAllUserRoles(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.UserRole, int64, error)
}

type userUsecase struct {
	userRepo repository.UserRepositoryInterface
	rabbitMQService service.RabbitMQServiceInterface
}

// AssignUserToRole implements [userUsecaseInterface].
func (u *userUsecase) AssignUserToRole(ctx context.Context, UserId uint, RoleId uint) error {
	return u.userRepo.AssignUserToRole(ctx, UserId, RoleId)
}

// CreateUser implements [userUsecaseInterface].
func (u *userUsecase) CreateUser(ctx context.Context, user model.User) error {
	password, err := conv.HashPassword(user.Password)
	if err != nil {
		log.Errorf("[UserUseCase] CreateUser - : %v", err)
	}

	uncryptedPassword := user.Password
	user.Password = password
	result, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Errorf("[UserUseCase] CreateUser - 2: %v", err)
		return err
	}

	emailPayload := service.EmailPayload{
		Email: result.Email,
		Password: uncryptedPassword,
		Type: "welcome_email",
		UserID: result.ID,
		Name: result.Name,
	}

	go func ()  {
		if err := u.rabbitMQService.PublisEmail(ctx, emailPayload); err != nil {
			log.Errorf("[UserUseCase] CreateUser - 3: %v", err)
		}
	}()

	return nil
}

// DeleteUser implements [userUsecaseInterface].
func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	_, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		log.Errorf("[UserUsecase] DeleteUser - 1: %v", err)
		return err
	}

	if err := u.userRepo.DeleteUser(ctx, id); err != nil {
		log.Errorf("[UserUsecase] DeleteUser - 2: %v", err)
		return err
	}
	return nil
}

// EditAssignUserToRole implements [userUsecaseInterface].
func (u *userUsecase) EditAssignUserToRole(ctx context.Context, assignRoleId uint, UserId uint, RoleId uint) error {
	return u.userRepo.EditAssignUserToRole(ctx, assignRoleId, UserId, RoleId)
}

// GetAllUser implements [userUsecaseInterface].
func (u *userUsecase) GetAllUser(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	// get users from repository
	users, totalRecords, err := u.userRepo.GetAllUser(ctx, page, limit, search, sortBy, sortOrder)
	if err != nil  {
		log.Errorf("[UserUsecase] GetAllUsers - 2: %v", err)
		return nil, 0, err
	}

	return users, totalRecords, nil
}

// GetAllUserRoles implements [userUsecaseInterface].
func (u *userUsecase) GetAllUserRoles(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.UserRole, int64, error) {
	return u.userRepo.GetAllUserRoles(ctx, page, limit, search, sortBy, sortOrder)
}

// GetUserByEmail implements [userUsecaseInterface].
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

// GetUserById implements [userUsecaseInterface].
func (u *userUsecase) GetUserById(ctx context.Context, id uint) (*model.User, error) {
	return u.userRepo.GetUserById(ctx, id)
}

// GetUserByRoleName implements [userUsecaseInterface].
func (u *userUsecase) GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error) {
	return u.userRepo.GetUserByRoleName(ctx, roleName)
}


// GetUserRoleById implements [userUsecaseInterface].
func (u *userUsecase) GetUserRoleById(ctx context.Context, assignRoleId uint) (*model.UserRole, error) {
	return u.userRepo.GetUserRoleById(ctx, assignRoleId)
}

// UpdateUser implements [userUsecaseInterface].
func (u *userUsecase) UpdateUser(ctx context.Context, user model.User) error {
	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		log.Errorf("[UserUsecase] UpdtaeUser - 2: %v", err)
		return err
	}
	
	return nil
}

func NewUserUsecase(userRepo repository.UserRepositoryInterface, rabbitMQService service.RabbitMQServiceInterface) UserUsecaseInterface {
	return &userUsecase{
		userRepo: userRepo,
		rabbitMQService: rabbitMQService,
	}
}
