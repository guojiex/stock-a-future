package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"stock-a-future/internal/service"
	"time"
)

func main() {
	// 命令行参数
	var (
		dataDir    string
		action     string
		backupPath string
	)
	flag.StringVar(&dataDir, "data", "data", "数据目录路径")
	flag.StringVar(&action, "action", "backup", "操作类型: backup, restore, info")
	flag.StringVar(&backupPath, "backup", "", "备份文件路径")
	flag.Parse()

	log.Printf("数据库工具启动...")
	log.Printf("数据目录: %s", dataDir)
	log.Printf("操作类型: %s", action)

	// 创建数据库服务
	dbService, err := service.NewDatabaseService(dataDir)
	if err != nil {
		log.Fatalf("创建数据库服务失败: %v", err)
	}
	defer dbService.Close()

	switch action {
	case "backup":
		if err := backupDatabase(dbService, dataDir); err != nil {
			log.Fatalf("备份数据库失败: %v", err)
		}
	case "restore":
		if backupPath == "" {
			log.Fatalf("恢复操作需要指定备份文件路径 (-backup)")
		}
		if err := restoreDatabase(dbService, dataDir, backupPath); err != nil {
			log.Fatalf("恢复数据库失败: %v", err)
		}
	case "info":
		if err := showDatabaseInfo(dbService, dataDir); err != nil {
			log.Fatalf("获取数据库信息失败: %v", err)
		}
	default:
		log.Fatalf("不支持的操作类型: %s", action)
	}
}

// backupDatabase 备份数据库
func backupDatabase(dbService *service.DatabaseService, dataDir string) error {
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(dataDir, "backups")

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	backupFile := filepath.Join(backupDir, fmt.Sprintf("favorites_%s.db", timestamp))

	// 获取数据库文件路径
	dbPath := filepath.Join(dataDir, "favorites.db")

	// 复制数据库文件
	if err := copyFile(dbPath, backupFile); err != nil {
		return fmt.Errorf("复制数据库文件失败: %v", err)
	}

	log.Printf("✓ 数据库备份完成: %s", backupFile)
	return nil
}

// restoreDatabase 恢复数据库
func restoreDatabase(dbService *service.DatabaseService, dataDir, backupPath string) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupPath)
	}

	// 获取数据库文件路径
	dbPath := filepath.Join(dataDir, "stock_future.db")

	// 如果原数据库存在，先备份
	if _, err := os.Stat(dbPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupDir := filepath.Join(dataDir, "backups")
		if err := os.MkdirAll(backupDir, 0o755); err != nil {
			return fmt.Errorf("创建备份目录失败: %v", err)
		}

		oldBackup := filepath.Join(backupDir, fmt.Sprintf("stock_future_old_%s.db", timestamp))
		if err := copyFile(dbPath, oldBackup); err != nil {
			return fmt.Errorf("备份原数据库失败: %v", err)
		}
		log.Printf("原数据库已备份到: %s", oldBackup)
	}

	// 恢复数据库
	if err := copyFile(backupPath, dbPath); err != nil {
		return fmt.Errorf("恢复数据库失败: %v", err)
	}

	log.Printf("✓ 数据库恢复完成: %s", dbPath)
	return nil
}

// showDatabaseInfo 显示数据库信息
func showDatabaseInfo(dbService *service.DatabaseService, dataDir string) error {
	dbPath := filepath.Join(dataDir, "stock_future.db")

	// 检查数据库文件
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("数据库文件不存在: %s", dbPath)
		return nil
	}

	// 获取文件信息
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return fmt.Errorf("获取数据库文件信息失败: %v", err)
	}

	log.Printf("数据库文件: %s", dbPath)
	log.Printf("文件大小: %d 字节", fileInfo.Size())
	log.Printf("修改时间: %s", fileInfo.ModTime().Format("2006-01-02 15:04:05"))

	// 获取数据库统计信息
	db := dbService.GetDB()

	var favoritesCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM favorite_stocks").Scan(&favoritesCount); err != nil {
		log.Printf("获取收藏数量失败: %v", err)
	} else {
		log.Printf("收藏数量: %d", favoritesCount)
	}

	var groupsCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM favorite_groups").Scan(&groupsCount); err != nil {
		log.Printf("获取分组数量失败: %v", err)
	} else {
		log.Printf("分组数量: %d", groupsCount)
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
