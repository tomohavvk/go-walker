package main

import (
	"errors"
	"fmt"
	"github.com/tomohavvk/go-walker/config"
	"github.com/tomohavvk/go-walker/db"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/service"
	"github.com/tomohavvk/go-walker/internal/web"
	"github.com/tomohavvk/go-walker/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.LoadConfig()

	if err := runDBMigration(cfg); err != nil {
		log.Fatalf("Failed to perform db migration: %v", err)
	}

	conn, err := initDBConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}

	logger := logging.NewLogger(os.Stdout, slog.Level(cfg.LogLevel))

	deviceRepository := repository.NewDeviceRepository(conn)
	groupRepository := repository.NewGroupRepository(conn)
	groupMessagesRepository := repository.NewGroupMessagesRepository(conn)
	deviceLocationRepository := repository.NewDeviceLocationRepository(conn)

	deviceService := service.NewDeviceService(logger, deviceRepository)
	groupService := service.NewGroupService(logger, groupRepository)
	groupMessagesService := service.NewGroupMessagesService(logger, groupMessagesRepository)
	deviceLocationService := service.NewDeviceLocationService(logger, deviceLocationRepository)

	wsHandler := web.NewWSMessageHandler(logger, deviceService, groupService, groupMessagesService, deviceLocationService)
	routes := web.NewRoutes(logger, wsHandler, deviceService)
	server := newHTTPServer(routes, cfg.HttpServer)

	interruptSignal := make(chan os.Signal)
	signal.Notify(interruptSignal, os.Interrupt)

	go func() {
		<-interruptSignal
		log.Println("interrupt signal received")

		sqlDb, err := conn.DB()
		if err != nil {
			log.Fatalf("failed to get sql db instance: %v", err)
		}

		if err := sqlDb.Close(); err != nil {
			log.Fatalf("failed to close db connection: %v", err)
		}

		log.Println("db connection successfully closed")

		if err := server.Close(); err != nil {
			log.Fatalf("failed to close http server: %v", err)
		}

		log.Println("http server successfully closed")
	}()

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed under interrupt signal")
		} else {
			log.Fatalf("unexpected close server: %v", err)
		}
	}
}

func initDBConnection(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.User, cfg.Password, cfg.Name, cfg.Host, cfg.Port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func runDBMigration(cfg config.AppConfig) error {
	return db.PerformMigration(cfg.DB)
}

func newHTTPServer(routes web.Routes, cfg config.HttpServer) http.Server {

	return http.Server{Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), Handler: routes.Setup().Handler()}
}
