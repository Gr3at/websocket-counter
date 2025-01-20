package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	AllowedOrigins map[string]bool
}

func New() *Config {
	cfg := &Config{}
	cfg.loadAllowedOrigins()
	return cfg
}

func (c *Config) loadAllowedOrigins() {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	c.AllowedOrigins = make(map[string]bool)

	if allowedOrigins == "" {
		log.Println("No ALLOWED_ORIGINS environment variable set")
		return
	}

	for _, origin := range strings.Split(allowedOrigins, ",") {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			c.AllowedOrigins[trimmedOrigin] = true
		}
	}
}
