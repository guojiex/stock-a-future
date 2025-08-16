package client

import (
	"testing"
)

func TestAKToolsClient_FormatDateForFrontend(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "标准日期格式 2025-08-15",
			input:    "2025-08-15",
			expected: "20250815",
		},
		{
			name:     "月份为一位数 2025-8-15",
			input:    "2025-8-15",
			expected: "20250815",
		},
		{
			name:     "日期为一位数 2025-08-5",
			input:    "2025-08-5",
			expected: "20250805",
		},
		{
			name:     "月份和日期都为一位数 2025-8-5",
			input:    "2025-8-5",
			expected: "20250805",
		},
		{
			name:     "只有年月 2025-08-",
			input:    "2025-08-",
			expected: "20250801",
		},
		{
			name:     "只有年月且月份为一位数 2025-8-",
			input:    "2025-8-",
			expected: "20250801",
		},
		{
			name:     "已经是8位数字格式",
			input:    "20250815",
			expected: "20250815",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "包含空格 2025-08-15 ",
			input:    "2025-08-15 ",
			expected: "20250815",
		},
		{
			name:     "包含换行符 2025-08-15\n",
			input:    "2025-08-15\n",
			expected: "20250815",
		},
		{
			name:     "无法解析的格式",
			input:    "invalid-date",
			expected: "invalid-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.formatDateForFrontend(tt.input)
			if result != tt.expected {
				t.Errorf("formatDateForFrontend(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAKToolsClient_IsNumeric(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "纯数字字符串",
			input:    "12345678",
			expected: true,
		},
		{
			name:     "包含字母的字符串",
			input:    "12345abc",
			expected: false,
		},
		{
			name:     "包含特殊字符的字符串",
			input:    "12345-67",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: true,
		},
		{
			name:     "单个数字",
			input:    "5",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("isNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAKToolsClient_CleanStockSymbol(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "上海市场代码",
			input:    "601225.SH",
			expected: "601225",
		},
		{
			name:     "深圳市场代码",
			input:    "000001.SZ",
			expected: "000001",
		},
		{
			name:     "北京市场代码",
			input:    "430139.BJ",
			expected: "430139",
		},
		{
			name:     "小写后缀",
			input:    "601225.sh",
			expected: "601225",
		},
		{
			name:     "无后缀代码",
			input:    "601225",
			expected: "601225",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.CleanStockSymbol(tt.input)
			if result != tt.expected {
				t.Errorf("CleanStockSymbol(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAKToolsClient_DetermineTSCode(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "上海市场 600xxx",
			input:    "600000",
			expected: "600000.SH",
		},
		{
			name:     "上海市场 601xxx",
			input:    "601225",
			expected: "601225.SH",
		},
		{
			name:     "上海市场 603xxx",
			input:    "603000",
			expected: "603000.SH",
		},
		{
			name:     "上海市场 688xxx",
			input:    "688001",
			expected: "688001.SH",
		},
		{
			name:     "深圳市场 000xxx",
			input:    "000001",
			expected: "000001.SZ",
		},
		{
			name:     "深圳市场 002xxx",
			input:    "002001",
			expected: "002001.SZ",
		},
		{
			name:     "深圳市场 300xxx",
			input:    "300001",
			expected: "300001.SZ",
		},
		{
			name:     "北京市场 430xxx",
			input:    "430001",
			expected: "430001.BJ",
		},
		{
			name:     "北京市场 830xxx",
			input:    "830001",
			expected: "830001.BJ",
		},
		{
			name:     "北京市场 870xxx",
			input:    "870001",
			expected: "870001.BJ",
		},
		{
			name:     "已有后缀的代码",
			input:    "601225.SH",
			expected: "601225.SH",
		},
		{
			name:     "未知代码格式",
			input:    "999999",
			expected: "999999.SH", // 默认上海市场
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.DetermineTSCode(tt.input)
			if result != tt.expected {
				t.Errorf("DetermineTSCode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
