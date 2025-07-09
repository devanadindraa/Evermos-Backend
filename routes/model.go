package routes

import (
	"context"
	"database/sql"

	"github.com/devanadindraa/Evermos-Backend/utils/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Dependency struct {
	handler *fiber.App
	db      *gorm.DB
}

func (d *Dependency) Close() {
	ctx := context.Background()

	db, err := d.db.DB()
	if err != nil {
		logger.Error(ctx, "Error closing database: %v", err)
		return
	}

	if err := db.Close(); err != nil {
		logger.Error(ctx, "Error closing database: %v", err)
	}
}

func (d *Dependency) GetHandler() *fiber.App {
	return d.handler
}

func (d *Dependency) GetDB() *sql.DB {
	ctx := context.Background()

	db, err := d.db.DB()
	if err != nil {
		logger.Error(ctx, "Error get database %v", err)
		return nil
	}

	return db
}
