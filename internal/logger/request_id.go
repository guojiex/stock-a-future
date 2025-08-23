package logger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// RequestIDGenerator 请求ID生成器
type RequestIDGenerator interface {
	Generate() string
}

// UUIDGenerator UUID风格的请求ID生成器
type UUIDGenerator struct{}

// Generate 生成UUID风格的请求ID
func (g *UUIDGenerator) Generate() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}

	// 设置版本和变体位
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// ShortIDGenerator 短ID生成器（性能更好）
type ShortIDGenerator struct {
	counter uint64
}

// Generate 生成短ID
func (g *ShortIDGenerator) Generate() string {
	// 使用时间戳（秒）+ 原子计数器
	timestamp := time.Now().Unix()
	counter := atomic.AddUint64(&g.counter, 1)

	// 生成8字节随机数
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用计数器
		randomBytes = []byte{
			byte(counter >> 24),
			byte(counter >> 16),
			byte(counter >> 8),
			byte(counter),
		}
	}

	return fmt.Sprintf("%x%04x%s",
		timestamp&0xffffffff,            // 时间戳的低32位
		counter&0xffff,                  // 计数器的低16位
		hex.EncodeToString(randomBytes)) // 4字节随机数
}

// 默认请求ID生成器
var defaultGenerator RequestIDGenerator = &ShortIDGenerator{}

// SetDefaultGenerator 设置默认请求ID生成器
func SetDefaultGenerator(generator RequestIDGenerator) {
	defaultGenerator = generator
}

// GenerateRequestID 生成请求ID
func GenerateRequestID() string {
	return defaultGenerator.Generate()
}
