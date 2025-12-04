@echo off
chcp 65001 >nul
echo ========================================
echo   Остановка сервера на порту 9190
echo ========================================
echo.

REM Останавливаем контейнер Docker (если запущен)
docker stop dnd-service 2>nul
docker rm dnd-service 2>nul

REM Ищем процесс Go, который может занимать порт
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :9190') do (
    echo [INFO] Найден процесс с PID: %%a
    echo [INFO] Останавливаю процесс...
    taskkill /F /PID %%a >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo [OK] Процесс остановлен
    ) else (
        echo [WARNING] Не удалось остановить процесс (возможно, нет прав)
    )
)

echo.
echo [INFO] Проверка порта...
timeout /t 2 >nul
netstat -ano | findstr :9190 >nul
if %ERRORLEVEL% EQU 0 (
    echo [WARNING] Порт 9190 все еще занят
    echo Попробуй вручную: taskkill /F /PID <номер_из_вывода_выше>
) else (
    echo [OK] Порт 9190 свободен!
)

echo.
pause






