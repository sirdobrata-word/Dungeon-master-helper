@echo off
chcp 65001 >nul
echo ========================================
echo   Исправление зависимостей
echo ========================================
echo.

cd /d "%~dp0"

echo [1] Проверка go.mod...
if not exist "go.mod" (
    echo   [ERROR] go.mod не найден!
    pause
    exit /b 1
)
echo   [OK] go.mod найден

echo.

echo [2] Обновление зависимостей...
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo   [ERROR] Ошибка при go mod tidy
    pause
    exit /b 1
)
echo   [OK] go mod tidy выполнен

echo.

echo [3] Загрузка зависимостей...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo   [ERROR] Ошибка при загрузке зависимостей
    pause
    exit /b 1
)
echo   [OK] Зависимости загружены

echo.

echo [4] Проверка зависимостей...
go mod verify
if %ERRORLEVEL% NEQ 0 (
    echo   [WARNING] Проблемы с зависимостями
) else (
    echo   [OK] Зависимости корректны
)

echo.

echo [5] Проверка компиляции...
go build ./cmd/server
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Компиляция успешна!
    if exist "server.exe" (
        del server.exe >nul 2>&1
    )
) else (
    echo   [ERROR] Ошибка компиляции
    echo   Проверь вывод выше
    pause
    exit /b 1
)

echo.
echo ========================================
echo   Все исправлено!
echo ========================================
echo.
echo Теперь можно запустить сервер:
echo   go run cmd/server/main.go
echo.
pause







