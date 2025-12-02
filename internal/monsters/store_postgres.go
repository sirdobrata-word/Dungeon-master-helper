package monsters

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// PostgresStore реализует Store для монстров в PostgreSQL.
// Схема: таблица monsters(id text primary key, data jsonb not null)
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore создаёт хранилище монстров в PostgreSQL и гарантирует,
// что таблица существует.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	const createTable = `
CREATE TABLE IF NOT EXISTS monsters (
	id   TEXT PRIMARY KEY,
	data JSONB NOT NULL
);`

	if _, err := db.Exec(createTable); err != nil {
		panic(fmt.Errorf("failed to create monsters table: %w", err))
	}

	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(monster Monster) (Monster, error) {
	if err := monster.Validate(); err != nil {
		return Monster{}, err
	}

	if monster.ID == "" {
		monster.ID = generateID()
	}

	data, err := json.Marshal(monster)
	if err != nil {
		return Monster{}, fmt.Errorf("failed to marshal monster: %w", err)
	}

	const insertQuery = `INSERT INTO monsters (id, data) VALUES ($1, $2::jsonb);`
	if _, err := s.db.Exec(insertQuery, monster.ID, data); err != nil {
		return Monster{}, fmt.Errorf("failed to insert monster: %w", err)
	}

	return monster, nil
}

func (s *PostgresStore) Get(id string) (Monster, error) {
	const selectQuery = `SELECT data FROM monsters WHERE id = $1;`

	var raw []byte
	err := s.db.QueryRow(selectQuery, id).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Monster{}, ErrNotFound
		}
		return Monster{}, fmt.Errorf("failed to get monster: %w", err)
	}

	var monster Monster
	if err := json.Unmarshal(raw, &monster); err != nil {
		return Monster{}, fmt.Errorf("failed to unmarshal monster: %w", err)
	}

	return monster, nil
}

func (s *PostgresStore) List() []Monster {
	const listQuery = `SELECT data FROM monsters ORDER BY id;`

	rows, err := s.db.Query(listQuery)
	if err != nil {
		return []Monster{}
	}
	defer rows.Close()

	var result []Monster
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var monster Monster
		if err := json.Unmarshal(raw, &monster); err != nil {
			continue
		}
		result = append(result, monster)
	}
	return result
}

func (s *PostgresStore) Update(id string, monster Monster) (Monster, error) {
	if err := monster.Validate(); err != nil {
		return Monster{}, err
	}

	monster.ID = id

	data, err := json.Marshal(monster)
	if err != nil {
		return Monster{}, fmt.Errorf("failed to marshal monster: %w", err)
	}

	const updateQuery = `UPDATE monsters SET data = $2::jsonb WHERE id = $1;`
	res, err := s.db.Exec(updateQuery, id, data)
	if err != nil {
		return Monster{}, fmt.Errorf("failed to update monster: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return Monster{}, ErrNotFound
	}

	return monster, nil
}

func (s *PostgresStore) Delete(id string) error {
	const deleteQuery = `DELETE FROM monsters WHERE id = $1;`

	res, err := s.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete monster: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return ErrNotFound
	}

	return nil
}





