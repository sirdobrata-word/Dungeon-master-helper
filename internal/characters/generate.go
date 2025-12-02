package characters

import (
	"fmt"

	"dice-service/internal/dice"
)

// GenerateCharacterSheet автоматически создаёт лист персонажа:
// - бросает характеристики 4d6 drop lowest
// - расставляет их по порядку STR, DEX, CON, INT, WIS, CHA
// - считает модификаторы, проф.бонус и базовые боевые поля.
func GenerateCharacterSheet(name, class, race, background string, level int) (CharacterSheet, error) {
	if level <= 0 {
		level = 1
	}

	scores, err := rollSixAbilities()
	if err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to generate abilities: %w", err)
	}

	abilities := AbilityScores{
		Strength:     scores[0],
		Dexterity:    scores[1],
		Constitution: scores[2],
		Intelligence: scores[3],
		Wisdom:       scores[4],
		Charisma:     scores[5],
	}

	dexMod := abilityModifier(abilities.Dexterity)
	conMod := abilityModifier(abilities.Constitution)

	prof := proficiencyBonus(level)
	maxHP := max(1, (8+conMod)*level) // упрощённая модель хитов

	sheet := CharacterSheet{
		Name:               name,
		Class:              class,
		Race:               race,
		Background:         background,
		Level:              level,
		AbilityScores:      abilities,
		ProficiencyBonus:   prof,
		ArmorClass:         10 + dexMod,
		Speed:              30,
		Initiative:         dexMod,
		MaxHitPoints:       maxHP,
		CurrentHitPoints:   maxHP,
		TemporaryHitPoints: 0,
	}

	if err := sheet.Validate(); err != nil {
		return CharacterSheet{}, err
	}
	return sheet, nil
}

func rollSixAbilities() ([6]int, error) {
	var scores [6]int
	for i := 0; i < 6; i++ {
		score, err := rollAbilityScore()
		if err != nil {
			return [6]int{}, err
		}
		scores[i] = score
	}
	return scores, nil
}

// 4d6 drop lowest.
func rollAbilityScore() (int, error) {
	expr, err := dice.ParseExpression("4d6")
	if err != nil {
		return 0, err
	}
	result, err := dice.Roll(expr)
	if err != nil {
		return 0, err
	}
	if len(result.Rolls) != 4 {
		return 0, fmt.Errorf("unexpected roll count: %d", len(result.Rolls))
	}
	lowest := result.Rolls[0]
	sum := 0
	for _, v := range result.Rolls {
		sum += v
		if v < lowest {
			lowest = v
		}
	}
	return sum - lowest, nil
}

func abilityModifier(score int) int {
	return (score - 10) / 2
}

func proficiencyBonus(level int) int {
	switch {
	case level <= 0:
		return 2
	case level <= 4:
		return 2
	case level <= 8:
		return 3
	case level <= 12:
		return 4
	case level <= 16:
		return 5
	default:
		return 6
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
