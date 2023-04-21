package config

type DB struct {
	DSN string `env:"POSTGRES_DSN" envDefault:"host=localhost port=5432 user=postgres password=DB_PASSWORD dbname=postgres sslmode=disable"`
}
