package database

import (
	"flight-booking/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and performs auto-migrations
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("flights.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.Flight{}, &models.Booking{})

	return db, nil
}
