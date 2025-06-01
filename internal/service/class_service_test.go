package service

import (
	"glofox/config"
	newError "glofox/errors"
	"glofox/models/dto"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the MapStore interface
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

func TestCreateClass_Success(t *testing.T) {
	// Create a valid class info
	classInfo := dto.Class{
		Name:      "Yoga Class",
		Capacity:  30,
		StartDate: "2025-06-01",
		EndDate:   "2025-06-10",
	}

	mockMapStore := new(MockMapStore)
	mockMapStore.On("Store", "Yoga Class", mock.Anything).Once()
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	mockMutex := &sync.Mutex{}
	svc := InitializeService(mockMapStore, mockMutex, cfg)

	// Call the CreateClass method
	err := svc.CreateClass(classInfo)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the Store method was called with the correct arguments
	mockMapStore.AssertExpectations(t)
}

func TestCreateClass_InvalidStartDate(t *testing.T) {
	// Create a class info with an invalid start date format
	classInfo := dto.Class{
		Name:      "Yoga Class",
		Capacity:  30,
		StartDate: "2025-06-01T00:00:00", // Invalid format
		EndDate:   "2025-06-10",
	}

	mockMapStore := new(MockMapStore)
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	mockMutex := &sync.Mutex{}
	svc := InitializeService(mockMapStore, mockMutex, cfg)

	// Call the CreateClass method
	err := svc.CreateClass(classInfo)

	// Assert that an error is returned (invalid date format)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "parsing time \"2025-06-01T00:00:00\": extra text: \"T00:00:00\"")
}

func TestCreateClass_InvalidEndDateBeforeStartDate(t *testing.T) {
	// Create a class info where the end date is before the start date
	classInfo := dto.Class{
		Name:      "Yoga Class",
		Capacity:  30,
		StartDate: "2025-06-10",
		EndDate:   "2025-06-01", // End date before start date
	}

	mockMapStore := new(MockMapStore)
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	mockMutex := &sync.Mutex{}
	svc := InitializeService(mockMapStore, mockMutex, cfg)

	// Call the CreateClass method
	err := svc.CreateClass(classInfo)

	// Assert that the correct error is returned (end date before start date)
	assert.Error(t, err)
	assert.Equal(t, err, newError.ErrEndTimeLessThanStartTime)
}

func TestCreateClass_InvalidStartDateAfterEndDate(t *testing.T) {
	// Create a class info where the end date is before the start date
	classInfo := dto.Class{
		Name:      "Yoga Class",
		Capacity:  30,
		StartDate: "2025-06-01",
		EndDate:   "2025-06-10",
	}

	mockMapStore := new(MockMapStore)
	mockMapStore.On("Store", "Yoga Class", mock.Anything).Once()
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	mockMutex := &sync.Mutex{}
	svc := InitializeService(mockMapStore, mockMutex, cfg)

	// Call the CreateClass method
	err := svc.CreateClass(classInfo)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the Store method was called with the correct arguments
	mockMapStore.AssertExpectations(t)
}
