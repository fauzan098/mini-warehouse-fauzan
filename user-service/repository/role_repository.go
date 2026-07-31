package repository

import (
	"context"
	"errors"

	"micro-warehouse/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	CreateRole(ctx context.Context, role model.Role) error
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleById(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type reloRepository struct {
	db *gorm.DB
}

// CreateRole implements [RoleRepositoryInterface].
func (r reloRepository) CreateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] CreateRole - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return r.db.WithContext(ctx).Create(&role).Error
	}
}

// DeleteRole implements [RoleRepositoryInterface].
func (r reloRepository) DeleteRole(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] DeleteRole - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&modelRole).Error; err != nil {
			log.Errorf("[RoleRepository] DeleteRole - 2: %v", err)
			return err
		}

		if len(modelRole.Users) > 0 {
			log.Errorf("[RoleRepository] DeleteRole - 3: Role with ID %v", "Role has users")
			return errors.New("role has associated users, cannot delete")
		}

		return r.db.WithContext(ctx).Delete(&modelRole).Error
	}
}

// GetAllRoles implements [RoleRepositoryInterface].
func (r reloRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] GetAllRoles - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRoles := []model.Role{}
		err := r.db.WithContext(ctx).Preload("Users").Find(&modelRoles).Error
		if err != nil {
			log.Errorf("[RoleRepository] GetAllRoles - 2: %v", err)
			return nil, err
		}
		return modelRoles, nil
	}
}

// GetRoleById implements [RoleRepositoryInterface].
func (r reloRepository) GetRoleById(ctx context.Context, id uint) (*model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] GetRoleById - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRole := model.Role{}
		err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&modelRole).Error
		if err != nil {
			log.Errorf("[RoleRepository] GetRoleById - 2: %v", err)
			return nil, err
		}
		return &modelRole, nil
	}
}

// UpdateRole implements [RoleRepositoryInterface].
func (r reloRepository) UpdateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] UpdateRole - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Where("id = ?", role.ID).First(&modelRole).Error; err != nil {
			log.Errorf("[RoleRepository] UpdateRole - 2: %v", err)
			return err
		}

		modelRole.Name = role.Name
		return r.db.WithContext(ctx).Save(&modelRole).Error
	}

}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return reloRepository{db: db}
}
