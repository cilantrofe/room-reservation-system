package app

import (
	"BookingSvc/internal/controller/handler"
	"BookingSvc/internal/infrastructure/repository"
	"BookingSvc/internal/service"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"time"
)

type App struct {
	dbPool  *pgxpool.Pool
	service *service.BookingService
	server  *http.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context) error {
	//инициализация grpc, handler, роутинг, адаптеров, репозиториев, кафка, коннекторов к другим микросервисам,
	connString := "postgres://postgres:stas7373@localhost:5433/bookingdata"

	pgxConn := &pgxpool.Pool{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Подключение к базе данных с использованием pgxpool
	pgxConn, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	repo := repository.NewPostgresRepository(pgxConn)

	a.service = service.NewBookingService(repo)
	return nil
}

func (a *App) Start(ctx context.Context) error {
	route := handler.SetupRoutes(a.service)
	a.server = &http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	// Запуск HTTP-сервера в отдельной горутине
	func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	log.Println("Server started on :8080")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	// Завершаем работу HTTP-сервера
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	// Закрытие пула соединений к базе данных
	a.dbPool.Close()
	log.Println("Database connection closed")

	return nil
}
