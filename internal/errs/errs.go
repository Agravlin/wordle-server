package errs

import (
	"errors"
)

const (
	msgLengthMismatch = "Guess length does not match target length"
	msgInvalidChar    = "Guess contains invalid non-ASCII characters"
)

// Static errors
var (
	ErrLengthMismatch = errors.New(msgLengthMismatch)
	ErrInvalidChar    = errors.New(msgInvalidChar)
)

/* // Dynamic errors
type RoomNotFoundError struct {
	RoomCode string
}

func (e *RoomNotFoundError) Error() string {
	return fmt.Sprintf(msgRoomNotFound, e.RoomCode)
} */
