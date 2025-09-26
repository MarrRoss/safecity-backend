package config

import (
	"awesomeProjectDDD/config/api"
	"awesomeProjectDDD/config/postgres"
	_ "github.com/joho/godotenv/autoload"
)

type Factory struct {
	API      api.Config
	Postgres postgres.Config
}

func New() *Factory {
	cfg := Load()
	return cfg
}
