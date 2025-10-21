
### 2. Создаем `build.bat` для Windows:

```batch
@echo off
chcp 65001 >nul
title 🏗️ SysPulse Builder

echo.
echo ========================================
echo    🏗️  Сборка SysPulse Monitor
echo ========================================
echo.

echo 📦 Очистка предыдущих сборок...
if exist dist rmdir /s /q dist

echo 📁 Создание структуры папок...
mkdir dist
mkdir dist\web
mkdir dist\web\static
mkdir dist\web\static\css
mkdir dist\web\static\js

echo 🛠️  Компиляция приложения...
go build -ldflags="-s -w -X main.version=1.0.0" -o dist\syspulse.exe ./cmd/syspulse-server

if not exist dist\syspulse.exe (
    echo ❌ Ошибка компиляции!
    pause
    exit 1
)

echo 📄 Копирование файлов...
copy web\static\index.html dist\web\static\ >nul
copy web\static\css\style.css dist\web\static\css\ >nul
copy web\static\js\app.js dist\web\static\js\ >nul

echo 🚀 Создание скрипта запуска...
echo @echo off > dist\run.bat
echo chcp 65001 ^>nul >> dist\run.bat
echo title SysPulse Monitor >> dist\run.bat
echo. >> dist\run.bat
echo echo 🚀 Запуск SysPulse Monitor... >> dist\run.bat
echo echo 📊 Системный мониторинг в реальном времени >> dist\run.bat
echo echo. >> dist\run.bat
echo start /B syspulse.exe >> dist\run.bat
echo timeout /t 2 /nobreak ^>nul >> dist\run.bat
echo start http://localhost:8080 >> dist\run.bat
echo echo. >> dist\run.bat
echo echo ✅ Сервер запущен на http://localhost:8080 >> dist\run.bat
echo echo ⚠️  Для остановки закройте это окно >> dist\run.bat
echo echo. >> dist\run.bat
echo pause >> dist\run.bat

echo 📋 Создание README...
echo # SysPulse Monitor > dist\README.txt
echo. >> dist\README.txt
echo 🚀 Запуск: >> dist\README.txt
echo   1. Запустите run.bat >> dist\README.txt
echo   2. Браузер откроется автоматически >> dist\README.txt
echo. >> dist\README.txt
echo 📊 Функции: >> dist\README.txt
echo   - Мониторинг CPU, памяти, диска в реальном времени >> dist\README.txt
echo   - Графики и оповещения >> dist\README.txt
echo   - Обновление каждые 500ms >> dist\README.txt
echo. >> dist\README.txt
echo ⚙️ Настройки порта: >> dist\README.txt
echo   set SYS_PULSE_PORT=9090 >> dist\README.txt
echo   run.bat >> dist\README.txt

echo.
echo 🎉 Сборка завершена!
echo.
echo 📁 Папка dist содержит:
dir dist
echo.
echo 🚀 Для запуска: 
echo   1. Перейдите в папку dist
echo   2. Запустите run.bat
echo.
echo 📦 Для переноса на другой компьютер:
echo   Скопируйте всю папку dist
echo.

pause
