package ji

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// 文件大小单位
const (
	Byte = 1       // 1 Byte
	KB   = 1 << 10 // 1 KB = 1024 Bytes
	MB   = 1 << 20 // 1 MB = 1024 KB
	GB   = 1 << 30 // 1 GB = 1024 MB
	TB   = 1 << 40 // 1 TB = 1024 GB
	PB   = 1 << 50 // 1 PB = 1024 TB
	EB   = 1 << 60 // 1 EB = 1024 PB
	ZB   = 1 << 70 // 1 ZB = 1024 EB
	YB   = 1 << 80 // 1 YB = 1024 ZB
)

// FileExist 检查文件是否存在且不是目录
// 参数:
//   - fullPathWithFilename: 文件的完整路径及文件名
//
// 返回值:
//   - fs.FileInfo: 文件信息。如果文件不存在或是目录，该值无效。
//   - bool: 如果文件存在且不是目录，则返回 true；否则返回 false。
func FileExist(fullPathWithFilename string) (fs.FileInfo, bool) {
	// 获取文件信息
	fileInfo, err := os.Stat(fullPathWithFilename)

	// 文件不存在
	if os.IsNotExist(err) {
		return nil, false
	}

	// 检查是否为目录
	if fileInfo.IsDir() {
		return fileInfo, false
	}

	return fileInfo, true
}

// CreateFile 创建文件，如果文件路径不存在则自动创建
// 参数: fullPathWithFilename - 文件的完整路径及文件名
// 返回值: 成功返回文件指针 *os.File，错误时返回 nil 和 error
func CreateFile(fullPathWithFilename string) (*os.File, error) {
	err := CreateDir(filepath.Dir(fullPathWithFilename))
	if err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	file, err := os.Create(fullPathWithFilename)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	return file, nil
}

// DirExist 检查目录是否存在
// 参数: dir - 目录路径
// 返回值: 如果目录存在且是目录，则返回 true；否则返回 false
func DirExist(dir string) bool {
	ff, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return ff.IsDir()
}

// CreateDir 创建目录及其父目录，如果目录不存在则自动创建
// 参数: dir - 目录路径
// 返回值: 创建成功返回 nil，错误时返回 error
func CreateDir(dir string) error {
	// 如果目录不存在，创建目录
	if !DirExist(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}
	return nil
}

// FileExtension 返回文件的扩展名
func FileExtension(filePath string) string {
	return filepath.Ext(filePath)
}

// FileSize 返回文件的大小，以字节为单位
func FileSize(filePath string) (int64, error) {
	fileInfo, exist := FileExist(filePath)
	if !exist {
		return 0, fmt.Errorf("文件：%s不存在", filePath)
	}
	return fileInfo.Size(), nil
}

// FileMimeType 返回文件的 MIME 类型
// 参数:
//   - filePath: 文件的完整路径及文件名
//
// 返回值:
//   - string: 文件的 MIME 类型
//   - error: 如果发生错误，则返回错误信息
func FileMimeType(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建一个缓冲区来读取文件的前512个字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != os.ErrClosed {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 使用 http.DetectContentType 来获取 MIME 类型
	return http.DetectContentType(buffer[:n]), nil
}

// RemoveFiles 删除多个路径的文件
// 参数: filePaths - 一个或多个文件路径
// 返回值: 成功返回 nil，错误时返回 error
func RemoveFiles(filePaths []string) error {
	var errors []string
	for _, filePath := range filePaths {
		if err := os.Remove(filePath); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("删除部分文件失败: %s", strings.Join(errors, "; "))
	}
	return nil
}

// RemoveDirs 删除多个目录
// 参数: dirPaths - 目录路径的数组
// 返回值: 成功返回 nil，错误时返回 error
func RemoveDirs(dirPaths []string) error {
	var errors []string

	// 遍历每个目录路径并尝试删除
	for _, dirPath := range dirPaths {
		err := os.Remove(dirPath)
		if err != nil {
			// 记录错误，但继续删除其他目录
			errors = append(errors, fmt.Sprintf("删除目录 %s 失败: %v", dirPath, err))
		}
	}

	// 如果有错误，返回所有错误信息
	if len(errors) > 0 {
		return fmt.Errorf("删除部分目录失败: %s", strings.Join(errors, "; "))
	}
	return nil
}

// FileHash 计算文件的 SHA-256 哈希值
func FileHash(file io.Reader) (string, error) {
	// 使用 SHA-256 计算哈希值
	hash := sha256.New()

	// 创建一个缓冲读取器，减少系统调用次数
	bufReader := bufio.NewReader(file)

	// 使用更大的缓冲区进行读取，例如 128KB 256KB
	buffer := make([]byte, 128*1024) // 缓冲区
	for {
		n, err := bufReader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// 更新哈希值
		hash.Write(buffer[:n])
	}

	// 返回哈希值的十六进制表示
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SafeFileName 处理文件名，转义掉 URL 中需要转义的字符
func SafeFileName(fileName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(fileName)

	// 获取文件名（去除扩展名）
	baseName := fileName[:len(fileName)-len(ext)]

	// 转义文件名（去除文件扩展名部分）
	safeBaseName := url.QueryEscape(baseName)

	// 返回转义后的文件名和扩展名
	return fmt.Sprintf("%s%s", safeBaseName, ext)
}
