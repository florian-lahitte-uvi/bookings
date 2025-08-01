package dbrepo

import (
	"context"
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// Give a context with a timeout
	// This is a good practice to avoid long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	// Prepare the SQL statement to insert a reservation
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(restriction models.RoomRestrictions) error {
	// Give a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to insert a room restriction
	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt, restriction.StartDate, restriction.EndDate, restriction.RoomID, restriction.ReservationID, restriction.RestrictionID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// return true if there are no room restrictions for the given dates
// return false if there are room restrictions for the given dates
func (m *postgresDBRepo) SearchAvaibilityByDatesByRoomID(RoomID int, start, end time.Time) (bool, error) {

	// Give a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	// Prepare the SQL statement to search for available rooms
	query := `select 
				count(id) 
			from 
				room_restrictions 
			where
				room_id = $1 
			and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, RoomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvaibilityForAllRooms returns a slice of available rooms for the given dates
// It returns an empty slice if there are no available rooms
func (m *postgresDBRepo) SearchAvaibilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// Give a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to search for available rooms
	query := `select r.id, r.room_name
	 from 
	 	rooms r 
	where r.id not in 
	(select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date) `

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// get RoomByID returns a room
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	// Give a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	// Prepare the SQL statement to get a room by ID
	stmt := `select id, room_name, created_at, updated_at from rooms where id = $1`
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return models.Room{}, err
	}
	return room, nil
}
