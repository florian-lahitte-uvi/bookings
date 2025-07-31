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
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	// Give a context with a timeout
	// This is a good practice to avoid long-running queries
	// that could block the database connection pool
	// and lead to performance issues.
	// Here we set a timeout of 3 seconds.
	// Adjust the timeout as needed for your application.
	// If the operation takes longer than this, it will be cancelled.
	// This is especially useful in web applications where you want to
	// avoid hanging requests.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL statement to insert a reservation
	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
