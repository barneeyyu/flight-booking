package models

import "gorm.io/gorm"

// Flight represents a flight in the system
type Flight struct {
	gorm.Model
	FlightNumber     string  `json:"flight_number"`
	DepartureAirport string  `json:"departure_airport"`
	ArrivalAirport   string  `json:"arrival_airport"`
	DepartureTime    string  `json:"departure_time"`
	ArrivalTime      string  `json:"arrival_time"`
	Airline          string  `json:"airline"`
	Price            float64 `json:"price"`
	AvailableSeats   int     `json:"available_seats"`
}

// Booking represents a booking made by a user
type Booking struct {
	gorm.Model
	FlightID      uint   `json:"flight_id"`
	PassengerName string `json:"passenger_name"`
	BookingStatus string `json:"booking_status"` // e.g., "Confirmed", "Waitlisted"
}
