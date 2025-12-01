# Диагностика проблем сервера

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Диагностика сервера D&D Service" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

$errors = @()
$warnings = @()

# Проверка 1: Go установлен
Write-Host "[1] Проверка Go..." -ForegroundColor Yellow
$goCommand = Get-Command go -ErrorAction SilentlyContinue
if ($goCommand) {
    $version = go version 2>&1
    Write-Host "  ✅ Go установлен: $version" -ForegroundColor Green
} else {
    Write-Host "  ❌ Go НЕ установлен!" -ForegroundColor Red
    $errors += "Go не установлен"
}

Write-Host ""

# Проверка 2: Структура проекта
Write-Host "[2] Проверка структуры проекта..." -ForegroundColor Yellow

$requiredFiles = @(
    "go.mod",
    "cmd/server/main.go",
    "internal/characters/sheet.go",
    "internal/characters/store.go",
    "internal/monsters/monster.go",
    "internal/monsters/store.go",
    "internal/company/company.go",
    "internal/company/store.go",
    "internal/dice/roller.go",
    "web/index.html",
    "web/styles.css"
)

foreach ($file in $requiredFiles) {
    $fullPath = Join-Path $scriptDir $file
    if (Test-Path $fullPath) {
        Write-Host "  ✅ $file" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $file - НЕ НАЙДЕН!" -ForegroundColor Red
        $errors += "Файл не найден: $file"
    }
}

Write-Host ""

# Проверка 3: Зависимости
Write-Host "[3] Проверка зависимостей..." -ForegroundColor Yellow
if (Test-Path "go.mod") {
    Write-Host "  ✅ go.mod существует" -ForegroundColor Green
    
    # Проверка go.sum
    if (Test-Path "go.sum") {
        Write-Host "  ✅ go.sum существует" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ go.sum отсутствует (будет создан при go mod download)" -ForegroundColor Yellow
        $warnings += "go.sum отсутствует"
    }
    
    # Попытка проверить зависимости
    Write-Host "  Проверка зависимостей..." -ForegroundColor Gray
    $modCheck = go mod verify 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Зависимости корректны" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ Проблемы с зависимостями" -ForegroundColor Yellow
        $warnings += "Проблемы с зависимостями: $modCheck"
    }
} else {
    Write-Host "  ❌ go.mod не найден!" -ForegroundColor Red
    $errors += "go.mod не найден"
}

Write-Host ""

# Проверка 4: Компиляция
Write-Host "[4] Проверка компиляции..." -ForegroundColor Yellow
if ($goCommand) {
    Write-Host "  Попытка компиляции..." -ForegroundColor Gray
    $buildOutput = go build ./cmd/server 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Компиляция успешна!" -ForegroundColor Green
        if (Test-Path "server.exe") {
            Write-Host "  ✅ server.exe создан" -ForegroundColor Green
        } elseif (Test-Path "server") {
            Write-Host "  ✅ server создан" -ForegroundColor Green
        }
    } else {
        Write-Host "  ❌ Ошибка компиляции!" -ForegroundColor Red
        Write-Host "  Вывод:" -ForegroundColor Red
        Write-Host $buildOutput -ForegroundColor Red
        $errors += "Ошибка компиляции: $buildOutput"
    }
} else {
    Write-Host "  ⚠️ Пропущено (Go не установлен)" -ForegroundColor Yellow
}

Write-Host ""

# Проверка 5: Синтаксис Go файлов
Write-Host "[5] Проверка синтаксиса..." -ForegroundColor Yellow
if ($goCommand) {
    $goFiles = Get-ChildItem -Path . -Recurse -Filter "*.go" -ErrorAction SilentlyContinue
    $syntaxErrors = 0
    foreach ($file in $goFiles) {
        $check = go fmt $file.FullName 2>&1
        if ($LASTEXITCODE -ne 0) {
            $syntaxErrors++
            Write-Host "  ⚠️ Проблемы в $($file.Name)" -ForegroundColor Yellow
        }
    }
    if ($syntaxErrors -eq 0) {
        Write-Host "  ✅ Синтаксис корректен" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ Найдено проблем: $syntaxErrors" -ForegroundColor Yellow
        $warnings += "Проблемы синтаксиса в $syntaxErrors файлах"
    }
}

Write-Host ""

# Проверка 6: Порт 9190
Write-Host "[6] Проверка порта 9190..." -ForegroundColor Yellow
$portCheck = Get-NetTCPConnection -LocalPort 9190 -ErrorAction SilentlyContinue
if ($portCheck) {
    $process = Get-Process -Id $portCheck.OwningProcess -ErrorAction SilentlyContinue
    Write-Host "  ⚠️ Порт 9190 занят процессом: $($process.ProcessName) (PID: $($portCheck.OwningProcess))" -ForegroundColor Yellow
    $warnings += "Порт 9190 занят"
} else {
    Write-Host "  ✅ Порт 9190 свободен" -ForegroundColor Green
}

Write-Host ""

# Проверка 7: Веб-файлы
Write-Host "[7] Проверка веб-файлов..." -ForegroundColor Yellow
$webFiles = @("index.html", "list.html", "monsters.html", "companies.html", "styles.css")
foreach ($file in $webFiles) {
    $webPath = Join-Path $scriptDir "web\$file"
    if (Test-Path $webPath) {
        $size = (Get-Item $webPath).Length
        Write-Host "  ✅ web\$file ($size байт)" -ForegroundColor Green
    } else {
        Write-Host "  ❌ web\$file - НЕ НАЙДЕН!" -ForegroundColor Red
        $errors += "Веб-файл не найден: web\$file"
    }
}

Write-Host ""

# Итоговый отчет
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Итоговый отчет" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($errors.Count -eq 0 -and $warnings.Count -eq 0) {
    Write-Host "✅ Все проверки пройдены успешно!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Сервер готов к запуску:" -ForegroundColor Yellow
    Write-Host "  go run cmd/server/main.go" -ForegroundColor Gray
    Write-Host "  или" -ForegroundColor Gray
    Write-Host "  .\start-server.ps1" -ForegroundColor Gray
} else {
    if ($errors.Count -gt 0) {
        Write-Host "❌ КРИТИЧЕСКИЕ ОШИБКИ ($($errors.Count)):" -ForegroundColor Red
        foreach ($error in $errors) {
            Write-Host "  - $error" -ForegroundColor Red
        }
        Write-Host ""
    }
    
    if ($warnings.Count -gt 0) {
        Write-Host "⚠️ ПРЕДУПРЕЖДЕНИЯ ($($warnings.Count)):" -ForegroundColor Yellow
        foreach ($warning in $warnings) {
            Write-Host "  - $warning" -ForegroundColor Yellow
        }
        Write-Host ""
    }
    
    Write-Host "Рекомендации:" -ForegroundColor Cyan
    if ($errors -like "*Go не установлен*") {
        Write-Host "  1. Установите Go с https://golang.org/dl/" -ForegroundColor Yellow
    }
    if ($errors -like "*go.mod*" -or $warnings -like "*зависимости*") {
        Write-Host "  2. Выполните: go mod download" -ForegroundColor Yellow
    }
    if ($warnings -like "*Порт 9190*") {
        Write-Host "  3. Закройте процесс на порту 9190 или используйте другой порт" -ForegroundColor Yellow
    }
}

Write-Host ""

