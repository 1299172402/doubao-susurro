package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"Doubao-input/internal/config"
)

// curlConfig 解析后的 curl 配置
type curlConfig struct {
	URL     string
	Params  map[string]string
	Headers map[string]string
	Cookies map[string]string
	Payload map[string]interface{}
}

// parseCurl 解析 curl 命令，返回配置
func parseCurl(curlStr string) (*curlConfig, error) {
	config := &curlConfig{
		Params:  make(map[string]string),
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
		Payload: make(map[string]interface{}),
	}

	// 提取 URL（支持单引号和双引号）
	urlRegex := regexp.MustCompile(`curl\s+'([^']+)'`)
	matches := urlRegex.FindStringSubmatch(curlStr)
	if matches == nil {
		urlRegex = regexp.MustCompile(`curl\s+"([^"]+)"`)
		matches = urlRegex.FindStringSubmatch(curlStr)
	}
	if matches == nil {
		return nil, fmt.Errorf("无法提取 URL")
	}
	fullURL := matches[1]

	// 分离 base URL 和 query params
	parsed, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("URL 解析失败: %w", err)
	}
	config.URL = fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
	for k, v := range parsed.Query() {
		if len(v) > 0 {
			config.Params[k] = v[0]
		}
	}

	// 提取 headers（-H 'key: value'）
	headerRegex := regexp.MustCompile(`-H\s+'([^']+)'`)
	for _, match := range headerRegex.FindAllStringSubmatch(curlStr, -1) {
		line := match[1]
		if idx := strings.Index(line, ":"); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			config.Headers[key] = val
		}
	}

	// 提取 cookies（-b 'key=val; key=val'）
	cookieRegex := regexp.MustCompile(`-b\s+'([^']+)'`)
	cookieMatch := cookieRegex.FindStringSubmatch(curlStr)
	if cookieMatch != nil {
		for _, pair := range strings.Split(cookieMatch[1], "; ") {
			if idx := strings.Index(pair, "="); idx != -1 {
				k := pair[:idx]
				v := pair[idx+1:]
				config.Cookies[k] = v
			}
		}
	}

	// 提取 payload（--data-raw '{...}'）
	payloadRegex := regexp.MustCompile(`--data-raw\s+'(\{.*\})'`)
	payloadMatch := payloadRegex.FindStringSubmatch(curlStr)
	if payloadMatch == nil {
		payloadRegex = regexp.MustCompile(`--data-raw\s+"(\{.*\})"`)
		payloadMatch = payloadRegex.FindStringSubmatch(curlStr)
	}
	if payloadMatch != nil {
		json.Unmarshal([]byte(payloadMatch[1]), &config.Payload)
	}

	return config, nil
}

// getConfig 从 config 中解析配置
func getConfig() (*curlConfig, error) {
	curlStr := config.GetConfig().Session

	config, err := parseCurl(curlStr)
	if err != nil {
		return nil, err
	}

	// 修改 payload 中的 direction 和 anchor_index
	if uplink, ok := config.Payload["uplink_body"].(map[string]interface{}); ok {
		if body, ok := uplink["pull_singe_chain_uplink_body"].(map[string]interface{}); ok {
			body["direction"] = 0
			body["anchor_index"] = 0
		}
	}

	return config, nil
}
