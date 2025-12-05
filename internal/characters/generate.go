package characters

import (
	"fmt"
	"sort"
	"strings"

	"dice-service/internal/dice"
)

// GenerateCharacterSheet автоматически создаёт лист персонажа:
// - бросает характеристики 4d6 drop lowest
// - расставляет их с учётом приоритетов класса (более высокие значения идут в приоритетные характеристики)
// - автоматически добавляет навыки класса в соответствии с уровнем
// - считает модификаторы, проф.бонус и базовые боевые поля.
func GenerateCharacterSheet(name, class, race, background, alignment string, level int, skills []string) (CharacterSheet, error) {
	if level <= 0 {
		level = 1
	}

	scores, err := rollSixAbilities()
	if err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to generate abilities: %w", err)
	}

	// Сортируем значения по убыванию
	sort.Sort(sort.Reverse(sort.IntSlice(scores[:])))

	// Получаем приоритетные характеристики для класса
	priorities := getClassPriorityAbilities(class)

	// Создаём мапу для отслеживания использованных значений
	used := make(map[int]bool)
	abilities := AbilityScores{}

	// Сначала присваиваем значения приоритетным характеристикам
	for _, priority := range priorities {
		for i, score := range scores {
			if !used[i] {
				switch priority {
				case "strength":
					abilities.Strength = score
				case "dexterity":
					abilities.Dexterity = score
				case "constitution":
					abilities.Constitution = score
				case "intelligence":
					abilities.Intelligence = score
				case "wisdom":
					abilities.Wisdom = score
				case "charisma":
					abilities.Charisma = score
				}
				used[i] = true
				break
			}
		}
	}

	// Затем присваиваем оставшиеся значения остальным характеристикам
	allAbilities := []struct {
		name  string
		value *int
	}{
		{"strength", &abilities.Strength},
		{"dexterity", &abilities.Dexterity},
		{"constitution", &abilities.Constitution},
		{"intelligence", &abilities.Intelligence},
		{"wisdom", &abilities.Wisdom},
		{"charisma", &abilities.Charisma},
	}

	for _, ab := range allAbilities {
		// Пропускаем уже заполненные характеристики (приоритетные)
		isPriority := false
		for _, p := range priorities {
			if p == ab.name {
				isPriority = true
				break
			}
		}
		if isPriority {
			continue
		}

		// Находим первое неиспользованное значение
		for i, score := range scores {
			if !used[i] {
				*ab.value = score
				used[i] = true
				break
			}
		}
	}

	dexMod := abilityModifier(abilities.Dexterity)
	conMod := abilityModifier(abilities.Constitution)

	prof := proficiencyBonus(level)
	maxHP := max(1, (8+conMod)*level) // упрощённая модель хитов

	// Получаем навыки класса в соответствии с уровнем
	classSkills := getClassSkills(class, level)
	
	// Объединяем навыки класса с переданными навыками (убираем дубликаты)
	allSkills := mergeSkills(classSkills, skills)

	sheet := CharacterSheet{
		Name:               name,
		Class:              class,
		Race:               race,
		Background:         background,
		Level:              level,
		Alignment:          alignment,
		AbilityScores:      abilities,
		ProficiencyBonus:   prof,
		ArmorClass:         10 + dexMod,
		Speed:              30,
		Initiative:         dexMod,
		MaxHitPoints:       maxHP,
		CurrentHitPoints:   maxHP,
		TemporaryHitPoints: 0,
		Skills:             allSkills,
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

// getClassPriorityAbilities возвращает список приоритетных характеристик для класса
// в порядке убывания приоритета (первая - самая важная)
func getClassPriorityAbilities(class string) []string {
	classLower := strings.ToLower(strings.TrimSpace(class))
	
	switch classLower {
	case "варвар", "barbarian":
		return []string{"strength", "constitution", "dexterity"}
	case "бард", "bard":
		return []string{"charisma", "dexterity", "constitution"}
	case "жрец", "cleric":
		return []string{"wisdom", "constitution", "strength"}
	case "друид", "druid":
		return []string{"wisdom", "constitution", "dexterity"}
	case "воин", "fighter":
		return []string{"strength", "constitution", "dexterity"}
	case "монах", "monk":
		return []string{"dexterity", "wisdom", "constitution"}
	case "паладин", "paladin":
		return []string{"strength", "charisma", "constitution"}
	case "следопыт", "ranger":
		return []string{"dexterity", "wisdom", "constitution"}
	case "плут", "rogue":
		return []string{"dexterity", "intelligence", "constitution"}
	case "чародей", "sorcerer":
		return []string{"charisma", "constitution", "dexterity"}
	case "колдун", "warlock":
		return []string{"charisma", "constitution", "dexterity"}
	case "волшебник", "wizard":
		return []string{"intelligence", "constitution", "dexterity"}
	case "изобретатель", "artificer":
		return []string{"intelligence", "constitution", "dexterity"}
	default:
		// Если класс не распознан, возвращаем стандартный порядок
		return []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"}
	}
}

// getClassSkills возвращает навыки класса в соответствии с уровнем
func getClassSkills(class string, level int) []string {
	classLower := strings.ToLower(strings.TrimSpace(class))
	
	// Базовые навыки для каждого класса (на 1 уровне)
	baseSkills := map[string][]string{
		"варвар":   {"Атлетика", "Выживание", "Запугивание", "Природа", "Внимательность", "Обращение с животными"},
		"barbarian": {"Атлетика", "Выживание", "Запугивание", "Природа", "Внимательность", "Обращение с животными"},
		
		"бард":     {"Акробатика", "Атлетика", "Обман", "История", "Проницательность", "Запугивание", "Расследование", "Медицина", "Природа", "Внимательность", "Выступление", "Убеждение", "Религия", "Ловкость рук", "Скрытность"},
		"bard":     {"Акробатика", "Атлетика", "Обман", "История", "Проницательность", "Запугивание", "Расследование", "Медицина", "Природа", "Внимательность", "Выступление", "Убеждение", "Религия", "Ловкость рук", "Скрытность"},
		
		"жрец":     {"История", "Медицина", "Проницательность", "Религия", "Убеждение"},
		"cleric":   {"История", "Медицина", "Проницательность", "Религия", "Убеждение"},
		
		"друид":    {"Магия", "Атлетика", "Обращение с животными", "История", "Проницательность", "Медицина", "Природа", "Внимательность", "Религия", "Выживание"},
		"druid":    {"Магия", "Атлетика", "Обращение с животными", "История", "Проницательность", "Медицина", "Природа", "Внимательность", "Религия", "Выживание"},
		
		"воин":     {"Акробатика", "Атлетика", "История", "Проницательность", "Запугивание", "Внимательность", "Выживание"},
		"fighter":  {"Акробатика", "Атлетика", "История", "Проницательность", "Запугивание", "Внимательность", "Выживание"},
		
		"монах":    {"Акробатика", "Атлетика", "История", "Проницательность", "Религия", "Скрытность"},
		"monk":     {"Акробатика", "Атлетика", "История", "Проницательность", "Религия", "Скрытность"},
		
		"паладин":  {"Атлетика", "Проницательность", "Запугивание", "Медицина", "Убеждение", "Религия"},
		"paladin":  {"Атлетика", "Проницательность", "Запугивание", "Медицина", "Убеждение", "Религия"},
		
		"следопыт": {"Атлетика", "Обращение с животными", "Проницательность", "Расследование", "Природа", "Внимательность", "Медицина", "Выживание", "Скрытность"},
		"ranger":   {"Атлетика", "Обращение с животными", "Проницательность", "Расследование", "Природа", "Внимательность", "Медицина", "Выживание", "Скрытность"},
		
		"плут":     {"Акробатика", "Атлетика", "Обман", "Проницательность", "Запугивание", "Расследование", "Внимательность", "Выступление", "Убеждение", "Ловкость рук", "Скрытность"},
		"rogue":    {"Акробатика", "Атлетика", "Обман", "Проницательность", "Запугивание", "Расследование", "Внимательность", "Выступление", "Убеждение", "Ловкость рук", "Скрытность"},
		
		"чародей":  {"Проницательность", "Запугивание", "Убеждение", "Религия", "Обман"},
		"sorcerer": {"Проницательность", "Запугивание", "Убеждение", "Религия", "Обман"},
		
		"колдун":   {"Магия", "Обман", "История", "Запугивание", "Расследование", "Природа", "Религия"},
		"warlock":  {"Магия", "Обман", "История", "Запугивание", "Расследование", "Природа", "Религия"},
		
		"волшебник": {"Магия", "История", "Проницательность", "Расследование", "Медицина", "Религия"},
		"wizard":   {"Магия", "История", "Проницательность", "Расследование", "Медицина", "Религия"},
		
		"изобретатель": {"Магия", "История", "Расследование", "Медицина", "Природа", "Внимательность"},
		"artificer": {"Магия", "История", "Расследование", "Медицина", "Природа", "Внимательность"},
	}
	
	// Количество навыков, которые персонаж выбирает на 1 уровне
	skillsCount := map[string]int{
		"варвар": 2, "barbarian": 2,
		"бард": 3, "bard": 3,
		"жрец": 2, "cleric": 2,
		"друид": 2, "druid": 2,
		"воин": 2, "fighter": 2,
		"монах": 2, "monk": 2,
		"паладин": 2, "paladin": 2,
		"следопыт": 3, "ranger": 3,
		"плут": 4, "rogue": 4,
		"чародей": 2, "sorcerer": 2,
		"колдун": 2, "warlock": 2,
		"волшебник": 2, "wizard": 2,
		"изобретатель": 2, "artificer": 2,
	}
	
	available, ok := baseSkills[classLower]
	if !ok {
		return []string{} // Если класс не найден, возвращаем пустой список
	}
	
	count, ok := skillsCount[classLower]
	if !ok {
		count = 2 // По умолчанию 2 навыка
	}
	
	// Выбираем первые N навыков из доступных (упрощённая логика)
	// В реальной игре игрок выбирает, но здесь мы автоматически выбираем первые
	selected := []string{}
	if len(available) > 0 {
		maxSelect := count
		if maxSelect > len(available) {
			maxSelect = len(available)
		}
		selected = available[:maxSelect]
	}
	
	// Некоторые классы получают дополнительные навыки на более высоких уровнях
	// Например, Бард получает Expertise на 3 уровне
	if level >= 3 {
		switch classLower {
		case "бард", "bard":
			// Бард получает Expertise (улучшение навыков), но не новые навыки
		case "плут", "rogue":
			// Плут получает Expertise на 1 уровне, но дополнительные навыки на 6 уровне
			if level >= 6 && len(selected) < len(available) {
				// Добавляем ещё один навык из доступных
				for _, skill := range available {
					hasSkill := false
					for _, s := range selected {
						if s == skill {
							hasSkill = true
							break
						}
					}
					if !hasSkill {
						selected = append(selected, skill)
						break
					}
				}
			}
		}
	}
	
	return selected
}

// mergeSkills объединяет навыки класса с переданными навыками, убирая дубликаты
func mergeSkills(classSkills []string, userSkills []string) []string {
	skillMap := make(map[string]bool)
	result := []string{}
	
	// Сначала добавляем навыки класса
	for _, skill := range classSkills {
		skillLower := strings.ToLower(strings.TrimSpace(skill))
		if skillLower != "" && !skillMap[skillLower] {
			skillMap[skillLower] = true
			result = append(result, skill)
		}
	}
	
	// Затем добавляем пользовательские навыки (если их ещё нет)
	for _, skill := range userSkills {
		skillLower := strings.ToLower(strings.TrimSpace(skill))
		if skillLower != "" && !skillMap[skillLower] {
			skillMap[skillLower] = true
			result = append(result, skill)
		}
	}
	
	return result
}

// LevelUp повышает уровень персонажа на 1 и пересчитывает все зависимые параметры
func LevelUp(sheet CharacterSheet) (CharacterSheet, error) {
	if sheet.Level >= 20 {
		return CharacterSheet{}, fmt.Errorf("character is already at maximum level (20)")
	}

	newLevel := sheet.Level + 1
	
	// Пересчитываем бонус мастерства
	newProf := proficiencyBonus(newLevel)
	
	// Вычисляем новые максимальные HP
	// В D&D 5e при повышении уровня добавляются хиты: Hit Die + модификатор телосложения
	// Упрощённая модель: среднее значение Hit Die + модификатор CON
	conMod := abilityModifier(sheet.AbilityScores.Constitution)
	
	// Определяем Hit Die для класса
	hitDie := getClassHitDie(sheet.Class)
	
	// Добавляем новые хиты (среднее значение Hit Die + модификатор CON, минимум 1)
	hpGain := max(1, (hitDie/2+1)+conMod) // среднее значение (например, для d8 это 5)
	newMaxHP := sheet.MaxHitPoints + hpGain
	
	// Увеличиваем текущие хиты на столько же (или до максимума)
	newCurrentHP := sheet.CurrentHitPoints + hpGain
	if newCurrentHP > newMaxHP {
		newCurrentHP = newMaxHP
	}
	
	// Получаем навыки для нового уровня
	newClassSkills := getClassSkills(sheet.Class, newLevel)
	
	// Объединяем старые навыки с новыми (если появились новые)
	allSkills := mergeSkills(newClassSkills, sheet.Skills)
	
	// Обновляем лист персонажа
	sheet.Level = newLevel
	sheet.ProficiencyBonus = newProf
	sheet.MaxHitPoints = newMaxHP
	sheet.CurrentHitPoints = newCurrentHP
	sheet.Skills = allSkills
	
	return sheet, nil
}

// getClassHitDie возвращает размер кости хитов для класса
func getClassHitDie(class string) int {
	classLower := strings.ToLower(strings.TrimSpace(class))
	
	switch classLower {
	case "варвар", "barbarian":
		return 12 // d12
	case "воин", "fighter", "паладин", "paladin", "следопыт", "ranger":
		return 10 // d10
	case "бард", "bard", "жрец", "cleric", "друид", "druid", "монах", "monk", "плут", "rogue", "колдун", "warlock":
		return 8 // d8
	case "чародей", "sorcerer", "волшебник", "wizard":
		return 6 // d6
	case "изобретатель", "artificer":
		return 8 // d8
	default:
		return 8 // по умолчанию d8
	}
}
