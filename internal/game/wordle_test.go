package game

import (
	"errors"
	"reflect"
	"testing"

	"github.com/agravlin/wordle-server/internal/errs"
)

func TestMakeGuessUpdatesBoardAndFinishesGame(t *testing.T) {
	g := NewGame("apple", []string{"alice"})

	feedback, err := g.MakeGuess("alice", "apple")
	if err != nil {
		t.Fatalf("MakeGuess() error = %v", err)
	}

	if !reflect.DeepEqual(feedback, []Feedback{StatusGreen, StatusGreen, StatusGreen, StatusGreen, StatusGreen}) {
		t.Fatalf("MakeGuess() feedback = %v, want all greens", feedback)
	}

	if g.State != StateFinished {
		t.Fatalf("MakeGuess() state = %v, want %v", g.State, StateFinished)
	}

	board, ok := g.GameState.Boards["alice"]
	if !ok || board == nil {
		t.Fatalf("MakeGuess() did not create a board for alice")
	}

	if board.CurrentRow != 1 {
		t.Fatalf("MakeGuess() current row = %d, want 1", board.CurrentRow)
	}

	wantRow := []int{int(StatusGreen), int(StatusGreen), int(StatusGreen), int(StatusGreen), int(StatusGreen)}
	if !reflect.DeepEqual(board.Grid[0], wantRow) {
		t.Fatalf("MakeGuess() board row = %v, want %v", board.Grid[0], wantRow)
	}
}

func TestCheckGuess(t *testing.T) {
	tests := []struct {
		name       string
		guess      string
		target     string
		wantResult []Feedback
		wantErr    error
	}{
		{
			name:       "Perfect match",
			guess:      "APPLE",
			target:     "APPLE",
			wantResult: []Feedback{StatusGreen, StatusGreen, StatusGreen, StatusGreen, StatusGreen},
			wantErr:    nil,
		},
		{
			name:       "No letters match",
			guess:      "GHOST",
			target:     "APPLE",
			wantResult: []Feedback{StatusGray, StatusGray, StatusGray, StatusGray, StatusGray},
			wantErr:    nil,
		},
		{
			name:       "Mixed match with yellows and greens",
			guess:      "MAPLE",
			target:     "APPLE",
			wantResult: []Feedback{StatusGray, StatusYellow, StatusGreen, StatusGreen, StatusGreen},
			wantErr:    nil,
		},
		{
			name:   "Priority check",
			guess:  "ALAMO", // Contains two 'A's
			target: "BLAME", // Contains one 'A' (at index 2)
			// The first 'A' must be Gray because the exact match at index 2 (Green) consumes the only available 'A' stock.
			wantResult: []Feedback{StatusGray, StatusGreen, StatusGreen, StatusGreen, StatusGray},
			wantErr:    nil,
		},
		{
			name:       "Error: Length mismatch",
			guess:      "LONGER",
			target:     "SHORT",
			wantResult: nil,
			wantErr:    errs.ErrLengthMismatch,
		},
		{
			name:       "Error: Invalid character (Non-ASCII)",
			guess:      "PLAÇT", // No length mismatch because ç is 2 bytes
			target:     "TARGET",
			wantResult: nil,
			wantErr:    errs.ErrInvalidChar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotErr := CheckGuess(tt.guess, tt.target)

			// 1. Error check
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("CheckGuess() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			// 2. Result (Feedback slice) check
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("CheckGuess() result = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
