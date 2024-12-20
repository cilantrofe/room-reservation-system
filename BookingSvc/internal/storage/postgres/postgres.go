package postgres

import (
	"context"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
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

func (r *Repository) GetUnavailableRoomsByHotelId(ctx context.Context, hotelID int, startDate, endDate time.Time) (map[int]struct{}, error) {
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
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roomID int
		if err := rows.Scan(&roomID); err != nil {
			return nil, fmt.Errorf("failed to scan room ID: %w", err)
		}
		unavailableRoomsID[roomID] = struct{}{}
	}
	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return unavailableRoomsID, nil
}

func (r *Repository) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error) {
	bookings := make([]*models.BookingInfo, 0)
	query := `
        SELECT UserID, RoomID, HotelID, StartDate, EndDate
        FROM bookings
        WHERE UserID = $1
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var booking models.BookingInfo
	for rows.Next() {
		if err := rows.Scan(&booking.UserID, &booking.RoomID, &booking.HotelID, &booking.StartDate, &booking.EndDate); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return bookings, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context, status string, bookingID int) error {
	query := `
		UPDATE bookings
		SET status = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, bookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	return nil
}

func (r *Repository) GetBookingsByHotelID(ctx context.Context, bookingID int) (*models.BookingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) DeleteBooking(ctx context.Context, bookingID int) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) CreateBooking(ctx context.Context, booking *models.BookingInfo) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
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
		return -1, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roomID int
		if err := rows.Scan(&roomID); err != nil {
			return -1, fmt.Errorf("failed to scan room ID: %w", err)
		}
		if roomID == booking.RoomID {
			return -1, fmt.Errorf("booking already exists")
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
		return -1, fmt.Errorf("failed to create booking: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return bookingID, nil
}
