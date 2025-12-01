package characters

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type Store interface {
	Create(sheet CharacterSheet) (CharacterSheet, error)
	Get(id string) (CharacterSheet, error)
	List() []CharacterSheet
	Update(id string, sheet CharacterSheet) (CharacterSheet, error)
}

type MemoryStore struct {
	mu    sync.RWMutex
	byID  map[string]CharacterSheet
	order []string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		byID: make(map[string]CharacterSheet),
	}
}

func (s *MemoryStore) Create(sheet CharacterSheet) (CharacterSheet, error) {
	if err := sheet.Validate(); err != nil {
		return CharacterSheet{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if sheet.ID == "" {
		sheet.ID = generateID()
	}
	if _, exists := s.byID[sheet.ID]; exists {
		return CharacterSheet{}, fmt.Errorf("character with id %s already exists", sheet.ID)
	}

	s.byID[sheet.ID] = sheet
	s.order = append(s.order, sheet.ID)
	return sheet, nil
}

func (s *MemoryStore) Get(id string) (CharacterSheet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sheet, ok := s.byID[id]
	if !ok {
		return CharacterSheet{}, ErrNotFound
	}
	return sheet, nil
}

func (s *MemoryStore) List() []CharacterSheet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]CharacterSheet, 0, len(s.byID))
	for _, id := range s.order {
		result = append(result, s.byID[id])
	}
	return result
}

// Update replaces an existing character sheet by id.
func (s *MemoryStore) Update(id string, sheet CharacterSheet) (CharacterSheet, error) {
	if err := sheet.Validate(); err != nil {
		return CharacterSheet{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[id]; !ok {
		return CharacterSheet{}, ErrNotFound
	}

	sheet.ID = id
	s.byID[id] = sheet
	return sheet, nil
}

func generateID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(fmt.Errorf("failed to generate id: %w", err))
	}
	return hex.EncodeToString(buf[:])
}
