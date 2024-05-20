package main

import (
	"fmt"

	"github.com/tanush-128/openzo_backend/online_order/config"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectToDB(cfg *config.Config) (*gorm.DB, error) {

	db := &gorm.DB{}
	err := error(nil)
	if cfg.MODE == "production" {
		dsn := cfg.DB_URL

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to open database connection: %w", err)
		}
		fmt.Println("Connected to database in production mode")
	} else {

		db, err = gorm.Open(
			sqlite.Open("test.db"),

			&gorm.Config{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to open database connection: %w", err)
		}
	}
	db.Migrator().AutoMigrate(&models.OnlineOrder{})
	db.Migrator().AutoMigrate(&models.OnlineCustomer{})
	db.Migrator().AutoMigrate(&models.OnlineOrderItem{})

	return db, nil
}
