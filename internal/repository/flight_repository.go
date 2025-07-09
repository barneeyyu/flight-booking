package repository

import (
	"flight-booking/internal/models"

	"gorm.io/gorm"
)

// FlightRepository defines the interface for flight data operations
type FlightRepository interface {
	FindAll(query *gorm.DB, page, pageSize int) ([]models.Flight, int64, error)
	FindByID(id uint) (*models.Flight, error)
	Create(flight *models.Flight) error
	Update(flight *models.Flight) error
}

// GORMFlightRepository is a concrete implementation of FlightRepository using GORM
type GORMFlightRepository struct {
	db *gorm.DB
}

// NewGORMFlightRepository creates a new GORMFlightRepository
func NewGORMFlightRepository(db *gorm.DB) *GORMFlightRepository {
	return &GORMFlightRepository{db: db}
}

// FindAll implements FlightRepository.FindAll
func (r *GORMFlightRepository) FindAll(query *gorm.DB, page, pageSize int) ([]models.Flight, int64, error) {
	var flights []models.Flight
	var total int64

	// Count total records
	if err := query.Model(&models.Flight{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and find records
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&flights).Error; err != nil {
		return nil, 0, err
	}

	return flights, total, nil
}

// FindByID implements FlightRepository.FindByID
func (r *GORMFlightRepository) FindByID(id uint) (*models.Flight, error) {
	var flight models.Flight
	if err := r.db.First(&flight, id).Error; err != nil {
		return nil, err
	}
	return &flight, nil
}

// Create implements FlightRepository.Create
func (r *GORMFlightRepository) Create(flight *models.Flight) error {
	return r.db.Create(flight).Error
}

// Update implements FlightRepository.Update
func (r *GORMFlightRepository) Update(flight *models.Flight) error {
	return r.db.Save(flight).Error
}
