# ⚡ Быстрый старт с PostgreSQL

## 🚀 Запуск через Docker (самый простой способ)

### 1. Обнови зависимости Go (один раз)

В `cmd` выполни:

```bash
cd /d "C:\Users\ALFA\OneDrive\Рабочий стол\тишина"
go mod tidy
```

Или запусти скрипт:
```bash
update-deps.bat
```

### 2. Запусти Docker Compose

```bash
docker-compose up --build
```

**Готово!** Сервис запущен с PostgreSQL. Все данные сохраняются в базе данных.

---

## 🖥️ Локальный запуск (без Docker)

### 1. Установи PostgreSQL

Скачай и установи: https://www.postgresql.org/download/windows/

### 2. Создай базу данных

Открой **SQL Shell (psql)** и выполни:

```sql
CREATE DATABASE dnd_service;
CREATE USER dnd_user WITH PASSWORD 'dnd_password';
GRANT ALL PRIVILEGES ON DATABASE dnd_service TO dnd_user;
\q
```

### 3. Настрой переменную окружения

В PowerShell:

```powershell
setx DATABASE_URL "postgres://dnd_user:dnd_password@localhost:5432/dnd_service?sslmode=disable"
```

**Важно:** Закрой и открой заново PowerShell/CMD!

### 4. Обнови зависимости и запусти

```bash
go mod tidy
go run cmd\server\main.go
```

---

## ✅ Проверка

Открой в браузере:
- http://localhost:9190/ui/index.html

Создай персонажа, монстра или компанию — они сохранятся в PostgreSQL!

---

## 📚 Подробные инструкции

- **Docker:** `DOCKER_INSTRUCTIONS.md`
- **PostgreSQL:** `POSTGRES_SETUP.md`






