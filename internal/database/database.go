package database

import (
	"fmt"
	"log"
	"time"

	apperrors "github.com/microservice-go/product-service/internal/errors"
	"github.com/microservice-go/product-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
)

type Config struct {
	Driver   string 
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config Config) (*gorm.DB, error) {
	if config.Driver == "" {
		return nil, apperrors.NewValidationError("driver", "database driver is required")
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var db *gorm.DB
	var err error

	switch config.Driver {
	case "postgres":
		if err := validatePostgresConfig(config); err != nil {
			return nil, err
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	case "sqlite":
		if config.DBName == "" {
			return nil, apperrors.NewValidationError("dbname", "database name is required for SQLite")
		}
		db, err = gorm.Open(sqlite.Open(config.DBName), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (supported: postgres, sqlite)", config.Driver)
	}

	if err != nil {
		return nil, apperrors.NewDatabaseError("connection", err)
	}

	if config.Driver == "postgres" {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, apperrors.NewDatabaseError("pool configuration", err)
		}
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetConnMaxLifetime(connMaxLifetime)
	}

	log.Printf("Database connection established (driver: %s)", config.Driver)
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return apperrors.NewValidationError("db", "database connection is nil")
	}

	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.Product{},
		&models.SubscriptionPlan{},
	)

	if err != nil {
		return apperrors.NewDatabaseError("migration", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func validatePostgresConfig(config Config) error {
	if config.Host == "" {
		return apperrors.NewValidationError("host", "host is required for PostgreSQL")
	}
	if config.User == "" {
		return apperrors.NewValidationError("user", "user is required for PostgreSQL")
	}
	if config.DBName == "" {
		return apperrors.NewValidationError("dbname", "database name is required for PostgreSQL")
	}
	return nil
}
