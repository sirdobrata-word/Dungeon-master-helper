@echo off
chcp 65001 >nul
echo ========================================
echo   Проверка порта 9190
echo ========================================
echo.

netstat -ano | findstr :9190

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [INFO] Порт 9190 занят!
    echo.
    echo Чтобы освободить порт, выполни:
    echo   taskkill /F /PID <номер_процесса>
    echo.
    echo Или используй другой порт для Docker:
    echo   docker-compose up --build -p 9191:9190
    echo.
) else (
    echo [OK] Порт 9190 свободен
    echo.
)

pause




