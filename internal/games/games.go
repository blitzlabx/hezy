package games

import (
	"math/rand"
	"strconv"
	"strings"
)

type TTTState struct {
	Board  [9]string
	Turn   string
	UserID int64
}

func NewTTT(userID int64) *TTTState {
	return &TTTState{
		Board:  [9]string{},
		Turn:   "X",
		UserID: userID,
	}
}

func (t *TTTState) Play(pos int) (string, bool) {
	if pos < 0 || pos > 8 || t.Board[pos] != "" {
		return "Invalid move", false
	}
	t.Board[pos] = t.Turn
	if winner := t.checkWinner(); winner != "" {
		return winner + " wins!", true
	}
	if t.isFull() {
		return "Draw!", true
	}
	if t.Turn == "X" {
		t.Turn = "O"
	} else {
		t.Turn = "X"
	}
	return "", false
}

func (t *TTTState) BotMove() (string, bool) {
	empty := []int{}
	for i, v := range t.Board {
		if v == "" {
			empty = append(empty, i)
		}
	}
	if len(empty) == 0 {
		return "Draw!", true
	}
	pos := empty[rand.Intn(len(empty))]
	t.Board[pos] = "O"
	if winner := t.checkWinner(); winner != "" {
		return winner + " wins!", true
	}
	if t.isFull() {
		return "Draw!", true
	}
	t.Turn = "X"
	return "", false
}

func (t *TTTState) checkWinner() string {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}
	for _, l := range lines {
		if t.Board[l[0]] != "" && t.Board[l[0]] == t.Board[l[1]] && t.Board[l[1]] == t.Board[l[2]] {
			return t.Board[l[0]]
		}
	}
	return ""
}

func (t *TTTState) isFull() bool {
	for _, v := range t.Board {
		if v == "" {
			return false
		}
	}
	return true
}

func (t *TTTState) Render() string {
	var b strings.Builder
	for i := 0; i < 9; i++ {
		cell := t.Board[i]
		if cell == "" {
			cell = "·"
		}
		b.WriteString(cell)
		if (i+1)%3 == 0 {
			b.WriteString("\n")
		} else {
			b.WriteString(" ")
		}
	}
	return b.String()
}

func RPS(userChoice string) string {
	choices := []string{"rock", "paper", "scissors"}
	bot := choices[rand.Intn(3)]
	userChoice = strings.ToLower(userChoice)

	result := "Draw"
	if userChoice == "rock" && bot == "scissors" ||
		userChoice == "paper" && bot == "rock" ||
		userChoice == "scissors" && bot == "paper" {
		result = "You win"
	} else if userChoice != bot {
		result = "You lose"
	}
	return "You: " + userChoice + "\nHezy: " + bot + "\n" + result
}

func Dice() string {
	return "You rolled: " + strconv.Itoa(rand.Intn(6)+1)
}
