package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type YouChatResponse struct {
	YouChatToken string `json:"youChatToken"`
}

type OpenAIStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Delta        Delta  `json:"delta"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

var modelMap = map[string]string{
	"deepseek-reasoner":  "deepseek_r1",
	"deepseek-chat":      "deepseek_v3",
	"o3-mini-high":       "openai_o3_mini_high",
	"o3-mini-medium":     "openai_o3_mini_medium",
	"o1":                 "openai_o1",
	"o1-mini":            "openai_o1_mini",
	"o1-preview":         "openai_o1_preview",
	"gpt-4o":             "gpt_4o",
	"gpt-4o-mini":        "gpt_4o_mini",
	"gpt-4-turbo":        "gpt_4_turbo",
	"gpt-3.5-turbo":      "gpt_3_5",
	"claude-3-opus":      "claude_3_opus",
	"claude-3-sonnet":    "claude_3_sonnet",
	"claude-3.5-sonnet":  "claude_3_5_sonnet",
	"claude-3.5-haiku":   "claude_3_5_haiku",
	"gemini-1.5-pro":     "gemini_1_5_pro",
	"gemini-1.5-flash":   "gemini_1_5_flash",
	"llama-3.2-90b":      "llama3_2_90b",
	"llama-3.1-405b":     "llama3_1_405b",
	"mistral-large-2":    "mistral_large_2",
	"qwen-2.5-72b":       "qwen2p5_72b",
	"qwen-2.5-coder-32b": "qwen2p5_coder_32b",
	"command-r-plus":     "command_r_plus",
}

func getReverseModelMap() map[string]string {
	reverse := make(map[string]string, len(modelMap))
	for k, v := range modelMap {
		reverse[v] = k
	}
	return reverse
}

func mapModelName(openAIModel string) string {
	if mappedModel, exists := modelMap[openAIModel]; exists {
		return mappedModel
	}
	return "deepseek_v3"
}

func reverseMapModelName(youModel string) string {
	reverseMap := getReverseModelMap()
	if mappedModel, exists := reverseMap[youModel]; exists {
		return mappedModel
	}
	return "deepseek-chat"
}

var originalModel string

func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle test endpoint first
	if r.URL.Path == "/test" || r.URL.Path == "/test/" {
		TestHandler(w, r)
		return
	}

	// 检查是否应该使用备用处理器
	if os.Getenv("USE_FALLBACK") == "true" || os.Getenv("FALLBACK_MODE") == "true" {
		FallbackHandler(w, r)
		return
	}

	if r.URL.Path != "/v1/chat/completions" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "You2Api Service Running...",
			"message": "MoLoveSze...",
		})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("Missing or invalid authorization header from IP: %s", r.RemoteAddr)
		http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	dsToken := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Request authenticated, token length: %d", len(dsToken))
	_ = dsToken // Token for future use

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(openAIReq.Messages) == 0 {
		log.Printf("Empty messages array received")
		http.Error(w, "Messages array cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Processing request: model=%s, messages=%d, stream=%v",
		openAIReq.Model, len(openAIReq.Messages), openAIReq.Stream)

	originalModel = openAIReq.Model
	lastMessage := openAIReq.Messages[len(openAIReq.Messages)-1].Content
	var chatHistory []map[string]interface{}
	for _, msg := range openAIReq.Messages {
		chatMsg := map[string]interface{}{
			"question": msg.Content,
			"answer":   "",
		}
		if msg.Role == "assistant" {
			chatMsg["question"] = ""
			chatMsg["answer"] = msg.Content
		}
		chatHistory = append(chatHistory, chatMsg)
	}

	chatHistoryJSON, _ := json.Marshal(chatHistory)

	// Construct the original URL
	originalURL := "https://you.com/api/streamingSearch"
	q := url.Values{}
	q.Add("q", lastMessage)
	q.Add("page", "1")
	q.Add("count", "10")
	q.Add("safeSearch", "Moderate")
	q.Add("mkt", "zh-HK")
	q.Add("enable_worklow_generation_ux", "true")
	q.Add("domain", "youchat")
	q.Add("use_personalization_extraction", "true")
	q.Add("pastChatLength", fmt.Sprintf("%d", len(chatHistory)-1))
	q.Add("selectedChatMode", "custom")
	q.Add("selectedAiModel", mapModelName(openAIReq.Model))
	q.Add("enable_agent_clarification_questions", "true")
	q.Add("use_nested_youchat_updates", "true")
	q.Add("chat", string(chatHistoryJSON))

	// Use a proxy to bypass Cloudflare
	proxyURL := "https://proxy.cors.sh/" + originalURL + "?" + q.Encode()

	youReq, err := http.NewRequest("GET", proxyURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	youReq.Header = http.Header{
		"x-cors-api-key":     {"live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"Accept":             {"text/event-stream"},
		"Accept-Language":    {"en-US,en;q=0.9"},
		"Accept-Encoding":    {"gzip, deflate, br"},
		"Referer":            {"https://you.com/"},
		"Origin":             {"https://you.com"},
		"DNT":                {"1"},
		"Connection":         {"keep-alive"},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-origin"},
		"sec-ch-ua":          {"\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
	}

	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	resp, err := client.Do(youReq)
	if err != nil {
		log.Printf("Request error: %v", err)
		http.Error(w, "Failed to send request to You.com API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("You.com response status: %d", resp.StatusCode)

	// Read a small part of the response body for debugging
	bodyPreview := make([]byte, 200)
	n, _ := resp.Body.Read(bodyPreview)
	log.Printf("Response body preview (first %d bytes): %s", n, string(bodyPreview[:n]))

	// Reset body reader
	resp.Body.Close()
	resp, err = client.Do(youReq)
	if err != nil {
		log.Printf("Failed to re-request: %v", err)
		http.Error(w, "Failed to re-request You.com API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("You.com API returned non-200 status: %d, trying fallback methods...", resp.StatusCode)
		// Try fallback method instead of failing immediately
		fallbackContent := tryMultipleMethods(lastMessage, originalModel)
		if fallbackContent == "" {
			fallbackContent = generateFallbackResponse(lastMessage)
			log.Printf("Using generated fallback response")
		} else {
			log.Printf("Successfully got content from fallback method, length: %d", len(fallbackContent))
		}

		if openAIReq.Stream {
			sendStreamResponse(w, fallbackContent, originalModel)
		} else {
			sendNormalResponse(w, fallbackContent, originalModel)
		}
		return
	}

	if openAIReq.Stream {
		content := handleStreamResponse(w, resp, originalModel)
		// If primary method returns empty content, try fallback
		if content == "" {
			log.Printf("Primary stream method returned empty content, trying fallback...")
			fallbackContent := tryMultipleMethods(lastMessage, originalModel)
			if fallbackContent == "" {
				fallbackContent = generateFallbackResponse(lastMessage)
			}
			sendStreamResponse(w, fallbackContent, originalModel)
		}
	} else {
		content := handleNonStreamResponse(w, resp, originalModel)
		// If primary method returns empty content, try fallback
		if content == "" {
			log.Printf("Primary non-stream method returned empty content, trying fallback...")
			fallbackContent := tryMultipleMethods(lastMessage, originalModel)
			if fallbackContent == "" {
				fallbackContent = generateFallbackResponse(lastMessage)
			}
			sendNormalResponse(w, fallbackContent, originalModel)
		}
	}
}

func handleStreamResponse(w http.ResponseWriter, resp *http.Response, model string) string {
	var totalContent strings.Builder
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	responseID := "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10)
	created := time.Now().Unix()
	isDebug := os.Getenv("DEBUG") == "true"

	for scanner.Scan() {
		line := scanner.Text()
		if isDebug {
			log.Printf("Raw line: %s", line)
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if isDebug {
				log.Printf("Data content: %s", data)
			}

			if data == "[DONE]" {
				// Send final chunk
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
				break
			} else if data != "" && data != "{}" {
				var youResp map[string]interface{}
				if err := json.Unmarshal([]byte(data), &youResp); err == nil {
					if isDebug {
						log.Printf("Parsed response: %+v", youResp)
					}

					// Try multiple possible field names for the content
					var content string
					if youChatToken, exists := youResp["youChatToken"].(string); exists && youChatToken != "" {
						content = youChatToken
					} else if text, exists := youResp["text"].(string); exists && text != "" {
						content = text
					} else if message, exists := youResp["message"].(string); exists && message != "" {
						content = message
					} else if delta, exists := youResp["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content = deltaContent
						}
					}

					if content != "" {
						totalContent.WriteString(content)
						chunk := OpenAIStreamResponse{
							ID:      responseID,
							Object:  "chat.completion.chunk",
							Created: created,
							Model:   model,
							Choices: []Choice{{
								Delta: Delta{Content: content},
								Index: 0,
							}},
						}
						if jsonData, err := json.Marshal(chunk); err == nil {
							fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
						}
					} else if isDebug {
						log.Printf("No content found in response: %+v", youResp)
					}
				} else if isDebug {
					log.Printf("Failed to parse JSON: %v, data: %s", err, data)
				}
			}
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	result := totalContent.String()
	log.Printf("handleStreamResponse returning content length: %d", len(result))
	return result
}

func handleNonStreamResponse(w http.ResponseWriter, resp *http.Response, model string) string {
	w.Header().Set("Content-Type", "application/json")

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	isDebug := os.Getenv("DEBUG") == "true"
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if isDebug {
			log.Printf("Non-stream line %d: %s", lineCount, line)
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if isDebug {
				log.Printf("Non-stream data: %s", data)
			}

			if data != "[DONE]" && data != "" && data != "{}" {
				var youResp map[string]interface{}
				if err := json.Unmarshal([]byte(data), &youResp); err == nil {
					if isDebug {
						log.Printf("Non-stream parsed: %+v", youResp)
					}

					// Log all available fields for debugging
					log.Printf("Available fields in response: %v", getMapKeys(youResp))

					// Try multiple possible field names for the content
					var content string
					if youChatToken, exists := youResp["youChatToken"].(string); exists && youChatToken != "" {
						content = youChatToken
					} else if text, exists := youResp["text"].(string); exists && text != "" {
						content = text
					} else if message, exists := youResp["message"].(string); exists && message != "" {
						content = message
					} else if delta, exists := youResp["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content = deltaContent
						}
					}

					if content != "" {
						fullResponse.WriteString(content)
					} else if isDebug {
						log.Printf("No content found in non-stream response: %+v", youResp)
					}
				} else if isDebug {
					log.Printf("Failed to parse non-stream JSON: %v, data: %s", err, data)
				}
			}
		}
	}

	finalContent := fullResponse.String()
	if isDebug {
		log.Printf("Final response content length: %d, content: %s", len(finalContent), finalContent)
	}

	// If no content was found, provide a fallback message
	if finalContent == "" {
		finalContent = "I apologize, but I couldn't process your request at this time. Please try again."
		if isDebug {
			log.Printf("Using fallback content")
		}
	}

	openAIResp := OpenAIResponse{
		ID:      "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []OpenAIChoice{{
			Message: Message{
				Role:    "assistant",
				Content: finalContent,
			},
			Index:        0,
			FinishReason: "stop",
		}},
	}

	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return ""
	}

	log.Printf("handleNonStreamResponse returning content length: %d", len(finalContent))
	return finalContent
}

// TestHandler - 简化的测试处理程序，用于调试You.com API
func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("TestHandler called with path: %s, method: %s", r.URL.Path, r.Method)

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

	response := map[string]interface{}{
		"test_message": testMessage,
		"direct_api":   result1,
		"cors_proxy":   result2,
		"timestamp":    time.Now().Unix(),
		"debug_info": map[string]interface{}{
			"path":   r.URL.Path,
			"method": r.Method,
			"headers": map[string]string{
				"content-type": r.Header.Get("Content-Type"),
			},
		},
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
	bodyBytes, err := bufio.NewReader(resp.Body).ReadBytes('\n')
	if err != nil && len(bodyBytes) == 0 {
		return map[string]interface{}{"error": err.Error()}
	}

	bodyStr := string(bodyBytes)
	log.Printf("Direct API response body (first 500 chars): %s", truncateString(bodyStr, 500))

	return map[string]interface{}{
		"status_code":  resp.StatusCode,
		"body_length":  len(bodyStr),
		"body_preview": truncateString(bodyStr, 200),
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
	content := parseTestStreamResponse(resp.Body)

	return map[string]interface{}{
		"status_code":    resp.StatusCode,
		"content":        content,
		"content_length": len(content),
	}
}

func parseTestStreamResponse(resp *http.Response) string {
	var content strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		log.Printf("Stream line %d: %s", lineCount, truncateString(line, 200))

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
							log.Printf("Found content in field '%s': %s", field, truncateString(value, 100))
							break
						}
					}

					// 检查嵌套对象
					if delta, exists := parsed["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content.WriteString(deltaContent)
							log.Printf("Found content in delta.content: %s", truncateString(deltaContent, 100))
						}
					}
				} else {
					log.Printf("Failed to parse JSON: %v", err)
					log.Printf("Raw data: %s", truncateString(data, 200))
				}
			}
		}
	}

	result := content.String()
	log.Printf("Final parsed content: %s", truncateString(result, 200))
	return result
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Enhanced error logging
func logError(context string, err error) {
	if err != nil {
		log.Printf("ERROR [%s]: %v", context, err)
	}
}

// Enhanced info logging
func logInfo(context string, message string, args ...interface{}) {
	log.Printf("INFO [%s]: %s", context, fmt.Sprintf(message, args...))
}
