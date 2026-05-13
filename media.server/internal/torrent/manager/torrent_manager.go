package manager

import (
	"sync"

	"fileserver/internal/torrent/model"
)

type TorrentManager struct {
	mu sync.RWMutex

	torrents map[string]*model.Torrent
}

func NewTorrentManager() *TorrentManager {

	return &TorrentManager{
		torrents: make(map[string]*model.Torrent),
	}
}

func (m *TorrentManager) Add(
	torrent *model.Torrent,
) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.torrents[torrent.ID.String()] = torrent
}

func (m *TorrentManager) Get(
	id string,
) (*model.Torrent, bool) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	torrent, exists := m.torrents[id]

	return torrent, exists
}

func (m *TorrentManager) List() []*model.Torrent {

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*model.Torrent, 0)

	for _, torrent := range m.torrents {
		result = append(result, torrent)
	}

	return result
}

func (m *TorrentManager) Remove(
	id string,
) (*model.Torrent, bool) {

	m.mu.Lock()
	defer m.mu.Unlock()

	torrent, exists := m.torrents[id]
	if !exists {
		return nil, false
	}

	delete(m.torrents, id)

	return torrent, true
}
