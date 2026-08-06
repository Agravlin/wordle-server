package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/agravlin/wordle-server/internal/api"
	"github.com/agravlin/wordle-server/internal/service"
	"github.com/agravlin/wordle-server/internal/ws"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	manager := ws.NewManager(logger)
	gameSvc := service.NewGameService(manager, logger)

	h := api.NewHandler(gameSvc, logger)
	http.HandleFunc("POST /api/create", h.HandleCreateRoom)
	http.HandleFunc("POST /api/join", h.HandleJoinRoomValidation)

	http.HandleFunc("/ws", h.ServeWS)

	log.Println("WebSocket server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("Server failed", slog.String("error", err.Error()))
	}
}
