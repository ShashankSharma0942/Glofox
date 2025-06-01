package handler

import (
	"glofox/constants"
	mapstore "glofox/core"
	newError "glofox/errors"
	"glofox/internal/service"
	"glofox/models/dto"
	"glofox/utils"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// BookingHandler defines the interface for handling booking-related HTTP requests.
type BookingHandler interface {
	CreateBooking(c *gin.Context)
}

// booking is the concrete implementation of BookingHandler.
// It holds the shared state, mutex, and business service required to process booking requests.
type booking struct {
	syMap   mapstore.MapStore
	lock    *sync.Mutex
	service service.BusinessService
}

// NewBookingHandler constructs and returns a new BookingHandler with injected dependencies.
func NewBookingHandler(syMap mapstore.MapStore, lock *sync.Mutex, services service.BusinessService) BookingHandler {
	return &booking{
		syMap:   syMap,
		lock:    lock,
		service: services,
	}
}

// CreateBooking handles the POST /booking endpoint.
// It validates and binds the request payload, delegates business logic to the service layer,
// and returns a structured JSON response.
func (booking *booking) CreateBooking(c *gin.Context) {
	var bookingInfo dto.BookingInfo

	// Attempt to bind the incoming JSON payload to the BookingInfo struct
	err := c.ShouldBindJSON(&bookingInfo)
	if err != nil {
		log.Println(newError.ErrUnmarshalling.Error(), err.Error())
		c.JSON(http.StatusBadRequest, utils.CreateResp(false, newError.ErrUnmarshalling.Error()))
		return
	}

	// Call the service layer to process the booking
	err = booking.service.CreateBooking(bookingInfo)
	if err != nil {
		log.Println(newError.ErrCreatingBooking.Error(), err.Error())
		c.JSON(http.StatusBadRequest, utils.CreateResp(false, err.Error()))
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, utils.CreateResp(true, constants.BookingSucces))
}
