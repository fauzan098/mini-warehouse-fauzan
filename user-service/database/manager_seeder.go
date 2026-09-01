package database

import (
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SeedManager(db *gorm.DB) {
	bytes, err := conv.HashPassword("manager123")
	if err != nil {
		log.Fatalf("[SeedManager] SeedManager - 1: %v", err)
	}

	modelRole := model.Role{}
	err = db.Where("name = ?", "Manager").First(&modelRole).Error
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	admin := model.User{
		Name: "Manager",
		Email: "managa@gmail.com",
		Password: bytes,
		Roles: []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "manager@gmail.com"}).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	} else {
		log.Infof("Admin %s created", admin.Name)
	}

}