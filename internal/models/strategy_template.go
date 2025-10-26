package models

import "time"

// StrategyTemplate 策略模板
type StrategyTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        StrategyType           `json:"strategy_type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Code        string                 `json:"code,omitempty"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ParameterDefinition 参数定义
type ParameterDefinition struct {
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name"`
	Type         string      `json:"type"` // int, float, string, bool, select
	DefaultValue interface{} `json:"default_value"`
	MinValue     interface{} `json:"min_value,omitempty"`
	MaxValue     interface{} `json:"max_value,omitempty"`
	Options      []string    `json:"options,omitempty"`
	Required     bool        `json:"required"`
	Description  string      `json:"description"`
}

// StrategyTypeDefinition 策略类型定义
type StrategyTypeDefinition struct {
	Type        StrategyType          `json:"type"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  []ParameterDefinition `json:"parameters"`
}
