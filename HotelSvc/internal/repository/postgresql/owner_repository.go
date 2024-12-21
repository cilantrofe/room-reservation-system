package postgresql

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresOwnerRepository struct {
	db *sql.DB
}

func NewPostgresOwnerRepository(db *sql.DB) *PostgresOwnerRepository {
	return &PostgresOwnerRepository{db: db}
}

func (repo *PostgresRoomRepository) GetOwnerIdByHotelId(ctx context.Context, hotelID int) (int, error) {
	var ownerId int
	err := repo.db.QueryRow("SELECT OwnerID from hotels WHERE id = $1 ", hotelID).Scan(&ownerId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("hotel not found")
	}
	return ownerId, err

}
