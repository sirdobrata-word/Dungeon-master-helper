package monsters

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type Store interface {
	Create(monster Monster) (Monster, error)
	Get(id string) (Monster, error)
	List() []Monster
	Update(id string, monster Monster) (Monster, error)
	Delete(id string) error
}

type MemoryStore struct {
	mu    sync.RWMutex
	byID  map[string]Monster
	order []string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		byID: make(map[string]Monster),
	}
}

func (s *MemoryStore) Create(monster Monster) (Monster, error) {
	if err := monster.Validate(); err != nil {
		return Monster{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if monster.ID == "" {
		monster.ID = generateID()
	}
	if _, exists := s.byID[monster.ID]; exists {
		return Monster{}, fmt.Errorf("monster with id %s already exists", monster.ID)
	}

	s.byID[monster.ID] = monster
	s.order = append(s.order, monster.ID)
	return monster, nil
}

func (s *MemoryStore) Get(id string) (Monster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	monster, ok := s.byID[id]
	if !ok {
		return Monster{}, ErrNotFound
	}
	return monster, nil
}

func (s *MemoryStore) List() []Monster {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Monster, 0, len(s.byID))
	for _, id := range s.order {
		result = append(result, s.byID[id])
	}
	return result
}

func (s *MemoryStore) Update(id string, monster Monster) (Monster, error) {
	if err := monster.Validate(); err != nil {
		return Monster{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[id]; !ok {
		return Monster{}, ErrNotFound
	}

	monster.ID = id
	s.byID[id] = monster
	return monster, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[id]; !ok {
		return ErrNotFound
	}

	delete(s.byID, id)
	// Удаляем из порядка
	for i, v := range s.order {
		if v == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	return nil
}

func generateID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(fmt.Errorf("failed to generate id: %w", err))
	}
	return hex.EncodeToString(buf[:])
}



