@echo off
chcp 65001 >nul
echo ========================================
echo   Диагностика сервера D&D Service
echo ========================================
echo.

cd /d "%~dp0"

REM Проверка 1: Go
echo [1] Проверка Go...
where go >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    go version
    echo   [OK] Go установлен
) else (
    echo   [ERROR] Go НЕ установлен!
    echo   Установите Go с https://golang.org/dl/
    goto :errors
)

echo.

REM Проверка 2: Файлы
echo [2] Проверка структуры проекта...
if exist "go.mod" (
    echo   [OK] go.mod
) else (
    echo   [ERROR] go.mod не найден!
    goto :errors
)

if exist "cmd\server\main.go" (
    echo   [OK] cmd\server\main.go
) else (
    echo   [ERROR] cmd\server\main.go не найден!
    goto :errors
)

if exist "internal\characters\sheet.go" (
    echo   [OK] internal\characters\sheet.go
) else (
    echo   [ERROR] internal\characters\sheet.go не найден!
    goto :errors
)

if exist "web\index.html" (
    echo   [OK] web\index.html
) else (
    echo   [ERROR] web\index.html не найден!
    goto :errors
)

echo.

REM Проверка 3: Зависимости
echo [3] Проверка зависимостей...
go mod verify >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Зависимости корректны
) else (
    echo   [WARNING] Проблемы с зависимостями
    echo   Выполните: go mod download
)

echo.

REM Проверка 4: Компиляция
echo [4] Проверка компиляции...
go build ./cmd/server >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Компиляция успешна
    if exist "server.exe" (
        echo   [OK] server.exe создан
        del server.exe >nul 2>&1
    )
) else (
    echo   [ERROR] Ошибка компиляции!
    echo   Попробуйте: go build ./cmd/server
    goto :errors
)

echo.

REM Проверка 5: Порт
echo [5] Проверка порта 9190...
netstat -ano | findstr ":9190" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [WARNING] Порт 9190 занят
    echo   Закройте процесс или используйте другой порт
) else (
    echo   [OK] Порт 9190 свободен
)

echo.
echo ========================================
echo   Итоговый отчет
echo ========================================
echo.
echo [OK] Все проверки пройдены!
echo.
echo Сервер готов к запуску:
echo   go run cmd/server/main.go
echo   или
echo   start-server.bat
echo.
goto :end

:errors
echo.
echo ========================================
echo   ОШИБКИ ОБНАРУЖЕНЫ!
echo ========================================
echo.
echo Исправьте ошибки и запустите диагностику снова.
echo.
pause
exit /b 1

:end
pause

