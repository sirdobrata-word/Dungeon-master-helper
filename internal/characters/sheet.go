package characters

import "errors"

var ErrNotFound = errors.New("character not found")

type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type CharacterSheet struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Class              string        `json:"class"`
	Race               string        `json:"race"`
	Background         string        `json:"background"`
	Level              int           `json:"level"`
	AbilityScores      AbilityScores `json:"abilityScores"`
	ProficiencyBonus   int           `json:"proficiencyBonus"`
	ArmorClass         int           `json:"armorClass"`
	Speed              int           `json:"speed"`
	Initiative         int           `json:"initiative"`
	MaxHitPoints       int           `json:"maxHitPoints"`
	CurrentHitPoints   int           `json:"currentHitPoints"`
	TemporaryHitPoints int           `json:"temporaryHitPoints"`
}

func (c CharacterSheet) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Class == "" {
		return errors.New("class is required")
	}
	if c.Level < 1 {
		return errors.New("level must be at least 1")
	}
	return nil
}
