package game

import (
	_ "embed"
	"math/rand"
	"strings"

	"github.com/agravlin/wordle-server/internal/errs"
)

//go:embed words.txt
var wordListFile string
var ValidWords []string
var validWordsMap = make(map[string]bool)

type Feedback int

const (
	StatusEmpty  Feedback = iota // Empty box
	StatusTyped                  // Letter typed (but not submitted)
	StatusGray                   // 0 (Absent)
	StatusYellow                 // 1 (Present)
	StatusGreen                  // 2 (Correct)
)

func init() {
	lines := strings.Split(wordListFile, "\n")

	for _, line := range lines {
		word := strings.TrimSpace(line)
		word = strings.ToLower(word)

		if word != "" {
			ValidWords = append(ValidWords, word)
			validWordsMap[word] = true
		}
	}
}

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
			result[i] = StatusGreen
			targetCounts[guess[i]]--
		}
	}

	// Mark rest as yellow or grey
	for i := range length {
		if result[i] == StatusGreen {
			continue
		}

		if targetCounts[guess[i]] > 0 {
			result[i] = StatusYellow
			targetCounts[guess[i]]--
		} else {
			result[i] = StatusGray
		}
	}

	return result, nil
}

func PickRandomWord() string {
	if len(ValidWords) == 0 {
		return ""
	}

	randomIndex := rand.Intn(len(ValidWords))
	return ValidWords[randomIndex]
}

func isLegitWord(guess string) bool {
	guess = strings.ToLower(strings.TrimSpace(guess))

	return validWordsMap[guess]
}
