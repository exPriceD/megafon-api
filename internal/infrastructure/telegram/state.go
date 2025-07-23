package telegram

import "sync"

type userState int

const (
	idle userState = iota
	awaitCity
	awaitPeriod
	awaitFrom
	awaitTo
)

type stateStore struct {
	mu   sync.RWMutex
	data map[int64]userState
}

func newStateStore() *stateStore {
	return &stateStore{data: make(map[int64]userState)}
}

func (s *stateStore) get(id int64) userState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[id]
}

func (s *stateStore) set(id int64, st userState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = st
}
