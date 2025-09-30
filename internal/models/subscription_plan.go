package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionPlan represents a subscription plan linked to a product
type SubscriptionPlan struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"` // Foreign key to Product
	PlanName  string    `gorm:"not null"`
	Duration  int       `gorm:"not null"` // Duration in days
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Belongs-to relationship with Product
	Product Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hook to generate UUID before creating a subscription plan
func (s *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for SubscriptionPlan model
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}
