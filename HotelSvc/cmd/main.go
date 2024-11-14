package main

import (
	"HotelSvc/api/grpc/hotelpb"
	handler "HotelSvc/api/http"
	"HotelSvc/repository/postgresql"
	"HotelSvc/service"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	hotelRepo := postgresql.NewPostgresHotelRepository(db)

	hotelService := service.NewHotelService(hotelRepo)

	roomRepo := postgresql.NewPostgresRoomRepository(db)

	roomService := service.NewRoomService(roomRepo)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := startHTTPServer(hotelService); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Запуск gRPC сервера в отдельной горутине
	go func() {
		defer wg.Done()
		if err := startGRPCServer(roomService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Ожидание завершения серверов
	wg.Wait()

	//if err := startHTTPServer(hotelService); err != nil {
	//	log.Fatalf("Failed to start HTTP server: %v", err)
	//}

	//if err := startGRPCServer(roomService); err != nil {
	//	log.Fatalf("Failed to start GRPC server: %v", err)
	//}
}

// initDB инициализирует подключение к базе данных PostgreSQL
func initDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	return sql.Open("postgres", dsn)
}

// startHTTPServer запускает HTTP сервер для обработки REST-запросов
func startHTTPServer(hotelService *service.HotelService) error {
	mux := http.NewServeMux()
	handler.RegisterHotelRoutes(mux, hotelService)

	addr := os.Getenv("HTTP_ADDR")
	log.Printf("Starting HTTP server on %s...", addr)
	return http.ListenAndServe(addr, mux)
}

type server struct {
	hotelpb.UnimplementedHotelServiceServer
	roomService *service.RoomService
}

func (s *server) GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error) {
	hotelId := req.GetHotelId()
	rooms, err := s.roomService.GetRoomsByHotelId(int(hotelId))
	if err != nil {
		return nil, err
	}

	return &hotelpb.GetRoomsResponse{Rooms: rooms}, nil

}

func startGRPCServer(roomService *service.RoomService) error {
	addr := os.Getenv("GRPC_ADDR")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	s := grpc.NewServer()
	hotelpb.RegisterHotelServiceServer(s, &server{roomService: roomService})

	reflection.Register(s)

	log.Printf("Starting GRPC server on %s", addr)
	return s.Serve(lis)
}
