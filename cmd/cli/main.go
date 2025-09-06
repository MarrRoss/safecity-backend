package main

import (
	"awesomeProjectDDD/config"
	_ "awesomeProjectDDD/docs"
	"awesomeProjectDDD/internal/adapter/db_storage"
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/logging"
	"awesomeProjectDDD/pkg/postgres"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

func main() {
	cfg := config.New()

	logger := logging.NewZeroLogger(zerolog.TraceLevel)
	observer := observability.New(logger)

	pg, err := postgres.New(
		cfg.Postgres.Dsn(),
		postgres.MaxPoolSize(10),
		postgres.ConnAttempts(20),
		postgres.ConnTimeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}

	//httpClient := httpResty.NewClient(cfg.API.HydraAdminURL)
	//hydraService := hydra.NewService(httpClient)

	dbUser, err := db_storage.NewUserRepositoryImpl(observer, pg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	addUserCommand := user.NewAddUserHandler(dbUser, observer)
	cmd := user.AddUserCommand{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "admin@admin.ru",
		Login:      "admin1337",
		ExternalID: uuid.New(),
		Tracking:   true,
	}
	res, err := addUserCommand.Handle(context.Background(), cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("User ID:", res.String())

}
