package game

import "fmt"

type State int

const (
	StateWaiting  State = iota // 0
	StatePlaying               // 1
	StateFinished              // 2
)

func (s State) String() string {
	switch s {
	case StateWaiting:
		return "WAITING"
	case StatePlaying:
		return "PLAYING"
	case StateFinished:
		return "FINISHED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

func (s State) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", s.String()), nil // []byte(fmt.Sprintf("%q", s.String()))
}
