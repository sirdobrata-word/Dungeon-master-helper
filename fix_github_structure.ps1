# Скрипт для исправления структуры репозитория на GitHub
# Выполни этот скрипт в PowerShell в папке проекта

Write-Host "=== Исправление структуры репозитория на GitHub ===" -ForegroundColor Cyan

# Проверка, что мы в правильной папке
if (-not (Test-Path "go.mod")) {
    Write-Host "ОШИБКА: go.mod не найден. Убедись, что ты в папке проекта!" -ForegroundColor Red
    exit 1
}

Write-Host "`n1. Проверка Git инициализации..." -ForegroundColor Yellow
if (-not (Test-Path ".git")) {
    Write-Host "Инициализация Git..." -ForegroundColor Yellow
    git init
}

Write-Host "`n2. Проверка remote репозитория..." -ForegroundColor Yellow
$remote = git remote get-url origin 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Remote не настроен. Добавляю..." -ForegroundColor Yellow
    git remote add origin https://github.com/sirdobrata-word/D-D-character-list-service.git
} else {
    Write-Host "Remote: $remote" -ForegroundColor Green
    Write-Host "Если remote неправильный, выполни:" -ForegroundColor Yellow
    Write-Host "  git remote remove origin" -ForegroundColor Gray
    Write-Host "  git remote add origin https://github.com/sirdobrata-word/D-D-character-list-service.git" -ForegroundColor Gray
}

Write-Host "`n3. Добавление всех файлов..." -ForegroundColor Yellow
git add .

Write-Host "`n4. Проверка статуса..." -ForegroundColor Yellow
git status

Write-Host "`n5. Создание коммита..." -ForegroundColor Yellow
git commit -m "Fix: Correct repository structure - move files to proper hierarchy"

Write-Host "`n6. Переименование ветки в main (если нужно)..." -ForegroundColor Yellow
git branch -M main

Write-Host "`n=== ВАЖНО ===" -ForegroundColor Red
Write-Host "Следующая команда перезапишет историю на GitHub!" -ForegroundColor Red
Write-Host "Это нормально, если это первый/единственный коммит." -ForegroundColor Yellow
Write-Host "`nВыполни вручную:" -ForegroundColor Cyan
Write-Host "  git push -u origin main --force" -ForegroundColor White

Write-Host "`n=== Готово! ===" -ForegroundColor Green
Write-Host "После выполнения git push структура на GitHub будет исправлена." -ForegroundColor Green








