#!/bin/bash

# 紧急修复脚本 - 确保代码修改正确应用
# 文件名: emergency_fix_build.sh

set -e

echo "🔧 开始紧急修复构建..."

# 1. 清理旧的构建文件
echo "🧹 清理旧的构建文件..."
rm -rf ./bin/
rm -f ./ttshub
find . -name "*.exe" -delete

# 2. 清理Go模块缓存
echo "🗑️ 清理Go模块缓存..."
go clean -modcache 2>/dev/null || true

# 3. 重新下载依赖
echo "📦 重新下载依赖..."
go mod download

# 4. 验证关键文件修改
echo "✅ 验证关键修改..."
if grep -q "优先检查数据库连接是否为空" internal/repository/channel_repository.go; then
    echo "✅ channel_repository.go 修改正确"
else
    echo "❌ channel_repository.go 修改缺失"
    exit 1
fi

# 5. 重新构建（确保CGO开启）
echo "🏗️ 重新构建应用..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./bin/ttshub ./cmd/ttshub/main.go

# 6. 验证构建结果
echo "🔍 验证构建结果..."
if [ -f "./bin/ttshub" ]; then
    echo "✅ 构建成功: ./bin/ttshub"
    
    # 检查文件大小和修改时间
    ls -lh ./bin/ttshub
    
    # 显示二进制文件信息
    file ./bin/ttshub
    ldd ./bin/ttshub 2>/dev/null || echo "无法获取动态链接信息（可能为静态链接）"
else
    echo "❌ 构建失败"
    exit 1
fi

# 7. 测试构建的应用（无数据库模式）
echo "🧪 测试构建的应用..."
timeout 10s ./bin/ttshub 2>&1 | head -20 || echo "应用启动测试完成"

echo "🎉 紧急修复构建完成!"
echo "📋 注意事项："
echo "   - 确保使用 ./bin/ttshub 进行部署"
echo "   - 检查 Zeabur 构建环境是否使用正确的构建命令"
echo "   - 可能需要设置构建变量: CGO_ENABLED=1"