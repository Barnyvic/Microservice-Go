package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product in the system
// Using GORM's polymorphic associations to support different product types
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	ProductType string    `gorm:"not null;index"` // For polymorphism support
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// One-to-many relationship with SubscriptionPlan
	SubscriptionPlans []SubscriptionPlan `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hook to generate UUID before creating a product
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

