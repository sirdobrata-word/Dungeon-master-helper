# Тестовый скрипт для проверки сервера

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Тест работоспособности сервера" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Проверка 1: Сервер запущен?
Write-Host "[1] Проверка доступности сервера..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:9190/healthz" -Method GET -TimeoutSec 2 -ErrorAction Stop
    if ($response.Content -eq "ok") {
        Write-Host "  ✅ Сервер работает!" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️ Сервер отвечает, но неожиданный ответ: $($response.Content)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ❌ Сервер НЕ запущен или недоступен!" -ForegroundColor Red
    Write-Host "  Запусти сервер командой: go run cmd/server/main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# Проверка 2: Статические файлы доступны?
Write-Host "[2] Проверка веб-интерфейса..." -ForegroundColor Yellow
$pages = @(
    "http://localhost:9190/ui/index.html",
    "http://localhost:9190/ui/list.html",
    "http://localhost:9190/ui/monsters.html",
    "http://localhost:9190/ui/companies.html"
)

foreach ($page in $pages) {
    try {
        $response = Invoke-WebRequest -Uri $page -Method GET -TimeoutSec 2 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            Write-Host "  ✅ $page" -ForegroundColor Green
        }
    } catch {
        Write-Host "  ❌ $page - недоступна!" -ForegroundColor Red
    }
}

Write-Host ""

# Проверка 3: API endpoints
Write-Host "[3] Проверка API endpoints..." -ForegroundColor Yellow

# Health check
try {
    $response = Invoke-WebRequest -Uri "http://localhost:9190/healthz" -Method GET -TimeoutSec 2
    Write-Host "  ✅ GET /healthz" -ForegroundColor Green
} catch {
    Write-Host "  ❌ GET /healthz" -ForegroundColor Red
}

# Characters list
try {
    $response = Invoke-WebRequest -Uri "http://localhost:9190/characters" -Method GET -TimeoutSec 2
    Write-Host "  ✅ GET /characters" -ForegroundColor Green
} catch {
    Write-Host "  ❌ GET /characters" -ForegroundColor Red
}

# Monsters list
try {
    $response = Invoke-WebRequest -Uri "http://localhost:9190/monsters" -Method GET -TimeoutSec 2
    Write-Host "  ✅ GET /monsters" -ForegroundColor Green
} catch {
    Write-Host "  ❌ GET /monsters" -ForegroundColor Red
}

# Companies list
try {
    $response = Invoke-WebRequest -Uri "http://localhost:9190/companies" -Method GET -TimeoutSec 2
    Write-Host "  ✅ GET /companies" -ForegroundColor Green
} catch {
    Write-Host "  ❌ GET /companies" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Открой в браузере:" -ForegroundColor Cyan
Write-Host "  http://localhost:9190/ui/index.html" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

