package ji

import (
	"regexp"
	"sync"
)

// 预编译正则表达式和初始化机制
var (
	alphaNumericRegex *regexp.Regexp // 字母数字正则
	emailRegex        *regexp.Regexp // 邮箱正则
	urlRegex          *regexp.Regexp // 链接正则
	domainRegex       *regexp.Regexp // 域名正则
	ipv4Regex         *regexp.Regexp // ipv4地址正则
	ipv6Regex         *regexp.Regexp // ipv6地址正则
	uuidHexRegex      *regexp.Regexp // uuid正则
	sqlInjectionRegex *regexp.Regexp // sql注入正则
	filenameRegex     *regexp.Regexp // 文件名正则
	once              sync.Once      // 正则初始化
)

// 初始化所有正则表达式
func initRegex() {
	once.Do(func() {
		alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
		emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		urlRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
		domainRegex = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z]{2,})$`)
		ipv4Regex = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
		ipv6Regex = regexp.MustCompile(`^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
		uuidHexRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		sqlInjectionRegex = regexp.MustCompile(`(?i)(union(\s+all)?(\s+select)?|select\s+.*\s+from\s+.*\s+where\s+.*--|insert\s+into\s+.*\s+values|update\s+.*\s+set\s+.*--|delete\s+from\s+.*--|drop\s+table\s+.*|alter\s+table\s+.*|create\s+table\s+.*|exec\s+|execute\s+|call\s+|sleep\s*\()`)
		filenameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,255}$`)
	})
}

// IsAlphaNumeric 检查字符串是否只包含字母和数字
func IsAlphaNumeric(input string) bool {
	return alphaNumericRegex.MatchString(input)
}

// IsValidEmail 检查电子邮件地址是否有效
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidURL 检查链接是否有效
func IsValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

// IsValidDomain 检查域名是否有效
func IsValidDomain(domain string) bool {
	return domainRegex.MatchString(domain)
}

// IsValidIPv4 检查 IPv4 地址是否有效
func IsValidIPv4(ip string) bool {
	return ipv4Regex.MatchString(ip)
}

// IsValidIPv6 检查 IPv6 地址是否有效
func IsValidIPv6(ip string) bool {
	return ipv6Regex.MatchString(ip)
}

// IsValidIP 检查 IP 地址是否有效（支持 IPv4 和 IPv6）
func IsValidIP(ip string) bool {
	return IsValidIPv4(ip) || IsValidIPv6(ip)
}

// IsUUID 检查字符串是否为有效的 UUID
func IsUUID(uuid string) bool {
	return uuidHexRegex.MatchString(uuid)
}

// ContainsSQLInjection 检查字符串是否包含潜在的 SQL 注入代码
func ContainsSQLInjection(input string) bool {
	return sqlInjectionRegex.MatchString(input)
}

// IsValidFilename 检查文件名是否有效
func IsValidFilename(filename string) bool {
	return filenameRegex.MatchString(filename)
}
