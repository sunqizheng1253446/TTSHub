@echo off
echo 🔧 TTSHub Windows 构建脚本
echo =========================

REM 检查 Go 是否安装
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Go 环境，请先安装 Go 1.23 或更高版本
    exit /b 1
)

echo ✅ Go 环境检查通过

REM 设置 CGO_ENABLED=1 以支持 SQLite
set CGO_ENABLED=1
echo 🏗️  设置 CGO_ENABLED=1 以支持 SQLite

REM 清理旧的构建文件
echo 🧹 清理旧的构建文件...
del /q ttshub.exe >nul 2>&1
del /q ttshub_new.exe >nul 2>&1

REM 下载依赖
echo 📦 下载依赖...
go mod download

REM 构建项目
echo 🏗️  构建项目...
go build -o ttshub_new.exe ./cmd/ttshub

REM 检查构建结果
if exist ttshub_new.exe (
    echo ✅ 构建成功!
    echo 📝 构建信息:
    dir ttshub_new.exe | findstr ttshub_new.exe
    echo.
    echo 🚀 可以通过以下命令运行:
    echo    .\ttshub_new.exe
    echo.
    echo 💡 提示: 如果要替换旧版本，可以执行:
    echo    move ttshub.exe ttshub_old.exe
    echo    move ttshub_new.exe ttshub.exe
) else (
    echo ❌ 构建失败!
    exit /b 1
)

echo 🎉 构建完成!