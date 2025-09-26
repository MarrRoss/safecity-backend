package postgres

import "fmt"

type Config struct {
	Host     string `env:"DATABASE_HOST,required"`
	Port     string `env:"DATABASE_PORT"              envDefault:"5432"`
	UserName string `env:"DATABASE_USERNAME,required"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Name     string `env:"DATABASE_NAME,required"`
}

func (c *Config) Dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		c.Host,
		c.Port,
		c.UserName,
		c.Password,
		c.Name,
	)
}
