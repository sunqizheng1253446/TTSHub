# TTSHub Deployment Check Script
Write-Host "TTSHub Deployment Check" -ForegroundColor Green
Write-Host "======================"

# Check if Go is installed
Write-Host "Checking if Go is installed..." -ForegroundColor Yellow
$goInstalled = $false
try {
    $goVersion = go version 2>$null
    if ($goVersion) {
        Write-Host "OK - Go is installed: $goVersion" -ForegroundColor Green
        $goInstalled = $true
    } else {
        Write-Host "WARNING - Go is not installed or not in PATH" -ForegroundColor Yellow
    }
} catch {
    Write-Host "WARNING - Go is not installed or not in PATH" -ForegroundColor Yellow
}

# Check required files (excluding ttshub.exe if Go is not available)
Write-Host "`nChecking required files..." -ForegroundColor Yellow
$requiredFiles = @(
    "configs\config.yaml",
    "internal\config\database.go",
    "internal\repository\channel_repository.go",
    "internal\handlers\router.go"
)

if ($goInstalled) {
    $requiredFiles += "ttshub.exe"
}

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    $fullPath = Join-Path $PSScriptRoot $file
    if (Test-Path $fullPath) {
        Write-Host "OK - $file" -ForegroundColor Green
    } else {
        Write-Host "MISSING - $file" -ForegroundColor Red
        $allFilesExist = $false
    }
}

# Analyze key code
Write-Host "`nAnalyzing key code..." -ForegroundColor Yellow

# Check database.go for NoDBMode handling
try {
    $databaseFile = Get-Content "internal\config\database.go" -Raw
    if ($databaseFile -match "NoDBMode") {
        Write-Host "OK - database.go contains NoDBMode flag" -ForegroundColor Green
    } else {
        Write-Host "ERROR - database.go missing NoDBMode flag" -ForegroundColor Red
    }
} catch {
    Write-Host "ERROR - Could not read database.go" -ForegroundColor Red
}

# Check channel_repository.go for NoDBMode handling
try {
    $repoFile = Get-Content "internal\repository\channel_repository.go" -Raw
    if ($repoFile -match "checkNoDBMode") {
        Write-Host "OK - channel_repository.go contains NoDBMode check function" -ForegroundColor Green
    }
    if ($repoFile -match "优先检查数据库连接是否为空") {
        Write-Host "OK - channel_repository.go contains key fix code" -ForegroundColor Green
    }
} catch {
    Write-Host "ERROR - Could not read channel_repository.go" -ForegroundColor Red
}

# Check router.go for NoDBMode handling
try {
    $routerFile = Get-Content "internal\handlers\router.go" -Raw
    if ($routerFile -match "config.NoDBMode") {
        Write-Host "OK - router.go contains NoDBMode handling logic" -ForegroundColor Green
    }
} catch {
    Write-Host "ERROR - Could not read router.go" -ForegroundColor Red
}

Write-Host "`nDeployment Recommendations:" -ForegroundColor Cyan
Write-Host "1. Ensure configs\config.yaml exists and is correct" -ForegroundColor White
Write-Host "2. If Go is installed, build the application with: go build -o ttshub.exe ./cmd/ttshub" -ForegroundColor White
Write-Host "3. Set environment variable for NoDB mode: `$env:NO_DB_MODE=1" -ForegroundColor White
Write-Host "4. Run the application: .\ttshub.exe" -ForegroundColor White

if (-not $goInstalled) {
    Write-Host "`nWARNING: Go is not installed. You cannot rebuild the application." -ForegroundColor Yellow
    Write-Host "Please install Go 1.23 or higher from https://golang.org/dl/" -ForegroundColor Yellow
}

Write-Host "`nApplication Status:" -ForegroundColor Yellow
Write-Host "The application has been designed to work in NoDB mode when database initialization fails." -ForegroundColor White
Write-Host "This prevents nil pointer panics and allows the service to continue running." -ForegroundColor White

Write-Host "`nNext Steps:" -ForegroundColor Cyan
Write-Host "1. If you have Go installed, try rebuilding with CGO support:" -ForegroundColor White
Write-Host "   `$env:CGO_ENABLED=1; go build -o ttshub.exe ./cmd/ttshub" -ForegroundColor White
Write-Host "2. If you don't have Go, you'll need to obtain a pre-built binary." -ForegroundColor White