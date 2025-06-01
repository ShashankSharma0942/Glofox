package service

import (
	newError "glofox/errors"
	"glofox/models/dto"
	"time"
)

// CreateClass processes the creation of a new class based on the provided input data.
// It validates date formats, ensures logical consistency of start and end dates,
// initializes the class information structure, and stores it in the shared map.
func (service *service) CreateClass(info dto.Class) error {
	//Time object
	startDate, err := time.Parse(service.cfg.DateFormat, info.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse(service.cfg.DateFormat, info.EndDate)
	if err != nil {
		return err
	}

	hrs := endDate.Sub(startDate).Hours()
	if hrs < 0 {
		return newError.ErrEndTimeLessThanStartTime
	}

	classInfo := dto.ClassInfo{
		AllowedCapacity: info.Capacity,

		EndDate:  endDate.Truncate(24 * time.Hour),
		Bookings: make(map[time.Time][]string),
	}

	service.lock.Lock()
	defer service.lock.Unlock()

	service.syMap.Store(info.Name, classInfo)

	return nil
}
