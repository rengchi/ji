// 参考 https://github.com/AlistGo/alist/blob/main/pkg/utils/json.go

package ji

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WriteJSONToFile 将结构体写入 JSON 文件
// 参数:
//   - dst: 目标文件路径
//   - data: 要写入的数据, 可以是结构体或其他支持 JSON 序列化的类型
//   - indent: 可选参数, 用于设置 JSON 缩进, 如果为 "", 则不缩进
//
// 返回值:
//   - error: 如果写入过程中发生错误，则返回错误信息
func WriteJSONToFile(dst string, data any, indent string) error {
	// 获取文件所在目录路径
	dir := filepath.Dir(dst)

	// 创建目录（如果不存在）
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录 %s 失败：%w", dir, err)
	}

	// 打开目标文件（如果文件已存在则清空）
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s：%w", dst, err)
	}
	defer file.Close()

	// 使用 bufio.Writer 进行缓冲写入，提高写入性能
	writer := bufio.NewWriter(file)

	// 使用 JSON 编码器写入文件
	encoder := json.NewEncoder(writer)
	if indent != "" {
		encoder.SetIndent("", indent)
	}

	// 编码数据并写入缓冲区
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("写入 JSON 数据失败：%w", err)
	}

	// 确保缓冲区的内容被写入文件
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区到文件失败：%w", err)
	}

	return nil
}

// ReadJSONFromFile 从 JSON 文件读取数据到结构体
// 参数:
//   - src: 源文件路径
//   - v: 目标结构体的指针
//
// 返回值:
//   - error: 如果读取过程中发生错误，则返回错误信息
func ReadJSONFromFile(src string, v any) error {
	file, err := os.Open(src)
	if err != nil {
		// 打开文件失败
		fmt.Printf("打开文件：%s 错误：%s\n", src, err.Error())
		return err
	}
	defer file.Close()

	// 使用标准库的 JSON 解码
	decoder := json.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		// JSON 解码错误
		fmt.Printf("json 数据解码 (decode) 错误：%s\n", err.Error())
		return err
	}
	return nil
}
