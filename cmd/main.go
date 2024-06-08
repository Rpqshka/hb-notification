package main

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	hb "hb-notification"
	"hb-notification/pkg/handler"
	"hb-notification/pkg/repository"
	"hb-notification/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	//migrations
	m, err := migrate.New(
		"file://schema",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME")))
	if err != nil {
		logrus.Fatalf("Error creating migration instance: %v", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error applying migrations: %v", err)
	}
	logrus.Debug("Migrations applied successfully")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	c := cron.New()
	jobID, err := c.AddFunc("0 21 * * *", func() { services.Notification.CheckBirthday() })
	if err != nil {
		logrus.Fatalf("Error creating cron job: %v", err)
	}
	handlers.Cron = c
	handlers.JobID = jobID

	c.Start()
	defer c.Stop()
	logrus.Printf("Happy Birthday Notification Cron Started")
	srv := new(hb.Server)
	go func() {
		if err := srv.Run(os.Getenv("SERV_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Printf("Happy Birthday Notification App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Happy Birthday Notification App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
