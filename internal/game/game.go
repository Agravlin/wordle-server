package game

import (
	"strings"

	"github.com/agravlin/wordle-server/internal/errs"
)

type Game struct {
	State      State
	TargetWord string
	Players    []string // Only nicknames
	GameState  *GameState
}

type Board struct {
	Grid       [][]int `json:"grid"`
	CurrentRow int
}

type GameState struct {
	Boards map[string]*Board `json:"boards"`
}

func NewGame(target string, players []string) *Game {
	target = strings.ToLower(strings.TrimSpace(target))

	gameState := &GameState{Boards: make(map[string]*Board)}
	for _, player := range players {
		gameState.Boards[player] = &Board{
			CurrentRow: 0,
			Grid:       make([][]int, 0),
		}
	}

	return &Game{
		State:      StatePlaying,
		TargetWord: target,
		Players:    players,
		GameState:  gameState,
	}
}

func (g *Game) MakeGuess(nick string, guess string) ([]Feedback, error) {
	if g.State != StatePlaying {
		return nil, errs.ErrGameOver
	}

	guess = strings.ToLower(strings.TrimSpace(guess))
	if !isLegitWord(guess) {
		return nil, errs.ErrNotLegitWord
	}

	feedback, err := CheckGuess(guess, g.TargetWord)
	if err != nil {
		return nil, err
	}

	if err := g.updateBoard(nick, feedback); err != nil {
		return nil, err
	}

	if isPerfectGuess(feedback) {
		g.State = StateFinished
	}

	return feedback, nil
}

func (g *Game) updateBoard(nick string, feedback []Feedback) error {
	if g.GameState == nil {
		g.GameState = &GameState{Boards: make(map[string]*Board)}
	}
	if g.GameState.Boards == nil {
		g.GameState.Boards = make(map[string]*Board)
	}

	board, ok := g.GameState.Boards[nick]
	if !ok {
		return &errs.NickNotFound{Nick: nick}
	}

	row := make([]int, len(feedback))
	for i, item := range feedback {
		row[i] = int(item)
	}

	board.Grid = append(board.Grid, row)
	board.CurrentRow = len(board.Grid)

	return nil
}

func isPerfectGuess(feedback []Feedback) bool {
	for _, item := range feedback {
		if item != StatusGreen {
			return false
		}
	}
	return len(feedback) > 0
}
