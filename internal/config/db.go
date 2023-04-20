package config

type DB struct {
	PostgresAddr string `env:"POSTGRES_ADDR" envDefault:"host=localhost port=5432 user=postgres password=DB_PASSWORD dbname=postgres sslmode=disable"`
}
