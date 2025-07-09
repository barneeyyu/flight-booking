package handler

import (
	"flight-booking/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FlightSearchItem represents a flight item in the search results
type FlightSearchItem struct {
	ID               uint    `json:"id"`
	DepartureAirport string  `json:"departure_airport"`
	ArrivalAirport   string  `json:"arrival_airport"`
	DepartureTime    string  `json:"departure_time"`
	ArrivalTime      string  `json:"arrival_time"`
	Airline          string  `json:"airline"`
	Price            float64 `json:"price"`
	// FlightNumber and AvailableSeats are intentionally omitted
}

// SearchFlightsResponse is the full response structure for flight search
type SearchFlightsResponse struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Data     []FlightSearchItem `json:"data"`
}

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

	if departure := c.Query("departure"); departure != "" {
		query = query.Where("departure_airport = ?", departure)
	}

	if arrival := c.Query("arrival"); arrival != "" {
		query = query.Where("arrival_airport = ?", arrival)
	}

	if airline := c.Query("airline"); airline != "" {
		query = query.Where("airline = ?", airline)
	}

	dateStr := c.Query("date")
	if dateStr != "" {
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format. Expected YYYY-MM-DD"})
			return
		}
		query = query.Where("DATE(departure_time) = ?", dateStr)
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "Invalid page parameter. Must be a positive integer."})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(400, gin.H{"error": "Invalid page_size parameter. Must be a positive integer."})
		return
	}

	flights, total, err := h.FlightRepo.FindAll(query, page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Convert models.Flight to FlightSearchItem to exclude specific fields
	var searchItems []FlightSearchItem
	for _, flight := range flights {
		searchItems = append(searchItems, FlightSearchItem{
			ID:               flight.ID,
			DepartureAirport: flight.DepartureAirport,
			ArrivalAirport:   flight.ArrivalAirport,
			DepartureTime:    flight.DepartureTime,
			ArrivalTime:      flight.ArrivalTime,
			Airline:          flight.Airline,
			Price:            flight.Price,
		})
	}

	c.JSON(200, SearchFlightsResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     searchItems,
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

	c.JSON(200, flight)
}
