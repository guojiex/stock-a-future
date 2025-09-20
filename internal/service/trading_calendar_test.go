package service

import (
	"testing"
	"time"
)

func TestTradingCalendar_IsTradingDay(t *testing.T) {
	tc := NewTradingCalendar()

	testCases := []struct {
		name     string
		date     string
		expected bool
		reason   string
	}{
		// 正常工作日
		{"正常工作日", "2024-09-20", true, "周五，非节假日"},
		{"正常工作日", "2024-09-19", true, "周四，非节假日"},

		// 周末
		{"周六", "2024-09-21", false, "周六"},
		{"周日", "2024-09-22", false, "周日"},

		// 2024年国庆节
		{"国庆节第一天", "2024-10-01", false, "国庆节"},
		{"国庆节第二天", "2024-10-02", false, "国庆节"},
		{"国庆节第三天", "2024-10-03", false, "国庆节"},
		{"国庆节第四天", "2024-10-04", false, "国庆节"},
		{"国庆节第五天", "2024-10-05", false, "国庆节"},
		{"国庆节第六天", "2024-10-06", false, "国庆节"},
		{"国庆节第七天", "2024-10-07", false, "国庆节"},

		// 国庆节后第一个工作日
		{"国庆节后", "2024-10-08", true, "国庆节后第一个工作日"},

		// 2024年春节
		{"春节", "2024-02-10", false, "春节"},
		{"春节", "2024-02-15", false, "春节"},

		// 2023年国庆节
		{"2023年国庆节", "2023-10-01", false, "2023年国庆节"},
		{"2023年国庆节", "2023-10-06", false, "2023年国庆节"},
	}

	for _, tc_case := range testCases {
		t.Run(tc_case.name, func(t *testing.T) {
			date, err := time.Parse("2006-01-02", tc_case.date)
			if err != nil {
				t.Fatalf("解析日期失败: %v", err)
			}

			result := tc.IsTradingDay(date)
			if result != tc_case.expected {
				t.Errorf("日期 %s (%s): 期望 %v, 实际 %v - %s",
					tc_case.date, tc_case.name, tc_case.expected, result, tc_case.reason)
			} else {
				t.Logf("✅ 日期 %s (%s): %v - %s",
					tc_case.date, tc_case.name, result, tc_case.reason)
			}
		})
	}
}

func TestTradingCalendar_GetTradingDaysInRange(t *testing.T) {
	tc := NewTradingCalendar()

	// 测试包含国庆节的时间段
	startDate, _ := time.Parse("2006-01-02", "2024-09-25")
	endDate, _ := time.Parse("2006-01-02", "2024-10-10")

	tradingDays := tc.GetTradingDaysInRange(startDate, endDate)

	t.Logf("2024年9月25日至10月10日期间的交易日:")
	for i, day := range tradingDays {
		t.Logf("  %d. %s (%s)", i+1, day.Format("2006-01-02"),
			[]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}[day.Weekday()])
	}

	// 验证国庆节期间没有交易日
	hasNationalDay := false
	for _, day := range tradingDays {
		if day.Month() == 10 && day.Day() >= 1 && day.Day() <= 7 {
			hasNationalDay = true
			t.Errorf("交易日列表中不应包含国庆节期间的日期: %s", day.Format("2006-01-02"))
		}
	}

	if !hasNationalDay {
		t.Logf("✅ 交易日列表正确过滤了国庆节期间")
	}

	// 验证交易日数量合理（应该少于总日历日数）
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
	tradingDaysCount := len(tradingDays)

	t.Logf("总日历日数: %d, 交易日数: %d", totalDays, tradingDaysCount)

	if tradingDaysCount >= totalDays {
		t.Errorf("交易日数量 (%d) 不应大于等于总日历日数 (%d)", tradingDaysCount, totalDays)
	}
}

func TestTradingCalendar_GetNextTradingDay(t *testing.T) {
	tc := NewTradingCalendar()

	testCases := []struct {
		name     string
		date     string
		expected string
	}{
		{"国庆节前最后一天", "2024-09-30", "2024-10-08"}, // 国庆节后第一个交易日
		{"周五", "2024-09-20", "2024-09-23"},       // 下周一
		{"周四", "2024-09-19", "2024-09-20"},       // 下一个工作日
	}

	for _, tc_case := range testCases {
		t.Run(tc_case.name, func(t *testing.T) {
			date, err := time.Parse("2006-01-02", tc_case.date)
			if err != nil {
				t.Fatalf("解析日期失败: %v", err)
			}

			expected, err := time.Parse("2006-01-02", tc_case.expected)
			if err != nil {
				t.Fatalf("解析期望日期失败: %v", err)
			}

			result := tc.GetNextTradingDay(date)
			if !result.Equal(expected) {
				t.Errorf("日期 %s 的下一个交易日: 期望 %s, 实际 %s",
					tc_case.date, tc_case.expected, result.Format("2006-01-02"))
			} else {
				t.Logf("✅ 日期 %s 的下一个交易日: %s",
					tc_case.date, result.Format("2006-01-02"))
			}
		})
	}
}
