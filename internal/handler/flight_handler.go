package handler

import (
	"flight-booking/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FlightHandler handles flight-related HTTP requests
type FlightHandler struct {
	FlightRepo repository.FlightRepository
	db         *gorm.DB // Still need db for query building
}

// NewFlightHandler creates a new FlightHandler
func NewFlightHandler(flightRepo repository.FlightRepository, db *gorm.DB) *FlightHandler {
	return &FlightHandler{FlightRepo: flightRepo, db: db}
}

// SearchFlights handles flight search requests
func (h *FlightHandler) SearchFlights(c *gin.Context) {
	query := h.db // Use the injected db for query building

	if departure := c.Query("departure_airport"); departure != "" {
		query = query.Where("departure_airport = ?", departure)
	}

	if arrival := c.Query("arrival_airport"); arrival != "" {
		query = query.Where("arrival_airport = ?", arrival)
	}

	if airline := c.Query("airline"); airline != "" {
		query = query.Where("airline = ?", airline)
	}

	if date := c.Query("date"); date != "" {
		query = query.Where("DATE(departure_time) = ?", date)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	flights, total, err := h.FlightRepo.FindAll(query, page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      flights,
	})
}

// GetFlight handles requests to get a single flight by ID
func (h *FlightHandler) GetFlight(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, err := h.FlightRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(200, gin.H{
		"available_seats": flight.AvailableSeats,
		"price":           flight.Price,
	})
}
