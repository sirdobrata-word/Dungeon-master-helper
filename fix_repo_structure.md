# Исправление структуры репозитория

## Проблема
Файлы находятся в папке "тишина" вместо корня репозитория.

## Решение

### Вариант 1: Исправить через Git (рекомендуется)

Если ты клонировал репозиторий или работаешь с ним через Git:

```bash
# 1. Клонируй репозиторий (если ещё не клонирован)
git clone https://github.com/sirdobrata-word/D-D-character-list-service.git
cd D-D-character-list-service

# 2. Перемести все файлы из папки "тишина" в корень
# Windows PowerShell:
Get-ChildItem -Path "тишина" -Recurse | Move-Item -Destination "." -Force
Remove-Item "тишина" -Force

# 3. Проверь структуру
ls

# 4. Добавь изменения
git add .
git commit -m "Fix: Move files from 'тишина' folder to root"
git push
```

### Вариант 2: Загрузить заново из локальной папки

Если у тебя есть локальная папка "тишина" с правильной структурой:

```bash
# 1. Перейди в папку "тишина"
cd "C:\Users\ALFA\OneDrive\Рабочий стол\тишина"

# 2. Убедись, что Git инициализирован
git init

# 3. Проверь remote
git remote -v

# Если remote неправильный, исправь:
git remote remove origin
git remote add origin https://github.com/sirdobrata-word/D-D-character-list-service.git

# 4. Добавь все файлы
git add .

# 5. Сделай коммит
git commit -m "Initial commit: D&D Character Sheet Service"

# 6. Загрузи на GitHub (возможно потребуется force push)
git push -u origin main --force
```

⚠️ **Внимание:** `--force` перезапишет историю на GitHub. Используй только если уверен.

### Вариант 3: Оставить как есть

Если структура с папкой "тишина" тебя устраивает, можно оставить так. 
Просто при клонировании нужно будет:
```bash
git clone https://github.com/sirdobrata-word/D-D-character-list-service.git
cd D-D-character-list-service/тишина
go run cmd/server/main.go
```

## Рекомендация

**Лучше использовать Вариант 2** - загрузить файлы напрямую из папки "тишина" в корень репозитория. 
Это стандартная практика для Go проектов.








