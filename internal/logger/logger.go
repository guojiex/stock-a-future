package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 统一日志接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// 带上下文的日志方法
	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	WarnCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)

	// 格式化日志方法（兼容现有代码）
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	// 获取子日志器
	With(fields ...Field) Logger
	WithRequestID(requestID string) Logger
}

// Field 日志字段
type Field = zap.Field

// 常用字段构造函数
var (
	String     = zap.String
	Int        = zap.Int
	Int64      = zap.Int64
	Float64    = zap.Float64
	Bool       = zap.Bool
	Duration   = zap.Duration
	Time       = zap.Time
	ErrorField = zap.Error
	Any        = zap.Any
	RequestID  = func(id string) Field { return zap.String("request_id", id) }
	UserID     = func(id string) Field { return zap.String("user_id", id) }
	Operation  = func(op string) Field { return zap.String("operation", op) }
	StockCode  = func(code string) Field { return zap.String("stock_code", code) }
)

// zapLogger zap日志器实现
type zapLogger struct {
	logger *zap.Logger
}

// Config 日志配置
type Config struct {
	Level      string `json:"level"`       // 日志级别: debug, info, warn, error
	Format     string `json:"format"`      // 日志格式: json, console, simple
	Output     string `json:"output"`      // 输出目标: stdout, stderr, file, both
	Filename   string `json:"filename"`    // 文件名（当output为file时）
	MaxSize    int    `json:"max_size"`    // 最大文件大小(MB)
	MaxBackups int    `json:"max_backups"` // 最大备份文件数
	MaxAge     int    `json:"max_age"`     // 最大保留天数
	Compress   bool   `json:"compress"`    // 是否压缩

	// 终端输出配置
	ConsoleFormat string `json:"console_format"` // 终端格式: simple, detailed
	ShowCaller    bool   `json:"show_caller"`    // 是否显示调用位置
	ShowTime      bool   `json:"show_time"`      // 是否显示时间
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:         "info",
		Format:        "console",
		Output:        "stdout",
		MaxSize:       100,
		MaxBackups:    3,
		MaxAge:        28,
		Compress:      true,
		ConsoleFormat: "simple",
		ShowCaller:    false,
		ShowTime:      true,
	}
}

// NewLogger 创建新的日志器
func NewLogger(config *Config) (Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 解析日志级别
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("无效的日志级别 %s: %w", config.Level, err)
	}

	var cores []zapcore.Core

	// 根据输出配置创建不同的核心
	switch config.Output {
	case "stdout", "stderr":
		core, err := createConsoleCore(config, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, core)
	case "file":
		core, err := createFileCore(config, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, core)
	case "both":
		// 终端输出：简化格式
		consoleCore, err := createConsoleCore(config, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, consoleCore)

		// 文件输出：详细格式
		fileCore, err := createFileCore(config, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, fileCore)
	default:
		return nil, fmt.Errorf("不支持的输出目标: %s", config.Output)
	}

	// 合并多个核心
	var core zapcore.Core
	if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// 创建日志器选项
	options := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}

	// 根据配置决定是否添加调用者信息
	if config.ShowCaller || config.Format == "json" || config.Output == "file" || config.Output == "both" {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	// 创建日志器
	logger := zap.New(core, options...)

	return &zapLogger{logger: logger}, nil
}

// createConsoleCore 创建终端输出核心
func createConsoleCore(config *Config, level zapcore.Level) (zapcore.Core, error) {
	// 终端编码器配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// 根据终端格式配置编码器
	if config.ConsoleFormat == "simple" {
		// 简化格式：只显示时间、级别和消息
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if config.ShowTime {
			encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		} else {
			encoderConfig.TimeKey = ""
		}
		if !config.ShowCaller {
			encoderConfig.CallerKey = ""
		} else {
			encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		}
	} else {
		// 详细格式
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		if config.ShowCaller {
			encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		} else {
			encoderConfig.CallerKey = ""
		}
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 输出目标
	var writeSyncer zapcore.WriteSyncer
	if config.Output == "stderr" || config.Output == "both" {
		writeSyncer = zapcore.AddSync(os.Stderr)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	return zapcore.NewCore(encoder, writeSyncer, level), nil
}

// createFileCore 创建文件输出核心
func createFileCore(config *Config, level zapcore.Level) (zapcore.Core, error) {
	if config.Filename == "" {
		return nil, fmt.Errorf("文件输出模式下必须指定文件名")
	}

	// 文件编码器配置（详细格式）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 文件输出
	file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, fmt.Errorf("无法打开日志文件 %s: %w", config.Filename, err)
	}
	writeSyncer := zapcore.AddSync(file)

	return zapcore.NewCore(encoder, writeSyncer, level), nil
}

// Debug 调试日志
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info 信息日志
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn 警告日志
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error 错误日志
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// DebugCtx 带上下文的调试日志
func (l *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...Field) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Debug(msg, fields...)
}

// InfoCtx 带上下文的信息日志
func (l *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Info(msg, fields...)
}

// WarnCtx 带上下文的警告日志
func (l *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...Field) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Warn(msg, fields...)
}

// ErrorCtx 带上下文的错误日志
func (l *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Error(msg, fields...)
}

// Debugf 格式化调试日志
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Infof 格式化信息日志
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warnf 格式化警告日志
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化错误日志
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Fatalf 格式化致命错误日志
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

// With 创建带字段的子日志器
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// WithRequestID 创建带请求ID的子日志器
func (l *zapLogger) WithRequestID(requestID string) Logger {
	return l.With(RequestID(requestID))
}

// addContextFields 从上下文中添加字段
func (l *zapLogger) addContextFields(ctx context.Context, fields []Field) []Field {
	// 从上下文中提取请求ID
	if requestID := GetRequestIDFromContext(ctx); requestID != "" {
		fields = append(fields, RequestID(requestID))
	}

	// 从上下文中提取用户ID
	if userID := GetUserIDFromContext(ctx); userID != "" {
		fields = append(fields, UserID(userID))
	}

	return fields
}

// 上下文键类型
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// WithRequestIDContext 在上下文中设置请求ID
func WithRequestIDContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithUserIDContext 在上下文中设置用户ID
func WithUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetRequestIDFromContext 从上下文中获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// 全局日志器实例
var globalLogger Logger

// InitGlobalLogger 初始化全局日志器
func InitGlobalLogger(config *Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		// 如果全局日志器未初始化，使用默认配置创建
		logger, err := NewLogger(DefaultConfig())
		if err != nil {
			panic(fmt.Sprintf("无法创建默认日志器: %v", err))
		}
		globalLogger = logger
	}
	return globalLogger
}

// 全局日志方法（兼容现有代码）
func Debug(msg string, fields ...Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	GetGlobalLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

func ErrorLog(msg string, fields ...Field) {
	GetGlobalLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	GetGlobalLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetGlobalLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetGlobalLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetGlobalLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetGlobalLogger().Fatalf(format, args...)
}
