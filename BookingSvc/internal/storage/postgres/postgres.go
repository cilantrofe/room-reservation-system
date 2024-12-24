package postgres

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateBooking(ctx context.Context, booking *models.BookingInfo) (int, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		status = "failed"
		return -1, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT RoomID
		FROM bookings
		where RoomID = $1
		and ($2 >= StartDate and $2 < EndDate) or
			($2 <= StartDate and $3 > StartDate);
	`
	rows, err := tx.Query(ctx, query, booking.RoomID, booking.StartDate, booking.EndDate)
	if err != nil {
		status = "failed"
		return -1, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roomID int
		if err := rows.Scan(&roomID); err != nil {
			return -1, fmt.Errorf("failed to scan room ID: %w", err)
		}
		if roomID == booking.RoomID {
			status = "failed"
			return -1, fmt.Errorf("in storage CreateBooking: %w", myerror.ErrBookingAlreadyExists)
		}
	}

	query = `
        INSERT INTO bookings (UserID, RoomID, HotelID, StartDate, EndDate, CreatedAt)
        VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id

    `
	var bookingID int

	err = tx.QueryRow(ctx, query, booking.UserID, booking.RoomID, booking.HotelID, booking.StartDate, booking.EndDate).Scan(&bookingID)
	if err != nil {
		status = "failed"
		return -1, fmt.Errorf("failed to create booking: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		status = "failed"
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return bookingID, nil
}

func (r *Repository) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()
	bookings := make([]*models.BookingInfo, 0)
	query := `
        SELECT UserID, RoomID, HotelID, StartDate, EndDate
        FROM bookings
        WHERE UserID = $1
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		status = "failed"
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var booking models.BookingInfo
	for rows.Next() {
		if err := rows.Scan(&booking.UserID, &booking.RoomID, &booking.HotelID, &booking.StartDate, &booking.EndDate); err != nil {
			status = "failed"
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		status = "failed"
		return nil, fmt.Errorf("rows iteration myerror: %w", err)
	}
	return bookings, nil
}

func (r *Repository) GetUnavailableRoomsByHotelId(ctx context.Context, hotelID int, startDate, endDate time.Time) (map[int]struct{}, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	unavailableRoomsID := make(map[int]struct{})
	query := `
		SELECT RoomID
		from bookings
		where HotelID = $1
		and ($2 >= StartDate and $2 < EndDate) or
			($2 <= StartDate and $3 > StartDate);
    `
	rows, err := r.db.Query(ctx, query, hotelID, startDate, endDate)
	if err != nil {
		status = "failed"
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roomID int
		if err := rows.Scan(&roomID); err != nil {
			status = "failed"
			return nil, fmt.Errorf("failed to scan room ID: %w", err)
		}
		unavailableRoomsID[roomID] = struct{}{}
	}
	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		status = "failed"
		return nil, fmt.Errorf("rows iteration myerror: %w", err)
	}
	return unavailableRoomsID, nil
}

func (r *Repository) GetBookingsByHotelID(ctx context.Context, hotelID int) ([]*models.BookingInfo, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()
	bookings := make([]*models.BookingInfo, 0)
	query := `
        SELECT UserID, RoomID, HotelID, StartDate, EndDate
        FROM bookings
        WHERE hotelID = $1
    `
	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		status = "failed"
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.BookingInfo
		if err = rows.Scan(&booking.UserID, &booking.RoomID, &booking.HotelID, &booking.StartDate, &booking.EndDate); err != nil {
			status = "failed"
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}
	if err = rows.Err(); err != nil {
		status = "failed"
		return nil, fmt.Errorf("rows iteration myerror: %w", err)
	}
	return bookings, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context, status string, bookingID int) error {
	start := time.Now()
	statusMetrics := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", statusMetrics, duration)
	}()
	query := `
		UPDATE bookings
		SET status = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, bookingID)
	if err != nil {
		statusMetrics = "failed"
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	return nil
}
