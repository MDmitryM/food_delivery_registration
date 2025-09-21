package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MDmitryM/food_delivery_registration/handler"
	"github.com/MDmitryM/food_delivery_registration/repository"
	"github.com/MDmitryM/food_delivery_registration/service"
	"github.com/MDmitryM/food_delivery_registration/telemetry"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg := repository.PostgresCfg{ //TODO: create .env
		Host:        os.Getenv("AUTH_DB_HOST"),
		Port:        os.Getenv("AUTH_DB_PORT"),
		PG_USER:     os.Getenv("POSTGRES_USER"),
		PG_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		PG_DB:       os.Getenv("POSTGRES_DB"),
		SSL_MODE:    os.Getenv("DB_SSL_MODE"),
	}

	rootCtx, cancel := context.WithCancel(context.Background())

	repo, err := repository.NewPostgresDB(rootCtx, cfg)
	if err != nil {
		logrus.Fatalf("%v", err.Error())
	}
	defer repo.Close()

	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Order system v0.0",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  10 * time.Second,
	})

	tracerCfg := telemetry.TracerCfg{
		ServiceName: "food_delivery_registration",
		JaegerUrl:   os.Getenv("JAEGER_URL"),
		JaegerPort:  os.Getenv("JAEGER_PORT"),
	}

	tp, err := telemetry.InitTelemetry(tracerCfg)
	if err != nil {
		logrus.Errorf("Can`t create tracer")
	} else {
		logrus.Info("Successfully created new OTLP tracer")

		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logrus.Errorf("Error shutting down tracer provider, %v", err)
			}
		}()

		app.Use(
			otelfiber.Middleware(),
		)
	}

	handler.InitRoutes(app)

	go func() {
		if err := app.Listen(":" + os.Getenv("AUTH_PORT")); err != nil {
			logrus.Fatalf("error while server start: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	cancel()
	if err := app.Shutdown(); err != nil {
		logrus.Fatalf("error while server shutdown, %s", err.Error())
	}

	logrus.Println("server gracefully stopped!")
}
