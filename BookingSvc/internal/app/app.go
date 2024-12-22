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
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	server *http.Server
	dbPool *pgxpool.Pool
	log    *zap.Logger
}

func NewApp() *App {
	return &App{}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	logger.Info("Connecting to database", zap.String("connection_string", connString))
	return pgxpool.Connect(ctx, connString)
}

func NewKafkaProducer(cfg *config.Config, logger *zap.Logger) *kafka.Producer {
	logger.Info("Initializing Kafka producer", zap.String("broker", cfg.KafkaBroker))
	return kafka.NewProducer([]string{cfg.KafkaBroker}, cfg.KafkaTopicClient, cfg.KafkaTopicHotel)
}

func NewHotelClient(cfg *config.Config, logger *zap.Logger) (*grpc.HotelSvcClient, error) {
	logger.Info("Initializing Hotel service client", zap.String("host", cfg.GRPCHost), zap.String("port", cfg.GRPCPort))
	return grpc.NewHotelClient(cfg.GRPCHost, cfg.GRPCPort)
}

func NewPaymentClient(cfg *config.Config, logger *zap.Logger) *paymentClient.Client {
	logger.Info("Initializing Payment service client", zap.String("url", cfg.PaymentSvcURL))
	return paymentClient.NewPaymentSvcClient(cfg.PaymentSvcURL)
}

func (a *App) Init(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("myerror initializing zap logger: %v", err)
	}
	a.log = logger

	a.log.Info("Loading configuration")
	// Тут можно сделать MustLoad ля-ля
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.log.Info("Initializing dependencies")
	kafkaProducer := NewKafkaProducer(cfg, a.log)

	hotelClient, err := NewHotelClient(cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize hotel client: %w", err)
	}

	paymentSvcClient := NewPaymentClient(cfg, a.log)
	dbPool, err := NewDatabasePool(ctx, cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize database pool: %w", err)
	}
	repo := postgres.NewPostgresRepository(dbPool)
	a.dbPool = dbPool

	service := service.NewBookingServiceImpl(repo, kafkaProducer, hotelClient, paymentSvcClient, a.log)
	bookingHandler := handler.NewBookingHandler(service)
	route := handler.SetupRoutes(bookingHandler)

	a.server = &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: route,
	}
	a.log.Debug("Initialization complete")
	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting HTTP server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Error in ListenAndServe", zap.Error(err))
			return fmt.Errorf("failed to serve HTTP server: %w", err)
		}
		a.log.Info("HTTP server stopped")
		return nil
	})

	group.Go(func() error {
		<-groupCtx.Done()
		return a.Stop(context.Background())
	})

	if err := group.Wait(); err != nil {
		a.log.Error("Error after wait", zap.Error(err))
		return err
	}
	a.log.Info("Server shutdown gracefully")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	a.log.Info("Shutting down HTTP server")
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("HTTP server shutdown myerror", zap.Error(err))
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	a.log.Info("HTTP server shutdown gracefully")

	if a.dbPool != nil {
		a.dbPool.Close()
		a.log.Info("Database connection closed")
	} else {
		a.log.Warn("Database pool is nil, skipping close")
	}
	return nil
}
