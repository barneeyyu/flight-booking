package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("flights.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Flight{}, &Booking{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Search for flights
	r.GET("/flights", searchFlights)

	// Book a flight
	r.POST("/bookings", createBooking)

	// Get booking status
	r.GET("/bookings/:id", getBooking)

	// Get flight details
	r.GET("/flights/:id", getFlight)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func searchFlights(c *gin.Context) {
	var flights []Flight
	query := db

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

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	var total int64
	query.Model(&Flight{}).Count(&total)

	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	if err := query.Offset((p - 1) * ps).Limit(ps).Find(&flights).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"total":     total,
		"page":      p,
		"page_size": ps,
		"data":      flights,
	})
}

func createBooking(c *gin.Context) {
	var booking Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var flight Flight
	if err := db.First(&flight, booking.FlightID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Flight not found"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	oversold := rand.Intn(100) < 20 // 20% chance of overselling

	if flight.AvailableSeats > 0 && !oversold {
		flight.AvailableSeats--
		db.Save(&flight)
		booking.BookingStatus = "Confirmed"
	} else {
		booking.BookingStatus = "Waitlisted"
	}

	if err := db.Create(&booking).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, booking)
}

func getBooking(c *gin.Context) {
	var booking Booking
	if err := db.First(&booking, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(200, booking)
}

func getFlight(c *gin.Context) {
	var flight Flight
	if err := db.First(&flight, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(200, gin.H{
		"available_seats": flight.AvailableSeats,
		"price":           flight.Price,
	})
}
