# Hezy by Blitz

Personal Telegram AI assistant built in Go.

Created by Blitz (@blitzlabx)

## Features

- Persistent conversation memory (SQLite)
- AI chat with system prompt (configurable via .env)
- Image generation (`/imagine <prompt>`)
- Interactive games
  - Tic-Tac-Toe (button board)
  - Rock Paper Scissors
  - Number guess
  - Dice
- Tools (live football data and more)
- Main menu with inline buttons
- Group support (only replies when mentioned or when a command is used)
- `/ping` HTTP endpoint for uptime monitors
- Full configuration through environment variables
- Single binary, Docker ready
- Designed for Render free tier

## Requirements

- Go 1.22+ (1.24 recommended)
- CGO enabled (for SQLite)
- Telegram bot token from @BotFather

## Quick start (local)

```bash
cp .env.example .env
# edit .env and put your TELEGRAM_TOKEN and ADMIN_ID

go mod tidy
CGO_ENABLED=1 go run ./cmd/hezy
```

## Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| TELEGRAM_TOKEN | Yes | Bot token from @BotFather |
| ADMIN_ID | No | Your Telegram user ID |
| DONATION_URL | No | Support / donation link shown in About |
| LOGO_URL | No | Optional logo URL |
| SYSTEM_PROMPT | No | Personality and instructions for the AI |
| PORT | No | HTTP port (default 8080) |
| POLLING | No | `true` / `false` (default true) |
| MEMORY_DB_PATH | No | SQLite file path |
| MAX_HISTORY | No | Max messages kept per user (default 20) |

## Commands

- `/start` or `/menu` – main menu
- `/ping` – health check (also available as HTTP `/ping`)
- `/clear` – wipe conversation memory for the current user
- `/imagine <prompt>` – generate an image
- `/help` – short help
- Any other text – AI conversation (memory is kept)

## Buttons

Main menu, games and tools use inline keyboards.  
Colored button styles (`primary` / `success` / `danger`) are supported by Telegram Bot API 9.4+.  
The current library version uses clear labels; style fields can be added when the library fully exposes them.

## Deploy on Render (free tier)

Go is natively supported on Render.

1. Push this folder to a GitHub repository.
2. On Render → New → Web Service.
3. Connect the repository.
4. Settings:
   - **Runtime**: Go
   - **Build Command**: `CGO_ENABLED=1 go build -o hezy ./cmd/hezy`
   - **Start Command**: `./hezy`
   - **Instance type**: Free
5. Add environment variables from `.env.example` (at least `TELEGRAM_TOKEN`).
6. Set `POLLING=true`.
7. Deploy.

After deploy, open the service URL + `/ping` to confirm it is alive.  
Use that URL with UptimeRobot or any cron/uptime service so the free instance does not sleep too aggressively.

### Alternative: Docker on Render

If you prefer Docker:

1. New → Web Service → Docker
2. Use the included `Dockerfile`
3. Set the same environment variables
4. Deploy

## Dockerfile

The project includes a multi-stage Dockerfile that produces a small runtime image with SQLite support.

```dockerfile
FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /hezy ./cmd/hezy

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /hezy .
RUN mkdir -p /app/data
ENV PORT=8080
EXPOSE 8080
CMD ["./hezy"]
```

## Project structure

```
hezy/
├── cmd/hezy/main.go
├── internal/
│   ├── config/
│   ├── memory/
│   ├── ai/
│   ├── keyboard/
│   ├── games/
│   └── handlers/
├── data/
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## Building a release binary

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hezy ./cmd/hezy
```

## Notes

- Conversation history is stored per user in SQLite.
- Free Render instances sleep after inactivity. Keep them awake with a `/ping` monitor.
- Image generation and some tools depend on external free endpoints. Availability can vary.
- Never commit your real `.env` file.

## License

Use freely for personal and portfolio projects.  
Credit Blitz (@blitzlabx) if you share it publicly.

---

Hezy by Blitz  
@blitzlabx
