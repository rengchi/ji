package ji

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Resp 定义响应结构体，Data字段使用空接口any表示可以接受任意类型数据
type Resp struct {
	Code            int    `json:"code"`                      // 接口返回的状态码
	Status          bool   `json:"status"`                    // 请求处理状态，true表示成功，false表示失败
	Message         string `json:"message,omitempty"`         // 返回的提示信息，若无可省略
	Data            any    `json:"data,omitempty"`            // 返回的实际数据，若为空则不返回该字段
	ExecTime        int64  `json:"exec_time,omitempty"`       // 请求处理时间，单位：毫秒
	IncludeExecTime bool   `json:"-"`                         // 控制是否添加执行时间字段（用于内部控制，不出现在 JSON 输出中）
	HTMLIsEscaped   bool   `json:"html_is_escaped,omitempty"` // 是否转义HTML
}

// WriteResponse 用于写入 HTTP 响应
func WriteResponse(w http.ResponseWriter, statusCode int, resp Resp) {
	// 获取请求开始时间，用于性能监控
	start := time.Now()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 判断是否需要包含执行时间
	if resp.IncludeExecTime {
		// 获取请求处理时间，并添加到响应体
		resp.ExecTime = time.Since(start).Milliseconds()
	}

	// 根据 HTMLIsEscaped 字段判断是否需要转义 HTML
	escapeHTML := !resp.HTMLIsEscaped // 默认为 true，表示需要转义 HTML

	// 判断客户端是否支持 Gzip 压缩
	if strings.Contains(w.Header().Get("Accept-Encoding"), "gzip") {
		// 启用 Gzip 压缩
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// 使用 json.NewEncoder 直接写入响应体
		encoder := json.NewEncoder(gz)
		encoder.SetEscapeHTML(escapeHTML) // 设置 HTML 转义
		if err := encoder.Encode(resp); err != nil {
			// 处理序列化错误
			http.Error(w, "JSON 编码错误", http.StatusInternalServerError)
			return
		}
	} else {
		// 如果不支持 Gzip 压缩，直接写入响应体
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(escapeHTML) // 设置 HTML 转义
		if err := encoder.Encode(resp); err != nil {
			// 处理序列化错误
			http.Error(w, "JSON 编码错误", http.StatusInternalServerError)
			return
		}
	}
}
