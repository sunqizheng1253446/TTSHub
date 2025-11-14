// 测试无数据库模式修复
package main

import (
	"fmt"
	"ttshub/internal/config"
	"ttshub/internal/handlers"
	"ttshub/internal/repository"
)

func main() {
	fmt.Println("=== 测试 TTSHub 无数据库模式修复 ===")
	
	// 模拟无数据库模式
	config.NoDBMode = true
	fmt.Printf("已设置无数据库模式: %t\n", config.NoDBMode)
	
	// 创建仓库实例（传入 nil 数据库）
	fmt.Println("创建无数据库模式仓库实例...")
	channelRepo := repository.NewChannelRepository(nil)
	
	// 测试 ListAll 方法
	fmt.Println("测试 ListAll 方法...")
	channels, err := channelRepo.ListAll()
	if err != nil {
		fmt.Printf("ListAll 错误: %v\n", err)
	} else {
		fmt.Printf("ListAll 成功，返回 %d 个渠道\n", len(channels))
	}
	
	// 测试配置缓存预加载
	fmt.Println("测试 PreloadConfigCache...")
	err = repository.PreloadConfigCache(channelRepo)
	if err != nil {
		fmt.Printf("PreloadConfigCache 错误: %v\n", err)
	} else {
		fmt.Println("PreloadConfigCache 成功")
	}
	
	// 测试路由设置
	fmt.Println("测试路由设置...")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("路由设置时发生 panic: %v\n", r)
		} else {
			fmt.Println("路由设置成功，没有 panic")
		}
	}()
	
	router := handlers.SetupRouter()
	if router != nil {
		fmt.Println("路由设置成功完成")
	}
	
	fmt.Println("=== 测试完成 ===")
}