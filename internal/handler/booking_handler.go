package handler

import (
	"flight-booking/internal/models"
	"flight-booking/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	BookingService service.BookingService
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{BookingService: bookingService}
}

// CreateBooking handles flight booking requests
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if booking.Quantity <= 0 {
		c.JSON(400, gin.H{"error": "Quantity must be a positive integer"})
		return
	}

	createdBooking, err := h.BookingService.CreateBooking(&booking)
	if err != nil {
		// Handle errors from the service layer
		if strings.Contains(err.Error(), "flight not found") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "not enough seats") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Internal Server Error: " + err.Error()})
		}
		return
	}

	c.JSON(200, createdBooking)
}

// GetBooking handles requests to get a single booking by ID
func (h *BookingHandler) GetBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.BookingService.GetBooking(uint(id))
	if err != nil {
		// Handle errors from the service layer
		if strings.Contains(err.Error(), "booking not found") {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Internal Server Error: " + err.Error()})
		}
		return
	}

	c.JSON(200, booking)
}
