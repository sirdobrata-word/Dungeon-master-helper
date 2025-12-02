package monsters

// GetSampleMonsters возвращает несколько примеров монстров для начальной загрузки
func GetSampleMonsters() []Monster {
	return []Monster{
		{
			Name:            "Красный Дракон",
			Type:            "Dragon",
			Size:            "Huge",
			Alignment:       "Chaotic Evil",
			ArmorClass:       22,
			HitPoints:       256,
			HitDice:          "19d12+133",
			Speed:            "40 ft., climb 40 ft., fly 80 ft.",
			AbilityScores:   map[string]int{"STR": 27, "DEX": 10, "CON": 25, "INT": 16, "WIS": 13, "CHA": 21},
			ChallengeRating: "17 (18,000 XP)",
			Description:      "Древний красный дракон - одно из самых могущественных существ в мире. Его огненное дыхание может испепелить целые армии.",
		},
		{
			Name:            "Лич",
			Type:            "Undead",
			Size:            "Medium",
			Alignment:       "Any Evil",
			ArmorClass:       17,
			HitPoints:       135,
			HitDice:          "18d8+54",
			Speed:            "30 ft.",
			AbilityScores:   map[string]int{"STR": 11, "DEX": 16, "CON": 16, "INT": 20, "WIS": 14, "CHA": 16},
			ChallengeRating: "21 (33,000 XP)",
			Description:      "Бессмертный некромант, обменявший свою душу на вечную жизнь. Обладает огромной магической силой и может воскрешать мертвых.",
		},
		{
			Name:            "Бехолдер",
			Type:            "Aberration",
			Size:            "Large",
			Alignment:       "Lawful Evil",
			ArmorClass:       18,
			HitPoints:       180,
			HitDice:          "19d10+76",
			Speed:            "0 ft., fly 20 ft. (hover)",
			AbilityScores:   map[string]int{"STR": 10, "DEX": 14, "CON": 18, "INT": 17, "WIS": 15, "CHA": 17},
			ChallengeRating: "13 (10,000 XP)",
			Description:      "Плавающий глаз с множеством щупалец. Каждое щупальце может использовать магический луч. Крайне параноидальное существо.",
		},
		{
			Name:            "Вампир",
			Type:            "Undead",
			Size:            "Medium",
			Alignment:       "Lawful Evil",
			ArmorClass:       16,
			HitPoints:       144,
			HitDice:          "17d8+68",
			Speed:            "30 ft.",
			AbilityScores:   map[string]int{"STR": 18, "DEX": 18, "CON": 18, "INT": 17, "WIS": 15, "CHA": 18},
			ChallengeRating: "13 (10,000 XP)",
			Description:      "Бессмертный вампир, питающийся кровью живых. Обладает способностью превращаться в туман, контролировать разум и регенерировать.",
		},
		{
			Name:            "Тролль",
			Type:            "Giant",
			Size:            "Large",
			Alignment:       "Chaotic Evil",
			ArmorClass:       15,
			HitPoints:       84,
			HitDice:          "8d10+40",
			Speed:            "30 ft.",
			AbilityScores:   map[string]int{"STR": 18, "DEX": 13, "CON": 20, "INT": 7, "WIS": 9, "CHA": 7},
			ChallengeRating: "5 (1,800 XP)",
			Description:      "Большое, злобное существо с мощной регенерацией. Может восстанавливать потерянные конечности. Слабость к огню и кислоте.",
		},
	}
}





