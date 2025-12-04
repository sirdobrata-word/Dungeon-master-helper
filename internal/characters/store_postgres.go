package characters

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// PostgresStore реализует Store для персонажей в PostgreSQL.
// Схема: таблица characters с отдельными колонками для поиска
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore создаёт хранилище персонажей в PostgreSQL и гарантирует,
// что таблица существует.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	const createTable = `
CREATE TABLE IF NOT EXISTS characters (
	id   TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	class TEXT NOT NULL,
	race TEXT NOT NULL,
	level INTEGER NOT NULL,
	data JSONB NOT NULL
);`

	if _, err := db.Exec(createTable); err != nil {
		panic(fmt.Errorf("failed to create characters table: %w", err))
	}

	// Миграция: добавляем колонки если их нет
	const migrateTable = `
DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'characters' AND column_name = 'name') THEN
		ALTER TABLE characters ADD COLUMN name TEXT;
		ALTER TABLE characters ADD COLUMN class TEXT;
		ALTER TABLE characters ADD COLUMN race TEXT;
		ALTER TABLE characters ADD COLUMN level INTEGER;
		
		-- Заполняем новые колонки из JSON
		UPDATE characters SET 
			name = data->>'name',
			class = data->>'class',
			race = data->>'race',
			level = (data->>'level')::INTEGER
		WHERE name IS NULL;
		
		-- Делаем колонки NOT NULL
		ALTER TABLE characters ALTER COLUMN name SET NOT NULL;
		ALTER TABLE characters ALTER COLUMN class SET NOT NULL;
		ALTER TABLE characters ALTER COLUMN race SET NOT NULL;
		ALTER TABLE characters ALTER COLUMN level SET NOT NULL;
		
		-- Создаём индексы для поиска
		CREATE INDEX IF NOT EXISTS idx_characters_name ON characters(name);
		CREATE INDEX IF NOT EXISTS idx_characters_class ON characters(class);
		CREATE INDEX IF NOT EXISTS idx_characters_race ON characters(race);
		CREATE INDEX IF NOT EXISTS idx_characters_level ON characters(level);
	END IF;
END $$;`

	if _, err := db.Exec(migrateTable); err != nil {
		panic(fmt.Errorf("failed to migrate characters table: %w", err))
	}

	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(sheet CharacterSheet) (CharacterSheet, error) {
	if err := sheet.Validate(); err != nil {
		return CharacterSheet{}, err
	}

	if sheet.ID == "" {
		sheet.ID = generateID()
	}

	data, err := json.Marshal(sheet)
	if err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to marshal character: %w", err)
	}

	const insertQuery = `INSERT INTO characters (id, name, class, race, level, data) VALUES ($1, $2, $3, $4, $5, $6::jsonb);`
	if _, err := s.db.Exec(insertQuery, sheet.ID, sheet.Name, sheet.Class, sheet.Race, sheet.Level, data); err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to insert character: %w", err)
	}

	return sheet, nil
}

func (s *PostgresStore) Get(id string) (CharacterSheet, error) {
	const selectQuery = `SELECT data FROM characters WHERE id = $1;`

	var raw []byte
	err := s.db.QueryRow(selectQuery, id).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CharacterSheet{}, ErrNotFound
		}
		return CharacterSheet{}, fmt.Errorf("failed to get character: %w", err)
	}

	var sheet CharacterSheet
	if err := json.Unmarshal(raw, &sheet); err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to unmarshal character: %w", err)
	}

	return sheet, nil
}

func (s *PostgresStore) List() []CharacterSheet {
	const listQuery = `SELECT data FROM characters ORDER BY id;`

	rows, err := s.db.Query(listQuery)
	if err != nil {
		// В случае ошибки возвращаем пустой список, а не паникуем
		return []CharacterSheet{}
	}
	defer rows.Close()

	var result []CharacterSheet
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var sheet CharacterSheet
		if err := json.Unmarshal(raw, &sheet); err != nil {
			continue
		}
		result = append(result, sheet)
	}
	return result
}

func (s *PostgresStore) Update(id string, sheet CharacterSheet) (CharacterSheet, error) {
	if err := sheet.Validate(); err != nil {
		return CharacterSheet{}, err
	}

	sheet.ID = id

	data, err := json.Marshal(sheet)
	if err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to marshal character: %w", err)
	}

	const updateQuery = `UPDATE characters SET name = $2, class = $3, race = $4, level = $5, data = $6::jsonb WHERE id = $1;`
	res, err := s.db.Exec(updateQuery, id, sheet.Name, sheet.Class, sheet.Race, sheet.Level, data)
	if err != nil {
		return CharacterSheet{}, fmt.Errorf("failed to update character: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return CharacterSheet{}, ErrNotFound
	}

	return sheet, nil
}



