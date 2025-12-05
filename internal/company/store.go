package company

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"dice-service/internal/characters"
	"dice-service/internal/monsters"
)

// Store интерфейс для хранилища компаний
type Store interface {
	Create(c Company) (Company, error)
	Get(id string) (Company, error)
	List() []CompanySummary
	Update(c Company) error
	Delete(id string) error
	AddCharacter(companyID string, char characters.CharacterSheet) error
	RemoveCharacter(companyID, characterID string) error
	AddMonster(companyID string, mon monsters.Monster) error
	RemoveMonster(companyID, monsterID string) error
}

// MemoryStore хранилище компаний в памяти
type MemoryStore struct {
	mu        sync.RWMutex
	companies map[string]Company
}

// NewMemoryStore создаёт новое хранилище компаний
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		companies: make(map[string]Company),
	}
}

// Create создаёт новую компанию
func (s *MemoryStore) Create(c Company) (Company, error) {
	if err := c.Validate(); err != nil {
		return Company{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем уникальность названия
	for _, existing := range s.companies {
		if strings.EqualFold(existing.Name, c.Name) {
			return Company{}, ErrDuplicateName
		}
	}

	if c.ID == "" {
		c.ID = generateID()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if c.Characters == nil {
		c.Characters = []characters.CharacterSheet{}
	}
	if c.Monsters == nil {
		c.Monsters = []monsters.Monster{}
	}

	s.companies[c.ID] = c
	return c, nil
}

// Get получает компанию по ID
func (s *MemoryStore) Get(id string) (Company, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.companies[id]
	if !ok {
		return Company{}, ErrNotFound
	}
	return c, nil
}

// List возвращает список всех компаний (краткая информация)
func (s *MemoryStore) List() []CompanySummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]CompanySummary, 0, len(s.companies))
	for _, c := range s.companies {
		result = append(result, c.ToSummary())
	}
	return result
}

// Update обновляет компанию
func (s *MemoryStore) Update(c Company) error {
	if err := c.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.companies[c.ID]
	if !ok {
		return ErrNotFound
	}

	// Проверяем уникальность названия, если оно изменилось (исключая текущую компанию)
	if !strings.EqualFold(existing.Name, c.Name) {
		for id, comp := range s.companies {
			if id != c.ID && strings.EqualFold(comp.Name, c.Name) {
				return ErrDuplicateName
			}
		}
	}

	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now()
	// Сохраняем персонажей и монстров если не переданы
	if c.Characters == nil {
		c.Characters = existing.Characters
	}
	if c.Monsters == nil {
		c.Monsters = existing.Monsters
	}

	s.companies[c.ID] = c
	return nil
}

// Delete удаляет компанию
func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.companies[id]; !ok {
		return ErrNotFound
	}

	delete(s.companies, id)
	return nil
}

// AddCharacter добавляет персонажа в компанию
func (s *MemoryStore) AddCharacter(companyID string, char characters.CharacterSheet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.companies[companyID]
	if !ok {
		return ErrNotFound
	}

	// Проверяем, нет ли уже такого персонажа
	for _, existing := range c.Characters {
		if existing.ID == char.ID {
			return nil // уже есть
		}
	}

	c.Characters = append(c.Characters, char)
	c.UpdatedAt = time.Now()
	s.companies[companyID] = c
	return nil
}

// RemoveCharacter удаляет персонажа из компании
func (s *MemoryStore) RemoveCharacter(companyID, characterID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.companies[companyID]
	if !ok {
		return ErrNotFound
	}

	for i, char := range c.Characters {
		if char.ID == characterID {
			c.Characters = append(c.Characters[:i], c.Characters[i+1:]...)
			c.UpdatedAt = time.Now()
			s.companies[companyID] = c
			return nil
		}
	}

	return nil // персонаж не найден, но это не ошибка
}

// AddMonster добавляет монстра в компанию
func (s *MemoryStore) AddMonster(companyID string, mon monsters.Monster) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.companies[companyID]
	if !ok {
		return ErrNotFound
	}

	// Проверяем, нет ли уже такого монстра
	for _, existing := range c.Monsters {
		if existing.ID == mon.ID {
			return nil // уже есть
		}
	}

	c.Monsters = append(c.Monsters, mon)
	c.UpdatedAt = time.Now()
	s.companies[companyID] = c
	return nil
}

// RemoveMonster удаляет монстра из компании
func (s *MemoryStore) RemoveMonster(companyID, monsterID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.companies[companyID]
	if !ok {
		return ErrNotFound
	}

	for i, mon := range c.Monsters {
		if mon.ID == monsterID {
			c.Monsters = append(c.Monsters[:i], c.Monsters[i+1:]...)
			c.UpdatedAt = time.Now()
			s.companies[companyID] = c
			return nil
		}
	}

	return nil // монстр не найден, но это не ошибка
}

// generateID генерирует уникальный ID используя crypto/rand
func generateID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(fmt.Errorf("failed to generate id: %w", err))
	}
	return hex.EncodeToString(buf[:])
}

