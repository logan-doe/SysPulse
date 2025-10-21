#!/bin/bash

echo "🏗️  Сборка SysPulse Monitor..."

# Очищаем предыдущую сборку
rm -rf dist
mkdir -p dist/web/static/{css,js}

echo "📦 Компиляция бинарника..."
go build -ldflags="-s -w -X main.version=1.0.0" -o dist/syspulse ./cmd/syspulse-server

if [ $? -ne 0 ]; then
    echo "❌ Ошибка сборки!"
    exit 1
fi

echo "✅ Билд успешно создан: dist/syspulse"

echo "📄 Копирование статических файлов..."
cp web/static/index.html dist/web/static/
cp web/static/css/style.css dist/web/static/css/
cp web/static/js/app.js dist/web/static/js/

echo "🚀 Создание скрипта запуска..."
cat > dist/run.sh << 'EOF'
#!/bin/bash

echo "🚀 Запуск SysPulse Monitor..."
echo "📊 Системный мониторинг в реальном времени"
echo ""

# Запускаем сервер
./syspulse &
SERVER_PID=$!

# Ждем запуска сервера
sleep 2

# Открываем браузер
echo "🌐 Открываю браузер..."
if command -v xdg-open > /dev/null; then
    xdg-open http://localhost:8080
elif command -v open > /dev/null; then
    open http://localhost:8080
else
    echo "📱 Откройте в браузере: http://localhost:8080"
fi

echo ""
echo "✅ Сервер запущен на http://localhost:8080"
echo "📊 Отслеживание метрик..."
echo "⏹️  Для остановки нажмите Ctrl+C"
echo ""

# Ждем завершения
wait $SERVER_PID
EOF

chmod +x dist/run.sh
chmod +x dist/syspulse

echo "🎉 Дистрибутив готов!"
echo ""
echo "📁 Содержимое dist:"
ls -la dist/
echo ""
echo "🚀 Для запуска:"
echo "   cd dist && ./run.sh"
