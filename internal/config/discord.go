package config

type Discord struct {
	NewDaosURL string `env:"DISCORD_NEW_DAOS_URL" require:"true"`
}
