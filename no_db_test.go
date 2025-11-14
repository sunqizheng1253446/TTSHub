package main

import (
	"fmt"
	"os"
	"ttshub/internal/config"
	"ttshub/internal/handlers"
	"ttshub/internal/utils"
)

func main() {
	fmt.Println("=== TTSHub 无数据库模式测试 ===")
	
	// 强制设置无数据库模式
	config.NoDBMode = true
	fmt.Printf("✅ 已设置无数据库模式: %t\n", config.NoDBMode)
	
	// 初始化日志系统
	logConfig := &config.LogConfig{
		Level:  "info",
		Path:   "./logs",
		Format: "console",
	}
	
	if err := utils.InitLogger(logConfig); err != nil {
		fmt.Printf("❌ 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 日志系统初始化完成")
	
	// 测试路由设置
	fmt.Println("🔧 测试路由设置...")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ 路由设置时发生 panic: %v\n", r)
			os.Exit(1)
		} else {
			fmt.Println("✅ 路由设置成功，没有 panic")
		}
	}()
	
	router := handlers.SetupRouter()
	if router != nil {
		fmt.Println("✅ 路由设置成功完成")
	}
	
	fmt.Println("🎉 无数据库模式测试完成!")
	fmt.Println("💡 你现在可以尝试运行:")
	fmt.Println("   go run no_db_test.go")
}