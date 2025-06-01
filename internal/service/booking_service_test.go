package service_test

import (
	"glofox/config"
	newError "glofox/errors"
	"glofox/internal/service"
	"glofox/models/dto"

	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMapStore struct {
	mock.Mock
}

func (m *MockMapStore) Load(key string) (value interface{}, ok bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockMapStore) Store(key string, value interface{}) {
	m.Called(key, value)
}
func (m *MockMapStore) Delete(key string) {
	m.Delete(key)
}

func TestInitializeService(t *testing.T) {
	mockMapStore := new(MockMapStore)

	mockMutex := &sync.Mutex{}
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	// Ensure the service is initialized correctly
	assert.NotNil(t, svc)
}

func TestCreateBooking_ValidBooking(t *testing.T) {
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	// Prepare valid booking info
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: time.Now().Format(cfg.DateFormat),
		ClassName:   "YogaClass",
	}

	bookingDate, _ := time.Parse(cfg.DateFormat, bookingInfo.BookingDate)

	// Mock ClassInfo data for a valid class
	classInfo := dto.ClassInfo{
		StartDate:       time.Now().Add(-24 * time.Hour), // Class starts 1 day ago
		EndDate:         time.Now().Add(24 * time.Hour),  // Class ends in 1 day
		AllowedCapacity: 5,
		Bookings:        make(map[time.Time][]string), // Booked slots will be stored in a map of time -> slice of bookings
	}

	// Set up mock MapStore to return classInfo when YogaClass is loaded
	mockMapStore := new(MockMapStore)
	mockMapStore.On("Load", "YogaClass").Return(classInfo, true).Once() // Expect Load to be called once
	mockMapStore.On("Store", "YogaClass", mock.MatchedBy(func(v interface{}) bool {
		// Check that the value being passed to Store is correct
		if bookings, ok := v.(dto.ClassInfo); ok {
			return len(bookings.Bookings[bookingDate]) == 1 // Ensure the booking is correctly added
		}
		return false
	})).Once() // Expect Store to be called once

	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	err := svc.CreateBooking(bookingInfo)

	// Assert that there is no error for valid booking
	assert.NoError(t, err)

	// Assert Store method was called with correct arguments
	mockMapStore.AssertExpectations(t)
}

func TestCreateBooking_InvalidDateFormat(t *testing.T) {
	// Invalid booking date format
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: "invalid_date", // Invalid date format
		ClassName:   "YogaClass",
	}
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	mockMapStore := new(MockMapStore)
	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	err := svc.CreateBooking(bookingInfo)

	// Assert that the error returned is due to the invalid date format
	assert.Error(t, err)
}

func TestCreateBooking_ClassNotExist(t *testing.T) {
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	// Booking info for a class that doesn't exist
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: time.Now().Format(cfg.DateFormat),
		ClassName:   "NonExistentClass", // Class that doesn't exist
	}

	mockMapStore := new(MockMapStore)

	mockMapStore.On("Load", "NonExistentClass").Return(nil, false).Once() // Simulate class not found
	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	err := svc.CreateBooking(bookingInfo)

	// Assert that the error returned is ErrClassNotExist
	assert.Equal(t, err, newError.ErrClassNotExist)

	// Assert Load was called with the correct class name
	mockMapStore.AssertExpectations(t)
}

func TestCreateBooking_BookingDateBeforeClassStart(t *testing.T) {
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	// Booking date is before the class start date
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: time.Now().Format(cfg.DateFormat), // Booking date before class start date
		ClassName:   "YogaClass",
	}

	// Class starts in 2 days, ends in 3 days
	classInfo := dto.ClassInfo{
		StartDate:       time.Now().Add(48 * time.Hour), // Class starts in 2 days
		EndDate:         time.Now().Add(72 * time.Hour), // Class ends in 3 days
		AllowedCapacity: 5,
		Bookings:        make(map[time.Time][]string),
	}

	// Set up mock map store
	mockMapStore := new(MockMapStore)
	mockMapStore.On("Load", "YogaClass").Return(classInfo, true).Once() // Expect the class to be loaded

	// Mutex and service setup
	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	// Run the service method
	err := svc.CreateBooking(bookingInfo)

	// Assert that the error returned is ErrBookingDatePassed (booking before class starts)
	assert.Equal(t, err, newError.ErrBookingDatePassed)

	// Assert expectations on mock MapStore
	mockMapStore.AssertExpectations(t)
}

func TestCreateBooking_BookingDateAfterClassEnd(t *testing.T) {
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	// Booking date is after the class end date
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: time.Now().Add(48 * time.Hour).Format(cfg.DateFormat), // Booking date after class end date
		ClassName:   "YogaClass",
	}

	classInfo := dto.ClassInfo{
		StartDate:       time.Now().Add(-48 * time.Hour), // Class starts 2 days ago
		EndDate:         time.Now().Add(24 * time.Hour),  // Class ends in 1 day
		AllowedCapacity: 5,
		Bookings:        make(map[time.Time][]string),
	}

	mockMapStore := new(MockMapStore)
	mockMapStore.On("Load", "YogaClass").Return(classInfo, true).Once()

	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)

	err := svc.CreateBooking(bookingInfo)

	// Assert that the error returned is ErrBookingDatePassed (booking after class ends)
	assert.Equal(t, err, newError.ErrBookingDatePassed)

	mockMapStore.AssertExpectations(t)
}

func TestCreateBooking_SlotsFullForTheDate(t *testing.T) {
	cfg := config.Config{
		DateFormat: "2006-01-02",
	}
	// Test when slots are full for the given date
	bookingInfo := dto.BookingInfo{
		UserName:    "john_doe",
		BookingDate: time.Now().Format(cfg.DateFormat), // Valid date format
		ClassName:   "YogaClass",
	}

	bookingDate, _ := time.Parse(cfg.DateFormat, bookingInfo.BookingDate)
	timeNow := time.Now()
	date := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC)
	// Class with one slot, and already full for the day
	classInfo := dto.ClassInfo{
		StartDate:       date.Add(-24 * time.Hour),
		EndDate:         time.Now().Add(24 * time.Hour),
		AllowedCapacity: 1,
		Bookings:        make(map[time.Time][]string),
	}
	// Pre-add a booking on the date
	classInfo.Bookings[bookingDate] = append(classInfo.Bookings[bookingDate], "existing_user")

	mockMapStore := new(MockMapStore)
	mockMapStore.On("Load", "YogaClass").Return(classInfo, true).Once()

	mockMutex := &sync.Mutex{}
	svc := service.InitializeService(mockMapStore, mockMutex, cfg)
	err := svc.CreateBooking(bookingInfo)

	// Assert that the error returned is ErrSlotsFullForTheDate
	assert.Equal(t, err, newError.ErrSlotsFullForTheDate)

	mockMapStore.AssertExpectations(t)
}
