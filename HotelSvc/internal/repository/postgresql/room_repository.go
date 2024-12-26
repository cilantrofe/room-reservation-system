package postgresql

import (
	"database/sql"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/models"
)

type PostgresRoomRepository struct {
	db *sql.DB
}

func NewPostgresRoomRepository(db *sql.DB) *PostgresRoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (repo *PostgresRoomRepository) GetRoomsByHotelId(id int) ([]*hotelpb.Room, error) {
	rows, err := repo.db.Query(
		`SELECT 
        r.id AS RoomId, 
        r.HotelId, 
        r.Number, 
        rt.Description, 
        rt.BasePrice AS Cost 
     FROM 
        rooms r 
     JOIN 
        room_type rt 
     ON 
        r.RoomTypeId = rt.id 
     WHERE 
        r.HotelId = $1;`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*hotelpb.Room
	for rows.Next() {
		var room hotelpb.Room
		if err := rows.Scan(&room.Id, &room.HotelId, &room.Number, &room.Description, &room.BasePrice); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, rows.Err()

}

func (repo *PostgresRoomRepository) AddRoom(room models.Room) error {
	_, err := repo.db.Exec(
		"INSERT INTO rooms (HotelId, RoomTypeId, Number) VALUES ($1, $2, $3)",
		room.HotelID, room.RoomTypeID, room.Number,
	)
	return err
}

func (repo *PostgresRoomRepository) AddRoomType(roomType models.RoomType) error {
	_, err := repo.db.Exec(
		"INSERT INTO room_type (Name, Description, BasePrice) VALUES ($1, $2, $3)",
		roomType.Name, roomType.Description, roomType.BasePrice,
	)
	return err
}
