package router

import (
	"flight-booking/internal/handler"
	"flight-booking/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter sets up all the API routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize repositories
	flightRepo := repository.NewGORMFlightRepository(db)
	bookingRepo := repository.NewGORMBookingRepository(db)

	// Initialize handlers with their respective repositories
	flightHandler := handler.NewFlightHandler(flightRepo, db)
	bookingHandler := handler.NewBookingHandler(bookingRepo, flightRepo)

	// Public routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Flight routes
	r.GET("/flights", flightHandler.SearchFlights)
	r.GET("/flights/:id", flightHandler.GetFlight)

	// Booking routes
	r.POST("/bookings", bookingHandler.CreateBooking)
	r.GET("/bookings/:id", bookingHandler.GetBooking)

	return r
}
