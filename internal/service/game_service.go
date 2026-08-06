package service

import (
	"log/slog"
	"strings"

	"github.com/agravlin/wordle-server/internal/errs"
	"github.com/agravlin/wordle-server/internal/util"
	"github.com/agravlin/wordle-server/internal/ws"
)

type GameService struct {
	manager *ws.Manager
	logger  *slog.Logger
}

func NewGameService(m *ws.Manager, l *slog.Logger) *GameService {
	return &GameService{
		manager: m,
		logger:  l,
	}
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
	roomCode := util.GenerateRoomCode()
	room := ws.NewRoom(roomCode, "", s.logger)
	s.manager.AddRoom(room)

	s.logger.Info("A new room is created",
		slog.String("room_id", roomCode),
		slog.String("creator", nick),
	)

	return room, nil
}

func (s *GameService) GetRoom(roomID string) (*ws.Room, error) {
	room, exists := s.manager.GetRoom(roomID)
	if !exists {
		return nil, &errs.RoomNotFound{RoomCode: roomID}
	}
	return room, nil
}
