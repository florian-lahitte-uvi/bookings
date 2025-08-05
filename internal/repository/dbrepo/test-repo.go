package dbrepo

import (
	"errors"
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("error inserting reservation")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(restriction models.RoomRestrictions) error {
	if restriction.RoomID == 1000 {
		return errors.New("error inserting room restriction")
	}
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
	if id > 2 {
		return room, errors.New("room not found")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	if id > 2 {
		return u, errors.New("User not found")
	}
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	if u.ID == 0 {
		return errors.New("User not found")
	}
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "test@example.com" && testPassword == "password" {
		return 1, "hashedPassword", nil
	}
	return 0, "", errors.New("invalid credentials")
}
