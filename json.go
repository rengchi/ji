// 参考 https://github.com/AlistGo/alist/blob/main/pkg/utils/json.go

package ji

import (
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
func WriteJSONToFile(dst string, data interface{}, indent string) error {
	var str []byte
	var err error

	// 使用标准库的 JSON 序列化
	if indent != "" {
		// 使用 MarshalIndent 进行格式化 JSON 序列化
		str, err = json.MarshalIndent(data, "", indent)
	} else {
		// 使用 Marshal 进行紧凑的 JSON 序列化
		str, err = json.Marshal(data)
	}
	if err != nil {
		// 错误的 JSON 序列化
		fmt.Printf("json 错误，%s\n", err.Error())
		return err
	}

	// 创建目标目录（如果不存在）
	dir := filepath.Dir(dst) // 使用 filepath.Dir 获取文件的目录路径
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		// 创建目录失败
		fmt.Printf("创建目录失败：%s\n", err.Error())
		return err
	}

	// 将格式化后的 JSON 数据写入文件
	err = os.WriteFile(dst, str, os.ModePerm)
	if err != nil {
		// 写入文件错误
		fmt.Printf("json 写入文件错误：%s； 目录：%s\n", err.Error(), dst)
		return err
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
func ReadJSONFromFile(src string, v interface{}) error {
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
