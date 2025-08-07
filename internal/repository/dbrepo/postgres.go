package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// Give a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	// Prepare the SQL statement to get a user by ID
	stmt := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to update a user
	stmt := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`
	_, err := m.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	// Prepare the SQL statement to get the user by email
	stmt := `select id, password from users where email = $1`
	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("invalid password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// Return a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.room_id, r.email, r.first_name, r.last_name, r.phone, r.start_date, r.end_date, r.created_at, r.updated_at, r.processed, rm.room_name
	from reservations r 
	left join rooms rm on (r.room_id = rm.id) 
	order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.RoomID, &r.Email, &r.FirstName, &r.LastName, &r.Phone, &r.StartDate, &r.EndDate, &r.CreatedAt, &r.UpdatedAt, &r.Processed, &r.Room.RoomName)
		if err != nil {
			return reservations, err
		}
		r.Room.ID = r.RoomID // Set the room ID manually
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.room_id, r.email, r.first_name, r.last_name, r.phone, r.start_date, r.end_date, r.created_at, r.updated_at, r.processed, rm.room_name
	from reservations r 
	left join rooms rm on (r.room_id = rm.id)
	where r.processed = 0
	order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.RoomID, &r.Email, &r.FirstName, &r.LastName, &r.Phone, &r.StartDate, &r.EndDate, &r.CreatedAt, &r.UpdatedAt, &r.Processed, &r.Room.RoomName)
		if err != nil {
			return reservations, err
		}
		r.Room.ID = r.RoomID // Set the room ID manually
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Return a reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `select r.id, r.room_id, r.email, r.first_name, r.last_name, r.phone, r.start_date, r.end_date, r.created_at, r.updated_at, r.processed, rm.room_name, rm.id
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where r.id = $1`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.RoomID, &res.Email, &res.FirstName, &res.LastName, &res.Phone, &res.StartDate, &res.EndDate, &res.CreatedAt, &res.UpdatedAt, &res.Processed, &res.Room.RoomName, &res.Room.ID)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Update reservations
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to update a Reservation
	stmt := `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5`
	_, err := m.DB.ExecContext(ctx, stmt, r.FirstName, r.LastName, r.Email, r.Phone, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Delete reservations by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to update a Reservation
	stmt := `delete from reservations where id = $1 `
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// update processed by reservation id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to update a Reservation
	stmt := `update reservations set processed = $1 where id = $2 `
	_, err := m.DB.ExecContext(ctx, stmt, processed, id)
	if err != nil {
		return err
	}
	return nil
}
