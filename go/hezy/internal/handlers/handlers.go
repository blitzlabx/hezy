package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/blitzlabx/hezy/internal/ai"
	"github.com/blitzlabx/hezy/internal/config"
	"github.com/blitzlabx/hezy/internal/games"
	"github.com/blitzlabx/hezy/internal/keyboard"
	"github.com/blitzlabx/hezy/internal/memory"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	cfg    *config.Config
	store  *memory.Store
	ai     *ai.Client
	ttt    map[int64]*games.TTTState
	guess  map[int64]int
	mu     sync.Mutex
}

func New(api *tgbotapi.BotAPI, cfg *config.Config, store *memory.Store) *Bot {
	return &Bot{
		api:   api,
		cfg:   cfg,
		store: store,
		ai:    ai.New(cfg.SystemPrompt),
		ttt:   make(map[int64]*games.TTTState),
		guess: make(map[int64]int),
	}
}

func (b *Bot) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.Chat.Type != "private" && !strings.Contains(msg.Text, "@"+b.api.Self.UserName) && !strings.HasPrefix(msg.Text, "/") {
		return
	}

	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	userID := msg.From.ID

	switch {
	case text == "/start" || text == "/menu":
		b.sendMain(chatID)
	case text == "/ping":
		b.api.Send(tgbotapi.NewMessage(chatID, "pong"))
	case text == "/clear":
		_ = b.store.Clear(userID)
		b.api.Send(tgbotapi.NewMessage(chatID, "Memory cleared."))
	case text == "/help":
		b.api.Send(tgbotapi.NewMessage(chatID, "Commands:\n/start - main menu\n/ping - health\n/clear - reset memory\n/imagine <prompt>\nJust chat normally for AI help."))
	case strings.HasPrefix(text, "/imagine "):
		prompt := strings.TrimPrefix(text, "/imagine ")
		b.handleImagine(chatID, prompt)
	case strings.HasPrefix(text, "/"):
		b.api.Send(tgbotapi.NewMessage(chatID, "Unknown command. Try /start"))
	default:
		b.handleChat(chatID, userID, text)
	}
}

func (b *Bot) handleCallback(cq *tgbotapi.CallbackQuery) {
	data := cq.Data
	chatID := cq.Message.Chat.ID
	userID := cq.From.ID

	_, _ = b.api.Request(tgbotapi.NewCallback(cq.ID, ""))

	switch {
	case data == "menu_main":
		b.sendMain(chatID)
	case data == "menu_chat":
		b.api.Send(tgbotapi.NewMessage(chatID, "Send me any message and I will reply as Hezy."))
	case data == "menu_imagine":
		b.api.Send(tgbotapi.NewMessage(chatID, "Send /imagine followed by your prompt.\nExample: /imagine a cyberpunk city at night"))
	case data == "menu_games":
		msg := tgbotapi.NewMessage(chatID, "Choose a game:")
		msg.ReplyMarkup = keyboard.GamesMenu()
		b.api.Send(msg)
	case data == "menu_tools":
		msg := tgbotapi.NewMessage(chatID, "Tools:")
		msg.ReplyMarkup = keyboard.ToolsMenu()
		b.api.Send(msg)
	case data == "menu_clear":
		_ = b.store.Clear(userID)
		b.api.Send(tgbotapi.NewMessage(chatID, "Memory cleared."))
	case data == "menu_about":
		about := "Hezy by Blitz (@blitzlabx)\nPersonal AI assistant with memory, games and tools."
		if b.cfg.DonationURL != "" {
			about += "\nSupport: " + b.cfg.DonationURL
		}
		b.api.Send(tgbotapi.NewMessage(chatID, about))
	case data == "game_rps":
		msg := tgbotapi.NewMessage(chatID, "Choose:")
		msg.ReplyMarkup = keyboard.RPSButtons()
		b.api.Send(msg)
	case strings.HasPrefix(data, "rps_"):
		choice := strings.TrimPrefix(data, "rps_")
		result := games.RPS(choice)
		b.api.Send(tgbotapi.NewMessage(chatID, result))
	case data == "game_dice":
		b.api.Send(tgbotapi.NewMessage(chatID, games.Dice()))
	case data == "game_ttt":
		b.mu.Lock()
		b.ttt[userID] = games.NewTTT(userID)
		b.mu.Unlock()
		msg := tgbotapi.NewMessage(chatID, "Tic-Tac-Toe\nYou are X\n"+b.ttt[userID].Render())
		msg.ReplyMarkup = keyboard.TTTBoard(b.ttt[userID].Board)
		b.api.Send(msg)
	case strings.HasPrefix(data, "ttt_"):
		b.handleTTT(chatID, userID, data)
	case data == "game_guess":
		b.mu.Lock()
		b.guess[userID] = 1 + int(userID%100)
		b.mu.Unlock()
		b.api.Send(tgbotapi.NewMessage(chatID, "Guess a number between 1 and 100. Just send the number."))
	case data == "tool_football":
		raw, err := b.ai.FootballScores()
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(chatID, "Could not fetch scores right now."))
			return
		}
		if len(raw) > 3500 {
			raw = raw[:3500] + "..."
		}
		b.api.Send(tgbotapi.NewMessage(chatID, "Football data:\n"+raw))
	case data == "tool_ss":
		b.api.Send(tgbotapi.NewMessage(chatID, "Send a URL and I will try to capture it (coming soon)."))
	}
}

func (b *Bot) handleTTT(chatID, userID int64, data string) {
	b.mu.Lock()
	state, ok := b.ttt[userID]
	b.mu.Unlock()
	if !ok {
		b.api.Send(tgbotapi.NewMessage(chatID, "No active game. Start one from the menu."))
		return
	}
	posStr := strings.TrimPrefix(data, "ttt_")
	pos, _ := strconv.Atoi(posStr)
	msgText, ended := state.Play(pos)
	if msgText != "" && !ended {
		b.api.Send(tgbotapi.NewMessage(chatID, msgText))
		return
	}
	if ended {
		b.api.Send(tgbotapi.NewMessage(chatID, state.Render()+"\n"+msgText))
		b.mu.Lock()
		delete(b.ttt, userID)
		b.mu.Unlock()
		return
	}
	botMsg, botEnded := state.BotMove()
	text := state.Render()
	if botEnded {
		text += "\n" + botMsg
		b.mu.Lock()
		delete(b.ttt, userID)
		b.mu.Unlock()
	}
	msg := tgbotapi.NewMessage(chatID, text)
	if !botEnded {
		msg.ReplyMarkup = keyboard.TTTBoard(state.Board)
	}
	b.api.Send(msg)
}

func (b *Bot) handleChat(chatID, userID int64, text string) {
	_ = b.store.Append(userID, "user", text)
	hist, _ := b.store.Get(userID)
	var historyLines []string
	for _, m := range hist {
		historyLines = append(historyLines, m.Role+": "+m.Content)
	}

	reply, err := b.ai.Chat(text, historyLines)
	if err != nil {
		reply = "Sorry, I had trouble thinking right now."
	}
	_ = b.store.Append(userID, "assistant", reply)
	b.api.Send(tgbotapi.NewMessage(chatID, reply))
}

func (b *Bot) handleImagine(chatID int64, prompt string) {
	msg := tgbotapi.NewMessage(chatID, "Generating image...")
	sent, _ := b.api.Send(msg)

	url, err := b.ai.GenerateImage(prompt)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(chatID, sent.MessageID, "Could not generate image right now.")
		b.api.Send(edit)
		return
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(url))
	photo.Caption = prompt
	b.api.Send(photo)
	del := tgbotapi.NewDeleteMessage(chatID, sent.MessageID)
	b.api.Request(del)
}

func (b *Bot) sendMain(chatID int64) {
	text := "Hezy by Blitz\nYour personal AI assistant.\nChoose an option or just send a message."
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard.MainMenu()
	b.api.Send(msg)
}

func (b *Bot) StartHTTP(port string) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hezy by Blitz is running")
	})
	log.Println("HTTP listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
