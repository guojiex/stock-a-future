package models

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewJSONDecimal(t *testing.T) {
	// 测试创建JSONDecimal
	value := decimal.NewFromFloat(123.45)
	jsonDecimal := NewJSONDecimal(value)

	if jsonDecimal.Decimal != value {
		t.Errorf("NewJSONDecimal应该正确设置Decimal值，期望 %v，实际 %v", value, jsonDecimal.Decimal)
	}
}

func TestJSONDecimal_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    decimal.Decimal
		expected string
	}{
		{
			name:     "正数",
			value:    decimal.NewFromFloat(123.45),
			expected: "123.45",
		},
		{
			name:     "零",
			value:    decimal.Zero,
			expected: "0",
		},
		{
			name:     "负数",
			value:    decimal.NewFromFloat(-123.45),
			expected: "-123.45",
		},
		{
			name:     "大数",
			value:    decimal.NewFromFloat(999999.99),
			expected: "999999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonDecimal := NewJSONDecimal(tt.value)
			data, err := jsonDecimal.MarshalJSON()

			if err != nil {
				t.Fatalf("MarshalJSON不应该返回错误: %v", err)
			}

			result := string(data)
			if result != tt.expected {
				t.Errorf("MarshalJSON结果不正确，期望 %s，实际 %s", tt.expected, result)
			}
		})
	}
}

func TestJSONDecimal_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected decimal.Decimal
		hasError bool
	}{
		{
			name:     "正数",
			input:    "123.45",
			expected: decimal.NewFromFloat(123.45),
			hasError: false,
		},
		{
			name:     "零",
			input:    "0",
			expected: decimal.Zero,
			hasError: false,
		},
		{
			name:     "负数",
			input:    "-123.45",
			expected: decimal.NewFromFloat(-123.45),
			hasError: false,
		},
		{
			name:     "整数",
			input:    "123",
			expected: decimal.NewFromInt(123),
			hasError: false,
		},
		{
			name:     "无效JSON",
			input:    "invalid",
			expected: decimal.Zero,
			hasError: true,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: decimal.Zero,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonDecimal := &JSONDecimal{}
			err := jsonDecimal.UnmarshalJSON([]byte(tt.input))

			if tt.hasError {
				if err == nil {
					t.Error("UnmarshalJSON应该返回错误")
				}
			} else {
				if err != nil {
					t.Fatalf("UnmarshalJSON不应该返回错误: %v", err)
				}

				if !jsonDecimal.Decimal.Equal(tt.expected) {
					t.Errorf("UnmarshalJSON结果不正确，期望 %v，实际 %v", tt.expected, jsonDecimal.Decimal)
				}
			}
		})
	}
}

func TestStockDaily_Fields(t *testing.T) {
	// 测试StockDaily结构体字段
	stock := StockDaily{
		TSCode:    "000001.SZ",
		TradeDate: "20240101",
		Open:      NewJSONDecimal(decimal.NewFromFloat(10.0)),
		High:      NewJSONDecimal(decimal.NewFromFloat(10.5)),
		Low:       NewJSONDecimal(decimal.NewFromFloat(9.8)),
		Close:     NewJSONDecimal(decimal.NewFromFloat(10.2)),
		Vol:       NewJSONDecimal(decimal.NewFromFloat(1000)),
	}

	if stock.TSCode != "000001.SZ" {
		t.Errorf("TSCode字段不正确，期望 000001.SZ，实际 %s", stock.TSCode)
	}

	if stock.TradeDate != "20240101" {
		t.Errorf("TradeDate字段不正确，期望 20240101，实际 %s", stock.TradeDate)
	}

	if !stock.Open.Decimal.Equal(decimal.NewFromFloat(10.0)) {
		t.Errorf("Open字段不正确，期望 10.0，实际 %v", stock.Open.Decimal)
	}

	if !stock.High.Decimal.Equal(decimal.NewFromFloat(10.5)) {
		t.Errorf("High字段不正确，期望 10.5，实际 %v", stock.High.Decimal)
	}

	if !stock.Low.Decimal.Equal(decimal.NewFromFloat(9.8)) {
		t.Errorf("Low字段不正确，期望 9.8，实际 %v", stock.Low.Decimal)
	}

	if !stock.Close.Decimal.Equal(decimal.NewFromFloat(10.2)) {
		t.Errorf("Close字段不正确，期望 10.2，实际 %v", stock.Close.Decimal)
	}

	if !stock.Vol.Decimal.Equal(decimal.NewFromFloat(1000)) {
		t.Errorf("Vol字段不正确，期望 1000，实际 %v", stock.Vol.Decimal)
	}
}

func TestStockBasic_Fields(t *testing.T) {
	// 测试StockBasic结构体字段
	stock := StockBasic{
		TSCode:   "000001.SZ",
		Symbol:   "000001",
		Name:     "平安银行",
		Area:     "深圳",
		Industry: "银行",
		Market:   "主板",
		ListDate: "19910403",
	}

	if stock.TSCode != "000001.SZ" {
		t.Errorf("TSCode字段不正确，期望 000001.SZ，实际 %s", stock.TSCode)
	}

	if stock.Symbol != "000001" {
		t.Errorf("Symbol字段不正确，期望 000001，实际 %s", stock.Symbol)
	}

	if stock.Name != "平安银行" {
		t.Errorf("Name字段不正确，期望 平安银行，实际 %s", stock.Name)
	}

	if stock.Area != "深圳" {
		t.Errorf("Area字段不正确，期望 深圳，实际 %s", stock.Area)
	}

	if stock.Industry != "银行" {
		t.Errorf("Industry字段不正确，期望 银行，实际 %s", stock.Industry)
	}

	if stock.Market != "主板" {
		t.Errorf("Market字段不正确，期望 主板，实际 %s", stock.Market)
	}

	if stock.ListDate != "19910403" {
		t.Errorf("ListDate字段不正确，期望 19910403，实际 %s", stock.ListDate)
	}
}
