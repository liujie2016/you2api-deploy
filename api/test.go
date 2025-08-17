package handler

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TestHandler - 简化的测试处理程序，用于调试You.com API
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 简单的测试请求
	testMessage := "Hello, how are you?"
	if r.Method == "POST" {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			if messages, ok := req["messages"].([]interface{}); ok && len(messages) > 0 {
				if msg, ok := messages[len(messages)-1].(map[string]interface{}); ok {
					if content, ok := msg["content"].(string); ok {
						testMessage = content
					}
				}
			}
		}
	}

	log.Printf("Testing with message: %s", testMessage)

	// 方法1: 直接调用You.com API（不使用代理）
	result1 := testDirectYouAPI(testMessage)

	// 方法2: 使用CORS代理
	result2 := testWithCORSProxy(testMessage)

	// 方法3: 测试不同的请求参数
	result3 := testWithDifferentParams(testMessage)

	response := map[string]interface{}{
		"test_message":     testMessage,
		"direct_api":       result1,
		"cors_proxy":       result2,
		"different_params": result3,
		"timestamp":        time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

func testDirectYouAPI(message string) map[string]interface{} {
	log.Printf("Testing direct You.com API call...")

	// 构建请求URL
	baseURL := "https://you.com/api/streamingSearch"
	params := url.Values{}
	params.Add("q", message)
	params.Add("page", "1")
	params.Add("count", "10")
	params.Add("safeSearch", "Moderate")
	params.Add("mkt", "zh-HK")
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", "gpt_4o")
	params.Add("selectedChatMode", "custom")

	fullURL := baseURL + "?" + params.Encode()
	log.Printf("Direct API URL: %s", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://you.com/")
	req.Header.Set("Origin", "https://you.com")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()

	log.Printf("Direct API response status: %d", resp.StatusCode)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	bodyStr := string(body)
	log.Printf("Direct API response body (first 500 chars): %s", truncate(bodyStr, 500))

	return map[string]interface{}{
		"status_code":  resp.StatusCode,
		"body_length":  len(bodyStr),
		"body_preview": truncate(bodyStr, 200),
		"headers":      formatHeaders(resp.Header),
	}
}

func testWithCORSProxy(message string) map[string]interface{} {
	log.Printf("Testing with CORS proxy...")

	// 使用CORS代理
	originalURL := "https://you.com/api/streamingSearch"
	params := url.Values{}
	params.Add("q", message)
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", "gpt_4o")

	proxyURL := "https://proxy.cors.sh/" + originalURL + "?" + params.Encode()
	log.Printf("CORS Proxy URL: %s", proxyURL)

	req, err := http.NewRequest("GET", proxyURL, nil)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	req.Header.Set("x-cors-api-key", "live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer resp.Body.Close()

	log.Printf("CORS Proxy response status: %d", resp.StatusCode)

	// 解析流式响应
	content := parseStreamResponse(resp.Body)

	return map[string]interface{}{
		"status_code":    resp.StatusCode,
		"content":        content,
		"content_length": len(content),
	}
}

func testWithDifferentParams(message string) map[string]interface{} {
	log.Printf("Testing with different parameters...")

	// 尝试不同的参数组合
	testCases := []map[string]string{
		{
			"selectedAiModel":  "deepseek_v3",
			"selectedChatMode": "default",
		},
		{
			"selectedAiModel":  "claude_3_5_sonnet",
			"selectedChatMode": "custom",
		},
		{
			"selectedAiModel":  "gemini_1_5_pro",
			"selectedChatMode": "creative",
		},
	}

	results := make([]map[string]interface{}, 0)

	for i, testCase := range testCases {
		log.Printf("Testing case %d: %+v", i+1, testCase)

		params := url.Values{}
		params.Add("q", message)
		params.Add("domain", "youchat")
		for key, value := range testCase {
			params.Add(key, value)
		}

		proxyURL := "https://proxy.cors.sh/https://you.com/api/streamingSearch?" + params.Encode()

		req, err := http.NewRequest("GET", proxyURL, nil)
		if err != nil {
			results = append(results, map[string]interface{}{
				"case":  testCase,
				"error": err.Error(),
			})
			continue
		}

		req.Header.Set("x-cors-api-key", "live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			results = append(results, map[string]interface{}{
				"case":  testCase,
				"error": err.Error(),
			})
			continue
		}

		content := parseStreamResponse(resp.Body)
		resp.Body.Close()

		results = append(results, map[string]interface{}{
			"case":           testCase,
			"status_code":    resp.StatusCode,
			"content":        content,
			"content_length": len(content),
		})
	}

	return map[string]interface{}{
		"test_cases": results,
	}
}

func parseStreamResponse(body io.Reader) string {
	var content strings.Builder
	scanner := bufio.NewScanner(body)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		log.Printf("Stream line %d: %s", lineCount, truncate(line, 200))

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			if data != "" && data != "{}" {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(data), &parsed); err == nil {
					log.Printf("Parsed JSON: %+v", parsed)

					// 尝试多种字段名
					fields := []string{"youChatToken", "text", "message", "content", "answer"}
					for _, field := range fields {
						if value, exists := parsed[field].(string); exists && value != "" {
							content.WriteString(value)
							log.Printf("Found content in field '%s': %s", field, truncate(value, 100))
							break
						}
					}

					// 检查嵌套对象
					if delta, exists := parsed["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content.WriteString(deltaContent)
							log.Printf("Found content in delta.content: %s", truncate(deltaContent, 100))
						}
					}
				} else {
					log.Printf("Failed to parse JSON: %v", err)
					log.Printf("Raw data: %s", truncate(data, 200))
				}
			}
		}
	}

	result := content.String()
	log.Printf("Final parsed content: %s", truncate(result, 200))
	return result
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func formatHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
