package database

import (
	"fmt"
	"time"

	"github.com/WebCraftersGH/Education-service/internal/config"
	problemcatalog "github.com/WebCraftersGH/Education-service/internal/repository/problem_catalog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err := db.AutoMigrate(&problemcatalog.Problem{}, &problemcatalog.ProblemContent{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
