package config

type Config struct {
	PG struct {
		DSN string `envconfig:"PG_DSN" default:"user=postgres dbname=cerebro password=postgres sslmode=disable host=127.0.0.1"`
	}
}
