package model

import "time"

type Product struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"type:varchar(100);not null"`
	Barcode    string     `json:"barcode" gorm:"type:varchar(100);uniqueIndex"`
	CategoryID uint       `json:"category_id"`
	Thumbnail  string     `json:"thumbnail"`
	About      string     `json:"about" gorm:"type:text"`
	Price      float64    `json:"price" gorm:"not null"`
	IsPopular  bool       `json:"is_popular" gorm:"type:boolean;default:false"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"update_at"`
	DeleteAt   *time.Time `json:"delete_at"`

	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
