# PowerShell скрипт для генерации PNG изображений из PlantUML диаграмм

Write-Host "🖼️  Generating PlantUML diagrams..." -ForegroundColor Green

# Создаем папку для изображений если её нет
if (!(Test-Path "diagrams\images")) {
    New-Item -ItemType Directory -Path "diagrams\images" | Out-Null
}

# Функция генерации через PlantUML Server
function Generate-Via-Server {
    Write-Host "📡 Generating via PlantUML Server..." -ForegroundColor Yellow
    
    $pumlFiles = Get-ChildItem -Path "diagrams\*.puml"
    
    foreach ($file in $pumlFiles) {
        $filename = $file.BaseName
        Write-Host "  Processing: $filename" -ForegroundColor Gray
        
        try {
            $content = Get-Content $file.FullName -Raw
            $response = Invoke-RestMethod -Uri "http://www.plantuml.com/plantuml/png/" -Method Post -Body $content -ContentType "text/plain"
            
            # Сохраняем изображение
            $outputPath = "diagrams\images\$filename.png"
            [System.IO.File]::WriteAllBytes($outputPath, $response)
            
            Write-Host "  ✅ Generated: $outputPath" -ForegroundColor Green
        }
        catch {
            Write-Host "  ❌ Failed: $filename - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
}

# Функция генерации через локальный PlantUML
function Generate-Via-Local {
    if (Get-Command java -ErrorAction SilentlyContinue) {
        if (Test-Path "plantuml.jar") {
            Write-Host "☕ Generating via local PlantUML..." -ForegroundColor Yellow
            java -jar plantuml.jar -tpng -o images diagrams\*.puml
        }
        else {
            Write-Host "⬇️  Downloading PlantUML..." -ForegroundColor Yellow
            Invoke-WebRequest -Uri "http://sourceforge.net/projects/plantuml/files/plantuml.jar/download" -OutFile "plantuml.jar"
            java -jar plantuml.jar -tpng -o images diagrams\*.puml
        }
        return $true
    }
    else {
        Write-Host "❌ Java not found for local generation" -ForegroundColor Red
        return $false
    }
}

# Пробуем генерацию через сервер
try {
    Generate-Via-Server
    Write-Host "✅ Successfully generated diagrams via server" -ForegroundColor Green
}
catch {
    Write-Host "❌ Server generation failed, trying local..." -ForegroundColor Red
    
    if (!(Generate-Via-Local)) {
        Write-Host "❌ All generation methods failed" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "📁 Generated images saved in: diagrams\images\" -ForegroundColor Cyan
Write-Host "🔗 You can now reference them in markdown as:" -ForegroundColor Cyan
Write-Host "   ![Diagram](diagrams/images/diagram_name.png)" -ForegroundColor Gray