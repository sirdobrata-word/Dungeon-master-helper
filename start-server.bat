@echo off
chcp 65001 >nul
echo ========================================
echo   D&D Character Sheet Service
echo ========================================
echo.

REM Переход в директорию скрипта
cd /d "%~dp0"

REM Проверка наличия Go
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go не найден в PATH!
    echo Установите Go с https://golang.org/dl/
    pause
    exit /b 1
)

REM Проверка наличия файла
if not exist "cmd\server\main.go" (
    echo [ERROR] Файл cmd\server\main.go не найден!
    echo Текущая директория: %CD%
    echo Убедитесь, что скрипт запущен из корня проекта
    pause
    exit /b 1
)

echo [INFO] Go найден
echo [INFO] Текущая директория: %CD%
echo [INFO] Запуск сервера...
echo.
echo ========================================
echo   Сервер будет доступен по адресу:
echo   http://localhost:9190
echo ========================================
echo.
echo Основные эндпоинты:
echo   - GET  http://localhost:9190/healthz
echo   - POST http://localhost:9190/roll
echo   - GET  http://localhost:9190/characters
echo   - POST http://localhost:9190/characters
echo   - GET  http://localhost:9190/monsters
echo   - GET  http://localhost:9190/companies
echo.
echo Веб-интерфейс:
echo   - http://localhost:9190/ui/index.html
echo   - http://localhost:9190/ui/list.html
echo   - http://localhost:9190/ui/monsters.html
echo   - http://localhost:9190/ui/companies.html
echo.
echo ========================================
echo Нажмите Ctrl+C для остановки сервера
echo ========================================
echo.

go run cmd\server\main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Ошибка запуска сервера!
    echo Проверьте, что:
    echo   1. Go установлен правильно
    echo   2. Все зависимости установлены (go mod download)
    echo   3. Файл cmd/server/main.go существует
    pause
    exit /b 1
)








