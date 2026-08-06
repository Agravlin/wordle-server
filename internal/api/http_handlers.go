package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/agravlin/wordle-server/internal/errs"
	"github.com/agravlin/wordle-server/internal/service"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	svc    *service.GameService
	logger *slog.Logger
}

func NewHandler(svc *service.GameService, l *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: l,
	}
}

// POST /api/create
func (h *Handler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := h.svc.CreateRoom(req.Nickname); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// POST /api/join
func (h *Handler) HandleJoinRoomValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RoomID   string `json:"room_id"`
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.ValidateJoinRequest(req.RoomID, req.Nickname); err != nil {
		var nickErr *errs.NickAlreadyExists
		if errors.As(err, &nickErr) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
