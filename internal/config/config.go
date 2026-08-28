package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	AdminID       int64
	DonationURL   string
	LogoURL       string
	SystemPrompt  string
	Port          string
	WebhookURL    string
	Polling       bool
	MemoryDBPath  string
	MaxHistory    int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	maxHist, _ := strconv.Atoi(os.Getenv("MAX_HISTORY"))
	if maxHist <= 0 {
		maxHist = 20
	}

	polling := true
	if v := strings.ToLower(os.Getenv("POLLING")); v == "false" || v == "0" {
		polling = false
	}

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		AdminID:       adminID,
		DonationURL:   os.Getenv("DONATION_URL"),
		LogoURL:       os.Getenv("LOGO_URL"),
		SystemPrompt:  os.Getenv("SYSTEM_PROMPT"),
		Port:          os.Getenv("PORT"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		Polling:       polling,
		MemoryDBPath:  os.Getenv("MEMORY_DB_PATH"),
		MaxHistory:    maxHist,
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.MemoryDBPath == "" {
		cfg.MemoryDBPath = "./data/hezy.db"
	}
	if cfg.SystemPrompt == "" {
		cfg.SystemPrompt = "You are Hezy, a helpful personal AI assistant created by Blitz (@blitzlabx). Be concise, useful and friendly. Remember conversation context."
	}

	return cfg, nil
}
