package app

import (
	"L0-wbtech/internal/config"
	"L0-wbtech/internal/handler"
	"L0-wbtech/internal/kafka"
	"L0-wbtech/internal/service"
	"L0-wbtech/pkg/logger/sl"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg          *config.Config
	log          *slog.Logger
	orderService service.Service
	consumer     *kafka.Consumer
	httpServer   *http.Server
}

func New(
	cfg *config.Config,
	orderService service.Service,
	consumer *kafka.Consumer,
	log *slog.Logger,
) *App {
	return &App{
		cfg:          cfg,
		orderService: orderService,
		consumer:     consumer,
		log:          log,
	}
}

func NewKafkaConsumer(
	cfg *config.Config,
	service service.Service,
	log *slog.Logger,
) *kafka.Consumer {
	return kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
		service,
		log,
	)
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.consumer.Start(ctx)

	go a.runHTTPServer()

	a.log.Info("Application started",
		"port", a.cfg.Server.Port,
		"kafka_topic", a.cfg.Kafka.Topic)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	const op = "app.Run"
	a.log.With(slog.String("op", op)).Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := a.httpServer.Shutdown(ctxShutdown); err != nil {
		a.log.Error("HTTP server shutdown error", sl.Err(err))
	}

	if err := a.consumer.Close(); err != nil {
		a.log.Error("Kafka consumer close error", sl.Err(err))
	}

	if err := a.orderService.Close(); err != nil {
		a.log.Error("Order service close error", sl.Err(err))
	}

	a.log.Info("Application stopped")
}

func (a *App) runHTTPServer() {
	const op = "app.runHTTPServer"
	log := a.log.With(slog.String("op", op))

	gin.SetMode(gin.ReleaseMode)
	if a.cfg.Env == "local" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(a.log))

	apiHandler := handler.New(a.orderService, a.log)
	apiHandler.RegisterRoutes(router)

	a.httpServer = &http.Server{
		Addr:    ":" + a.cfg.Server.Port,
		Handler: router,
	}

	log.Info("Starting HTTP server", "port", a.cfg.Server.Port)
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("HTTP server error", sl.Err(err))
		os.Exit(1)
	}
}

func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		dataFetchStart := time.Now()
		c.Set("data_fetch_start", dataFetchStart)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		dataFetchTime := time.Duration(0)
		if startVal, exists := c.Get("data_fetch_start"); exists {
			dataFetchStart = startVal.(time.Time)
			dataFetchTime = end.Sub(dataFetchStart)
		}

		log.Info("request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", latency,
			"data_fetch_time", dataFetchTime,
		)
	}
}
