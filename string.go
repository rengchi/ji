package ji

import (
	"math/rand"
	"strings"
)

// 随机字符生成
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// globalRand 时间种子
var (
	gr = rand.New(rand.NewSource(CurrentTimestampMilli()))
)

// Charset Charset 生成
type Charset struct {
	ch string
	r  *rand.Rand
}

// Generate 字符生成器实例
func Generate() *Charset {
	return &Charset{
		ch: chars,
		r:  gr,
	}
}

// NewCharset 创建一个支持自定义字符集的生成器
func NewCharset(charset string) *Charset {
	if charset == "" {
		charset = chars // 默认字符集
	}
	return &Charset{
		ch: charset,
		r:  gr,
	}
}

// Random 生成指定长度的随机字符串
func (c *Charset) Random(length int) string {
	var sb strings.Builder
	sb.Grow(length) // 预分配内存，避免多次分配
	for i := 0; i < length; i++ {
		randomIndex := c.r.Intn(len(c.ch))
		sb.WriteByte(c.ch[randomIndex])
	}
	return sb.String()
}

// Salt 16位密码盐
func (c *Charset) Salt() string {
	return c.Random(16)
}

// Assemble 拼接字符串
func Assemble(str ...string) string {
	var builder strings.Builder
	for _, s := range str {
		builder.WriteString(s)
	}
	return builder.String()
}

// LeftPad 在左侧填充字符
func LeftPad(s string, padStr string, totalLength int) string {
	padCount := totalLength - len(s)
	if padCount > 0 {
		return strings.Repeat(padStr, padCount) + s
	}
	return s
}

// RightPad 在右侧填充字符
func RightPad(s string, padStr string, totalLength int) string {
	padCount := totalLength - len(s)
	if padCount > 0 {
		return s + strings.Repeat(padStr, padCount)
	}
	return s
}

// BothPad 在左右两侧填充字符
func BothPad(s string, padStr string, totalLength int) string {
	padCount := totalLength - len(s)
	if padCount > 0 {
		leftPad := padCount / 2
		rightPad := padCount - leftPad
		return strings.Repeat(padStr, leftPad) + s + strings.Repeat(padStr, rightPad)
	}
	return s
}

// ShuffleString 打乱字符串的字符顺序
func ShuffleString(s string) string {
	// 将字符串转换为 rune 切片，支持Unicode字符
	runes := []rune(s)
	n := len(runes)

	// 使用 Fisher-Yates 算法进行打乱
	for i := n - 1; i > 0; i-- {
		// 高效生成随机索引
		j := gr.Intn(i + 1)
		// 交换字符
		runes[i], runes[j] = runes[j], runes[i]
	}

	// 返回打乱后的字符串
	return string(runes)
}
