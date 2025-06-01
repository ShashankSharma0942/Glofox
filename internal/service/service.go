package service

import (
	"glofox/config"
	mapstore "glofox/core"
	"glofox/models/dto"
	"sync"
)

// service is the concrete implementation of BusinessService interface.
// It holds a thread-safe map store (syMap) and a mutex lock to manage concurrent access.
type service struct {
	syMap mapstore.MapStore
	lock  *sync.Mutex
	cfg   config.Config
}

// BusinessService defines the business logic interface for class and booking operations.
type BusinessService interface {
	CreateClass(info dto.Class) error
	CreateBooking(bookingInfo dto.BookingInfo) error
}

// InitializeService creates and returns a new instance of BusinessService
// injecting the shared map store and mutex for thread-safe operations.
func InitializeService(syMap mapstore.MapStore, mu *sync.Mutex, cfg config.Config) BusinessService {
	return &service{
		syMap: syMap,
		lock:  mu,
		cfg:   cfg,
	}
}
