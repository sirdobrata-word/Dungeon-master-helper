# D&D Character Sheet Service

Сервис для создания и управления листами персонажей D&D, а также для броска кубиков.

## Запуск сервиса

### Вариант 1: Через скрипт запуска (рекомендуется)
```bash
# Windows CMD
start-server.bat

# Windows PowerShell
.\start-server.ps1
```

### Вариант 2: Через go run
```bash
go run cmd/server/main.go
```

### Вариант 3: Сборка и запуск
```bash
go build -o dice-service.exe cmd/server/main.go
./dice-service.exe
```

### Вариант 4: С указанием порта
```bash
# Windows PowerShell
$env:PORT="3000"; go run cmd/server/main.go

# Windows CMD
set PORT=3000 && go run cmd/server/main.go
```

Сервис запустится на порту **9190** по умолчанию (или на указанном в переменной окружения `PORT`).

## API Эндпоинты

### 1. Проверка здоровья
- **GET** `http://localhost:9190/healthz`
- Ответ: `ok`

### 2. Бросок кубиков
- **POST** `http://localhost:9190/roll`
- Тело запроса:
```json
{
  "expression": "2d6 + 3"
}
```
- Ответ:
```json
{
  "expression": "2d6 + 3",
  "rolls": [4, 5],
  "modifier": 3,
  "total": 12
}
```

### 3. Создание персонажа
- **POST** `http://localhost:9190/characters`
- Тело запроса:
```json
{
  "name": "Арагорн",
  "class": "Ranger",
  "race": "Human",
  "background": "Folk Hero",
  "level": 5,
  "abilityScores": {
    "strength": 16,
    "dexterity": 14,
    "constitution": 15,
    "intelligence": 12,
    "wisdom": 13,
    "charisma": 10
  },
  "proficiencyBonus": 3,
  "armorClass": 16,
  "speed": 30,
  "initiative": 2,
  "maxHitPoints": 45,
  "currentHitPoints": 45,
  "temporaryHitPoints": 0
}
```
- Ответ: созданный персонаж с автоматически сгенерированным `id`

### 4. Получение списка всех персонажей
- **GET** `http://localhost:9190/characters`
- Ответ: массив всех персонажей

### 5. Получение персонажа по ID
- **GET** `http://localhost:9190/characters/{id}`
- Пример: `http://localhost:9190/characters/a1b2c3d4e5f6g7h8`
- Ответ: объект персонажа или ошибка 404

## Использование в Postman

### Настройка коллекции

1. **Создайте новую коллекцию** в Postman (например, "D&D Service")

2. **Создайте переменную окружения:**
   - Base URL: `http://localhost:9190`
   - В Postman: Environments → Create Environment → добавить переменную `base_url` = `http://localhost:9190`

### Примеры запросов

#### 1. Health Check
- Method: `GET`
- URL: `{{base_url}}/healthz`
- Headers: не требуются

#### 2. Roll Dice
- Method: `POST`
- URL: `{{base_url}}/roll`
- Headers:
  - `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "expression": "1d20 + 5"
}
```

#### 3. Create Character
- Method: `POST`
- URL: `{{base_url}}/characters`
- Headers:
  - `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "name": "Леголас",
  "class": "Fighter",
  "race": "Elf",
  "background": "Noble",
  "level": 3,
  "abilityScores": {
    "strength": 13,
    "dexterity": 18,
    "constitution": 14,
    "intelligence": 11,
    "wisdom": 15,
    "charisma": 12
  },
  "proficiencyBonus": 2,
  "armorClass": 17,
  "speed": 30,
  "initiative": 4,
  "maxHitPoints": 28,
  "currentHitPoints": 28,
  "temporaryHitPoints": 0
}
```

#### 4. List All Characters
- Method: `GET`
- URL: `{{base_url}}/characters`
- Headers: не требуются

#### 5. Get Character by ID
- Method: `GET`
- URL: `{{base_url}}/characters/{{character_id}}`
- Headers: не требуются
- Примечание: `character_id` нужно получить из ответа при создании персонажа

### Тестирование в Postman

1. Запустите сервис: `go run cmd/server/main.go`
2. Откройте Postman
3. Сначала проверьте health: `GET http://localhost:9190/healthz`
4. Создайте персонажа: `POST http://localhost:9190/characters`
5. Скопируйте `id` из ответа
6. Получите персонажа: `GET http://localhost:9190/characters/{id}`
7. Получите список: `GET http://localhost:9190/characters`

## Примеры выражений для кубиков

- `1d20` - один 20-гранный кубик
- `2d6` - два 6-гранных кубика
- `1d20 + 5` - один 20-гранный кубик + модификатор 5
- `3d8 - 2` - три 8-гранных кубика - модификатор 2
- `4d6` - четыре 6-гранных кубика (стандарт для генерации характеристик)

## Тестирование

Запуск всех тестов:
```bash
go test ./...
```

Запуск тестов с покрытием:
```bash
go test -cover ./...
```

## Развёртывание на GitHub

### Первоначальная настройка

1. **Установи Git** (если ещё не установлен):
   - Скачай с [git-scm.com](https://git-scm.com/downloads)
   - Или через пакетный менеджер: `winget install Git.Git`

2. **Настрой Git** (если первый раз):
   ```bash
   git config --global user.name "Твоё Имя"
   git config --global user.email "твой@email.com"
   ```

### Создание репозитория на GitHub

1. Зайди на [github.com](https://github.com) и войди в аккаунт
2. Нажми кнопку **"New"** (или **"+"** → **"New repository"**)
3. Заполни:
   - **Repository name**: `dnd-character-service` (или любое другое имя)
   - **Description**: "D&D Character Sheet Service на Go"
   - **Visibility**: Public или Private (на твой выбор)
   - **НЕ** ставь галочки на "Initialize with README", "Add .gitignore", "Choose a license" (у нас уже есть файлы)
4. Нажми **"Create repository"**

### Загрузка кода на GitHub

Выполни эти команды в терминале в папке проекта:

```bash
# 1. Инициализируй Git репозиторий (если ещё не инициализирован)
git init

# 2. Добавь все файлы в staging
git add .

# 3. Сделай первый коммит
git commit -m "Initial commit: D&D Character Sheet Service"

# 4. Добавь удалённый репозиторий (замени USERNAME и REPO_NAME на свои)
git remote add origin https://github.com/USERNAME/REPO_NAME.git

# 5. Переименуй ветку в main (если нужно)
git branch -M main

# 6. Загрузи код на GitHub
git push -u origin main
```

**Важно:** Замени `USERNAME` на свой GitHub username и `REPO_NAME` на имя репозитория, которое ты создал.

### Обновление кода на GitHub

Когда вносишь изменения, используй:

```bash
# 1. Проверь статус изменений
git status

# 2. Добавь изменённые файлы
git add .

# 3. Сделай коммит с описанием изменений
git commit -m "Описание того, что изменилось"

# 4. Загрузи изменения на GitHub
git push
```

### Что будет загружено

- ✅ Весь исходный код Go
- ✅ HTML/CSS/JS файлы веб-интерфейса
- ✅ README.md с документацией
- ✅ postman_collection.json с тестами
- ✅ go.mod с зависимостями
- ✅ .gitignore (исключает бинарные файлы и временные файлы)

### Что НЕ будет загружено (благодаря .gitignore)

- ❌ Скомпилированные `.exe` файлы
- ❌ Временные файлы
- ❌ IDE настройки
- ❌ Логи

### Клонирование репозитория

Если кто-то захочет склонировать твой проект:

```bash
git clone https://github.com/USERNAME/REPO_NAME.git
cd REPO_NAME
go run cmd/server/main.go
```

