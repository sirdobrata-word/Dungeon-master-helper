package company

import (
	"errors"
	"time"

	"dice-service/internal/characters"
	"dice-service/internal/monsters"
)

var ErrNotFound = errors.New("company not found")

// Company представляет склад/кампанию, объединяющую персонажей и монстров
type Company struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	CreatedAt   time.Time                    `json:"createdAt"`
	UpdatedAt   time.Time                    `json:"updatedAt"`
	Characters  []characters.CharacterSheet  `json:"characters"`  // персонажи в компании
	Monsters    []monsters.Monster           `json:"monsters"`    // монстры в компании
}

// CompanySummary краткая информация о компании для списка
type CompanySummary struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CharacterCount  int       `json:"characterCount"`
	MonsterCount    int       `json:"monsterCount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ToSummary конвертирует Company в CompanySummary
func (c Company) ToSummary() CompanySummary {
	return CompanySummary{
		ID:             c.ID,
		Name:           c.Name,
		Description:    c.Description,
		CharacterCount: len(c.Characters),
		MonsterCount:   len(c.Monsters),
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func (c Company) Validate() error {
	if c.Name == "" {
		return errors.New("company name is required")
	}
	return nil
}

// AddCharacterRequest запрос на добавление персонажа в компанию
type AddCharacterRequest struct {
	CharacterID string `json:"characterId"`
}

// AddMonsterRequest запрос на добавление монстра в компанию
type AddMonsterRequest struct {
	MonsterID string `json:"monsterId"`
}





