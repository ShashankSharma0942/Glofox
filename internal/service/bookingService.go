package service

import (
	newError "glofox/errors"
	"glofox/models/dto"
	"time"
)

// CreateBooking handles booking a user into a class on a specific date.
// It performs several checks including class existence, booking window validity,
// capacity limits, and then stores the booking in the system.
func (service *service) CreateBooking(bookingInfo dto.BookingInfo) error {

	bookingDate, err := time.Parse(service.cfg.DateFormat, bookingInfo.BookingDate)
	if err != nil {
		return err
	}
	service.lock.Lock()
	defer service.lock.Unlock()

	classDate, exist := service.syMap.Load(bookingInfo.ClassName)
	if !exist {
		err = newError.ErrClassNotExist
		return err
	}

	// Type assert the loaded value to ClassInfo type
	typeCastData := classDate.(dto.ClassInfo)

	if bookingDate.Before(typeCastData.StartDate) {
		return newError.ErrBookingDatePassed
	}
	if bookingDate.After(typeCastData.EndDate) {
		return newError.ErrBookingDatePassed
	}
	if len(typeCastData.Bookings[bookingDate]) >= typeCastData.AllowedCapacity {
		return newError.ErrSlotsFullForTheDate
	}

	typeCastData.Bookings[bookingDate] = append(typeCastData.Bookings[bookingDate], bookingInfo.UserName)

	service.syMap.Store(bookingInfo.ClassName, typeCastData)

	return err
}
