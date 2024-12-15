package app

import (
	"context"
	"errors"
	"fmt"
	grpc "github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/grpc/hotelsvc"
	paymentClient "github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/http/paymentsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/kafka"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/config"
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
	server *http.Server
}

func NewApp() *App {
	return &App{}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return pgxpool.Connect(ctx, connString)
}

func NewKafkaProducer(cfg *config.Config) *kafka.Producer {
	return kafka.NewProducer([]string{cfg.KafkaBroker}, cfg.KafkaTopicClient, cfg.KafkaTopicHotel)
}

func NewHotelClient(cfg *config.Config) (*grpc.HotelSvcClient, error) {
	return grpc.NewHotelClient(cfg.GRPCHost, cfg.GRPCPort)
}

func NewPaymentClient(cfg *config.Config) *paymentClient.Client {
	return paymentClient.NewPaymentSvcClient(cfg.PaymentSvcURL)
}

func (a *App) Init(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	kafkaProducer := NewKafkaProducer(cfg)
	hotelClient, err := NewHotelClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize hotel client: %w", err)
	}
	paymentSvcClient := NewPaymentClient(cfg)

	dbPool, err := NewDatabasePool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database pool: %w", err)
	}
	repo := postgres.NewPostgresRepository(dbPool)

	service := service.NewBookingService(repo, kafkaProducer, hotelClient, paymentSvcClient)
	bookingHandler := handler.NewBookingHandler(service)
	route := handler.SetupRoutes(bookingHandler)

	a.server = &http.Server{
		Addr:    ":" + cfg.HTTPPort,
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

	return nil
}
