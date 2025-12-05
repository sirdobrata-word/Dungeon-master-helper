package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"dice-service/internal/characters"
	"dice-service/internal/company"
	"dice-service/internal/dice"
	"dice-service/internal/monsters"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9190"
	}

	var (
		charStore    characters.Store
		monsterStore monsters.Store
		companyStore company.Store
	)

	// Проверяем наличие DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Printf("DATABASE_URL is not set, using in-memory stores (данные не сохраняются между перезапусками)")
		charStore = characters.NewMemoryStore()
		monsterStore = monsters.NewMemoryStore()
		companyStore = company.NewMemoryStore()
	} else {
		// Подключаемся к PostgreSQL
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to open PostgreSQL connection: %v", err)
		}
		defer db.Close()

		// Проверяем подключение
		if err := db.Ping(); err != nil {
			log.Fatalf("failed to ping PostgreSQL: %v", err)
		}

		log.Printf("Connected to PostgreSQL, using persistent stores")
		
		// Создаём Postgres-хранилища (они автоматически создадут таблицы)
		charStore = characters.NewPostgresStore(db)
		monsterStore = monsters.NewPostgresStore(db)
		companyStore = company.NewPostgresStore(db)
	}

	api := newServer(charStore, monsterStore, companyStore)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           loggingMiddleware(api.routes()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Printf("Starting dice service on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type rollRequest struct {
	Expression string `json:"expression"`
}

type rollResponse struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Modifier   int    `json:"modifier"`
	Total      int    `json:"total"`
}

type server struct {
	characterStore characters.Store
	monsterStore   monsters.Store
	companyStore   company.Store
}

func newServer(charStore characters.Store, monStore monsters.Store, compStore company.Store) *server {
	return &server{
		characterStore: charStore,
		monsterStore:   monStore,
		companyStore:   compStore,
	}
}

func (s *server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	// Register /roll first to ensure it's not overridden
	mux.Handle("/roll", http.HandlerFunc(s.handleRoll))
	mux.HandleFunc("/healthz", handleHealth)
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("web"))))
	mux.Handle("/characters/generate", http.HandlerFunc(s.handleGenerateCharacter))
	mux.Handle("/characters", http.HandlerFunc(s.handleCharactersCollection))
	mux.Handle("/characters/", http.HandlerFunc(s.handleCharacterByID))
	mux.Handle("/monsters", http.HandlerFunc(s.handleMonstersCollection))
	mux.Handle("/monsters/", http.HandlerFunc(s.handleMonsterByID))
	mux.Handle("/monsters/load-samples", http.HandlerFunc(s.handleLoadSampleMonsters))
	
	// Company endpoints
	// Используем точное совпадение для /companies
	mux.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что путь точно /companies (без слеша)
		if r.URL.Path != "/companies" {
			writeError(w, http.StatusNotFound, "404 page not found")
			return
		}
		s.handleCompaniesCollection(w, r)
	})
	mux.Handle("/companies/", http.HandlerFunc(s.handleCompanyByID))
	return mux
}

func (s *server) handleRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var payload rollRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	if strings.TrimSpace(payload.Expression) == "" {
		writeError(w, http.StatusBadRequest, "expression must not be empty")
		return
	}

	expr, err := dice.ParseExpression(payload.Expression)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := dice.Roll(expr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := rollResponse{
		Expression: payload.Expression,
		Rolls:      result.Rolls,
		Modifier:   expr.Modifier,
		Total:      result.Total,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *server) handleCharactersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createCharacter(w, r)
	case http.MethodGet:
		s.listCharacters(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) handleCharacterByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/characters/")
	if path == "" {
		writeError(w, http.StatusBadRequest, "missing character id")
		return
	}

	// Проверяем специальный путь для повышения уровня
	if strings.HasSuffix(path, "/levelup") {
		if r.Method == http.MethodPost {
			id := strings.TrimSuffix(path, "/levelup")
			if id == "" {
				writeError(w, http.StatusBadRequest, "missing character id")
				return
			}
			s.levelUpCharacter(w, r, id)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Обычный путь с ID
	id := path
	switch r.Method {
	case http.MethodGet:
		s.getCharacter(w, r, id)
	case http.MethodPut:
		s.updateCharacter(w, r, id)
	case http.MethodDelete:
		s.deleteCharacter(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) createCharacter(w http.ResponseWriter, r *http.Request) {
	var sheet characters.CharacterSheet
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&sheet); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	sheet.ID = ""

	created, err := s.characterStore.Create(sheet)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *server) listCharacters(w http.ResponseWriter, r *http.Request) {
	items := s.characterStore.List()
	writeJSON(w, http.StatusOK, items)
}

func (s *server) getCharacter(w http.ResponseWriter, r *http.Request, id string) {
	sheet, err := s.characterStore.Get(id)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sheet)
}

func (s *server) updateCharacter(w http.ResponseWriter, r *http.Request, id string) {
	var sheet characters.CharacterSheet
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&sheet); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	// путь определяет идентификатор
	sheet.ID = id

	updated, err := s.characterStore.Update(id, sheet)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		// ошибки валидации и т.п. считаем ошибкой запроса
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (s *server) deleteCharacter(w http.ResponseWriter, r *http.Request, id string) {
	err := s.characterStore.Delete(id)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "character deleted"})
}

func (s *server) levelUpCharacter(w http.ResponseWriter, r *http.Request, id string) {
	// Получаем текущего персонажа
	sheet, err := s.characterStore.Get(id)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Повышаем уровень
	leveledUp, err := characters.LevelUp(sheet)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Сохраняем обновлённого персонажа
	updated, err := s.characterStore.Update(id, leveledUp)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

type generateCharacterRequest struct {
	Name       string   `json:"name"`
	Class      string   `json:"class"`
	Race       string   `json:"race"`
	Background string   `json:"background"`
	Alignment  string   `json:"alignment"`
	Level      int      `json:"level"`
	Skills     []string `json:"skills"`
}

func (s *server) handleGenerateCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var payload generateCharacterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	sheet, err := characters.GenerateCharacterSheet(
		payload.Name,
		payload.Class,
		payload.Race,
		payload.Background,
		payload.Alignment,
		payload.Level,
		payload.Skills,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := s.characterStore.Create(sheet)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *server) handleMonstersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createMonster(w, r)
	case http.MethodGet:
		s.listMonsters(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) handleMonsterByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/monsters/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing monster id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getMonster(w, r, id)
	case http.MethodPut:
		s.updateMonster(w, r, id)
	case http.MethodDelete:
		s.deleteMonster(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) createMonster(w http.ResponseWriter, r *http.Request) {
	var monster monsters.Monster
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&monster); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	monster.ID = ""

	created, err := s.monsterStore.Create(monster)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *server) listMonsters(w http.ResponseWriter, r *http.Request) {
	items := s.monsterStore.List()
	writeJSON(w, http.StatusOK, items)
}

func (s *server) getMonster(w http.ResponseWriter, r *http.Request, id string) {
	monster, err := s.monsterStore.Get(id)
	if err != nil {
		if errors.Is(err, monsters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "monster not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, monster)
}

func (s *server) updateMonster(w http.ResponseWriter, r *http.Request, id string) {
	var monster monsters.Monster
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&monster); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	monster.ID = id

	updated, err := s.monsterStore.Update(id, monster)
	if err != nil {
		if errors.Is(err, monsters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "monster not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (s *server) deleteMonster(w http.ResponseWriter, r *http.Request, id string) {
	err := s.monsterStore.Delete(id)
	if err != nil {
		if errors.Is(err, monsters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "monster not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "monster deleted"})
}

func (s *server) handleLoadSampleMonsters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	samples := monsters.GetSampleMonsters()
	loaded := 0

	for _, sample := range samples {
		_, err := s.monsterStore.Create(sample)
		if err == nil {
			loaded++
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "sample monsters loaded",
		"count":   loaded,
	})
}

// ===== Company Handlers =====

func (s *server) handleCompaniesCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleCompaniesCollection: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		s.createCompany(w, r)
	case http.MethodGet:
		s.listCompanies(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) handleCompanyByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleCompanyByID: %s %s", r.Method, r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/companies/")
	parts := strings.Split(path, "/")
	id := parts[0]
	
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing company id")
		return
	}

	// Проверяем дополнительные пути для персонажей и монстров
	if len(parts) >= 2 {
		switch parts[1] {
		case "characters":
			if len(parts) == 2 {
				s.handleCompanyCharacters(w, r, id)
			} else {
				s.handleCompanyCharacterByID(w, r, id, parts[2])
			}
			return
		case "monsters":
			if len(parts) == 2 {
				s.handleCompanyMonsters(w, r, id)
			} else {
				s.handleCompanyMonsterByID(w, r, id, parts[2])
			}
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		s.getCompany(w, r, id)
	case http.MethodPut:
		s.updateCompany(w, r, id)
	case http.MethodDelete:
		s.deleteCompany(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) createCompany(w http.ResponseWriter, r *http.Request) {
	var comp company.Company
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload: "+err.Error())
		return
	}
	comp.ID = ""
	// Инициализируем пустые слайсы если они nil
	if comp.Characters == nil {
		comp.Characters = []characters.CharacterSheet{}
	}
	if comp.Monsters == nil {
		comp.Monsters = []monsters.Monster{}
	}

	created, err := s.companyStore.Create(comp)
	if err != nil {
		if errors.Is(err, company.ErrDuplicateName) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *server) listCompanies(w http.ResponseWriter, r *http.Request) {
	items := s.companyStore.List()
	log.Printf("listCompanies: возвращено %d компаний", len(items))
	writeJSON(w, http.StatusOK, items)
}

func (s *server) getCompany(w http.ResponseWriter, r *http.Request, id string) {
	comp, err := s.companyStore.Get(id)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comp)
}

func (s *server) updateCompany(w http.ResponseWriter, r *http.Request, id string) {
	var comp company.Company
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	comp.ID = id

	err := s.companyStore.Update(comp)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		if errors.Is(err, company.ErrDuplicateName) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, _ := s.companyStore.Get(id)
	writeJSON(w, http.StatusOK, updated)
}

func (s *server) deleteCompany(w http.ResponseWriter, r *http.Request, id string) {
	err := s.companyStore.Delete(id)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "company deleted"})
}

// Добавление персонажа в компанию
func (s *server) handleCompanyCharacters(w http.ResponseWriter, r *http.Request, companyID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req company.AddCharacterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	// Получаем персонажа из основного хранилища
	char, err := s.characterStore.Get(req.CharacterID)
	if err != nil {
		if errors.Is(err, characters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "character not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавляем в компанию
	err = s.companyStore.AddCharacter(companyID, char)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "character added to company"})
}

// Удаление персонажа из компании
func (s *server) handleCompanyCharacterByID(w http.ResponseWriter, r *http.Request, companyID, characterID string) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	err := s.companyStore.RemoveCharacter(companyID, characterID)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "character removed from company"})
}

// Добавление монстра в компанию
func (s *server) handleCompanyMonsters(w http.ResponseWriter, r *http.Request, companyID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req company.AddMonsterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	// Получаем монстра из основного хранилища
	mon, err := s.monsterStore.Get(req.MonsterID)
	if err != nil {
		if errors.Is(err, monsters.ErrNotFound) {
			writeError(w, http.StatusNotFound, "monster not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Добавляем в компанию
	err = s.companyStore.AddMonster(companyID, mon)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "monster added to company"})
}

// Удаление монстра из компании
func (s *server) handleCompanyMonsterByID(w http.ResponseWriter, r *http.Request, companyID, monsterID string) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	err := s.companyStore.RemoveMonster(companyID, monsterID)
	if err != nil {
		if errors.Is(err, company.ErrNotFound) {
			writeError(w, http.StatusNotFound, "company not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "monster removed from company"})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, map[string]string{"error": message})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, duration)
	})
}

