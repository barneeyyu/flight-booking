package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"flight-booking/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockBookingService is a mock implementation of BookingService interface
type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) CreateBooking(booking *models.Booking) (*models.Booking, error) {
	args := m.Called(booking)
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockBookingService) GetBooking(id uint) (*models.Booking, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Booking), args.Error(1)
}

// SetupRouter for testing
func setupTestRouter(bookingHandler *BookingHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/bookings", bookingHandler.CreateBooking)
	r.GET("/bookings/:id", bookingHandler.GetBooking)
	return r
}

// TestCreateBooking_Success tests a successful booking creation
func TestCreateBooking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingReq := models.Booking{
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      2,
	}
	expectedBooking := models.Booking{
		Model:         gorm.Model{ID: 1},
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      2,
		TotalPrice:    200.0,
		BookingStatus: "Confirmed",
	}

	mockService.On("CreateBooking", &bookingReq).Return(&expectedBooking, nil).Once()

	jsonValue, _ := json.Marshal(bookingReq)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBooking models.Booking
	json.Unmarshal(w.Body.Bytes(), &responseBooking)
	assert.Equal(t, expectedBooking.ID, responseBooking.ID)
	assert.Equal(t, expectedBooking.PassengerName, responseBooking.PassengerName)
	assert.Equal(t, expectedBooking.Quantity, responseBooking.Quantity)
	assert.Equal(t, expectedBooking.TotalPrice, responseBooking.TotalPrice)
	assert.Equal(t, expectedBooking.BookingStatus, responseBooking.BookingStatus)

	mockService.AssertExpectations(t)
}

// TestCreateBooking_InvalidQuantity tests booking creation with invalid quantity
func TestCreateBooking_InvalidQuantity(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingReq := models.Booking{
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      0, // Invalid quantity
	}

	jsonValue, _ := json.Marshal(bookingReq)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Quantity must be a positive integer")

	mockService.AssertNotCalled(t, "CreateBooking", mock.Anything)
}

// TestCreateBooking_FlightNotFound tests booking creation when flight is not found
func TestCreateBooking_FlightNotFound(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingReq := models.Booking{
		FlightID:      999, // Non-existent flight
		PassengerName: "Test User",
		Quantity:      1,
	}

	mockService.On("CreateBooking", &bookingReq).Return((*models.Booking)(nil), errors.New("flight not found or not enough seats available")).Once()

	jsonValue, _ := json.Marshal(bookingReq)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "flight not found")

	mockService.AssertExpectations(t)
}

// TestCreateBooking_NotEnoughSeats tests booking creation when not enough seats are available
func TestCreateBooking_NotEnoughSeats(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingReq := models.Booking{
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      10, // Requesting more than available (even with oversell)
	}

	mockService.On("CreateBooking", &bookingReq).Return((*models.Booking)(nil), errors.New("not enough seats available (oversell limit reached)")).Once()

	jsonValue, _ := json.Marshal(bookingReq)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not enough seats")

	mockService.AssertExpectations(t)
}

// TestCreateBooking_InternalError tests booking creation when an internal error occurs
func TestCreateBooking_InternalError(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingReq := models.Booking{
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      1,
	}

	mockService.On("CreateBooking", &bookingReq).Return((*models.Booking)(nil), errors.New("database connection error")).Once()

	jsonValue, _ := json.Marshal(bookingReq)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")

	mockService.AssertExpectations(t)
}

// TestGetBooking_Success tests successful retrieval of a booking
func TestGetBooking_Success(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingID := uint(1)
	expectedBooking := models.Booking{
		Model:         gorm.Model{ID: bookingID},
		FlightID:      1,
		PassengerName: "Test User",
		Quantity:      1,
		TotalPrice:    100.0,
		BookingStatus: "Confirmed",
	}

	mockService.On("GetBooking", bookingID).Return(&expectedBooking, nil).Once()

	req, _ := http.NewRequest("GET", "/bookings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBooking models.Booking
	json.Unmarshal(w.Body.Bytes(), &responseBooking)
	assert.Equal(t, expectedBooking.ID, responseBooking.ID)
	assert.Equal(t, expectedBooking.PassengerName, responseBooking.PassengerName)
	assert.Equal(t, expectedBooking.Quantity, responseBooking.Quantity)
	assert.Equal(t, expectedBooking.TotalPrice, responseBooking.TotalPrice)
	assert.Equal(t, expectedBooking.BookingStatus, responseBooking.BookingStatus)

	mockService.AssertExpectations(t)
}

// TestGetBooking_InvalidID tests retrieval with an invalid booking ID format
func TestGetBooking_InvalidID(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	req, _ := http.NewRequest("GET", "/bookings/abc", nil) // Invalid ID
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid booking ID")

	mockService.AssertNotCalled(t, "GetBooking", mock.Anything)
}

// TestGetBooking_NotFound tests retrieval of a non-existent booking
func TestGetBooking_NotFound(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingID := uint(999)
	mockService.On("GetBooking", bookingID).Return((*models.Booking)(nil), errors.New("booking not found")).Once()

	req, _ := http.NewRequest("GET", "/bookings/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "booking not found")

	mockService.AssertExpectations(t)
}

// TestGetBooking_InternalError tests retrieval when an internal error occurs
func TestGetBooking_InternalError(t *testing.T) {
	// Given
	mockService := new(MockBookingService)
	handler := NewBookingHandler(mockService)

	router := setupTestRouter(handler)

	bookingID := uint(1)
	mockService.On("GetBooking", bookingID).Return((*models.Booking)(nil), errors.New("database connection error")).Once()

	req, _ := http.NewRequest("GET", "/bookings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")

	mockService.AssertExpectations(t)
}
