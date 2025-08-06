package repository

import (
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(restriction models.RoomRestrictions) error
	SearchAvaibilityByDatesByRoomID(RoomID int, start, end time.Time) (bool, error)
	SearchAvaibilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
}
