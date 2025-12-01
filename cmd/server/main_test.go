package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dice-service/internal/characters"
)

func TestHandleRollSuccess(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	body := `{"expression":"1d4 + 1"}`
	req := httptest.NewRequest(http.MethodPost, "/roll", strings.NewReader(body))
	rec := httptest.NewRecorder()

	srv.handleRoll(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var payload rollResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Expression != "1d4 + 1" {
		t.Fatalf("expression mismatch, got %q", payload.Expression)
	}
	if payload.Total < 2 || payload.Total > 5 {
		t.Fatalf("total out of range: %d", payload.Total)
	}
}

func TestHandleRollInvalidMethod(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/roll", nil)
	rec := httptest.NewRecorder()

	srv.handleRoll(rec, req)

	if rec.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Result().StatusCode)
	}
}

func TestHandleRollInvalidJSON(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/roll", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	srv.handleRoll(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	assertErrorBody(t, rec.Body, "invalid JSON payload")
}

func TestHandleRollEmptyExpression(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/roll", strings.NewReader(`{"expression":"   "}`))
	rec := httptest.NewRecorder()

	srv.handleRoll(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	assertErrorBody(t, rec.Body, "expression must not be empty")
}

func TestHandleRollInvalidExpression(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/roll", strings.NewReader(`{"expression":"1dx"}`))
	rec := httptest.NewRecorder()

	srv.handleRoll(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleHealth(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "ok" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestCreateCharacter(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	payload := `{"name":"Aria","class":"Rogue","race":"Human","background":"Criminal","level":3,"abilityScores":{"strength":10,"dexterity":16,"constitution":12,"intelligence":13,"wisdom":11,"charisma":14},"proficiencyBonus":2,"armorClass":15,"speed":30,"initiative":3,"maxHitPoints":24,"currentHitPoints":20,"temporaryHitPoints":0}`
	req := httptest.NewRequest(http.MethodPost, "/characters", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	srv.handleCharactersCollection(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var sheet characters.CharacterSheet
	if err := json.NewDecoder(rec.Body).Decode(&sheet); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if sheet.ID == "" {
		t.Fatalf("expected generated id")
	}
}

func TestGetCharacterNotFound(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/characters/missing", nil)
	rec := httptest.NewRecorder()

	srv.handleCharacterByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	assertErrorBody(t, rec.Body, "character not found")
}

func TestGenerateCharacter(t *testing.T) {
	t.Parallel()

	srv := newTestServer()

	body := `{"name":"Таша","class":"Wizard","race":"Human","background":"Sage","level":5}`
	req := httptest.NewRequest(http.MethodPost, "/characters/generate", strings.NewReader(body))
	rec := httptest.NewRecorder()

	srv.handleGenerateCharacter(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var sheet characters.CharacterSheet
	if err := json.NewDecoder(rec.Body).Decode(&sheet); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if sheet.ID == "" {
		t.Fatalf("expected generated id")
	}
	if sheet.Name == "" || sheet.Class == "" {
		t.Fatalf("expected basic fields to be set")
	}
	if sheet.AbilityScores.Strength < 3 || sheet.AbilityScores.Strength > 18 {
		t.Fatalf("unexpected strength score: %d", sheet.AbilityScores.Strength)
	}
}

func newTestServer() *server {
	return newServer(characters.NewMemoryStore())
}

func assertErrorBody(t *testing.T, r io.Reader, want string) {
	t.Helper()
	var payload map[string]string
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		t.Fatalf("failed to decode error: %v", err)
	}
	if payload["error"] != want {
		t.Fatalf("unexpected error: %q, want %q", payload["error"], want)
	}
}
