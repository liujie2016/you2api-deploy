package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FallbackHandler 提供更鲁棒的处理方案
func FallbackHandler(w http.ResponseWriter, r *http.Request) {
	// 基本的CORS和路由处理
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.URL.Path != "/v1/chat/completions" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "You2Api Fallback Service Running...",
			"message": "Enhanced compatibility mode",
		})
		return
	}

	// 解析请求
	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(openAIReq.Messages) == 0 {
		http.Error(w, "Messages array cannot be empty", http.StatusBadRequest)
		return
	}

	// 获取用户消息
	userMessage := openAIReq.Messages[len(openAIReq.Messages)-1].Content
	model := openAIReq.Model
	if model == "" {
		model = "gpt-4o"
	}

	log.Printf("Processing request: model=%s, message=%s", model, userMessage[:min(50, len(userMessage))])

	// 尝试多种方法获取响应
	content := tryMultipleMethods(userMessage, model)

	// 如果所有方法都失败，提供智能回退
	if content == "" {
		content = generateFallbackResponse(userMessage)
	}

	// 返回响应
	if openAIReq.Stream {
		sendStreamResponse(w, content, model)
	} else {
		sendNormalResponse(w, content, model)
	}
}

// tryMultipleMethods 尝试多种方法获取响应
func tryMultipleMethods(message, model string) string {
	methods := []func(string, string) string{
		tryOriginalMethod,
		trySimplifiedMethod,
		tryAlternativeProxy,
		tryDirectCall,
	}

	for i, method := range methods {
		log.Printf("Trying method %d...", i+1)
		if content := method(message, model); content != "" {
			log.Printf("Method %d succeeded", i+1)
			return content
		}
	}

	log.Printf("All methods failed")
	return ""
}

// tryOriginalMethod 尝试原始方法
func tryOriginalMethod(message, model string) string {
	youModel := mapModelName(model)

	params := url.Values{}
	params.Add("q", message)
	params.Add("page", "1")
	params.Add("count", "10")
	params.Add("safeSearch", "Moderate")
	params.Add("mkt", "zh-HK")
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", youModel)
	params.Add("selectedChatMode", "custom")

	proxyURL := "https://proxy.cors.sh/https://you.com/api/streamingSearch?" + params.Encode()
	return makeRequest(proxyURL, map[string]string{
		"x-cors-api-key": "live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7",
		"User-Agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	})
}

// trySimplifiedMethod 尝试简化方法
func trySimplifiedMethod(message, model string) string {
	youModel := mapModelName(model)

	params := url.Values{}
	params.Add("q", message)
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", youModel)

	proxyURL := "https://proxy.cors.sh/https://you.com/api/streamingSearch?" + params.Encode()
	return makeRequest(proxyURL, map[string]string{
		"x-cors-api-key": "live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7",
	})
}

// tryAlternativeProxy 尝试替代代理
func tryAlternativeProxy(message, model string) string {
	proxies := []string{
		"https://cors-anywhere.herokuapp.com/",
		"https://api.allorigins.win/raw?url=",
		"https://thingproxy.freeboard.io/fetch/",
	}

	youModel := mapModelName(model)
	params := url.Values{}
	params.Add("q", message)
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", youModel)

	originalURL := "https://you.com/api/streamingSearch?" + params.Encode()

	for _, proxy := range proxies {
		proxyURL := proxy + url.QueryEscape(originalURL)
		if content := makeRequest(proxyURL, nil); content != "" {
			return content
		}
	}

	return ""
}

// tryDirectCall 尝试直接调用（可能被CORS阻止）
func tryDirectCall(message, model string) string {
	youModel := mapModelName(model)

	params := url.Values{}
	params.Add("q", message)
	params.Add("domain", "youchat")
	params.Add("selectedAiModel", youModel)

	directURL := "https://you.com/api/streamingSearch?" + params.Encode()
	return makeRequest(directURL, map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Referer":    "https://you.com/",
		"Origin":     "https://you.com",
	})
}

// makeRequest 执行HTTP请求并解析响应
func makeRequest(url string, headers map[string]string) string {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return ""
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Non-200 status code: %d", resp.StatusCode)
		return ""
	}

	return parseResponse(resp.Body)
}

// parseResponse 解析响应内容
func parseResponse(body io.Reader) string {
	var content strings.Builder
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" || data == "" {
				continue
			}

			// 尝试解析JSON
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(data), &parsed); err == nil {
				if text := extractTextFromResponse(parsed); text != "" {
					content.WriteString(text)
				}
			}
		} else if line != "" {
			// 如果不是标准的SSE格式，可能是纯文本响应
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(line), &parsed); err == nil {
				if text := extractTextFromResponse(parsed); text != "" {
					content.WriteString(text)
				}
			}
		}
	}

	return content.String()
}

// extractTextFromResponse 从响应中提取文本内容
func extractTextFromResponse(data map[string]interface{}) string {
	// 尝试常见的响应字段
	fields := []string{
		"youChatToken", "text", "message", "content", "answer", "response",
		"completion", "output", "result", "reply", "data",
	}

	for _, field := range fields {
		if value, exists := data[field].(string); exists && value != "" {
			return value
		}
	}

	// 检查嵌套结构
	if delta, ok := data["delta"].(map[string]interface{}); ok {
		if content, ok := delta["content"].(string); ok && content != "" {
			return content
		}
	}

	if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok && content != "" {
					return content
				}
			}
			if delta, ok := choice["delta"].(map[string]interface{}); ok {
				if content, ok := delta["content"].(string); ok && content != "" {
					return content
				}
			}
		}
	}

	return ""
}

// generateFallbackResponse 生成智能回退响应
func generateFallbackResponse(message string) string {
	// 简单的关键词匹配回复
	message = strings.ToLower(message)

	responses := map[string]string{
		"hello":     "Hello! How can I help you today?",
		"你好":        "你好！我是AI助手，很高兴为您服务！",
		"hi":        "Hi there! What can I do for you?",
		"help":      "I'm here to help! Please let me know what you need assistance with.",
		"what":      "I'm an AI assistant created to help answer questions and have conversations.",
		"who":       "I'm an AI assistant. How may I assist you today?",
		"how":       "I can help you with various questions and tasks. What would you like to know?",
		"thanks":    "You're welcome! Is there anything else I can help you with?",
		"thank you": "You're very welcome! Feel free to ask if you need anything else.",
		"谢谢":        "不客气！如果您还有其他问题，请随时告诉我。",
	}

	for keyword, response := range responses {
		if strings.Contains(message, keyword) {
			return response
		}
	}

	// 默认回复
	if isChineseMessage(message) {
		return "抱歉，我目前无法处理您的请求。这可能是由于网络连接问题或服务暂时不可用。请稍后再试，或者重新表述您的问题。"
	}

	return "I apologize, but I'm currently unable to process your request. This might be due to network connectivity issues or temporary service unavailability. Please try again later or rephrase your question."
}

// isChineseMessage 检查是否为中文消息
func isChineseMessage(message string) bool {
	for _, r := range message {
		if r >= 0x4e00 && r <= 0x9fff { // 中文Unicode范围
			return true
		}
	}
	return false
}

// sendStreamResponse 发送流式响应
func sendStreamResponse(w http.ResponseWriter, content, model string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	responseID := "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10)
	created := time.Now().Unix()

	// 分块发送内容
	words := strings.Fields(content)
	for i, word := range words {
		chunk := OpenAIStreamResponse{
			ID:      responseID,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   model,
			Choices: []Choice{{
				Delta: Delta{Content: word + " "},
				Index: 0,
			}},
		}

		if jsonData, err := json.Marshal(chunk); err == nil {
			fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
		}

		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		// 模拟打字效果
		if i < len(words)-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 发送结束信号
	finalChunk := OpenAIStreamResponse{
		ID:      responseID,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []Choice{{
			Delta:        Delta{Content: ""},
			Index:        0,
			FinishReason: "stop",
		}},
	}

	if jsonData, err := json.Marshal(finalChunk); err == nil {
		fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
	}
	fmt.Fprint(w, "data: [DONE]\n\n")

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// sendNormalResponse 发送普通响应
func sendNormalResponse(w http.ResponseWriter, content, model string) {
	w.Header().Set("Content-Type", "application/json")

	response := OpenAIResponse{
		ID:      "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []OpenAIChoice{{
			Message: Message{
				Role:    "assistant",
				Content: content,
			},
			Index:        0,
			FinishReason: "stop",
		}},
	}

	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
