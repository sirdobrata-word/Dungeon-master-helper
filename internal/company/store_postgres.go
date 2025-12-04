package company

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"dice-service/internal/characters"
	"dice-service/internal/monsters"
)

// PostgresStore реализует Store для компаний в PostgreSQL.
// Схема: таблица companies с отдельными колонками для поиска
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore создаёт хранилище компаний в PostgreSQL и гарантирует,
// что таблица существует.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	const createTable = `
CREATE TABLE IF NOT EXISTS companies (
	id   TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	data JSONB NOT NULL
);`

	if _, err := db.Exec(createTable); err != nil {
		panic(fmt.Errorf("failed to create companies table: %w", err))
	}

	// Миграция: добавляем колонки если их нет
	const migrateTable = `
DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'companies' AND column_name = 'name') THEN
		ALTER TABLE companies ADD COLUMN name TEXT;
		ALTER TABLE companies ADD COLUMN created_at TIMESTAMP;
		ALTER TABLE companies ADD COLUMN updated_at TIMESTAMP;
		
		-- Заполняем новые колонки из JSON
		UPDATE companies SET 
			name = data->>'name',
			created_at = (data->>'createdAt')::TIMESTAMP,
			updated_at = (data->>'updatedAt')::TIMESTAMP
		WHERE name IS NULL;
		
		-- Делаем колонки NOT NULL
		ALTER TABLE companies ALTER COLUMN name SET NOT NULL;
		ALTER TABLE companies ALTER COLUMN created_at SET NOT NULL;
		ALTER TABLE companies ALTER COLUMN updated_at SET NOT NULL;
		
		-- Создаём индексы для поиска
		CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(name);
		CREATE INDEX IF NOT EXISTS idx_companies_created_at ON companies(created_at);
	END IF;
END $$;`

	if _, err := db.Exec(migrateTable); err != nil {
		panic(fmt.Errorf("failed to migrate companies table: %w", err))
	}

	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(c Company) (Company, error) {
	if err := c.Validate(); err != nil {
		return Company{}, err
	}

	if c.ID == "" {
		c.ID = generateID()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	data, err := json.Marshal(c)
	if err != nil {
		return Company{}, fmt.Errorf("failed to marshal company: %w", err)
	}

	const insertQuery = `INSERT INTO companies (id, name, created_at, updated_at, data) VALUES ($1, $2, $3, $4, $5::jsonb);`
	if _, err := s.db.Exec(insertQuery, c.ID, c.Name, c.CreatedAt, c.UpdatedAt, data); err != nil {
		return Company{}, fmt.Errorf("failed to insert company: %w", err)
	}

	return c, nil
}

func (s *PostgresStore) Get(id string) (Company, error) {
	const selectQuery = `SELECT data FROM companies WHERE id = $1;`

	var raw []byte
	err := s.db.QueryRow(selectQuery, id).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Company{}, ErrNotFound
		}
		return Company{}, fmt.Errorf("failed to get company: %w", err)
	}

	var c Company
	if err := json.Unmarshal(raw, &c); err != nil {
		return Company{}, fmt.Errorf("failed to unmarshal company: %w", err)
	}

	return c, nil
}

func (s *PostgresStore) List() []CompanySummary {
	const listQuery = `SELECT data FROM companies ORDER BY id;`

	rows, err := s.db.Query(listQuery)
	if err != nil {
		return []CompanySummary{}
	}
	defer rows.Close()

	result := []CompanySummary{}
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var c Company
		if err := json.Unmarshal(raw, &c); err != nil {
			continue
		}
		result = append(result, c.ToSummary())
	}
	return result
}

func (s *PostgresStore) Update(c Company) error {
	if err := c.Validate(); err != nil {
		return err
	}

	// Загружаем существующую компанию, чтобы сохранить CreatedAt и, при необходимости, связанные сущности.
	existing, err := s.Get(c.ID)
	if err != nil {
		return err
	}

	if c.Characters == nil {
		c.Characters = existing.Characters
	}
	if c.Monsters == nil {
		c.Monsters = existing.Monsters
	}
	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now()

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal company: %w", err)
	}

	const updateQuery = `UPDATE companies SET name = $2, updated_at = $3, data = $4::jsonb WHERE id = $1;`
	res, err := s.db.Exec(updateQuery, c.ID, c.Name, c.UpdatedAt, data)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStore) Delete(id string) error {
	const deleteQuery = `DELETE FROM companies WHERE id = $1;`

	res, err := s.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStore) AddCharacter(companyID string, char characters.CharacterSheet) error {
	c, err := s.Get(companyID)
	if err != nil {
		return err
	}

	for _, existing := range c.Characters {
		if existing.ID == char.ID {
			return nil
		}
	}

	c.Characters = append(c.Characters, char)
	c.UpdatedAt = time.Now()

	return s.Update(c)
}

func (s *PostgresStore) RemoveCharacter(companyID, characterID string) error {
	c, err := s.Get(companyID)
	if err != nil {
		return err
	}

	for i, ch := range c.Characters {
		if ch.ID == characterID {
			c.Characters = append(c.Characters[:i], c.Characters[i+1:]...)
			c.UpdatedAt = time.Now()
			return s.Update(c)
		}
	}

	return nil
}

func (s *PostgresStore) AddMonster(companyID string, mon monsters.Monster) error {
	c, err := s.Get(companyID)
	if err != nil {
		return err
	}

	for _, existing := range c.Monsters {
		if existing.ID == mon.ID {
			return nil
		}
	}

	c.Monsters = append(c.Monsters, mon)
	c.UpdatedAt = time.Now()

	return s.Update(c)
}

func (s *PostgresStore) RemoveMonster(companyID, monsterID string) error {
	c, err := s.Get(companyID)
	if err != nil {
		return err
	}

	for i, m := range c.Monsters {
		if m.ID == monsterID {
			c.Monsters = append(c.Monsters[:i], c.Monsters[i+1:]...)
			c.UpdatedAt = time.Now()
			return s.Update(c)
		}
	}

	return nil
}


