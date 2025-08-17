package main

import (
	"flag"
	"log"
	"path/filepath"
	"stock-a-future/internal/service"
)

func main() {
	// 命令行参数
	var dataDir string
	flag.StringVar(&dataDir, "data", "data", "数据目录路径")
	flag.Parse()

	log.Printf("开始数据迁移...")
	log.Printf("数据目录: %s", dataDir)

	// 创建数据库服务
	dbService, err := service.NewDatabaseService(dataDir)
	if err != nil {
		log.Fatalf("创建数据库服务失败: %v", err)
	}
	defer dbService.Close()

	// 迁移路径
	favoritesPath := filepath.Join(dataDir, "favorites", "favorites.json")
	groupsPath := filepath.Join(dataDir, "favorites", "groups.json")

	log.Printf("收藏文件路径: %s", favoritesPath)
	log.Printf("分组文件路径: %s", groupsPath)

	// 执行迁移
	if err := dbService.MigrateFromJSON(favoritesPath, groupsPath); err != nil {
		log.Fatalf("数据迁移失败: %v", err)
	}

	log.Printf("✓ 数据迁移完成！")
	log.Printf("数据库文件: %s", filepath.Join(dataDir, "stock_future.db"))
}
