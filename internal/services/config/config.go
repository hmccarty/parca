package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Discord
	DiscordToken   string `yaml:"discord_token"`
	DiscordAppID   string `yaml:"discord_app_id"`
	DiscordGuildID string `yaml:"discord_guild_id"`

	// Google Calendar
	GoogleServiceEmail   string `yaml:"google_service_email"`
	GooglePrivatePEMFile string `yaml:"google_private_pem_file"`

	// SMTP
	SMTPEmail    string `yaml:"smtp_email"`
	SMTPPassword string `yaml:"smtp_password"`
	SMTPServer   string `yaml:"smtp_server"`
	SMTPPort     string `yaml:"smtp_port"`
}

func NewConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
