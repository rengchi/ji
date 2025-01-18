package ji

import (
	"errors"
	"fmt"
	"strings"
)

// Miss 创建一个新的错误对象，包含拼接后的错误消息
// 参数:
//   - str: 要拼接的多个字符串
//
// 返回值:
//   - error: 新创建的错误对象，包含拼接后的错误信息
func Miss(str ...string) error {
	// 如果没有传递任何参数，则返回一个空的错误
	if len(str) == 0 {
		return errors.New("未提供错误信息")
	}

	// 使用 strings.Builder 拼接多个字符串，提高性能
	var builder strings.Builder
	for i, s := range str {
		if i > 0 {
			// 非第一个字符串前添加空格
			builder.WriteString(" ")
		}
		builder.WriteString(s)
	}

	// 返回拼接后的错误信息
	return fmt.Errorf("错误: %s", builder.String())
}
