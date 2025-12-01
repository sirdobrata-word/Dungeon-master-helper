# Быстрая инструкция по загрузке на GitHub

## Шаг 1: Создай репозиторий на GitHub

1. Зайди на https://github.com
2. Нажми **"New repository"**
3. Назови репозиторий (например: `dnd-character-service`)
4. **НЕ** ставь галочки на "Initialize with README"
5. Нажми **"Create repository"**

## Шаг 2: Выполни команды в терминале

Открой PowerShell или CMD в папке проекта и выполни:

```bash
# Инициализация Git
git init

# Добавление всех файлов
git add .

# Первый коммит
git commit -m "Initial commit: D&D Character Sheet Service"

# Добавление удалённого репозитория (ЗАМЕНИ USERNAME и REPO_NAME!)
git remote add origin https://github.com/USERNAME/REPO_NAME.git

# Переименование ветки в main
git branch -M main

# Загрузка на GitHub
git push -u origin main
```

## Важно!

**Замени в команде `git remote add origin`:**
- `USERNAME` → твой GitHub username
- `REPO_NAME` → имя репозитория, которое ты создал

Например, если твой username `john` и репозиторий `dnd-service`, команда будет:
```bash
git remote add origin https://github.com/john/dnd-service.git
```

## После загрузки

Твой код будет доступен по адресу:
```
https://github.com/USERNAME/REPO_NAME
```

## Обновление кода

Когда вносишь изменения:

```bash
git add .
git commit -m "Описание изменений"
git push
```

## Если возникли проблемы

### Ошибка "remote origin already exists"
```bash
git remote remove origin
git remote add origin https://github.com/USERNAME/REPO_NAME.git
```

### Нужно авторизоваться
GitHub может попросить ввести логин и пароль. Используй:
- **Username**: твой GitHub username
- **Password**: Personal Access Token (не обычный пароль!)

Как создать токен:
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token
3. Выбери права: `repo`
4. Скопируй токен и используй его как пароль

