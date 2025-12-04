package monsters

import "errors"

var ErrNotFound = errors.New("monster not found")

type Monster struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`            // например: "Beast", "Undead", "Dragon"
	Size            string            `json:"size"`            // Tiny, Small, Medium, Large, Huge, Gargantuan
	Alignment       string            `json:"alignment"`       // например: "Chaotic Evil"
	ArmorClass      int               `json:"armorClass"`
	HitPoints       int               `json:"hitPoints"`
	HitDice          string            `json:"hitDice"`         // например: "10d8+30"
	Speed            string            `json:"speed"`           // например: "30 ft., fly 60 ft."
	AbilityScores    map[string]int    `json:"abilityScores"`  // STR, DEX, CON, INT, WIS, CHA
	Skills           map[string]int    `json:"skills"`          // например: {"Perception": 5, "Stealth": 4}
	SavingThrows     map[string]int    `json:"savingThrows"`   // например: {"DEX": 6, "CON": 8}
	DamageResistances []string         `json:"damageResistances"`
	DamageImmunities  []string         `json:"damageImmunities"`
	ConditionImmunities []string       `json:"conditionImmunities"`
	Senses            string            `json:"senses"`         // например: "darkvision 60 ft."
	Languages         []string         `json:"languages"`
	ChallengeRating   string            `json:"challengeRating"` // например: "5 (1,800 XP)"
	Traits            []string         `json:"traits"`          // особенности
	Actions           []string         `json:"actions"`         // действия
	LegendaryActions  []string         `json:"legendaryActions"` // легендарные действия
	Description       string           `json:"description"`
}

func (m Monster) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.Type == "" {
		return errors.New("type is required")
	}
	if m.ArmorClass < 0 || m.ArmorClass > 30 {
		return errors.New("armor class must be between 0 and 30")
	}
	if m.HitPoints < 1 {
		return errors.New("hit points must be at least 1")
	}
	return nil
}







