package dto

import "time"

type Class struct {
	Name      string `json:"className" validate:"required"`
	Capacity  int    `json:"classCapacity" validate:"required"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}

type ClassInfo struct {
	AllowedCapacity int                    `json:"allowedCapacity"`
	StartDate       time.Time              `json:"classStartDt"`
	EndDate         time.Time              `json:"classEndDt"`
	Bookings        map[time.Time][]string `json:"bookings"`
}
