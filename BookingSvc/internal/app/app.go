package app

import (
	"BookingSvc/internal/controller/handler"
	"BookingSvc/internal/infrastructure/repository"
	"BookingSvc/internal/service"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Получение переменных для базы данных
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	//connString := "postgres://postgres:stas7373@localhost:5433/bookingdata"

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
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // значение по умолчанию
	}
	log.Printf("Starting HTTP server on port %s", httpPort)

	route := handler.SetupRoutes(a.service)
	a.server = &http.Server{
		Addr:    ":" + httpPort,
		Handler: route,
	}

	// Настройка graceful shutdown с использованием errgroup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)

	// Запуск HTTP-сервера в отдельной горутине
	group.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error in ListenAndServe: %v", err)
			return fmt.Errorf("failed to serve HTTP server: %w", err)
		}
		log.Println("HTTP server stopped")
		return nil
	})

	// Обработка shutdown по сигналу
	group.Go(func() error {
		<-groupCtx.Done()
		return a.Stop(context.Background())
	})

	// Ожидание завершения работы сервера или ошибки
	if err := group.Wait(); err != nil {
		log.Printf("Error after wait: %v", err)
		return err
	}

	log.Println("Server shutdown gracefully")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	// Завершение работы HTTP-сервера с graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	log.Println("HTTP server shutdown gracefully")

	// Закрытие пула соединений к базе данных
	a.dbPool.Close()
	log.Println("Database connection closed")

	return nil
}
