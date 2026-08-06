package service

import (
	"strings"

	"github.com/agravlin/wordle-server/internal/errs"
	"github.com/agravlin/wordle-server/internal/util"
	"github.com/agravlin/wordle-server/internal/ws"
)

type GameService struct {
	manager *ws.Manager
}

func NewGameService(m *ws.Manager) *GameService {
	return &GameService{manager: m}
}

func (s *GameService) ValidateJoinRequest(roomID, nick string) error {
	nick = strings.TrimSpace(nick)
	if nick == "" {
		return errs.ErrEmptyNick
	}

	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room.IsNickTaken(nick) {
		return &errs.NickAlreadyExists{RoomCode: roomID, Nick: nick}
	}

	return nil
}

func (s *GameService) CreateRoom(nick string) (*ws.Room, error) {
	room := ws.NewRoom(util.GenerateRoomCode(), "")
	s.manager.AddRoom(room)

	return room, nil
}

func (s *GameService) GetRoom(roomID string) (*ws.Room, error) {
	room, exists := s.manager.GetRoom(roomID)
	if !exists {
		return nil, &errs.RoomNotFound{RoomCode: roomID}
	}
	return room, nil
}
