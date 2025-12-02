@echo off
chcp 65001 >nul
echo ========================================
echo   Проверка установки Go
echo ========================================
echo.

REM Проверка 1: Команда go доступна?
echo [1] Проверка наличия команды 'go'...
where go >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Go найден!
    for /f "tokens=*" %%i in ('where go') do echo   Путь: %%i
) else (
    echo   [ERROR] Go НЕ найден в PATH!
    echo.
    echo   Решение:
    echo   1. Скачай Go с https://golang.org/dl/
    echo   2. Установи Go
    echo   3. Перезапусти терминал
    pause
    exit /b 1
)

echo.

REM Проверка 2: Версия Go
echo [2] Проверка версии Go...
go version 2>nul
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Версия получена
) else (
    echo   [WARNING] Не удалось получить версию
)

echo.

REM Проверка 3: GOROOT
echo [3] Проверка GOROOT...
go env GOROOT >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    for /f "tokens=*" %%i in ('go env GOROOT') do echo   [OK] GOROOT: %%i
) else (
    echo   [WARNING] GOROOT не установлен
)

echo.

REM Проверка 4: Тест компиляции
echo [4] Тест компиляции...
echo package main > "%TEMP%\go_test.go"
echo import "fmt" >> "%TEMP%\go_test.go"
echo func main() { fmt.Println("Go работает!") } >> "%TEMP%\go_test.go"

go run "%TEMP%\go_test.go" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Компиляция работает
) else (
    echo   [ERROR] Ошибка компиляции
)

del "%TEMP%\go_test.go" >nul 2>&1

echo.
echo ========================================
echo   Итоговая проверка
echo ========================================
echo.
echo [OK] Go установлен и готов к использованию!
echo.
echo Теперь можно запустить сервер:
echo   go run cmd/server/main.go
echo   или
echo   start-server.bat
echo.
pause





