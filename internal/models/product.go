package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	ProductType string    `gorm:"not null;index"` 
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	SubscriptionPlans []SubscriptionPlan `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}


func (Product) TableName() string {
	return "products"
}
