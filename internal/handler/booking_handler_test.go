package handler

import (
	"bytes"
	"glofox/constants"
	newError "glofox/errors"
	"glofox/models/dto"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the BusinessService
type MockBusinessService struct {
	mock.Mock
}

func (m *MockBusinessService) CreateBooking(bookingInfo dto.BookingInfo) error {
	args := m.Called(bookingInfo)
	return args.Error(0)
}
func (m *MockBusinessService) CreateClass(classData dto.Class) error {
	args := m.Called(classData)
	return args.Error(0)
}

// Mocking the MapStore (No actual behavior needed for this test)
type MockMapStore struct {
	mock.Mock
}

func (m *MockMapStore) Store(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockMapStore) Load(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}
func (m *MockMapStore) Delete(key string) {
	m.Delete(key)
}

func (m *MockMapStore) PrintMap() {
	m.Called()
}

// Test cases for CreateBooking handler
func TestCreateBooking_ValidInput(t *testing.T) {
	// Prepare mock service
	mockService := new(MockBusinessService)
	mockService.On("CreateBooking", mock.AnythingOfType("dto.BookingInfo")).Return(nil).Once()

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewBookingHandler(mockMapStore, mockMutex, mockService)

	// Create a new Gin context with the booking info as the body
	body := `{"UserName":"john_doe","BookingDate":"2025-05-10","ClassName":"YogaClass"}`
	w := performRequestBookingHandler("POST", "/booking", body, handler)

	// Check the response code and message
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), constants.BookingSucces)

	// Assert that CreateBooking was called once with the correct argument
	mockService.AssertExpectations(t)
}

func TestCreateBooking_InvalidPayload(t *testing.T) {
	// Prepare mock service (won't be called in this case)
	mockService := new(MockBusinessService)

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewBookingHandler(mockMapStore, mockMutex, mockService)

	// Create an invalid booking info (malformed JSON)
	body := `{"UserName": "john_doe", "BookingDate": "2025-05-10"`

	// Create a new Gin context with the invalid body
	w := performRequestBookingHandler("POST", "/booking", body, handler)

	// Check that the response code is BadRequest (400) and contains the error message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), newError.ErrUnmarshalling.Error())
}

func TestCreateBooking_ServiceError(t *testing.T) {
	// Prepare mock service
	mockService := new(MockBusinessService)
	mockService.On("CreateBooking", mock.AnythingOfType("dto.BookingInfo")).Return(newError.ErrBookingDatePassed).Once()

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewBookingHandler(mockMapStore, mockMutex, mockService)

	// Create a valid booking info (request body)
	body := `{"UserName":"john_doe","BookingDate":"2025-05-10","ClassName":"YogaClass"}`
	w := performRequestBookingHandler("POST", "/booking", body, handler)

	// Check the response code and message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), newError.ErrBookingDatePassed.Error())

	// Assert that CreateBooking was called once with the correct argument
	mockService.AssertExpectations(t)
}

func performRequestBookingHandler(method, url, body string, handler BookingHandler) *httptest.ResponseRecorder {
	// Set up the Gin router and add the handler
	r := gin.Default()
	r.POST(url, handler.CreateBooking)

	// Create a new HTTP request
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
