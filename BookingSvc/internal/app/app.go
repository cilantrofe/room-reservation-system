package app

import (
	"context"
	"errors"
	"fmt"
	grpc "github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/grpc/hotelsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/kafka"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/controller/handler"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/service"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/storage/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
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
	//if err := godotenv.Load(); err != nil {
	//	log.Println("Warning: No .env file found")
	//}

	// Получение переменных для базы данных
	dbHost := os.Getenv("BOOKING_DB_HOST")
	dbPort := os.Getenv("BOOKING_DB_PORT")
	dbUser := os.Getenv("BOOKING_DB_USER")
	dbPassword := os.Getenv("BOOKING_DB_PASSWORD")
	dbName := os.Getenv("BOOKING_DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Подключение к базе данных с использованием pgxpool
	pgxConn, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	a.dbPool = pgxConn
	repo := postgres.NewPostgresRepository(pgxConn)
	grpcHost := os.Getenv("BOOKING_GRPC_HOST")
	grpcPort := os.Getenv("BOOKING_GRPC_PORT")
	hotelClient, err := grpc.NewHotelClient(grpcHost, grpcPort)
	if err != nil {
		log.Fatalf("Unable to connect to hotel: %v\n", err)
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	kafkaProducer := kafka.NewProducer([]string{kafkaBroker}, kafkaTopic)

	a.service = service.NewBookingService(repo, kafkaProducer, hotelClient)
	bookingHandler := handler.NewBookingHandler(a.service)
	route := handler.SetupRoutes(bookingHandler)

	httpPort := os.Getenv("BOOKING_HTTP_PORT")
	a.server = &http.Server{
		Addr:    ":" + httpPort,
		Handler: route,
	}
	return nil
}

func (a *App) Start(ctx context.Context) error {
	log.Printf("Starting HTTP server")
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
	if a.dbPool != nil {
		a.dbPool.Close()
		log.Println("Database connection closed")
	} else {
		log.Println("Database pool is nil, skipping close")
	}

	return nil
}
