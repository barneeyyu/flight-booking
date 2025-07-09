package repository

import (
	"flight-booking/internal/models"

	"gorm.io/gorm"
)

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Create(booking *models.Booking) error
	FindByID(id uint) (*models.Booking, error)
	Update(booking *models.Booking) error
}

// GORMBookingRepository is a concrete implementation of BookingRepository using GORM
type GORMBookingRepository struct {
	db *gorm.DB
}

// NewGORMBookingRepository creates a new GORMBookingRepository
func NewGORMBookingRepository(db *gorm.DB) *GORMBookingRepository {
	return &GORMBookingRepository{db: db}
}

// Create implements BookingRepository.Create
func (r *GORMBookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

// FindByID implements BookingRepository.FindByID
func (r *GORMBookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

// Update implements BookingRepository.Update
func (r *GORMBookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}
