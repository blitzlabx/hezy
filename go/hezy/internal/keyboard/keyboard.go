package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Chat", "menu_chat"),
			tgbotapi.NewInlineKeyboardButtonData("Imagine", "menu_imagine"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Games", "menu_games"),
			tgbotapi.NewInlineKeyboardButtonData("Tools", "menu_tools"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Clear Memory", "menu_clear"),
			tgbotapi.NewInlineKeyboardButtonData("About", "menu_about"),
		),
	)
}

func GamesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Tic-Tac-Toe", "game_ttt"),
			tgbotapi.NewInlineKeyboardButtonData("Rock Paper Scissors", "game_rps"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Number Guess", "game_guess"),
			tgbotapi.NewInlineKeyboardButtonData("Dice", "game_dice"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "menu_main"),
		),
	)
}

func ToolsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Football Scores", "tool_football"),
			tgbotapi.NewInlineKeyboardButtonData("Screenshot", "tool_ss"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "menu_main"),
		),
	)
}

func RPSButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Rock", "rps_rock"),
			tgbotapi.NewInlineKeyboardButtonData("Paper", "rps_paper"),
			tgbotapi.NewInlineKeyboardButtonData("Scissors", "rps_scissors"),
		),
	)
}

func ConfirmButtons(yesData, noData string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", yesData),
			tgbotapi.NewInlineKeyboardButtonData("No", noData),
		),
	)
}

func TTTBoard(board [9]string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < 3; i++ {
		var row []tgbotapi.InlineKeyboardButton
		for j := 0; j < 3; j++ {
			idx := i*3 + j
			label := board[idx]
			if label == "" {
				label = " "
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(label, "ttt_"+string(rune('0'+idx))))
		}
		rows = append(rows, row)
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Quit", "menu_games"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
