package service

import (
	"errors"
	"flight-booking/internal/models"
	"flight-booking/internal/repository"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	BookingStatusConfirmed  = "Confirmed"
	BookingStatusWaitlisted = "Waitlisted"
)

type BookingService interface {
	CreateBooking(booking *models.Booking) (*models.Booking, error)
	GetBooking(id uint) (*models.Booking, error)
}

type BookingServiceImpl struct {
	BookingRepo   repository.BookingRepository
	DB            *gorm.DB
	OversellLimit int
}

func NewBookingService(bookingRepo repository.BookingRepository, db *gorm.DB, oversellLimit int) BookingService {
	return &BookingServiceImpl{
		BookingRepo:   bookingRepo,
		DB:            db,
		OversellLimit: oversellLimit,
	}
}

func (s *BookingServiceImpl) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	oversellLimit := s.OversellLimit

	// Start a transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var flight models.Flight
		// Select flight with pessimistic lock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", booking.FlightID).
			First(&flight).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("flight not found")
			}
			return fmt.Errorf("failed to lock flight: %w", err)
		}

		// Check available seats with oversell logic
		// TODO: 超賣邏輯需要再優化，這裡只是做個簡單的範例
		if flight.AvailableSeats >= booking.Quantity {
			booking.BookingStatus = BookingStatusConfirmed
		} else if flight.AvailableSeats+oversellLimit >= booking.Quantity {
			booking.BookingStatus = BookingStatusWaitlisted
		} else {
			return fmt.Errorf("not enough seats: available=%d, oversell limit=%d", flight.AvailableSeats, oversellLimit)
		}

		// Deduct seats (can go negative due to oversell)
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
		// TODO: 訂單建立失敗，可設定發信通知用戶
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
