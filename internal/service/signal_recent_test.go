package service

import (
	"fmt"
	"testing"
	"time"

	"stock-a-future/internal/models"
)

// TestGetRecentUpdatedSignals_MultipleSignalDates 测试多个信号日期的场景
func TestGetRecentUpdatedSignals_MultipleSignalDates(t *testing.T) {
	// 创建测试服务
	dbService, err := NewDatabaseService("../../data")
	if err != nil {
		t.Skipf("无法连接数据库，跳过测试: %v", err)
		return
	}
	defer dbService.Close()

	// 创建信号服务实例（用于测试）
	signalService := &SignalService{
		db: dbService,
	}

	t.Log("=== 测试GetRecentUpdatedSignals方法 ===")

	// 获取最近的信号
	recentSignals, err := signalService.GetRecentUpdatedSignals(100)
	if err != nil {
		t.Fatalf("获取最近信号失败: %v", err)
	}

	t.Logf("\n✅ 成功获取 %d 个最近的BUY/SELL信号", len(recentSignals))

	// 统计信号日期分布
	signalDateCount := make(map[string]int)
	stockMap := make(map[string]*models.StockSignal)

	for _, signal := range recentSignals {
		signalDateCount[signal.SignalDate]++
		stockMap[signal.TSCode] = signal

		// 记录前20个信号的详细信息
		if len(stockMap) <= 20 {
			t.Logf("  %d. %s (%s) - 信号日期: %s, 类型: %s, 强度: %s",
				len(stockMap), signal.Name, signal.TSCode,
				signal.SignalDate, signal.SignalType, signal.SignalStrength)
		}
	}

	t.Logf("\n📊 信号日期分布:")
	// 按日期排序显示
	dates := make([]string, 0, len(signalDateCount))
	for date := range signalDateCount {
		dates = append(dates, date)
	}
	// 简单排序（降序）
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] < dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	for _, date := range dates {
		count := signalDateCount[date]
		// 格式化日期显示
		formattedDate := date
		if len(date) == 8 {
			formattedDate = fmt.Sprintf("%s-%s-%s", date[0:4], date[4:6], date[6:8])
		}
		t.Logf("  %s: %d 个信号", formattedDate, count)
	}

	// 特别检查中国东航
	t.Log("\n🔍 检查中国东航信号:")
	chinaEasternSignal, found := stockMap["600115.SH"]
	if found {
		t.Logf("  ✅ 找到中国东航信号:")
		t.Logf("     股票代码: %s", chinaEasternSignal.TSCode)
		t.Logf("     股票名称: %s", chinaEasternSignal.Name)
		t.Logf("     信号日期: %s", chinaEasternSignal.SignalDate)
		t.Logf("     信号类型: %s", chinaEasternSignal.SignalType)
		t.Logf("     信号强度: %s", chinaEasternSignal.SignalStrength)
		t.Logf("     置信度: %s", chinaEasternSignal.Confidence.Decimal.String())
	} else {
		t.Log("  ⚠️  未找到中国东航信号")

		// 直接查询数据库检查中国东航的信号
		t.Log("\n📋 直接查询数据库中的中国东航信号:")
		query := `
			SELECT ts_code, name, trade_date, signal_date, signal_type, signal_strength, confidence
			FROM stock_signals 
			WHERE ts_code = '600115.SH' OR name LIKE '%东航%'
			ORDER BY signal_date DESC
			LIMIT 5
		`
		rows, err := dbService.GetDB().Query(query)
		if err != nil {
			t.Logf("  查询失败: %v", err)
		} else {
			defer rows.Close()
			count := 0
			for rows.Next() {
				var tsCode, name, tradeDate, signalDate, signalType, signalStrength string
				var confidence float64
				err = rows.Scan(&tsCode, &name, &tradeDate, &signalDate, &signalType, &signalStrength, &confidence)
				if err != nil {
					continue
				}
				count++
				t.Logf("  %d. 代码: %s, 名称: %s, 信号日期: %s, 类型: %s, 强度: %s, 置信度: %.2f",
					count, tsCode, name, signalDate, signalType, signalStrength, confidence)
			}
			if count == 0 {
				t.Log("  ⚠️  数据库中也没有中国东航的信号记录")
			}
		}
	}

	// 验证修复是否生效
	t.Log("\n✅ 验证修复效果:")
	if len(signalDateCount) > 1 {
		t.Logf("  ✅ 成功获取了多个日期的信号（%d个不同日期）", len(signalDateCount))
		t.Log("  ✅ 修复生效：不再局限于单一最新日期")
	} else if len(signalDateCount) == 1 {
		t.Log("  ⚠️  只获取到一个日期的信号，可能所有信号都在同一天")
	}

	// 检查是否有不同日期的信号
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("20060102")
	signalsBeforeToday := 0
	for date := range signalDateCount {
		if date < time.Now().Format("20060102") && date >= sevenDaysAgo {
			signalsBeforeToday++
		}
	}
	if signalsBeforeToday > 0 {
		t.Logf("  ✅ 包含了历史信号（非今天的信号数: %d）", signalsBeforeToday)
	}
}

// TestGetRecentUpdatedSignals_EachStockOnlyOnce 测试每只股票只返回一次
func TestGetRecentUpdatedSignals_EachStockOnlyOnce(t *testing.T) {
	dbService, err := NewDatabaseService("../../data")
	if err != nil {
		t.Skipf("无法连接数据库，跳过测试: %v", err)
		return
	}
	defer dbService.Close()

	signalService := &SignalService{
		db: dbService,
	}

	t.Log("=== 测试每只股票只返回一次 ===")

	recentSignals, err := signalService.GetRecentUpdatedSignals(100)
	if err != nil {
		t.Fatalf("获取最近信号失败: %v", err)
	}

	// 检查是否有重复的股票代码
	stockCodes := make(map[string]bool)
	duplicates := 0

	for _, signal := range recentSignals {
		if stockCodes[signal.TSCode] {
			duplicates++
			t.Errorf("  ❌ 发现重复股票: %s (%s)", signal.Name, signal.TSCode)
		}
		stockCodes[signal.TSCode] = true
	}

	if duplicates == 0 {
		t.Logf("  ✅ 通过：所有股票代码都是唯一的（%d只股票）", len(stockCodes))
	} else {
		t.Errorf("  ❌ 失败：发现 %d 个重复的股票", duplicates)
	}
}

// TestSignalServiceIntegration_ChinaEastern 集成测试：验证中国东航信号
func TestSignalServiceIntegration_ChinaEastern(t *testing.T) {
	dbService, err := NewDatabaseService("../../data")
	if err != nil {
		t.Skipf("无法连接数据库，跳过测试: %v", err)
		return
	}
	defer dbService.Close()

	t.Log("=== 集成测试：中国东航信号 ===")

	// 1. 检查数据库中是否存在中国东航的信号
	t.Log("\n1️⃣  检查数据库中的中国东航信号:")
	var hasSignal bool
	var latestSignalDate string
	query := `
		SELECT signal_date 
		FROM stock_signals 
		WHERE (ts_code = '600115.SH' OR name LIKE '%东航%')
		  AND signal_type IN ('BUY', 'SELL')
		ORDER BY signal_date DESC
		LIMIT 1
	`
	err = dbService.GetDB().QueryRow(query).Scan(&latestSignalDate)
	if err == nil {
		hasSignal = true
		t.Logf("  ✅ 找到中国东航信号，最新信号日期: %s", latestSignalDate)
	} else {
		t.Logf("  ⚠️  未找到中国东航的BUY/SELL信号")
	}

	if !hasSignal {
		t.Skip("数据库中没有中国东航信号，跳过后续测试")
		return
	}

	// 2. 使用GetRecentUpdatedSignals获取信号
	t.Log("\n2️⃣  使用GetRecentUpdatedSignals获取信号:")
	signalService := &SignalService{
		db: dbService,
	}

	recentSignals, err := signalService.GetRecentUpdatedSignals(100)
	if err != nil {
		t.Fatalf("获取最近信号失败: %v", err)
	}

	t.Logf("  获取到 %d 个信号", len(recentSignals))

	// 3. 检查中国东航是否在返回的信号中
	t.Log("\n3️⃣  检查中国东航是否在返回结果中:")
	found := false
	for _, signal := range recentSignals {
		if signal.TSCode == "600115.SH" || signal.Name == "中国东航" {
			found = true
			t.Log("  ✅ 成功找到中国东航信号！")
			t.Logf("     股票代码: %s", signal.TSCode)
			t.Logf("     股票名称: %s", signal.Name)
			t.Logf("     信号日期: %s", signal.SignalDate)
			t.Logf("     信号类型: %s", signal.SignalType)
			t.Logf("     信号强度: %s", signal.SignalStrength)
			break
		}
	}

	if !found {
		t.Errorf("  ❌ 未找到中国东航信号！这表明修复可能还有问题")

		// 额外调试信息
		t.Log("\n🔍 调试信息:")
		t.Logf("  数据库中最新信号日期: %s", latestSignalDate)
		t.Logf("  7天前日期: %s", time.Now().AddDate(0, 0, -7).Format("20060102"))

		// 检查信号日期是否在7天范围内
		sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("20060102")
		if latestSignalDate >= sevenDaysAgo {
			t.Log("  ⚠️  信号日期在7天范围内，应该被查询到")
		} else {
			t.Log("  ℹ️  信号日期超过7天，不在查询范围内")
		}
	}
}
