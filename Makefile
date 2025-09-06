# Makefile for Smart Home Pro project

.PHONY: diagrams diagrams-server diagrams-local clean-diagrams help

# Генерация диаграмм
diagrams: diagrams-server

# Генерация через PlantUML сервер (по умолчанию)
diagrams-server:
	@echo "🖼️  Generating PlantUML diagrams via server..."
	@mkdir -p diagrams/images
	@for file in diagrams/*.puml; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file" .puml); \
			echo "  Processing: $$filename"; \
			curl -X POST --data-binary "@$$file" \
				"http://www.plantuml.com/plantuml/png/" \
				-o "diagrams/images/$$filename.png" --silent; \
			if [ $$? -eq 0 ]; then \
				echo "  ✅ Generated: diagrams/images/$$filename.png"; \
			else \
				echo "  ❌ Failed: $$filename"; \
			fi; \
		fi; \
	done
	@echo "✅ Diagrams generated in diagrams/images/"

# Генерация через локальный PlantUML
diagrams-local:
	@echo "☕ Generating PlantUML diagrams locally..."
	@mkdir -p diagrams/images
	@if command -v java >/dev/null 2>&1; then \
		if [ ! -f plantuml.jar ]; then \
			echo "⬇️  Downloading PlantUML..."; \
			wget http://sourceforge.net/projects/plantuml/files/plantuml.jar/download -O plantuml.jar; \
		fi; \
		java -jar plantuml.jar -tpng -o images diagrams/*.puml; \
		echo "✅ Diagrams generated locally"; \
	else \
		echo "❌ Java not found. Please install Java or use 'make diagrams-server'"; \
		exit 1; \
	fi

# Очистка сгенерированных изображений
clean-diagrams:
	@echo "🗑️  Cleaning generated diagrams..."
	@rm -rf diagrams/images/*.png
	@echo "✅ Cleaned diagrams/images/"

# Docker команды
docker-build:
	@echo "🐳 Building Docker containers..."
	@cd apps && docker-compose build

docker-up:
	@echo "🚀 Starting Docker containers..."
	@cd apps && docker-compose up -d

docker-down:
	@echo "🛑 Stopping Docker containers..."
	@cd apps && docker-compose down

docker-logs:
	@echo "📋 Showing Docker logs..."
	@cd apps && docker-compose logs -f

# Помощь
help:
	@echo "Smart Home Pro - Available commands:"
	@echo ""
	@echo "📊 Diagrams:"
	@echo "  make diagrams        - Generate PNG images from PlantUML files (via server)"
	@echo "  make diagrams-server - Generate diagrams via PlantUML server (requires internet)"
	@echo "  make diagrams-local  - Generate diagrams locally (requires Java)"
	@echo "  make clean-diagrams  - Remove generated diagram images"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-build    - Build all Docker containers"
	@echo "  make docker-up       - Start all services"
	@echo "  make docker-down     - Stop all services"
	@echo "  make docker-logs     - Show container logs"
	@echo ""
	@echo "💡 Usage examples:"
	@echo "  make diagrams && git add diagrams/images/ && git commit -m 'Update diagrams'"
	@echo "  make docker-up && sleep 10 && curl http://localhost:8000/health"