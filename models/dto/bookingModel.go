package dto

import "time"

type BookingInfo struct {
	ClassName   string `json:"className" validate:"required"`
	UserName    string `json:"userName" validate:"required"`
	BookingDate string `json:"bookingDate"`
}

type Booking struct {
	UserName    string    `json:"userName"`
	BookingDate time.Time `json:"bookingDate"`
}
