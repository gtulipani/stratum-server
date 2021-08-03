package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stratum-server/config"
	"stratum-server/controller"
	"stratum-server/repository"
	"stratum-server/service"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	envFile = ".env"
)

func main() {
	if _, err := os.Stat(envFile); !os.IsNotExist(err) {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("error loading .env file: %v", err.Error())
		}
	}

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err.Error())
	}

	postgres := repository.NewRepository(cfg.PostgreSQLConfig)
	svc := service.NewService(postgres, cfg.SubscriptionsTable)
	handler := controller.NewHandler(svc)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("error on server shutdown: %s", err.Error())
		}
	}()

	log.Printf("HTTP listener started on :%s @ %s", cfg.HTTPPort, time.Now().Format(time.RFC3339))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start http server: %s", err.Error())
	}
}
