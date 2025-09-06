#!/bin/bash

# Скрипт для генерации PNG изображений из PlantUML диаграмм

echo "🖼️  Generating PlantUML diagrams..."

# Создаем папку для изображений если её нет
mkdir -p diagrams/images

# Метод 1: Через PlantUML Server (требует интернет)
generate_via_server() {
    echo "📡 Generating via PlantUML Server..."
    
    for puml_file in diagrams/*.puml; do
        if [ -f "$puml_file" ]; then
            filename=$(basename "$puml_file" .puml)
            echo "  Processing: $filename"
            
            # Отправляем файл на PlantUML сервер и сохраняем PNG
            curl -X POST \
                --data-binary "@$puml_file" \
                "http://www.plantuml.com/plantuml/png/" \
                -o "diagrams/images/${filename}.png" \
                --silent
                
            if [ $? -eq 0 ]; then
                echo "  ✅ Generated: diagrams/images/${filename}.png"
            else
                echo "  ❌ Failed: $filename"
            fi
        fi
    done
}

# Метод 2: Через локальный PlantUML (требует Java и plantuml.jar)
generate_via_local() {
    if command -v java >/dev/null 2>&1; then
        if [ -f "plantuml.jar" ]; then
            echo "☕ Generating via local PlantUML..."
            java -jar plantuml.jar -tpng -o images diagrams/*.puml
        else
            echo "⬇️  Downloading PlantUML..."
            wget http://sourceforge.net/projects/plantuml/files/plantuml.jar/download -O plantuml.jar
            java -jar plantuml.jar -tpng -o images diagrams/*.puml
        fi
    else
        echo "❌ Java not found for local generation"
        return 1
    fi
}

# Метод 3: Через npm планtuml пакет
generate_via_npm() {
    if command -v plantuml >/dev/null 2>&1; then
        echo "📦 Generating via npm plantuml..."
        plantuml -tpng -o images diagrams/*.puml
    else
        echo "📦 Installing plantuml via npm..."
        npm install -g plantuml
        plantuml -tpng -o images diagrams/*.puml
    fi
}

# Пробуем разные методы
if generate_via_server; then
    echo "✅ Successfully generated diagrams via server"
elif generate_via_npm; then
    echo "✅ Successfully generated diagrams via npm"
elif generate_via_local; then
    echo "✅ Successfully generated diagrams via local PlantUML"
else
    echo "❌ All generation methods failed"
    exit 1
fi

echo ""
echo "📁 Generated images saved in: diagrams/images/"
echo "🔗 You can now reference them in markdown as:"
echo "   ![Diagram](diagrams/images/diagram_name.png)"