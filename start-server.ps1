# D&D Character Sheet Service - PowerShell Startup Script
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  D&D Character Sheet Service" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Переход в директорию скрипта
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

# Проверка наличия Go
$goPath = Get-Command go -ErrorAction SilentlyContinue
if (-not $goPath) {
    Write-Host "[ERROR] Go не найден в PATH!" -ForegroundColor Red
    Write-Host "Установите Go с https://golang.org/dl/" -ForegroundColor Yellow
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

# Проверка наличия файла
$mainFile = Join-Path $scriptDir "cmd\server\main.go"
if (-not (Test-Path $mainFile)) {
    Write-Host "[ERROR] Файл cmd\server\main.go не найден!" -ForegroundColor Red
    Write-Host "Текущая директория: $scriptDir" -ForegroundColor Yellow
    Write-Host "Ожидаемый путь: $mainFile" -ForegroundColor Yellow
    Write-Host "Убедитесь, что скрипт запущен из корня проекта" -ForegroundColor Yellow
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

Write-Host "[INFO] Go найден: $($goPath.Source)" -ForegroundColor Green
Write-Host "[INFO] Текущая директория: $scriptDir" -ForegroundColor Green
Write-Host "[INFO] Запуск сервера..." -ForegroundColor Green
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Сервер будет доступен по адресу:" -ForegroundColor Cyan
Write-Host "  http://localhost:9190" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Основные эндпоинты:" -ForegroundColor White
Write-Host "  - GET  http://localhost:9190/healthz" -ForegroundColor Gray
Write-Host "  - POST http://localhost:9190/roll" -ForegroundColor Gray
Write-Host "  - GET  http://localhost:9190/characters" -ForegroundColor Gray
Write-Host "  - POST http://localhost:9190/characters" -ForegroundColor Gray
Write-Host "  - GET  http://localhost:9190/monsters" -ForegroundColor Gray
Write-Host "  - GET  http://localhost:9190/companies" -ForegroundColor Gray
Write-Host ""

Write-Host "Веб-интерфейс:" -ForegroundColor White
Write-Host "  - http://localhost:9190/ui/index.html" -ForegroundColor Gray
Write-Host "  - http://localhost:9190/ui/list.html" -ForegroundColor Gray
Write-Host "  - http://localhost:9190/ui/monsters.html" -ForegroundColor Gray
Write-Host "  - http://localhost:9190/ui/companies.html" -ForegroundColor Gray
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Нажмите Ctrl+C для остановки сервера" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Запуск сервера
try {
    go run cmd/server/main.go
    if ($LASTEXITCODE -ne 0) {
        throw "Go вернул код ошибки: $LASTEXITCODE"
    }
} catch {
    Write-Host ""
    Write-Host "[ERROR] Ошибка запуска сервера!" -ForegroundColor Red
    Write-Host "Ошибка: $_" -ForegroundColor Red
    Write-Host "Проверьте, что:" -ForegroundColor Yellow
    Write-Host "  1. Go установлен правильно" -ForegroundColor Yellow
    Write-Host "  2. Все зависимости установлены (go mod download)" -ForegroundColor Yellow
    Write-Host "  3. Файл cmd/server/main.go существует" -ForegroundColor Yellow
    Write-Host "  4. Текущая директория: $scriptDir" -ForegroundColor Yellow
    Read-Host "Нажмите Enter для выхода"
    exit 1
}

