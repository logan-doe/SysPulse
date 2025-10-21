.PHONY: build build-linux build-windows build-darwin clean run

# Сборка для текущей платформы
build:
	@echo "🏗️  Building SysPulse Monitor..."
	@mkdir -p dist/web/static/{css,js}
	@go build -ldflags="-s -w -X main.version=1.0.0" -o dist/syspulse ./cmd/syspulse-server
	@cp web/static/index.html dist/web/static/
	@cp web/static/css/style.css dist/web/static/css/
	@cp web/static/js/app.js dist/web/static/js/
	@echo "✅ Build completed: dist/syspulse"

# Сборка для Linux
build-linux:
	@echo "🐧 Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=1.0.0" -o dist/syspulse-linux ./cmd/syspulse-server

# Сборка для Windows
build-windows:
	@echo "🪟 Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=1.0.0" -o dist/syspulse.exe ./cmd/syspulse-server

# Сборка для Mac
build-darwin:
	@echo "🍎 Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=1.0.0" -o dist/syspulse-macos ./cmd/syspulse-server

# Сборка всех платформ
build-all: build-linux build-windows build-darwin
	@echo "🌍 Multi-platform build completed"

# Очистка
clean:
	@echo "🧹 Cleaning..."
	@rm -rf dist/
	@echo "✅ Clean completed"

# Запуск в режиме разработки
run:
	@echo "🚀 Starting development server..."
	@go run ./cmd/syspulse-server

# Помощь
help:
	@echo "Available commands:"
	@echo "  build       - Build for current platform"
	@echo "  build-linux - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-darwin - Build for macOS"
	@echo "  build-all   - Build for all platforms"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Run in development mode"
