package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
	PlanName  string    `gorm:"not null"`
	Duration  int       `gorm:"not null"` 
	Price     float64   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Product Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}

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
