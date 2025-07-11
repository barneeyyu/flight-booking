package models

import "gorm.io/gorm"

// Flight represents a flight in the system

// TODO: 若查詢常用 (DepartureAirport, ArrivalAirport, DepartureTime)，可加複合索引提升效能

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
// TODO: 可加複合唯一鍵 (flight_id, passenger_name) 避免重複訂購
// TODO: 若查詢常用 (FlightID, PassengerName)，可加複合索引提升效能

// TODO: 若未來需支援付款流程，可新增付款相關欄位（如 payment_status, payment_time 等）
// TODO: 若需通知用戶，可考慮加上 email 或 notification 欄位

type Booking struct {
	gorm.Model
	FlightID      uint    `json:"flight_id"`
	PassengerName string  `json:"passenger_name"`
	Quantity      int     `json:"quantity"`
	TotalPrice    float64 `json:"total_price"`
	BookingStatus string  `json:"booking_status"` // e.g., "Confirmed", "Waitlisted"
	// PaymentStatus string // TODO: 付款狀態（如 unpaid, paid, refunded）
	// NotificationSent bool // TODO: 是否已通知用戶
}
