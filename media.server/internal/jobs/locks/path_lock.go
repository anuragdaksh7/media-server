package locks

import "sync"

type PathLockManager struct {
	mu sync.Mutex

	locks map[string]*sync.RWMutex
}

func NewPathLockManager() *PathLockManager {

	return &PathLockManager{
		locks: make(map[string]*sync.RWMutex),
	}
}

func (m *PathLockManager) getLock(
	path string,
) *sync.RWMutex {

	m.mu.Lock()
	defer m.mu.Unlock()

	lock, exists := m.locks[path]

	if !exists {

		lock = &sync.RWMutex{}

		m.locks[path] = lock
	}

	return lock
}

func (m *PathLockManager) Lock(
	path string,
) {

	lock := m.getLock(path)

	lock.Lock()
}

func (m *PathLockManager) Unlock(
	path string,
) {

	lock := m.getLock(path)

	lock.Unlock()
}

func (m *PathLockManager) RLock(
	path string,
) {

	lock := m.getLock(path)

	lock.RLock()
}

func (m *PathLockManager) RUnlock(
	path string,
) {

	lock := m.getLock(path)

	lock.RUnlock()
}
