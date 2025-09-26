package api

type Config struct {
	Port          string `env:"API_PORT" envDefault:"8000"`
	HydraAdminURL string `env:"HYDRA_URL"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`
}
