package errs

import (
	"errors"
	"fmt"
)

const (
	msgLengthMismatch    = "Guess length does not match target length"
	msgInvalidChar       = "Guess contains invalid non-ASCII characters"
	msgEmptyNick         = "Nick can not be empty"
	msgRoomNotFound      = "Room %s not found"
	msgNickAlreadyExists = "Nick %s already exists in room %s"
	msgGameOver          = "Game is already finished"
	msgNotLegitWord      = "Not a word"
	msgNickNotFound      = "Nick %s was not found in room %s"
)

// Static errors
var (
	ErrLengthMismatch = errors.New(msgLengthMismatch)
	ErrInvalidChar    = errors.New(msgInvalidChar)
	ErrEmptyNick      = errors.New(msgEmptyNick)
	ErrGameOver       = errors.New(msgGameOver)
	ErrNotLegitWord   = errors.New(msgNotLegitWord)
)

// Dynamic errors

// RoomNotFound
type RoomNotFound struct {
	RoomCode string
}

func (e *RoomNotFound) Error() string {
	return fmt.Sprintf(msgRoomNotFound, e.RoomCode)
}

// NickAlreadyExists
type NickAlreadyExists struct {
	Nick     string
	RoomCode string
}

func (e *NickAlreadyExists) Error() string {
	return fmt.Sprintf(msgNickAlreadyExists, e.Nick, e.RoomCode)
}

// NickNotFound
type NickNotFound struct {
	Nick     string
	RoomCode string
}

func (e *NickNotFound) Error() string {
	return fmt.Sprintf(msgNickNotFound, e.Nick, e.RoomCode)
}
