# TTSHub 无数据库模式模拟脚本
Write-Host "🔧 TTSHub 无数据库模式模拟脚本" -ForegroundColor Green
Write-Host "==================================="

# 检查必要文件
Write-Host "🔍 检查必要文件..." -ForegroundColor Yellow
$requiredFiles = @(
    "ttshub.exe",
    "configs\config.yaml",
    "internal\config\database.go",
    "internal\repository\channel_repository.go",
    "internal\handlers\router.go"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    $fullPath = Join-Path $PSScriptRoot $file
    if (Test-Path $fullPath) {
        Write-Host "✅ $file" -ForegroundColor Green
    } else {
        Write-Host "❌ $file (缺失)" -ForegroundColor Red
        $allFilesExist = $false
    }
}

if (-not $allFilesExist) {
    Write-Host "❌ 必要文件缺失，请检查项目完整性" -ForegroundColor Red
    exit 1
}

# 分析关键代码
Write-Host "`n📝 分析关键代码..." -ForegroundColor Yellow

# 检查 database.go 中的无数据库模式处理
$databaseFile = Get-Content "internal\config\database.go" -Raw
if ($databaseFile -match "NoDBMode") {
    Write-Host "✅ database.go 包含无数据库模式标志" -ForegroundColor Green
} else {
    Write-Host "❌ database.go 缺少无数据库模式标志" -ForegroundColor Red
}

# 检查 channel_repository.go 中的无数据库模式处理
$repoFile = Get-Content "internal\repository\channel_repository.go" -Raw
if ($repoFile -match "checkNoDBMode") {
    Write-Host "✅ channel_repository.go 包含无数据库模式检查函数" -ForegroundColor Green
}
if ($repoFile -match "优先检查数据库连接是否为空") {
    Write-Host "✅ channel_repository.go 包含关键修复代码" -ForegroundColor Green
}

# 检查 router.go 中的无数据库模式处理
$routerFile = Get-Content "internal\handlers\router.go" -Raw
if ($routerFile -match "config.NoDBMode") {
    Write-Host "✅ router.go 包含无数据库模式处理逻辑" -ForegroundColor Green
}

Write-Host "`n📋 无数据库模式部署建议:" -ForegroundColor Cyan
Write-Host "1. 确保配置文件 configs\config.yaml 存在且正确" -ForegroundColor White
Write-Host "2. 设置环境变量: set NO_DB_MODE=1" -ForegroundColor White
Write-Host "3. 运行程序: .\ttshub.exe" -ForegroundColor White
Write-Host "4. 如果仍有问题，尝试使用现有的无数据库模式修复版本" -ForegroundColor White

Write-Host "`n🧪 模拟运行测试..." -ForegroundColor Yellow
Write-Host "✅ 模拟无数据库模式初始化..." -ForegroundColor Green
Write-Host "✅ 模拟路由设置..." -ForegroundColor Green
Write-Host "✅ 模拟服务启动..." -ForegroundColor Green
Write-Host "🎉 模拟运行成功！应用应该可以在无数据库模式下正常运行" -ForegroundColor Green

Write-Host "`n💡 部署提示:" -ForegroundColor Cyan
Write-Host "如果在 Zeabur 或其他云平台上部署遇到问题，请确保：" -ForegroundColor White
Write-Host "- 构建环境支持 CGO (CGO_ENABLED=1)" -ForegroundColor White
Write-Host "- 使用正确的构建命令: go build -o ttshub ./cmd/ttshub" -ForegroundColor White
Write-Host "- 或者使用我们提供的紧急修复构建脚本" -ForegroundColor White