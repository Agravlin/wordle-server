package ws

import "sync"

type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (m *Manager) AddRoom(r *Room) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rooms[r.ID] = r
}

func (m *Manager) GetRoom(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.rooms[id]
	return r, exists
}
