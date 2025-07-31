package repository

import (
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(restriction models.RoomRestrictions) error
	SearchAvaibilityByDates(RoomID int, start, end time.Time) (bool, error)
}
