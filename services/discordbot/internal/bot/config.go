package bot

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application.
type Config struct {
	BotToken string
	GuildID  string
	// EnableDefaultIntents controla se devemos usar um fallback de intents padrão
	// quando nenhuma feature declara intents explicitamente.
	// Controlado via variável de ambiente DEFAULT_INTENTS (true/false).
	EnableDefaultIntents bool
}

// LoadConfig reads configuration from environment variables.
// It no longer reads from a .env file directly; this is expected to be handled
// by the execution environment (e.g., Docker Compose).
func LoadConfig() (*Config, error) {
	botToken := strings.TrimSpace(os.Getenv("BOT_TOKEN"))
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN environment variable is not set")
	}

	cfg := &Config{
		BotToken: botToken,
		GuildID:  strings.TrimSpace(os.Getenv("GUILD_ID")),
	}

	if v := strings.TrimSpace(os.Getenv("DEFAULT_INTENTS")); v != "" {
		parsed, err := parseBoolEnv(v)
		if err != nil {
			return nil, fmt.Errorf("DEFAULT_INTENTS: %w", err)
		}
		cfg.EnableDefaultIntents = parsed
	}

	return cfg, nil
}

func parseBoolEnv(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "1", "true", "on", "yes":
		return true, nil
	case "0", "false", "off", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid value %q (expected true/false, 1/0, on/off, yes/no)", value)
	}
}


