package handler

import (
	"encoding/json"
	"errors"
	"flight-booking/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockFlightRepository is a mock implementation of FlightRepository interface
type MockFlightRepository struct {
	mock.Mock
}

func (m *MockFlightRepository) FindAll(query *gorm.DB, page, pageSize int) ([]models.Flight, int64, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).([]models.Flight), args.Get(1).(int64), args.Error(2)
}

func (m *MockFlightRepository) FindByID(id uint) (*models.Flight, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Flight), args.Error(1)
}

func (m *MockFlightRepository) Create(flight *models.Flight) error {
	args := m.Called(flight)
	return args.Error(0)
}

func (m *MockFlightRepository) Update(flight *models.Flight) error {
	args := m.Called(flight)
	return args.Error(0)
}

// SetupRouter for testing
func setupFlightTestRouter(flightHandler *FlightHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/flights", flightHandler.SearchFlights)
	r.GET("/flights/:id", flightHandler.GetFlight)
	return r
}

func successFindAll(query *gorm.DB, page, pageSize int) ([]models.Flight, int64, error) {
	return []models.Flight{
		{
			Model:            gorm.Model{ID: 1},
			DepartureAirport: "Taipei",
			ArrivalAirport:   "Tokyo",
			DepartureTime:    "2025-08-01 10:00",
			ArrivalTime:      "2025-08-01 14:00",
			Airline:          "EVA Air",
			FlightNumber:     "BR101",
			Price:            500,
			AvailableSeats:   100,
		},
	}, 1, nil
}

// TestSearchFlights_Success tests a successful flight search
func TestSearchFlights_Success(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewFlightHandler(mockRepo, db) // Pass a dummy DB instance for handler's internal query building

	router := setupFlightTestRouter(handler)

	expectedFlights := []models.Flight{
		{
			Model:            gorm.Model{ID: 1},
			DepartureAirport: "Taipei",
			ArrivalAirport:   "Tokyo",
			DepartureTime:    "2025-08-01 10:00",
			ArrivalTime:      "2025-08-01 14:00",
			Airline:          "EVA Air",
			FlightNumber:     "BR101",
			Price:            500,
			AvailableSeats:   100,
		},
	}
	expectedTotal := int64(1)

	// Mock the FindAll method. The first argument (query *gorm.DB) is hard to match precisely,
	mockRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(successFindAll(db, 1, 10)).Once()

	req, _ := http.NewRequest("GET", "/flights?departure=Taipei&date=2025-08-01&page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	var response SearchFlightsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedTotal, response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, expectedFlights[0].ID, response.Data[0].ID)
	assert.Equal(t, expectedFlights[0].DepartureAirport, response.Data[0].DepartureAirport)
	assert.Equal(t, expectedFlights[0].Price, response.Data[0].Price)

	mockRepo.AssertExpectations(t)
}

// TestSearchFlights_InvalidDate tests flight search with an invalid date format
func TestSearchFlights_InvalidDate(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	req, _ := http.NewRequest("GET", "/flights?date=invalid-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid date format")

	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything)
}

// TestSearchFlights_InvalidPageParams tests flight search with invalid page or page_size
func TestSearchFlights_InvalidPageParams(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	// Test invalid page
	req, _ := http.NewRequest("GET", "/flights?page=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid page parameter")

	// Test invalid page_size
	req, _ = http.NewRequest("GET", "/flights?page_size=xyz", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid page_size parameter")

	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything)
}

// TestSearchFlights_InternalError tests flight search when an internal error occurs in repository
func TestSearchFlights_InternalError(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	mockRepo.On("FindAll", mock.Anything, 1, 10).Return(([]models.Flight)(nil), int64(0), errors.New("database error")).Once()

	req, _ := http.NewRequest("GET", "/flights?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")

	mockRepo.AssertExpectations(t)
}

// TestGetFlight_Success tests successful retrieval of a flight
func TestGetFlight_Success(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	flightID := uint(1)
	expectedFlight := &models.Flight{
		Model:            gorm.Model{ID: flightID},
		FlightNumber:     "BR101",
		DepartureAirport: "Taipei",
		ArrivalAirport:   "Tokyo",
		Price:            500.00,
		AvailableSeats:   100,
	}

	mockRepo.On("FindByID", flightID).Return(expectedFlight, nil).Once()

	req, _ := http.NewRequest("GET", "/flights/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	var responseFlight models.Flight
	json.Unmarshal(w.Body.Bytes(), &responseFlight)
	assert.Equal(t, expectedFlight.ID, responseFlight.ID)
	assert.Equal(t, expectedFlight.FlightNumber, responseFlight.FlightNumber)

	mockRepo.AssertExpectations(t)
}

// TestGetFlight_InvalidID tests retrieval with an invalid flight ID format
func TestGetFlight_InvalidID(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	req, _ := http.NewRequest("GET", "/flights/abc", nil) // Invalid ID
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid flight ID")

	mockRepo.AssertNotCalled(t, "FindByID", mock.Anything)
}

// TestGetFlight_NotFound tests retrieval of a non-existent flight
func TestGetFlight_NotFound(t *testing.T) {
	// Given
	mockRepo := new(MockFlightRepository)
	handler := NewFlightHandler(mockRepo, &gorm.DB{})

	router := setupFlightTestRouter(handler)

	flightID := uint(999)
	mockRepo.On("FindByID", flightID).Return((*models.Flight)(nil), gorm.ErrRecordNotFound).Once()

	req, _ := http.NewRequest("GET", "/flights/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Flight not found")

	mockRepo.AssertExpectations(t)
}
