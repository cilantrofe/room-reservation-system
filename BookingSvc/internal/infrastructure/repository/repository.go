package repository

import (
	"BookingSvc/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	query := `
        SELECT id, user_id, room_id, hotel_id, start_date, end_date, status
        FROM bookings
        WHERE user_id = $1
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.HotelID, &booking.StartDate, &booking.EndDate, &booking.Status); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookings, nil
}

func (r *PostgresRepository) GetBooking(ctx context.Context, bookingID int) (*models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresRepository) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresRepository) DeleteBooking(ctx context.Context, bookingID int) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	query := `
        INSERT INTO bookings (user_id, room_id, hotel_id, start_date, end_date, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING id
    `

	var bookingID int
	err := r.db.QueryRow(ctx, query, booking.UserID, booking.RoomID, booking.HotelID, booking.StartDate, booking.EndDate, booking.Status).Scan(&bookingID)
	if err != nil {
		log.Fatalf("failed to create booking: %v", err)
	}
	return nil
}
