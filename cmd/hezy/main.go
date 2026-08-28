package main

import (
	"log"
	"os"

	"github.com/blitzlabx/hezy/internal/config"
	"github.com/blitzlabx/hezy/internal/handlers"
	"github.com/blitzlabx/hezy/internal/memory"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	_ = os.MkdirAll("./data", 0755)

	store, err := memory.New(cfg.MemoryDBPath, cfg.MaxHistory)
	if err != nil {
		log.Fatal("memory:", err)
	}
	defer store.Close()

	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	api.Debug = false
	log.Printf("Authorized as @%s", api.Self.UserName)

	bot := handlers.New(api, cfg, store)

	go bot.StartHTTP(cfg.Port)

	if cfg.Polling {
		log.Println("Starting long polling")
		bot.StartPolling()
	} else {
		log.Println("Polling disabled. Set POLLING=true or configure webhook.")
		select {}
	}
}
