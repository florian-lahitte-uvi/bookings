package dbrepo

import (
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(restriction models.RoomRestrictions) error {
	return nil
}

// return true if there are no room restrictions for the given dates
// return false if there are room restrictions for the given dates
func (m *testDBRepo) SearchAvaibilityByDatesByRoomID(RoomID int, start, end time.Time) (bool, error) {
	return false, nil
}

// SearchAvaibilityForAllRooms returns a slice of available rooms for the given dates
// It returns an empty slice if there are no available rooms
func (m *testDBRepo) SearchAvaibilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// get RoomByID returns a room
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
