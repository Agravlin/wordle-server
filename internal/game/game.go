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
	Grid [][]int `json:"grid"`
}

type GameState struct {
	Boards map[string]*Board `json:"boards"`
}

func NewGame(target string, players []string) *Game {
	return &Game{
		State:      StatePlaying,
		TargetWord: target,
		Players:    players,
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

	return feedback, nil
}
