# Скрипт проверки установки Go

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Проверка установки Go" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Проверка 1: Команда go доступна?
Write-Host "[1] Проверка наличия команды 'go'..." -ForegroundColor Yellow
$goCommand = Get-Command go -ErrorAction SilentlyContinue

if ($goCommand) {
    Write-Host "  ✅ Go найден!" -ForegroundColor Green
    Write-Host "  Путь: $($goCommand.Source)" -ForegroundColor Gray
} else {
    Write-Host "  ❌ Go НЕ найден в PATH!" -ForegroundColor Red
    Write-Host ""
    Write-Host "  Решение:" -ForegroundColor Yellow
    Write-Host "  1. Скачай Go с https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host "  2. Установи Go" -ForegroundColor Yellow
    Write-Host "  3. Перезапусти терминал" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# Проверка 2: Версия Go
Write-Host "[2] Проверка версии Go..." -ForegroundColor Yellow
try {
    $version = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ $version" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ Не удалось получить версию" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ❌ Ошибка при проверке версии: $_" -ForegroundColor Red
}

Write-Host ""

# Проверка 3: GOROOT
Write-Host "[3] Проверка GOROOT..." -ForegroundColor Yellow
try {
    $goroot = go env GOROOT 2>&1
    if ($LASTEXITCODE -eq 0 -and $goroot) {
        Write-Host "  ✅ GOROOT: $goroot" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ GOROOT не установлен" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ⚠️ Не удалось получить GOROOT" -ForegroundColor Yellow
}

Write-Host ""

# Проверка 4: GOPATH
Write-Host "[4] Проверка GOPATH..." -ForegroundColor Yellow
try {
    $gopath = go env GOPATH 2>&1
    if ($LASTEXITCODE -eq 0 -and $gopath) {
        Write-Host "  ✅ GOPATH: $gopath" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ GOPATH не установлен (это нормально для Go 1.11+)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ⚠️ Не удалось получить GOPATH" -ForegroundColor Yellow
}

Write-Host ""

# Проверка 5: Go модули
Write-Host "[5] Проверка поддержки модулей..." -ForegroundColor Yellow
try {
    $modules = go env GOMOD 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Модули поддерживаются" -ForegroundColor Green
    }
} catch {
    Write-Host "  ⚠️ Не удалось проверить модули" -ForegroundColor Yellow
}

Write-Host ""

# Проверка 6: Тест компиляции
Write-Host "[6] Тест компиляции простой программы..." -ForegroundColor Yellow
$testFile = "$env:TEMP\go_test_$(Get-Random).go"
$testContent = @"
package main
import "fmt"
func main() {
    fmt.Println("Go работает!")
}
"@

try {
    Set-Content -Path $testFile -Value $testContent -Encoding UTF8
    $output = go run $testFile 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Компиляция работает: $output" -ForegroundColor Green
    } else {
        Write-Host "  ❌ Ошибка компиляции: $output" -ForegroundColor Red
    }
} catch {
    Write-Host "  ❌ Ошибка: $_" -ForegroundColor Red
} finally {
    if (Test-Path $testFile) {
        Remove-Item $testFile -Force -ErrorAction SilentlyContinue
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Итоговая проверка" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($goCommand) {
    Write-Host "✅ Go установлен и готов к использованию!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Теперь можно запустить сервер:" -ForegroundColor Yellow
    Write-Host "  go run cmd/server/main.go" -ForegroundColor Gray
    Write-Host "  или" -ForegroundColor Gray
    Write-Host "  .\start-server.ps1" -ForegroundColor Gray
} else {
    Write-Host "❌ Go не установлен!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Установи Go:" -ForegroundColor Yellow
    Write-Host "  1. Скачай с https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host "  2. Установи (по умолчанию в C:\Program Files\Go)" -ForegroundColor Yellow
    Write-Host "  3. Перезапусти терминал" -ForegroundColor Yellow
}



