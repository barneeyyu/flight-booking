package handler

import (
	"flight-booking/internal/models"
	"flight-booking/internal/repository"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	BookingRepo repository.BookingRepository
	FlightRepo  repository.FlightRepository // Need FlightRepo to update available seats
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(bookingRepo repository.BookingRepository, flightRepo repository.FlightRepository) *BookingHandler {
	return &BookingHandler{BookingRepo: bookingRepo, FlightRepo: flightRepo}
}

// CreateBooking handles flight booking requests
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	flight, err := h.FlightRepo.FindByID(booking.FlightID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Flight not found"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	oversold := rand.Intn(100) < 20 // 20% chance of overselling

	if flight.AvailableSeats > 0 && !oversold {
		flight.AvailableSeats--
		h.FlightRepo.Update(flight)
		booking.BookingStatus = "Confirmed"
	} else {
		booking.BookingStatus = "Waitlisted"
	}

	if err := h.BookingRepo.Create(&booking); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, booking)
}

// GetBooking handles requests to get a single booking by ID
func (h *BookingHandler) GetBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.BookingRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(200, booking)
}
