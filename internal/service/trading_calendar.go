package service

import (
	"time"
)

// TradingCalendar 交易日历服务
type TradingCalendar struct {
	holidays map[string]bool // 存储节假日，格式为 "YYYY-MM-DD"
}

// NewTradingCalendar 创建新的交易日历服务
func NewTradingCalendar() *TradingCalendar {
	tc := &TradingCalendar{
		holidays: make(map[string]bool),
	}

	// 初始化多年的节假日
	tc.initHolidays2023()
	tc.initHolidays2024()
	tc.initHolidays2025()

	return tc
}

// initHolidays2023 初始化2023年的节假日
func (tc *TradingCalendar) initHolidays2023() {
	// 2023年节假日（A股休市日）
	holidays2023 := []string{
		// 元旦
		"2023-01-01", "2023-01-02", "2023-01-03",

		// 春节
		"2023-01-21", "2023-01-22", "2023-01-23", "2023-01-24",
		"2023-01-25", "2023-01-26", "2023-01-27",

		// 清明节
		"2023-04-05",

		// 劳动节
		"2023-05-01", "2023-05-02", "2023-05-03",

		// 端午节
		"2023-06-22", "2023-06-23", "2023-06-24",

		// 中秋节、国庆节
		"2023-09-29", "2023-09-30", "2023-10-01", "2023-10-02",
		"2023-10-03", "2023-10-04", "2023-10-05", "2023-10-06",
	}

	for _, holiday := range holidays2023 {
		tc.holidays[holiday] = true
	}
}

// initHolidays2024 初始化2024年的节假日
func (tc *TradingCalendar) initHolidays2024() {
	// 2024年节假日（A股休市日）
	holidays2024 := []string{
		// 元旦
		"2024-01-01",

		// 春节
		"2024-02-10", "2024-02-11", "2024-02-12", "2024-02-13",
		"2024-02-14", "2024-02-15", "2024-02-16", "2024-02-17",

		// 清明节
		"2024-04-04", "2024-04-05", "2024-04-06",

		// 劳动节
		"2024-05-01", "2024-05-02", "2024-05-03", "2024-05-04", "2024-05-05",

		// 端午节
		"2024-06-10",

		// 中秋节
		"2024-09-15", "2024-09-16", "2024-09-17",

		// 国庆节
		"2024-10-01", "2024-10-02", "2024-10-03", "2024-10-04",
		"2024-10-05", "2024-10-06", "2024-10-07",
	}

	for _, holiday := range holidays2024 {
		tc.holidays[holiday] = true
	}
}

// IsWorkingDay 判断是否为工作日（周一到周五，且非节假日）
func (tc *TradingCalendar) IsWorkingDay(date time.Time) bool {
	// 检查是否为周末
	weekday := date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// 检查是否为节假日
	dateStr := date.Format("2006-01-02")
	if tc.holidays[dateStr] {
		return false
	}

	return true
}

// IsTradingDay 判断是否为交易日（工作日的别名，更符合股票交易语境）
func (tc *TradingCalendar) IsTradingDay(date time.Time) bool {
	return tc.IsWorkingDay(date)
}

// GetNextTradingDay 获取下一个交易日
func (tc *TradingCalendar) GetNextTradingDay(date time.Time) time.Time {
	nextDay := date.AddDate(0, 0, 1)
	for !tc.IsTradingDay(nextDay) {
		nextDay = nextDay.AddDate(0, 0, 1)
	}
	return nextDay
}

// GetPreviousTradingDay 获取前一个交易日
func (tc *TradingCalendar) GetPreviousTradingDay(date time.Time) time.Time {
	prevDay := date.AddDate(0, 0, -1)
	for !tc.IsTradingDay(prevDay) {
		prevDay = prevDay.AddDate(0, 0, -1)
	}
	return prevDay
}

// GetTradingDaysInRange 获取指定日期范围内的所有交易日
func (tc *TradingCalendar) GetTradingDaysInRange(startDate, endDate time.Time) []time.Time {
	var tradingDays []time.Time

	current := startDate
	for !current.After(endDate) {
		if tc.IsTradingDay(current) {
			tradingDays = append(tradingDays, current)
		}
		current = current.AddDate(0, 0, 1)
	}

	return tradingDays
}

// CountTradingDays 计算指定日期范围内的交易日数量
func (tc *TradingCalendar) CountTradingDays(startDate, endDate time.Time) int {
	return len(tc.GetTradingDaysInRange(startDate, endDate))
}

// initHolidays2025 初始化2025年的节假日
func (tc *TradingCalendar) initHolidays2025() {
	// 2025年节假日（A股休市日）
	holidays2025 := []string{
		// 元旦
		"2025-01-01",

		// 春节
		"2025-01-28", "2025-01-29", "2025-01-30", "2025-01-31",
		"2025-02-01", "2025-02-02", "2025-02-03", "2025-02-04",

		// 清明节
		"2025-04-04", "2025-04-05", "2025-04-06",

		// 劳动节
		"2025-05-01", "2025-05-02", "2025-05-03", "2025-05-04", "2025-05-05",

		// 端午节
		"2025-05-31", "2025-06-01", "2025-06-02",

		// 中秋节
		"2025-10-06", "2025-10-07", "2025-10-08",

		// 国庆节
		"2025-10-01", "2025-10-02", "2025-10-03", "2025-10-04", "2025-10-05",
	}

	for _, holiday := range holidays2025 {
		tc.holidays[holiday] = true
	}
}
