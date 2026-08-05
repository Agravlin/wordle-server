package game

import (
	_ "embed"
	"strings"

	"github.com/agravlin/wordle-server/internal/errs"
)

//go:embed words.txt
var wordListFile string

var ValidWords []string

type Feedback int

const (
	Gray   Feedback = iota // 0 (Absent)
	Yellow                 // 1 (Present)
	Green                  // 2 (Correct)
)

func CheckGuess(guess, target string) ([]Feedback, error) {
	if len(guess) != len(target) {
		return nil, errs.ErrLengthMismatch
	}

	length := len(guess) // Now always 5, potentially for future
	result := make([]Feedback, length)
	targetCounts := make(map[byte]int)

	// Invalid char check
	for i := range length {
		if guess[i] > 127 {
			return nil, errs.ErrInvalidChar
		}
		targetCounts[target[i]]++
	}

	// Match correct letters
	for i := range length {
		if guess[i] == target[i] {
			result[i] = Green
			targetCounts[guess[i]]--
		}
	}

	// Mark rest as yellow or grey
	for i := range length {
		if result[i] == Green {
			continue
		}

		if targetCounts[guess[i]] > 0 {
			result[i] = Yellow
			targetCounts[guess[i]]--
		} else {
			result[i] = Gray
		}
	}

	return result, nil
}

func init() {
	cleanFile := strings.ReplaceAll(wordListFile, "\r", "") // Windows' CRLF
	ValidWords = strings.Split(strings.TrimSpace(cleanFile), "\n")
}
