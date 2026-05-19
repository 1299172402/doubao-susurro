package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// messageResponse API 响应结构
type messageResponse struct {
	DownlinkBody struct {
		PullSingleChainDownlinkBody struct {
			Messages []struct {
				UserType   int    `json:"user_type"`
				MessageID  string `json:"message_id"`
				TTSContent string `json:"tts_content"`
			} `json:"messages"`
		} `json:"pull_singe_chain_downlink_body"`
	} `json:"downlink_body"`
}

// GetLatestMessage 获取最新一条用户消息，返回消息 ID 和文本内容
func GetLatestMessage() (string, string, error) {
	config, err := getConfig()
	if err != nil {
		return "", "", fmt.Errorf("配置加载失败: %w", err)
	}

	// 构建请求体
	payloadBytes, err := json.Marshal(config.Payload)
	if err != nil {
		return "", "", fmt.Errorf("序列化 payload 失败: %w", err)
	}

	// 构建完整 URL（带参数）
	reqURL := config.URL
	if len(config.Params) > 0 {
		params := url.Values{}
		for k, v := range config.Params {
			params.Set(k, v)
		}
		reqURL += "?" + params.Encode()
	}

	// 创建请求
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 headers
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	// 设置 cookies
	for k, v := range config.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("请求失败: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var msgResp messageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 查找用户消息（user_type == 1）
	messages := msgResp.DownlinkBody.PullSingleChainDownlinkBody.Messages
	for _, msg := range messages {
		if msg.UserType == 1 {
			return msg.MessageID, msg.TTSContent, nil
		}
	}

	return "", "", fmt.Errorf("未找到用户消息：%s", string(body))
}
