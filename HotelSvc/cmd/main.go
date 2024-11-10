package main

import (
	handler "HotelSvc/api/http"
	"HotelSvc/repository/postgresql"
	"HotelSvc/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	if err := startHTTPServer(hotelService); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

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
