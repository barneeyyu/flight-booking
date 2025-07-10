package service

import (
	"errors"
	"flight-booking/internal/models"
	"flight-booking/internal/repository"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingService interface {
	CreateBooking(booking *models.Booking) (*models.Booking, error)
	GetBooking(id uint) (*models.Booking, error)
}

type BookingServiceImpl struct {
	BookingRepo repository.BookingRepository
	DB          *gorm.DB
}

func NewBookingService(bookingRepo repository.BookingRepository, db *gorm.DB) BookingService {
	return &BookingServiceImpl{
		BookingRepo: bookingRepo,
		DB:          db,
	}
}

func (s *BookingServiceImpl) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	// Start a transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// TODO: 下單前先查詢是否已存在同一乘客同一航班的訂單，避免重複訂購
		// TODO: 像是-> existing, _ := s.BookingRepo.FindByFlightAndPassenger(booking.FlightID, booking.PassengerName)
		var flight models.Flight
		// Select flight with pessimistic lock and check available seats in DB
		// Use tx for all operations within the transaction
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND available_seats >= ?", booking.FlightID, booking.Quantity).
			First(&flight).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("flight not found or not enough seats available")
			}
			return fmt.Errorf("failed to lock flight: %w", err)
		}

		// Deduct seats
		flight.AvailableSeats -= booking.Quantity

		// Update flight within the transaction
		if err := tx.Save(&flight).Error; err != nil {
			return fmt.Errorf("failed to update flight seats: %w", err)
		}

		// Calculate total price
		booking.TotalPrice = float64(booking.Quantity) * flight.Price

		// Create booking within the transaction
		if err := tx.Create(&booking).Error; err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		return nil // Commit transaction
	})

	// TODO: 若需付款，這裡可串接金流並更新訂單狀態
	// TODO: 訂單建立成功後可透過 Queue 發送 email 或通知用戶

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingServiceImpl) GetBooking(id uint) (*models.Booking, error) {
	booking, err := s.BookingRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return booking, nil
}
