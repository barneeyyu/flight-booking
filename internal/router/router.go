package router

import (
	"flight-booking/internal/handler"
	"flight-booking/internal/repository"
	"flight-booking/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter sets up all the API routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize repositories
	flightRepo := repository.NewGORMFlightRepository(db)
	bookingRepo := repository.NewGORMBookingRepository(db)

	// Initialize services
	bookingService := service.NewBookingService(bookingRepo, db, 10) // 設定超賣上限為 10 張

	// Initialize handlers with their respective repositories/services
	flightHandler := handler.NewFlightHandler(flightRepo, db)
	bookingHandler := handler.NewBookingHandler(bookingService)

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
