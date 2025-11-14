@echo off
echo TTSHub 构建测试脚本
echo ===================

REM 检查 Go 是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go 环境，请先安装 Go 1.23 或更高版本
    exit /b 1
)

echo Go 环境检查通过

REM 设置 CGO_ENABLED=1 以支持 SQLite
set CGO_ENABLED=1
echo 设置 CGO_ENABLED=1 以支持 SQLite

REM 清理旧的构建文件
echo 清理旧的构建文件...
del /q bin\server >nul 2>&1

REM 创建 bin 目录（如果不存在）
if not exist bin mkdir bin

REM 下载依赖
echo 下载依赖...
go mod download

REM 构建项目
echo 构建项目...
go build -o ./bin/server ./cmd/ttshub/main.go

REM 检查构建结果
if exist bin\server.exe (
    echo 构建成功!
    echo 构建信息:
    dir bin\server.exe
) else (
    echo 构建失败!
    exit /b 1
)

echo 构建完成!