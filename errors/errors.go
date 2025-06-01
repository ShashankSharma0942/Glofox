package newError

import (
	"errors"
	"fmt"
)

var (
	ErrUnmarshalling            = errors.New("error while unamrshalling")
	ErrCreatingBooking          = errors.New("Error while creating booking:")
	ErrClassNotExist            = errors.New(fmt.Sprintf("Please Check Your Class Name"))
	ErrBookingDatePassed        = errors.New("booking for the mentioned date is not allowed for the class")
	ErrSlotsFullForTheDate      = errors.New("booking full for the requested class on the mentioned date")
	ErrEndTimeLessThanStartTime = errors.New("class end date can not be less than start end date")
)
