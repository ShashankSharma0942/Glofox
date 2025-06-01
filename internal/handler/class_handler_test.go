package handler

import (
	"bytes"
	"glofox/constants"
	newError "glofox/errors"
	"net/http"
	"sync"
	"testing"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test cases for CreateClass handler
func TestCreateClass_ValidInput(t *testing.T) {
	// Prepare mock service
	mockService := new(MockBusinessService)
	mockService.On("CreateClass", mock.AnythingOfType("dto.Class")).Return(nil).Once()

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewClassHandler(mockMapStore, mockMutex, mockService)

	// Create a new Gin context with the class data as the body
	body := `{"Name":"Yoga Class","Capacity":30,"StartDate":"2025-06-01","EndDate":"2025-06-10"}`
	w := performRequest("POST", "/class", body, handler)

	// Check the response code and message
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), constants.ClassSuccess)

	// Assert that CreateClass was called once with the correct argument
	mockService.AssertExpectations(t)
}

func TestCreateClass_InvalidPayload(t *testing.T) {
	// Prepare mock service (won't be called in this case)
	mockService := new(MockBusinessService)

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewClassHandler(mockMapStore, mockMutex, mockService)

	// Create an invalid class info (malformed JSON)
	body := `{"Name": "Yoga Class", "Capacity": 30, "StartDate": "2025-06-01"`

	// Create a new Gin context with the invalid body
	w := performRequest("POST", "/class", body, handler)

	// Check that the response code is BadRequest (400) and contains the error message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), newError.ErrUnmarshalling.Error())
}

func TestCreateClass_ServiceError(t *testing.T) {
	// Prepare mock service
	mockService := new(MockBusinessService)
	mockService.On("CreateClass", mock.AnythingOfType("dto.Class")).Return(newError.ErrEndTimeLessThanStartTime).Once()

	// Prepare mock MapStore and lock
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}

	// Initialize the handler with mock dependencies
	handler := NewClassHandler(mockMapStore, mockMutex, mockService)

	// Create a valid class info (request body)
	body := `{"Name":"Yoga Class","Capacity":30,"StartDate":"2025-06-01","EndDate":"2025-05-10"}`
	w := performRequest("POST", "/class", body, handler)

	// Check the response code and message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), newError.ErrEndTimeLessThanStartTime.Error())

	// Assert that CreateClass was called once with the correct argument
	mockService.AssertExpectations(t)
}

func performRequest(method, url, body string, handler ClassHandler) *httptest.ResponseRecorder {
	// Set up the Gin router and add the handler
	r := gin.Default()
	r.POST(url, handler.CreateClass)

	// Create a new HTTP request
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
