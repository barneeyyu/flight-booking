package models

import "gorm.io/gorm"

// Flight represents a flight in the system
type Flight struct {
	gorm.Model
	FlightNumber     string  `json:"flight_number" gorm:"index"`
	DepartureAirport string  `json:"departure_airport" gorm:"index:idx_flight_search"`
	ArrivalAirport   string  `json:"arrival_airport" gorm:"index:idx_flight_search"`
	DepartureTime    string  `json:"departure_time" gorm:"index:idx_flight_search"`
	ArrivalTime      string  `json:"arrival_time"`
	Airline          string  `json:"airline" gorm:"index"`
	Price            float64 `json:"price" gorm:"index"`
	AvailableSeats   int     `json:"available_seats"`
}

// Booking represents a booking made by a user
// TODO: 若未來需支援付款流程，可新增付款相關欄位（如 payment_status, payment_time 等）
// TODO: 若需通知用戶，可考慮加上 email 或 notification 欄位
type Booking struct {
	gorm.Model
	FlightID      uint    `json:"flight_id" gorm:"index:idx_booking_search"`
	PassengerName string  `json:"passenger_name" gorm:"index:idx_booking_search"`
	Quantity      int     `json:"quantity"`
	TotalPrice    float64 `json:"total_price"`
	BookingStatus string  `json:"booking_status" gorm:"index"` // e.g., "Confirmed", "Waitlisted"
	// PaymentStatus string // TODO: 付款狀態（如 unpaid, paid, refunded）
	// NotificationSent bool // TODO: 是否已通知用戶
}
